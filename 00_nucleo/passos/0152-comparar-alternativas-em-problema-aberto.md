# 0152 — Comparar alternativas em problema concreto e aberto

**Data:** 2026-11-08. **Protocolo ex-ante:** dois agentes independentes, mesmo problema, sem indicar mecanismos ou fontes.

## Objetivo

Verificar se a jornada Bateia produz alternativas relevantes que não aparecem na jornada XERJ, evitando a pergunta anterior que já sugeria proveniência, aplicabilidade e contraexemplos. Não pressupor vantagem nem exclusividade do Bateia.

## Problema

Uma cooperativa organiza peças doadas em oito oficinas, com 60 voluntários, 45 mil peças e celulares Android simples. Entradas, deslocamentos, separação para consertos, consumo e correções podem ocorrer sem conexão por até72 horas e ser registrados por pessoas diferentes sobre a mesma peça. Planilhas levam a viagens perdidas, promessas de itens indisponíveis e perda de observações nas correções. Fotos e histórico importam. Há dois mantenedores generalistas, orçamento recorrente pequeno e necessidade de exportar os dados. A pergunta pede mudanças no trabalho e no software que reduzam esses problemas, sem indicar arquitetura, algoritmo, linguagem ou produto.

## Desenho e isolamento

Agentes novos, mesmo modelo/configuração herdados, sem histórico, fontes prontas, índices anteriores ou acesso ao outro braço. Cada um registra antes da pesquisa atéquatro opções iniciais em baseline.json; depois usa apenas o produto atribuído para descoberta/leitura. Memória do agente pode sugerir URLs para aquisição, mas a origem deve ser registrada e não conta como descoberta do produto.

Manuais genéricos e binários são instrumentos. Transporte Git por URL exata e preparação/leitura dos resultados são permitidos; pesquisa web/GitHub search externa, leitura direta dos clones, execução das fontes, feedback, publicação e timers ficam proibidos. XERJ usa Hub público e diretórios novos; Bateia pode usar catálogos nativos e coleta. Até20 minutos ou40 chamadas de pesquisa por braço; ajuda, status, sync e lifecycle ficam registrados à parte. Atéseis aquisições. Falhas/setup/tempo de indexação permanecem nos resultados.

## Critérios fixados antes dos resultados

- Relevância de0–4: atende às restrições concretas, reconhece condições difíceis e explicita o trabalho humano restante; não prometer disponibilidade global sem sustentação no contexto sem conexão.
- Evidência de0–3: ausência/mera indicação; descrição superficial; trecho localizado que sustenta o mecanismo; múltiplos trechos ou código/testes específicos. Teste lido não é teste executado.
- Especificidade decisória de0–2: slogan; escolha parcial; mudança concreta de organização/garantias/fluxo com trade-off verificável.
- Alternativas com o mesmo efeito decisório pertencem à mesma família, mesmo se usam linguagens, bibliotecas ou nomes diferentes. Equivalência exige comparar descrição e garantias, não palavras-chave.
- Candidata relevante exclusiva nesta execução: relevância≥3, evidência≥2, especificidade=2 e sem equivalente no outro braço. Isso não prova que o outro produto seja incapaz de encontrá-la.
- Separar alternativas já presentes no baseline das acrescentadas ou materialmente modificadas pela pesquisa. Uma opção inicialmente lembrada pelo agente não vira descoberta da ferramenta por receber citação depois.

## Avaliação

Depois de congelar os resultados, um avaliador independente recebe o problema e pacotes anonimizados, sem nomes dos braços ou resultados anteriores. Agrupa famílias, confere apoio nos trechos e aplica os critérios. O agente principal registra também as medidas e limites. O anonimato é editorial, não garantia absoluta de cegamento; formato de evidência pode sugerir a origem.

Uma execução por braço é exploratória e não permite estimar superioridade geral, novidade, ganho de criatividade ou custo de tokens. Quantidade de referências/arquivos não substitui relevância. Se não houver alternativa exclusiva fundamentada do Bateia, relatar isso diretamente.

## Artefatos e aceite

Protocolo, pergunta, runner, baselines, consultas/saídas, resultados independentes, pacotes de avaliação, parecer e síntese em `lab/bancada/insights/problema-aberto-2026-11-08`. Nenhuma alteração produtiva ou reescrita dos ensaios anteriores. Conferir testemunhos, origem das fontes e verificação aplicável do repositório. Nós temporários devem ser encerrados; nenhum deploy ou comunicação externa faz parte do teste.
