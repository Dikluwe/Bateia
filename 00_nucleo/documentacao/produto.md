# Visão do produto Bateia

O Bateia ajuda quem vai iniciar ou evoluir um software a investigar soluções existentes, identificar oportunidades de reúso e obter insights de engenharia. Isso inclui mecanismos observados, alternativas, trade-offs e hipóteses de melhoria com planos de validação. Reúso possui uma jornada própria; as dependências e os limites atuais de insights estão no [mapa arquitetural](../arquitetura/README.md).

A intenção está registrada em [D-001](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [D-002](../desejos/D-002-alternativas-de-software-com-evidencia.md).

## Investigação independente de insights

O recorte [0105](../passos/0105-investigar-insights-sem-decisao-de-reuso.md) recebe pergunta, alvo e evidências congeladas, sem depender de IdeaSpec ou decisão de reúso. Produz mecanismos vinculados aos bytes fornecidos, hipóteses conservadoras e plano de validação. As avaliações são declaradas e não refeitas; efeito e aplicabilidade no alvo permanecem desconhecidos. Consultar o [exemplo ilustrativo](../../lab/bancada/insights/0105/README.md).

## Entrada e resultado

A jornada principal recebe uma `IdeaSpec` com problema, capacidades e restrições. A política de reúso determina critérios como licenças admitidas e limites de composição. Um alvo opcional permite comparar alternativas de um software conhecido.

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

## Como avaliar utilidade

A qualidade depende de encontrar candidatos relevantes, comprovar capacidades, explicar bloqueios e ajudar a escolher um ponto de partida. Testes de infraestrutura e quantidade de arquivos extraídos são sinais técnicos, mas a avaliação de produto exige ideias reais, julgamentos independentes e cobertura declarada.

A aprovação de uma etapa não promove automaticamente as demais. O [índice de passos](../passos/README.md) registra, por exemplo, descoberta aprovada com validação downstream pendente no passo 0023.
