# Passo 0073 — Localizar a ocorrência probatória dentro do arquivo

## Estado e autoridade

- Estado: contrato candidato pré-implementação.
- Regime: materialização segregada Tekt completa, como ensaio com segregação
  procedural no mesmo filesystem.
- Papel desta autoria: autor do contrato e dos testes contratuais; não lê nem altera a
  implementação candidata futura ou o gate adversarial deste passo.
- Baseline Git: `03fd5f1`, acrescido dos artefatos públicos materializados pelo passo
  0072 e enumerados em **Entradas congeladas**.
- Escrita autorizada nesta fase: este passo, prompts proprietários e testes novos com
  sufixo `step0073_contract_test.go`.

## Problema congelado

O passo 0072 encontra de forma confiável `/usr/include/c++/13/span` e
`bits/basic_string.h`, mas a avaliação ainda pode citar a primeira linha que casa com
qualquer padrão. Assim, um comentário `@file basic_string.h` na linha 25 vence a
ocorrência substancial de `_M_limit` na linha 407; o mesmo ocorre com `span` antes das
declarações de `subspan` nas linhas 388/420.

O arquivo correto é recuperado, porém a fonte pública aponta para o trecho errado. Este
passo separa termos de recuperação de padrões probatórios e torna linha, trecho, padrão
casado e classe da correspondência observáveis.

## Política explícita de Unknown

- nome de arquivo, basename, stem e path são somente sinais de recuperação; nunca geram
  `EvidenceSource`, confiança, componente reutilizável ou promoção de status;
- ocorrência apenas em comentário, string literal ou documentação embutida em arquivo de
  código não é ocorrência probatória forte;
- substring de identificador não casa: `_M_limit` não casa `_M_limited`, e `subspan` não
  casa `subspanner`;
- macro, concatenação de tokens, declaração dividida ou sintaxe que o localizador lexical
  suportado não consegue atribuir com segurança permanece sem fonte forte; não se inventa
  parser C++ nem se promove a correspondência por aproximação;
- sem outra evidência admissível, esses casos seguem `INCONCLUSIVE`/`INSUFFICIENT_EVIDENCE`
  ou a causa Unknown anterior aplicável (`COLLECTION_INCOMPLETE`, `CATALOG_GAP` ou
  `HYPOTHESIS_ONLY`), nunca `ABSENT` por opacidade do localizador.

## Contrato público aditivo

Adicionar em `01_core/domain/evidence.go`:

```go
type EvidenceMatchClass string

const (
    EvidenceMatchStrongContent     EvidenceMatchClass = "STRONG_CONTENT"
    EvidenceMatchSupportingContent EvidenceMatchClass = "SUPPORTING_CONTENT"
    EvidenceMatchDocumentation     EvidenceMatchClass = "DOCUMENTATION_CONTENT"
    EvidenceMatchPathOnly          EvidenceMatchClass = "PATH_ONLY"
)
```

`EvidenceSource` ganha:

```go
MatchedPattern string             `json:"matched_pattern,omitempty"`
MatchClass     EvidenceMatchClass `json:"match_class,omitempty"`
```

Para toda fonte criada por regra de capacidade, `MatchedPattern`, `MatchClass`, `Line` e
`Snippet` devem referir-se à mesma ocorrência. `Line` é um índice baseado em 1 no conteúdo
imutável coletado. `Snippet` é a linha dessa ocorrência, normalizada apenas em espaços e
truncada pelo limite público já existente. `MatchedPattern` conserva o padrão canônico da
regra que efetivamente casou, e não o termo de path usado para recuperar o arquivo.

`PATH_ONLY` é valor público para diagnósticos de seleção, mas não é uma fonte admissível:
nenhum `CapabilityEvidence.Sources` ou `ReusableComponent.Evidence` pode contê-lo.

## Separação entre recuperação e prova

- `CapabilityRule.TargetHints` alimenta recuperação de caminhos.
- `CapabilityRule.StrongEvidence` e `SupportingEvidence` alimentam exclusivamente
  correspondências de conteúdo.
- `CapabilityTerms` continua servindo à identidade/alias da capacidade e não se torna
  prova por si só.
- `EvidenceRequirementQuery.Terms` contém somente termos de recuperação derivados de
  `TargetHints`; padrões `_M_limit`, `subspan` ou `size() - __pos` não vazam para a query
  salvo se também forem declarados explicitamente como `TargetHints`.
- a presença do basename no comentário `@file` não converte o termo de recuperação em
  padrão probatório.

## Localizador lexical C++ suportado

Sem prometer um parser C++, o recorte mínimo suportado reconhece em linhas de código:

1. identificadores completos de funções e membros seguidos de `(`, inclusive nomes com
   `_`, qualificadores, templates, `const` e `noexcept` posteriores;
