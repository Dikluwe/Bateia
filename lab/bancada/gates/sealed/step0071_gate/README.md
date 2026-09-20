# Gate adversarial do passo 0071

Gate pré-candidato, derivado exclusivamente da spec congelada, dos três contratos
`*_step0071_contract_test.go`, de seus prompts e do baseline público enumerado.

Execute:

```sh
04_wiring/step0071_gate/run.sh
```

O runner primeiro confere `frozen-inputs.sha256` e depois executa a suíte. Antes da
implementação candidata, a compilação deve falhar pela ausência dos símbolos públicos
do passo 0071. Depois da implementação, cada mutante M01–M14 deve ser rejeitado; o gate
só aceita `14/14`, positivos preservados, negativos violados e opacos mantidos como
Unknown. O teste de determinismo executa vinte reordenações e compara JSON byte a byte.

