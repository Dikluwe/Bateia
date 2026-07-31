/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/java.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePomXml_Basic(t *testing.T) {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>spring-boot-api</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>3.2.0</version>
        </dependency>
        <dependency>
            <groupId>org.junit.jupiter</groupId>
            <artifactId>junit-jupiter-api</artifactId>
            <version>5.10.0</version>
            <scope>test</scope>
        </dependency>
    </dependencies>
</project>`

	ref := domain.CommitRef{Owner: "spring-projects", Repo: "spring-boot-api"}
	manifest, err := ParsePomXml(ref, content)
	if err != nil {
		t.Fatalf("ParsePomXml falhou: %v", err)
	}

	if manifest.ModuleName != "spring-boot-api" {
		t.Errorf("Esperado ModuleName 'spring-boot-api', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "java" {
		t.Errorf("Esperado Language 'java', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}

	if manifest.Dependencies[0].Type != "direct" {
		t.Errorf("Esperado tipo 'direct' para spring-boot-starter-web, obtido '%s'", manifest.Dependencies[0].Type)
	}
	if manifest.Dependencies[1].Type != "dev" {
		t.Errorf("Esperado tipo 'dev' para junit-jupiter-api, obtido '%s'", manifest.Dependencies[1].Type)
	}
}

func TestParseBuildGradle_Basic(t *testing.T) {
	content := `
plugins {
    id 'org.springframework.boot' version '3.2.0'
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
`

	ref := domain.CommitRef{Owner: "gradle", Repo: "kotlin-app"}
	manifest, err := ParseBuildGradle(ref, content)
	if err != nil {
		t.Fatalf("ParseBuildGradle falhou: %v", err)
	}

	if manifest.Language != "java" {
		t.Errorf("Esperado Language 'java', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}
