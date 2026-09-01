# Prompt espelhado — contrato adversarial L2 do Passo 0021

Exija uma CLI offline que prepare pacote cego e sele um oráculo a partir de JSON. Verifique
que a saída preparada não vaza braço, query, rank ou URL, que conteúdo de candidato com
`ignore previous instructions` permanece dado, que revisão pendente vira `UNKNOWN`, que
corpus vazio resulta em `INCONCLUSIVE_NO_CANDIDATES` e que duas execuções equivalentes
produzem bytes idênticos. Entradas inválidas não devem produzir artefato.

Não leia nem corrija implementação candidata.
