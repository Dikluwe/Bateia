# Regressões de despacho MCP

Reproduza três prepares com budget 1: somente a primeira deve executar.
Teste as 16 combinações das quatro extensões em duas ordens de registro,
orçamento comum entre labels e ferramentas legadas, tools/list sem consumo,
registro idempotente e inventário completo. Force linhas de 128 KiB com e sem LF:
o reader não pode consumir mais de 64 KiB+1 e nenhum handler executa.
Rejeite profundidade e quantidade de tokens antes dos handlers e prove recuperação
na próxima linha. Prove que contexto previamente cancelado não lê e que operações
duráveis recebem o contexto original mesmo com todas as extensões compostas.
