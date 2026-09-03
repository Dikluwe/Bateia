/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/runstore/file_store.md
 * @layer L3
 * @updated 2026-09-01
 */

package runstore

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type RunStoreLimits struct {
	MaxArtifactBytes int64
	MaxRunBytes      int64
	MaxTotalBytes    int64
}

type runAlias domain.RecommendationRun

type diskRecord struct {
	SchemaVersion string                   `json:"schema_version"`
	Run           runAlias                 `json:"run"`
	Idempotency   domain.IdempotencyRecord `json:"idempotency"`
	Checksum      string                   `json:"checksum"`
}

type FileStore struct {
	mu      sync.Mutex
	root    string
	records string
	blobs   string
	secret  []byte
	limits  RunStoreLimits
	runs    map[domain.RunID]domain.RecommendationRun
	idem    map[string]domain.IdempotencyRecord
	closed  bool
}

var (
	storeRunIDPattern = regexp.MustCompile(`^run_[0-9a-f]{32}$`)
	digestPattern     = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

func NewFileStore(root string, cursorSecret []byte, limits RunStoreLimits) (*FileStore, error) {
	if root == "" || len(cursorSecret) != 32 || limits.MaxArtifactBytes <= 0 || limits.MaxRunBytes < limits.MaxArtifactBytes || limits.MaxTotalBytes < limits.MaxRunBytes {
		return nil, fmt.Errorf("invalid run store configuration")
	}
	if info, err := os.Lstat(root); err == nil {
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return nil, domain.ErrIntegrity
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, domain.ErrRunStoreUnavailable
	} else if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, domain.ErrRunStoreUnavailable
	}
	if err := os.Chmod(root, 0o700); err != nil {
		return nil, domain.ErrRunStoreUnavailable
	}
	store := &FileStore{
		root: filepath.Clean(root), records: filepath.Join(root, "records"), blobs: filepath.Join(root, "blobs"),
		secret: append([]byte(nil), cursorSecret...), limits: limits,
		runs: map[domain.RunID]domain.RecommendationRun{}, idem: map[string]domain.IdempotencyRecord{},
	}
	for _, directory := range []string{store.records, store.blobs} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return nil, domain.ErrRunStoreUnavailable
		}
		if err := os.Chmod(directory, 0o700); err != nil {
			return nil, domain.ErrRunStoreUnavailable
		}
	}
	if err := store.auditTree(); err != nil {
		return nil, err
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *FileStore) auditTree() error {
	return filepath.WalkDir(s.root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return domain.ErrIntegrity
		}
		info, err := os.Lstat(path)
		if err != nil {
			return domain.ErrIntegrity
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return domain.ErrIntegrity
		}
		if entry.IsDir() {
			if info.Mode().Perm()&0o077 != 0 {
				if err := os.Chmod(path, 0o700); err != nil {
					return domain.ErrRunStoreUnavailable
				}
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return domain.ErrIntegrity
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Nlink != 1 {
			return domain.ErrIntegrity
		}
		if info.Mode().Perm()&0o077 != 0 {
			if err := os.Chmod(path, 0o600); err != nil {
				return domain.ErrRunStoreUnavailable
			}
		}
		return nil
	})
}

