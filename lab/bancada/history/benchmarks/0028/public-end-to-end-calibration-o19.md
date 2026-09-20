# Calibração pública end-to-end O19

Este corpus é uma projeção pública sanitizada do protocolo, não uma fixture,
mutação ou ataque. Ele não contém identidade de caso, hash externo, marker,
canário, expectation ou resultado incorporado ao input.

O replay deriva em memória os campos contratuais de integridade (`sha256`,
`size` e `canonical_json`) a partir de cada documento público, constrói os bytes
completos de `contract_input` e `contract_output` e chama o oracle binding-free:

```sh
GOCACHE=/tmp/bateia-go-cache go run 00_nucleo/benchmark/0028/public-end-to-end-calibration-o19-replay.go
```

Observação independente em 2026-09-02 contra O19: sete classificações
`Preserved`, zero `Unknown` e zero `Violated`. Esta observação é diagnóstico
público separado; não constitui expectation dentro do corpus.
