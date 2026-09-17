# Passo 0072 — Recuperar evidência local por requisito antes da amostra geral

## Estado e autoridade

- Estado: contrato candidato pré-implementação.
- Regime: materialização segregada Tekt completa.
- Papel desta autoria: autor do contrato e dos oráculos observáveis; sem acesso a
  implementação candidata, que ainda não existe.
- Baseline: commit `03fd5f1` mais os arquivos
  públicos enumerados em **Entradas congeladas**.
- Alterações autorizadas nesta fase: este passo, os prompts proprietários dos testes e
  testes novos com sufixo `step0072_contract_test.go`. Implementação L1/L2/L3 existente
  e gates adversariais são somente leitura.
- Sucessor causal: implementação candidata somente depois do registro dos hashes desta
  fase.

## Problema congelado

Em árvores locais grandes, a amostra geral pode consumir `CodeFiles` antes de examinar
headers plausivelmente relacionados às regras. O adaptador local ainda pode cortar a
árvore silenciosamente em 20.000 entradas; o coletor rejeita headers C++ sem extensão,
como `<span>`; e desempates lexicais escolhem incidentalmente arquivos sem relação com o
requisito. Assim, `/usr/include/c++/13` pode omitir `span` e `bits/basic_string.h`, embora
um repositório pequeno com o mesmo conteúdo funcione.

O passo torna a recuperação de caminhos uma fase determinística, barata e orientada por
requisito, executada antes da amostra geral. A prioridade de caminho localiza candidatos;
ela nunca constitui evidência de implementação.

## Política explícita de Unknown

`COLLECTION_INCOMPLETE` é a causa pública de Unknown quando há caminho plausível não
examinado, busca de caminhos incompleta, quota do requisito esgotada, falha de leitura ou
superfície local truncada.

- `COLLECTION_INCOMPLETE` projeta obrigatoriamente `EvidenceInconclusive`, nunca
  `EvidenceAbsent`, e não entra em `CandidateAssessment.Gaps`.
- `NO_PLAUSIBLE_PATH` só é admissível quando a enumeração foi completa; ele não prova
  ausência da capacidade por si só.
- `TARGETED_EXAMINED` significa apenas que os caminhos candidatos selecionados foram
  lidos; o conteúdo ainda precisa satisfazer as regras de autoridade do passo 0071.
- correspondência do nome do arquivo não promove `PROVEN`, `PARTIAL` ou componente.
- ausência de suporte do provedor, budget esgotado ou erro não vira sucesso por default.
- catálogo `CATALOG_GAP` permanece anterior e distinto: sem regra segura não há consulta
  dirigida a inventar.

## Contrato público aditivo

Adicionar em `01_core/domain/evidence.go`:

```go
type RequirementCollectionStatus string

const (
    RequirementTargetedExamined RequirementCollectionStatus = "TARGETED_EXAMINED"
    RequirementNoPlausiblePath RequirementCollectionStatus = "NO_PLAUSIBLE_PATH"
    RequirementCollectionIncomplete RequirementCollectionStatus = "COLLECTION_INCOMPLETE"
)

type EvidenceRequirementQuery struct {
    RequirementKind string   `json:"requirement_kind"`
    Requirement     string   `json:"requirement"`
    RuleID          string   `json:"rule_id"`
    Terms           []string `json:"terms"`
}

type RequirementCollectionCoverage struct {
    RequirementKind string                      `json:"requirement_kind"`
    Requirement     string                      `json:"requirement"`
    RuleID          string                      `json:"rule_id"`
    Status          RequirementCollectionStatus `json:"status"`
    SearchComplete  bool                        `json:"search_complete"`
    CandidatePaths  []string                    `json:"candidate_paths"`
    ExaminedPaths   []string                    `json:"examined_paths"`
    SkippedPaths    []string                    `json:"skipped_paths"`
    ReasonCode      string                      `json:"reason_code"`
}
```

`ExtractionBudget` ganha `RequirementFiles int` com JSON `requirement_files`.
`EvidenceExtractionPlan` ganha `RequirementQueries []EvidenceRequirementQuery` com JSON
`requirement_queries`. `EvidenceBundle` ganha
`RequirementCoverage []RequirementCollectionCoverage` com JSON
`requirement_coverage`. Todos os arrays novos serializam como `[]`, nunca `null`.

`PlanEvidenceForIdea*` compila uma consulta para cada requisito `COVERED`, usando somente
a regra declarada correspondente. Os termos canônicos são derivados de
`CapabilityTerms`, `TargetHints` e padrões fortes/de apoio: tokens técnicos e basenames
plausíveis, sem strings vazias, duplicatas, termos genéricos ou identidade do candidato.
Consultas e termos são ordenados deterministicamente e não modificam a ideia nem o
catálogo.

## Recuperação anterior à amostra

