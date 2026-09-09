# Resultado 0060 — Tradução de caminho na consulta local

A pergunta sobre interromper notificações não recuperava Remove mesmo com limite 10000. A tradução isolada caminho → path a colocou em 5º, com match somente em path no corpo. O problema diagnosticado é uma lacuna lexical da consulta; a função e seus bytes já estavam disponíveis. Fontes: benchmark/0060/original.json, path-control.json e after.json. Declaração fsnotify.go:333, SHA256 c721c84f2a2c6734ec0e82c0c9f4133ba3f62761457b52cc2a26aaf511e2c78a.

## Correção

Consultas naturais sem marcadores de negação recebem caminho/caminhos → path, registrada em Expansion com origem, regra 0060-v1:path:translation e status hypothesis. Símbolos, modos mixed e palavras em backticks não recebem essa regra. A consulta original e as restrições permanecem preservadas. Nenhuma alteração nos pesos, orçamento, indexação ou busca pública. EvidenceVersion v22 e checkpoint local/live 0060 invalidam reutilização da evidência anterior.

## Comparação no corpus congelado

| Alvo | Antes 0059 | Depois 0060 |
| --- | --- | --- |
| Add | documentação | documentação |
| Remove | ausente | código |
| BlendLab | documentação | documentação |
| LinearRgbToXyz | código | código |
| Intersection | código | código |
| Union | código | código |

Alvos recuperados: 5/6 → 6/6, sendo quatro por código e dois por documentação. Negativos: 0/3 indicações e 0/3 conclusões indevidas. Identidade e disponibilidade de bytes não comprovam comportamento. Fonte: benchmark/0060/comparison.json e collected/summary.json. O corpus, alvos e métricas não foram alterados.

## Limitação central

O mesmo termo path recupera Add, registros por caminho e Mkfifo; ele não distingue a ação de parar notificações. A consulta original tinha um resultado documental; após a correção há 24 resultados no limite normal. Portanto, a melhora medida é de recuperação do alvo, não de precisão semântica. Nenhuma dessas funções foi executada para verificar o pedido.

A amostra foi usada no diagnóstico e na calibração; não é holdout. Encontrar 6/6 alvos conhecidos não justifica declarar busca resolvida ou análise semântica completa. Add e BlendLab continuam originados da documentação. O próximo problema a avaliar é discriminação de ação entre candidatos do mesmo domínio, com perguntas novas e respostas esperadas congeladas antes de calibrar regras.

## Processo

Tekt A/B sem atestação física, ambiente compartilhado. A implementou sem ler testes/corpus/relatórios; B escreveu controles a partir do contrato sem ler produção candidata. Prompts anteriores ao código e hashes congelados antes da execução dos testes. Integrador conhece os alvos. Verificações e recibos em benchmark/0060. Não houve busca pública ou execução dos projetos. Tempos são telemetria do Bateia.

Validação final: build, suíte Go completa, controles independentes, nove testes do runner e crystalline-lint passaram. Replay reproduziu o resumo; hashes do corpus, produção e testes conferiram com os congelados. Não houve revisão do candidato ou dos testes após a execução.
