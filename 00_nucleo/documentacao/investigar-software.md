# Investigar software com o Bateia

O primeiro resultado é um pacote de referências para investigar uma pergunta. A IA lê esse pacote, interpreta as fontes e formula possibilidades; o Bateia organiza a coleta, os localizadores e as tarefas. Para este fluxo, não é necessário criar política de reúso nem avaliações manuais.

## Preparar o executável

Na raiz, com Go 1.25 ou superior, execute `GOCACHE=/tmp/bateia-go-cache go build -o /tmp/bateia ./04_wiring`. Use um binário recém-compilado: um executável antigo deixado na raiz pode não conter os comandos atuais. O cache e o binário ficam fora do código do produto.

## Encontrar referências

```sh
/tmp/bateia idea insights discover \
  --question 'backpressure queue cancellation' \
  --familiar 'thread pool' \
  --output /tmp/descoberta.json
```

Precisa de rede. Sem `--sources`, consulta o catálogo GitHub, seleciona referências e tenta ler READMEs em commits fixados. Token é opcional, sujeito a quotas do serviço. As consultas são lexicais e não traduzem a pergunta. Prefira uma pergunta específica na língua das fontes; `--familiar` registra abordagens que você já conhece.

A descoberta usa `context-v2` por padrão: combina termos da pergunta em consultas menores, busca metadados de repositórios no GitHub (nome, descrição e tópicos) e prioriza correspondência literal antes de variedade. A revisão imutável é resolvida nas referências escolhidas para leitura. O ranking é lexical: confirme a pertinência nos documentos e reformule consultas por capacidades quando necessário. `--query-policy lexical-v1` reproduz a política anterior; não use nomes lembrados como prova de descoberta nativa.

Para outras fontes, use `--sources github,codeberg,gitlab --max-calls 15`. A leitura automática e a disponibilidade variam por fonte. A [referência completa](referencia-comandos.md) descreve credenciais, formatos e limites. Não confundir identidade no catálogo com conteúdo efetivamente lido.

Abra o JSON e confira `attempts`, os bundles e os contextos efetivamente presentes, os cartões/tarefas e as lacunas. Os nomes e a localização variam entre o formato GitHub e o multifonte. Falha, timeout ou resultado vazio não demonstram ausência de soluções.

## Continuar uma trilha

Uma descoberta multifonte conserva memória e referências citadas nos documentos. Para continuar a partir dela, sem nova busca de catálogo:

```sh
/tmp/bateia idea insights expand \
  --input /tmp/descoberta-multifonte.json \
  --focus https://github.com/org/repo \
  --output /tmp/expansao.json
```

Substitua o foco por uma referência existente na memória. Na próxima chamada, use a saída anterior como `--input`; o foco é herdado. Cada chamada realiza um salto limitado. Fronteira esgotada no foco não significa que todas as abordagens possíveis foram exploradas. Veja o [contrato de vizinhanças](../contratos/exploracao-de-vizinhancas-v1.md).

## Explorar material coletado

Se já existe uma lista JSON de bundles emitida por `idea evidence collect`, o passo é offline:

```sh
/tmp/bateia idea insights explore \
  --question 'How can corpus evidence distinguish retrieval relevance from verified behavior?' \
  --bundles /tmp/bundles.json \
  --familiar 'BM25' --familiar 'commit pinning' \
  --output /tmp/exploracao.json
```

O Bateia verifica vínculos com os bytes fornecidos e seleciona passagens. Entregue à IA a pergunta, `cards`, `tasks`, `ai_instructions`, `gaps` e `limitations`. Peça a ela para distinguir observação, interpretação e hipótese e citar os cartões usados. Explore não autentica por si só a origem dos bundles, não chama um modelo e não executa os projetos lidos.

Para experimentar sem coleta nova, substitua `/tmp/bundles.json` por `lab/bancada/insights/0136/exploration-bundles.json`. É um corpus histórico conhecido, com incompletudes documentadas, e não uma investigação atual da rede.

