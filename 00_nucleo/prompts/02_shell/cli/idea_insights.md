# Comando offline de insights

**Camada:** L2  
**Arquivo proprietário:** `02_shell/cli/idea_insights.go`  
**Contrato:** passo 0105

Expor idea insights analyze --input pedido.json [--output resultado.json]. JSON estrito e limitado, rejeição de desconhecidos/duplicatas/valores extras. Proteger input contra sobrescrita/aliases e usar writer atômico/bloqueio de segredos/short writes já existentes na casca. Apenas chamar o caso de uso puro; sem clientes L3, runstore ou política de reúso.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

0134: ler o mesmo pedido estrito uma vez. schema_version 2 seleciona InvestigateSourceInsights; demais versões seguem a validação legada. A casca não valida semântica de observações nem recalcula hashes. Preservar JSON estrito, proteção de input, escrita atômica/segredos e short writes; saída v1 permanece idêntica.

0135: aceitar compare com o mesmo parser/proteções; ler InsightComparisonRequest estrito e delegar a CompareInsightBaseline. Analyze conserva seleção v1/v2 e saídas existentes. Nenhuma regra de baseline na casca.

0136: encaminhar explore para runInsightExploration, com suas flags próprias, antes do parse legado. Preservar analyze e compare sem alterações.

0137: discover pela API sem portas informa necessidade de configuração, sem criar fallback concreto. O executável chama RunInsightDiscoveryCommand com portas injetadas.

## Evolução0142

Aplicar passo0142: memória validada antes de efeitos, followups de um salto com limites e origem própria, --memory/--follow-links na casca e feedback offline. FollowReader apenas GitHub, demais forges reaproveitam RepositoryReader. Memória em snapshot único de saída, sem releitura automática de examinados/descartados, sem promover histórico a evidência atual. Modo sem --sources preservado.


## Evolução0143

Passo0143: expand exige portas configuradas pelo executável como discover; demais ramos offline preservados.

## 0155 — Ideias de algoritmos
Adicionar formulate --input --output pelo parser estrito/proteções já existentes; chamar somente caso de uso L1. Sem efeitos novos ou clientes L3. Preservar analyze/compare.
