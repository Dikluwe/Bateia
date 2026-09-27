# Captura consistente local — 0097

Antes de publicar dossiê, validar as raízes/output pela política0096 e capturar
arquivos na MESMA leitura que alimenta tree-sha256-v2. Não hashear e depois
copiar/reler. Sink opcional no walker; framing continua só no hash; somente os
bytes de conteúdo seguem também para a cópia. Escrever O_EXCL em os.Root privado;
arquivo0400 ao concluir, diretórios0700. Falha de escrita/Close/cancelamento ou
quota retorna erro sem digest/candidatos parciais. Preservar todos os limites,
DFS, exclusões, .git regular e diretórios vazios (estes últimos fora do digest).

Compartilhar o orçamento entre raízes. Até20 raízes;20k entradas;64MiB conteúdo;
16MiB/arquivo;64diretórios;paths4096. Recusar duplicatas lexicais e identidades
owner/repo duplicadas, mesmo que haja colisão no prefixo curto do path-hash.

Container temporário criado com MkdirTemp, fora das raízes para não ingerir a
própria captura. Raízes amplas que incluem armazenamento temporário são recusadas.
Limpeza em erro, timeout, recusa de output, sucesso e resume. Antes de remover,
conferir SameFile e ausência de symlink no container criado; não remover alvo
substituído. Erro de limpeza é reportado, não escondido. Não copiar permissões
executáveis, symlinks, hardlinks ou diretórios ignorados.

API privada capturedLocalRun fornece candidatos originais, SnapshotBindings,
mapa consultivo de identidades Git e Close. Usar local.NewSnapshotProvider na CLI
local sempre, inclusive com ContentProvider injetado; fake local não pode fazer
fallback silencioso à origem. Provider mantém identidades/URLs públicas, usa só
cópias e não fabrica BaseCommit. Live/replay/search simples não mudam.

Git remoto para remapeamento de target é capturado somente se target existe,
antes do pipeline; config até64KiB/raiz pelo leitor confinado existente. Após
captura, lookup usa somente mapa congelado. local_metadata é fingerprint dos
IDs consultivos normalizados, não do config completo nem autenticação da origem.
Fingerprint entra em InputHashes e na pré-imagem RunID, com campo omitempty que
preserva live/replay. Não sustenta capacidade nem exige recoleta de evidência de
arquivos idênticos; ranking é recalculado com os hints desta admissão.

Manter local_snapshot=tree-sha256-v2 e adicionar local_collection=frozen-local-v1
à RuleVersions/comparação de checkpoint. Coleta antiga não é reutilizada, mesmo
com digest idêntico. Stages discover/collect locais registram checkpoint_reused
como contagem0/1; coleta só marca1 quando bundles foram reaproveitados.

Captura é imutável pelo workflow, não contra processo com mesmo UID/administrador.
Não é snapshot atômico de árvore, prova de estado Git, assinatura ou autenticação.
Diretórios vazios e metadata não integram hash; orçamentos de análise podem deixar
arquivos capturados sem exame. Cancelamento de filesystem é cooperativo.
