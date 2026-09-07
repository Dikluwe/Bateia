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
   - `RepositoryCandidate`: Candidato retornado no topo do funil com campos de rastreabilidade `CommitID`, `ConsultadoEm`, `Stale`, `Archived`, `DefaultBranch` e `UpdatedAt`.
   - `Dependency`: Dependência declarada (`Name`, `Version`, `Indirect` bool, `Type` string: `"direct"`, `"dev"`, `"peer"`).
   - `Manifest`: Manifesto triado (`Ref`, `ModuleName`, `Language` string, `GoVersion`, `Dependencies`, `ConsultadoEm`, `Content` []byte).
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
| 2026-07-31 | Expansão Multi-Linguagem: adiciona Type em Dependency e Language em Manifest | types.go |

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0035 — Análise semântica C/C++

Cumprir `00_nucleo/passos/0035-corrigir-lint-e-analisar-c-cpp-com-clang.md`. Porta L1 SemanticProvider, execução Clang limitada em L3, injeção L4. Contexto explícito em compile_commands.arguments, sem executar command, plugins ou response files; capturar hashes de configuração, fonte e headers, argumentos, versão, símbolos e chamadas com intervalos. Falta de evidência nunca vira complete. Exposição L2 separa resolução estática e adequação à consulta.

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.

## Passo 0037
Aplicar o contrato de resposta por mecanismo em `00_nucleo/passos/0037-resposta-fundamentada-por-mecanismo.md`. Contexto único de compilação, trechos exatos, limites explícitos e estados desconhecidos preservados.

## Passo 0038
Aplicar contrato do protótipo AST+CFG: adaptador opcional explicitamente configurado, hashes e limites, fluxo intraprocedural e estados não avaliados.