func (s *FileStore) load() error {
	entries, err := os.ReadDir(s.records)
	if err != nil {
		return domain.ErrRunStoreUnavailable
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			return domain.ErrIntegrity
		}
		encoded, err := os.ReadFile(filepath.Join(s.records, entry.Name()))
		if err != nil {
			return domain.ErrIntegrity
		}
		run, record, err := decodeDiskRecord(encoded)
		if err != nil || !validRunID(run.RunID) || !digestPattern.MatchString(record.KeyDigest) || !digestPattern.MatchString(record.RequestHash) || record.RunID != run.RunID || record.RequestHash != run.RequestHash {
			return domain.ErrIntegrity
		}
		if expected := recordFileName(run.RunID); entry.Name() != expected {
			return domain.ErrIntegrity
		}
		if _, exists := s.runs[run.RunID]; exists {
			return domain.ErrIntegrity
		}
		if _, exists := s.idem[record.KeyDigest]; exists {
			return domain.ErrIntegrity
		}
		s.runs[run.RunID] = run
		s.idem[record.KeyDigest] = record
	}
	for _, run := range s.runs {
		if err := s.validatePublishedArtifacts(run); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

func (s *FileStore) ResolveIdempotency(ctx context.Context, digest string) (domain.IdempotencyRecord, bool, error) {
	if err := ctx.Err(); err != nil {
		return domain.IdempotencyRecord{}, false, err
	}
	if !digestPattern.MatchString(digest) {
		return domain.IdempotencyRecord{}, false, domain.ErrInvalidRunInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.IdempotencyRecord{}, false, domain.ErrRunStoreUnavailable
	}
	record, ok := s.idem[digest]
	return record, ok, nil
}

func (s *FileStore) CreateOrResolve(ctx context.Context, run domain.RecommendationRun, record domain.IdempotencyRecord) (domain.RecommendationRun, bool, error) {
	if err := ctx.Err(); err != nil {
		return domain.RecommendationRun{}, false, err
	}
	if !validInitialRun(run) || !digestPattern.MatchString(run.RequestHash) || !digestPattern.MatchString(record.KeyDigest) || record.RunID != run.RunID || record.RequestHash != run.RequestHash {
		return domain.RecommendationRun{}, false, domain.ErrInvalidRunInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.RecommendationRun{}, false, domain.ErrRunStoreUnavailable
	}
	if prior, ok := s.idem[record.KeyDigest]; ok {
		return cloneRun(s.runs[prior.RunID]), false, nil
	}
	if _, exists := s.runs[run.RunID]; exists {
		return domain.RecommendationRun{}, false, domain.ErrIntegrity
	}
	encoded, err := encodeDiskRecord(run, record)
	if err != nil {
		return domain.RecommendationRun{}, false, domain.ErrRunStoreUnavailable
	}
	if int64(len(encoded)) > s.limits.MaxRunBytes || s.logicalSizeLocked()+int64(len(encoded)) > s.limits.MaxTotalBytes {
		return domain.RecommendationRun{}, false, domain.ErrQuotaExceeded
	}
	if err := atomicWrite(filepath.Join(s.records, recordFileName(run.RunID)), encoded); err != nil {
		return domain.RecommendationRun{}, false, domain.ErrRunStoreUnavailable
	}
	s.runs[run.RunID] = cloneRun(run)
	s.idem[record.KeyDigest] = record
	return cloneRun(run), true, nil
}

func (s *FileStore) Get(ctx context.Context, id domain.RunID) (domain.RecommendationRun, error) {
	if err := ctx.Err(); err != nil {
		return domain.RecommendationRun{}, err
	}
	if !validRunID(id) {
		return domain.RecommendationRun{}, domain.ErrInvalidRunInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.RecommendationRun{}, domain.ErrRunStoreUnavailable
	}
	run, ok := s.runs[id]
	if !ok {
		return domain.RecommendationRun{}, domain.ErrRunNotFound
	}
	if err := s.validatePublishedArtifacts(run); err != nil {
		return domain.RecommendationRun{}, err
	}
	return cloneRun(run), nil
}

func (s *FileStore) Transition(ctx context.Context, transition domain.RunTransition) (domain.RecommendationRun, error) {
	if err := ctx.Err(); err != nil {
		return domain.RecommendationRun{}, err
	}
	if !validRunID(transition.RunID) {
		return domain.RecommendationRun{}, domain.ErrInvalidRunInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.RecommendationRun{}, domain.ErrRunStoreUnavailable
	}
	run, ok := s.runs[transition.RunID]
	if !ok {
		return domain.RecommendationRun{}, domain.ErrRunNotFound
	}
	if run.Version != transition.ExpectedVersion {
		return domain.RecommendationRun{}, domain.ErrRunVersionConflict
	}
	if err := usecase.ValidateRecommendationRunTransition(run, transition); err != nil {
		return domain.RecommendationRun{}, err
	}
	updated := cloneRun(run)
	updated.Lifecycle = transition.Lifecycle
	updated.Stage = transition.Stage
	updated.Progress = transition.Progress
	updated.RecommendationStatus = transition.RecommendationStatus
	updated.MarketDiscoveryStatus = transition.MarketDiscoveryStatus
	updated.ImprovementStatus = transition.ImprovementStatus
	updated.Failure = transition.Failure
	updated.Version++
	if err := s.persistUpdatedRunLocked(updated); err != nil {
		return domain.RecommendationRun{}, err
	}
	s.runs[updated.RunID] = updated
	return cloneRun(updated), nil
}

func (s *FileStore) PublishArtifact(ctx context.Context, publication domain.ArtifactPublication) (domain.RecommendationRun, error) {
	if err := ctx.Err(); err != nil {
		return domain.RecommendationRun{}, err
	}
	entries, err := validateArtifactDocument(publication)
	if err != nil {
		return domain.RecommendationRun{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.RecommendationRun{}, domain.ErrRunStoreUnavailable
	}
	run, ok := s.runs[publication.RunID]
	if !ok {
		return domain.RecommendationRun{}, domain.ErrRunNotFound
	}
	if run.Version != publication.ExpectedVersion {
		return domain.RecommendationRun{}, domain.ErrRunVersionConflict
	}
	for _, artifact := range run.Artifacts {
		if artifact.Name == publication.Name {
			return domain.RecommendationRun{}, domain.ErrArtifactImmutable
		}
	}
	if int64(len(publication.Document)) > s.limits.MaxArtifactBytes {
		return domain.RecommendationRun{}, domain.ErrQuotaExceeded
	}
	digest := sha256.Sum256(publication.Document)
	hash := "sha256:" + hex.EncodeToString(digest[:])
	updated := cloneRun(run)
	updated.Version++
	updated.Artifacts = append(updated.Artifacts, domain.RunArtifactIndex{
		Name: publication.Name, SchemaVersion: "1", Hash: hash,
		SizeBytes: int64(len(publication.Document)), EntryCount: len(entries), Version: updated.Version, State: domain.ArtifactPublished,
	})
	recordBytes, err := s.encodedRecordForRunLocked(updated)
	if err != nil {
		return domain.RecommendationRun{}, domain.ErrRunStoreUnavailable
	}
	runBytes := int64(len(recordBytes)) + s.blobBytesForRunLocked(run.RunID) + int64(len(publication.Document))
	totalDelta := int64(len(recordBytes)-s.recordSizeLocked(run.RunID)) + int64(len(publication.Document))
	if runBytes > s.limits.MaxRunBytes || s.logicalSizeLocked()+totalDelta > s.limits.MaxTotalBytes {
		return domain.RecommendationRun{}, domain.ErrQuotaExceeded
	}
	blobPath := filepath.Join(s.blobs, blobFileName(publication.RunID, publication.Name))
	if err := atomicWrite(blobPath, publication.Document); err != nil {
		return domain.RecommendationRun{}, domain.ErrRunStoreUnavailable
	}
	readBack, err := os.ReadFile(blobPath)
	if err != nil || !hmac.Equal(readBack, publication.Document) {
		return domain.RecommendationRun{}, domain.ErrIntegrity
	}
	if err := atomicWrite(filepath.Join(s.records, recordFileName(run.RunID)), recordBytes); err != nil {
		return domain.RecommendationRun{}, domain.ErrRunStoreUnavailable
	}
	s.runs[updated.RunID] = updated
	return cloneRun(updated), nil
}

func (s *FileStore) ReadArtifact(ctx context.Context, request domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	if err := ctx.Err(); err != nil {
		return domain.ArtifactPage{}, err
	}
	if !validArtifactRequest(request) {
		return domain.ArtifactPage{}, domain.ErrInvalidRunInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.ArtifactPage{}, domain.ErrRunStoreUnavailable
	}
	run, ok := s.runs[request.RunID]
	if !ok {
		return domain.ArtifactPage{}, domain.ErrRunNotFound
	}
	var decoded *cursorPayload
	if request.Cursor != "" {
		payload, err := s.decodeCursor(request.Cursor)
		if err != nil || payload.SchemaVersion != "1" || payload.RunID != request.RunID || payload.Artifact != request.Artifact {
			return domain.ArtifactPage{}, domain.ErrCursorInvalid
		}
		decoded = &payload
	}
	index, ok := findArtifact(run, request.Artifact)
	if !ok || index.State != domain.ArtifactPublished {
		if requiredArtifact(run, request.Artifact) {
			return domain.ArtifactPage{}, domain.ErrIntegrity
		}
		return domain.ArtifactPage{}, domain.ErrArtifactNotReady
	}
	document, entries, err := s.readAndValidateBlob(run.RunID, index)
	if err != nil {
		return domain.ArtifactPage{}, err
	}
	_ = document
	offset := 0
	if decoded != nil {
		if decoded.ArtifactVersion != index.Version || decoded.Hash != index.Hash || decoded.Offset < 0 || decoded.Offset > len(entries) {
			return domain.ArtifactPage{}, domain.ErrCursorInvalid
		}
		offset = decoded.Offset
	}
	end := offset + request.Limit
	if end > len(entries) {
		end = len(entries)
	}
	pageEntries := make([][]byte, end-offset)
	for index := range pageEntries {
		pageEntries[index] = append([]byte(nil), entries[offset+index]...)
	}
	var next *string
	if end < len(entries) {
		cursor, err := s.encodeCursor(cursorPayload{SchemaVersion: "1", RunID: request.RunID, Artifact: request.Artifact, ArtifactVersion: index.Version, Hash: index.Hash, Offset: end})
		if err != nil {
			return domain.ArtifactPage{}, domain.ErrRunStoreUnavailable
		}
		next = &cursor
	}
	return domain.ArtifactPage{SchemaVersion: "1", RunID: request.RunID, Artifact: request.Artifact, ArtifactHash: index.Hash, Entries: pageEntries, NextCursor: next}, nil
}

func (s *FileStore) InterruptNonTerminal(ctx context.Context) ([]domain.RecommendationRun, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, domain.ErrRunStoreUnavailable
	}
	ids := make([]string, 0, len(s.runs))
	for id := range s.runs {
		ids = append(ids, string(id))
	}
	sort.Strings(ids)
	changed := []domain.RecommendationRun{}
	for _, rawID := range ids {
		run := s.runs[domain.RunID(rawID)]
		if run.Lifecycle != domain.RunQueued && run.Lifecycle != domain.RunRunning && run.Lifecycle != domain.RunCancelling {
			continue
		}
		updated := cloneRun(run)
		updated.Lifecycle = domain.RunInterrupted
		updated.Stage = domain.RunStageTerminal
		updated.RecommendationStatus = nil
		updated.Failure = nil
		updated.Version++
		if err := s.persistUpdatedRunLocked(updated); err != nil {
			return changed, err
		}
		s.runs[updated.RunID] = updated
		changed = append(changed, cloneRun(updated))
	}
	return changed, nil
}

type artifactEnvelope struct {
	SchemaVersion string                 `json:"schema_version"`
	Artifact      domain.RunArtifactName `json:"artifact"`
	Entries       []json.RawMessage      `json:"entries"`
}

func validateArtifactDocument(publication domain.ArtifactPublication) ([]json.RawMessage, error) {
	if !validRunID(publication.RunID) || !validArtifactName(publication.Name) || publication.SchemaVersion != "1" || len(publication.Document) == 0 {
		return nil, domain.ErrInvalidRunInput
	}
	decoder := json.NewDecoder(strings.NewReader(string(publication.Document)))
	decoder.DisallowUnknownFields()
	var envelope artifactEnvelope
	if err := decoder.Decode(&envelope); err != nil || decoder.Decode(&struct{}{}) != io.EOF || envelope.SchemaVersion != "1" || envelope.Artifact != publication.Name || envelope.Entries == nil {
		return nil, domain.ErrInvalidRunInput
	}
	for _, entry := range envelope.Entries {
		if !json.Valid(entry) {
			return nil, domain.ErrInvalidRunInput
		}
	}
	return envelope.Entries, nil
}

func validRunID(id domain.RunID) bool { return storeRunIDPattern.MatchString(string(id)) }

func validArtifactName(name domain.RunArtifactName) bool {
	for _, allowed := range []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance, domain.ArtifactConnectorCapabilities, domain.ArtifactMarketPlan, domain.ArtifactMarketSnapshot, domain.ArtifactMarketCandidates, domain.ArtifactMarketGaps, domain.ArtifactMechanismCatalog, domain.ArtifactMechanismGaps, domain.ArtifactArchitectureComparison, domain.ArtifactImprovementHypotheses, domain.ArtifactHypothesisGaps, domain.ArtifactValidationPlans, domain.ArtifactSynthesisProvenance} {
		if name == allowed {
			return true
		}
	}
	return false
}

