# web_experiment_run

Executar experimento N/F1/F2 em L1: admissão, 12 queries por braço, ordenação natural, capabilities, retomada por porta, budgets/deadlines por runtime, intercalação, selagem, recibo e veredicto. Preservar defaults, 36 pedidos, limites, rate margin, mensagens e ordem. Não ler/gravar arquivos. L2 fornece expansões e leitor de snapshots e publica projeções.
Contrato 0131/ADR-0014.

A descoberta observável simples também recebe expansões por dados e executa consultas em L1; L2 não invoca SearchOne. Preservar schema, hash e estado FAILED para erro do provedor.

## Retomada íntegra e ligada ao plano — 0170

Antes de SearchOne, admitir todos os snapshots da retomada autenticada com schemas1/2, ReplayWebSnapshot válido e sem MissingFields. Conferir RunID, Arm, Provider, GrammarHash e QuerySetHash contra a execução planejada. Schema2 exige Authorization e Budget iguais aos atuais; schema1 legado pode omiti-los, mas valores presentes devem coincidir. Não reparar selos ou traduzir Unknown em EMPTY.

Cada resposta deve ter QueryID único conhecido e request exatamente igual à planejada: RequestID0020-%02d-ARM pela posição natural, QueryID, Query, MaxResults efetivo e Locale en-US; conferir Provider, Arm, Position e AttemptNumber=1. Selo recalculado não dispensa vínculo ao plano. Subconjunto válido pode retomar SUCCESS/EMPTY; queries ausentes ou desconhecidas continuam buscando conforme o contrato anterior. Admissão inválida retorna erro de integridade/reprodutibilidade e output zero antes de novas buscas. Validar braços N/F1/F2 únicos, IDs de query não vazios/únicos por braço e Text não vazio. Copiar Queries antes da ordenação. Preservar entrada, sequência intercalada, contadores, orçamento, estados e hashes do controle íntegro.
