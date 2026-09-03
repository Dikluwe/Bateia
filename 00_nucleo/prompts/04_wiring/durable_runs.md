# Durable recommendation wiring

Em L4, componha a via assíncrona do contrato 0026: raiz local, limites de store,
gerador criptográfico de run ID, store L3, fila/worker L3, o mesmo
`RecommendationService` certificado, serviço L1 de runs e servidor/CLI L2.

Valide workers `[1,4]`, capacidade `[1,32]` e relações de quota na inicialização.
Abra store e execute recovery de não terminais antes de iniciar workers ou aceitar
requests. No shutdown, pare reservas, cancele ativos, aguarde e feche recursos.
Injete tudo; não implemente hash, validação de chave/ID, transição, cursor,
paginação, quota, cancelamento ou regra semântica em L4.

A raiz e limites são configuração local, nunca argumentos MCP. Preserve a mesma
instância de `RecommendationService` para `bateia_recommend` e worker assíncrono.
Este prompt refina composição sem editar prompts certificados do Passo 0025.

## Refinamento aditivo 0027

Generalize o handler da mesma fila para despachar pelo `run_kind` persistido e
injete o módulo proprietário `market_discovery.go`. Não crie segunda store, fila,
root, idempotência, secret, recovery ou shutdown; preserve identidade de instância
e observáveis 0026 conforme o contrato 0027.
