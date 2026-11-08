# Laudo 0142 — Aprofundamento e memória de exploração

**Data:** 2026-10-28. **Resultado:** materialização e verificação aprovadas para o [contrato0142](../passos/0142-aprofundar-trilhas-e-conservar-memoria.md), sob [ADR-0021](../ADR/ADR-0021-memoria-de-exploracao-e-aprofundamento-limitado.md).

## Entrega

O modo multifonte aprofunda atédois destinos por padrão (--follow-links0–4) e retorna wrappers followups com identidade, vínculo de trilha, recibo, bundle e cartões. O leitor dedicado GitHub confirma metadata pública/identidade, fixa commit, lê árvore raiz e Contents com limites, tamanho/base64/blob verificados, sem fallback raw. Demais forges usam os leitores imutáveis já compostos. Quotas de bytes/documentos/cartões e deadlines são compartilhadas com a coleta primária; nenhum descendente é seguido no mesmo run.

exploration_memory registra caminhos de repositório nas cinco forges com leitor, estados PENDING/EXAMINED/DISCARDED, tentativas, motivo e fingerprint da leitura anterior. --memory carrega snapshot standalone ou de relatório novo. A publicação protegida contém resultado e memória no mesmo arquivo, sem estado global ou segunda escrita divergente. --output protege a entrada contra sobrescrita/alias. O comando offline insights memory descarta ou reabre alvo conhecido com motivo; não fabrica estado examinado.

Pendentes/desconhecidos têm precedência na seleção, e examinados/descartados não são relidos automaticamente. Falhas mantêm pendência e parciais; progresso das leituras da memória é conservado mesmo quando catálogo falha. Histórico importado é declarado, não conteúdo relido; fingerprints não viram cartões atuais. L1 conserva política/estado, L2 transporte, L3 efeitos e L4 composição.

## Verificação

make verify passou: lint e11 controles negativos, auditoria de387 arquivos produtivos, testes do módulo raiz, vet, build e diff-check. A primeira solicitação de acesso a loopback expirou na revisão automática antes de execução; a repetição autorizada passou. Não houve bloqueio de ação restante.

Testes puros verificam validação/hash/cópia, estados/motivos, capacidade, seleção orientada por memória, duas rodadas sem repetição, quotas de128KiB/8documentos/max-passages, fonte/leitor ausente e cancelamento parcial. Transporte GitHub em RoundTripper controlado confirma quatro GETs do adaptador, SHA, publicidade, identidade e rejeição de auth/404/redirect/blob/tamanho/symlink. CLI testa entrada standalone/relatório, validação antes de factory, segredo em motivo e proteção da saída; o binário final também executou descarte/reabertura offline, com snapshot igual ao produzido pelo núcleo.

O replay autoral usa dois READMEs de origem congelados e destinos sintéticos, explicitamente marcados. A primeira rodada exercita duas leituras (incluindo o destino Serendip-o-matic) e gera14 entradas:4 examinadas/10pendentes. A segunda pula pais/destinos examinados, lê dois outros e chega a6examinadas/8pendentes. Feedback conserva motivo. Isso verifica o fluxo, não o conteúdo público ou os mecanismos dos destinos reais.

Os seis replays históricos0106 e0134–0138 passaram com o binário final. Os289 arquivos de evidência anteriores e todos os testes Go anteriores permanecem idênticos. Oito arquivos produtivos existentes foram alterados; foram acrescentados cinco arquivos produtivos e seis de testes, além do replay Lab. Prompts e hashes estão na linhagem nova.

## Limites

Não houve ensaio público de aprofundamento autenticado/atual neste recorte. Metadados permanecem pistas e EXAMINED significa documento coletado, sem comprovar mecanismo ou benefício. Memória não é auditoria completa, cache de conteúdo nem identidade autenticada; até256 repositórios/1MiB, com limites e margem explícitos. Catálogos fora das cinco forges não entram nesta memória inicial. Revisões posteriores de examinados requerem reabertura com motivo, sem expiração automática. Novidade NOT_ASSESSED e efeito UNVERIFIED permanecem; nenhum candidato foi executado.

Consultar [uso/bancada](../../lab/bancada/insights/0142/README.md), [recibo](../../lab/bancada/insights/0142/verification.json), [linhagem](../../lab/bancada/insights/0142/lineage.json) e [controle CLI](../../lab/bancada/insights/0142/cli-workflow-check.json).
