# 0137 — Descobrir referências diversificadas automaticamente

**Data:** 2026-10-28. **Contrato ex-ante.** Autorização: usuário pede descoberta diversificada automática para auxiliar a IA a encontrar abordagens fora do habitual. Investigação e políticas em L1; flags/saída em L2; efeitos GitHub existentes L3; somente composição em L4.

## Superfície

`idea insights discover --question TEXTO [--familiar TEXTO ...] --output arquivo.json [--max-per-query N] [--max-references N] [--max-passages N] [--timeout duração] [--attempt-timeout duração]`. Output obrigatório antes de efeitos. Executável injeta ProjectCatalog e porta estreita de leitura de árvore/arquivos do cliente GitHub existente. Nenhuma instanciação concreta em L1/L2, novo runtime ou alteração dos comandos anteriores. A API CLI sem portas informa configuração ausente; a API pública injetada permite fixtures/integradores.

## Plano e budgets

Pedido v1: pergunta (1000 bytes), até oito rótulos familiares (120 bytes), máximo por consulta 1–5 (padrão 2), referências 1–8 (padrão 4), trechos 1–16 (padrão 8), total 1s–5min (padrão 90s), tentativa 1s–30s (padrão 10s, até total). Normalizar como 0136; rejeitar pergunta sem termos lexicais aproveitáveis antes de efeitos. Runtime explícito por contexto; sem relógio/I/O concreto em L1.

Plano determinístico com cinco rotas independentes, derivação LEXICAL_EXPLORATION_HEURISTIC: DIRECT (até três primeiros termos da pergunta), BROADENED (primeiro termo fora dos rótulos familiares; fallback primeiro termo), ALTERNATIVES (esse termo + alternative), TRADEOFFS (último termo fora dos rótulos familiares; fallback último termo + tradeoffs), COUNTEREXAMPLES (primeiro termo fora dos rótulos familiares + failure). Usar tokens de 0136, removendo dos termos exploratórios tokens dos rótulos familiares; a consulta direta preserva o caminho habitual. Sem tradução, sinônimos semânticos, negação automática de familiaridade ou alegação de outro domínio. Queries iguais podem existir com rotas distintas e conservam identidade/proveniência próprias. Esses sufixos são heurísticas de pesquisa, não explicações de mecanismos encontrados.

Executar sequencialmente todas as rotas enquanto houver tempo, mesmo após sucesso da direta; uma chamada SearchCatalog por rota, sem retry oculto no core. Timeout total e por operação via porta de runtime; um máximo de cinco chamadas lógicas de catálogo e oito coletas. Adaptador GitHub pode fazer requisições internas de metadados sob o mesmo deadline; a quota não significa número exato de requisições HTTP. Cancelamento impede novas chamadas e conserva resultados já admitidos. Falha por rota não impede as restantes dentro do orçamento.

## Resultados e seleção

Conferir papel PROJECT_CATALOG, identidade do provedor, query da tentativa, cardinalidade e candidatos com owner/repo válidos. Não admitir candidatos de resposta inválida; registrar INTEGRITY_ERROR. Falha/empty/rate-limit/timeout são estados de operação, sem inferência de ausência de abordagens. Não copiar mensagens arbitrárias do provedor: somente código de erro seguro. Preservar contagens e query/rota/ID/attempt. Ordenar candidatos por identidade; desduplicar owner/repo sem caixa, unir origens e registrar conflitos entre commits imutáveis observados. Preferir uma projeção canônica de metadados quando duplicados, sem resolver conflito de revisão por conveniência.

Selecionar até max-references pelo ganho de consultas ainda não representadas, autor ainda não representado e menos rótulos familiares encontrados literalmente em nome/descrição/tópicos; desempate por owner/repo. Estrelas e popularidade não decidem seleção; match familiar não exclui contraexemplo. Registrar razões e rotas sem referência selecionada. Cobertura de queries/autores é diversidade estrutural; não classificar relevância, novidade, equivalência ou domínio semântico.

## Coleta limitada e contexto

Para cada selecionado, exigir commit hexadecimal imutável sem conflito antes de ler. Não usar HEAD ou resolver silenciosamente uma revisão conflitante. Ler a árvore pelo provedor existente; escolher primeiro README de raiz regular segundo prioridade README.md, README.rst, README.txt, README, ignorando caixa e com desempate lexical. Exigir path relativo válido. Extrair só esse arquivo, conferir ref/identidade/path/cardinalidade da resposta, UTF-8 textual e até 32 KiB; hash dos bytes completos. Não truncar, clonar ou executar. Os limites do adaptador continuam aplicáveis a envelopes HTTP; o teto de 32 KiB é de admissão após leitura.

Produzir bundles v1 com escopo ROOT_README_ONLY e complete/tree_complete false, pois não representam auditoria do projeto. Registrar uma entrada de cobertura `document`/`root README`, regra `insight-readme-v1`, TARGETED_EXAMINED, com o path examinado; isso registra leitura integral do documento escolhido, sem avaliação de capacidade. A entrada explícita também evita a alternância legada de lista vazia/null em requirement_coverage no round-trip JSON, preservando hashes dos contratos antigos sem alterá-los. Recebimentos distinguem COLLECTED, OPAQUE_REVISION, REVISION_CONFLICT, README_NOT_FOUND, READ_FAILED, INTEGRITY_ERROR, CONTENT_LIMIT, CANCELLED/NOT_STARTED. Falha de uma coleta não apaga candidatos nem ativa coletas extras além da seleção. Alimentar ExploreInsights com os bundles admitidos, pergunta e familiaridade, gerando cartões/tarefas sem curadoria manual. Reportar a coleta limitada e incompletude, sem promover catálogo a evidência de código executado.

## Relatório e publicação

Relatório insight-discovery-v1 contém pedido/hash/plano estável, capacidades do catálogo, tentativas, candidatos com origens/seleção/matches, recebimentos, bundles, contexto de exploração, lacunas e limitações. ID do relatório inclui observações/recibos e contexto; live pode mudar metadados/durações. Plano mantém ID determinístico para o pedido normalizado.

Estados: REFERENCES_COLLECTED quando operações de busca/coletas selecionadas concluíram sem falha e há leitura; PARTIAL se restam falhas/lacunas operacionais após alguma rota válida, inclusive quando outras buscas falharam e nenhum candidato pôde ser admitido; NO_CANDIDATES para todas as buscas válidas vazias; FAILED se nenhuma rota válida concluiu; CANCELLED/TIMED_OUT se interrompido. Efeito UNVERIFIED e novidade NOT_ASSESSED sempre. Preservar resultados parciais antes de retornar erro de interrupção/falha total. Input inválido/configuração ausente não publica nem chama portas. L2 publica JSON protegido atomicamente, bloqueia segredos/short writes, só depois devolve erro operacional seguro.

## Aceite

Testes puros: rotas independentes, todas executadas após sucesso, determinismo/causalidade/imutabilidade, conflitos/duplicatas/ordem/popularidade ignorada, quotas/identidades/deadline/cancelamento, resultados parciais e integridade dos READMEs. CLI com portas falsas: validação antes de efeitos, output obrigatório/protegido, flags repetidas/posições, segredos e publicação parcial. Integração HTTP em L3 com cliente GitHub real/servidor local e wiring de composição; fixture fora de L1/L2. Ensaio online público limitado se disponível, conservando inclusive falhas e quotas, sem executar candidatos. Replays de 0106/0134/0135/0136 intactos. Atualizar documentação, linhagem e laudo; executar make verify, registros novos em lab/bancada/insights/0137.
