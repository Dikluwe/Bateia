# Evidence coverage diagnosis tests

Atacar a taxonomia de perdas com fixtures mínimas e determinísticas. Cobrir alvos
não planejados, falhas, budget, linguagem sem suporte, ausência de padrões,
ligação estrutural, conflito negativo, aliases e invariância de ordem. Um alias da
mesma regra nunca pode aumentar o denominador do recall.

Os testes não acessam rede, não alteram regras e devem falhar se uma causa mais
interpretativa esconder uma falha operacional anterior.
