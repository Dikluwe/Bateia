# 0103 Organização de responsabilidades sob Tekt

**Data**: 2026-10-06  
**Recorte**: worktree com base em `9d00e2f`, incluindo adições locais anteriores  
**Estado**: primeira fatia verificada; avaliação geral da hierarquia permanece em andamento

## Resultado

A organização passa a ligar responsabilidade, camada, código, contrato e verificação no [mapa arquitetural](../arquitetura/README.md). `AGENTS.md` na raiz fornece uma entrada para agentes e `.agentrules` incorpora o propósito de reúso e insights, preservando as fronteiras Tekt. Os escopos A1 continuam organizações candidatas, avaliadas por seus contratos e propriedades.

No código, `improvement_analysis.go` passa a possuir a derivação e as projeções de hipóteses, lacunas, planos, proveniência e estado de negócio. Recebe pedido normalizado e mecanismos validados, com sinal de degradação, sem store, fila ou serviço de runs. `improvement_synthesis.go` conserva normalização, integridade e resolução de fontes pelas portas, admissão, lifecycle, cancelamento e publicação. O estado de negócio vem da análise; sucesso operacional não promove a hipótese.

Ambos permanecem no mesmo pacote L1. A fronteira é de responsabilidade e contrato, não de isolamento de imports entre pacotes. Formatos públicos, envelopes, IDs, defaults, ordem de transições e publicação foram preservados no recorte verificado.

## Verificação

| Verificação | Resultado | Condições |
| --- | --- | --- |
| Caracterização antes/depois | Passou | Digest da projeção capturado antes da extração e comparado após a separação |
| Cenários da análise | Passaram | Pedido não mutado, degradação, determinismo de mapas, zero hipóteses e abstenção de adoção |
| `crystalline-lint --fail-on warning .` | Passou sem violações | Inclui os prompts/cabeçalhos regularizados das adições locais |
| `go test ./...` | Passou | Loopback permitido; assinatura Git desativada somente no processo de fixtures; alguns pacotes em cache |
| `go vet ./...` | Passou | Execução autorizada para usar o cache Go existente |
| Build de `./04_wiring` | Passou | Binário temporário `/tmp/bateia-arquitetura-0103` |
| Links e whitespace | Passaram | Links locais dos novos artefatos e `git diff --check` |
| Fronteiras inspecionadas | Sem ocorrências proibidas no recorte | Nenhuma fonte/executor em L0 ou import de `bateia/lab` nas camadas produtivas |

O vet e build iniciais encontraram cache Go somente leitura no sandbox; a execução autorizada com o mesmo cache passou. Não foi alterada a configuração Git persistente. O módulo separado de gates selados não foi reproduzido.

A [bancada do corte](../../lab/bancada/arquitetura/0103/README.md) registra o baseline e seu teste de reprodução. A caracterização é de um caso congelado; outros comportamentos são cobertos pelos cenários e contratos ativos. Não constitui prova universal de equivalência.

## Regularização de adições locais

Foram adicionados oito prompts proprietários e cabeçalhos de linhagem à CLI de pesquisa, seus testes, Bitbucket e pesquisa JOSS/zbMATH. Os contratos foram escritos após leitura da implementação existente e não são apresentados como sua causa histórica original. Nenhum comportamento desses conectores foi alterado nesta regularização.

## Limites e próximos cortes

- Insights ainda dependem de fontes produzidas por runs de recomendação; independência de ponta a ponta não foi implementada.
- Runtime ainda acumula configuração e coordenação de diferentes capacidades. O primeiro corte retira a derivação e classificação de negócio, sem generalizar todo o executor.
- Relógio em L1 e I/O de dossiês na CLI permanecem pendentes de contrato próprio.
- A composição atual de síntese exige exatamente três conectores, enquanto o registro padrão da worktree pode trazer cinco ou mais. É um conflito preexistente de composição registrado no mapa, não resolvido ou coberto por um teste de startup nesta fatia.
- Não foram relaxadas triagem, política de reúso ou força probatória para produzir mais hipóteses.
- A hierarquia não foi avaliada por tarefas comparáveis de IA. Organização mais explícita não demonstra, por si, ganho de produtividade ou qualidade.

A próxima fatia deve tratar o conflito de composição ou reduzir a resolução/publicação misturada às capacidades do runtime, com prompt e teste de preservação específicos. A árvore de diretórios não foi reestruturada em massa.
