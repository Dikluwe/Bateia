# Contrato de oráculo black-box independente — Passo 0028 O8

**Revisão:** O8 vinculada ao contrato R4 e catálogo v3  
**Regime:** protocolo Tekt completo  
**Papel:** autoridade segregada de ORÁCULOS  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

O8 recebe somente `contract_input` e `mutant_contract_output`. ID, binding,
descrição ou target de ataque, record, configuração executável, marker, expected
e candidato não são entradas do avaliador.

## Entradas causais públicas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo `tekt.mechanism-rules.v3` | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Invalidação R3/v2 → R4/v3 | `77e00c1eed39b972c6e415a45c208303e7af548740c8f7b68d78bb6c3ad33289` |
| Protocolo executável O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Invalidação do protocolo O5 → O8 | `33fc337b62de599278e90d25e5c0bd428201bc60ed47d48b0bfd1e7c0511c79f` |
| Manifest protegido predecessor O6 | `f48e81288a30ca220169ab4b733af1166be0f87d436cc7695302f416f997ff5e` |

## Classificação

- `Preserved`: todas as invariantes R4 aplicáveis estão satisfeitas.
- `Violated`: ao menos uma invariante decidível falhou.
- `Unknown`: nenhuma invariante decidível falhou e existe opacidade admitida.

Violação decidível precede `Unknown`. O join externo e o mutation score continuam
reservados ao verificador.

## Refinamento documental R4

O catálogo deve ser literalmente v3 e hash-bound. Cada classe possui exatamente
uma regra `DOCUMENTARY`, com RuleID R4 explícito, `witness_kind=DOCUMENTATION` e
`status_max=PARTIAL`. A source ligada ao bundle deve ter kind `documentation`.

Uma source documental catalogada, ligada e íntegra pode produzir mecanismo
`PARTIAL` sob status pai `PROVEN|PARTIAL`. Ela nunca conta para os mínimos
`IMPLEMENTATION`, `INDEPENDENT` ou dois hashes distintos de `PROVEN`; multiplicar
documentos não altera esse teto. Shape, role, kind, RuleID ou teto divergente
constituem violação ou gap conforme R4, nunca coerção.

As demais invariantes de extração e as outras oito portas preservam a calibração
O6 onde R4 não a substitui: safety, lifecycle, artifacts, start, cancel, replay,
tools e refinamentos continuam binding-free.

## API pública

```go
func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

O8 não emite receipt, seal, certificado, laudo ou veredito.
