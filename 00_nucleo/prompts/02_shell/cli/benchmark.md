# Prompt — CLI offline do benchmark ponta a ponta

Implemente `bateia benchmark validate|mutate|run` como uma fronteira L2 estritamente
offline. Os comandos devem consumir artefatos JSON congelados, produzir JSON
reprodutível e não construir clientes, ler credenciais ou acessar a rede.

`validate` verifica correspondência entre manifesto e corpus; `mutate` mede se os
ataques negativos são rejeitados; `run` executa o motor real e aceita filtro de
split. Entradas obrigatórias ausentes, JSON inválido e subcomandos desconhecidos
devem falhar explicitamente.
