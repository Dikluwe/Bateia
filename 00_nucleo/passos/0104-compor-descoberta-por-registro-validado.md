# 0104 Compor descoberta por registro validado

**Estado**: IMPLEMENTADO E VERIFICADO LOCALMENTE  
**Data**: 2026-10-06  
**Relacionado**: ADR-0008, passos 0027/0028/0034/0103

## Problema e resultado

O construtor de síntese rejeita qualquer quantidade diferente de três conectores, apesar de a composição padrão oferecer cinco ou mais. O construtor de descoberta exige no mínimo três, enquanto o caso de uso seleciona um ou mais IDs por pedido. A apresentação MCP também conserva enum estático mesmo ao receber um registro ampliado.

L1 será proprietário da admissão do registro configurado: não vazio, capacidades válidas e IDs únicos. L4 delegará essa admissão e conservará apenas composição e lifecycle. L2 projetará os IDs configurados, sem inferir cobertura ou comprovação pelo número de fontes. A quantidade de conectores registrada não limita o catálogo; o contrato de pedidos mantém seleção máxima de dez.

## Preservação

Assinaturas dos construtores L4, store/fila compartilhados, ordem Recover antes de Start, shutdown, catálogo fechado de mecanismos e erros públicos permanecem. A composição basal Apple/GitHub/Sourcegraph conserva o schema original por sua identidade completa. A lista de tools e suas posições não mudam. Novas composições unitárias e ampliadas passam a ser admitidas explicitamente; vazio/duplicado/inválido continuam recusados.

## Verificação

Testes puros de admissão em L1 e composição real de ambos os servidores em L4, com leitura de tools/list e list_connectors. Asserções de enum por identidade, limite de dez mesmo com onze registrados, consulta offline sem Discover, shutdown/reabertura e rejeição antes de criar a raiz para inputs inválidos. Reproduzir a falha antes da correção e executar lint, suíte completa, vet e build.

## Fora de escopo e falhas

Não alterar provedores, política de reúso/insights, schemas sem relação com registro ou evidência histórica. Não interpretar construtor válido como disponibilidade remota das fontes. O marco falha se o schema anunciar fontes diferentes das registradas, se recursos forem criados para registro inválido ou se a composição basal mudar seu contrato.


## Execução

A [síntese de execução](../laudos/0104-composicao-por-registro-validado.md) registra regressões reproduzidas, composição ampliada, preservação das capturas basais e verificações. A bancada de startup fica em `lab/bancada/arquitetura/0104`. O teste de composição inclui o conjunto padrão real e Bitbucket configurado sem consultar as APIs remotas.
