/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/repository_sources.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func ParseRepositoryIdentity(value string) (*domain.RepositoryIdentity, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	repository := value
	lower := strings.ToLower(value)
	switch {
	case strings.HasPrefix(lower, "git@github.com:"):
		repository = value[len("git@github.com:"):]
	case strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "ssh://"):
		parsed, err := url.Parse(value)
		if err != nil || !strings.EqualFold(parsed.Hostname(), "github.com") || parsed.RawQuery != "" || parsed.Fragment != "" {
			return nil, fmt.Errorf("target-repo inválido %q; use owner/repo ou URL GitHub", value)
		}
		repository = parsed.Path
	}
	repository = strings.Trim(strings.TrimSpace(repository), "/")
	if strings.HasSuffix(strings.ToLower(repository), ".git") {
		repository = repository[:len(repository)-4]
	}
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || invalidTargetRepositoryPart(parts[0]) || invalidTargetRepositoryPart(parts[1]) {
		return nil, fmt.Errorf("target-repo inválido %q; use owner/repo ou URL GitHub", value)
	}
	return &domain.RepositoryIdentity{Owner: strings.ToLower(parts[0]), Repo: strings.ToLower(parts[1])}, nil
}

func invalidTargetRepositoryPart(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || strings.ContainsAny(value, " \t\r\n@:#?")
}

func ParseGitRepositoryHints(data []byte) []domain.RepositoryIdentity {
	seen := map[string]domain.RepositoryIdentity{}
	remoteSection := false
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			remoteSection = strings.HasPrefix(strings.ToLower(line), "[remote ")
			continue
		}
		if !remoteSection {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || !strings.EqualFold(strings.TrimSpace(key), "url") {
			continue
		}
		identity, parseErr := ParseRepositoryIdentity(strings.TrimSpace(value))
		if parseErr != nil || identity == nil {
			continue
		}
		seen[ideaRunRepositoryKey(identity.Owner, identity.Repo)] = *identity
	}
	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]domain.RepositoryIdentity, 0, len(keys))
	for _, key := range keys {
		out = append(out, seen[key])
	}
	return out
}
func LocalRepositoryCandidate(root, base, digest string) domain.RepositoryCandidate {
	sum := sha256.Sum256([]byte(root))
	key := hex.EncodeToString(sum[:])
	return domain.RepositoryCandidate{Owner: "local", Repo: strings.ToLower(base) + "-" + key[:8], URL: "file://" + root, CommitID: digest}
}
