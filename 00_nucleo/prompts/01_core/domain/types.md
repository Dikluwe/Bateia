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

# 0044 — Porta e evidência de recuperação de implementação
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md (contrato congelado).
Tipos de domínio para orçamento, relatório, inputs hashados e hits de AST; porta opcional de recuperação sem IO concreto em L1. Busca pode ser partial/not_evaluated; hits não são comprovação de capacidade. Preservar hashes da varredura e conferir identidade antes de exibir definição. Campos opcionais em resultados de extração/bundles mantêm compatibilidade com consumidores existentes.

## Passo 0055 — hipótese rastreável

Aplicar o contrato congelado 00_nucleo/passos/0055-planejar-consultas-de-codigo.md, que substitui a tentativa única por original seguido de no máximo uma hipótese apenas após vazio válido. QueryPlan, RawCandidates e RejectedCandidates registram somente a derivada; contagens finais são admitidos. Plano preserva pergunta/restrições não verificadas. Ambas portas usam catálogo puro; linguagem e escopo regex literal persistem. Admitir somente FileMatch com caminho e ação+objeto em snippets efetivamente retidos do mesmo arquivo; reter todos os previews limitados individualmente a 240 runes, sem descrição/nome. Não combinar arquivos. Preservar limites, erros, timeouts, cancelamento e rastros nos candidatos. Validar envelope GraphQL antes de considerá-lo vazio. EvidenceVersion muda para evidence-v15-planned-discovery.

Adendo 0055: aplicar 00_nucleo/passos/0055-adendo-evidencia.md. SearchAttempt.AdmissionEvidence ([]ProviderCandidateMatch, JSON admission_evidence,omitempty) somente na derivada, preserva matches admitidos e os snippets efetivamente examinados da porta Discover, também nos rastros anexados aos candidatos. Não publicar rejeitados como evidência de admissão.

## Passo 0056 — cobertura observada

Aplicar o contrato congelado `00_nucleo/passos/0056-explicar-cobertura-e-perdas.md`. Somente na derivada: count limitado20..200 sem overflow, coleta de até N repositórios distintos na ordem da resposta e evidências posteriores desses mesmos repositórios, admissão lexical antes do limite final. Cobertura unknown, matchCount nullable, contagens observadas e diagnósticos com motivos, no máximo32 previews retidos de240 runes normalizadas por repositório; omissões e truncamentos explícitos. Diagnósticos usam somente evidência de arquivo retida, exigindo ação e objeto no mesmo arquivo. Selected/insufficient_evidence/candidate_limit distinguem admissão e seleção. Propagar Coverage e CandidateDiagnostics nas tentativas e rastros; Matches/AdmissionEvidence somente selecionados. Literal permanece compatível. EvidenceVersion evidence-v17-discovery-coverage.