func validArtifactRequest(request domain.ArtifactReadRequest) bool {
	return validRunID(request.RunID) && validArtifactName(request.Artifact) && request.Limit >= 1 && request.Limit <= 50 && len(request.Cursor) <= 512
}

func cloneRun(run domain.RecommendationRun) domain.RecommendationRun {
	run.Request.Idea.Users = cloneStrings(run.Request.Idea.Users)
	run.Request.Idea.Capabilities = cloneStrings(run.Request.Idea.Capabilities)
	run.Request.Idea.Constraints = cloneStrings(run.Request.Idea.Constraints)
	run.Request.Idea.NonGoals = cloneStrings(run.Request.Idea.NonGoals)
	run.Request.Idea.PreferredStack = cloneStrings(run.Request.Idea.PreferredStack)
	run.Request.Policy.AllowedLicenses = cloneStrings(run.Request.Policy.AllowedLicenses)
	run.Request.Search.Locales = cloneStrings(run.Request.Search.Locales)
	if run.MarketRequest != nil {
		value := *run.MarketRequest
		value.Idea.Users = cloneStrings(value.Idea.Users)
		value.Idea.Capabilities = cloneStrings(value.Idea.Capabilities)
		value.Idea.Constraints = cloneStrings(value.Idea.Constraints)
		value.Idea.NonGoals = cloneStrings(value.Idea.NonGoals)
		value.Idea.PreferredStack = cloneStrings(value.Idea.PreferredStack)
		value.QueryProposals = append([]domain.MarketQueryProposal{}, value.QueryProposals...)
		value.Markets = cloneStrings(value.Markets)
		value.ConnectorIDs = cloneStrings(value.ConnectorIDs)
		run.MarketRequest = &value
	}
	if run.ImprovementRequest != nil {
		value := *run.ImprovementRequest
		value.Idea.Users = cloneStrings(value.Idea.Users)
		value.Idea.Capabilities = cloneStrings(value.Idea.Capabilities)
		value.Idea.Constraints = cloneStrings(value.Idea.Constraints)
		value.Idea.NonGoals = cloneStrings(value.Idea.NonGoals)
		value.Idea.PreferredStack = cloneStrings(value.Idea.PreferredStack)
		value.SourceArtifacts = append([]domain.SourceArtifactRef{}, value.SourceArtifacts...)
		value.Proposals = append([]domain.ImprovementProposal{}, value.Proposals...)
		for index := range value.Proposals {
			value.Proposals[index].MechanismRefs = cloneStrings(value.Proposals[index].MechanismRefs)
			value.Proposals[index].ProductObservationRefs = cloneStrings(value.Proposals[index].ProductObservationRefs)
			value.Proposals[index].Assumptions = cloneStrings(value.Proposals[index].Assumptions)
		}
		run.ImprovementRequest = &value
	}
	if run.Artifacts != nil {
		artifacts := make([]domain.RunArtifactIndex, len(run.Artifacts))
		copy(artifacts, run.Artifacts)
		run.Artifacts = artifacts
	}
	if run.RecommendationStatus != nil {
		value := *run.RecommendationStatus
		run.RecommendationStatus = &value
	}
	if run.MarketDiscoveryStatus != nil {
		value := *run.MarketDiscoveryStatus
		run.MarketDiscoveryStatus = &value
	}
	if run.ImprovementStatus != nil {
		value := *run.ImprovementStatus
		run.ImprovementStatus = &value
	}
	if run.Failure != nil {
		value := *run.Failure
		run.Failure = &value
	}
	return run
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	out := make([]string, len(values))
	copy(out, values)
	return out
}

