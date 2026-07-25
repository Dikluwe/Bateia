# Prompt: Armazenamento e Cache em Disco (cache/cache_store)

**Camada**: L₃ — Infraestrutura  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `03_infra/cache/cache_store.go`, `03_infra/cache/cache_store_test.go`  

---

## Contexto

Conforme ADR-0001 e ADR-0002, o cache é um componente de primeira classe do Bateia. O conteúdo extraído de manifestos e código é imutável quando endereçado por `repo@commit` e nunca expira. Os resultados de descoberta no topo do funil expiram por TTL.

---

## Restrições Estruturais

- **L₃ (Infraestrutura)**: Implementa a interface `domain.CacheStore` definida em $L_1$.
- Realiza I/O físico no File System (`os`, `filepath`, `json`).
- Proibido importar $L_2$ ou $L_4$.

---

## Instrução

1. Implementar o struct `DiskCacheStore`:
   ```go
   type DiskCacheStore struct {
       baseDir string
       mu      sync.RWMutex
   }
   ```
2. Criar `NewDiskCacheStore(baseDir string) (*DiskCacheStore, error)` que garante a criação dos diretórios `discovery`, `manifests` e `extracted`. Se `baseDir` for vazio, utilizar `~/.cache/bateia`.
3. Implementar métodos de `domain.CacheStore`:
   - `GetDiscovery(queryKey string, ttl time.Duration)`: Sanitizar a chave e verificar se o arquivo existe em `discovery/`. Se expirar pelo TTL, retornar `hit: false`.
   - `SaveDiscovery(queryKey string, candidates []domain.RepositoryCandidate)`: Salvar JSON com timestamp.
   - `GetManifest(ref domain.CommitRef)` / `SaveManifest(...)`: Ler/gravar em `manifests/` pela chave imutável `ref.Key()`.
   - `GetExtracted(ref domain.CommitRef)` / `SaveExtracted(...)`: Ler/gravar em `extracted/` pela chave imutável `ref.Key()`.

---

## Critérios de Verificação

```
Dado um candidato salvo no cache com chave "foo/bar@c1"
Quando GetManifest for chamado com a mesma ref
Então retorna os dados cacheados e hit == true.

Dado um cache de descoberta com TTL expirado
Quando GetDiscovery for chamado
Então retorna hit == false.
```

---

## Resultado Esperado

- `03_infra/cache/cache_store.go`: Implementação do `DiskCacheStore`.
- `03_infra/cache/cache_store_test.go`: Suite de testes unitários com diretório temporário (`os.MkdirTemp`).

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | cache_store.go, cache_store_test.go |
