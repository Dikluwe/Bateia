# 0135 — conc × baseline da fila do Bateia

Pedido curado de comparação com investigação v2 congelada de `sourcegraph/conc`
e snapshot dos bytes atuais de `03_infra/runworker/local_queue.go` do Bateia.
O snapshot é de conteúdo fornecido, sem atribuição a HEAD e sem garantia de
completude do alvo. O comando não busca fontes ou executa esses arquivos.

## Leitura manual do caso

| Aspecto | Declaração sobre conc | Declaração sobre Bateia |
| --- | --- | --- |
| Limite | Workers, ilimitado por padrão | Capacidade contabiliza reservas/estados ativos; canal configurado com buffer |
| Configuração | `WithMaxGoroutines` antes do uso | 1–4 workers, capacidade 1–32, handler não nulo |
| Saturação | Submissão pode bloquear | `Reserve` retorna `ErrQueueFull` |
| Compromisso | Chamador precisa poder aguardar | Chamador precisa tratar a recusa |
| Limites | Teste inspecionado não demonstra cancelamento da submissão | Cancelamento cooperativo não força encerramento de handler que ignore contexto |

Essa leitura é curadoria autoral, apoiada nos trechos do pedido/resultado.
O programa agrupa declarações por classe/tipo, sem deduzir essas diferenças
semanticamente. As cinco linhas BOTH_DECLARED indicam disponibilidade;
nenhuma equivale a mecanismo ou efeito comprovado no alvo. A comparação
propõe conferir a política já declarada antes de alterar o Bateia e planejar
experimento para uma diferença concreta, sem emitir hipótese genérica de
adicionar fila limitada. Não houve nova coleta ou execução de conc.

## Reprodução offline

```sh
go build -o /tmp/bateia-insights-0135 ./04_wiring
/tmp/bateia-insights-0135 idea insights compare \
  --input lab/bancada/insights/0135/comparison-input.json \
  --output /tmp/comparacao-insights.json
python3 lab/bancada/insights/0135/verify.py --binary /tmp/bateia-insights-0135
python3 lab/bancada/insights/0134/verify.py --binary /tmp/bateia-insights-0135
python3 lab/bancada/insights/0106/verify.py --binary /tmp/bateia-insights-0135
```

O manifesto do snapshot é uma lista ordenada por path de objetos com chaves
`path`, `content_hash`, nessa ordem, hashes `sha256:`, JSON UTF-8 compacto sem
newline/escaping HTML. O digest SHA-256 desses bytes não inclui observações;
a identidade do relatório inclui também as observações canonicalizadas.
O contrato define limites, campos e tratamento de desconhecidos.

## Registros

- `baseline.json`: hashes anteriores à alteração, incluindo Go e artefatos históricos.
- `comparison-input.json`, `comparison-result.json`: pedido e resultado reais.
- `comparison-check.json`: replay, manifesto, localizadores, referências,
  estabilidade/causalidade, adulteração e controle sem baseline.
- `legacy-v1-check.json`, `source-v2-check.json`: replays históricos preservados.
- `focused-tests.log`, `lint-focused.log`: validação da capacidade.
- `initial-attempt/`: gate que apontou ausência do arquivo de teste espelhado
  do normalizador. Os testes existentes foram separados no arquivo próprio.
- `final/verification.json`: gate canônico final e logs.
- `lineage.json`, `verification.json`: linhagem e verificação consolidadas.

Este caso conhecido é diagnóstico/calibração, sem holdout ou prova de
generalização para outros projetos. Declarações podem ser incompletas ou
incorretas mesmo com hash/localizador válido. Nenhum efeito foi medido.

[Contrato](../../../../00_nucleo/passos/0135-comparar-insights-com-baseline-do-alvo.md)
· [ADR](../../../../00_nucleo/ADR/ADR-0016-baseline-de-conteudo-na-comparacao-de-insights.md)
· [Laudo](../../../../00_nucleo/laudos/0135-comparacao-de-insights-com-alvo.md)
