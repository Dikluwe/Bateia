# Passo 0057 — Inspeção pura dos bytes

InspectImplementationBehavior(report, files) retorna um registro por hit, sem alterar entradas. Exige definition, .go, arquivo único por path, tamanho igual quando positivo, SHA256 dos bytes e parse sem erros. Referência deve casar FuncDecl com corpo, nome simples usado na recuperação e linhas inicial/final exatas; linhas físicas ignoram diretivas //line. A combinação de nome e intervalo distingue receptores. Inválidos têm not_evaluated, motivo e zero fatos. AST percorre assignment (cada LHS formatado), increment_decrement e return em ordem fonte, sem entrar em FuncLit. Snippet são bytes exatos da instrução. Observed inclui limites explícitos de fluxo/execução, aliases/chamadas/efeitos indiretos, ausência de escrita não prova preservação e intenção não verificada. Não executar código, I/O ou alterar ranking.


## Adendo 0057 — Origem documental
Antes de qualquer inspeção, Retrieval=documentation retorna not_evaluated com Facts vazio e motivo explícito de origem documental, mesmo com bytes válidos. Outro hit de canal independente continua avaliável. Conforme 0057-adendo-origem-documental.md.