func recordChecksum(run domain.RecommendationRun, record domain.IdempotencyRecord) (string, error) {
	encoded, err := json.Marshal(struct {
		Run         runAlias                 `json:"run"`
		Idempotency domain.IdempotencyRecord `json:"idempotency"`
	}{runAlias(run), record})
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func encodeDiskRecord(run domain.RecommendationRun, record domain.IdempotencyRecord) ([]byte, error) {
	checksum, err := recordChecksum(run, record)
	if err != nil {
		return nil, err
	}
	return json.Marshal(diskRecord{SchemaVersion: "1", Run: runAlias(run), Idempotency: record, Checksum: checksum})
}

func decodeDiskRecord(encoded []byte) (domain.RecommendationRun, domain.IdempotencyRecord, error) {
	decoder := json.NewDecoder(strings.NewReader(string(encoded)))
	decoder.DisallowUnknownFields()
	var envelope diskRecord
	if err := decoder.Decode(&envelope); err != nil || decoder.Decode(&struct{}{}) != io.EOF || envelope.SchemaVersion != "1" {
		return domain.RecommendationRun{}, domain.IdempotencyRecord{}, domain.ErrIntegrity
	}
	run := domain.RecommendationRun(envelope.Run)
	checksum, err := recordChecksum(run, envelope.Idempotency)
	if err != nil || !hmac.Equal([]byte(checksum), []byte(envelope.Checksum)) {
		return domain.RecommendationRun{}, domain.IdempotencyRecord{}, domain.ErrIntegrity
	}
	return run, envelope.Idempotency, nil
}

func recordFileName(id domain.RunID) string { return opaqueName("run\x00"+string(id)) + ".json" }
func blobFileName(id domain.RunID, name domain.RunArtifactName) string {
	return opaqueName("blob\x00"+string(id)+"\x00"+string(name)) + ".blob"
}
func opaqueName(value string) string {
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}

func atomicWrite(path string, data []byte) error {
	directory := filepath.Dir(path)
	temporary, err := os.OpenFile(filepath.Join(directory, ".write-"+opaqueName(path+fmt.Sprint(len(data)))+".tmp"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	clean := false
	defer func() {
		_ = temporary.Close()
		if !clean {
			_ = os.Remove(temporaryPath)
		}
	}()
	if _, err := temporary.Write(data); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	dir, err := os.Open(directory)
	if err != nil {
		return err
	}
	err = dir.Sync()
	_ = dir.Close()
	if err != nil {
		return err
	}
	clean = true
	return nil
}

func (s *FileStore) recordForRunLocked(id domain.RunID) (domain.IdempotencyRecord, bool) {
	for _, record := range s.idem {
		if record.RunID == id {
			return record, true
		}
	}
	return domain.IdempotencyRecord{}, false
}

func (s *FileStore) encodedRecordForRunLocked(run domain.RecommendationRun) ([]byte, error) {
	record, ok := s.recordForRunLocked(run.RunID)
	if !ok {
		return nil, domain.ErrIntegrity
	}
	return encodeDiskRecord(run, record)
}

func (s *FileStore) persistUpdatedRunLocked(run domain.RecommendationRun) error {
	encoded, err := s.encodedRecordForRunLocked(run)
	if err != nil {
		return domain.ErrIntegrity
	}
	delta := int64(len(encoded) - s.recordSizeLocked(run.RunID))
	if int64(len(encoded))+s.blobBytesForRunLocked(run.RunID) > s.limits.MaxRunBytes || s.logicalSizeLocked()+delta > s.limits.MaxTotalBytes {
		return domain.ErrQuotaExceeded
	}
	if err := atomicWrite(filepath.Join(s.records, recordFileName(run.RunID)), encoded); err != nil {
		return domain.ErrRunStoreUnavailable
	}
	return nil
}

func (s *FileStore) recordSizeLocked(id domain.RunID) int {
	info, err := os.Stat(filepath.Join(s.records, recordFileName(id)))
	if err != nil {
		return 0
	}
	return int(info.Size())
}

func (s *FileStore) blobBytesForRunLocked(id domain.RunID) int64 {
	var total int64
	for _, artifact := range s.runs[id].Artifacts {
		if artifact.State == domain.ArtifactPublished {
			total += artifact.SizeBytes
		}
	}
	return total
}

func (s *FileStore) logicalSizeLocked() int64 {
	var total int64
	for id := range s.runs {
		total += int64(s.recordSizeLocked(id)) + s.blobBytesForRunLocked(id)
	}
	return total
}

func findArtifact(run domain.RecommendationRun, name domain.RunArtifactName) (domain.RunArtifactIndex, bool) {
	for _, artifact := range run.Artifacts {
		if artifact.Name == name {
			if artifact.Hash == "" {
				artifact.Hash = artifact.ArtifactHash
			}
			return artifact, true
		}
	}
	return domain.RunArtifactIndex{}, false
}

func requiredArtifact(run domain.RecommendationRun, name domain.RunArtifactName) bool {
	switch run.Lifecycle {
	case domain.RunSucceeded:
		for _, required := range requiredArtifactNames(run) {
			if name == required {
				return true
			}
		}
		return false
	case domain.RunFailed:
		return name == domain.ArtifactRequest || name == domain.ArtifactFailures
	case domain.RunCancelled:
		return name == domain.ArtifactRequest
	}
	return false
}

func requiredArtifactNames(run domain.RecommendationRun) []domain.RunArtifactName {
	if run.RunKind == domain.RunKindMarketDiscovery {
		return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactConnectorCapabilities, domain.ArtifactMarketPlan, domain.ArtifactMarketSnapshot, domain.ArtifactMarketCandidates, domain.ArtifactMarketGaps}
	}
	if run.RunKind == domain.RunKindImprovementSynthesis {
		return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactArchitectureComparison, domain.ArtifactImprovementHypotheses, domain.ArtifactHypothesisGaps, domain.ArtifactValidationPlans, domain.ArtifactSynthesisProvenance}
	}
	if run.RunKind == domain.RunKindRecommendation {
		return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance, domain.ArtifactMechanismCatalog, domain.ArtifactMechanismGaps}
	}
	return []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance}
}

