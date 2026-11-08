# web_experiment_run

Executar experimento N/F1/F2 em L1: admissão, 12 queries por braço, ordenação natural, capabilities, retomada por porta, budgets/deadlines por runtime, intercalação, selagem, recibo e veredicto. Preservar defaults, 36 pedidos, limites, rate margin, mensagens e ordem. Não ler/gravar arquivos. L2 fornece expansões e leitor de snapshots e publica projeções.
Contrato 0131/ADR-0014.

A descoberta observável simples também recebe expansões por dados e executa consultas em L1; L2 não invoca SearchOne. Preservar schema, hash e estado FAILED para erro do provedor.
