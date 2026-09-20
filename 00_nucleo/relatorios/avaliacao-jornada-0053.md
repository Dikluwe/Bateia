# Execução 0053 — Perguntas novas e perdas por estágio

2026-09-08. O [passo 0053](../passos/0053-avaliar-jornada-em-perguntas-novas.md) entregou um corpus congelado, executor por estágios, controles de métricas e execução local/pública. **A avaliação não confirma qualidade geral: recuperou somente uma das seis funções esperadas, embora tenha extraído cinco dos seis arquivos-alvo.** Nenhum algoritmo de produção foi alterado para melhorar esse resultado.

## Corpus e protocolo

Autor B escolheu fsnotify v1.9.0, go-colorful v1.2.0 e bitset v1.24.2: monitoramento de arquivos, ciência de cores e conjuntos de bits. Escreveu duas perguntas funcionais em português e uma negativa por projeto, sem consultar ranking, catálogos, testes ou benchmarks antigos. Os alvos são declarações exatas, sustentadas por arquivos/intervalos, não somente caminhos. Versão, licença, proveniência e hashes de 181 arquivos estão no [corpus](../../lab/bancada/history/benchmarks/0053/corpus.json), congelado em [freeze.json](../../lab/bancada/history/benchmarks/0053/freeze.json) antes da execução.

O autor A escreveu o executor sem ler esse corpus; o integrador testou métricas com fontes sintéticas antes da avaliação real. Regime Tekt A/B sem atestação física: filesystem compartilhado, sem selo de refinamento. Projetos vieram do cache de módulos disponível, logo são uma amostra de conveniência, não representativa do ecossistema. As perguntas deixam de ser holdout para futuras mudanças após esta publicação.

Cada caso valida inventário e hashes antes/depois, chama a recuperação anterior à seleção e executa `analyze` pela CLI real. Arquivo presente, símbolo recuperado, arquivo extraído e símbolo confirmado são observáveis separados. `body_available` refere-se ao corpo de um símbolo recuperado, com intervalo e bytes conferidos; não conta implicitamente todas as funções contidas em um arquivo extraído. Nenhum projeto externo foi executado; comportamento permanece `not_evaluated`.

## Resultado local

Seis positivas e três negativas avaliáveis; nenhum caso local desconhecido.

| Observável nas positivas | Resultado |
| --- | --- |
| Arquivo-alvo examinado na recuperação | 6/6 |
| Arquivo-alvo nos primeiros 3 caminhos distintos dos hits | 4/6 |
| Arquivo-alvo nos primeiros 12 caminhos distintos dos hits | 5/6 |
| Arquivo-alvo efetivamente extraído | 5/6 |
| Símbolo esperado recuperado antes da seleção | 1/6 |
| Símbolo esperado confirmado após extração | 1/6 |
| Corpo do símbolo recuperado disponível e conferido | 1/6 |

Três negativas sem hits ou alegações explícitas de satisfação detectadas pelo executor. Isso não é auditoria exaustiva de toda linguagem do relatório. Não houve prova comportamental em nenhuma etapa.

| Pergunta/objetivo | Alvo congelado | Arquivo extraído | Função recuperada |
| --- | --- | --- | --- |
| Monitorar alterações em diretório sem subdiretórios | fsnotify `Add` | Sim | Não |
| Parar notificações de caminho monitorado | fsnotify `Remove` | Não | Não |
| Interpolar cores em L*a*b* | colorful `BlendLab` | Sim | Não |
| Converter RGB linear para CIE XYZ | colorful `LinearRgbToXyz` | Sim | Sim |
| Posições presentes em ambos os conjuntos | bitset `Intersection` | Sim | Não |
| Posições presentes em qualquer conjunto | bitset `Union` | Sim | Não |

[Métricas finais locais](../../lab/bancada/history/benchmarks/0053/local-final/summary.json) e [casos com canais, fontes e lacunas](../../lab/bancada/history/benchmarks/0053/local-final/cases.json). Os recibos completos da coleta foram preservados em `local-collected/`.

## Onde a recuperação perdeu qualidade

Os hits ajudam a explicar as perdas, sem tornar a explicação uma prova causal completa:

- Monitoramento: a primeira pergunta privilegia correspondências por `file`; a segunda retorna apenas `watch` por `receive`. A função `Remove` não entrou nos hits.
- Interpolação: funções com termos `l` e `b` ocupam posições fortes, mas a operação `BlendLab` não aparece.
- Conjuntos de bits: contagens `popcnt…` aparecem pelo termo `bits`; as operações pedidas de interseção e união não são recuperadas.

Cinco positivas atingiram o limite de 24 hits; a pergunta de parar notificações retornou apenas um hit. Limites podem contribuir para perdas, mas aumentá-los não resolve uma consulta sem correspondência útil. Nesta amostra, nenhuma função-alvo recuperada se perdeu na extração: a perda dominante já existia antes dessa etapa. Medir somente arquivo-alvo teria dado 5/6 e escondido a recuperação de funções em 1/6.

