/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/capability_compiler.md
 * @layer L1
 * @updated 2026-09-24
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"fmt"
	"strings"
)

// CompileCapabilityPlan turns declared needs into verification questions.
// It neither searches candidates nor treats a question as evidence.
func CompileCapabilityPlan(intake domain.RecommendationIntake) domain.CapabilityCompilation {
	idea := normalizeIdea(intake.Request.Idea)
	out := domain.CapabilityCompilation{
		SchemaVersion: "1", IntakeHash: intake.InputHash, IntakeStatus: intake.Status,
		Goal: idea.Problem, Users: append([]string{}, idea.Users...),
		Requirements: []domain.CompiledRequirement{}, Questions: []domain.CompiledEvidenceQuestion{},
	}
	authorities := map[string]domain.InputAuthority{}
	for _, source := range intake.Provenance {
		authorities[source.Path] = source.Authority
	}
	seen := map[string]bool{}
	add := func(kind, text, path, priority string) {
		text = normalizeIdeaText(text)
		if text == "" {
			return
		}
		id := requirementID(kind, text)
		if seen[id] {
			return
		}
		seen[id] = true
		authority := authorities[path]
		if authority == "" {
			authority = domain.AuthorityUnresolved
		}
		r := domain.CompiledRequirement{
			ID: id, Kind: kind, Text: text, SourcePath: path,
			Authority: authority, Priority: priority,
		}
		out.Requirements = append(out.Requirements, r)
		for _, channel := range capabilityChannels(kind) {
			q := compiledEvidenceQuestion(r, channel)
			out.Questions = append(out.Questions, q)
		}
	}
	// Provenance is attached to the whole user-supplied list, not individual indices.
	addDeclaredList := func(kind string, values []string, path, priority string) {
		for i, value := range values {
			authorities[fmt.Sprintf("%s[%d]", path, i)] = authorities[path]
			add(kind, value, fmt.Sprintf("%s[%d]", path, i), priority)
		}
	}
	addDeclaredList("capability", intake.Request.Idea.Capabilities, "idea.capabilities", "MUST")
	addDeclaredList("quality", intake.Request.Idea.Qualities, "idea.qualities", "MUST")
	addDeclaredList("operating_condition", intake.Request.Idea.OperatingConditions, "idea.operating_conditions", "MUST")
	addDeclaredList("constraint", intake.Request.Idea.Constraints, "idea.constraints", "MUST")
	addDeclaredList("non_goal", intake.Request.Idea.NonGoals, "idea.non_goals", "MUST_NOT")
	addDeclaredList("preference", intake.Request.Idea.PreferredStack, "idea.preferred_stack", "PREFER")
	if intake.Request.Policy.SchemaVersion != "" {
		if len(intake.Request.Policy.AllowedLicenses) > 0 {
			add("license_policy", strings.Join(intake.Request.Policy.AllowedLicenses, ", "), "policy.allowed_licenses", "MUST")
		}
		if intake.Request.Policy.RequireLocal {
			add("local_policy", "execução local obrigatória", "policy.require_local", "MUST")
		}
		if !intake.Request.Policy.AllowArchived {
			add("archive_policy", "repositório não arquivado", "policy.allow_archived", "MUST")
		}
	}
	return out
}

func capabilityChannels(kind string) []string {
	switch kind {
	case "capability", "constraint", "non_goal":
		return []string{"documentation", "code", "tests", "execution"}
	case "quality", "operating_condition":
		return []string{"documentation", "tests", "execution"}
	case "preference":
		return []string{"documentation", "code"}
	case "license_policy":
		return []string{"license"}
	case "local_policy":
		return []string{"documentation", "code", "execution"}
	case "archive_policy":
		return []string{"maintenance"}
	}
	return nil
}

func compiledEvidenceQuestion(r domain.CompiledRequirement, channel string) domain.CompiledEvidenceQuestion {
	q := domain.CompiledEvidenceQuestion{
		RequirementID: r.ID, Channel: channel, Status: "PENDING",
		Limit: "Pergunta de verificação, não evidência; ausência de menção não prova ausência.",
	}
	sum := sha256.Sum256([]byte(r.ID + "\x00" + channel))
	q.ID = fmt.Sprintf("eq-%x", sum[:8])
	switch channel {
	case "documentation":
		q.Question = fmt.Sprintf("A documentação do candidato declara %q e suas condições ou limites?", r.Text)
		q.Support = "Declaração localizada, com versão e escopo."
		q.Refute = "Declaração explícita de incompatibilidade ou limitação."
		q.Limit = "Declaração documental não comprova implementação nem execução."
	case "code":
		q.Question = fmt.Sprintf("Que caminho de código implementa ou contradiz %q?", r.Text)
		q.Support = "Implementação localizada e ligada a uma interface utilizável."
		q.Refute = "Caminho explícito que viola o requisito."
		q.Limit = "Inspeção estática não comprova comportamento em execução; silêncio não prova ausência."
	case "tests":
		q.Question = fmt.Sprintf("Que teste cobre %q, incluindo falha ou limite relevante?", r.Text)
		q.Support = "Teste localizado, asserções e resultado vinculados à versão."
		q.Refute = "Teste que demonstra violação na condição declarada."
		q.Limit = "Existência do teste sem resultado não comprova execução nem cobertura completa."
	case "execution":
		q.Question = fmt.Sprintf("Em qual ensaio controlado se observa %q nas condições declaradas?", r.Text)
		q.Support = "Entrada, ambiente, versão, resultado observado e critério de aceitação."
		q.Refute = "Resultado reproduzível que viola o critério declarado."
		q.Limit = "Sem critério de aceitação ou autorização para executar, manter inconclusivo."
	case "license":
		q.Question = fmt.Sprintf("As licenças do candidato e componentes relevantes são compatíveis com a lista permitida: %s?", r.Text)
		q.Support = "Arquivos e metadados de licença identificados por versão e componente."
		q.Refute = "Licença identificada fora da política permitida."
		q.Limit = "Licença não demonstra capacidade; ambiguidade exige análise própria."
	case "maintenance":
		q.Question = "O repositório está arquivado na versão e origem examinadas?"
		q.Support = "Estado de arquivamento observado na origem do candidato."
		q.Refute = "Estado arquivado confirmado."
		q.Limit = "Não estar arquivado não prova manutenção ativa ou qualidade."
	}
	if r.Kind == "non_goal" {
		q.Question = fmt.Sprintf("Há ocorrência explícita do comportamento excluído %q neste canal?", r.Text)
		q.Support = "Negação delimitada ou ensaio que demonstra o bloqueio do comportamento."
		q.Refute = "Ocorrência positiva e reproduzível do comportamento excluído."
		q.Limit = "Ausência de menção ou de teste não prova ausência do comportamento."
	}
	return q
}
