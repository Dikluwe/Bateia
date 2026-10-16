# Bancada da separação de análise e execução

`baseline.json` registra o digest capturado antes da extração da análise pura de hipóteses. A fixture determinística e o verificador estão em `01_core/usecase/improvement_analysis_test.go`. O caso congela as cinco projeções, incluindo IDs e ordenação; não inclui estado de execução.

Reprodução: `go test ./01_core/usecase -run '^TestImprovementAnalysis' -count=1`. O teste de preservação compara o digest do resultado com o valor congelado. Cenários adicionais verificam degradação, ordem de mapas, não mutação do pedido e distinção entre estado de negócio e adoção.

Contrato: `00_nucleo/passos/0103-testar-hierarquia-de-arquiteturas-sob-tekt.md`. Esta bancada não altera ou reemite recibos dos gates selados.
