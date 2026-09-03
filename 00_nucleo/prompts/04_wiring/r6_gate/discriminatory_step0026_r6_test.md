# Prompt espelhado — variantes compiladas R6 do Passo 0026

Construa uma referência congelada que implemente as portas públicas substitutas do
Passo 0026 e uma única suíte black-box de conformidade. A suíte não recebe ID, categoria,
texto adversarial, expectativa ou seletor de comportamento. Ela deve executar controles
positivos de concorrência, reopen, cursor, cancelamento e legado, além das obrigações
negativas observáveis.

Materialize cada mutação sobrevivente como um pacote Go permanente e compilável, com
fonte distinta e `New(root)` próprio. Cada variante deve alterar em compile time a porta
chamada pela suíte comum; são proibidos campos ou argumentos de seleção runtime,
comparações por Mxx, `observeMutation`, tabelas de resultado e digest/texto como decisão.
A referência deve passar em processo separado e cada variante deve falhar em processo
separado, com o conjunto de subtestes black-box falhos como witness normalizado.

Para create, transition e publish percorra os cinco limites reais open, write, file-sync,
rename e dir-sync em diretórios temporários, usando OpenFile exclusivo, Write, Sync e
Rename. Para arquitetura, audite os imports da própria variante compilada. Execute os
goldens search, benchmark e recommend do predecessor certificado por uma lista de fontes
explícita que não compile candidato.

Execute as variantes nas ordens lexical, reversa e permutação determinística, três vezes
por ordem. Exija referência preservada, cada processo mutante rejeitado e o mesmo conjunto
de subtestes falhos nas nove rodadas. Registre SHA-256 da fonte e do diff de cada variante;
O01–O05 continuam Unknown e os 30 casos já causalmente validados permanecem creditados.
