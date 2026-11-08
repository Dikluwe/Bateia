# Reteste 0153 — Bateia com coleta ligada à revisão de alternativas

**A correção foi implementada, verificada e usada pelo agente. A omissão de réplicas colaborativas/Automerge não se repetiu. Os dois braços chegaram às mesmas duas políticas operacionais centrais; não houve alternativa exclusiva nova atribuível à pesquisa.** Reaproveitar ODK permaneceu uma diferença relevante de responsabilidade de implementação no braço Bateia, já prevista em seu baseline.

Este reteste usa a pergunta de negócio de [0152](../problema-aberto-2026-11-08/README.md): oito oficinas, 45 mil peças, 60 voluntários, Android simples, até 72 horas sem rede, dois mantenedores, fotos, histórico e exportação. O cenário é hipotético concreto. Não executamos projetos-fonte nem um piloto operacional.

## O que mudou no produto

A correção trata a ligação entre novos documentos e exploração, além da comparação dentro da mesma fonte. Regras e coordenação estão em L1; flags e publicação em L2. Não foram inseridos nomes de bibliotecas, soluções ou mecanismos esperados do ensaio.

| Antes | Modo novo opt-in |
|---|---|
| `collect` devolvia bundles; explorar os novos bytes exigia outro comando | `collect --question` devolve bundles e exploração atualizada na mesma saída |
| Um trecho selecionado por arquivo | Revisão conserva os trechos iniciais e usa vagas para até duas janelas adicionais não sobrepostas por arquivo |
| Tarefas gerais de alternativas/inversão/recombinação | Acrescenta comparação de alternativas na mesma fonte e revisão explícita de cobertura, descartes e lacunas |
| Fonte não resolvida não sustentava leitura | Continua sem sustentar leitura; fica explícita também em `unexplored_repositories` |

`explore --review-alternatives` ativa o modo sobre bundles existentes. O teto global de cartões permanece; mais janelas não significam mais alternativas comprovadas. O modo padrão e o array legado de coleta permanecem. O produto prepara o contexto e as tarefas; o agente interpreta. Não há classificação semântica de cobertura nem chamada de modelo embutida.

O [registro da correção](change.md) explica ownership, testes e limites. O [guia atualizado](../../../../00_nucleo/documentacao/investigar-software.md) documenta a jornada. Fonte da coordenação: [insight_collection.go](../../../../01_core/usecase/insight_collection.go).

## Verificação da mudança

O [replay separado](regression-replay/comparison.json), nunca entregue aos pesquisadores, usou os documentos históricos de 0152. O modo padrão novo teve saída idêntica byte a byte à do binário anterior. A revisão passou de 12 para 16 cartões, acrescentando quatro janelas e duas tarefas (cinco para sete). É uma verificação do retorno do produto, sem julgar criatividade.

Os testes novos cobrem bytes coletados/contexto, preservação da primeira passagem, janela distante e literal, ausência de sobreposição, limite global, imutabilidade do pedido, cancelamento, erro antes de coleta, falha não promovida a leitura, formato opt-in versus legado e proteção do plano. A [verificação canônica](verification/verification.json) passou em lint, contratos arquiteturais, testes completos, vet, build e diff. O binário desse build teve o mesmo SHA-256 do congelado para o ensaio.

## Como o ensaio foi refeito

Dois agentes novos receberam a mesma [pergunta](problem.txt), sem referências, mecanismos ou histórico. Cada baseline foi congelado antes de qualquer ação de produto, e todos os audits preservam seu hash. Usaram apenas seu produto para descoberta e leitura, com aquisição Git por URL exata auditada e origem `model` ou `tool`. Mesmos tetos: 20 minutos, 40 chamadas de pesquisa e seis aquisições; setup separado. Guias específicos das interfaces foram preservados nas pastas dos braços.

Ambos receberam a mesma exigência de registrar opções consideradas, apresentadas, descartadas e lacunas. O guia Bateia apresentou a nova interface genericamente, sem citar fontes esperadas. O [protocolo](protocol.json) registra que novos agentes e essa exigência comum impedem isolar causalmente o efeito do código numa comparação antes/depois.

Um terceiro agente novo avaliou os pacotes X/Y com nomes de produtos removidos, problema e [rubrica](review-rubric.md), sem consultar o repo, fontes externas ou o ensaio anterior. O anonimato é editorial e imperfeito. A correspondência, não fornecida ao avaliador, é X=Bateia, Y=XERJ. O [parecer completo](assessment.md) conserva notas, divergências com as autoavaliações, agrupamentos e incertezas.

## Resultado por decisão

| Decisão | Bateia | XERJ | Avaliação |
|---|---|---|---|
| Captura livre offline, pedidos provisórios, preservação de observações e conciliação antes de promessa | A1 e A3 | A | Mesmo núcleo de garantia; implementações e trabalho de manutenção variam |
| Custódia local restringe quem confirma compromissos, mantendo observações livres | A2 | B | Mesma política central; transferência e recuperação ainda são protocolos propostos |
| Réplicas colaborativas e conflitos visíveis, com Automerge como opção | A3 | Opção interna de A | Considerada pelos dois; três títulos contra dois não significam maior cobertura |
| Reaproveitar captura Android/formulários/auditoria ODK | A1 | Não investigado; declarado como lacuna | Diferença parcial relevante: configurar captura existente versus construir o cliente |

