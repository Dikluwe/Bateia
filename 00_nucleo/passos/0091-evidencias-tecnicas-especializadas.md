# Passo 0091 — Evidências técnicas especializadas

## Objetivo

Criar um contrato comum para achados de busca semântica, mineração, APIs,
arqueologia, recuperação arquitetural, impacto de mudanças, SBOM/SCA,
compatibilidade de licença, saúde do projeto e proveniência de build. Esses
achados enriquecem o modo avançado sem alterar `search` nem promover uma
observação externa a prova determinística de capacidade.

## Contrato

`TechnicalEvidenceBatch` (`schema_version: "1"`) contém `findings`. Cada
achado exige ID único, técnica, candidato, alegação, resultado
`SUPPORTS|REFUTES|INCONCLUSIVE`, referência de fonte, revisão ou digest,
método e data RFC3339. `valid_until` e limitações são opcionais. A fonte deve
apontar para um artefato verificável; o contrato valida forma e consistência,
**não confirma a veracidade da fonte**.

O conciliador agrupa achados por `(candidate, claim)`, separa achados vencidos,
marca suporte/refutação simultâneos como `CONFLICTING` e ordena a saída para
reprodutibilidade. Um achado vencido não sustenta uma alegação atual. O campo
`decision_impact` é sempre `ADVISORY_ONLY` nesta etapa.

## Comandos

`bateia idea evidence reconcile --input achados-tecnicos.json --as-of 2026-09-24T00:00:00Z --output revisao.json`

`bateia idea decide --idea ideia.json --assessment assessment.json --policy reuse-policy.json --technical-evidence achados-tecnicos.json --evidence-as-of 2026-09-24T00:00:00Z`

O segundo comando anexa `technical_evidence` ao relatório de decisão. Sem a
opção, mantém o JSON antigo. A decisão determinística e seus gates não mudam.
`lab/bancada/corpora/technical-evidence-illustrative.json` é apenas um exemplo
de formato; seus achados são fictícios e **não são evidência sobre um projeto real**.

## Próxima etapa

Criar adaptadores verificáveis por técnica, capazes de conferir a fonte e
converter alegações específicas em `EvidenceSource` admissíveis. Só então
avaliar mudanças nas regras de decisão, com testes de regressão e curva
risco–cobertura. Compatibilidade de licença permanece sujeita ao modo concreto
de integração e à revisão apropriada.

A primeira fatia verificável foi implementada no passo 0092.