2. expressões fortes declaradas literalmente pela regra, com espaços normalizados, como
   `size() - __pos`;
3. limites de identificador antes e depois do padrão;
4. comentários `//` e `/* ... */` bem formados e literais de string/caractere como regiões
   não probatórias.

Construções fora desse recorte não recebem interpretação semântica. O localizador deve
devolver a ocorrência textual suportada, não alegar verificação de comportamento.

## Precedência e canonicidade

A ordenação de fontes é total e independente da ordem de itens, regras e padrões:

```text
STRONG_CONTENT
  > SUPPORTING_CONTENT
  > DOCUMENTATION_CONTENT
  > PATH_ONLY (não emitido como fonte)
```

Dentro da mesma classe, desempatar por path lexical, linha crescente, padrão canônico e
snippet. Repetições legítimas podem ser preservadas, mas a fonte principal e o componente
reutilizável usam a primeira ocorrência segundo essa ordem. Duplicatas idênticas são
deduplicadas sem perder a linha.

## Escopo observável

Coberto:

- valores e tags JSON dos campos públicos novos;
- `_M_limit` citado na linha 407, e não `@file basic_string.h` na linha 25;
- `subspan` citado na primeira ocorrência substancial determinística (linha 388 no corpus),
  e não `@file span`;
- expressão `size() - __pos` citada na sua própria linha;
- funções/membros C++ sem keywords `func`, `type`, `class`, `def` ou `fn`;
- prioridade entre forte, apoio, documentação e path-only;
- repetição, reordenação e inputs imutáveis;
- rejeição de comentários falsos, strings, basename-only e substrings.

Fora de escopo:

- parser C++ completo, resolução de overload/template ou expansão de macros;
- prova comportamental, compilação ou execução do código citado;
- aumentar automaticamente a confiança além das regras do passo 0071;
- mudar recuperação, budgets ou cobertura de coleta do passo 0072;
- alterar catálogo, descoberta, licença, manutenção ou decisão de adoção.

## Corpus contratual pré-candidato

### Positivos

1. `basic_string.h` com `@file` na linha 25 e `_M_limit(...)` na linha 407 cita linha 407,
   snippet da declaração, padrão `_M_limit` e classe `STRONG_CONTENT`;
2. `span` com `@file` na linha 25 e `subspan(...)` nas linhas 388/420 escolhe de forma
   determinística a linha 388;
3. `if (__off > size() - __pos)` cita a linha da expressão e o padrão completo;
4. membro C++ com `const`/`noexcept` é localizado sem keyword de linguagem Go/Python/Rust.

### Negativos

1. `@file basic_string.h`, path `span` e basename-only não são fontes probatórias;
2. `_M_limit`/`subspan` apenas em comentário ou literal não geram fonte forte;
3. `_M_limited` e `subspanner` não casam;
4. `PATH_ONLY` nunca aparece em sources ou componentes.

### Opacos

1. macro/token pasting e ocorrência dividida que não cabe no recorte lexical permanecem
   Unknown sem fonte inventada;
2. coleta incompleta e catálogo opaco conservam suas causas anteriores.

### Determinismo

Executar vinte permutações de itens, regras, grupos e padrões. O assessment serializado
deve ser idêntico byte a byte, e ideia, catálogo e bundle permanecem imutáveis.

## Mutantes que o contrato deve eliminar

- M01: continuar usando `firstLine` com termos de path;
- M02: tratar `@file` como evidência substancial;
- M03: emitir basename/path como `EvidenceSource`;
- M04: não reconhecer função/membro C++ sem keyword Go/Python/Rust;
- M05: citar linha zero ou a primeira linha do arquivo;
- M06: preencher `MatchedPattern` com termo diferente da ocorrência citada;
- M07: aceitar substring de identificador;
- M08: aceitar comentário como strong content;
- M09: aceitar literal de string/caractere como strong content;
- M10: escolher apoio/documentação antes de strong content;
- M11: reordenação alterar fonte principal ou JSON;
- M12: deduplicação colapsar ocorrências em linhas distintas de forma não determinística;
- M13: macro opaca ser promovida como ocorrência suportada;
- M14: modificar ideia, catálogo ou bundle durante avaliação.

Mutation score exigido para selo: `14/14 = 1.0`. Unknown somente satisfaz casos opacos
explicitamente declarados.

## Entradas congeladas

- prompt humano que relata o arquivo correto e a linha incorreta;
- skill `tekt-materializacao-segregada` e suas referências;
- baseline Git `03fd5f19d0643d80818cad0334ff09279e91e17d`;
- artefatos públicos do passo 0072 presentes no worktree: passo, prompts, testes de
  contrato e os arquivos materializados `01_core/domain/evidence.go`,
  `01_core/usecase/evidence.go` e `03_infra/local/client.go`, identificados pelos hashes
  registrados abaixo.

