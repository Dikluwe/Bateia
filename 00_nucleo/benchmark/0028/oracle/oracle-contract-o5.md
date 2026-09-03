# Contrato de oráculo black-box independente — Passo 0028 O5

**Revisão:** O5 vinculada ao contrato R3  
**Regime:** protocolo Tekt completo  
**Papel:** autoridade segregada de ORÁCULOS  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

O5 classifica somente pares formados por `contract_input` e
`mutant_contract_output` do protocolo público. O oráculo não recebe nem consulta
ID de caso, hash de ataque, commitment, descrição/target de mutação, configuração
do runner, expectativa, marcador, candidato ou veredito anterior.

## Entradas congeladas

| Entrada pública | SHA-256 |
|---|---|
| Contrato observável R3 | `f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b` |
| Catálogo `tekt.mechanism-rules.v2` | `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5` |
| Protocolo executável O5 | `455d6beec665765a0cca96b209ed5e651e7103c9a9dc5bebba51faa1b23b1fdb` |

Nenhuma outra entrada semântica foi lida. Hash identifica os bytes, mas não
atesta isolamento do workspace compartilhado.

## Fronteira e classificação

`Evaluate` faz decode estrito, aplica predicados contratuais globais e predicados
da porta indicada no próprio input. Não existe tabela por caso. O resultado é:

- `Violated`: ao menos um predicado decidível do R3 falhou;
- `Unknown`: nenhum predicado falhou e a própria construção contratual expõe
  opacidade real — identidade/evidence insuficiente, classe fora do catálogo,
  source histórica indisponível, target/baseline desconhecido, compatibilidade
  não medida, licença não reconhecida ou budget de contraevidência esgotado;
- `Preserved`: todos os predicados aplicáveis são decidíveis e satisfeitos.

`Unknown` não é fallback de parse, erro de infraestrutura ou ausência de
predicado. É excluído do numerador e denominador de mutation score. Uma mutação
negativa jamais é aceita como `Unknown` pelo simples fato de ser inesperada.

## Predicados globais

O output deve ser JSON estrito do protocolo, com `return_value` e `error`
mutuamente exclusivos, erro somente em pares código/mensagem R3 e arrays/objetos
explícitos. Rede, leitura/escrita de filesystem, subprocesso e mutação externa são
zero. Nenhum campo inseguro cruza publication. Cada componente de runtime tem no
máximo uma identidade observada e a identidade permanece estável na sequência.

`EvaluateSequence` exige saída canônica idêntica para o mesmo input canônico,
independentemente de repetição/ordem. O canonicalizador usa
`TEKT_CANONICAL_JSON_V1`; nenhum binding participa.

## Cobertura por porta

1. `MECHANISM_EXTRACT`: catálogo/hash exatos, content hash, join unívoco,
   identidade imutável, role/kind catalogados, source única, status e novidade,
   gaps, provenance, ordenação e ausência de inferência por conteúdo;
2. `SYNTHESIS_START`: schema/limites, allowlist de source, duplicatas, política,
   validação e erro pré-fila, idempotência, reserva/ID/create e receipt;
3. `SYNTHESIS_EXECUTE`: ausência de retorno inline, lifecycle sem efeito externo,
   cinco artifacts seguros e publicação atômica;
4. `RUN_READ`: run kind, lifecycle, três status nullable, progress e separação
   `SUCCEEDED`/status;
5. `ARTIFACT_READ`: enum/run binding, canonical bytes, SHA-256, size, cursor e
   erros controlados;
6. `RUN_CANCEL`: CAS/version, ausência de novo trabalho/publicação e erros
   controlados;
7. `TOOLS_LIST`: dez tools na ordem, nove nomes basais e tool 10 exata no fim;
8. `REPLAY`: catálogo/capture/hashes, zero rede, output determinístico e erro de
   integridade sem reparo;
9. `REFINEMENT_AUDIT`: receipts 0025–0027, runtime único e findings exatos de
   camadas.

## API pública

```go
type Classification string
const (
    Preserved Classification = "Preserved"
    Violated  Classification = "Violated"
    Unknown   Classification = "Unknown"
)

type ContractInput struct {
    Port    Port
    Request json.RawMessage
}

type ContractOutput struct {
    ReturnValue            json.RawMessage
    Error                  *ControlledError
    EffectObservation      EffectObservation
    PublicationObservation PublicationObservation
}

type Exchange struct {
    ContractInput  json.RawMessage
    ContractOutput json.RawMessage
}

type Finding struct {
    Classification Classification
    Port           Port
    Predicates     int
    Failures       []string
    Unknowns       []string
}

func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

O join posterior com records, IDs e hashes pertence exclusivamente ao
verificador. Esta autoridade não calcula mutation score e não emite receipt,
seal, certificado ou veredito.
