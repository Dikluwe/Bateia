# 0163 — Fechar completude, bytes e admissão de lotes

## Intenção anterior à materialização

Corrigir os três contraexemplos da revisão de bb9e14a. Responsáveis L1: evidence para completude/classificação de orçamento, evidence_recovery para métricas e histórico e insight_collection para admissão do lote. Estender ao coletor geral a regra de tamanhos já exigida na releitura técnica, mantendo os controles desta. Preservar schemas, registros/hashes e assertivas históricas; atualizar apenas contratos ativos e acrescentar regressões.

## Comportamento esperado

Antes de truncar por orçamento ou produzir hash de conteúdo, confrontar o tamanho positivo de ExtractedFile e o conhecido da árvore com os bytes brutos devolvidos. Zero é desconhecido legado; negativo ou conflito positivo da árvore também invalida o arquivo. Não deixar canonicalização da árvore apagar esse conflito. Arquivo divergente é omitido, registra FAILED/falha de completude e mantém Complete=false. Não pode produzir licença PROVEN nem fonte dirigida examinada. Reutilizar o predicado de coerência de tamanhos na releitura técnica, sem mudar seu contrato. Tamanho referente aos bytes brutos não é comparado depois do corte legítimo por orçamento.

Se o provedor devolveu um arquivo e o saldo de bytes impede sua admissão integral, registrar SKIPPED_BUDGET, motivo de orçamento e incompletude. Omitido integralmente tem Bytes0 e não cria item; truncado mantém os bytes parciais e hash correspondente, com SKIPPED_BUDGET e Bytes admitidos. Ambos conservam cobertura incompleta e diagnóstico BUDGET_EXHAUSTED quando aplicável; não fabricar NOT_FOUND. Somente ausência real da resposta mantém NOT_FOUND. Histórico ativo permite nova tentativa de alvos limitados por bytes. TargetsExecuted conta tentativa com fragmento admitido (Bytes>0), sem contar omissão integral; ProviderCalls continua contando portas reais. Não promover fragmento a licença/capacidade comprovada.

Na coleta coordenada, validar todos os planos e a unicidade das origens planejadas antes de qualquer porta. Usar owner/repo sem distinção de caixa e commit do consumidor; vazio/HEAD representam a mesma origem não resolvida. Recusar duplicatas conhecidas com ErrInvalidRunInput e relatório vazio, sem consolidar implicitamente budgets/targets. Revisões fixas distintas continuam permitidas. Colisões só conhecidas após resolução continuam sujeitas à rejeição de normalização; não prometer prever o retorno externo de HEAD.

## Aceitação e limites

Regressões públicas de tamanho0/exato, maior/menor/negativo, somente árvore/extração e conflito bruto da árvore; coletor, assessor de licença e cobertura dirigida sem promoção. Bytes esgotados após metadados/arquivo anterior, arquivo truncado e arquivo realmente ausente; status/bytes/motivo, diagnóstico e retomada de LICENSE. Verificar métricas de tentativa parcial e omissão integral, sem mudar controles anteriores de saldo. Lote duplicado nas duas ordens, caixa de owner/repo, vazio/HEAD e controle de dois pins distintos; zero portas para rejeição conhecida e entradas intactas.

Executar regressões antes/depois, make verify com V0–V27 selecionadas, disabledRules vazio e análise completa; comparar SHA256 do acervo, conferir links/linhagem e diff. Gates selados e hashes antigos permanecem nos commits/módulos registrados; não recalcular resultados para aparentar conformidade. Logs temporários no /tmp, nunca em L0. Nenhuma nova inferência de comportamento, autenticidade ou benefício.

## Resultado e verificação

Os três contraexemplos falharam nas regressões anteriores ao patch e passaram após a correção. O coletor geral agora compara tamanhos antes do corte por orçamento; arquivos inválidos não entram nos itens, produzem FAILED/falha e não promovem licença ou cobertura dirigida. A regra de tamanho é compartilhada com a releitura técnica; seus controles anteriores continuam passando. A árvore bruta conserva conflitos que a canonicalização descartaria.

Omissão integral e fragmento por saldo recebem SKIPPED_BUDGET com bytes/motivo, preservando cobertura incompleta e causa BUDGET_EXHAUSTED. Ausência verdadeira mantém NOT_FOUND. O probe LICENSE voltou a planejar LICENSE, em vez de COPYING. Fragmento admitido conta como tentativa na métrica, mas não bloqueia retomada no histórico ativo. Essa distinção mantém os controles anteriores de saldo/métricas sem editar suas assertivas. Fragmento não ganhou autoridade de licença/capacidade.

A admissão coordenada recusa duplicatas conhecidas por chave estruturada de owner/repo sem distinção de caixa e commit, com vazio/HEAD equivalentes. O probe com dois planos idênticos passou de seis chamadas para zero; controles de um plano e dois pins distintos continuam válidos. Não houve consolidação implícita de planos nem nova promessa de prever a resolução externa de HEAD.

make verify passou em lint, arquitetura, go test ./..., go vet ./..., build e diff-check, com loopback autorizado para fixtures HTTP e assinatura Git desativada somente no processo de fixtures pelo runner canônico. Logs finais/binário em /tmp/bateia-0163-verification. Regressões em /tmp/bateia-0163-before.log e /tmp/bateia-0163-focused.log; consumidores em /tmp/bateia-0163-core.log; reaplicação dos probes em /tmp/bateia-0163-after-probes.log. São diagnósticos locais, não novos recibos selados.

O alcance efetivo confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true e ausência de violações/erros de leitura/sintaxe. Os 11 controles do parser e dois controles V3 ligado/desligado passaram. Exclusões históricas explícitas permanecem. Seleção completa não implica autenticação remota ou prova semântica universal.

A comparação SHA256 confirmou 4.690 arquivos históricos de Lab, passos e laudos intactos; somente o índice ativo de passos foi excluído por integrar a atualização. Fontes/patches do linter e acervo de gates conservaram hashes. Nenhum teste anterior foi alterado; apenas foram acrescentadas regressões. Links locais, linhagem e diff foram conferidos. Não foram recalculados hashes/snapshots antigos nem reproduzidos módulos selados como resultados novos.
