# Prompt: Adaptador de Conteúdo GitHub (github/client)

**Camada**: L₃ — Infraestrutura  
**Criado em**: 2026-07-25  
**Atualizado em**: 2026-08-24
**Arquivos gerados**: `03_infra/github/client.go`, `03_infra/github/client_test.go`  

---

## Contexto

Conforme ADR-0001 (rev 1), o meio e fundo do funil (Triagem e Extração) utilizam a Contents API autenticada do GitHub (5.000 req/h). Se ocorrer bloqueio ou requisição não-autenticada, o adaptador ativa o fallback interno para `raw.githubusercontent.com` limitado. O fallback por clone foi removido na revisão de segurança de 2026-09-26.

---

## Restrições Estruturais

- **L₃ (Infraestrutura)**: Implementa a interface `domain.ContentProvider` de $L_1$.
- Realiza requisições HTTP REST (`api.github.com`), requisições raw (`raw.githubusercontent.com`) limitadas; não executa Git, checkout, hooks ou filtros locais.
- Proibido importar $L_2$ ou $L_4$.

---

## Instrução

1. Definir `Client` que implementa `domain.ContentProvider`:
   ```go
   type Client struct {
       token      string
       httpClient *http.Client
   }
   ```
2. Implementar `NewClient(token string) *Client`.
3. Implementar `FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error)`:
   - Consultar Contents API para `/go.mod`.
   - Se retornar 403/429/unauthorized, ativar `fetchRawFallback`.
   - Retornar erro `MANIFEST_NOT_FOUND (404/410)` se o manifesto não existir no repositório.
   - Invocar o parser puro `parser.ParseGoMod`.
4. Implementar `FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error)`:
   - Com SHA resolvido, usar Git Trees recursiva e rejeitar respostas truncadas.
   - Sem SHA, preservar a leitura da raiz para o funil legado.
   - Propagar falhas; nunca fabricar uma árvore substituta.
5. Implementar `ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error)`:
   - Validar candidatos e caminhos, conferir modos regulares na árvore Git e tentar HTTP (Contents API / Raw fallback). Symlinks, submodules, árvore incompleta e limites excedidos são terminais. Nunca clonar nem executar Git; indisponibilidade de HTTP permanece erro.
6. Implementar `FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error)`:
   - Consultar `GET https://api.github.com/repos/{owner}/{repo}`.
   - Atualizar `Stars`, `Description`, estado arquivado, branch e atividade.
   - Resolver o commit da branch padrão quando a entrada ainda usar `HEAD`.
   - Tratar rate limit (403/429) e erro 404 com fallback gracioso mantendo os metadados existentes.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | client.go, client_test.go |
| 2026-07-31 | Adição de FetchRepoMetadata | client.go, interfaces.go |
| 2026-08-24 | Árvore recursiva imutável e metadados de manutenção | client.go, client_test.go |

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Revisão de segurança 2026-09-26

Aplicar `security.go`: redirect bloqueado, envelopes JSON limitados antes de decodificar (Contents 4 MiB, árvores/descoberta 8 MiB, metadados 1 MiB), arquivo 2 MiB, até 100 arquivos e 8 MiB por extração. Falhas de limite e JSON inválido não podem ativar fallback mais permissivo. Validar identidade, revisão, caminho, cardinalidade e modo regular antes de admitir evidência. Clone fallback desabilitado: indisponibilidade de API/árvore/raw produz coleta inconclusiva; nenhuma execução Git, herança de ambiente/configuração, hook/filtro ou leitura de worktree ocorre no provedor. O legado aceita HEAD mas o preserva como mutável.

## Instante de coleta — 0128

Ao ler manifesto remoto, adquirir instante UTC em L3 e fornecer a ParseManifestAt. Não alterar coleta de bytes, fallback, validações ou enriquecimento de README/go.sum.
