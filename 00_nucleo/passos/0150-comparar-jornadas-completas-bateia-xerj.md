# 0150 — Comparar jornadas completas Bateia e XERJ

**Data:** 2026-11-08. **Protocolo ex-ante:** duas jornadas independentes, começando na pergunta e sem referências fornecidas.

## Intenção

Comparar o que um agente usando XERJ e outro usando Bateia encontram e propõem. O ensaio0149 começava com corpus preparado pelo Bateia e não avaliava descoberta ponta a ponta. O usuário solicita agora os dois agentes e a jornada completa.

## Desenho

Mesma pergunta dirigida do caso0148, mesmo alvo XERJ, mesmo modelo/configuração herdados e até20 minutos/24 chamadas de produto por braço. Agentes sem histórico, sem acesso às propostas anteriores ou ao outro braço. Nenhum corpus, catálogo local previamente baixado, resultado anterior, clone ou fonte de referência é fornecido. Binários e instruções genéricas de operação são instrumentos, não corpus de pesquisa.

Permitir descoberta nativa do produto e seu catálogo público padrão. Aquisição por URL exata via Git serve apenas de transporte; o conteúdo deve ser pesquisado/coletado pelo produto do braço. Fontes sugeridas pela memória do agente são permitidas, mas registradas como tal; não contam como descoberta do produto. Proibir web/search/GitHub search externos, leitura direta de fontes e uso cruzado das ferramentas. Permitir leitura de resultados do próprio produto e documentos locais do experimento. Não executar fontes, publicar, enviar feedback ou instalar timers.

Cada braço recebe diretórios novos de índices, caches e resultados, sem credenciais fornecidas. O feedback do XERJ fica desativado; consultas ao Hub público são permitidas. Falhas e intervenções manuais são parte do resultado, não excluídas depois. Comandos, duração, saídas e origem declarada da aquisição ficam registrados em wrapper comum.

## Critérios antes dos resultados

Comparar fontes descobertas, material efetivamente lido, proveniência/localizadores, pertinência ao problema, especificidade das ideias, reconhecimento de capacidades existentes, condições/contraexemplos, experimentos propostos, trabalho manual e bloqueios. Não pontuar mera quantidade de arquivos como qualidade. Distinguir documentos com referência recuperada de execução verificada.

O avaliador principal conhece a conversa, mas os agentes não; a avaliação não é cega. Uma execução por braço e uma pergunta dirigida não estimam desempenho geral, criatividade ou economia de tokens. Registrar métricas operacionais realmente observadas; não inferir tokens da quantidade de texto.

## Artefatos e aceite

Protocolo, runner, transcrições, resultados independentes e comparação ficam em `lab/bancada/insights/jornada-completa-2026-11-08`. Preservar0148/0149. Conferir links, origens e plausibilidade das propostas sem fornecer referências aos agentes durante a execução. Registrar casos inconclusivos. Executar a verificação aplicável do repositório, sem alterar produto ou regras de evidência para favorecer um braço.

## Resultado

Os dois braços terminaram sem referências externas prontas. O XERJ descobriu e recuperou admission-compare pelo Hub; o Bateia descobriu e leu Hypothesis e PredictiveEcology/reproducible. Nenhum braço declarou referência externa de memória. Ambos adquiriram o alvo explícito, na mesma revisão, e produziram três propostas convergentes com condições e experimentos, sem executá-los.

O XERJ usou 24 chamadas CLI, incluindo ajuda/setup, e o Bateia 10, incluindo 15 consultas internas de catálogo. Essas unidades e os tempos não são uma comparação causal de eficiência. Falhas, cobertura, orçamento e uma colisão inicial da auditoria do Bateia estão na [comparação](../../lab/bancada/insights/jornada-completa-2026-11-08/README.md). Os dois encontraram contratos de evidência já existentes no XERJ, qualificando as propostas anteriores como extensões a investigar.

Lint, auditorias, suíte raiz, vet, build e diff-check do Bateia passaram. A ampliação por linguagem pedida posteriormente é uma fase suplementar do passo0151, sem alterar os resultados destes agentes.
