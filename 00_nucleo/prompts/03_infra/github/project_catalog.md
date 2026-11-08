# Prompt cristalino — GitHub como catálogo de projetos

Pesquisar unidades de repositório em nome, descrição e tópicos, limitar resultados, registrar rate-limit e resolver owner/repo para commit imutável. Nunca classificar popularidade como correspondência funcional.

## Revisão de segurança 2026-09-26

Aplicar `security.go`: redirect bloqueado, envelopes JSON limitados antes de decodificar (Contents 4 MiB, árvores/descoberta 8 MiB, metadados 1 MiB), arquivo 2 MiB, até 100 arquivos e 8 MiB por extração. Falhas de limite e JSON inválido não podem ativar fallback mais permissivo. Validar identidade, revisão, caminho, cardinalidade e modo regular antes de admitir evidência. Clone fallback desabilitado: indisponibilidade de API/árvore/raw produz coleta inconclusiva; nenhuma execução Git, herança de ambiente/configuração, hook/filtro ou leitura de worktree ocorre no provedor. O legado aceita HEAD mas o preserva como mutável.

## 0154 — Descoberta contextual
Quando MetadataOnly: escopo in:name,description,topics. Quando DeferRevisionResolution: conservar metadados do search sem FetchRepoMetadata por candidato; revisões serão resolvidas nos escolhidos pela porta existente. Campos falsos mantêm fluxo/escopo legado. Não ignorar falhas de integridade nem retornar dados secretos.
