# Crash-safe local recommendation run store

Implemente em L3 a porta `RecommendationRunStore` do Passo 0026 sobre uma raiz
local recebida na construção. Cumpra `CreateOrResolve` concorrente, CAS por versão,
índice idempotente por digest, publicação PENDING→blob validado→PUBLISHED,
reopen, interrupção de órfãos, quotas e leitura paginada exatamente como
`00_nucleo/benchmark/0026/contract.md`.

Use nomes internos derivados por SHA-256, nunca concatene entrada externa em
caminho. Rejeite traversal, symlink, hardlink inesperado e arquivo não regular;
use diretórios 0700, arquivos 0600, temporário exclusivo, fsync e rename na mesma
pasta. Valide tamanho e SHA-256 antes de publicar. Não sobrescreva terminal ou
artifact publicado e não trate corrupção como ausência.

Codifique cursores com payload canônico e HMAC-SHA-256 usando segredo persistente
privado, vinculando run, nome, artifact version, hash e offset. Nunca registre
chave clara, cursor, segredo, conteúdo, token ou erro físico em superfície pública.
Não importe L2/L4 e não implemente regra semântica de recomendação.

