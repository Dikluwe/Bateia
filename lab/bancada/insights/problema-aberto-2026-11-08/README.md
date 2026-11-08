# Bateia × XERJ: alternativas para um problema aberto

**Resultado:** a jornada Bateia produziu uma alternativa relevante ausente no resultado XERJ: configurar uma plataforma de formulários ODK com estado por peça e revisão operacional. A jornada XERJ também produziu variantes relevantes ausentes no resultado Bateia. Houve diferença de decisões, mas **não foi demonstrada vantagem de descoberta nativa do Bateia**: os projetos úteis de ambos vieram da memória dos agentes. As fontes refinaram famílias já previstas nos baselines.

Esta é uma investigação exploratória de um cenário hipotético concreto, realizada em 2026-11-08; nenhum sistema foi implantado, nenhum projeto-fonte ou piloto foi executado. A avaliação qualifica propostas condicionais, não soluções comprovadas em campo.

## Problema e controle do ensaio

Uma cooperativa tem oito oficinas, 45 mil peças individuais doadas e 60 voluntários com Android simples. Há até 72 horas sem rede. Pessoas registram ações e correções concorrentes; ocorrem viagens perdidas, promessas de peças indisponíveis e perda de observações. Fotos, histórico e exportação importam. Dois generalistas mantêm o sistema, com orçamento recorrente pequeno. A pergunta pede mudanças no trabalho e no software, alternativas e um ensaio pequeno antes de escolher.

O [enunciado integral](problem.txt) foi idêntico nos dois braços, sem mecanismos, fontes ou corpus fornecidos. Cada agente novo registrou até quatro opções **antes da primeira ferramenta**. O runner congelou e conferiu seu hash em cada ação. Não receberam histórico da conversa nem resultados do outro braço. Cada um só podia descobrir e ler fontes usando o produto atribuído; Git por URL pública exata servia apenas de transporte. URLs lembradas eram permitidas, mas registradas como origem `model`, separadas de descoberta `tool`.

O [protocolo](protocol.json) limitou cada braço a 20 minutos ou 40 chamadas de pesquisa e seis aquisições; ajuda, sincronização do catálogo e lifecycle foram registrados separadamente. O [runner](run.py), as tarefas e os guias genéricos estão preservados. Os guias diferem conforme a interface: Bateia recebeu exemplos de formatos porque a ajuda é parcial; XERJ recebeu indicação de suas ajudas nativas. Não se mediu custo de tokens nem eficiência dos modelos.

Um terceiro agente novo recebeu somente problema, [rubrica prévia](review-rubric.md) e pacotes com os nomes dos produtos removidos. O anonimato é editorial; formatos podem sugerir origem. A avaliação e seus critérios foram congelados antes da síntese principal. Pacote X corresponde a XERJ e Y a Bateia; o mapa não foi entregue ao avaliador.

## Comparação por decisão

| Família operacional | Bateia | XERJ | O que muda |
|---|---|---|---|
| Diário imutável e projeção de estado | A1, PouchDB | A2, PouchDB | Operações/correções acrescentam registros; estado é reconstruível; custódia confirma fisicamente a reserva. Núcleo equivalente. |
| Formulários e estado por entidade | A2, ODK | Ausente no resultado | Configurar captura existente; revisar atualizações recebidas e corrigir o estado da entidade separadamente da submissão histórica. |
| Ficha corrente e reconciliação documental | Ausente no resultado | A1, PouchDB | Editar a ficha atual, conciliar conflitos e manter arquivo histórico separado; não derivar todo estado de operações. |
| Colaboração descritiva com mescla | Ausente no resultado | A3, Automerge | Mesclar estruturas descritivas e expor valores concorrentes por campo; decisões físicas continuam separadas. |

O [parecer anônimo completo](assessment.md) qualificou as cinco propostas usando relevância ≥3/4, evidência ≥2/3 e especificidade decisória =2/2. Não somou notas nem declarou vencedor. A1 Bateia e A2 XERJ receberam 4/3/2 e foram agrupadas como equivalentes. ODK e Automerge receberam 3/3/2: são relevantes sob condições, com lacunas importantes. A ficha XERJ recebeu 4/3/2.

**ODK é uma diferença operacional relevante, não apenas outra biblioteca.** A responsabilidade muda de manter aplicação/projeção próprias para configurar formulários e administrar/revisar estado de uma plataforma existente. A leitura mostrou que editar uma submissão não atualiza automaticamente a entidade relacionada, o que muda o fluxo de correção. Mas fotos históricas acessíveis offline por peça, carga de 45 mil registros, filtros, telas adicionais e esforço de manutenção permanecem por validar. O próprio parecer considerou alguns excertos curtos insuficientes para sustentar limites adicionais, como versão mínima e retenção de cinco dias; essas alegações não foram promovidas a comportamento verificado.

**As variantes XERJ também mudam decisões.** Ficha mutável e diário reconstruível não são a mesma obrigação de manutenção/reconciliação. A alternativa Automerge muda a unidade da colaboração descritiva e da inspeção de conflitos, ao custo de transporte, persistência e interface próprios. Sua utilidade depende de haver colaboração descritiva suficiente para compensar esse trabalho aos dois mantenedores.

