# 0155 — Descobrir princípios e decisões no XERJ

## Intenção e contrato antes da implementação

O Bateia deve ajudar a investigar ideias de algoritmos e outros caminhos, não apenas localizar implementações. Capacidade de investigação sob Tekt: L1 possui contexto e registro de raciocínio, L2 adapta flags/JSON, L3 continua realizando leitura. Não há chamada de modelo embutida nem catálogo fechado de algoritmos.

Modo opt-in algorithm-ideas em explore/discover/collect: preparar perguntas sobre entradas/estado/passos/invariantes, condições/custos/contraexemplos e decisão corrente versus alternativa. Ativa revisão de alternativas com passagens múltiplas. Preservar modos e hashes históricos quando omitido. Tarefas são guias, não achados nem execução de direções de pesquisa.

Novo comando offline formulate recebe alvo congelado, fontes congeladas e ideias interpretadas pelo agente. Cada ideia exige princípio, dois ou mais passos, evidência literal da fonte, condições de aplicação, falhas, decisão atual/proposta/benefício/custo, relação declarada com alvo, evidência do alvo e observável/critério de rejeição. Verificar identidade/revisão/path/hash/linhas/trecho contra os bytes, quotas, cancelamento, cópia e identidade determinística; rejeitar citações inventadas ou papéis trocados. Não classificar semântica como verdade nem afirmar ausência de capacidade pelo recorte. Relação aceita: ALREADY_PRESENT, DIFFERENT_APPROACH, NOT_ESTABLISHED; sempre declaração não reavaliada. Efeito permanece UNVERIFIED; nenhuma adoção ou superioridade automática.

## Teste sobre XERJ

Um pesquisador novo usa Bateia para examinar o XERJ (alvo indicado pelo usuário) e descobrir referências nativas, sem fontes lembradas, histórico ou mecanismos sugeridos. Baseline antes das ferramentas, pin do alvo, consultas/coletas auditadas e requisito de citar comportamento atual. Formular até quatro insights; menos se houver suporte insuficiente. Excluir o que já existe ou registrar como confirmação, sem anunciar melhoria nova. A investigação deve expor princípio, evidência, condição e decisão mudada; execução operacional e interpretação são separadas.

Revisor independente verifica cada cadeia e classifica utilidade/novidade relativa, condições e lacunas, sem reescrever o resultado para produzir vantagem. Dados de bancada e eventual experimento limitado ficam no Lab. Nenhuma alteração no XERJ, PR, publicação ou email enviado. Gate canônico make verify; relatório e texto utilizável no email ao Ivan com limites explícitos.

## Resultado

Materialização e gate canônico concluídos. O ensaio no alvo XERJ registrou uma hipótese condicional de formato/acesso e uma confirmação ALREADY_PRESENT de WAND. Quinze vínculos de fonte/alvo foram conferidos em retornos e Git; a revisão encontrou cache FTS existente e restringiu a aplicabilidade da hipótese. Nenhuma fonte foi executada e nenhum ganho medido. Parecer, dados e redação utilizável no email na [bancada 0155](../../lab/bancada/insights/algoritmos-xerj-2026-11-08/README.md).
