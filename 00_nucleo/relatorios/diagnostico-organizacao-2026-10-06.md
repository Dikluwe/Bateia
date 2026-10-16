# Diagnóstico de organização do Bateia

Em 6 de outubro de 2026, a leitura da worktree com base no commit `9d00e2f` identificou um produto orientado à decisão de reúso, organizado em camadas Tekt e acompanhado por um acervo amplo de contratos e avaliações. A principal oportunidade é conectar esse acervo às capacidades atuais e distinguir implementação, conformidade estrutural e utilidade demonstrada.

## Escopo da leitura

Foram examinados propósito, regras locais, configuração do lint, dependências Go, ADRs de contenção/configuração/runtime, índices de passos e relatórios, composição, portas, jornada de recomendação, funil, CLI, MCP e persistência durável. A leitura é estrutural e dirigida aos fluxos principais; não constitui auditoria semântica de cada componente nem avaliação nova de utilidade do produto.

A worktree já continha modificações em domínio, CLI, composição, testes, README e novos adaptadores de Bitbucket e pesquisa acadêmica. Esses arquivos participam da verificação, portanto os resultados se referem à worktree e não apenas ao commit base.

## Achados e encaminhamentos

| Achado | Evidência local | Encaminhamento |
| --- | --- | --- |
| Propósito explícito de decisão de reúso | Desejos D-001/D-002, README e `reuse_decision.go` | Organizar documentação em torno da jornada principal |
| Contratos e relatos dispersos por evolução histórica | Índices de passos, laudos e relatórios | Criar mapa por capacidade sem reescrever histórico |
| Pureza completa de L1 ainda não corresponde ao código | Relógio em `idea_run.go` e parsers; concorrência/estado no funil | Contratar uma convergência por fatia com metadados temporais separados |
| CLI concentra persistência e projeção do dossiê | `02_shell/cli/idea_recommend.go` | Avaliar ownership e compatibilidade antes de extrair |
| Container durável já tem ownership em L3 | ADR-0007 e `03_infra/runcontainer/container.go` | Preservar lifecycle `Recover` antes de `Start` |
| Composições de descoberta têm alcance diferente | `04_wiring/main.go` e `market_discovery.go` | Documentar ativação efetiva; registro de capacidade não equivale a busca |
| Adições locais estão sem linhagem obrigatória | Lint V1 em oito arquivos novos | Regularizar prompts e cabeçalhos em uma mudança própria |
| Regras de agentes aparecem em dois arquivos | `.agentrules` e `.agents/AGENTS.md` | Conhecer consumidores antes de consolidar |
| Verificação depende de condições da máquina | Loopback bloqueado no sandbox e assinatura nos commits de fixture | Usar ambiente com loopback e configuração efêmera de teste |

## Verificação da worktree

| Verificação | Resultado | Condições |
| --- | --- | --- |
| `go vet ./...` | Passou | Go 1.26.5 instalado; módulo declara Go 1.25.0 |
| Build de `./04_wiring` | Passou | Binário temporário em `/tmp/bateia-organizacao-20261006` |
| `go test ./...` inicial | Falhou | Assinatura de fixtures em CLI/local e socket bloqueado em GitHub |
| Suíte repetida | Passou | Execução com loopback permitido e `commit.gpgsign=false` somente no ambiente do processo |
| `crystalline-lint --fail-on warning .` | Falhou | Oito violações V1 de linhagem causal |
| Fontes executáveis/mascarados em L0 | Nenhum encontrado nas extensões inspecionadas | Varredura dos caminhos e permissões |
| Imports de `bateia/lab` pelo produto | Nenhum encontrado | Varredura de imports das camadas produtivas |
| Diretórios de bancada em L4 | Nenhum encontrado | Inspeção de diretórios |

A suíte aprovada cobre o módulo raiz, incluindo gates ativos alcançados por `./...`. Não representa reprodução do módulo separado de gates selados nem pesquisa remota real dos novos provedores. Alguns pacotes utilizaram resultados em cache.

O lint aponta `02_shell/cli/idea_research.go`, `idea_research_test.go`, `03_infra/bitbucket/connector.go`, `connector_test.go` e, em `03_infra/research`, `common.go`, `joss.go`, `research_test.go` e `zbmath.go`. As violações estavam presentes antes da organização documental.

## Organização realizada

Foi adicionada a [documentação transversal](../documentacao/README.md), com visão do produto, arquitetura implementada, guia de engenharia e sequência proposta de organização. O README principal ganhou um ponto de entrada e alinhamento da descrição de L0 ao ADR-0004. Não houve alteração de código produtivo ou testes nesta etapa.

A próxima fatia recomendada é regularizar a especificação e a linhagem das adições locais. Depois, consolidar um inventário de capacidades e maturidade, começando pela recomendação ponta a ponta. A sequência e os critérios estão no [plano de organização](../documentacao/organizacao.md).
