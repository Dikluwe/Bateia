# Laudo 0119 — Referência comportamental do candidato

A bancada executou TestPool e ExamplePool originais de conc sobre bytes do
commit congelado em 0106. Fontes, testes e go.mod preservaram os hashes; módulo
temporário reconstruído não certifica checkout completo. Go 1.26.5 linux/amd64,
testify 1.8.1 e dependências efetivas registrados no recibo.

TestPool e seus subtestes passaram em três repetições com detector de race;
ExamplePool passou uma vez. São 14 nomes de testes/subtestes, mais um exemplo.
Não houve relato de race nessas execuções. Os asserts vieram do candidato e
não foram reescritos a partir de nossas classificações.

Quatro comportamentos coincidem com a revisão 0118: configuração antes de Go
funciona no exemplo; após Go provoca panic; após Go seguido de Wait também
provoca panic; outro lote depois de Wait completa suas tarefas. Reutilizar o
pool para executar tarefas e reconfigurá-lo continuam permissões distintas.

As expectativas antes apenas lidas agora possuem testemunhos de execução nesses
casos. Os registros 0118 permanecem históricos, com seus limites originais.
Não se promoveu resultado do modelo 0117 a defeito do candidato.

As quatro demais alegações não foram avaliadas pelos testes selecionados.
Suporte à configuração simultânea ao primeiro Go e atomicidade permanecem
inconclusivos; Wait antes de Go não foi executado. Aprovação do detector de
race nos casos selecionados não prova segurança universal ou atomicidade.

A origem dos testes é externa à investigação. A seleção e a interpretação
continuam autorais. Não se mediu precisão semântica do classificador nem se
produziu corpus representativo com revisão independente. Essa distinção evita
contar concordância em casos escolhidos como qualidade geral do produto.

crystalline-lint, testes completos do módulo raiz, go vet, build do produto e
diff check passaram. Assinatura Git foi desativada apenas no processo das
fixtures. Gates selados não reproduzidos. Produto sem imports do Lab; nenhuma
mudança no produto ou observador v8 nesta fatia.

[Comparação](../../lab/bancada/concorrencia/0119/results/comparison.json),
[execução](../../lab/bancada/concorrencia/0119/results/execution.json),
[verificação](../../lab/bancada/concorrencia/0119/results/verification.json)
e [reprodução/limites](../../lab/bancada/concorrencia/0119/README.md).
