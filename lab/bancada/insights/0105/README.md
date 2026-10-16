# Bancada de insights independentes

O exemplo `queue-documentation-input.json` pergunta como uma fila limitada poderia orientar o tratamento de sobrecarga no Bateia. A fonte é uma fixture documental ilustrativa; owner/repo e SHA são identidades de teste, sem afirmação de origem externa autenticada. A avaliação PARTIAL e a licença declarada não foram verificadas remotamente.

A execução gera um mecanismo PARTIAL e hipótese INCONCLUSIVE com alvo UNKNOWN. O plano pede baseline, aplicabilidade e experimento. Nenhuma decisão de reúso, adoção de código ou efeito real é produzida.

Reprodução: `go build -o /tmp/bateia-insights ./04_wiring`, seguido de `/tmp/bateia-insights idea insights analyze --input lab/bancada/insights/0105/queue-documentation-input.json`. A saída registrada está em `queue-documentation-result.json`.

`offline-check.json` registra o startup em diretório temporário, com HOME ausente e raiz de runtime inválida. O comando passou sem criar home/cache ou runtime, apesar de haver `.env` no diretório. Foram fornecidas somente evidências congeladas; nenhum provedor ou candidato foi executado.
