# Prompt espelhado — oráculo de composição L4 dos runs duráveis 0026

## Autoridade

Teste independente da composição observável, baseado no Passo 0026, contrato R1 e fake offline já certificado no Passo 0025.

## Obrigações

- Compor store, recovery, fila limitada, serviço de runs e MCP somente em L4 por fábrica exportada; nenhum nome ou detalhe privado é congelado.
- Injetar o mesmo `RecommendationService` usado pela via síncrona.
- Exercitar por stdio em memória `tools/list → start → get → get_artifact → cancel`.
- Confirmar `QUEUED/v1`, terminal `SUCCEEDED+INCONCLUSIVE`, artefato JSON seguro e terminal imutável ao cancelar.
- Aguardar o terminal por polling de 5 ms com deadline de 30 s, absorvendo contenção do gate completo sem aceitar estado não terminal.
- Fechar e reabrir a mesma raiz temporária, recuperando exatamente lifecycle, versão e status.

## Limites

Provedores são fakes offline, o store usa `t.TempDir`, e nenhum token, rede ou diretório externo entra na superfície MCP.
