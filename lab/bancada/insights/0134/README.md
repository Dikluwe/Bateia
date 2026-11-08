# 0134 — Observações específicas no relatório de insights

Caso curado a partir do diagnóstico 0133 e dos bytes congelados no ensaio 0106.
Fonte: `sourcegraph/conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b`.
O pedido v2 conserva o bundle e as avaliações anteriores, acrescentando seis
observações manuais com onze referências. Sem nova consulta remota, execução
da fonte ou reavaliação automática da interpretação.

## O que o resultado conserva

- Limite de workers exige configurar `WithMaxGoroutines` antes do primeiro uso;
  o padrão permite concorrência sem teto.
- Na saturação, a submissão pode bloquear; `Go` não recebe contexto para
  cancelar essa espera.
- O canal de tarefas é sem buffer. Limitar workers não constitui uma fila
  durável com N posições.
- O teste inspecionado observa concorrência máxima/conclusão; não demonstra
  cancelamento ou rejeição da submissão. O teste da fonte não foi executado.

Cada declaração continua `DECLARED_NOT_REASSESSED`. A definição genérica do
catálogo aparece em `catalog_definition`; linguagens/ecossistemas aparecem
em `extraction_scope`. A pré-condição da hipótese vem da declaração sobre
configuração e seu trade-off vem da declaração sobre espera, com IDs em
`evidence_refs`. O mecanismo original continua PARTIAL, a hipótese
INCONCLUSIVE e o estado do alvo UNKNOWN. Não houve experimento de efeito
no Bateia ou autorização de adoção.

## Reprodução offline

Na raiz do repositório:

```sh
go build -o /tmp/bateia-insights-0134 ./04_wiring
python3 lab/bancada/insights/0134/verify.py --binary /tmp/bateia-insights-0134
python3 lab/bancada/insights/0106/verify.py --binary /tmp/bateia-insights-0134
/tmp/bateia-insights-0134 idea insights analyze \
  --input lab/bancada/insights/0134/source-insight-input.json \
  --output /tmp/source-insight.json
```

O verificador confere hashes dos 29 itens e trechos/localizadores das seis
observações, reprodução exata de v2/v1, estabilidade por reordenação e prefixo
do hash, mudança causal de IDs, recusa de hash/linha/conteúdo adulterados com
preservação da saída e conservação das observações sem avaliações, sem hipótese.
Esse caso conhecido é diagnóstico/calibração, sem holdout ou prova de
generalização para outras fontes/perguntas.

## Registros

- `source-insight-input.json`, `source-insight-result.json`: pedido e saída v2.
- `source-replay-check.json`, `legacy-replay-check.json`: verificadores reais
  executados sobre o binário deste recorte.
- `baseline.json`: hashes anteriores à alteração dos arquivos Go e corpus 0106.
- `lineage.json`: hashes de código/prompts, comparação com baseline e corpus.
- `core-tests.log`, `cli-tests.log`, `focused-tests.log`: testes da capacidade.
- `sandbox-attempt/`: primeira tentativa do gate; apenas duas fixtures HTTP
  falharam por bloqueio de sockets locais pelo sandbox.
- `loopback-tests.log`: essas fixtures passaram com loopback habilitado.
- `sandbox-retry/`: repetição ainda restrita; GitHub precisou reexecutar uma
  fixture e encontrou o mesmo bloqueio. Sem alteração ou relaxamento do teste.
- `final/verification.json` e logs: gate canônico final.
- `verification.json`: índice consolidado de verificação e limites.

Contrato: [passo 0134](../../../../00_nucleo/passos/0134-conservar-observacoes-especificas-nos-insights.md).
Decisão: [ADR-0015](../../../../00_nucleo/ADR/ADR-0015-observacoes-localizadas-na-investigacao-de-insights.md).
Resultado: [laudo 0134](../../../../00_nucleo/laudos/0134-observacoes-especificas-nos-insights.md).
