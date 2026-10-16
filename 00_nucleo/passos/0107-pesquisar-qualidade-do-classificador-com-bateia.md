# 0107 Pesquisar qualidade do classificador com Bateia

**Data**: 2026-10-06  
**Estado**: PESQUISA EXPLORATÓRIA CONCLUÍDA

Usar a ideia de avaliação semântica como entrada do Bateia para pesquisar
software acadêmico (JOSS/zbMATH) e implementações públicas no GitHub.
Registrar consultas, origem manual de sementes, tentativas, falhas, candidatos
e critérios de seleção. Não apresentar os catálogos de software como busca
exaustiva de artigos ou os scores de triagem como qualidade científica.

Coletar commits/bytes de candidatos selecionados, com limites pequenos e
explícitos. Usar a avaliação lexical existente para localizar capacidades,
sem promover correspondência textual a comportamento comprovado. Confrontar
observações com documentação primária e implementação, conservando lacunas.
Não forçar metodologias de avaliação nas quatro classes de mecanismos do
catálogo. Exemplos, recibos e resultados ficam no Lab.

O produto desta pesquisa é uma comparação fundamentada e uma fatia de
experimento para o Bateia. Nenhum código de candidato será executado e nenhuma
política produtiva será alterada nesta investigação. Avaliação independente
do classificador continua exigindo referência revisada fora do próprio produtor.

## Sonda controlada da admissão

Aplicar pelo CLI insights a fonte congelada do passo 0106: baseline, reordenação,
duplicação de testemunha, remoção das avaliações, hash adulterado, regra fora
do catálogo, promoção apenas de status declarado e substituição apenas de
RuleIDs por outra classe catalogada. Os primeiros casos verificam invariantes
do contrato; os dois últimos caracterizam a confiança em avaliações recebidas,
sem fabricar uma taxa de acerto semântico ou referência independente.
Conservar entradas, saídas, exit codes e a distinção entre conformance e qualidade.

## Execução

O [relatório 0107](../relatorios/pesquisa-classificador-0107.md) registra a
descoberta real, falhas, fallback externo, avaliação e sonda de admissão.
Referência independente e qualidade semântica permanecem pendentes.