func (s *FileStore) validatePublishedArtifacts(run domain.RecommendationRun) error {
	for _, artifact := range run.Artifacts {
		if !validArtifactName(artifact.Name) || artifact.State != domain.ArtifactPublished {
			return domain.ErrIntegrity
		}
		if _, _, err := s.readAndValidateBlob(run.RunID, artifact); err != nil {
			return err
		}
	}
	for _, name := range requiredArtifactNames(run) {
		if requiredArtifact(run, name) {
			if _, ok := findArtifact(run, name); !ok {
				return domain.ErrIntegrity
			}
		}
	}
	return nil
}

func validInitialRun(run domain.RecommendationRun) bool {
	if !validRunID(run.RunID) || run.SchemaVersion != "1" || run.Version != 1 || run.Lifecycle != domain.RunQueued || run.Stage != domain.RunStageQueued || run.Artifacts == nil {
		return false
	}
	if run.RunKind == domain.RunKindMarketDiscovery {
		return run.MarketRequest != nil && run.ImprovementRequest == nil && run.Progress == (domain.RecommendationRunProgress{Planned: 4}) && run.RecommendationStatus == nil && run.MarketDiscoveryStatus == nil && run.ImprovementStatus == nil
	}
	if run.RunKind == domain.RunKindImprovementSynthesis {
		return run.ImprovementRequest != nil && run.MarketRequest == nil && run.Progress == (domain.RecommendationRunProgress{Planned: 5}) && run.RecommendationStatus == nil && run.MarketDiscoveryStatus == nil && run.ImprovementStatus == nil
	}
	return run.MarketRequest == nil && run.ImprovementRequest == nil && run.Progress == (domain.RecommendationRunProgress{Planned: 8}) && run.MarketDiscoveryStatus == nil && run.ImprovementStatus == nil
}

