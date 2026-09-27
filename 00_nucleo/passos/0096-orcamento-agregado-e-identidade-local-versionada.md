# Passo 0096 — Orçamento agregado e identidade local versionada

Continuação de 0094/0095, restrita à admissão local de `idea recommend`.
Pesquisa simples permanece disponível. Sem dependência nova, acesso externo,
execução de candidatos ou alteração de certificados históricos.

## Contraexemplo de identidade

O hash antigo concatenava path + NUL + conteúdo sem delimitar registros.
As árvores `{a:"Xb\0Y"}` e `{a:"X", b:"Y"}` produziam a mesma pré-imagem
`a\0Xb\0Y`, com SHA-256
`fcdcf93cd1ef536fba84fd05e3ef0fa57922c0577e69cb9214110e8ca306c889`.
É colisão estrutural de serialização, não quebra criptográfica do SHA-256.

O formato `tree-sha256-v2` usa domínio `bateia-local-tree-v2` seguido de NUL.
Cada arquivo tem tag 1, comprimento do path em uint64 big-endian, path relativo
em slash, comprimento do conteúdo em uint64 big-endian e os bytes. Tag 0 termina
a árvore. O conteúdo efetivamente lido deve ter exatamente o tamanho anunciado;
alteração detectada durante leitura retorna erro. A ordem continua DFS lexical
por diretório, não ordenação global de paths. Diretórios vazios e symlinks
ignorados não compõem essa identidade; ela representa o inventário de arquivos
admitidos, não todo o estado do filesystem.

`RuleVersions.local_snapshot` identifica o novo formato e participa da comparação
de checkpoints. Retomadas locais antigas não reaproveitam evidências: `--resume`
recalcula e atualiza o dossiê. Para preservar o histórico, execute com um novo
`--output`. O oráculo v1 do teste de 0095 é explicitamente supersedido. Hashes
SHA-256 de arquivos individuais/replay não mudam; certificados anteriores não
foram reescritos.

## Orçamento e confinamento

- Até 20 raízes, sem duplicatas lexicais normalizadas.
- Até 20.000 entradas e 64 MiB de conteúdo, compartilhados entre todos `--local`.
- Até 16 MiB por arquivo, 64 níveis de diretórios e paths de 4096 bytes.
- Diretórios lidos em lotes de 128 antes de ordenar; conteúdo em buffers de 32 KiB.
- Nomes ignorados/excluídos contam como entradas; suas subárvores não são lidas.
- `.git`, `target` e `node_modules` ignorados somente como subdiretórios;
  `.git` regular e raiz com esses nomes continuam sendo considerados.
- Leitura por um `os.Root` retido, sem subprocessos, com descritores não
  bloqueantes/no-follow. Arquivos não regulares e hardlinks são recusados.

O prazo único inicia antes da admissão e segue até o pipeline; não é renovado
por raiz. Timeout não positivo é recusado. Falha de hashing retorna erro sem
hash/lista parcial e antes de criar o dossiê. Output descendente é excluído;
output igual/ancestral da raiz é recusado. Raiz/saída com componentes symlink
são recusadas conservadoramente para evitar aliases na exclusão. Não resolve
todos os aliases de mounts/bind mounts. Usuários precisam fornecer caminhos
diretos, sem aliases symlink, dentro das quotas.

Na validação de dossiê existente, não enumerar todos os arquivos: sem `--resume`,
ler uma única entrada em descritor no-follow e conferir o diretório retido;
com `--resume`, somente conferir o tipo. Isso fecha a leitura sem orçamento de
uma saída excluída do hash. A admissão preserva os artefatos existentes.

Cancelamento é cooperativo entre operações/leitura, não deadline imposto pelo
SO em syscall já bloqueada. Conferências de size/modtime/nlink no descritor não
constituem snapshot atômico, detecção de toda mutação adversarial ou autenticação
de origem. A árvore pode mudar entre hash e coleta posterior. A validação lexical
de output também não é uma sandbox contra alterações concorrentes de ancestrais.

## Verificação

Corpus permanente em `security_tree_hash_test.go`: colisão legada/v2, oráculo
independente, DFS/exclusões, quotas exatas/excedidas, cancelamento, FIFO/hardlink,
raiz retida após rename, versões de checkpoint e falhas antes de publicação.
`security_tree_hash_integration_test.go` verifica o orçamento real entre raízes e
o comportamento de retomada ponta a ponta; `security_dossier_test.go` cobre a
admissão da saída sem alterar os artefatos existentes.

Validação final em 2026-09-26: `go test ./...`, `go vet ./...` e
`git diff --check` passaram. `go test -race` passou no corpus de hash/dossiê,
estabilidade de fileHash e regressões CLI de local/resume/identidade.
Isso não certifica ausência de vulnerabilidades ou autenticidade da evidência.
