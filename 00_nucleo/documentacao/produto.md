# Visão do produto Bateia

O foco atual do Bateia é ajudar quem vai iniciar ou evoluir um software a obter insights de engenharia investigando soluções existentes. Isso inclui estratégias observadas, alternativas, condições, trade-offs e hipóteses de melhoria com planos de validação. O [passo 0133](../passos/0133-priorizar-insights-de-engenharia.md) registra a prioridade reafirmada pelo usuário. Reúso permanece como capacidade e possível resultado, com jornada e políticas próprias; as dependências e os limites atuais estão no [mapa arquitetural](../arquitetura/README.md).

A intenção está registrada em [D-001](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [D-002](../desejos/D-002-alternativas-de-software-com-evidencia.md).

## Resultado prioritário desejado

A investigação deve ajudar a responder: como outros projetos tratam este problema e o que podemos aprender para nosso contexto? A entrada orientadora é uma pergunta, acompanhada de alvo e restrições quando disponíveis. O resultado desejado contém a estratégia concreta, suas fontes, condições necessárias, limites, diferenças entre abordagens e um próximo passo verificável. Uma observação de incompatibilidade também pode ensinar algo útil.

Observação da fonte, interpretação, hipótese de benefício e efeito medido no alvo conservam escopos distintos. IA pode auxiliar consultas e interpretação, com referências conferíveis e hipóteses identificadas; não funciona como confirmação automática. A investigação não depende de uma decisão de adoção e não muda os critérios de licença/reúso existentes.

Essa é a direção de produto. Busca automática, comparação semântica geral e validação de aplicabilidade ainda não são entregas completas da jornada de insights atual.

O usuário explicitou que a ferramenta deve auxiliar a IA a sair do caminho habitual. O [recorte 0136](../passos/0136-preparar-exploracao-criativa-para-ia.md) prepara cartões de leitura e tarefas de exploração diretamente da pergunta e de bundles coletados, sem assessments/observations. A IA recebe exemplos referenciáveis e perspectivas de alternativas, inversão, analogia, recombinação e contraexemplos. A escolha automática usa variedade estrutural e matches lexicais; não mede novidade ou criatividade, não traduz perguntas e não gera conclusões por modelo. Sem evidência, o pacote pede pesquisa e conserva as lacunas.

O [recorte 0137](../passos/0137-descobrir-referencias-diversificadas-automaticamente.md) acrescenta `idea insights discover`: automatiza cinco consultas independentes no catálogo GitHub, seleção por cobertura de consultas/autores e coleta limitada de READMEs fixados, alimentando explore no mesmo resultado. As origens, falhas e limites de cada etapa ficam explícitos. Esse primeiro provedor não representa busca em todos os conectores; as heurísticas lexicais não comprovam diversidade conceitual. O [ensaio público](../../lab/bancada/insights/0137/README.md) produziu contexto sem curadoria de avaliações, conservando um documento rejeitado por tamanho e o estado parcial.

O [recorte 0138](../passos/0138-ampliar-descoberta-de-insights-para-outros-sitios.md) permite `--sources all` ou lista de sítios, reutilizando os conectores de forges, lojas e software acadêmico. Inclui Gitee e AtomGit, mantém orçamento global e diferencia falha/autenticação de vazio. A IA recebe pistas de catálogo com identidade/proveniência por sítio, separadas das leituras de código/documentos. O recorte0138 começou com leitura imutável GitHub; a [evolução0139](../passos/0139-ler-readmes-de-gitee-e-atomgit.md) acrescenta Gitee/AtomGit com resolução de commit e contextos de origem separados. O [recorte0140](../passos/0140-ler-readmes-de-codeberg-e-gitlab.md) amplia esses leitores para Codeberg/GitLab, incluindo subgrupos do GitLab, pelas APIs públicas. A IA recebe os cartões desses contextos nas mesmas tarefas; outros catálogos não ganham SHA presumido. Novidade/relevância ainda precisam de interpretação e validação, e o ensaio registra ruído e fontes indisponíveis.

## Investigação independente de insights

O recorte [0105](../passos/0105-investigar-insights-sem-decisao-de-reuso.md) recebe pergunta, alvo e evidências congeladas, sem depender de IdeaSpec ou decisão de reúso. Produz mecanismos vinculados aos bytes fornecidos, hipóteses conservadoras e plano de validação. As avaliações são declaradas e não refeitas; efeito e aplicabilidade no alvo permanecem desconhecidos. Consultar o [exemplo ilustrativo](../../lab/bancada/insights/0105/README.md).

O [ensaio 0106](../laudos/0106-correcoes-e-ensaio-github.md) registrou detalhes de uma implementação real numa revisão manual. O [recorte 0134](../passos/0134-conservar-observacoes-especificas-nos-insights.md) conserva esses detalhes quando fornecidos em pedidos de versão 2: cada observação declara tipo, texto e referências por arquivo/hash/linhas, com trechos conferíveis na saída. Condições operacionais e trade-offs da hipótese vêm dessas declarações; linguagens e ecossistemas ficam no escopo do extrator. A definição genérica do catálogo é identificada separadamente.

Isso ainda exige observações e avaliações fornecidas: conferir bytes/localizadores não reavalia a interpretação. Sem mecanismo classificado, as observações permanecem INCONCLUSIVE e não geram hipóteses por si sós; sem observações, a lacuna é explícita. O [exemplo 0134](../../lab/bancada/insights/0134/README.md) usa os bytes congelados de `conc`, sem nova coleta ou execução da fonte. A versão 1 permanece reproduzível.

O [recorte 0135](../passos/0135-comparar-insights-com-baseline-do-alvo.md) acrescenta `idea insights compare`: observações das fontes e de um snapshot do alvo ficam agrupadas por classe/tipo, com referências separadas e perguntas para validar diferenças. O programa indica se há declarações de ambos os lados, de um lado ou de nenhum; isso não é veredito de equivalência, presença de capacidade ou ganho. Considerar a política já declarada no alvo precede a proposta de alteração. A baseline cobre os arquivos fornecidos, admite conteúdo não commitado e conserva a origem declarada. Consultar o [caso conc × Bateia](../../lab/bancada/insights/0135/README.md).

## Entrada e resultado

A jornada implementada de reúso recebe uma `IdeaSpec` com problema, capacidades e restrições. A política de reúso determina critérios como licenças admitidas e limites de composição. Um alvo opcional permite comparar alternativas de um software conhecido.

O resultado reúne descoberta, triagem, planos de coleta, bundles de evidência, avaliação, alternativas e decisão. O dossiê deve permitir localizar a origem das afirmações, identificar lacunas e compreender por que a decisão foi tomada.

| Decisão | Significado |
| --- | --- |
| `FORK` | Usar um projeto como base de uma derivação |
| `EXTEND` | Ampliar uma base existente |
| `COMPOSE` | Combinar capacidades de vários componentes |
| `INTEGRATE` | Incorporar um componente à solução |
| `REFERENCE` | Aproveitar conhecimento como referência |
| `BUILD` | Construir sob as condições e evidências avaliadas |
| `INCONCLUSIVE` | Evidência insuficiente para uma decisão sustentada |

Os critérios executáveis pertencem ao [caso de uso de decisão](../../01_core/usecase/reuse_decision.go). Estas descrições não substituem seus limiares e precondições.

## Capacidades e papel no produto

| Capacidade | Papel |
| --- | --- |
| Expansão lexical e funcional | Transformar intenção em consultas rastreáveis |
| Descoberta de repositórios, produtos e software acadêmico | Encontrar candidatos com proveniência |
| Triagem e calibração | Priorizar a investigação e tornar perdas observáveis |
| Coleta, recuperação e avaliação de evidências | Verificar capacidades e explicitar desconhecidos |
| Decisão e ranking de alternativas | Relacionar evidência às opções de reúso |
| Rastreabilidade e revisão | Localizar declarações e planejar sua revisão |
| CLI, MCP e execuções duráveis | Tornar a jornada acessível a pessoas e agentes |

## Limites do que se pode concluir

Correspondência lexical, presença em catálogo, dependência em manifesto e observação sintática têm forças probatórias diferentes. Uma biblioteca declarada não comprova uma capacidade; um símbolo encontrado não comprova comportamento executado. A ausência de um resultado em uma busca limitada não demonstra inexistência.

Há parsers de manifestos de muitos ecossistemas, extração AST de Go e análise opcional C/C++ por Clang. Isso não implica cobertura semântica equivalente entre linguagens. Os registros acadêmicos trazem sinais de catálogo/revisão com escopo próprio; o vínculo ao código precisa de verificação adicional.

Os modos `local`, `live` e `replay` têm fontes e condições distintas. Replay permite reprocessar artefatos congelados; não autentica automaticamente a fonte nem promete igualdade de tempos e metadados operacionais. O progresso declarado de revisão também não altera ou aprova a decisão histórica.

Ao concluir `idea recommend --resume`, o diagnóstico informa se a evidência anterior foi reutilizada e o primeiro impedimento observado quando recusada. Distingue manifesto indisponível, mudanças causais, versões e artefatos ilegíveis/divergentes/sem registro. A mensagem não expõe valores ou erros internos e não altera a decisão, a política de aceitação ou a identidade causal. Sem flag, a saída permanece como antes. O [passo 0132](../passos/0132-explicar-retomada-de-evidencias.md) delimita essa melhoria de clareza.

## Como avaliar utilidade

A qualidade depende de encontrar candidatos relevantes, comprovar capacidades, explicar bloqueios e ajudar a escolher um ponto de partida. Testes de infraestrutura e quantidade de arquivos extraídos são sinais técnicos, mas a avaliação de produto exige ideias reais, julgamentos independentes e cobertura declarada.

Para insights, avaliar primeiro relevância à pergunta, especificidade da estratégia, fontes conferíveis, condições/contraexemplos e utilidade do próximo passo. Um relatório não precisa autorizar adoção para ser útil; precisa reduzir o trabalho de investigação ou permitir uma escolha/experimento melhor. Registrar tempo, trabalho restante e julgamentos antes de atribuir ganho. A aprovação de testes técnicos não demonstra esse ganho.

A aprovação de uma etapa não promove automaticamente as demais. O [índice de passos](../passos/README.md) registra, por exemplo, descoberta aprovada com validação downstream pendente no passo 0023.

## Ampliação de possibilidades

O usuário reafirmou que o critério central é ampliar possibilidades que a IA considera, inclusive após pesquisas anteriores. O [passo0141](../passos/0141-ampliar-possibilidades-com-perspectivas-e-trilhas.md) distribui operações entre perspectivas/fontes e extrai trilhas documentais para repositórios citados. Oportunidade de consulta, link citado, conteúdo lido e hipótese continuam distintos. O caso congelado da coleção de museus produz caminhos adicionais como Serendip-o-matic com a origem exata, sem presumir mecanismo implementado, relevância ou novidade. A IA recebe as trilhas em suas tarefas; a leitura do destino continua uma etapa posterior, sem rede adicional escondida nesta extração.

O [passo0142](../passos/0142-aprofundar-trilhas-e-conservar-memoria.md) acrescenta atéduas leituras de destinos e memória por snapshot dos repositórios nas cinco forges com leitor. A próxima rodada com --memory prioriza pendentes/desconhecidos, evitando releitura de examinados/descartados; feedback offline com motivo permite reabrir caminhos. Histórico declarado e documentos atuais permanecem separados.


O [recorte0143](../passos/0143-expandir-vizinhancas-a-partir-da-memoria.md) acrescenta `idea insights expand --input relatório --output próximo`: continuar a fronteira da memória sem formular outra pergunta ou consultar catálogos. O foco opcional restringe a trilha; documentos congelados alimentam referências e somente leituras desta chamada geram cartões atuais. Cada chamada avança um salto limitado e conserva pendências/falhas; esgotar o foco não comprova esgotar o universo de abordagens.