`CollectEvidence` deve:

1. enumerar e normalizar caminhos antes de selecionar a amostra geral;
2. casar consultas por basename, stem ou token de identificador completo; substring
   interna não casa (`span` não casa `spanner`);
3. reservar até `RequirementFiles` arquivos distintos por requisito antes de gastar
   `CodeFiles` com a amostra geral; zero usa um default positivo e limitado;
4. desempatar por correspondência exata de basename, stem, menor profundidade e caminho
   lexical, nessa ordem;
5. reconhecer `.h`, `.hpp`, `.hh`, `.hxx`, `.c`, `.cc`, `.cpp`, `.cxx` e headers C++ sem
   extensão quando o basename completo for termo da consulta, incluindo `span`;
6. deduplicar um arquivo compartilhado sem cobrar leitura ou bytes duas vezes, mas
   registrá-lo na cobertura de todos os requisitos atendidos;
7. somente depois preencher a quota geral com os critérios legados.

A quota por requisito é reserva justa, não multiplicador ilimitado: todos os arquivos
continuam sujeitos ao orçamento total de bytes, tamanho máximo, confinamento e commit
imutável. Se esses limites impedirem uma leitura dirigida, o diagnóstico lista o caminho
em `skipped_paths` e usa `COLLECTION_INCOMPLETE`.

O adaptador local `FetchTree` não pode devolver uma árvore truncada como sucesso. Deve
enumerar toda a superfície acessível dentro das regras de confinamento ou retornar erro;
o chamador então marca `TreeComplete=false`. O antigo corte silencioso em 20.000 entradas
é proibido. Ordem do filesystem não pode alterar o conjunto ou a saída canônica.

## Avaliação e diagnóstico

Para cada consulta, a cobertura registra caminhos candidatos, examinados e omitidos em
ordem lexical, sem `null`. Um requisito com `SearchComplete=false`, `skipped_paths` não
vazio ou candidato plausível não lido é `COLLECTION_INCOMPLETE`.

`AssessEvidence` consulta essa cobertura antes de concluir ausência. Quando a cobertura
do requisito for incompleta e não houver evidência admissível positiva/conflitante, o
resultado é `INCONCLUSIVE`, confiança zero e `ReasonCode == "COLLECTION_INCOMPLETE"`.
Isso não mascara `CATALOG_GAP` ou `HYPOTHESIS_ONLY`, que conservam precedência semântica.

## Escopo observável

Coberto:

- valores, tags JSON, arrays não nulos e round-trip dos tipos novos;
- compilação determinística de consultas a partir de regras declaradas;
- prioridade de `span` e `bits/basic_string.h` antes de distratores;
- headers C++ extensionless relevantes, sem aceitar qualquer arquivo extensionless;
- quota justa por requisito e orçamento total preservado;
- diagnóstico de caminhos candidatos, examinados e omitidos;
- `COLLECTION_INCOMPLETE` preservado como Unknown;
- enumeração local acima do antigo corte de 20.000 entradas;
- ausência de regressão na coleta pequena e focada;
- repetição/reordenação com saída canônica e inputs imutáveis.

Fora de escopo:

- inferência semântica por nome de arquivo;
- aumentar confiança ou provar capacidade por path match;
- seguir symlinks externos, ler arquivos especiais ou remover limites de bytes;
- busca remota sem enumeração de árvore;
- alterar catálogo, descoberta, licença, manutenção ou políticas de adoção.

## Corpus contratual pré-candidato

### Positivos

1. árvore grande com distratores coleta `span` e `bits/basic_string.h` sob quota pequena;
2. duas capacidades recebem ao menos uma leitura dirigida cada antes da amostra geral;
3. árvore pequena Go conserva a coleta anterior;
4. árvore local com mais de 20.000 entradas contém o último arquivo na resposta.

### Negativos

1. `spanner`, `span.txt`, `basic_string.md` e arquivo extensionless sem termo não são
   promovidos como headers C++;
2. caminho localizado não produz `PROVEN` sem evidência de conteúdo admissível;
3. amostra geral não consome a reserva antes dos requisitos;
4. bytes e arquivos compartilhados não são cobrados duas vezes.

### Opacos

1. caminho plausível omitido por quota/bytes gera `COLLECTION_INCOMPLETE`;
2. erro/truncamento de árvore conserva `TreeComplete=false` e Unknown;
3. requisito sem regra permanece `CATALOG_GAP`, sem consulta sintética.

### Determinismo

Executar no mínimo vinte permutações da árvore, consultas e termos. Planos, cobertura e
paths coletados devem produzir JSON byte a byte idêntico; inputs permanecem imutáveis.

## Mutantes que o contrato deve eliminar

