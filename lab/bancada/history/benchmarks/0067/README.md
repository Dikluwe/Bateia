# Experimento 0067 — Fronteira mínima do radar

O passo autoritativo é
`../../passos/0067-planejar-fronteira-minima-do-radar-de-ideias.md`; `contract.md`
congela seu fragmento observável. Estes artefatos precedem implementação e oráculos 0067.

O protótipo não muda a convergência. Ele converte o outcome atual do radar numa lista
ordenada de probes estruturais mínimos, sem executá-los nem inferir conteúdo textual.
`baseline.json` contém somente fixtures e efeitos simulados do harness externo; nada
disso entra no código produtivo.

A hipótese pré-registrada é reduzir probes sem perder capacidade de remover o bloqueador
atual: 11 probes contra 15 no baseline ingênuo, ambos 8/8, redução esperada de 4/15. O
controle C09 impede confundir remoção de conflito com convergência. Um futuro veredito
`OBSERVED` valerá somente para este corpus sintético e não demonstrará economia real.

Estado: protótipo implementado e verificado em modo sombra. Os oráculos independentes
confirmaram o corpus, 16/16 mutações reais foram eliminadas e a suíte completa do
repositório passou. O veredito `OBSERVED` vale somente para redução de probes no corpus
sintético com respostas simuladas; economia real ainda precisa ser medida. Protocolo
completo como ensaio, com segregação procedural e sem atestação física de isolamento.
