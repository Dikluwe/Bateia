# Prompt: Adaptador de Conteúdo GitHub (github/client)

**Camada**: L₃ — Infraestrutura  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `03_infra/github/client.go`, `03_infra/github/client_test.go`  

---

## Contexto

Conforme ADR-0001 (rev 1), o meio e fundo do funil (Triagem e Extração) utilizam a Contents API autenticada do GitHub (5.000 req/h). Se ocorrer bloqueio ou requisição não-autenticada, o adaptador ativa o fallback interno para `raw.githubusercontent.com` ou clone local shallow (`git clone --depth 1`).

---

## Restrições Estruturais

- **L₃ (Infraestrutura)**: Implementa a interface `domain.ContentProvider` de $L_1$.
- Realiza requisições HTTP REST (`api.github.com`), requisições raw (`raw.githubusercontent.com`) e execuções de comandos git locais (`os/exec`).
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
   - Ler os nós de diretório raiz.
5. Implementar `ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error)`:
   - Tentar via HTTP (Contents API / Raw fallback). Se falhar ou for leitura ampla, executar clone shallow local `git clone --depth 1` em diretório temporário isolado.

---

## Critérios de Verificação

```
Dado uma resposta 403 Forbidden da Contents API para go.mod
Quando FetchManifest for executado
Então aciona o fallback via raw.githubusercontent.com e retorna o Manifest preenchido.
```

---

## Resultado Esperado

- `03_infra/github/client.go`: Adaptador GitHub `Client`.
- `03_infra/github/client_test.go`: Testes unitários com simulador de fallback.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | client.go, client_test.go |
