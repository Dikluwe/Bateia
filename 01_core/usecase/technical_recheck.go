/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_recheck.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const technicalRecheckMaxBytes = 4 << 20

// RecheckTechnicalSources never accepts caller-supplied bundle bytes. The
// configured provider reads an explicit GitHub candidate at a full commit SHA.
// This is a provider observation, not a cryptographic certificate on exported JSON.
func RecheckTechnicalSources(ctx context.Context, req domain.TechnicalSourceRecheckRequest, provider domain.ContentProvider, asOf time.Time) (domain.TechnicalSourceRecheckReport, error) {
	if ctx == nil {
		return domain.TechnicalSourceRecheckReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.TechnicalSourceRecheckReport{}, err
	}
	initial, err := VerifyTechnicalEvidence(req.Batch, nil, asOf)
	if err != nil {
		return domain.TechnicalSourceRecheckReport{}, err
	}
	if provider == nil || len(req.Candidates) > 10 || len(req.Batch.Findings) > 100 {
		return domain.TechnicalSourceRecheckReport{}, fmt.Errorf("recheck exige provedor e limita a 10 candidatos e 100 achados")
	}
	candidates := map[string]domain.RepositoryCandidate{}
	for _, c := range req.Candidates {
		if !validTechnicalRemoteCandidate(c) {
			return domain.TechnicalSourceRecheckReport{}, fmt.Errorf("recheck exige identidade GitHub e commit SHA completo: %s/%s", c.Owner, c.Repo)
		}
		ref := domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}
		if _, duplicate := candidates[ref.Key()]; duplicate {
			return domain.TechnicalSourceRecheckReport{}, fmt.Errorf("candidato duplicado: %s", ref.Key())
		}
		c.URL = "https://github.com/" + c.Owner + "/" + c.Repo
		candidates[ref.Key()] = c
	}
	active := map[string]bool{}
	for _, verification := range initial.Verifications {
		if verification.Status == domain.TechnicalUnverified {
			active[verification.FindingID] = true
		}
	}
	pathsByCandidate := map[string]map[string]domain.EvidenceKind{}
	for _, f := range req.Batch.Findings {
		if !active[f.ID] || f.Result != domain.TechnicalSupports {
			continue
		}
		kind := domain.EvidenceManifest
		if f.Predicate == domain.TechnicalLicenseID && f.Technique == domain.TechniqueLicense {
			kind = domain.EvidenceLicense
		} else if f.Predicate != domain.TechnicalSBOMComponent || f.Technique != domain.TechniqueSBOM {
			continue
		}
		c, ok := candidates[f.Candidate]
		if !ok || c.CommitID != f.Source.Revision || !validMechanismPath(f.Source.Ref) {
			return domain.TechnicalSourceRecheckReport{}, fmt.Errorf("achado %s não está ligado a candidato e caminho fixados", f.ID)
		}
		if kind == domain.EvidenceLicense && !authoritativeLicensePath(f.Source.Ref) {
			return domain.TechnicalSourceRecheckReport{}, fmt.Errorf("achado %s exige documento raiz de licença", f.ID)
		}
		if pathsByCandidate[f.Candidate] == nil {
			pathsByCandidate[f.Candidate] = map[string]domain.EvidenceKind{}
		}
		pathsByCandidate[f.Candidate][f.Source.Ref] = kind
	}
	report := domain.TechnicalSourceRecheckReport{SchemaVersion: "1", Basis: "PROVIDER_RECHECK", Bundles: []domain.EvidenceBundle{}}
	keys := make([]string, 0, len(pathsByCandidate))
	for key := range pathsByCandidate {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var totalBytes int
	for _, key := range keys {
		if err := ctx.Err(); err != nil {
			return domain.TechnicalSourceRecheckReport{}, err
		}
		c := candidates[key]
		ref := domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}
		paths := pathsByCandidate[key]
		needsLicense := false
		for _, kind := range paths {
			needsLicense = needsLicense || kind == domain.EvidenceLicense
		}
		bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Plan: domain.EvidenceExtractionPlan{Repository: ref, Candidate: c}, Items: []domain.CollectedEvidence{}}
		fail := func(message string) {
			failure := domain.ExecutionFailure{Scope: key, Stage: "technical_source_recheck", Message: message}
			bundle.Failures = append(bundle.Failures, failure)
			report.Failures = append(report.Failures, failure)
		}
		// Inspect the full tree so another root license cannot be omitted from
		// the fresh bundle to manufacture an unambiguous adoption gate.
		sizes := map[string]int64{}
		tree, treeErr := provider.FetchTree(ctx, c)
		if err := ctx.Err(); err != nil {
			return domain.TechnicalSourceRecheckReport{}, err
		}
		if treeErr != nil {
			fail("não foi possível conferir árvore da revisão: " + treeErr.Error())
		} else {
			for _, node := range tree {
				_, requested := paths[node.Path]
				licensePath := needsLicense && authoritativeLicensePath(node.Path)
				if !requested && !licensePath {
					continue
				}
				if node.Type != "file" {
					fail("documento solicitado não regular exige revisão de escopo")
					continue
				}
				if node.Size < 0 {
					fail("árvore declara tamanho negativo")
				}
				if node.Size > 0 {
					if sizes[node.Path] > 0 && sizes[node.Path] != node.Size {
						fail("árvore declara tamanhos conflitantes")
					}
					sizes[node.Path] = node.Size
				}
				if licensePath {
					paths[node.Path] = domain.EvidenceLicense
				}
			}
		}
		ordered := make([]string, 0, len(paths))
		for path := range paths {
			ordered = append(ordered, path)
		}
		sort.Strings(ordered)
		if len(ordered) > 20 {
			fail("releitura excede limite de 20 arquivos por candidato")
		} else if len(bundle.Failures) == 0 {
			if err := ctx.Err(); err != nil {
				return domain.TechnicalSourceRecheckReport{}, err
			}
			extracted, readErr := provider.ExtractRepoFiles(ctx, c, ordered)
			if err := ctx.Err(); err != nil {
				return domain.TechnicalSourceRecheckReport{}, err
			}
			if readErr != nil || extracted == nil {
				fail("provedor não devolveu os arquivos solicitados")
			} else if extracted.Ref != ref || extracted.Stale {
				fail("provedor devolveu revisão divergente ou conteúdo stale")
			} else {
				seen := map[string]bool{}
				for _, file := range extracted.Files {
					kind, expected := paths[file.Path]
					if !expected || seen[file.Path] {
						fail("provedor devolveu caminho inesperado ou duplicado")
						continue
					}
					seen[file.Path] = true
					if !validEvidenceFileSize(file, sizes[file.Path]) {
						fail("tamanho declarado diverge dos bytes extraídos")
						continue
					}
					if !utf8.ValidString(file.Content) {
						fail("conteúdo extraído não é UTF-8 válido")
						continue
					}
					totalBytes += len(file.Content)
					if totalBytes > technicalRecheckMaxBytes {
						fail("conteúdo excede limite total de 4 MiB")
						break
					}
					bundle.Items = append(bundle.Items, domain.CollectedEvidence{Kind: kind, Ref: ref, Path: file.Path, Content: file.Content, ContentHash: contentHash(file.Content)})
				}
				if len(seen) != len(paths) {
					fail("releitura incompleta dos caminhos solicitados")
				}
			}
		}
		bundle.Complete = len(bundle.Failures) == 0
		sort.Slice(bundle.Items, func(i, j int) bool { return bundle.Items[i].Path < bundle.Items[j].Path })
		report.Bundles = append(report.Bundles, bundle)
	}
	if err := ctx.Err(); err != nil {
		return domain.TechnicalSourceRecheckReport{}, err
	}
	report.Review, err = VerifyTechnicalEvidence(req.Batch, report.Bundles, asOf)
	if err := ctx.Err(); err != nil {
		return domain.TechnicalSourceRecheckReport{}, err
	}
	for i := range report.Review.Verifications {
		if report.Review.Verifications[i].IntegrityVerified {
			report.Review.Verifications[i].SourceAuthentication = "PROVIDER_OBSERVED"
		}
	}
	return report, err
}

func validTechnicalRemoteCandidate(c domain.RepositoryCandidate) bool {
	sha, err := hex.DecodeString(c.CommitID)
	if err != nil || len(sha) != 20 {
		return false
	}
	for _, part := range []string{c.Owner, c.Repo} {
		if part == "" || part == "." || part == ".." {
			return false
		}
		for _, char := range part {
			if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == '-' || char == '_' || char == '.') {
				return false
			}
		}
	}
	return c.URL == "" || strings.TrimSuffix(c.URL, "/") == "https://github.com/"+c.Owner+"/"+c.Repo
}
