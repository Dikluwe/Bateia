# 0121 — Ligar documentação a evidência estática

Intenção ex-ante, sob Tekt A0: nova capacidade pura no Lab 0121, consumindo
observador v8 sem alterá-lo. Localizar comentários Go Doc anexados a funções,
métodos e métodos declarados em interfaces. Conservar texto, bytes/hash,
posição física e dono sintático. Declarações não são garantias de execução.

Documentos externos exigem associação fornecida com intervalo de linhas e
localizador exato do dono AST. Validar origem dos bytes e dono; classificar
associação como fornecida, não descoberta automática. Não casar métodos pelo
nome para inferir implementação de interface. Não classificar significado,
suporte de uso, contradição ou atomicidade nesta fatia.

Ligar fatos emitidos pelo observador ao corpo do dono apenas por contenção
lexical e mesma identidade de fonte. Isso não estabelece concordância entre
texto e código, causalidade ou alcance. Preservar todos os resultados v8 dentro
da projeção; limites de entradas/fatos e parsing parcial explícitos.

Especificações precedem Go. Verificar associação errada, hashes, comentários
soltos, interface/implementação distintas, diretivas line, fonte incompleta,
reordenação e publicação exclusiva. Reproduzir em fontes conhecidas 0120 e conc
como diagnóstico, não avaliação cega nem precisão semântica independente.

Produto sem imports do Lab; lint, testes raiz, vet e build obrigatórios.

Materializado: função pura Inspect e CLI própria, sem alterar observador v8.
132 grupos Go Doc e três associações fornecidas em fontes conhecidas; texto e
fatos mantêm autoridades distintas. Interpretação semântica permanece fora da
fatia. [Laudo](../laudos/0121-documentacao-e-evidencia-estatica.md) e
[bancada](../../lab/bancada/concorrencia/0121/README.md).
