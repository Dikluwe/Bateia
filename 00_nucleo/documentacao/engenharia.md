# Guia de engenharia do Bateia

Uma mudança no Bateia deve preservar a relação entre intenção, contrato, implementação e evidência. A unidade de trabalho recomendada é uma capacidade observável, com limites e critérios de aceite definidos antes da alteração.

## Ciclo de uma mudança

1. Localize a capacidade na [arquitetura](arquitetura.md), leia seu prompt e os testes de contrato e identifique consumidores CLI/MCP.
2. Declare a incerteza, comportamento esperado, falhas, observável e fora de escopo em um passo. Atualize a especificação proprietária em `00_nucleo/prompts` antes de materializar uma evolução relevante.
3. Registre um ADR quando mudar fronteiras, ownership, formatos ou uma decisão transversal. Não crie um ADR apenas para narrar uma edição rotineira.
4. Implemente uma fatia pequena. Mantenha transformação/regras em L1, efeitos nos adaptadores e composição em L4; respeite o comportamento vigente de L2 até uma migração contratada.
5. Preserve o cabeçalho `Crystalline Lineage`, com `@prompt`, `@layer` e `@updated`, em cada arquivo Go criado ou modificado.
6. Execute as verificações e relate exatamente seu resultado. Documente limites no laudo ou síntese; resultados brutos e fixtures pertencem ao Lab.
7. Atualize a documentação transversal e os exemplos afetados. Diferencie disponibilidade de código, aprovação de testes e utilidade demonstrada.

Essas obrigações derivam de [.agentrules](../../.agentrules). A organização documental desta pasta não dispensa a tríade prompt, código/testes e lint nas mudanças de implementação.

## Verificação técnica

Na raiz, execute `make verify` para a pilha completa ou `make lint` para o lint. O gate usa a fonte/patch fixados em `lab/ferramentas/crystalline-lint/manifest.json`, compilados e validados fora do produto; não exige checkout externo nem altera a instalação global. Pré-requisitos e cache estão em [lab/ferramentas/README.md](../../lab/ferramentas/README.md). Os comandos individuais da pilha são:

```bash
crystalline-lint --fail-on warning .
go test ./...
go vet ./...
go build -o /dev/shm/bateia ./04_wiring
```

`make architecture` e a etapa arquitetural de `make verify` incluem gravidade/efeitos (0129), ownership anterior (0131) e a guarda da investigação (0145). Esta última verifica proprietários, controles negativos e preservação dos corpos dos consumidores realocados. Gates selados continuam no regime separado descrito abaixo.

Em um ambiente que não permita escrever em `/dev/shm`, use um destino temporário autorizado e registre a substituição. Go deve atender ao `go.mod`; integrações opcionais como Clang, ferramenta de fluxo e govulncheck dependem da capacidade exercitada.

O módulo raiz inclui gates ativos do Lab na expansão `./...`. Os gates selados em `lab/bancada/gates/sealed` possuem outro `go.mod`; não ficam automaticamente cobertos por esse comando. Sua reprodução deve seguir o README e o commit registrado no próprio gate.

Testes que usam `httptest` precisam de permissão para escutar em loopback. Fixtures que criam commits podem herdar assinatura Git obrigatória da máquina. Para essas fixtures, desative assinatura somente no processo de teste:

```bash
GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=commit.gpgsign GIT_CONFIG_VALUE_0=false go test ./...
```

Esse comando não altera configuração persistente. Se `GIT_CONFIG_COUNT` já estiver em uso, acrescente a opção ao conjunto existente. Diferencie falha do produto de bloqueio do ambiente, mantendo o comando original e o resultado ajustado no registro de verificação.

Para análise C/C++ local, Clang procura `compile_commands.json` na raiz, em
`build/` e depois em `../build/`. O primeiro banco existente governa o contexto:
erros de leitura ou JSON não autorizam fallback. O limite é 16 MiB. Entradas
`command` do Meson são tokenizadas como argumentos literais; somente o Clang
configurado é executado em modo estático, após filtrar opções. Shell, plugins
e response files continuam recusados. Caminho/hash do banco e hashes de
fontes/headers acompanham as unidades. Seleção limitada a 12 unidades não
comprova cobertura ou comportamento. Consulte o
[passo 0125](../passos/0125-investigar-clipping-e-transicao-gs.md).

