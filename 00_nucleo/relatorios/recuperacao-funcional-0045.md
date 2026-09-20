# Execução 0045 — Recuperação por perguntas funcionais

2026-09-07. [Passo e contrato](../passos/0045-recuperar-perguntas-funcionais.md). Implementado; validação limitada descrita abaixo.

## Mudança e motivo

O 0044 encontrava símbolos exatos, mas perguntas funcionais ainda selecionavam tipos e comentários associados às palavras da consulta. O 0045 pontua conjuntos de identificadores AST do nome e do corpo, priorizando funções pertinentes com corpo. Repetição de comentários e literais não aumenta a nota. O glossário PT/EN é finito, genérico e explícito; não contém nomes de projetos. A pontuação continua sendo inferência lexical.

A busca distribui sementes entre requisitos e expande até um salto de chamadas diretas, reservando no máximo um terço dos hits para expansão. Considera somente identificadores com destino único no diretório/pacote; exclui seletores, function values, closures e nomes sombreados. Ambiguidade e falhas de parse impedem relações no diretório. Não há resolução dinâmica ou typecheck completo.

Cada relação registra chamadora, linha da chamada e hash. A releitura confirma bytes e a existência da chamada AST na função informada; se a referência falha, ela é removida. Um resultado sustentado exclusivamente por chamada sem referência válida é descartado. O índice de chamadas da releitura é calculado uma vez por arquivo. Markdown mostra no máximo três referências por hit; o JSON conserva o inventário.

`analyze` e `idea recommend` local usam o ranking novo pela porta existente. O método é `go_ast_functional_v2`, o contrato `evidence-v7-functional-retrieval` e checkpoints locais registram `implementation_search: 0045`. Não houve alteração de detectores para aprovar capacidades.

## Avaliação independente de autoria

Autor B congelou um corpus de três projetos sintéticos em domínios distintos, com controle por símbolo exato, ausência e parser inválido. Métrica: arquivos esperados entre os primeiros três caminhos distintos de cada consulta, agregados apenas nos casos com alvo positivo.

| Medida no corpus | Baseline 0044 | Candidata 0045 |
| --- | --- | --- |
| Arquivos esperados recuperados | 3/4 (75%) | 4/4 (100%) |
| Precisão dos caminhos retornados | 3/10 (30%) | 4/4 (100%) |
| Controle exato e negativos | Passaram | Passaram |

O corpus, runner e esperados não mudaram para acomodar a candidata. Isso demonstra ganho nesta amostra pequena, não recall público ou compreensão semântica geral. [Avaliação B](../../lab/bancada/history/benchmarks/0045/evaluation-b.md) e [comparação](../../lab/bancada/history/benchmarks/0045/comparison.json).

Os ataques adicionais exercitam profundidade, orçamento, determinismo, chamada ambígua/sombreada e identidade das referências. A primeira revisão falhou: aceitava referência com hash correto na linha do `package`, fora da função chamadora. A correção passou a exigir a chamada AST no local declarado. Segunda execução passou sem alterar os testes. Depois, memoização da releitura e limite de referências exibidas foram incluídos; a suíte completa confirmou os testes novamente.

## Execução real no Bateia

A pergunta `como a busca preserva requisitos originais e rejeita capacidades sem evidência` agora extrai `idea_discovery.go` e `investigation_plan.go`, ausentes na execução 0044 dessa consulta. Entre as definições recuperadas estão `DiscoverExpansionPlanWithPolicy` e `BuildInvestigationPlan`. Ainda não selecionou `reuse_decision.go` nessa pergunta. Localizar esses arquivos é melhoria de recuperação, não demonstra que a resposta funcional está completa.

A consulta `SanitizeAssessmentEvidence` continua encontrando `evidence_authority.go:125–198` com identidade conferida. Ambas as varreduras examinaram 224 arquivos Go, 1.455.741 bytes, e conservaram 24 hits; estado parcial por orçamento. Esses números são cobertura/telemetria do Bateia, não benchmark do projeto.

A recomendação da ideia do Bateia manteve as dez capacidades INCONCLUSIVE. Foram 43 alvos coletados, 855 ignorados por orçamento e quatro não encontrados; a coleta ainda tem limitações. Replay reproduziu assessment, decisão e alternativas byte a byte. [Inventário, parâmetros, hashes e comparação de replay](../../lab/bancada/history/benchmarks/0045/observations.json); [relatório natural preservado](../../lab/bancada/history/benchmarks/0045/natural.md).

## Verificação e custódia

- `GOCACHE=/tmp/bateia-0045-cache go test ./...`: passou, incluindo gates existentes. A suíte usa servidores HTTP locais, executada com autorização de sandbox correspondente.
- `crystalline-lint .`: `✓ No violations found`.
- `GOCACHE=/tmp/bateia-0045-cache go build -o /tmp/bateia-0045 ./04_wiring`: passou.
- `git diff --check`: passou.

Regime Tekt **A/B sem atestação de isolamento físico**: autor do parser e autor dos testes separados; coordenador escreveu contrato e integração. Filesystem compartilhado. Corpus métrico congelado e executado antes da candidata; ataques adicionais congelados após liberação da implementação, sem acesso ao código candidato pelo autor B. A avaliação B registra exposição acidental limitada a nomes/linhas de testes antigos não relacionados durante busca de assinaturas.

Houve colisão do nome `baseline.json`: o baseline métrico substituiu o registro inicial de identidade. Os hashes anteriores de produção foram recuperados do registro 0044 já existente; a procedência está explícita em [contract-baseline.json](../../lab/bancada/history/benchmarks/0045/contract-baseline.json). Não se apresenta reconstrução de registro como selo prévio. O conteúdo do contrato não foi reescrito. Gates e certificados existentes não foram alterados.

## Limites e continuidade

A seleção permanece lexical, com glossário manual e relações AST delimitadas. Não interpreta negação como capacidade comprovada, não verifica execução e não generaliza para C/C++. Funções grandes ainda podem acumular termos pertinentes sem conter a resposta. O próximo ensaio deve ampliar o corpus com projetos reais e medir qualidade por requisito, especialmente os arquivos relevantes que ainda ficam fora do orçamento, antes de ampliar promessas. Os problemas prévios de gRPC/Protobuf e Bleve permanecem fora deste passo.
