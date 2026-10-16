# 0119 — Referência comportamental original de conc

Os testes originais `TestPool` e `ExamplePool` foram executados em módulo
temporário reconstruído dos bytes congelados de 0106, commit
`5f936abd7ae87036af1f75c95fb9d0daaf00116b`. Os 28 arquivos de fonte/material
fornecidos, incluindo go.mod e testes, tiveram hashes conferidos antes/depois.
Nenhum assert foi substituído. A coleta não certifica checkout completo.

Go 1.26.5 linux/amd64, `GOWORK=off`, `GOTOOLCHAIN=local`, dependências do go.mod
original (incluindo testify 1.8.1). A execução usou `-race -count=3` e passou:
14 nomes de testes/subtestes em três repetições, mais ExamplePool uma vez.
Go não repete exemplos com `-count`. O detector de race não relatou corrida
nessas execuções; isso não é garantia universal.

| Alegação do passo 0118 | Evidência agora executada |
| --- | --- |
| Configurar antes de Go é a ordem demonstrada | ExamplePool executou essa ordem com sucesso |
| Configurar depois de Go provoca panic | Teste original de configuração antes de Wait confirmou panic |
| Wait após Go permite reconfigurar | Teste original depois de Wait confirmou panic, contrariando a alegação |
| Outro lote depois de Wait é possível | Teste original completou dois lotes, contagens 10 e 20 |
| Configuração concorrente ao primeiro Go é suportada | Não avaliada pelos testes selecionados |
| Guarda anterior prova atomicidade | Não avaliada pelos testes selecionados |
| Modelo demonstra defeito em conc | Não avaliada pelos testes selecionados |
| tasks necessariamente indica goroutines iniciadas | Não avaliada pelos testes selecionados |

Quatro resultados comportamentais são consistentes com a revisão anterior.
Os outros quatro mantêm os julgamentos/limites de 0118, sem promoção. Nenhum
teste selecionado cobre configurar simultaneamente ao primeiro Go ou Wait
antes do primeiro Go. Testes de limite também rodaram como parte de TestPool,
mas a comparação não transforma seus casos finitos em garantia global.

## O que é independente

Os asserts e a implementação vêm do projeto candidato e precedem nossos
julgamentos. Essa é uma referência comportamental externa ao Bateia. A seleção
dos testes e a interpretação ainda são desta investigação; não houve revisão
humana independente nem corpus semântico representativo. Não há medida de
precisão do classificador: as oito alegações de 0118 são autorais e não saídas
automáticas classificadas por um oráculo independente.

## Reprodução e recibos

Na raiz do Bateia:

```sh
python3 lab/bancada/concorrencia/0119/run_reference.py
python3 lab/bancada/concorrencia/0119/compare.py
```

O primeiro comando cria módulo em `/tmp`, executa conc e pode baixar dependências
do go.mod. O segundo confere hashes/eventos e compara os quatro casos selecionados.
Não importa ou incorpora conc ao produto. Os comandos substituem resultados
da bancada atual, sem modificar os recibos anteriores 0106–0118.

`results/upstream-tests.jsonl` registra eventos brutos; `execution.json` registra
comando, ambiente, hashes e contagens; `go.sum` e `modules.jsonstream` registram
dependências; `comparison.json` distingue casos observados dos não avaliados;
`verification.json` registra os checks do Bateia. O produto e o observador v8
não receberam alterações nesta fatia.
