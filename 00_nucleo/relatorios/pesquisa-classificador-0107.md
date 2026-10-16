# Pesquisa 0107 — Qualidade do classificador pelo próprio Bateia

**Data**: 2026-10-06  
**Estado**: PESQUISA EXPLORATÓRIA CONCLUÍDA; QUALIDADE SEMÂNTICA NÃO MEDIDA  
**Contrato**: [passo 0107](../passos/0107-pesquisar-qualidade-do-classificador-com-bateia.md)  
**Evidência**: [bancada e reprodução](../../lab/bancada/pesquisa/0107/README.md)

## Conclusão

A primeira fatia recomendada é avaliar propriedades específicas e a autoridade
das afirmações antes de ampliar o catálogo de mecanismos. A agregação atual
confere identidade, bytes e localizadores, mas depende do julgamento semântico
fornecido. Seu contrato deve continuar explícito. Um estudo de compreensão
precisa avaliar o produtor desse julgamento e a jornada completa, além do agregador.

O Bateia foi usado para expandir a ideia, pesquisar catálogos acadêmicos e GitHub,
triar candidatos, tentar coleta, avaliar bytes congelados, calcular métricas e
investigar insights. A pesquisa adicional em fontes primárias e a síntese foram
feitas pelo agente, não automaticamente pelo Bateia. Não houve decisão de adoção.

## O que a execução encontrou

`idea plan` produziu 11 expansões. Seis consultas manuais discriminadas ao
GitHub produziram seis candidatos em SHAs imutáveis. CheckList, Snorkel e
FActScore são sementes conhecidas da pesquisa anterior: essa busca confirma
implementações e não constitui descoberta independente nem revisão exaustiva.

Quatro consultas acadêmicas produziram oito tentativas: dois sucessos JOSS,
duas respostas JOSS 429 e quatro falhas de parsing zbMATH por conteúdo não JSON.
Os seis candidatos JOSS incluem três de incerteza científica e dois de genética,
além do Causal Testing Framework. A expressão mutation testing teve ambiguidade
de domínio; descoberta de software científico não equivale a busca geral de artigos.

A triagem de catálogo marcou CheckList, Snorkel, FActScore e Clair ADVANCE com
score 0,7. Clair não foi selecionado: correspondência textual com static analysis
não demonstra adequação à pergunta sobre mecanismos. GraphicsFuzz ficou REJECT
e Torch-Uncertainty REVIEW; ambos permaneceram sem inspeção de código. Esses
vereditos são sinais para pesquisa, não medidas de qualidade dos projetos.

A coleta Bateia de quatro candidatos ficou incompleta por 403 e falha de
resolução do CTF. O recibo público confirmou quota core 60/60 usada, restante
zero. Dez arquivos dos três SHAs já conhecidos foram então lidos por fallback
externo público, totalizando 158.844 bytes. Os bundles de revisão mantêm falhas,
targets e complete=false originais. A origem complementar não foi atribuída
ao coletor Bateia. CTF foi aprofundado pelo artigo; seu código não foi coletado.

## Soluções observadas e limites de transferência

| Referência | Observação localizada | Aplicação proposta ao Bateia | Limite |
| --- | --- | --- | --- |
| CheckList | Tipos MFT, INV e DIR; expectativas explícitas | Separar casos mínimos, invariâncias e mudanças esperadas | Os autores dos casos precisam justificar as expectativas |
| Snorkel | Diagnóstico por função: cobertura, sobreposição e conflito | Instrumentar a qualidade das regras antes de combinar votos | Dependência entre indícios pode distorcer confiança |
| FActScore | Afirmações atômicas, decisão de suporte por afirmação e respond_ratio | Decompor descrições e conferir sustentação individual | Verificador por recuperação/modelo estima suporte textual, não execução |
| Causal Testing Framework | Propriedades, modelo causal e dados de execução separados | Desenhar intervenções que identifiquem o que altera uma conclusão | Inferência depende de premissas causais e adequação dos dados |
| Classificação seletiva | Relação entre erro das conclusões aceitas e cobertura | Medir excesso de confiança e abstenção separadamente | Confiança precisa de avaliação; score lexical não é probabilidade calibrada |

