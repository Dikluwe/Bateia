# Passo 0073 — contrato público de localização intrarquivo

Papel: autor segregado do contrato pré-candidato. Baseline Git `03fd5f1` acrescido dos
artefatos públicos congelados do passo 0072. Este teste precede a implementação e não pode
ser adaptado ao patch candidato nem ao gate adversarial.

Materialize em `01_core/domain/evidence_step0073_contract_test.go` os valores públicos
`STRONG_CONTENT`, `SUPPORTING_CONTENT`, `DOCUMENTATION_CONTENT` e `PATH_ONLY`, além dos
campos JSON `matched_pattern` e `match_class` em `EvidenceSource`. Exija round-trip sem
perder linha, snippet ou classificação. `PATH_ONLY` existe para diagnóstico, mas o caso de
uso não pode emiti-lo como fonte probatória.

Autoridade completa: `00_nucleo/passos/0073-localizar-a-ocorrencia-probatoria-dentro-do-arquivo.md`.

