# Catálogo ativo de regras de mecanismo

Materialize em L1 o catálogo v3 ativo por embed de compilação. Exponha uma função
que devolva os bytes e o SHA-256 canônico, criando um clone defensivo novo em cada
chamada. Mantenha o blob e o backing embutido privados.

Não leia filesystem, ambiente ou rede em runtime. Não faça parsing, decisão ou
composição MCP nessa função. Preserve byte a byte o catálogo congelado no passo
0078 e mantenha uma única autoridade produtiva sobre seu hash.
