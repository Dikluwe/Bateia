# Prompt — refinamento independente O9

Derive controles somente de R4, catálogo v3, protocolo O8 e diagnósticos públicos
sanitizados O8. Verifique genericamente que propostas não inventam evidence nem
cruzam namespaces de mecanismo/produto, que referências de source preservam
coerência interna antes da fila e que cancelamento terminal é idempotente enquanto
uma transição CAS avança somente uma versão. Reconheça estados opacos explícitos
em artifacts somente quando bytes fixos, hash, size, ordem e safety permanecem
válidos. Não use case ID, record, configuração adversarial, marker, expected ou
candidato.

## Invalidação O11

O11 substitui somente expectativas O9 refutadas publicamente: shapes com
`novelty|transfer_mode=UNKNOWN` sem incerteza causal são positivos, terminal
cancel pode representar uma única transição CAS e o preimage documental usa o
RuleID documental efetivo. O histórico O9 permanece congelado e não é regravado.
