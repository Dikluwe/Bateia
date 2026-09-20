# Execução do passo 0054

## Mudança e limite
A descoberta Sourcegraph usa uma única consulta integral em SearchCandidatesWithTrace e Discover. O filtro de linguagem permanece; o escopo de repositório da descoberta permanece. Discriminadores não substituem a pergunta. Consultas vazias e cancelamentos prévios não consomem chamada; erros não geram buscas amplas nem sucesso vazio. Rastros registram a consulta efetiva e os erros com redação de token.

A versão `evidence-v14-intent-preserving-discovery` invalida caches associados à versão anterior. O texto preservado continua sendo interpretado pelo provedor: não há alegação de tradução semântica de perguntas naturais. O passo remove deliberadamente o fallback defeituoso; não o substitui por um novo planejador. Ranking e extração local não foram alterados.

## Autoria e custódia
Contrato e prompts foram escritos antes da implementação. A implementou somente produção; B escreveu controles sem ler produção/patch e atualizou os contratos legados que exigiam o defeito anterior. Integração executou os gates. Regime A/B executado sem atestação de isolamento físico, sem selo formal de refinamento. Baseline e congelamentos estão em `../../lab/bancada/history/benchmarks/0054/`.

Corpus, runner r2 e expectativas do passo 0053 são as mesmas entradas para a comparação. Os arquivos de 0053 são conferidos por hash antes/depois. A descoberta pública e a análise local de snapshots são medições separadas; a jornada remota integrada e o comportamento executável dos projetos permanecem não avaliados.

## Verificação
Suíte Go completa passou, incluindo os controles HTTP independentes das duas portas. Compilação da CLI/auxiliar, linteiro (zero violações), nove controles do avaliador e comparação offline passaram. Candidato e testes mantêm os hashes congelados. A primeira execução do linteiro ocorreu durante a autoria do teste e mostrou um prompt temporariamente órfão; após materialização, a verificação final passou.

Os testes legados deixaram de exigir várias tentativas, substituição por discriminador, timeout duplicado e sucesso posterior a uma resposta vazia. A alteração é uma substituição explícita do contrato defeituoso, não evidência de equivalência com a estratégia anterior. Os controles novos cobrem termos tardios, negação/símbolos, linguagem/escopo, sucesso/vazio, falha HTTP/GraphQL/JSON, cancelamento e redação de token.

[Recibo e hashes](../../lab/bancada/history/benchmarks/0054/verification.json), [controles congelados](../../lab/bancada/history/benchmarks/0054/tests-freeze.json), [logs](../../lab/bancada/history/benchmarks/0054/full-tests.log).

## Resultado medido

| Medida | 0053 | 0054 |
|---|---:|---:|
| Projetos esperados na descoberta pública | 0/2 avaliáveis, 1 desconhecido | 0/3 avaliáveis |
| Tentativas por pergunta pública | até 3, com redução e retirada de linguagem | 1, pergunta integral com `lang:Go` |
| Funções-alvo recuperadas localmente | 1/6 | 1/6 |
| Arquivos-alvo extraídos | 5/6 | 5/6 |

As três consultas públicas retornaram zero candidatos, sem erro nesta coleta. Não apareceram listas provenientes de “Como criar” ou “Como obter”, pois essas consultas deixaram de ser emitidas. Não houve ganho demonstrado de recall. Todas as métricas locais permaneceram iguais; reavaliar os mesmos bytes offline reproduziu as métricas. O corpus e todos os artefatos 0053 mantiveram seus hashes.

[Métricas e tentativas](../../lab/bancada/history/benchmarks/0054/comparison.json), [respostas brutas e recibos](../../lab/bancada/history/benchmarks/0054/collected/receipts.json), [reprodução](../../lab/bancada/history/benchmarks/0054/README.md).

## Conclusão proporcional à evidência

**Fato observado:** a perda textual e a remoção de filtro desapareceram nas duas portas sob teste e nas três buscas públicas. A descoberta útil continua em zero neste conjunto.

**Inferência:** enviar a pergunta natural integral também é insuficiente nesta amostra; apenas eliminar o fallback não resolve a recuperação.

**Hipótese para o próximo passo:** construir consultas de código com ações, objetos e termos discriminantes, registrando transformações como hipóteses e avaliando os candidatos contra o pedido original. O teste de ganho deve manter perguntas e alvos congelados, além de reservar perguntas novas para evitar ajuste ao corpus.

**Não verificado:** compreensão de negações pelo provedor, melhoria geral em outros domínios, comportamento executável dos projetos e jornada remota integrada. Tempos registrados são telemetria da ferramenta/provedor.

