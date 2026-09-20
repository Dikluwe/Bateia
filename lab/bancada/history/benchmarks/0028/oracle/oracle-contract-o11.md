# Contrato de oráculo black-box independente — Passo 0028 O11

**Revisão:** O11 vinculada a R4, catálogo v3 e protocolo O8  
**Regime:** protocolo Tekt completo  
**Papel:** autoridade segregada de ORÁCULOS  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

O avaliador recebe exclusivamente `contract_input` e `contract_output` e deriva
`Preserved|Violated|Unknown` de predicados normativos. IDs de caso, records,
configs, bindings, hashes adversariais, markers, expected e candidato não entram
na API nem na decisão.

## Entradas públicas congeladas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo v3 | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Protocolo O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Coverage gap público O10 | `9f40154c07727341ca5cca983630b8e1a8e62d6fcadf8839f45d17a33cc5b05f` |
| Falha discriminatória pública O10 | `46df2393574b4895e327b5358375721466210fae752f8c52883fdb23ca8ec682` |
| Manifest protegido predecessor O9 | `b7960417c3e9de2c12b42fa56b5f9b73a12bc87585541800b3bd92d1596f4b0c` |

## Refinamentos normativos O11

`R4.SYNTHESIS_START.NO_INVENTED_EVIDENCE` exige que start não publique artifacts,
hashes, índices ou `layer_findings`: proposta externa não se converte em evidence
por efeito lateral observável. Todos os refs continuam sujeitos a namespace,
categoria de source e validação pré-fila.

`R4.SYNTHESIS_START.SOURCE_ARTIFACT_INTERNAL_REFERENCE_INTEGRITY` preserva um
único kind por source run, associação coerente de artifact/hash e identidade
histórica. O `run_id` novo do receipt não pode colidir com nenhum `run_id`
declarado em `source_artifacts`.

Em R4, mecanismo sustentado somente pela forma literal `DOCUMENTARY` usa o RuleID
documental efetivamente ligado no preimage causal. Evidência técnica ou mista
continua usando o RuleID pai. O teto documental permanece `PARTIAL`.

`REAL_OPACITY` não é inferida apenas da presença de campos que aceitam `UNKNOWN`.
`novelty=UNKNOWN` e `transfer_mode=UNKNOWN` isolados são resultados conservadores,
não opacidade do gate. `Unknown` requer identidade/classe/testemunha insuficiente,
baseline ausente, compatibilidade não medida, licença não reconhecida, budget
esgotado, disponibilidade não evidenciada ou outra incerteza R4 concreta. Uma
incerteza real permanece opaca quando coexiste com entradas resolvidas; toda
violação decidível tem precedência.

`R4.RUN_CANCEL.IDEMPOTENT_TERMINAL_TRANSITION` aceita versão esperada ou uma única
transição CAS para a versão seguinte em estado terminal. Avanço adicional é
violação.

Os predicados anteriores de pureza, provenance, status, novelty, canonical bytes,
hash/size/paginação, source resolution, admission, lifecycle, publicação, replay,
runtime, refinamento e tools permanecem ativos.

## API pública

```go
func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

O11 não emite receipt, seal, certificado ou veredito.
