# 0133 — Prioridade em insights e primeiro diagnóstico do Bateia

Contrato ex-ante: [passo 0133](../../../../00_nucleo/passos/0133-priorizar-insights-de-engenharia.md). Resultado legível e limites: [laudo 0133](../../../../00_nucleo/laudos/0133-foco-em-insights-e-primeira-investigacao.md).

O usuário escolheu o próprio Bateia como primeiro alvo. Pergunta do ensaio: o relatório preserva condições e diferenças específicas de uma implementação real para orientar uma melhoria do alvo? O caso usa fontes congeladas e já conhecidas do passo 0106; a investigação é autoral e não representa holdout, avaliação semântica independente ou busca automática nova.

## Reprodução

Na raiz:

```sh
BATEIA_VERIFY_OUTPUT="$PWD/lab/bancada/insights/0133/final" GOCACHE=/tmp/bateia-review-go-cache make verify
python3 lab/bancada/insights/0106/verify.py --binary "$PWD/lab/bancada/insights/0133/final/bateia"
python3 lab/bancada/insights/0133/investigate.py --binary "$PWD/lab/bancada/insights/0133/final/bateia"
GOCACHE=/tmp/bateia-review-go-cache go test -v ./03_infra/runworker -run '^TestStep0026LocalQueue(ValidatesLimitsAndCommitCannotFailForCapacity|BoundsConcurrencyAndCancelsRunningContext|ConcurrentReservationsNeverOversubscribe)$' -count=1
```

O runner novo confere baseline de fontes Go, bindings dos testemunhos e projeções de campos; não faz interpretação semântica automática dos trechos. Executa apenas o Bateia. Os testes focados são da fila do Bateia, não do projeto consultado. O corpus 0106 permanece intacto; todos os resultados novos ficam nesta bancada.

## Artefatos

- `baseline.json`: documentos e arquivos Go antes da mudança de foco.
- `final/verification.json`: comandos/saídas da pilha completa.
- `replay-check.json`: reprodução atual dos controles 0106, sem sobrescrever seus recibos.
- `current-insight-result.json`: saída atual do comando; idêntica à saída congelada 0106.
- `source-witnesses.json`: cinco observações manuais anteriores, com trechos, intervalos e hashes conferidos.
- `target-witnesses.json`: código atual do alvo com intervalos e hashes, incluindo fila, extração e hipótese.
- `projection-observations.json`: comparação explícita da estratégia/pré-condições com os metadados do catálogo, estados e limites.
- `target-queue-tests.log`: três contratos existentes da fila executados neste ensaio.
- `verification.json`: consolidação do recorte documental e diagnóstico.

## Resultado observado

A estratégia projetada é a do catálogo, os campos de pré-condições reproduzem `language_limits`, e trade-offs/modelo de falha do mecanismo estão vazios. O relatório preserva PARTIAL/INCONCLUSIVE/UNKNOWN e não executa a fonte. A revisão compara isso com configuração de teto, submissão bloqueante e canal sem buffer na fonte congelada, e com reserva/rejeição/cancelamento cooperativo já existentes no alvo.

Essas observações sustentam o próximo contrato proposto: conservar condições específicas e sua proveniência no relatório, separadas da cobertura do extrator e da avaliação do alvo. Nenhuma mudança de semântica, formato, política ou algoritmo foi implementada aqui. Não se mediu ganho ao usuário ou efeito no alvo.
