# Corpus de hash local v2 — 0096

Fixtures sintéticas, sem rede, Git ou execução de candidatos. Oráculo independente
constrói a pré-imagem v2 com domínio, tags e comprimentos big-endian; nunca usa
walker/helper produtivo para gerar expectativas. Reproduzir a colisão estrutural
legada {a:Xb-NUL-Y} versus {a:X,b:Y} e exigir hashes v2 distintos.

Preservar DFS lexical por diretório, .git como arquivo, raiz chamada target,
exclusão de output descendente e diretórios ignorados; symlinks estáticos são
ignorados. Verificar quotas exatas/+1 de bytes, entradas compartilhadas entre
raízes, profundidade, paths e quantidade de raízes. Zero-byte file cabe com zero
bytes restantes; nunca retornar hash/lista parcial. Ordem das raízes não altera
candidatos; duplicação é erro. Testar FIFO sem bloqueio, hardlinks, descritor de
raiz retido após rename, contexto cancelado e cancelamento entre leituras.

Output igual/ancestral ou com ancestrais symlink é recusado; raízes com aliases
symlink também. CLI não publica dossiê/stdout quando admissão excede quota/falha
ou timeout não positivo/expirado. Checkpoint só reaproveita local_snapshot v2,
mesmo quando todos os demais hashes e artefatos estão corretos.

Checagem do dossiê preserva arquivos, aceita vazio/resume e recusa saída file,
FIFO e symlink. Não enumerar diretório inteiro para detectar não-vazio.

0096 supersede explicitamente o oráculo de árvore v1 de 0095 porque sua pré-imagem
é ambígua. Hash de arquivo isolado é preservado. Não reescrever certificados
históricos nem enfraquecer critérios para acomodar falhas.

0097 supersedes only the local collection spy: CLI uses real captured bytes, not fixtureContent. End-to-end test requires README bytes/hash, no fake calls and checkpoint_reused0/1. Missing local_collection forces recollection. Additional contract: local_snapshot_capture_test.md.
