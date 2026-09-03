# Oráculo independente de extração de mecanismos do Passo 0028

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/mechanism_extraction_step0028_oracle_test.go`

Use o split TEST protegido e ligue o sujeito somente pela interface pública
`MechanismExtractor`. Não conheça nem imponha construtor, tipo concreto ou campo
privado. Confira integridade dos fixtures por SHA-256 e execute controles de
append-only log, fila limitada, documentação isolada, testemunha cross-repo e
contraevidência central.

O join é exatamente `(kind,ref,path,content_hash)`; conteúdo serve somente para
recomputar hash. RuleID, symbol, line, role e witness kind vêm da avaliação e do
catálogo. Reordene bundle e assessments e exija resultado idêntico. Nunca aceite
cross-repo como `PROVEN`, documentação como prova completa ou canário/raw content
na projeção. Casos deliberadamente opacos permanecem `Unknown` e não contam no
mutation score.
