# Passo 0061 — 02_shell/report/action_assessment

WriteActionAssessments(sb *strings.Builder, assessments []domain.RequirementActionAssessment) apresenta ação/status/motivo/limites, todos campos de fonte e evidências. Limites Markdown12 candidatos e4 evidências, remanescentes no JSON. Toda string externa em code HTML com escaping HTML e representação citada (novas linhas escapadas); impedir Markdown/HTML injetado ou snippets como instruções. Não mudar métricas/status.

## Adendo r1
Antes do corte em 12, ordenar uma cópia dos candidatos estavelmente por aligned_hint, conflicting_hints, opposed_hint, unknown e qualquer valor desconhecido por último. Preservar ordem de entrada nos empates. Não mutar entrada, JSON ou ranking. A ordem serve à leitura e não representa confiança semântica nem comprovação.
