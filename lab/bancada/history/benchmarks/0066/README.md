# Experimento 0066 — Radar de convergência de ideias

Contrato autoritativo em
`../../passos/0066-prototipar-radar-de-convergencia-de-ideias.md`; recorte canônico em
`contract.md`. Estes artefatos nascem antes da implementação Go e dos oráculos.

O protótipo associa projeções explicitamente estruturadas e mede se três perspectivas de
três origens declaradas encontram o mesmo candidato/conceito. Não gera projeções, não usa
LLM/embeddings, não comprova implementação e não altera decisões.

`baseline.json` contém ground truth e predições isoladas apenas para o harness externo.
Produção nunca recebe esses rótulos. O gate compara o radar à melhor projeção isolada e
só declara `OBSERVED` sob os thresholds pré-fixados; a conclusão vale apenas para o
corpus de seis casos e três famílias.

`manifest.json` registra hashes e autoridades; `budget.json` fixa limites e condições de
parada. Autor de oráculos e adversário devem congelar suas saídas antes do implementador.

Estado: protótipo implementado e verificado em modo sombra. O oráculo independente
confirmou o corpus congelado, o ensaio pós-selo matou 14/14 mutações e a suíte completa
do repositório passou. O veredito `OBSERVED` significa somente ganho no corpus sintético
de seis casos; medição em projetos reais continua pendente. Protocolo completo como
ensaio, com segregação procedural e sem atestação física de isolamento.
