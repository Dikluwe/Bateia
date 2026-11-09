# Declarações repetidas e ambíguas do oráculo assistido

Testar as APIs públicas `ValidateAssistedJudgments` e
`DetermineHumanEscalations` sem I/O. Declarações distintas do mesmo candidato não
podem ser reduzidas por ordem de chegada: variar veredito, confiança, razões e
ordem estrutural dos vínculos da IA; variar ação, motivo e veredito da revisão.

Conservar todas as ocorrências, inclusive repetições idênticas, em históricos
ordenados pelo JSON completo. Permutar entradas e exigir relatório idêntico sem
mutá-las. Ambiguidade torna Effective UNKNOWN, exige escalonamento e
INCONCLUSIVE_HUMAN_REVIEW, com AMBIGUOUS_AI_JUDGMENTS ou
AMBIGUOUS_HUMAN_REVIEWS. Nenhum candidato ambíguo incrementa Reviewed,
DefinitiveReviewed, DefinitiveAgreements, DefinitiveOverrides ou
ResolvedUncertain; Escalated conta o candidato. Julgamentos distintos e válidos
podem conservar grounding 1, mas qualquer vínculo inválido torna o grounding do
candidato zero e FAIL_EVIDENCE_GROUNDING prevalece independentemente da ordem.

Declarações estruturalmente idênticas são idempotentes em métricas, efetivo e
escalonamento, mas continuam presentes nos históricos. Conservar controle com
uma declaração/revisão aceita. Usar ainda corpus residual fora da amostra de 10%
para provar que duplicidade distinta impõe escalonamento, sem depender da escolha
da amostra. Não inferir cronologia nem resolver conflito por timestamp inexistente.
