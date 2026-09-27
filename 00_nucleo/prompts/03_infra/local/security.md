# Prompt: limites e metadata confinada do provedor local

Camada L3. Gerar `03_infra/local/security.go`.

Percorrer diretórios com descritores rooted e lotes pequenos (`ReadDir(128)`), sem `fs.WalkDir`/`fs.ReadDir` que alocam uma raiz inteira antes de checar orçamento. Limitar varredura a 20.000 entradas, com erro explícito e sem árvore truncada apresentada como completa. Diretórios excluídos não são percorridos; symlinks não são percorridos no inventário.

Expor tipos symlink, hardlink e nonregular honestamente. Leitura de bytes via `os.Root` continua permitindo aliases internos confinados; recusar symlinks externos e qualquer arquivo regular com nlink diferente de 1. Conferir regularidade, nlink, tamanho e modtime no descritor antes/depois da leitura. Isso recusa também hardlinks legítimos como tradeoff conservador e não prova origem histórica do inode.

Máximo de 100 arquivos e 8 MiB por extração; arquivo 2 MiB. Diretório de manifestos também usa leitura em lotes e limite de probes. Busca de implementação local mantém orçamento separado, limitado a 2048 arquivos/32 MiB/512 KiB por arquivo, com teto de árvore 20.000; jamais analisar hardlinks ou aliases protegidos.

Remover toda execução Git. Consultar BaseCommit apenas por leitura confinada `.git/HEAD`, refs soltas e packed-refs, validando SHA completo. HEAD/refs/ponteiro .git até 4096 bytes; packed-refs até 2 MiB/20.000 linhas antes de dividir. HEAD detached e refs empacotadas são positivos. `.git` externo, symlink metadata, gitdir fora da root, metadados em hardlinks e worktrees que dependem de diretório comum externo ficam indisponíveis, sem IO externo ou execução de commands.

BaseCommit é consultivo: não prova que bytes atuais da worktree correspondem ao commit; hash de snapshot continua dos bytes lidos. Em extração, metadata indisponível deixa BaseCommit vazio e registra ChecksNotRun/reason; resolução de metadata retorna erro explícito. Nenhuma garantia de sandbox completo de filesystem, disponibilidade de FUSE ou autenticação de metadata é prometida.
