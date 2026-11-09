# 0167 — Isolar calibração e validar fontes de síntese

## Intenção anterior à materialização

Corrigir os quatro achados da revisão de c4013b8: grupos da mesma ideia em treino/teste, custo dependente da ordem de conexões repetidas, RunID divergente no retorno do store e cursor repetido não consecutivo. Responsáveis L1: screening_calibration, composition e improvement_sources. Contrato: [CAL-01/CMP-01/SYN-01/02](../contratos/integridade-de-entradas-e-replay-v1.md#cal-01--uma-identidade-por-grupo-de-calibração). Preservar registros, hashes, schemas e resultados anteriores.

## Comportamento esperado

Rejeitar grupos com identidade vazia, espaços nas extremidades ou repetida antes de treinar; não agrupar silenciosamente. Preservar dados válidos, uncertain e separação de priors. Repetições de conexão para o mesmo par ordenado só são admitidas quando estruturalmente idênticas; divergências de status/fontes/custo são erro, sem escolha implícita de custo. Conservar conexões únicas e duplicatas idênticas.

Na resolução de fontes, conferir RunID exato após Get e recusar divergência antes de páginas/idempotência/fila. Lembrar cursores por artefato e rejeitar qualquer repetição antes de reler página. Preservar tipo, terminalidade, hashes, paginação50/1000, erros e base zero em falha; não produzir resultados parciais ou promover hipóteses.

## Aceitação e preservação

Regressões antes/depois: grupos repetidos, vazios e controle válido com hash anterior; permutações de conexões com custos, fontes ou estados divergentes, duplicatas idênticas e custo desconhecido; Get com ID vazio/divergente, paginação A→B→A, ciclos maiores e controles normais, com contagem de portas e ausência de efeitos posteriores. Reexecutar consumidores anteriores, make verify com V0–V27 selecionadas e análise completa. Conferir SHA256 do acervo, links, linhagem e diff. Evidências/binários em /tmp; não alterar testes anteriores nem reescrever recibos/snapshots ou promover gates selados.

## Resultado e verificação

As regressões reproduziram os quatro defeitos antes do patch e passaram após a correção. Grupos inválidos são recusados antes de treinar, com relatório zero; os grupos válidos preservaram o hash literal anterior do relatório, folds e tratamento de uncertain. Conexões divergentes foram recusadas em ambas as ordens, enquanto conexão única e duplicatas idênticas mantiveram o relatório e seu hash anterior, inclusive custo desconhecido.

Get com RunID vazio/divergente passou a falhar com ErrIntegrity sem ler páginas ou iniciar idempotência/fila. Cursores repetidos não consecutivos foram recusados na terceira/quarta leitura, antes da revisita; cursores iguais em artefatos distintos continuaram permitidos. Os controles de identidade da página, erro de porta, hash fallback, teto1000 e projeções congeladas anteriores também passaram. A API pública foi exercitada com sentinela contada após admissão válida, sem criar jobs ou publicar artefatos. Testes anteriores permaneceram intactos.

Logs locais em /tmp/bateia-0167-before.log e /tmp/bateia-0167-focused.log; probes originais reaplicados em /tmp/bateia-0167-after-{calibration,composition,synthesis}.log. O probe de composição conserva as conexões única/idêntica e agora recusa ambas as ordens ambíguas. São diagnósticos locais, não recibos selados.

make verify passou em lint, arquitetura, go test ./..., go vet ./..., build de ./04_wiring e diff-check, com loopback permitido para fixtures HTTP e assinatura Git desativada somente no processo das fixtures pelo runner canônico. Evidências finais/binário em /tmp/bateia-0167-verification. O alcance do linter confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true, 1.186 fontes e zero violações/erros de leitura/sintaxe. Os 11 controles do parser e dois controles V3 ligado/desligado passaram; três exclusões históricas de oráculos permaneceram declaradas.

A comparação SHA256 confirmou 4.694 arquivos históricos de Lab, passos e laudos intactos; somente o índice ativo de passos foi excluído por integrar a atualização. Links, linhagem e diff foram conferidos. Nenhum registro, hash, snapshot ou resultado histórico foi recalculado; fontes/patches do linter e gates selados mantiveram seus bytes, commits e recibos próprios.
