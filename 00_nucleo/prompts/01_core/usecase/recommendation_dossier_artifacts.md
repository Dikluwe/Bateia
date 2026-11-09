# recommendation_dossier_artifacts

Possuir pré-imagens, manifesto causal e estágios de dossiê/trace em L1. SHA-256 de bytes via leitor de artefato lógico; resolução física só no adaptador. Preservar ordem JSON, versões e status, inclusive perfil opt-in ausente. Erros de sidecar nunca anunciam estágio completo.
Contrato 0131/ADR-0014.


## Integridade 0169

Canonicalizar SearchLocales antes da pré-imagem, rejeitando locale inválido; perfil vazio preserva hash anterior. Construção dos estágios mantém assinatura pública: falha de hash ou hash vazio marca o estágio FAILED, sem Outputs; dependentes de entrada inválida também falham. Não anunciar COMPLETE com dependência ausente. Caminho bem-sucedido preserva bytes e hashes.
