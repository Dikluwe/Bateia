# Contrato de oráculo black-box independente — Passo 0028 O9

**Revisão:** O9 vinculada a R4, catálogo v3 e protocolo O8  
**Regime:** protocolo Tekt completo  
**Papel:** autoridade segregada de ORÁCULOS  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

O avaliador continua recebendo somente `contract_input` e
`mutant_contract_output`. IDs, records, bindings, payloads, configs, markers,
expected e candidato permanecem fora da fronteira.

## Entradas públicas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo v3 | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Protocolo O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Coverage gap público O8 | `e66400db2d6bebe5b3be4d284058b51aa2c93f9d1f8ecd3544057e09d62330e1` |
| Falha discriminatória pública O8 | `23341b5e768e8f0031421e9bd0c4105426c327a27fc55e970149ff11648c2c09` |
| Manifest protegido predecessor O8 | `30a5af21d089938c35f5fed4b001cd6e714cdfae850057409ec0d47d5e8b916c` |

## Refinamentos O9

`R4.SYNTHESIS_START.NO_INVENTED_EVIDENCE` exige que refs de mecanismo e de
produto permaneçam em namespaces disjuntos e que cada categoria referenciada
tenha uma source declarada capaz de fornecê-la. Proposta continua proposta e não
se torna evidence por prefixo ou posição.

`SOURCE_ARTIFACT_INTERNAL_REFERENCE_INTEGRITY` exige um único `run_kind` por
`run_id`, pares permitidos e coerência do artifact causal associado ao mesmo
hash. A validação permanece anterior a resolution, idempotência e fila.

`R4.RUN_CANCEL.IDEMPOTENT_TERMINAL_TRANSITION` aceita uma única transição CAS
para cancelamento em curso ou de run enfileirado e exige estabilidade de versão
quando o estado observado já é terminal. Cancel não publica nem inicia trabalho.

`ARTIFACT_READ` reconhece como `Unknown` estados R4 explícitos de novelty,
transferência ou disponibilidade, além dos opacos O8 já suportados, somente após
bytes de ordem fixa, identidade, hash, size, paginação, ordenação e safety terem
passado. Uma falha decidível continua precedendo opacidade.

Todos os predicados O8 de provenance, novelty, source resolution, canonical
artifact bytes, cancel stability, replay e tools permanecem ativos.

## API pública

```go
func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

O9 não emite receipt, seal, certificado ou veredito.
