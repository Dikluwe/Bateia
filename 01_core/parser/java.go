/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/java.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"
	"strings"

	"bateia/01_core/domain"
)

var (
	pomArtifactRegex   = regexp.MustCompile(`(?s)<artifactId>(.*?)</artifactId>`)
	pomDependencyRegex = regexp.MustCompile(`(?s)<dependency>(.*?)</dependency>`)
	pomGroupRegex      = regexp.MustCompile(`(?s)<groupId>(.*?)</groupId>`)
	pomVersionRegex    = regexp.MustCompile(`(?s)<version>(.*?)</version>`)
	pomScopeRegex      = regexp.MustCompile(`(?s)<scope>(.*?)</scope>`)
	gradleDepRegex     = regexp.MustCompile(`(?m)^\s*(implementation|testImplementation|api|testCompile|runtimeOnly)\s*\(?['"]([^'"]+)['"]\)?`)
)

// ParsePomXml analisa o conteúdo de um manifesto Maven pom.xml.
func ParsePomXml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "java",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	// Extrair artifactId principal
	if matches := pomArtifactRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = strings.TrimSpace(matches[1])
	} else {
		manifest.ModuleName = ref.Repo
	}

	// Extrair bloco <dependency>
	depBlocks := pomDependencyRegex.FindAllStringSubmatch(content, -1)
	for _, block := range depBlocks {
		if len(block) < 2 {
			continue
		}
		inner := block[1]

		artifactMatches := pomArtifactRegex.FindStringSubmatch(inner)
		if len(artifactMatches) < 2 {
			continue
		}
		artifact := strings.TrimSpace(artifactMatches[1])

		group := ""
		if groupMatches := pomGroupRegex.FindStringSubmatch(inner); len(groupMatches) >= 2 {
			group = strings.TrimSpace(groupMatches[1])
		}

		version := ""
		if versionMatches := pomVersionRegex.FindStringSubmatch(inner); len(versionMatches) >= 2 {
			version = strings.TrimSpace(versionMatches[1])
		}

		depType := "direct"
		if scopeMatches := pomScopeRegex.FindStringSubmatch(inner); len(scopeMatches) >= 2 {
			if strings.TrimSpace(scopeMatches[1]) == "test" {
				depType = "dev"
			}
		}

		depName := artifact
		if group != "" {
			depName = group + ":" + artifact
		}

		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     strings.ToLower(depName),
			Version:  version,
			Indirect: false,
			Type:     depType,
		})
	}

	return manifest, nil
}

// ParseBuildGradle analisa o conteúdo de um manifesto Gradle (build.gradle ou build.gradle.kts).
func ParseBuildGradle(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "java",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	matches := gradleDepRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		scope := match[1]
		depString := strings.TrimSpace(match[2])

		depType := "direct"
		if strings.Contains(strings.ToLower(scope), "test") {
			depType = "dev"
		}

		// Quebrar 'group:artifact:version' se houver
		parts := strings.Split(depString, ":")
		name := depString
		version := ""
		if len(parts) >= 2 {
			name = parts[0] + ":" + parts[1]
		}
		if len(parts) >= 3 {
			version = parts[2]
		}

		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     strings.ToLower(name),
			Version:  version,
			Indirect: false,
			Type:     depType,
		})
	}

	return manifest, nil
}
