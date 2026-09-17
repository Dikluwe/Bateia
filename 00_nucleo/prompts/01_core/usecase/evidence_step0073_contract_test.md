# Passo 0073 — contrato do localizador intrarquivo

Papel: autor segregado do contrato pré-candidato. Baseline Git `03fd5f1` acrescido dos
artefatos públicos congelados do passo 0072. Não ler ou alterar implementação candidata ou
gate adversarial do passo 0073.

Materialize em `01_core/usecase/evidence_step0073_contract_test.go` observações black-box
de `PlanEvidenceForIdea` e `AssessEvidence`:

- queries contêm somente `TargetHints`, sem padrões probatórios implícitos;
- `_M_limit`, `subspan` e `size() - __pos` produzem padrão, classe, linha e snippet da
  mesma ocorrência forte;
- `@file`, basename/path, comentários, strings e substrings não viram strong content;
- função/membro C++ não depende de keywords Go/Python/Rust;
- strong content precede apoio e documentação;
- comentário/string-only conserva Unknown lexical sem inventar ausência;
- vinte permutações preservam JSON e inputs.

Não exigir parser C++ completo. Macro, token-pasting e sintaxe fora do recorte permanecem
Unknown. Autoridade completa:
`00_nucleo/passos/0073-localizar-a-ocorrencia-probatoria-dentro-do-arquivo.md`.

