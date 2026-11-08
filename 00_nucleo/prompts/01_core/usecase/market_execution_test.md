# market_execution_test

Testar o executor de mercado configurado sendo chamado pelo runtime e os contratos de budget/cancelamento/artefatos. Sem I/O concreto, usar registry, store e conectores em memória; provar que a porta configurada não é ignorada.
Contrato 0131/ADR-0014.

A porta configurada também deve governar o resultado, não só receber a chamada: controlar um executor que abstém candidatos/retorna INCONCLUSIVE e exigir sua projeção sem reconstruir conclusões a partir das observações.
