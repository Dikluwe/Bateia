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
