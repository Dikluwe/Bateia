# Prompt: Tipos e Interfaces de Domínio (domain/types)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `01_core/domain/types.go`, `01_core/domain/interfaces.go`  

---

## Contexto

O motor de busca Bateia processa repositórios através de um funil de três estágios (Descoberta -> Triagem -> Extração).
Este componente define a representação agnóstica de todas as entidades do domínio, os objetos de valor rastreáveis (`repo@commit` e `consultado_em`) e os contratos das interfaces dos serviços que serão implementados na infraestrutura ($L_3$).

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Zero I/O, zero dependências externas (apenas a stdlib do Go, ex: `time`, `fmt`).
- Deve definir apenas as estruturas de dados e as interfaces (`SearchProvider`, `ContentProvider`, `CacheStore`).
- Não deve conhecer detalhes de como o cache é salvo em disco ou como requisições HTTP são feitas.
- Proibido importar qualquer código de $L_2$, $L_3$ ou $L_4$.

---

## Instrução

1. Em `01_core/domain/types.go`, definir os tipos:
   - `CommitRef`: Objeto de valor com `Owner`, `Repo`, `Commit` e método `Key() string` (`owner/repo@commit`).
   - `RepositoryCandidate`: Candidato retornado no topo do funil com campos de rastreabilidade `CommitID`, `ConsultadoEm` (time.Time) e `Stale` (bool).
   - `Dependency`: Dependência declarada (`Name`, `Version`, `Indirect`).
   - `Manifest`: Manifesto triado (`Ref`, `ModuleName`, `GoVersion`, `Dependencies`, `ConsultadoEm`).
   - `DirectoryNode`: Nó na árvore de arquivos (`Path`, `Type`, `Size`).
   - `QueryFilter`: Filtros informados pelo usuário (`Query`, `Language`, `MaxResults`, `MustHave`).
   - `TriagedRepository`: Candidato aprovado na triagem com pontuação/motivo.
   - `ExtractedFile`: Conteúdo e tamanho do arquivo extraído.
   - `ExtractedRepository`: Repositório final extraído com seus arquivos.
   - `SearchReport`: Relatório compilado dos 3 estágios.

2. Em `01_core/domain/interfaces.go`, definir as interfaces puras:
   - `SearchProvider`: `SearchCandidates(ctx context.Context, filter QueryFilter) ([]RepositoryCandidate, error)`
   - `ContentProvider`:
     - `FetchManifest(ctx context.Context, candidate RepositoryCandidate) (*Manifest, error)`
     - `FetchTree(ctx context.Context, candidate RepositoryCandidate) ([]DirectoryNode, error)`
     - `ExtractRepoFiles(ctx context.Context, candidate RepositoryCandidate, paths []string) (*ExtractedRepository, error)`
   - `CacheStore`:
     - `GetDiscovery(queryKey string, ttl time.Duration) ([]RepositoryCandidate, bool, error)`
     - `SaveDiscovery(queryKey string, candidates []RepositoryCandidate) error`
     - `GetManifest(ref CommitRef) (*Manifest, bool, error)`
     - `SaveManifest(ref CommitRef, manifest *Manifest) error`
     - `GetExtracted(ref CommitRef) (*ExtractedRepository, bool, error)`
     - `SaveExtracted(ref CommitRef, extracted *ExtractedRepository) error`

---

## Critérios de Verificação

```
Dado um CommitRef com Owner "foo", Repo "bar" e Commit "abc123"
Quando Key() for chamado
Então retorna "foo/bar@abc123"

Dado um CommitRef com Commit vazio ""
Quando Key() for chamado
Então retorna "foo/bar"
```

---

## Resultado Esperado

- `01_core/domain/types.go`: Entidades puras do domínio Bateia.
- `01_core/domain/interfaces.go`: Contratos de serviços agnósticos.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | types.go, interfaces.go |
