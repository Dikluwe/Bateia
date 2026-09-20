# Contrato de oráculos independentes — Passo 0028

**Revisão:** O4 vinculada ao contrato R3  
**Papel:** autoridade segregada de ORÁCULOS  
**Implementação candidata lida:** não  
**Casos adversariais lidos:** não  
**Atestação máxima possível:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

Este artefato congela expectativas independentes sem conter ataques, receipts,
seals, certificado ou veredito. O gate decide apenas as relações observáveis aqui
delimitadas; não decide adoção nem equivalência geral.

## 1. Entradas causais

| Entrada | SHA-256 |
|---|---|
| Passo 0028 | `9105b27e0e8e9165b0b133a647669e858b977ea2ab8c7d951ed0f217a26ac1a5` |
| Manifesto 0028 | `0b8fbcebf9f3af491ae1851f59715c248fa63181617cc336dc39b6969fb555ad` |
| Contrato R3 | `f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b` |
| Catálogo de regras v2 | `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5` |
| ADR-0002 | `b457fadd55d6328c2afce969b92a527b461edab9094b1bef66fcfebe71b79a68` |
| Feedback de cobertura O1 | `7eca960c91c865e2a24c697c2568c3a99c11f90a597be16848d9c28f8f7d00e1` |
| Falha discriminatória O1 | `d85a18daa25d523f2dca9d5225e2789929a375181032a089fcca302a578a67ec` |
| Feedback de cobertura O2 | `92da449660fdd1f0fba1906cfe9077168820d1e838c671633c33fd0c83c69f3b` |
| Falha discriminatória O2 | `7652afa4080d3ff273d499c56d508bfd3df65a8cf0a0b0b6406e77fd6e82edd4` |
| Feedback de cobertura O3 | `f97c12b13b36d0aa741cbe2db3910d3b6b671f03bb7cbad5a14e644b1e912876` |
| Falha discriminatória O3 | `2a65e760238afb1f0d01ada36555007a986a1ecc8673f48bd959e791fb948646` |
| Commitments públicos O4 | `0408fa27c802444c3131c138c7e09500b802f2109ba52a5826c3e1bd834d14e5` |
| Raiz das 182 entradas comprometidas | `a486a67b330946daea5d9bfb41f7764ce35a1e300f94bb77210341d3893a4a3a` |

Hash identifica bytes, mas não prova isolamento; o workspace é compartilhado.
Qualquer alteração nessas entradas ou no split protegido invalida O4.

## 2. Invalidação O3 e regime O4

O3 foi invalidado antes de selo: seu `VariantSubject` executava o subject de
referência e depois sobrescrevia um campo. Isso era uma diferença pós-observação,
não uma mutação comportamental da porta, e resultou em 0/160 mutações
qualificadas. Os arquivos executores O3 permanecem append-only, mas têm build tag
`step0028_o3_invalidated` e não participam da compilação normal.

O4 usa o protocolo Tekt completo. A autoridade leu somente as entradas públicas
allowlisted acima. Não leu implementação candidata, casos adversariais, ataques,
receipts, seals, certificado ou veredito. Escreveu somente contrato de oráculo,
API/harness independente, testes, prompts, matriz, autoria, invalidação e
manifesto do split protegido.

Cada `M001–M160` tem um `semantic_variant_id` e um hash de entrada canônica
fornecidos pelo commitment público. `O4SemanticVariantSubject` transforma uma
entrada tipada antes da Observation e a entrega a uma implementação independente
da porta. Ele não chama `ReferenceSubject`, não lê `Case.Expected`, não seleciona
expectativa por ID e não injeta ID/hash na Observation. A suite valida os
predicados normativos e somente os campos requeridos da entrada comprometida.

O gate exige, para cada M, baseline `Preserved`, binding exato do hash/semantic
ID, hash da entrada mutante distinto, stage `INPUT_BEFORE_OBSERVATION` e
exatamente um required observable divergente. A execução ocorre nas ordens
normal, reversa e rotacionada. Os 160 hashes de entrada mutante são distintos.

`P01–P12` executam subjects de referência reais e devem ser `Preserved`.
`O01–O10` executam uma porta opaca explícita e devem ser `Unknown`; nunca entram
no numerador ou denominador de mutation score.

## 3. Domínio observável

`Input` e `Observation` representam identidade `owner/repo@commit/path`,
linguagem/ecossistema, hashes e kinds de source, RuleIDs e roles, limites, IDs e
colisão, proposal strictness, paginação/snapshot, state machine e hashes de
artifact, lifecycle/status, cancelamento/recuperação, replay, identidades de
runtime, schemas, tools, refinamentos 0025–0027 e findings de camada.

As portas são `EXTRACT`, `HYPOTHESIS`, `START`, `PAGINATION`, `RUNTIME`,
`LIFECYCLE`, `REPLAY`, `SAFETY`, `TOOLS`, `SCHEMA`, `ARTIFACT`, `REFINEMENT` e
`FULL_REFINEMENT`. Uma operation não ligada falha; não recebe `Preserved`.

