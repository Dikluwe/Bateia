# 0166 — Validar bytes, testemunhos e políticas

## Intenção anterior à materialização

Corrigir os quatro achados reproduzidos na revisão de 47434cf. Responsáveis L1: algorithm_insight, evidence_recovery, technical_recheck e reuse_decision. Declarar contratos antes do patch; preservar limites, interpretações, formatos públicos, testes e registros/hashes anteriores.

## Comportamento esperado

Formulação valida bytes originais, sem substituição causada pela cópia JSON. Copiar somente o estado mutável necessário, mantendo quotas e isolamento do caller; rejeitar conteúdo UTF-8 inválido e hash dos bytes substituídos, preservando U+FFFD originalmente válido e IDs legítimos.

Na recuperação, metadados isolados não testemunham tentativa de código/teste derivado. Sem resultado pertinente, somente item de caminho pertinente pode testemunhar; target real de metadados e resultados explícitos NOT_FOUND/FAILED conservam sua regra. Não criar um resultado de ausência só para justificar a contagem. Preservar projeção legada, regra da última ocorrência e agregação0165.

Releitura técnica rejeita conteúdo bruto UTF-8 inválido antes da admissão/hash, registra falha e incompletude e não o substitui. Fonte válida conserva bytes/hashes na publicação e reabertura, inclusive U+FFFD literal. Preservar tamanhos, budgets e distinção entre integridade e interpretação.

Decisão rejeita limiares não finitos e não ignora erro de serialização da política. Políticas finitas conservam hash e comportamento; não alterar avaliações ou decisões históricas.

## Aceitação e preservação

Regressões públicas antes/depois dos quatro casos: campos e evidência inválidos nos dois papéis da formulação, controles válidos e inputs intactos; metadados sem resultados, árvores vazia/sem matches e controles de caminhos/status pertinentes; recheck com bytes inválidos e replay de fontes válidas; NaN/infinitos em cada limiar e limites finitos. Executar make verify com V0–V27 selecionadas, disabledRules vazio e análise completa. Conferir SHA256 do acervo, links, linhagem e diff. Logs/binários em /tmp; não reescrever recibos nem recalcular hashes históricos ou promover módulos selados.

## Resultado e verificação

As regressões reproduziram os quatro achados antes do patch e passaram depois. A formulação usa cópias próprias dos slices alterados, sem roundtrip JSON que substitua bytes antes da validação. Preservou o ID válido congelado antes do patch e o controle de U+FFFD literal. Metadados isolados deixaram de registrar execução ou bloquear retomada de código/teste; itens pertinentes e resultados explícitos conservaram a política anterior. A agregação0165 e os consumidores de última ocorrência também passaram.

A releitura técnica recusou o conteúdo inválido com falha/incompletude e sem item, preservando três roundtrips dos controles válidos e sua revalidação. Políticas com NaN/infinitos foram recusadas, os limites finitos continuaram válidos e o hash literal anterior permaneceu. Todos os testes anteriores ficaram intactos. Logs das regressões em /tmp/bateia-0166-before.log e /tmp/bateia-0166-focused.log; probes da revisão reaplicados em /tmp/bateia-0166-after-{formulation,witness,recheck,policy}.log. São diagnósticos locais, não recibos selados.

O primeiro make verify passou em lint, arquitetura, vet e build, mas um teste HTTP GitHub foi impedido pelo sandbox de abrir socket de loopback. A repetição autorizada com loopback passou em todos os gates: lint, arquitetura, go test ./..., go vet ./..., build de ./04_wiring e diff-check. Evidências finais/binário em /tmp/bateia-0166-verification-final; tentativa restrita preservada em /tmp/bateia-0166-verification. Assinatura Git das fixtures foi desativada somente no processo pelo runner canônico.

O alcance efetivo do linter confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true, 1.183 fontes lidas e zero violações/erros de leitura/sintaxe. Os 11 controles do parser e dois controles de V3 ligado/desligado passaram; três exclusões históricas de oráculos permaneceram declaradas.

A comparação SHA256 confirmou 4.693 arquivos históricos de Lab, passos e laudos intactos; somente o índice ativo de passos foi excluído da comparação por integrar a atualização. Links, linhagem e diff foram conferidos. Fontes/patches do linter, registros, snapshots e hashes históricos permaneceram intactos; módulos selados conservam seus próprios commits e recibos.
