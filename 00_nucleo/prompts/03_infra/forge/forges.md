# Especificação de catálogos de forjas

**Camada:** L3
**Proprietário:** `03_infra/forge/forges.go`
**Consumidores:** descoberta de mercado L1, composição L4 e testes do adaptador.

## Contrato 0126

- Gitee (`gitee-repository-catalog`) consulta `https://gitee.com/api/v5/search/repositories`; AtomGit (`atomgit-repository-catalog`) consulta `https://api.atomgit.com/api/v5/search/repositories`. Preservar construtor Gitee sem token; oferecer variante autenticada. AtomGit recebe token explicitamente. L3 não lê ambiente.
- Gitee aceita token opcional pelo parâmetro documentado `access_token`; AtomGit usa Bearer conforme introdução da API e exige token antes de consultar. Tokens não aparecem em resultados, erros ou hashes de resposta.
- Manter Codeberg e GitLab. Todos anunciam somente identidade de repositório (`owner_repo`, `repository_identity`); catálogo não observa código nem resolve commit. Novas versões distinguem a mudança sem reescrever artefatos anteriores.
- Uma requisição por consulta, primeira página, limite de 1–10 (default 10), timeout limitado e corpo até 2 MiB; nunca seguir redirects. Cancelamento respeita contexto.
- V5 retorna array; Codeberg também admite envelope `data`. Resposta inválida não equivale a array vazio. HTTP/transportes/credencial ausente geram UNKNOWN_PROVIDER com classe sanitizada; payload inválido gera INTEGRITY_ERROR; cancelamento gera CANCELLED. Somente array vazio válido representa EMPTY_OBSERVED.
- Identidade inclui conector, namespace/repositório e URL HTTPS do provedor. Aceitar namespaces aninhados; não inventar identidade a partir de nome de exibição. AtomGit admite URLs `gitcode.com` documentadas como legado, preservando a origem. Recusar URL com credenciais, query, fragmento, porta ou caminho incompatível.
- Gitee exige `private=false`; AtomGit exige também `public=true` e ausência de indicador interno ativo. Flags booleanas ou inteiros 0/1 são aceitos conforme documentação. Dados privados ou sem visibilidade pública demonstrada são descartados; descarte integral é INTEGRITY_ERROR, não ausência de resultados.
- `Commit` permanece vazio e vínculo UNVERIFIED: metadados não verificam conteúdo, revisão, licença ou reúso. Não consolidar espelhos por semelhança de nome. Guardar SHA-256 da resposta e instante UTC para proveniência.
- Testes offline exercitam URL, autenticação, limites, isolamento de credenciais, visibilidade, namespaces, identidade, falhas, corpo inválido, cancelamento e composição MCP.

## Referências consultadas em 2026-10-18

- https://gitee.com/openeuler/go-gitee/blob/3284d947aa89e38815cb6ffd3cdfe40566e512aa/docs/SearchApi.md
- https://docs.atomgit.com/docs/apis/
- https://docs.atomgit.com/docs/apis/get-api-v-5-search-repositories/

A documentação AtomGit lista flags inteiras no schema e booleanas nos exemplos. Testes locais não comprovam disponibilidade remota ou utilidade dos candidatos.
