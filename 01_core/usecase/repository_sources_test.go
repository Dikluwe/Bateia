/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/repository_sources_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func TestRepositorySourcesArePureAndCanonical(t *testing.T) {
	for _, ref := range []string{"Org/Repo", "https://github.com/Org/Repo.git", "git@github.com:Org/Repo.git"} {
		got, err := ParseRepositoryIdentity(ref)
		if err != nil || !reflect.DeepEqual(got, &domain.RepositoryIdentity{Owner: "org", Repo: "repo"}) {
			t.Fatalf("bad identity: %v %v", got, err)
		}
	}
	if _, err := ParseRepositoryIdentity("https://other.example/org/repo"); err == nil {
		t.Fatal("foreign authority accepted")
	}
	data := []byte("[remote \"origin\"]\nurl = https://github.com/Org/Repo.git\n[remote \"mirror\"]\nurl = git@github.com:org/repo.git\n")
	hints := ParseGitRepositoryHints(data)
	if len(hints) != 1 || hints[0].Repo != "repo" {
		t.Fatalf("hints: %+v", hints)
	}
	a := LocalRepositoryCandidate("/tmp/Sample", "Sample", "abc")
	b := LocalRepositoryCandidate("/tmp/Sample", "Sample", "abc")
	if !reflect.DeepEqual(a, b) || a.Owner != "local" || a.URL != "file:///tmp/Sample" || a.CommitID != "abc" {
		t.Fatal("local identity changed")
	}
}
