# Despacho MCP e limites de segurança — 2026-09-26

As extensões avançada, durável, scouting e síntese devem registrar handlers no
servidor, nunca executar chamadas dentro de `io.Reader`. Todas as ferramentas
passam uma única vez pelo orçamento `ToolCallBudget` e recebem o contexto de
`Serve`, inclusive leitura de artefatos. Registro repetido é idempotente; a ordem
de registro não pode esconder ferramentas. Mantenha enums, schemas, goldens,
validadores e erros públicos existentes. A lista base continua contendo três tools.

Leia o transporte original com scanner de capacidade limitada, antes de qualquer
decodificação ou extensão: mensagens de até 64 KiB, no máximo 32 níveis JSON e
8192 tokens JSON. Uma linha maior, com ou sem LF, encerra a sessão sem drenar
entrada ilimitada. Limites estruturais retornam erro de parse e permitem continuar
na próxima linha. Chamadas market malformadas reconhecíveis mantêm Invalid params,
mas agora consomem o orçamento comum. EOF sem LF continua válido.

Cancelamento prévio não lê o transporte. Operações recebem o contexto original
para interromper I/O e trabalho cooperativamente; nenhuma extensão substitui-o por
Background. Um `io.Reader` genérico bloqueado não pode ser interrompido de modo
portável sem suporte do próprio transporte; não crie goroutines abandonadas para
simular essa garantia.

## Supersessão explícita dos hashes de implementação

Os oráculos históricos fixavam o corpo do servidor em vez de apenas seu contrato.
A falha de orçamento exige alterar esses corpos. Esta revisão substitui somente
os hashes abaixo nos três testes históricos; comparações profundas de schemas e
goldens permanecem intactas, além de novas regressões de segurança.

| Arquivo | Hash anterior | Hash desta revisão |
| --- | --- | --- |
| server.go | e0c8605039740c895e733c8edc7cb7b6e6137a23b91ce93a1dbfa34fe18a0d42 | ec8fc8c4461d73984c29a73b8b5d65f5ffcb6c3a35b7ad907529b5fc613c4ff9 |
| server.md (prompt) | b5075f7cbb40a86a15d90ff690b410e4a47a8b3d806c5def3794dd8c4b76f564 | 70222330ecb281f58af0996eb9ec7690bf9f1e4fc48e8d665c0cbf708c080410 |
| durable_runs.go | 436989b07761406908297d733e8c07cb890c12aae92613cc541a232d8e1ec29e | ccd07f1e8832f02895be7274f88a65661537d332b58ff03ca83a6c7a0a1ac627 |

Isso é uma correção de segurança com evidência local, não recertificação histórica
nem afirmação de que o novo código corresponde aos snapshots anteriores.

## Projeção do registro no passo 0104

A lista pública de ferramentas deve usar o registro de descoberta injetado, inclusive quando síntese estiver habilitada. Projetar enum dos IDs configurados e limite de seleção de até dez. Preservar schema basal quando o registro contiver exatamente Apple/GitHub/Sourcegraph e comportamento legado quando não houver registro. O dispatcher continua único, com orçamento/contexto existentes; não adicionar parsing ou política de negócio.
