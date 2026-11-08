# CLI de rastreabilidade consultiva — passo 0098

Implementar `bateia idea trace build --input build.json [--output graph.json]`
e `bateia idea trace impact --input change.json [--output plan.json]`.
Usar contratos domain e funções puras do core. Inputs schema_version 1;
grafo interno decision-trace-v1. A saída descreve dependências declaradas,
nunca valida autenticidade, recalcula decisão, executa candidato ou acessa rede.

Flags estritas: input obrigatório único, output opcional único, valores não vazios,
sem posicionais ou parâmetros desconhecidos. Ler com readStrictIntakeJSON para
aplicar os mesmos limites, recusa de FIFO, chaves duplicadas e campos desconhecidos.
Input inválido não pode produzir stdout nem criar/substituir output.

`validateDecisionTraceOutput(output, inputs...)` também é utilizado pelo fluxo
avançado: recusar colisões por path, inode e aliases contra qualquer input e
recusar links simbólicos em output ou seus ancestrais. Diretório pai deve existir;
output preexistente deve ser arquivo regular. Substituição explícita de resultado
regular distinto é permitida. Não alterar inputs.

Admissão antecipada pode aceitar pais e destinos protegidos ainda ausentes;
resolver o prefixo existente dos destinos protegidos para identificar aliases
mesmo quando leaf não existe. A escrita continua exigindo pai existente.
Links pendentes em destinos protegidos são recusados, nunca reclassificados como
folha ausente. Destinos protegidos ausentes com componente '..' são recusados
por ambiguidade após futura criação de diretórios; paths existentes continuam válidos.

O writer global writeIdeaJSON usa os.WriteFile e não é atômico neste passo.
O writer específico trace serializa completamente, limita a 16 MiB, fixa o pai
com os.Root, cria temporário aleatório exclusivo, sincroniza e renomeia no root.
Limpar temporários em todos os caminhos. Revalidar leaf/inode antes do rename.
Comparar identidade do parent antes/depois de OpenRoot, e refazer validação de
paths antes do rename. Aplicar rejectSerializedSecrets antes de stdout ou arquivo.
Não declarar proteção contra toda mutação concorrente do namespace externo.

Testar caminhos positivos build/impact, stdout versus arquivo, rejeições de flags,
JSON desconhecido/duplicado/ambiguidades, FIFO, aliases, adulteração de graph_hash,
preservação de inputs e ausência de efeitos em falhas.
## Explain 0099

Aceitar também `idea trace explain --input request.json [--output index.json]`.
Strict decode de DecisionTraceExplainInput schema1 seguido de ExplainDecisionTrace.
A requisição contém grafo congelado e build_input exato; seleção node_ids opcional.
Recusar mismatch, hash ausente/adulterado e IDs desconhecidos/duplicados antes de
stdout/arquivo. Reutilizar admissão/writer dos outros trace comandos. Documento
report/idea/rule_versions e JSON Pointer são declarações indexadas, não paths a
buscar ou provas de autenticidade. Sem recalcular/promover decisão ou executar IO
no núcleo.

## Compare 0100

Aceitar `idea trace compare --input request.json [--output comparison.json]`.
Decodificar estritamente DecisionTraceCompareInput schema1 com before e after,
ambos grafos congelados; usar CompareDecisionTraces sem IO no núcleo.
Reutilizar o parser, admissão antecipada, reader limitado e writer atômico trace.
Hashes ausentes/adulterados ou grafos inválidos são recusados antes de qualquer
stdout/arquivo. O resultado distingue historical_decision e candidate_decision:
esta última é apenas o veredito fornecido no segundo snapshot, não aprovação.
decision_impact permanece NONE e o plano de impacto conserva a decisão histórica.
Comparar declarações/arestas registradas não comprova mudança real, revogação,
autenticidade ou suficiência das evidências. Sem rede, execução de candidatos,
releitura de fontes, carregamento de providers/env ou alteração de recommend e
artefatos default.

## Review 0101

Aceitar idea trace review --input pedido.json [--output registro.json]. Entrada
DecisionTraceReviewInput schema1: impact_input e records opcional. Reutilizar
TrackDecisionTraceReview, leitor JSON estrito e writer atômico protegido. Recusar
hashes/estados/recibos inválidos sem stdout ou substituição de arquivo. Sem rede,
.env, provedores, execução ou promoção de decisões. Busca simples inalterada.

## Review-context 0102

Aceitar idea trace review-context --input pedido.json [--output contexto.json].
Entrada DecisionTraceReviewContextInput schema1 com review_input e build_input.
Usar BuildDecisionTraceReviewContext, leitor estrito e writer protegido/atômico.
Não aceitar índices/planos arbitrários nem promover COMPLETED. Recusar documentos
que não correspondem ao grafo congelado, inclusive sem tarefas. Sem rede/env,
execução ou alteração da pesquisa simples.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