## Descoberta pública

Executadas as três perguntas marcadas para descoberta, sem fornecer nomes de repositórios. As consultas completas retornaram zero candidatos. Os fallbacks reduziram o pedido aos dois primeiros termos e, depois, removeram o filtro de linguagem:

- fsnotify: `Como começar`; a última tentativa terminou em timeout. Estado final da medição: **unknown**.
- colorful: `Como criar`; dez candidatos, sem `lucasb-eyer/go-colorful` entre eles.
- bitset: `Como obter`; dez candidatos, sem `bits-and-blooms/bitset` entre eles.

Resultado: **0/2 projetos esperados nos primeiros dez candidatos avaliáveis; uma busca desconhecida**. A lista de colorful inclui `Benevanio/Jobs_Scraper_Global`; a de bitset inclui `OpenAiTx/OpenAiTx`. Não houve julgamento exaustivo dos demais candidatos, mas os termos efetivos mostram a perda da informação funcional da pergunta.

[Métricas públicas corrigidas](../../lab/bancada/history/benchmarks/0053/public-final/summary.json) e [tentativas/candidatos preservados](../../lab/bancada/history/benchmarks/0053/public-final/cases.json). Os registros originais estão em `public-collected/`. Nenhum projeto esperado foi encontrado nas buscas avaliáveis; não substituímos consultas por repositórios injetados para simular sucesso. A jornada remota integrada permanece **não avaliada**. Descoberta pública e análise local em snapshot congelado não foram multiplicadas em um suposto sucesso ponta a ponta no mesmo commit.

## Correções no executor e rastreabilidade

Os controles do integrador detectaram dois defeitos no avaliador, corrigidos pelo autor A sem ler corpus/testes:

1. Antes de executar o corpus, conteúdo adulterado com hash declarado preservado ainda podia contar como confirmado. A revisão r1 passou a recalcular o SHA dos bytes extraídos; a primeira falha e a implementação inicial foram preservadas.
2. Após a descoberta pública, erro dentro de `attempts` sem erro no topo era contado como resultado vazio avaliável. A revisão r2 passou a classificar esse caso como unknown. As respostas de rede foram preservadas e reavaliadas offline: a métrica inicial 0/3 foi corrigida para 0/2, mais um unknown. Não houve nova busca para obter resultado mais favorável.

O corpus/perguntas/alvos permaneceram intactos. A revisão r2 do executor ocorreu depois de o integrador observar resultados, portanto não é apresentada como avaliador final integralmente congelado antes da coleta. Há congelamentos encadeados e fontes anteriores para auditar a mudança. O autor A recebeu apenas os contraexemplos de schema/contabilidade.

Os testes sintéticos também tiveram correções explícitas de fixture: acrescentar o campo obrigatório FileEvidence e tratar `channels` como lista, sem alterar as expectativas de integridade. Os nove controles finais passaram, incluindo hash adulterado, erro fora do denominador, caminho sem símbolo, documentação sem prova e descoberta omitida. Os arquivos entregues inicialmente em `benchmark/0053` foram movidos mecanicamente para `lab/bancada/history/benchmarks/0053`, preservando bytes e hashes; o recibo de localização registra isso.

As métricas locais foram idênticas em duas coletas. Reavaliar as observações congeladas com o executor r2 também preservou as métricas locais. [Replicação](../../lab/bancada/history/benchmarks/0053/replication.json), [instruções e scripts](../../lab/bancada/history/benchmarks/0053/README.md), [verificação e hashes](../../lab/bancada/history/benchmarks/0053/verification.json).

Build, suíte Go completa e linteiro passaram; testes HTTP e descoberta pública usaram as permissões correspondentes. Os arquivos Go de produção mantêm os hashes da baseline. A entrega implementa a avaliação, não reivindica melhora do algoritmo. Durações são telemetria da ferramenta.

## Próxima correção indicada pelos dados

**Priorizar a construção da consulta e o fallback público:** preservar os termos funcionais, retirar a fórmula interrogativa e impedir que “Como criar” seja usado como substituto útil da intenção original. Candidatos de consultas relaxadas devem continuar sendo avaliados contra o pedido inteiro. Este benchmark fornece uma baseline para comparar essa mudança sem trocar perguntas.

Na recuperação local, o próximo experimento deve representar a operação pedida — monitorar/remover, interpolar, interseção/união — e avaliar o símbolo correto, não apenas palavras de objeto como `file` e `bits`. Adicionar vocabulário ou subir limites sem novos controles independentes repetiria o risco de otimizar exemplos conhecidos. Seleção/extração e custo remoto permanecem relevantes, mas não foram a perda dominante de funções nesta amostra.