A implementação CheckList exige labels ou uma função expect no teste mínimo,
oferece expectativa de invariância e recebe explicitamente a expectativa
direcional. Para nós, reordenar testemunhas não deve alterar a classe; remover
uma guarda necessária pode alterar a propriedade inferida. Uma transformação
que muda bytes precisa atualizar hashes e localizadores antes de testar semântica.
[Tipos no SHA coletado](https://github.com/marcotcr/checklist/blob/4e6e5e33a26f30c20ed602b2050f6c73e123cc23/checklist/test_types.py).

Snorkel distingue dados sem rótulo, sobreposição e conflito. O LabelModel
coletado declara independência condicional das funções dado o rótulo verdadeiro.
Para o Bateia, hashes diferentes provam diferença de conteúdo, não independência
semântica. A recomendação é começar pelo diagnóstico das regras; não importar
um modelo probabilístico supondo que todos os indícios sejam independentes.
[Análise](https://github.com/snorkel-team/snorkel/blob/45824f9867228c6cb6c9fdc2afd2b7e5bb7fd0b3/snorkel/labeling/analysis.py),
[premissa do modelo](https://github.com/snorkel-team/snorkel/blob/45824f9867228c6cb6c9fdc2afd2b7e5bb7fd0b3/snorkel/labeling/model/label_model.py).

FActScore conserva decisões individuais e proporção de respostas, além de
uma média de suporte. Isso inspira uma avaliação por afirmação: limite de
workers, configuração padrão, reserva antes de efeito, rejeição, bloqueio e
cancelamento são propriedades separadas. Sua recuperação e verificação por
modelo não podem ser tratadas como prova automática de semântica de código.
[Scorer no SHA coletado](https://github.com/shmsw25/FActScore/blob/f28272deffcf33efc1f1117d5479c10bb75221a9/factscore/factscorer.py).

O CTF foi uma pista nova descoberta pelo conector JOSS. Seu artigo usa modelo
causal, propriedades esperadas e dados de execução para estimar efeitos.
A adaptação inicial proposta é mais simples: intervenções controladas nos
insumos do classificador. Não se infere causalidade apenas desenhando um grafo;
o uso estatístico do framework exigiria dados e premissas próprios.
[Artigo primário](https://joss.theoj.org/papers/10.21105/joss.07739.pdf).

A literatura de classificação seletiva trata erro e cobertura conjuntamente.
Esse princípio pode orientar métricas sem exigir adotar uma rede neural.
INCONCLUSIVE precisa ter motivo: evidência insuficiente, coleta incompleta,
conflito, categoria desconhecida e falta de aplicabilidade são situações distintas.
[SelectiveNet](https://proceedings.mlr.press/v97/geifman19a.html).

Uma revisão de 2025 sobre modelos de código reúne 45 trabalhos e examina
transformações que preservam semântica, como renomeação. Isso sustenta incluir
robustez a representação na bancada, sem confundir estabilidade com correção:
um classificador consistentemente errado também pode ser estável.
[Revisão](https://arxiv.org/abs/2507.22610).

## O que o Bateia conseguiu avaliar

Foram aplicadas seis regras HIPOTÉTICAS para localizar indícios nos três
bundles de revisão. As regras foram informadas pelos próprios candidatos:
não são referência independente e seus resultados não medem generalização.
Entre 18 avaliações de capacidades e três de alinhamento, quatro ficaram
PARTIAL/HYPOTHESIS_ONLY e 17 INCONCLUSIVE; nenhuma ficou PROVEN.

`idea evidence evaluate` retornou precisão null, sample_sufficient=false,
inconclusive_rate=17/21, source_trace_rate=1 e budget_completion=0. Os campos
de referência estão preenchidos, mas a coleta Bateia permanece incompleta.
Não fornecemos julgamentos independentes nem fabricamos uma taxa de acerto.

Na jornada `idea insights analyze`, o resultado foi INCONCLUSIVE, com zero
mecanismos e 24 lacunas RULE_NOT_IN_MECHANISM_CATALOG. As regras de avaliação
pesquisadas não pertencem às quatro classes fechadas. Isso mostra uma limitação
da interface: observações úteis deveriam sobreviver à falta de categoria adequada.
Não demonstra que as metodologias pesquisadas sejam inaplicáveis ao Bateia.

Foi observada também uma limitação operacional: `--rules` carrega o catálogo
padrão antes de desserializar o arquivo. Omitir hash conserva o hash padrão
incompatível; o arquivo desta bancada declara hash vazio explicitamente.
Nenhum código produtivo foi alterado para contornar essa condição.

## Sonda controlada do classificador atual

Oito casos foram executados pelo CLI sobre a fonte congelada 0106. Reordenação
e duplicação preservaram os mecanismos; remoção de avaliações e regra desconhecida
produziram INCONCLUSIVE sem mecanismo; adulteração de bytes foi recusada.

Com bytes inalterados, promover somente o status declarado produziu mecanismo
PROVEN. Trocar também os RuleIDs produziu APPEND_ONLY_LOG_RECOVERY PROVEN sobre
os mesmos trechos de concorrência. As hipóteses continuaram no máximo PROPOSED
e o alvo UNKNOWN. A sonda caracteriza a confiança em avaliações declaradas:
não é uma taxa de falsos positivos medida contra corpus independente, nem quebra
do contrato atual de agregação. Demonstra por que a validação semântica precisa
ser estudada em um responsável próprio.

## Fatia recomendada para estudar qualidade

Começar com propriedades de concorrência e sobrecarga, usando fonte, configuração
e caminho específico como unidade, em vez do repositório inteiro. Cada afirmação
precisa conservar propriedade, escopo, condição, evidências, contraevidências,
procedimento de avaliação e autoridade. A categoria de mecanismo pode permanecer
desconhecida enquanto há observações úteis. Essa é uma proposta, ainda não implementada.

Construir 12 casos-base exploratórios que distingam limite de workers, fila
bufferizada, modo ilimitado, reserva antes/depois de efeito, rejeição, bloqueio
e alcance do cancelamento. Derivar pares de invariância e mudança necessária.
Os casos derivados são correlacionados; não contá-los como amostras independentes.
Revisar os rótulos separadamente do produtor, registrar discordância e congelar
uma parte para avaliação sem ajuste posterior. Nenhuma referência independente
desse tipo foi produzida nesta pesquisa.

Comparar a avaliação declarada atual com um produtor estrutural limitado e,
somente depois, uma proposta assistida por IA. Medir localização de indícios,
correção de cada afirmação, preservação de condições, excesso de autoridade,
erro entre conclusões aceitas, cobertura e fidelidade da descrição. Manter
resultados separados por produtor e pela jornada composta. Casos fora do escopo
devem gerar abstenção motivada, não ausência universal.

O avaliador existente ainda não cobre toda essa necessidade: evidence_precision
revisa somente resultados PROVEN; source_trace_rate conta campos de referência;
sample_sufficient depende da contagem de dez PROVEN, não da quantidade de
julgamentos independentes. Um protocolo novo precisa denominadores, revisão,
escopo imutável e tratamento de PARTIAL explícitos, preservando o formato legado.

Sob Tekt, regras e contratos puros pertencem a L1; coleta, execução de ferramentas
externas e recibos concretos ficam em L3; entrada/apresentação em L2; composição
em L4. O corpus, intervenções e medições ficam no Lab. Não há necessidade de
substituir o agregador ou expandir classes antes de demonstrar o ganho dessa fatia.

## Verificação e limites

Os cinco estágios offline reproduziram os arquivos byte a byte; os dez arquivos
externos conferiram com hashes e recibos. A sonda de oito casos passou, lint
cristalino e diff check passaram. O código Go permanece o verificado no passo
0106; não foi alterado nem a suíte Go repetida nesta pesquisa de dados e documentação.

Não foram executados candidatos, instaladas dependências, treinados modelos,
consultadas APIs de IA ou alteradas políticas produtivas. Esta pesquisa reduz
incerteza sobre como desenhar a avaliação; não estabelece a qualidade semântica
do classificador nem um ganho comprovado para o alvo.
