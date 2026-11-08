# Jornada completa: pergunta → fontes → ideias

**Neste ensaio, os dois agentes chegaram a três propostas fundamentadas começando sem referências prontas. O XERJ conseguiu descobrir material pelo próprio Hub, sem participação do Bateia.** As propostas convergiram em recibos por alegação, condições de aplicação e contraexemplos/replay; as fontes e o trabalho de cada percurso foram diferentes.

As [buscas adicionais em sete linguagens](multilinguagem/README.md), solicitadas depois, não integram a comparação original. O [protocolo](protocol.json) foi registrado antes das execuções. Os resultados independentes foram preservados: [braço XERJ](xerj/report.md) e [braço Bateia](bateia/report.md).

## O que cada um encontrou

| Aspecto | Agente + Bateia | Agente + XERJ |
| --- | --- | --- |
| Descoberta externa | Catálogo GitHub: Hypothesis e PredictiveEcology/reproducible, ambos lidos | Hub nativo: corpus admission-compare, recuperado localmente |
| Origem das referências | Produto; nenhuma referência externa classificada como memória do agente | Produto; nenhuma referência externa classificada como memória do agente |
| Próprio alvo | Adquirido pela URL comum, seguido de coletas por plano | Adquirido pela mesma URL comum, seguido de autoindex e buscas |
| Versão do alvo | `1557de2f5be981d9db441969201c92a8af6d2de3` | Mesma revisão, conferida pelo avaliador após concluir |
| Evidência mais específica do alvo | Contrato completo de `hub_evidence.rs`, abstenção/calibração e limites do Facts API | PinnedCitation/from_bytes, recibos CJK/startup e limites da verificação lexical |
| Transferência externa | Geração/redução de exemplos do Hypothesis e replay/cache do reproducible | Admissão separada de autoridade, evidência expirada e fonte fraca no admission-compare |
| Proveniência externa | URLs públicas, commits, hashes e trechos nos bundles | Corpus, pin, licença e arquivo/linha; o resultado não forneceu URL upstream |
| Chamadas CLI do produto | 10, incluindo 15 consultas internas de descoberta | 24, incluindo 5 ajudas, setup e nó local |
| Aquisição manual por URL | 2: alvo e Hypothesis após falha de leitura da revisão no catálogo | 1: alvo; o corpus externo veio pelo produto |
| Resultado | Três hipóteses e experimentos propostos | Três hipóteses e experimentos propostos |

O endereço do alvo era dado comum, não uma referência descoberta. Nenhum dos dois recebeu c-ares, Feynman, Hypothesis, admission-compare ou qualquer lista de projetos. Ambos prepararam consultas e interpretaram o material recuperado; nenhum produto produziu as ideias sem trabalho do agente.

## Ideias que surgiram

| Tema | Desenvolvimento pelo Bateia | Desenvolvimento pelo XERJ |
| --- | --- | --- |
| Recibos de verificação | Ampliar o TraceEvidenceReport existente com avaliações independentes de integridade, relevância e resultado de experimento | Compor recuperação, integridade, execução e comportamento na alegação, preservando sinais separados |
| Aplicabilidade | Ficha de condições e pares mínimos, usando a fronteira people/devices já documentada no Query-1 | Predicados supported/contradicted/unknown para versão, fixture, modo lexical/neural, ambiente e validade |
| Contraexemplos/replay | Gerar e reduzir casos preservando a propriedade que falha, com oracle e identidade no cache | Cápsulas pinadas com entrada, política, tick/ambiente, resultado esperado/observado e histórico |

Os dois reconheceram mecanismos já existentes e não demonstraram a ausência global das extensões propostas. A pergunta já direcionava os temas de proveniência, aplicabilidade e contraexemplos; a convergência não mede criatividade espontânea. Nenhum experimento de ganho ou defeito novo foi executado pelos agentes.

## Correção importante em relação aos recortes anteriores

