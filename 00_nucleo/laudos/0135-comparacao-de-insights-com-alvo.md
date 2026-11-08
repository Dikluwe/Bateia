# 0135 — Comparação de insights com o alvo

**Data:** 2026-10-28. **Contrato:** [passo 0135](../passos/0135-comparar-insights-com-baseline-do-alvo.md). **Decisão:** [ADR-0016](../ADR/ADR-0016-baseline-de-conteudo-na-comparacao-de-insights.md). **Evidências:** [bancada](../../lab/bancada/insights/0135/README.md).

## Entrega e responsabilidade

`idea insights compare --input ... [--output ...]` recebe uma investigação v2 e baseline opcional do alvo. A baseline conserva arquivos/observações e digest de um manifesto de conteúdo; não atribui conteúdo de worktree a um commit Git. O relatório `insight-baseline-comparison-v1` conserva as análises da fonte e observações do alvo com trechos, hashes e referências em namespaces separados.

As linhas por fonte/classe/tipo indicam BOTH_DECLARED, SOURCE_ONLY, TARGET_ONLY ou NO_DECLARATIONS. Esses estados descrevem disponibilidade de declarações, sem inferir equivalência, diferença, presença/ausência de capacidade ou benefício. Lacunas e perguntas de validação orientam a conferência da política já declarada antes de propor alteração. A comparação não projeta hipóteses genéricas de adicionar uma capacidade e não altera as hipóteses dos relatórios analyze.

L1 possui normalização de snapshot em `insight_baseline.go` e comparação consultiva em `insight_comparison.go`. A investigação de fontes e a normalização/projeção de observações são reutilizadas. L2 acrescenta transporte e ajuda; não há novo efeito, runtime, cliente concreto ou alteração de política de reúso. README, produto, mapa e índice foram atualizados.

## conc × fila do Bateia

O caso curado conserva integralmente o pedido v2 de 0134, com seis observações da fonte congelada. Acrescenta cinco declarações sobre os bytes atuais de `03_infra/runworker/local_queue.go`: configuração, reserva de capacidade, recusa na saturação, necessidade de tratar a recusa e limite do cancelamento cooperativo. Os trechos/hash pertencem a um snapshot de conteúdo fornecido, com escopo SUPPLIED_FILES_ONLY.

Na leitura manual, a submissão bloqueante de conc contrasta com a reserva/rejeição já declarada no Bateia. Isso oferece uma pergunta mais específica sobre política de admissão e encerramento da espera. O programa agrupa as declarações e seus localizadores; não realizou essa interpretação semântica. As cinco dimensões possuem declarações de ambos os lados, mas o resultado permanece PARTIAL, com efeito UNVERIFIED e NOT_AN_ADOPTION_DECISION. Não houve nova coleta ou execução de conc.

## Verificação

Os testes focados verificaram integridade de snapshot/arquivos/localizadores, quotas, canonicalização, causalidade, imutabilidade e isolamento entre fontes/classes. Cobriram baseline ausente, declarações ausentes ou sem classificação de mecanismo, sem promoção de estado. A CLI foi exercitada com JSON estrito, aliases, short writes, preservação de saída em falha e bloqueio de segredo contido em trecho do alvo.

O verificador do caso conferiu o manifesto e trechos, replay exato, estabilidade por reordenação/prefixo, IDs sensíveis à declaração e recusa de conteúdo/snapshot/linhas adulterados. O controle sem baseline permaneceu INCONCLUSIVE. Alterar o significado declarado não mudou a classificação de disponibilidade, nem comprovou efeito. Os replays v1 e v2 anteriores continuaram idênticos byte a byte.

Comparada à baseline, a mudança preservou todos os testes Go preexistentes e os quinze artefatos históricos inspecionados. Dos 809 arquivos Go anteriores, somente dois arquivos da CLI foram alterados. Sete arquivos Go novos incluem os três componentes produtivos e quatro arquivos de teste; os nove arquivos possuem prompts espelhados e [linhagem](../../lab/bancada/insights/0135/lineage.json).

`make verify` passou integralmente: lint sem violações, onze controles do linter, auditorias sobre 359 arquivos produtivos, suíte raiz, vet, build e diff-check. O [recibo final](../../lab/bancada/insights/0135/final/verification.json) registra comandos e códigos de saída. O binário final passou também nos três verificadores de replay (comparação, analyze v2 e v1). A primeira tentativa do gate apontou a exigência de arquivo de teste correspondente ao normalizador; os testes que o cobriam foram separados sem relaxar regras. O gate usou loopback para as fixtures HTTP e configuração de assinatura Git restrita ao processo.

## Limites

O snapshot só cobre os arquivos fornecidos e sua associação ao alvo é declarada. Hash/localizador correto não comprova a interpretação. O catálogo continua fechado; busca, coleta, interpretação comparativa e medição de efeito não foram automatizadas neste corte. A comparação ainda exige curadoria e revisão, e o caso conhecido é diagnóstico/calibração, sem holdout ou prova de generalização. Gates selados com módulos próprios não são promovidos pela suíte raiz.
