# Resultado 0059 — Documentação preservada no orçamento de busca

A seleção agora reserva vagas para candidatos documentais exclusivos antes da expansão de chamadas. O limite total permanece igual; declarações exatas têm prioridade. A reserva nominal é max(1, limite/4), reduzida conforme disponibilidade documental, exatos e vaga mínima de código. Essa é uma política proposta, não um ótimo demonstrado.

## Diagnóstico observado
Com limite 10000, Add aparecia em 31º: 28 hits de código, seguido por documentação; Add era o terceiro documental. BlendLab aparecia em 79º: 76 hits de código, também terceiro documental. No limite normal de 24, nenhum recebia vaga. Fontes: benchmark/0059/fsnotify-p1-full-before.json e colorful-p1-full-before.json. O código em documentation_retrieval.go preenchia apenas vagas restantes, após corte em implementation_search.go.

Remove não aparecia mesmo sem corte na pergunta original (fsnotify-p2-before.json); não atribuímos esse caso ao orçamento.

## Resultado no corpus congelado 0053

| Caso | Alvo | Antes 0058 | Depois 0059 |
| --- | --- | --- | --- |
| fsnotify-p1 | Add | ausente | documental |
| fsnotify-p2 | Remove | ausente | ausente |
| colorful-p1 | BlendLab | ausente | documental |
| colorful-p2 | LinearRgbToXyz | código | código |
| bitset-p1 | Intersection | código | código |
| bitset-p2 | Union | código | código |

Alvos recuperados e identidade confirmada: 3/6 → 5/6. Os três alvos por código permanecem 3/6; os outros dois são indícios documentais. Três negativos: 0 indicações e 0 conclusões indevidas. Fonte: benchmark/0059/comparison.json e collected/summary.json. A métrica body_available indica bytes disponíveis, não implementação analisada.

A comparação inclui cada hit adicionado/removido. Em fsnotify-p1,24 hits de código/chamada tornam-se 20 mais 4 documentais. Em colorful-p1,24 tornam-se 18 mais 6 documentais. Também há substituições em outros positivos; nenhum dos três alvos de código foi perdido. Isso não demonstra ausência de regressões fora dessa amostra.

## Qualidade e limites da evidência
Add conserva comentário 278–313 e declaração 314 em fsnotify.go, SHA256 c721c84f2a2c6734ec0e82c0c9f4133ba3f62761457b52cc2a26aaf511e2c78a. Seu match documental é apenas file. BlendLab conserva comentário 733–734 e declaração 735–741 em colors.go, SHA256 a75c3edb5dd7e957c49694a6ac3c8901ec7aba2fb5dededca552510eb533e37c; matches b e l. São sinais fracos: a correção torna os candidatos visíveis, não melhora por si só a correspondência semântica desses termos.

Ambos mantêm behavior.status=not_evaluated e facts=[], conferidos no JSON final e preservados em documentary-behavior-check.json. Encontrar a declaração por comentário não prova comportamento nem satisfação da pergunta. Descoberta pública e execução dos projetos não foram avaliadas. Tempos são telemetria do Bateia.

## Processo de verificação
Contrato e prompts anteriores ao candidato; A implementou e B escreveu controles sem ler a produção candidata. Regime A/B sem atestação física, ambiente compartilhado; integrador conhece os alvos. Casos são calibração, não holdout.

O primeiro teste de relações falhou porque sua amostra de controle tinha 4 declarações contra 6 na amostra mista, alterando IDF. B corrigiu somente a amostra para manter todas as declarações e neutralizar o comentário. Teste, hash e falha r0 foram preservados; r1 congelado antes de execução. Produção não mudou nessa revisão. Logs e verificação final: benchmark/0059.

Lacunas restantes: recuperação de Remove e relação semântica entre perguntas e implementações, inclusive os dois novos candidatos documentais. O ganho deste passo é impedir exclusão sistemática de um canal de evidência; não é um ganho de comprovação funcional.

Validação final: build, suíte Go completa r1, seis controles novos, nove testes do runner e crystalline-lint passaram. Replay reproduziu o resumo; corpus, contrato, produção e testes r1 conferem com os hashes congelados.
