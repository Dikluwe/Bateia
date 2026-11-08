# Hash local confinado com orçamento compartilhado

Antes de publicar dossiê, calcular hash local por os.Root retido, sem processos.
Compartilhar até 20 raízes, vinte mil entradas e 64 MiB de conteúdo na admissão
inteira, não por --local. Arquivo até 16 MiB; profundidade64 e paths4096.
Ler diretórios em lotes128 e controlar orçamento antes de sort. Preservar DFS
lexicográfico por diretório do filepath.WalkDir, não sort global de paths.
0096 supersede a pré-imagem histórica ambígua (path + NUL + bytes). Usar formato
tree-sha256-v2: domínio bateia-local-tree-v2 seguido de NUL; para cada arquivo,
tag 1, uint64 big-endian do comprimento do path, path em slash, uint64 big-endian
do tamanho do conteúdo, bytes; tag 0 final. Tamanho lido deve coincidir com
anunciado. fileHash permanece SHA256 dos mesmos bytes. Registrar local_snapshot
na RuleVersions e compará-lo no reaproveitamento: checkpoint antigo/ausente não
é reutilizado. Não reescrever certificados históricos.

Preservar exclusão de output descendente e dirs .git/target/node_modules;
ignorar symlinks estáticos. Recusar output igual/ancestral de raiz e componentes
symlink de raiz/output (inclusive ancestrais de saída inexistente). Recusar raízes
duplicadas normalizadas. Contar nomes excluídos
no orçamento, não enumerar suas subárvores.

Ler conteúdo em streaming de 32 KiB, com limite +1 e contexto da execução.
Recusar arquivo não regular/hardlink e symlink surgido durante resolução;
os.Root é a fronteira contra escapes sob races. Conferir descritor antes/depois
por regular/nlink/SameFile/size/modtime. Qualquer falha retorna hash vazio e erro,
nunca digest parcial ou lista parcial de candidatos. Limites são política privada
do host, não ampliáveis por metadados do candidato.

prepareDossier confirma tipo; sem resume lê somente ReadDir(1) em descritor
no-follow/nonblocking com SameFile. Em resume não enumerar a saída. Não usar
os.ReadDir completo apenas para detectar se está vazio.

Cancelamento é cooperativo entre operações/leitura: não impõe deadline do SO
em syscall de filesystem já bloqueada. Hash não é snapshot atômico de árvore,
assinatura/autenticação de origem nem prova de equivalência do commit Git.

0097: variantes To permitem sink privado na mesma io.CopyBuffer do conteúdo, sem framing nos arquivos. os.Root/O_EXCL, arquivo0400; erro de sink/Close/cancelamento aborta sem digest. Diretórios vazios retidos não mudam hash. Ciclo de vida em local_snapshot_capture.md.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
