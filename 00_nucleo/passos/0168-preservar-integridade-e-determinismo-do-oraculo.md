# 0168 — Preservar integridade e determinismo do oráculo

## Intenção anterior à materialização

Fechar os seis achados da revisão de a218335. Responsáveis L1: assisted_judgment (selo, variantes, precedência e grounding), multicriteria_decision (aritmética de escalas) e market_discovery (consolidação). Critérios ASS-01–04/MCDA-01/MKT-01 do [contrato](../contratos/integridade-de-entradas-e-replay-v1.md#ass-01--snapshot-íntegro-antes-da-união-cega). Não alterar schemas, rubrica, gates de reúso, registros ou hashes históricos.

## Comportamento esperado

Rejeitar snapshot sem replay íntegro antes da união, sem recalcular selo; conservar corpus vazio. Preservar variantes de title/description de cada URL em ordem canônica, com IDs legados no primeiro texto e sufixos hash nas variantes. Gate de evidência não pode ser mascarado por revisão inválida, e grounding conta vínculo válido, não presença de histórico. Preservar cegamento, texto hostil como dado, revisões e desconhecido.

Escalas finitas cuja diferença transborda usam a mesma normalização sobre valores/extremos divididos por2, mantendo caminho anterior nos casos ordinários. Sucesso produz JSON finito e ausência permanece intervalo completo. Consolidação escolhe representante por JSON canônico e agrega procedência/conjuntos sem depender da ordem, inclusive witnesses e desempates de origens/links. Sem novas equivalências ou promoção de força.

## Aceitação e preservação

Regressões antes/depois: selos ausentes/adulterados, permutações com variantes, controles únicos/idênticos com hashes anteriores; evidência e revisões válidas/inválidas em matrizes e ordens; grounding0/misto/1/ausente/vazio; escalas extremas nos dois sentidos com endpoints, interior, intervalos e ausência; permutações de metadados/procedência de mercado e inputs intactos. Dois controles ativos possuíam fixtures positivas não seladas: atualizar somente sua preparação e prompts correspondentes, mantendo assertivas; isso não é reparação de snapshots históricos. Nenhum arquivo do acervo Lab/passos/laudos anterior deve ser alterado.

Executar make verify com V0–V27 selecionadas, disabledRules vazio, análise completa, testes, vet/build e diff. Comparar SHA256 do acervo e conferir links/linhagem. Evidências/binários em /tmp; módulos/gates selados mantêm seus próprios commits e recibos.

## Resultado e verificação

As regressões reproduziram os seis achados antes do patch e passaram depois. Prepare exige replay íntegro e não produz pacote/saída CLI para selo ausente ou conteúdo adulterado. Variantes antes sobrescritas agora permanecem citáveis, com pacote idêntico sob permutações; controles únicos/idênticos conservaram o hash literal anterior. Novos pacotes com variantes têm nova identidade causal, sem substituição dos pacotes históricos.

Falha de evidência prevaleceu na matriz de revisão válida/inválida e sob reordenação dos candidatos; grounding passou a0 para referência inválida e1/3 no controle misto com julgamento ausente. Históricos e Effective=UNKNOWN nos casos de revisão inválida foram conservados. As assertivas dos controles0021 anteriores permaneceram iguais; somente as duas fixtures positivas ativas foram explicitamente preparadas com selos e identidades, e seus prompts registram a distinção em relação ao acervo preservado.

A normalização multicritério passou nos dois sentidos em escalas ordinárias, ±1e308 e ±MaxFloat64, incluindo endpoints, interior, intervalos e ausência; scores/contribuições e JSON ficaram finitos. O relatório ordinário conservou o hash anterior. Mercado passou nas permutações com variantes e procedência, incluindo witnesses e empates de origens/repo links; controles canônicos único/idêntico mantiveram seus hashes. Inputs e identidades entre conectores foram preservados.

Logs locais em /tmp/bateia-0168-before.log, /tmp/bateia-0168-focused.log e /tmp/bateia-0168-cli.log. Probes originais reaplicados em /tmp/bateia-0168-after-{assisted,mcda,market}.log: snapshot adulterado agora retorna ErrIntegrity; os dois sentidos da união têm hash igual e todos os textos; a falha não é mascarada e grounding inválido é0; escala extrema serializa; metadados de mercado ficam estáveis. São diagnósticos locais, não recibos selados.

make verify passou em lint, arquitetura, go test ./..., go vet ./..., build de ./04_wiring e diff-check. Loopback foi permitido para fixtures HTTP; o runner canônico desativou assinatura Git somente no processo das fixtures. Evidências/binário em /tmp/bateia-0168-verification. O alcance efetivo confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true, 1.190 fontes lidas e zero violações/erros de leitura/sintaxe. Os 11 controles do parser e dois controles de V3 ligado/desligado passaram; as três exclusões históricas de oráculos permaneceram declaradas.

A comparação SHA256 confirmou 4.695 arquivos históricos de Lab, passos e laudos intactos; somente o índice ativo de passos foi excluído por integrar a atualização. Links, linhagem e diff foram conferidos. Nenhum arquivo, registro, snapshot, selo ou hash do acervo anterior foi recalculado; fontes/patches do linter e gates selados conservaram seus bytes, commits e recibos próprios.