## Estratégia de testes

| Nível | O que verificar |
| --- | --- |
| Domínio e transformações | Invariantes, entradas inválidas, fronteiras e desconhecidos |
| Casos de uso | Decisões, falhas parciais, cancelamento e orçamentos com portas substituídas |
| Adaptadores | Protocolos, identidade, limites de leitura, erros e persistência |
| CLI e MCP | Formatos públicos, defaults, validação e projeção fiel |
| Composição | Serviços corretos, features configuradas e ordem de lifecycle |
| Avaliação de produto | Relevância, cobertura de evidência, decisão e utilidade em corpus declarado |

Use testes focados durante a implementação e a pilha completa antes de declarar conclusão. Um teste que apenas espelha a estrutura do código não comprova preservação de comportamento. Testes históricos que fixam localização por texto devem ser avaliados diante de ADRs, sem alterar recibos antigos.

## Revisão de contratos e segurança operacional

Revise schema, mensagens, defaults, ordenação, identidade do alvo, hashes e política de abstenção. Falha de descoberta, leitura incompleta ou identidade opaca não podem virar inexistência comprovada por conveniência.

Credenciais ficam no ambiente ou no `.env`, nunca em documentação, fixtures ou relatórios. Não versione caches e `.bateia`. Para novos provedores, explicite configuração, escopo pesquisável, limitações, identidade, proveniência e distinção entre zero resultados e falha.

Antes de integrar, verifique fronteiras de importação, ausência de executores e fontes mascarados em L0, ausência de bancada em L4 e isolamento do Lab. Um lint aprovado verifica suas regras conhecidas; não prova toda a semântica dessas fronteiras.

## Critério de conclusão

Uma evolução está pronta quando entrega o observável contratado, preserva consumidores, apresenta verificação reproduzível e deixa claros seus limites. A utilidade de uma capacidade exige sua avaliação própria; não deve ser inferida apenas de compilação ou testes de infraestrutura.

## Reprodução histórica das fronteiras — passo 0129

Para reproduzir especificamente o ensaio histórico 0129, use `bash lab/bancada/arquitetura/0129/reproduce-lint.sh` (o gate atual `make verify` já inclui fonte/patch fixados no projeto): ele extrai a revisão declarada do linter para `/tmp`, aplica o patch do parser, executa os testes da ferramenta e 11 controles adversariais e roda `crystalline-lint --fail-on warning .` com o binário validado no PATH somente desse processo. A fonte externa e a instalação global permanecem intactas. Requer Rust/Cargo, crates em cache e o repositório do linter na revisão disponível; o caminho pode ser informado por `BATEIA_LINTER_SOURCE`.

Execute também `go test -v ./lab/bancada/arquitetura/0129/boundaries`. O auditor verifica efeitos concretos, gravidade produtiva, contratos em L4, isolamento do Lab e executáveis em L0, com controles negativos. Não substitui avaliação semântica das regras de negócio.

Ao compor a CLI como biblioteca, instale `cli.ConfigureCommandHost(host.New())` no startup L4. Configure `WithExecutionRuntime` nas fachadas do funil/recomendação, ou use contexto com runtime nas chamadas L1. Não há driver implícito em L1/L2. Os testes externos `cli_test` compõem o mesmo adaptador sem inseri-lo no build produtivo. O seam de configuração da CLI é de startup; não deve ser alterado durante comandos concorrentes.

## Conclusão dos pontos da revisão — 0131

A jornada do dossiê, identidade, manifesto e estágios pertencem a `recommendation_dossier*.go` em L1. A CLI fornece adaptadores de nomes lógicos e apresentação, recebendo o serviço composto em L4; o caminho compatível da biblioteca delega à mesma jornada. A investigação web e o veredicto pertencem a `web_experiment_run.go`; a investigação de mercado implementa a porta de executor em `market_execution.go`, enquanto o runtime publica e controla CAS/lifecycle.

Os testes puros de códigos/locale ficam em L1 e o auditor AST geral em `lab/bancada/arquitetura/0131/contracts`. Resultados brutos Clang ficam em `lab/bancada/pesquisa/ast-cfg`. O gate inclui testes de ownership e recusa material não declarado em L0, além da auditoria de efeitos/imports. Relatórios históricos de 0129/0130 descrevem seus próprios recortes, não a árvore posterior.
