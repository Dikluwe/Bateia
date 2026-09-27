# Prompt cristalino — superfícies README e tópicos do GitHub

README usa busca de código apenas com filename documental e aceita somente caminhos raiz permitidos. Tópicos exigem igualdade do discriminator normalizado. Ambos resolvem owner/repo para commit; rate-limit permanece explícito.

## Revisão de segurança 2026-09-26

Aplicar `security.go`: redirect bloqueado, envelopes JSON limitados antes de decodificar (Contents 4 MiB, árvores/descoberta 8 MiB, metadados 1 MiB), arquivo 2 MiB, até 100 arquivos e 8 MiB por extração. Falhas de limite e JSON inválido não podem ativar fallback mais permissivo. Validar identidade, revisão, caminho, cardinalidade e modo regular antes de admitir evidência. Clone fallback desabilitado: indisponibilidade de API/árvore/raw produz coleta inconclusiva; nenhuma execução Git, herança de ambiente/configuração, hook/filtro ou leitura de worktree ocorre no provedor. O legado aceita HEAD mas o preserva como mutável.