O avaliador qualificou as cinco propostas (relevância 4/4 e especificidade 2/2). Para as duas propostas de custódia atribuiu evidência 2/3: as fontes apoiam primitivas e limites, não comprovam o protocolo de autoridade. Para as outras, 3/3 indica apoio estático de mecanismo/limite, sem comprovar o sistema completo. As notas não foram somadas nem usadas para declarar vencedor.

### O que melhorou frente à omissão anterior

O agente Bateia desta vez investigou CouchDB e Automerge, explicitou conflitos, vencedores determinísticos e retenção independente de revisões compactáveis. A colaboração foi considerada, com condição e piloto, em vez de ficar fora das fontes. Também entregou a revisão de todas as opções do baseline e razões para descartar atalhos como tratar last-write-wins como garantia de estoque.

A correção foi efetivamente exercitada: sete chamadas de coleta com pergunta; seis saídas com cartões e uma apenas com lacunas de resolução. As coleções de documentação ODK tiveram oito passagens adicionais por saída; uma coleção Automerge teve uma. As tarefas de revisão estavam presentes em todas. Ver [uso da correção](correction-usage.json). Não somar cartões de saídas sobrepostas como descobertas independentes.

Isso demonstra que o novo contexto entrou na jornada e que a omissão de colaboração não se repetiu. Não demonstra que o código causou a escolha de Automerge: o próprio baseline Bateia já previa réplicas colaborativas, e a URL foi lembrada pelo agente. Tampouco prova que toda subalternativa possível passou a ser comparada: o desenho de ficha mutável versus projeção de log não foi isolado como uma comparação experimental própria neste reteste. A revisão orienta o agente, sem certificar exaustividade semântica.

### O que não foi demonstrado

Nenhuma família exclusiva nova foi acrescentada pela pesquisa. ODK muda a responsabilidade de manutenção e constitui uma exclusiva parcial observada do Bateia, mas o baseline já previa adaptar aplicativo de campo offline. O parecer considerou essa relação **mesma família**, apesar de o agente a chamar de materialmente modificada; preservamos as duas avaliações.

Ambos revisaram detalhes relevantes a partir das fontes, sobretudo não depender de revisões compactáveis para conservar histórico. Essas revisões não ficaram sem equivalente no outro braço. A descoberta nativa continuou fraca nos dois; todos os projetos úteis foram adquiridos a partir de URLs lembradas pelos agentes. A correção não resolveu relevância de catálogo nem demonstrou repertório superior ao de um agente com XERJ.

## Evidência e auditoria

| Medida | Bateia | XERJ |
|---|---:|---:|
| Chamadas de pesquisa | 9 | 20 |
| Setup | 1 | 7 |
| Aquisições por URL de memória | 4 | 2 |
| Fontes úteis descobertas nativamente e adquiridas | 0 | 0 |
| Trechos distintos citados nas alternativas conferidos em Git e retorno do produto | 18/18 | 10/10 |
| Divergências nas linhas declaradas desses trechos | 0 | 0 |

A quantidade de consultas/aquisições e a escala dos contextos diferem; a tabela não mede eficiência ou custo de tokens. O número de títulos não é a unidade de comparação. A [conferência de evidência](evidence-check.json) inclui hashes dos blobs e presença nos retornos efetivamente registrados; as [métricas](metrics.json) preservam origens, pins e pendências de lifecycle. O relato Bateia menciona 19 testemunhos no inventário de fontes; a tabela conta os 18 trechos distintos efetivamente citados nas alternativas.

Relatórios originais: [Bateia](bateia/report.md), [alternativas](bateia/alternatives.json), [cobertura](bateia/coverage.json), [fontes](bateia/sources.json); [XERJ](xerj/report.md), [alternativas](xerj/alternatives.json), [cobertura](xerj/coverage.json), [fontes](xerj/sources.json). Planos, saídas e audits ficam nessas pastas, sem incorporar clones, índices ou binários. O [setup](setup.json) registra os binários; fingerprints dos fontes alterados estão em `change-fingerprints.json`.

O XERJ expôs inesperadamente nomes de journals de experimentos anteriores num comando de status no encerramento (audit/028), apesar da URL/prefixo e diretórios novos. O agente declarou não ter aberto nem usado os arquivos, e não há aquisição/leitura posterior de fonte na trilha. Isso é uma imperfeição do isolamento que mantemos registrada, sem tratar o ensaio como isolamento absoluto. O nó exclusivo foi encerrado; audit/029 confirmou endpoint inacessível. O audit/012 de startup ficou `pending` após SIGINT/KeyboardInterrupt do runner; não foi reescrito. Ver [lifecycle](xerj/lifecycle.json).

## Limites práticos

Literalidade não garante que o trecho sustente toda a afirmação: o parecer apontou recortes Bateia curtos para algumas alegações de exportação, LWW, mídia e acesso histórico. Essas extensões não foram usadas como garantias. Nenhum protocolo de custódia, retenção Android, custo, carga de 45 mil peças, exportação/restauração completa ou disciplina dos voluntários foi validado em execução.

Os pilotos continuam propostas futuras. A evidência permite avaliar escolhas condicionais; não demonstra redução de viagens, adequação econômica ou superioridade geral. Uma execução por braço e baselines breves não isolam criatividade nem efeito causal da ferramenta. Ausência num resultado não prova incapacidade de encontrar uma alternativa.

O [passo 0153](../../../../00_nucleo/passos/0153-ligar-coleta-a-revisao-de-alternativas.md) registra o contrato da mudança; `comparison.json` conserva a síntese estruturada. O manifesto SHA-256 acompanha os artefatos finais.
