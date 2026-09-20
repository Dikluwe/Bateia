# Contrato de oráculo black-box independente — Passo 0028 O6

**Revisão:** O6 vinculada ao contrato R3  
**Regime:** protocolo Tekt completo  
**Papel:** autoridade segregada de ORÁCULOS  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

O6 preserva a fronteira binding-free de O5: recebe somente `contract_input` e
`mutant_contract_output`. ID de caso, hash/descrição/target de ataque, record,
configuração adversarial, marker, expected e candidato não entram no avaliador.

## Entradas causais

| Entrada pública | SHA-256 |
|---|---|
| Contrato R3 | `f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b` |
| Catálogo v2 | `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5` |
| Protocolo executável O5 | `455d6beec665765a0cca96b209ed5e651e7103c9a9dc5bebba51faa1b23b1fdb` |
| Coverage gap sanitizado O5 | `bd93b949127bd80495af2840d9bd270c8aa36f6091086cd2c4550b40d2f396bb` |
| Falha discriminatória sanitizada O5 | `dc9f22846b8effe85e6fc97ebe0c9f27b50f567977ed5fcf865a58456fd4a96d` |

Também são inputs somente os artefatos O5 produzidos pela própria autoridade.
Nenhum artefato adversarial privado foi lido.

## Semântica de classificação

- `Preserved`: todas e somente as invariantes R3 aplicáveis à troca observada
  foram satisfeitas. Representações admitidas pelo contrato não são rejeitadas
  por convenções auxiliares do oráculo.
- `Violated`: ao menos uma invariante R3 decidível falhou. Uma porta não pode
  sobreviver apenas porque seu output tem shape válido.
- `Unknown`: nenhuma invariante decidível falhou e o input/output preserva uma
  opacidade expressamente admitida pelo R3. Opacidade não encobre violação
  simultânea e continua fora do mutation score.

## Calibração O6

Extração distingue invariantes contratuais de detalhes representacionais:
catálogo/hash, join, conteúdo, role/kind, source imutável, prova mínima, status,
novidade, provenance, safety e gaps permanecem obrigatórios; observações neutras
de identidade e formas canônicas permitidas não são efeitos nem falhas.

Artifact read valida request/enum, nome, bytes, size, SHA-256, documento compacto
com ordem fixa do schema do artifact, paginação/cursor e ausência de campos
inseguros. `TEKT_CANONICAL_JSON_V1` é usado para hashes semânticos do protocolo,
não para reordenar bytes de artifacts cujo contrato usa structs de ordem fixa.

Synthesis start valida integralmente textos, enums, limites, pares, refs,
duplicatas normalizadas e política. Falhas de source/hash acontecem antes de
idempotência/fila. Sucesso aceita somente os dois caminhos de admissão do R3:
reuso idempotente ou criação fresca ordenada.

Run read aplica exclusividade dos status, progress/lifecycle e separa
`SUCCEEDED` de `COMPLETE|SUPPORTED`. Cancel não executa admissão, worker nem
publicação e admite apenas as representações de retorno preservadas pelo
contrato predecessor.

Replay não confunde hashes de sources congeladas com hashes de artifacts
produzidos; exige integridade, zero rede e estabilidade em sequência.
Refinement reconhece receipt/layer/runtime ausente ou não demonstrado como
`Unknown` somente quando não há sobrealegação. Tools list conserva dez nomes na
ordem e exige a superfície exata demonstrável pelo R3, sem inventar baseline
não presente no input.

## API pública

```go
func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

O join dos findings com records permanece reservado ao verificador. O6 não
calcula mutation score e não emite receipt, seal, certificado ou veredito.
