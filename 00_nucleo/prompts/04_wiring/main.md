## Passo 0094 — scanner opcional

Despachar security scan antes de .env, caches ou clientes de pesquisa; argumentos
após scan seguem para RunSecurityScanCommand e vulnerability.NewRunner. Não instalar
scanner, compilar ou executar candidato. O scanner pode consultar a base oficial.

# Prompt: Ponto de Entrada e Fiação (wiring/main)

**Camada**: L₄ — Fiação  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `04_wiring/main.go`  

---

## Contexto

Conforme ADR-0002, a camada $L_4$ (Fiação) é o ponto de montagem único da aplicação. Ela substitui a subpasta redundante `cmd/bateia/` por uma fiação direta em `04_wiring/main.go`.

---

## Restrições Estruturais

- **L₄ (Fiação)**: Absorve a complexidade de injeção de dependências.
- **Pode Importar Todas as Camadas**: Importa $L_1$ (`domain`, `usecase`), $L_2$ (`cli`, `report`) e $L_3$ (`cache`, `sourcegraph`, `github`).
- **Zero Lógica de Negócio**: Não pode conter estruturas condicionais de negócio ou algoritmos de decisão. Apenas instanciação e repasse de controle.

---

## Instrução

1. Em `04_wiring/main.go`, definir a função `main()`.
2. Ler tokens de ambiente (`GITHUB_TOKEN`, `SOURCEGRAPH_TOKEN`).
3. Instanciar os adaptadores concretos de $L_3$:
   - `DiskCacheStore`
   - Sourcegraph `Client`
   - GitHub `Client`
4. Instanciar o caso de uso `SearchFunnel` de $L_1$, injetando os adaptadores de $L_3$ como implementações das interfaces de $L_1$.
5. Instanciar e executar a `CLIApp` de $L_2$, passando o orquestrador e o gerador de relatórios Markdown.
6. No caminho `mcp`, interpretar opcionalmente
   `BATEIA_MCP_MAX_TOOL_CALLS` como inteiro decimal positivo. Quando presente,
   aplicar o orçamento process-local ao servidor MCP antes de `Serve`; quando
   ausente, preservar comportamento ilimitado. Valor vazio, zero, negativo ou
   inválido deve falhar fechado antes de servir requisições. L₄ apenas traduz a
   configuração e injeta o limite; a contenção pertence a L₂.

---

## Critérios de Verificação

```
Dado o comando "go run ./04_wiring/main.go --help"
Quando executado no terminal
Então exibe as opções de ajuda da CLI do Bateia sem erros de compilação ou pânico.
```

---

## Resultado Esperado

- `04_wiring/main.go`: Ponto de entrada executável do Bateia.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 (ADR-0002) | main.go |
| 2026-09-03 | Orçamento MCP process-local para o piloto de temperatura do Passo 0029 | main.go |

## Passo 0035 — Análise semântica C/C++

Cumprir `00_nucleo/passos/0035-corrigir-lint-e-analisar-c-cpp-com-clang.md`. Porta L1 SemanticProvider, execução Clang limitada em L3, injeção L4. Contexto explícito em compile_commands.arguments, sem executar command, plugins ou response files; capturar hashes de configuração, fonte e headers, argumentos, versão, símbolos e chamadas com intervalos. Falta de evidência nunca vira complete. Exposição L2 separa resolução estática e adequação à consulta.

## Passo 0038
Aplicar contrato do protótipo AST+CFG: adaptador opcional explicitamente configurado, hashes e limites, fluxo intraprocedural e estados não avaliados.

## Passo 0062
Despachar effects antes de carregar .env ou clientes; injetar effects.NewRunner em cli.RunEffectsCommand.

## Passo 0068

Instanciar o expansor lexical produtivo com limite três em todos os três caminhos atuais
de `main.go`. Isto apenas permite as três propostas do catálogo técnico; regras legadas
continuam com suas duas saídas. Não adicionar decisão de negócio a L4.

## Passo 0070

Despachar `algorithm history analyze` antes de carregar `.env`, cache, expansores ou
clientes L3. L4 apenas encaminha argumentos e stdout para
`cli.RunAlgorithmHistoryCommand`; a análise histórica permanece em L1 e a casca offline
em L2.

## Passo 0093

Despachar `idea evidence recheck` para CLI com cliente GitHub existente e
credencial opcional. Não usar o roteador local neste comando: worktree não
autentica bytes de um commit declarado. Preservar despacho e busca legados.


## Passo 0098 — Trace offline

Despachar idea trace antes da criação do expansor lexical e antes de .env,
cache ou clientes remotos. Encaminhar apenas os argumentos e stdout a
cli.RunIdeaCommand, sem injetar provedor ou executar candidatos. Preservar
os caminhos existentes dos demais subcomandos.


## Passo 0105

Despache idea insights antes de carregar .env, vocabulário, caches e clientes, como trace. Não criar runtime para este comando offline.

## Passo 0137

Exceção de transporte explicitamente online: despachar idea insights discover antes do ramo offline, instanciar cliente GitHub existente (token opcional do ambiente/.env) e injetar como ProjectCatalog e porta estreita de conteúdo em cli.RunInsightDiscoveryCommand. Não criar runtime durável/cache ou plano/política na composição. Os demais insights continuam offline como antes.

0138: injetar adicionalmente factory lazy insightDiscoverySites para --sources; sem opção nova não construir o conjunto federado. Composição da lista no owner L4 próprio, política permanece L1.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

## Passo 0129 — configuração de entrada
A interpretação do limite MCP pertence a mcp.ParseMaxToolCalls em L2. A fiação lê a configuração e delega; a função privada histórica permanece como ponte para testes de composição. Deadlines de shutdown são adquiridos pelo host L3.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.

## Consolidação 0131
A interpretação do limite MCP e seus oráculos estão consolidados em 02_shell/mcp/config.go e config_test.go. A ponte privada e o teste duplicado em L4 foram removidos após comparação de igualdade dos corpos dos oráculos. O caminho produtivo de recomendação injeta serviço com runtime explícito.


## Evolução0143

Passo0143: despachar idea insights expand antes do ramo offline, carregar configuração online e injetar a factory insightDiscoverySites em RunInsightNeighborhoodCommand; sem decidir foco/fronteira.


## Reorganização0145

Delegar interpretação literal de --no-token a cli.WebTokenRequested; composição conserva carregamento de ambiente/cliente e ordem histórica.