- M01: manter corte local silencioso em 20.000 entradas;
- M02: selecionar amostra geral antes das consultas;
- M03: tratar todo arquivo extensionless como código;
- M04: ignorar `span` por não possuir extensão;
- M05: substring interna faz `span` casar `spanner`;
- M06: uma capacidade consome toda a quota das demais;
- M07: path match promove capacidade;
- M08: budget omitido vira ausência;
- M09: diagnóstico omite caminhos pulados;
- M10: `COLLECTION_INCOMPLETE` entra em gaps do candidato;
- M11: ordem da árvore altera seleção;
- M12: deduplicação cobra o mesmo arquivo mais de uma vez;
- M13: catálogo sem regra gera consulta dinâmica declarada;
- M14: árvore pequena focada deixa de coletar arquivos legados.

Mutation score exigido para selo: `14/14 = 1.0`. Unknown somente satisfaz casos opacos
explicitamente declarados.

## Entradas congeladas

- este passo e seus prompts proprietários;
- baseline público de `01_core/domain/evidence.go`, `01_core/domain/interfaces.go`,
  `01_core/usecase/evidence.go`, `03_infra/local/client.go` e testes preexistentes;
- relato reproduzível: GSL pequena é coletável; `/usr/include/c++/13` omite headers
  relevantes sob amostra global.

## Gates reproduzíveis

Pré-candidato, devem falhar por símbolos/comportamentos ainda não materializados:

```sh
go test ./01_core/domain -run '^TestStep0072' -count=1
go test ./01_core/usecase -run '^TestStep0072' -count=1
go test ./03_infra/local -run '^TestStep0072' -count=1
```

Após candidato, repetir cada gate duas vezes, executar os pacotes afetados, race e
`git diff --check`. Qualquer alteração posterior destes artefatos invalida o selo.

## Registro pré-candidato

- Budget de calibração: duas revisões focais; no máximo três execuções do corpus completo.
- Autor do contrato: agente segregado `cpp_collection_contract`.
- Leitura autorizada: prompt humano, skill Tekt, baseline e artefatos públicos até 0071.
- Escrita autorizada: somente este passo, prompts e testes `step0072_contract_test.go`.
- Contexto proibido: implementação candidata e gate adversarial do passo 0072.
- Baseline Git: `03fd5f1`.
- Cache reproduzível dos gates: `GOCACHE=/tmp/bateia-step0072-contract-cache`.

### Hashes congelados

| Artefato | SHA-256 |
|---|---|
| `00_nucleo/prompts/01_core/domain/evidence_step0072_contract_test.md` | `9eecabaafdcb1c61dc681e606ff5f04d1737f66196bca0df3884731f1365e28f` |
| `00_nucleo/prompts/01_core/usecase/evidence_step0072_contract_test.md` | `409efd859e3fcb574a8fa1558b135515f65d8ba73cb8e20c975d32bf674fdbf3` |
| `00_nucleo/prompts/03_infra/local/client_step0072_contract_test.md` | `81219fc33b0f82ebb820d7d0ccb94054ef7005653cadf51f8d8eb7d0f1a5dba5` |
| `01_core/domain/evidence_step0072_contract_test.go` | `7fe1d59f3f690aefcaeaa07647642e316d7d21ec90166b782b34514dde4737a7` |
| `01_core/usecase/evidence_step0072_contract_test.go` | `844aaca8e22ff15de099afd30fc70cc19759c813ff26bd645e0a3e7f4420eefe` |
| `03_infra/local/client_step0072_contract_test.go` | `1a948fc484e2bd97e6978a5ad2d4a5ad07b1d9cb22e916638d3f8fa9ca083c33` |
| baseline `01_core/domain/evidence.go` | `a3387e44d328ef0a1e1239226fa16c3ecac2e136a293bab6af4ffe47dd230e8b` |
| baseline `01_core/domain/interfaces.go` | `8292e9a0b1425aabdad38fbf6a6c4b6964d146116e4b8ba5c1f3b982aee50ad4` |
| baseline `01_core/usecase/evidence.go` | `2659cf5db5fba30e5003064aa3371aca11ee3673d7665db4b035f6d1d9c56050` |
| baseline `03_infra/local/client.go` | `29b8176ab30bc07ec9ea818e7a2a6608fa0409b8d0255fa5b12b6268581bb012` |

### Resultado pré-candidato

Executados em 2026-09-16 contra `03fd5f1`:

- domínio: falhou por símbolos `Requirement*`, `EvidenceRequirementQuery` e campos de
  budget/plano ainda inexistentes;
- caso de uso: falhou pelos mesmos tipos/campos e por `RequirementCoverage` inexistente;
- adaptador local: compilou e rejeitou o baseline com `entries=20000 span=false`,
  testemunhando o corte silencioso.

Esse vetor é a falha pré-candidato esperada. Não constitui selo nem veredito sobre uma
implementação futura.
