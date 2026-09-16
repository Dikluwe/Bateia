# Passo 0069 — Exportar radar para Archify

## Intenção congelada

Adicionar a exportação opcional `bateia idea radar export --input radar-input.json
[--output radar.archify.json]`. Bateia continua autoridade semântica: recalcula
`AnalyzeIdeaRadar` e `PlanIdeaRadarProbes` e emite um único documento Archify Dataflow
`schema_version: 1`; não executa Node, Archify nem renderização e não altera 0066/0067.

Regime: protocolo completo limitado à autoria contratual L0, como ensaio procedural sem
atestação física de isolamento. Autor da intenção: usuário. Autor do contrato: este
executor, autorizado a escrever somente os nove artefatos L0 de 0069 e sem autoria de
implementação, oráculos, ataques, selo ou veredito.

Revisão contratual r1: evidência do renderer no commit Archify
`6db72a9aea3d0f67a6a034e41f8a5491476a11c1` invalidou r0. O layout Dataflow efetivo
oferece somente rows 0..4, centers de stage fixos e largura default de node 112; validar
apenas o JSON Schema não demonstra renderabilidade. A hipótese r1 é que limitar cada
stage a cinco nodes e abreviar labels elimina clipping/overlap estrutural do exemplar.

## Contrato observável

O futuro apresentador puro expõe:

```text
report.ExportIdeaRadarArchify(domain.IdeaRadarInput) ([]byte, error)
```

Ele chama explicitamente as duas autoridades L1. Se radar ou plano retornar `UNKNOWN`,
falha com a razão pública, sem documento parcial. Relatórios `COMPLETE` devem concordar
um-a-um em `TargetID` e outcome; divergência falha fechada.

O JSON contém somente propriedades admitidas pelos schemas públicos Archify Dataflow v1:

Versão externa congelada para esta materialização: commit Archify
`6db72a9aea3d0f67a6a034e41f8a5491476a11c1`. Compatibilidade futura com `main`
não é presumida.

- `schema_version: 1`, `diagram_type: "dataflow"`;
- `meta.title: "Bateia Idea Radar"`;
- stages, nesta ordem: `Origins`, `Perspectives`, `Targets`;
- nodes de origem, projeção e target; todos usam `type: "external"`, token visual sem
  alegação arquitetural, com `stage` 0, 1 ou 2 e `row` canônica;
- um flow origem→perspectiva por projeção válida;
- um flow perspectiva→target por observação canônica retida por 0066;
- um card por target, com outcome/reason e probes de 0067 em texto legível.

Probes nunca viram flows. Não criar relações entre targets, probes, evidência ou origens.
Origem `RECORDED` é deduplicada pelo `OriginID` opaco exato. `UNKNOWN_ORIGIN` recebe um
placeholder exclusivo por projeção, rotulado como desconhecido; nunca é coalescido nem
tratado como origem física. Flows de observação carregam ID e stance; variantes visuais
podem ser `default`/`security`/`dashed` para support/conflict/ambiguous.

IDs Archify são `prefixo + hex(bytes UTF-8 da identidade canônica)`: `origin_`,
`origin_unknown_`, `perspective_`, `target_`, `flow_projection_` e
`flow_observation_`. Prefixos separam namespaces e hex é injetivo, portanto todo ID
obedece `^[a-zA-Z][a-zA-Z0-9_-]*$` sem colisão ou sanitização lossy.

Labels de node não repetem a identidade integral. Defina `short(s)`: substitua cada byte
ASCII fora de `[A-Za-z0-9]` por `_`; se o resultado tiver até 8 bytes, conserve-o; caso
contrário use os primeiros 7 bytes seguidos de `~`. Labels exatas são `O ` + `short` do
OriginID (ou `O unknown`), `P ` + um de `PROB|BEHV|ARCH|CONS|OPER`, e `T ` + `short` do
TargetID. Assim cada label tem no máximo 10 bytes e cabe no node default 112. A identidade
completa permanece no ID hexadecimal; flows exibem projection/origin/observation completos
em seus labels, e card title/items preservam TargetID e provenance legíveis.

Ordenar nodes por stage e identidade, flows de projeção por projection ID, flows de
observação por observation ID, cards por TargetID e itens de probe pela ordem de 0067.
Serializar com indentação de dois espaços e newline final. Permutações equivalentes,
duplicatas semânticas e nil/vazio textual produzem bytes idênticos. O input e todas as
slices alcançáveis permanecem imutáveis.

## CLI, limites e falha fechada

`RunIdeaRadarExportCommand(args, stdout)` aceita exatamente `export`, `--input` obrigatório
e `--output` opcional. Decodifica um único JSON, rejeita campos desconhecidos/trailing
JSON e, sem `--output`, escreve somente o documento em stdout. Com `--output`, grava
atomicamente e deixa stdout vazio. Rejeita output que identifique o próprio input.
Nenhum arquivo é criado/truncado antes do cálculo completo.

Herdar todos os limites 0066/0067 e acrescentar o limite de renderabilidade antes da
materialização: no máximo 5 origens únicas (placeholder desconhecido conta), 5 perspectivas
e 5 targets. Logo: até 15 nodes, 1029 flows (5 projeções e 1024 observações), 5 cards,
1024 probes e 4 MiB serializados. Acima de qualquer limite, falhar fechado com
`render_budget_exceeded`, sem parcial. Rows são exatamente 0..N-1 em cada stage. O trabalho
é linear nas coleções limitadas; não há power set, rede, relógio, aleatoriedade, LLM,
filesystem no apresentador ou dependência Archify/Node.

## Materialização futura e gates

Novos arquivos propostos: `02_shell/report/idea_radar_archify.go`, seu `_test.go`,
`02_shell/cli/idea_radar_export.go` e seu `_test.go`; o dispatcher `idea` recebe somente
o encaminhamento aditivo necessário. Não editar arquivos 0066/0067.

Testes futuros cobrem schema/forma pública, golden, IDs adversariais sem colisão,
origem desconhecida, deduplicação, legibilidade de outcomes/probes, ausência de relações
inventadas, UNKNOWN de radar/plano, inconsistência, limites, 20 permutações, nil/vazio,
imutabilidade, stdout/output e proteção do input. Gates finais futuros: testes focais,
`go test ./...`, `go vet ./02_shell/...`, build, `crystalline-lint .` e
`git diff --check`. O exemplar congelado de `benchmark/0069` deve passar `archify validate`
e `archify deliver` usando exatamente o commit fixado; isso não impõe `showcase` a toda
exportação. Nesta rodada nenhum gate de implementação foi executado e nenhum
veredito foi produzido.