Os casos congelados cobrem extração pura, witness cross-repo/cross-commit,
docs-only, produto-only, target `UNKNOWN`, licença `UNRECOGNIZED`,
popularidade/novidade, composição, contraevidência, source hash e validação
pré-idempotência/pré-fila, normalização/hash/idempotency, lifecycle/status,
cancel/recovery, replay offline, runtime único, publicação/paginação/schemas,
refinamentos e preservação exata das nove tools com a tool 10 no final.

## 4. Regras normativas essenciais

Extração aceita somente join exato `(kind,ref,path,content_hash)` entre source
avaliada e evidence coletada. `EvidenceSource.RuleID` liga literalmente o
catálogo v2; conteúdo bruto serve apenas para recomputar hash. `PROVEN` requer
implementação e witness independente na mesma source imutável, sob pai
`PROVEN`, com hashes distintos e sem contraevidência central. Docs-only tem teto
`PARTIAL`; produto-only não vira mecanismo; binding ambíguo, regra desconhecida,
kind incompatível ou hash divergente gera gap controlado.

Novidade é posterior ao status e não o altera. Popularidade não é suporte.
Target `UNKNOWN` não vira ausência ou suporte. Licença `UNRECOGNIZED` permanece
visível e bloqueia reuse automático. Composição sem witness mantém cinco
compatibilidades desconhecidas. Contraevidência e escopo de busca permanecem
visíveis.

Source refs e canonical document são validados no store antes de
`ResolveIdempotency`, reserva, ID, criação ou fila. Replay usa apenas captures,
catálogo e artifacts congelados, faz zero chamadas de rede e preserva bytes,
IDs, hashes e índices sob repetição/reordenação semântica.

`tools/list` preserva deep-equal as nove tools basais certificadas e acrescenta
`bateia_start_improvement_synthesis` somente na posição final. A captura basal
real é vinculada ao SHA-256
`ee630cff5eada0dee4afe82ad9513dea7b760f0beb7905e915d2714a9114a596`;
captura externa não vinculada não satisfaz o gate.

## 5. API pública permitida ao implementador

As APIs de produção são exatamente as do contrato R3:

```go
type MechanismExtractionInput struct {
    Bundle               EvidenceBundle
    Assessments          []CandidateAssessment
    ConventionalBaseline *ConventionalBaseline
    RuleCatalogJSON      []byte
    RuleCatalogSHA256    string
}

type MechanismExtractionResult struct {
    Mechanisms []ArchitecturalMechanism `json:"mechanisms"`
    Gaps       []MechanismGap           `json:"gaps"`
}

type MechanismExtractor interface {
    ExtractMechanisms(context.Context, MechanismExtractionInput) (MechanismExtractionResult, error)
}

type ImprovementSynthesisStarter interface {
    StartImprovementSynthesis(context.Context, string, ImprovementSynthesisRequest) (ImprovementSynthesisStartReceipt, error)
}

type ImprovementSynthesisExecutor interface {
    ExecuteImprovementSynthesis(context.Context, RunID) error
}
```

Tipos, constantes e tags JSON do contrato R3, inclusive
`ImprovementLicenseAssessment`, permanecem exatos.

A superfície pública do harness independente é:

```go
type Subject interface { Observe(context.Context, Case) (Observation, error) }
type ObservePort func(context.Context, Case) (Observation, error)

type Ports struct {
    Extract, Start, Lifecycle, Replay, Runtime ObservePort
    Tools, Refinement, Hypothesis, Safety      ObservePort
    Artifact, Pagination, Schema, Full         ObservePort
}

type O4ExecutableSubject interface {
    Execute(context.Context, Case, O4Commitment) (O4Execution, error)
}

type O4ReferenceSubject struct{}
type O4SemanticVariantSubject struct{}
type O4OpaqueSubject struct{}

func DecodeCases([]byte) ([]Case, error)
func DecodeO4Commitments([]byte) (O4CommitmentDocument, error)
func O4EntryCommitmentRoot([]O4Commitment) (string, error)
func O4PositiveCase(O4Commitment) (Case, error)
func O4NegativeCase(O4Commitment) (Case, error)
func O4OpaqueCase(O4Commitment) (Case, error)
func EvaluateO4Negative(O4Commitment, Observation, O4Execution, O4Execution) O4PredicateFinding
func EvaluateO4Positive(O4Commitment, Observation, O4Execution) O4PredicateFinding
func EvaluateO4Opaque(O4Commitment, O4Execution) O4PredicateFinding
```

`Case`, `Input`, `Observation`, `Operation`, `Obligation`, `O4Predicate`,
`O4Commitment`, `O4CommitmentDocument`, `O4Execution` e `O4PredicateFinding` são
serializáveis. Rótulos TEST e valores `Expected` permanecem protegidos.

## 6. Execução reservada

O verificador deve validar JSON/hashes, executar referência e cada mutante duas
vezes, ordem embaralhada, race, vet, build, diff, scan anti-tautologia, replay e
gates de refinamento. Somente o verificador calcula mutation score ou emite
receipt, seal, certificado ou veredito.
