# Verificação simultânea do harness de catálise

**Arquivo proprietário:** `01_core/experiment/catalysis/harness_test.go`

Materialize verificações locais mínimas, puras e determinísticas para o contrato
do harness do Passo 0030. Confirme que os estados observacionais têm grafia
estável, que o excesso do alvo editorial não invalida estruturalmente uma saída
JSON válida e que uma violação do protocolo de braço é exposta como `Violated`.

Este arquivo é um self-test do implementador e não substitui o oráculo externo,
independente e adversarial do Passo 0030.