func (s *FileStore) readAndValidateBlob(id domain.RunID, index domain.RunArtifactIndex) ([]byte, []json.RawMessage, error) {
	encoded, err := os.ReadFile(filepath.Join(s.blobs, blobFileName(id, index.Name)))
	if err != nil {
		return nil, nil, domain.ErrIntegrity
	}
	digest := sha256.Sum256(encoded)
	hash := "sha256:" + hex.EncodeToString(digest[:])
	if hash != index.Hash || int64(len(encoded)) != index.SizeBytes {
		return nil, nil, domain.ErrIntegrity
	}
	entries, err := validateArtifactDocument(domain.ArtifactPublication{RunID: id, Name: index.Name, SchemaVersion: "1", Document: encoded})
	if err != nil || len(entries) != index.EntryCount {
		return nil, nil, domain.ErrIntegrity
	}
	return encoded, entries, nil
}

type cursorPayload struct {
	SchemaVersion   string                 `json:"schema_version"`
	RunID           domain.RunID           `json:"run_id"`
	Artifact        domain.RunArtifactName `json:"artifact"`
	ArtifactVersion uint64                 `json:"artifact_version"`
	Hash            string                 `json:"hash"`
	Offset          int                    `json:"offset"`
}

func (s *FileStore) encodeCursor(payload cursorPayload) (string, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write(encoded)
	return base64.RawURLEncoding.EncodeToString(encoded) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (s *FileStore) decodeCursor(cursor string) (cursorPayload, error) {
	if len(cursor) > 512 || strings.Count(cursor, ".") != 1 {
		return cursorPayload{}, domain.ErrCursorInvalid
	}
	parts := strings.Split(cursor, ".")
	if parts[0] == "" || parts[1] == "" || strings.Contains(cursor, "=") {
		return cursorPayload{}, domain.ErrCursorInvalid
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return cursorPayload{}, domain.ErrCursorInvalid
	}
	macBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || len(macBytes) != sha256.Size {
		return cursorPayload{}, domain.ErrCursorInvalid
	}
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write(payloadBytes)
	if !hmac.Equal(mac.Sum(nil), macBytes) {
		return cursorPayload{}, domain.ErrCursorInvalid
	}
	decoder := json.NewDecoder(strings.NewReader(string(payloadBytes)))
	decoder.DisallowUnknownFields()
	var payload cursorPayload
	if err := decoder.Decode(&payload); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return cursorPayload{}, domain.ErrCursorInvalid
	}
	return payload, nil
}
