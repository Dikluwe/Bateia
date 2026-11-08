# 0139 — Leitura automática de README em Gitee e AtomGit

**Data:** 2026-10-28. **Contrato ex-ante.** Capacidade A1-INV sob Tekt. Usuário autoriza ampliar leitura para Gitee/AtomGit. Manter modo GitHub legado, schemas e hashes históricos. Não executar software remoto.

## Responsáveis e fronteiras

L1 declara uma porta InsightSiteRepositoryReader (ResolveRepository por ProductIdentity e métodos do leitor de conteúdo existente). Binding ganha RepositoryReader opcional, admitido somente em Gitee/AtomGit, nunca simultâneo ao Reader GitHub. L3 implementa resolução e leitura com APIs fixas de cada forge; L4 injeta leitor próprio e token próprio. L2 permanece sem mudanças de política.

A evidência legada CommitRef não tem sítio. Portanto não misturar bundles de forges diferentes no agregado legado: acrescentar site_contexts opcionais ao relatório multifonte, cada um com candidate_id, ProductIdentity, bundle congelado e exploração própria. Todos os cartões desses contextos recebem IDs derivados da identidade do sítio e do cartão original; tarefas top-level referenciam esses IDs junto aos cartões GitHub. Não alterar owner/repo para simular namespace. Bundles/exploration legados continuam escopo GitHub. Cada contexto novo contém no máximo um README e um cartão; orçamento total de cartões respeita MaxPassages incluindo GitHub. Metadados continuam não verificados mesmo após ler documentação.

## Resolução e limites

Revalidar identidade HTTPS/host/caminho com exatamente owner/repo, sem credenciais/query/fragmento/port/escape/traversal, antes de qualquer rede. Buscar metadata do repositório e exigir identidade correspondente, público e default_branch válida. Resolver essa branch pelo endpoint de commit para SHA hexadecimal de 40 ou 64 caracteres. Árvore raiz e conteúdo recebem sempre esse SHA, nunca branch mutável. Não seguir download_url nem redirects. Apenas modos de arquivo regular (100644/100755), nunca symlink/submodule; árvore não recursiva, até100 entradas, recusar truncamento/paginação incompleta. Seleção/validação de README pelo coletor L1 existente (README.md/rst/txt/sem extensão); não baixar outros arquivos.

Até quatro GETs por referência (metadata, commit, árvore raiz, conteúdo), sem retries, sob deadline da leitura/total. Respostas metadata/commit/árvore limitadas a256KiB e conteúdo JSON a64KiB. Base64 estrito, tamanho declarado igual ao decodificado, SHA de blob Git verificado (SHA1 para40, SHA256 para64), conteúdo máximo32KiB, UTF8/texto validado pelo núcleo. Total admitido128KiB compartilhado com GitHub. Até8 referências; max-calls continua orçamento de catálogo e não conta esses GETs. Branch movida após resolução não muda ref de árvore/conteúdo. Erros serializados são códigos allowlisted; tokens/URLs de erros nunca publicados. AtomGit sem token não faz HTTP; Gitee opcional. Autenticação, ausência, rate limit, integridade, cancelamento e timeout preservam parcial sem inventar leitura.

## Verificação

Fixtures de protocolo com ambos hosts: token separado, SHA fixo mesmo após mudança de branch, conteúdo/hashes e origem corretos, mesmo owner/repo/commit/content em hosts distintos continua contextos distintos. Casos negativos sem rede para identidade inválida/token ausente; erros HTTP, redirects, body excedido, JSON/base64/tamanho/blob inválidos, árvore truncada/paginada, symlink e cancelamento. Testes puros de admissão, wrapper/refs, budgets e parcial, composição L4 e compatibilidade GitHub. Prompts espelhados antes do código, lint/test/vet/build via make verify, recibos novos em lab/bancada/insights/0139. Disponibilidade remota e autenticação real não são comprovadas por fixtures.

Protocolos consultados: [Gitee SDK](https://gitee.com/sdk/gitee5j/blob/main/docs/RepositoriesApi.md), [AtomGit conteúdo](https://docs.gitcode.com/en/docs/apis/get-api-v-5-repos-owner-repo-contents-path/), [árvore](https://docs.gitcode.com/en/docs/apis/get-api-v-5-repos-owner-repo-git-trees-sha/), [commit](https://docs.gitcode.com/en/docs/apis/get-api-v-5-repos-owner-repo-commits-sha/). Manter endpoints/auth das integrações existentes; nenhum fallback entre hosts.