## O que a pesquisa acrescentou

Os [baselines Bateia](bateia/baseline.json) e [XERJ](xerj/baseline.json) já continham observações imutáveis, reconciliação e custódia. O diário comum não é descoberta da pesquisa. O baseline Bateia também já propunha uma plataforma existente de coleta offline. O parecer julgou que o seguimento por entidades e a independência entre correção de submissão e estado representam **revisão material desse caminho**, com confiança moderada; não uma família inteiramente nova.

No XERJ, o parecer considerou a mescla descritiva uma revisão material, também com confiança moderada. Para a ficha e a fila de conflitos, julgou que o incremento material alegado pelo agente não ficou demonstrado frente ao baseline: havia reconciliação explícita desde o início. Essa divergência foi mantida no parecer, sem reescrever o resultado do pesquisador.

Portanto, a resposta depende do alcance da afirmação:

- **A jornada Bateia chegou a uma alternativa relevante que o outro resultado não apresentou? Sim, ODK.** A exclusividade é observada nesta execução e há um refinamento decisório apoiado por fonte.
- **Bateia descobriu essa alternativa que o agente não teria lembrado? Não foi demonstrado.** A família já constava no baseline e a URL ODK foi lembrada pelo agente.
- **Isso demonstra vantagem sobre XERJ? Não.** Também houve variantes relevantes só no XERJ e não houve repetição suficiente para estimar frequência ou causalidade.

## Descoberta, leitura e auditoria

| Medida registrada | Bateia | XERJ |
|---|---:|---:|
| Chamadas de pesquisa | 9 | 24 |
| Chamadas de setup | 1 | 8 |
| Aquisições Git | 2 | 2 |
| Aquisições úteis com origem em descoberta nativa | 0 | 0 |
| Trechos distintos conferidos literalmente no pin | 17/17 | 14/14 |
| Localizadores de linhas imprecisos | 0 | 3 |

Essas contagens descrevem caminhos distintos, não eficiência comparável. O XERJ inclui a verificação de encerramento entre suas pesquisas. Os snapshots de baseline não contam como chamadas. Ver [métricas](metrics.json) e [conferência de evidência](evidence-check.json).

Bateia tentou quatro descobertas, incluindo multifonte; houve baixa aderência e uma saída parcial sem bundles. O agente não tratou os resultados irrelevantes como soluções. A coleta remota com HEAD falhou; depois PouchDB e ODK docs, lembrados pelo agente, foram adquiridos e lidos pelo produto. XERJ sincronizou catálogo de 514 corpora, recebeu ambiguidade/no_coverage e não aproveitou hits irrelevantes; adquiriu PouchDB e Automerge de memória e pesquisou os índices locais pelo produto. A recuperação fundamentou condições e limites nos dois braços.

Os resultados originais estão em [relatório Bateia](bateia/report.md), [alternativas Bateia](bateia/alternatives.json), [fontes Bateia](bateia/sources.json), [relatório XERJ](xerj/report.md), [alternativas XERJ](xerj/alternatives.json) e [fontes XERJ](xerj/sources.json). Consultas, saídas, origens e audits ficam nas respectivas pastas. A conferência principal ocorreu depois de congelados os resultados: os pesquisadores não receberam fontes ou correções do avaliador. Todos os literais foram encontrados nos arquivos fixados; três referências XERJ apontam para o início da passagem, fora da linha exata do literal. Os originais foram preservados e os flags entregues ao avaliador.

O nó XERJ exclusivo foi encerrado, conforme [lifecycle](xerj/lifecycle.json) e audit/034. O audit/013 de startup permanece `pending` porque o runner terminou por interrupção; isso é uma lacuna do instrumento, não prova de processo ativo. Nenhum feedback, timer, publicação ou comunicação a terceiros foi realizado.

## Limites e verificação

É um cenário, uma execução por braço, mesmo modelo/configuração e baselines breves. O resultado mistura raciocínio do agente, fontes lembradas, descoberta e recuperação; não isola o efeito causal da ferramenta. Não encontrar uma alternativa não demonstra incapacidade de encontrá-la. Ler documentação/código não comprova comportamento em Android, custo, adesão dos voluntários ou sucesso em campo. O ensaio de 72 horas, celulares baratos, falhas concorrentes e exportação proposto pelos agentes continua futuro.

Não houve alteração funcional para favorecer o ensaio. A verificação canônica `make verify` passou em lint com falha em warning, contratos arquiteturais, `go test ./...`, `go vet ./...`, build de `./04_wiring` e diff. [Resultados e comandos](verification/verification.json), com logs na mesma pasta. Esses gates verificam o repositório; não promovem a força das propostas.

O [protocolo do passo 0152](../../../../00_nucleo/passos/0152-comparar-alternativas-em-problema-aberto.md), os pacotes X/Y, o parecer JSON e [síntese estruturada](comparison.json) permitem auditar a conclusão. Um manifesto SHA-256 acompanha os artefatos desta bancada; binários, clones e índices temporários não são incorporados ao repositório.