As duas jornadas encontraram `engine/crates/xerj-common/src/hub_evidence.rs`, que já distingue integridade de entailment e contém citações pinadas. O braço Bateia aprofundou seu contrato de calibração, cobertura e abstenção. Esse arquivo não integrava o corpus inicial de 17 arquivos do caso0148. Portanto, não cabe apresentar “adicionar proveniência/recibos” como novidade genérica para o XERJ. As propostas precisam mostrar uma diferença concreta de exposição, composição ou validação sobre o que já existe. A integração efetiva em todos os consumidores ainda não foi verificada.

## Dificuldades e trabalho restante

**Bateia.** A primeira descoberta multifonte gastou orçamento em catálogos de apps e outras fontes sem trazer material pertinente. O agente precisou reformular consultas e restringir GitHub. Houve README grande, revisões opacas, timeout e leitura recusada. A concatenação de bundles preparada pelo agente foi rejeitada na exploração; uma entrada menor funcionou. Os cartões não priorizaram o contrato central, e o agente leu os bytes retornados pelo coletor. A tentativa de Anserini não obteve conteúdo e não fundamentou propostas.

**XERJ.** O primeiro find não tinha catálogo; foi preciso sync. O caminho hosted foi recusado nesse build e o fallback local funcionou. O primeiro resultado de catálogo não foi considerado automaticamente relevante. A indexação do alvo levou112,235s e terminou com conteúdo pesquisável e exclusões; exit3 não foi classificado como falha geral. Houve buscas durante a indexação, queries sem resultado e ruído. O agente atingiu o limite de24 chamadas, incluindo setup; esse limite restringe o que se pode concluir sobre descoberta adicional.

As durações observadas entre primeira chamada e relatório foram aproximadamente405s no Bateia e455s no XERJ. Não são estimativas comparáveis de latência ou custo: incluem raciocínio, comandos sobrepostos e transporte variável (o mesmo clone levou23,3s em um braço e69,8s no outro). Não medimos tokens. Quantidade de arquivos, registros ou chamadas não é pontuação de qualidade.

## Isolamento e instrumentação

Os agentes usaram o mesmo modelo/configuração herdados, contexto novo sem histórico e diretórios próprios sem corpus/cache inicial. Apenas os manuais genéricos de operação foram fornecidos. Bateia recebeu uma descrição operacional com exemplo de plano porque sua ajuda é parcial; XERJ recebeu instrução para usar a ajuda nativa. Essa assistência desigual é um limite de usabilidade do ensaio, não uma intervenção de descoberta.

Pesquisa e aquisição passaram pelo [runner](run.py), que limpou credenciais e desativou feedback. Leitura direta de fontes e buscadores externos ficaram proibidos por protocolo; não houve isolamento criptográfico dos agentes. O avaliador conhece a conversa, portanto a avaliação não é cega.

Duas operações iniciais concorrentes do braço Bateia colidiram no audit/001. O registro de aquisição sobrepôs o de descoberta. A saída do produto e os fatos observados foram preservados na [auditoria suplementar](bateia/supplemental-audit.json); a reserva antecipada corrigiu o runner durante o ensaio. A contagem de10 chamadas inclui a operação recuperada; não apresentamos a transcrição persistente como completa nesse ponto.

Os comandos, saídas e durações estão em `bateia/audit` e `xerj/audit`; [metrics.json](metrics.json) conserva as medidas. [setup.json](setup.json) registra os hashes dos binários. Os projetos não foram executados; somente os produtos de investigação. O nó temporário XERJ foi encerrado. Não houve publicação, comunicação ou feedback.

## Repetir a comparação

Use diretórios novos, binários com os hashes de setup.json e dois agentes sem histórico. Entregue a cada agente apenas seu `task.txt`, `tool-guide.txt` e o runner, com o mesmo alvo e pergunta. Não forneça os resultados, sources.json, ideas.json, relatórios ou as pastas adquiridas. O catálogo e os repositórios podem mudar; uma nova execução não promete as mesmas fontes ou ideias. Preservar a própria tentativa é parte do procedimento.

O gate do Bateia passou em lint, auditorias, testes, vet, build e diff-check; os registros estão em `verification/`. A conferência de hashes dos bundles verifica seus vínculos internos, não a validade semântica das propostas. Uma pergunta e uma execução por braço não estabelecem superioridade geral de nenhuma ferramenta.
