# Análise da jornada principal de reúso — passo 0039

Data: 2026-09-07. Diagnóstico executado no worktree atual, sem alterações de lógica de produção e sem alegação de avaliação independente.

## Resultado

A jornada principal mantém decisões finais conservadoras nos casos examinados, mas ainda transforma indícios lexicais em capacidades comprovadas e aceita fontes sem autoridade para certas conclusões. A saída não oferece evidência suficientemente confiável para uma decisão geral de reúso. Isso não demonstra ausência das capacidades do Bateia: demonstra limites e falsos positivos do seu avaliador.

O propósito adotado é receber uma ideia, descobrir alternativas e ajudar a decidir FORK, EXTEND, COMPOSE, INTEGRATE, REFERENCE, BUILD ou INCONCLUSIVE com evidência. A investigação concentrou-se em `RecommendationService` → `RunIdeaPipeline`. As melhorias recentes de Clang no `SearchFunnel` não demonstram, por si, melhoria nessa jornada.

## Execução e limites

Foram executadas cinco recomendações CLI locais: exportação Python, eventos Go, variante com negação, ideia pública do Bateia do passo 0012 e comparação dos dois candidatos sintéticos. Também foi executado um replay. Um adaptador controlado exercitou descoberta vazia e falha de provedor através do serviço real e do MCP real; cada cenário registrou duas requisições ao provedor.

Todos os comandos CLI concluíram com código 0. Os testes existentes de `01_core/usecase`, `02_shell/cli` e `02_shell/mcp` passaram, com resultados em cache. Isso confirma os contratos atualmente testados, não a qualidade semântica das conclusões abaixo.

Não foi medida descoberta pública, recall entre ecossistemas ou qualidade geral. Coin não foi candidato desta amostra; arquivos históricos sobre ele apareceram incidentalmente na autoanálise. Não foram executados benchmarks dos projetos. Tempos presentes nos logs são telemetria do Bateia. Os exemplos sintéticos foram inspecionados estaticamente; a existência de um teste não significa que o Bateia o executou.

Entradas, dossiês e logs estão no [pacote 0039](../benchmark/0039/README.md). O [resultado estruturado](../benchmark/0039/results.json) registra estados, caminhos coletados e comparações. O manifesto registra commit-base, hashes de fontes, binário e pacote; hashes identificam o estado, mas não substituem a preservação do checkout para recompilá-lo.

## Casos observados

| Caso | Resultado observado | Interpretação |
| --- | --- | --- |
| Go que apenas devolve a entrada | Capacidade assíncrona PROVEN; decisão INCONCLUSIVE | Falso positivo de capacidade; manutenção desconhecida impediu adoção |
| Exportação de notas em Python | INCONCLUSIVE por falta de regra | Avaliador admite limite, mas não cobre a capacidade apresentada |
| Ideia pública do Bateia | Nove capacidades sem regra; uma PARTIAL | Cobertura insuficiente para avaliar a própria proposta |
| Dois candidatos locais distintos | Mesmo ranking 0,6033 e alinhamentos 1 | Descrição derivada da consulta contamina triagem |
| Documento que nega upload | CONFLICT com o não objetivo `upload code` | Correspondência literal ignora negação |
| Provedor vazio ou com erro | INCONCLUSIVE no serviço e MCP | Proteção preservada; erro permanece registrado |
| Replay | Quatro artefatos iguais byte a byte | Reprodutibilidade observada nesse caso, excluindo telemetria |

## Achados e causas

### P1 — Nomes e presença de testes viram comprovação de comportamento

**Fato observado:** em `out-events/assessment-report.json`, `asynchronous event processing` recebe PROVEN. A implementação completa de `events.go` é `func EventProcessor(value string) string { return value }`; o teste apenas verifica esse eco. As fontes recebem força 0,9 pelo símbolo e 1,0 pelo teste. O avaliador ainda cria uma parte reutilizável.

**Causa inspecionada:** em [evidence.go](../../01_core/usecase/evidence.go), `DefaultEvidenceRules`, `assessCapability` e `strengthFor` associam padrões de nomes a forças suficientes para PROVEN. O limiar admite máximo ≥ 0,85. Não há demonstração do comportamento assíncrono nesse exemplo.

**Inferência:** a evidência pode contaminar comparação e recomendações de reúso. Não houve FORK indevido nesta execução: a manutenção INCONCLUSIVE bloqueou a decisão final. Esse bloqueio não corrige a comprovação falsa da capacidade.

### P1 — Arquivo de teste tratado como licença do projeto

**Fato observado:** `out-self/assessment-report.json` identifica MIT como PROVEN usando `01_core/parser/license_checker_test.go`. O trecho apresentado é uma referência de prompt, sem sustentar a licença do projeto.

**Causa inspecionada:** `classifyEvidencePath`, em [evidence.go](../../01_core/usecase/evidence.go), classifica nomes iniciados por `license` antes de reconhecer testes. `assessLicense` interpreta marcadores encontrados nesse conteúdo como licença.