## Gates reproduzíveis

Pré-candidato, devem falhar pelos campos/valores públicos e pelo localizador ainda não
materializados:

```sh
GOCACHE=/tmp/bateia-step0073-contract-cache go test ./01_core/domain -run '^TestStep0073' -count=1
GOCACHE=/tmp/bateia-step0073-contract-cache go test ./01_core/usecase -run '^TestStep0073' -count=1
```

Após candidato, repetir duas vezes, executar pacotes afetados, race e
`git diff --check`. Alteração posterior de qualquer entrada contratual invalida o selo.

## Registro pré-candidato

- Budget de calibração: duas revisões focais e no máximo três execuções do corpus.
- Autor do contrato: agente segregado `cpp_locator_contract`.
- Leitura autorizada: intenção congelada, skill, baseline e artefatos públicos até 0072.
- Escrita autorizada: apenas este passo, prompts e testes `step0073_contract_test.go`.
- Contexto proibido: implementação candidata e gate adversarial do passo 0073.
- Regime de atestação: segregação procedural; sem alegação de isolamento físico.

### Hashes congelados

| Artefato contratual 0073 | SHA-256 |
|---|---|
| `00_nucleo/prompts/01_core/domain/evidence_step0073_contract_test.md` | `dacffebb254e65153045a3694a035489511625a944e93f355b2602858175eeb4` |
| `00_nucleo/prompts/01_core/usecase/evidence_step0073_contract_test.md` | `1a06d521230773aa12f7ae95613fa8a52da5c17bdd721469f20535309c65916f` |
| `01_core/domain/evidence_step0073_contract_test.go` | `79d3aad764c26afb5a81bfd79297ad261225066a0bbefe172c03ced9264e0b55` |
| `01_core/usecase/evidence_step0073_contract_test.go` | `bbcf261ad882ea5b118c649f679a15b7da909bfa9b269d0a0214eedef1e0b33c` |

| Entrada pública materializada 0072 | SHA-256 |
|---|---|
| `00_nucleo/passos/0072-recuperar-evidencia-local-por-requisito-antes-da-amostra-geral.md` | `dc82ac35e6416d63e8c4b5836c1e1a2e97746e186e73e6a8439628583c2feb9f` |
| `00_nucleo/prompts/01_core/domain/evidence_step0072_contract_test.md` | `9eecabaafdcb1c61dc681e606ff5f04d1737f66196bca0df3884731f1365e28f` |
| `00_nucleo/prompts/01_core/usecase/evidence_step0072_contract_test.md` | `409efd859e3fcb574a8fa1558b135515f65d8ba73cb8e20c975d32bf674fdbf3` |
| `00_nucleo/prompts/01_core/usecase/evidence_step0072_implementation_test.md` | `45f55d7ba9d83c9227f2efcd823a603e70c45899af03ac6c2f47b17d9e153fde` |
| `00_nucleo/prompts/03_infra/local/client_step0072_contract_test.md` | `81219fc33b0f82ebb820d7d0ccb94054ef7005653cadf51f8d8eb7d0f1a5dba5` |
| `01_core/domain/evidence_step0072_contract_test.go` | `7fe1d59f3f690aefcaeaa07647642e316d7d21ec90166b782b34514dde4737a7` |
| `01_core/usecase/evidence_step0072_contract_test.go` | `844aaca8e22ff15de099afd30fc70cc19759c813ff26bd645e0a3e7f4420eefe` |
| `01_core/usecase/evidence_step0072_implementation_test.go` | `0c009ec3dc10d942287de2f6209c2af68cf21e7f8607dc0ee0119e376d1b6668` |
| `03_infra/local/client_step0072_contract_test.go` | `1a948fc484e2bd97e6978a5ad2d4a5ad07b1d9cb22e916638d3f8fa9ca083c33` |
| `01_core/domain/evidence.go` | `a158a9c7c757fa58ba5a275fefbf8e1eccdafcf281518731e7596042538932b9` |
| `01_core/usecase/evidence.go` | `c4e3bb761cb2e4560b5affddb688be80a085fe9fe5807ddd2614cf94513ca9f8` |
| `03_infra/local/client.go` | `fa0794a1b74562171ea0136d80538072ff48de7cf84813bba3caa2d11f69638a` |

### Resultado pré-candidato

Executado em 2026-09-16 contra o baseline congelado:

- domínio: falhou na compilação porque `EvidenceMatchClass`, seus quatro valores e os
  campos `MatchedPattern`/`MatchClass` ainda não existem;
- caso de uso: falhou na compilação pelos mesmos campos e valores ausentes, antes de poder
  observar o comportamento intrarquivo.

Esse vetor é a falha pré-candidato esperada. Não constitui selo, mutation score ou
veredito sobre implementação futura.