## Atualizar a exploração depois de coletar

Ao selecionar novas fontes úteis, faça a coleta já com a pergunta:

```sh
/tmp/bateia idea evidence collect \
  --plan /tmp/plano.json \
  --question 'What changes could address the operating constraints?' \
  --max-passages 16 \
  --output /tmp/coleta-explorada.json
```

Este modo devolve um objeto com `bundles` e `exploration`, com revisão de alternativas ativa. Sem `--question`, a saída continua sendo a lista histórica de bundles. Identidades sem commit imutável ficam em `unexplored_repositories`; falhas não viram leitura. O plano continua escolhendo fontes e alvos explicitamente; o comando não faz uma busca adicional nem chama IA.

Para bundles já existentes, acrescente `--review-alternatives` ao comando `explore`. O modo conserva os cartões iniciais e preenche vagas com até duas passagens adicionais não sobrepostas por arquivo, sob o mesmo teto global. As tarefas pedem comparar decisões dentro da mesma fonte e registrar opções consideradas, descartes justificados e lacunas. Leia os documentos completos quando as passagens não bastarem. A revisão não certifica que todas as alternativas foram encontradas; direções de pesquisa permanecem planos até serem executadas.

## Investigar princípios de algoritmos e decisões

Use `--algorithm-ideas` em `discover`, `explore` ou `evidence collect --question` para preparar a investigação de princípios, passos/invariantes, condições e decisões. Exemplo:

```sh
/tmp/bateia idea insights discover \
  --question 'code search index performance correctness' \
  --algorithm-ideas \
  --output /tmp/fontes-algoritmos.json
```

O modo acrescenta tarefas e ativa a revisão com passagens adicionais. Leia a fonte e o alvo: uma ideia pode já existir, aplicar-se apenas a outro caminho ou não ter evidência suficiente de transferência. A seleção de trechos continua lexical e o agente faz a interpretação.

Depois de formular as hipóteses, registre-as com:

```sh
/tmp/bateia idea insights formulate \
  --input /tmp/ideias-com-evidencias.json \
  --output /tmp/ideias-conferidas.json
```

O pedido contém `target` (bundle do alvo), `sources` (bundles das fontes) e `ideas`. Cada ideia exige princípio, passos, evidências da fonte e do alvo, condições/falhas, decisão atual/proposta, benefício/custo e observável/critério de rejeição. Citações precisam conferir com identidade, commit, path, hash, linhas e literal dos bundles. `target_relation` registra ALREADY_PRESENT, DIFFERENT_APPROACH ou NOT_ESTABLISHED como declaração do agente. O relatório mantém interpretação e relação com o alvo como DECLARED_NOT_REASSESSED e efeito UNVERIFIED. Aprovar o registro confere os vínculos; não certifica o algoritmo, sua novidade ou benefício.

## Transformar o contexto em propostas

Para cada ideia, registre problema concreto, fonte e trecho, mecanismo observado, diferença em relação ao alvo, condições de aplicação e experimento que poderia rejeitar a proposta. Confira se o alvo já oferece a capacidade antes de sugeri-la como nova. O [caso XERJ](../../lab/bancada/insights/xerj-2026-11-08/README.md) aplica esse procedimento.

`analyze` e `compare` atendem outra necessidade: organizar avaliações e observações fornecidas num contrato específico, com catálogo limitado. Não são um comando geral que interpreta qualquer repositório automaticamente. `idea recommend` acrescenta políticas e uma decisão de reúso. Consulte esses fluxos na [referência de comandos](referencia-comandos.md).

## Guardar e retomar

Resultados temporários ficam em `/tmp` ou numa pasta local da investigação. Uma evidência que sustenta uma conclusão preservada pertence a `lab/bancada`, com origem, versão, comandos e limites. Contratos e sínteses declarativas ficam em L0; não coloque scripts ou saídas brutas nessa camada. Investigar um projeto público não publica o Bateia nem envia propostas aos seus responsáveis.