**Inferência:** exemplos de um detector podem alimentar uma condição eliminatória de adoção. Este diagnóstico não determina qual é a licença real do Bateia. A correção deve exigir fonte de licença com autoridade e escopo explícitos, excluindo exemplos e testes do detector.

### P1 — A consulta é reapresentada como descrição do candidato local

**Fato observado:** no caso comparativo, exportador Python e eco Go recebem a mesma descrição `asynchronous event processing asynchronous event processing`, o mesmo ranking 0,603333 e sinais `description_alignment=1` e `capability_coverage=1`.

**Causa inspecionada:** `localCandidatesExcluding`, em [idea_recommend.go](../../02_shell/cli/idea_recommend.go), preenche `Description` com `idea.Problem` e `idea.Capabilities`. A triagem em [probabilistic_screening.go](../../01_core/usecase/probabilistic_screening.go) usa esses campos como sinais.

**Inferência:** a intenção do usuário retorna como aparente evidência do candidato. A inclusão explícita de um candidato local pode justificar examiná-lo, mas não justificar alinhamento observado igual a 1. A avaliação posterior distinguiu o Python como ABSENT para eventos; isso não torna a triagem correta.

### P2 — Catálogo e seleção limitam a generalidade

**Fato observado:** a capacidade de exportação Markdown recebe “nenhuma regra versionada cobre esta capacidade”. Na ideia pública do Bateia, nove das dez capacidades recebem a mesma explicação; a restante é PARTIAL. O catálogo padrão contém 14 regras.

**Fato observado:** a autoanálise coleta documentação de benchmarks, scripts e tipos, mas não coleta `01_core/usecase/recommendation.go`, `idea_run.go` nem `evidence.go`. A lista efetiva consta em `results.json`, campo `self_collected_paths`.

**Causa inspecionada:** `prioritizeEvidenceTree` usa relevância lexical booleana seguida de ordem alfabética. `capabilitiesForPath` aceita termos de quatro ou mais caracteres. O orçamento padrão separa código, documentação e testes, mas não assegura que as implementações centrais da pergunta sejam examinadas antes dos artefatos de validação.

**Inferência:** adicionar regras para cada exemplo isolado não basta. É necessário representar capacidade não mapeada, melhorar seleção por responsabilidade/implementação e distinguir limitação de cobertura de ausência funcional. Não se conclui que o Bateia não implemente as capacidades que seu avaliador deixou de comprovar.

### P2 — Negação e localização da fonte não são preservadas

**Fato observado:** a frase `This tool does not upload code.` é usada como CONFLICT para o não objetivo `upload code`. A fonte registra linha 1, embora a frase esteja na linha 3 do README sintético.

**Causa inspecionada:** `assessNonGoal` faz correspondência literal sem polaridade; `sourceFromItem` atribui linha 1. Ambos estão em [evidence.go](../../01_core/usecase/evidence.go).

**Inferência:** trecho e hash não tornam correta uma interpretação nem um intervalo. Negação deve ser considerada, e localização deve corresponder ao conteúdo efetivamente citado; quando desconhecida, não deve ser inventada.

### P2 — INCONCLUSIVE não explica suficientemente a ação seguinte

**Fato observado:** `out-events/decision.md` apresenta razão genérica e lista processamento assíncrono como lacuna, apesar de o assessment considerá-lo PROVEN. A causa bloqueadora observada é manutenção desconhecida. `alternatives.md` explica esse ponto melhor: “cobertura comprovada: 1/1; manutenção inconclusiva”.

**Causa inspecionada:** [reuse_decision.go](../../01_core/usecase/reuse_decision.go) atribui razão genérica e capacidades como lacunas quando há ambiguidade; a projeção principal em [idea_recommend.go](../../02_shell/cli/idea_recommend.go) repete essa estrutura.

**Inferência:** o dossiê deve separar lacuna de capacidade, evidência faltante e condição de adoção desconhecida, informando qual observação resolveria cada incerteza.

## Proteções confirmadas

O serviço e o MCP mantêm INCONCLUSIVE com zero candidatos e com erro de descoberta. A falha não é silenciosamente convertida em BUILD; os erros aparecem como falhas de discovery. O replay preservou `reuse-decision.json`, `assessment-report.json`, `alternatives.json` e `expansion-plan.json` byte a byte. Essas são evidências positivas limitadas aos casos executados.

## Próximo incremento recomendado

Priorizar o contrato de comprovação e autoridade na jornada principal: separar declaração, implementação localizada, teste localizado e comportamento verificado; impedir PROVEN apenas por nomes; aceitar licença apenas de fonte apropriada; retirar a consulta dos metadados pontuáveis do candidato. Validar contra os contraexemplos deste pacote antes de ampliar o catálogo.

Em seguida, melhorar a seleção das implementações, preservar negação e intervalos, e tornar INCONCLUSIVE acionável. O critério de avanço deve ser uma decisão de reúso mais justificável em domínios diferentes. Quantidade de símbolos extraídos ou sucesso num único repositório não mede esse resultado.

As correções acima são propostas, não implementações realizadas neste passo. A qualidade da descoberta pública e o ganho comparativo após futuras correções permanecem não verificados.
