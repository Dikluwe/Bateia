# Execução do passo 0056

## Correção implementada
A hipótese derivada pede uma quantidade explícita de resultados ao provedor, limitada a200. Esse count não é número de repositórios. A seleção examina um conjunto limitado de repositórios antes de aplicar MaxResults aos candidatos com evidência, corrigindo o corte prematuro que podia ocultar um candidato válido atrás de outro sem evidência. A consulta literal, o catálogo0055-v2, os bloqueios e as regras lexicais não mudaram.

JSON e Markdown separam resultados recebidos, repositórios observados/examinados/adiados, elegíveis e selecionados. Diagnósticos distinguem ausência de snippets, termos faltantes, termos em arquivos diferentes e limite de candidatos. Truncamentos para240 caracteres e omissões após32 snippets por repositório são registrados. Apenas conteúdo retido participa da admissão; nome, caminho e descrição não substituem conteúdo. O relatório limita a apresentação dos diagnósticos e informa o que omitiu, mantendo JSON integral.

Cobertura do índice permanece unknown, mesmo quando o provedor informa matchCount. Não há downloads adicionais nem execução de código. Analisador evidence-v17-discovery-coverage e checkpoint live0056-v1 invalidam os resultados da estratégia antiga.

## Observação pública

| Medida | 0055 | 0056 |
|---|---:|---:|
| Projetos esperados admitidos no corpus | 0/3 | 0/3 |
| Funções-alvo recuperadas localmente | 1/6 | 1/6 |
| Repositórios examinados na hipótese de bits | 1 | 3 |
| Candidatos admitidos nessa hipótese | 0 | 3 |

Aumentar candidatos admitidos não demonstra melhora semântica. A pergunta de bits solicitou count:100. O provedor informou100 correspondências e devolveu3 registros, de3 repositórios; todos foram examinados e selecionados. Zero repositórios foram adiados pelo limite local. O projeto esperado bits-and-blooms/bitset não estava nesse conjunto. Nesta execução a perda do alvo ocorreu antes da admissão local; não há evidência de que baixar o conteúdo desses três repositórios recuperaria esse alvo.

Os três registros são periph/devices, soypat/gsdf e majestrate/XD. Houve5 snippets omitidos em periph/devices pelo teto32; zero truncamentos registrados. Os demais não tiveram omissões. As outras duas perguntas públicas permaneceram literais pelos bloqueios de negação/sintaxe já existentes.

[Métricas e diagnósticos](../../lab/bancada/history/benchmarks/0056/comparison.json), [resposta completa da pergunta de bits](../../lab/bancada/history/benchmarks/0056/collected/case-07/discovery.stdout), [recibos de execução](../../lab/bancada/history/benchmarks/0056/collected/receipts.json).

## Relação com a intenção original
**Fatos observados nos snippets:** periph/devices mostra r.Intersect(i.Rect) em ssd1306/image1bit/image1bit.go. soypat/gsdf descreve interseção de SDFs em operations.go. majestrate/XD contém Reserved.Intersect em lib/bittorrent/handshake.go, com comentário de que limpa bits ausentes do outro operando. Esses arquivos/trechos estão no recibo citado acima; não foram executados.

**Inferência:** a ocorrência de “intersect” e “bit” no mesmo arquivo reúne domínios distintos. O candidato de handshake tem indício relacionado à operação sobre bits, mas isso não demonstra que cria um conjunto novo preservando os originais, como a pergunta exige. Não classificar qualquer um deles como solução exata. O alvo fixo do benchmark não mede todas as alternativas potencialmente úteis.

**Hipótese para próxima validação:** verificar a relação entre ação, objeto e comportamento na implementação candidata (mesmo símbolo, entradas, retorno e mutação), mantendo aquisição limitada de conteúdo quando necessária. Antes de elevar limites novamente, distinguir insuficiência do índice/ordenação de perda por snippets. Nesta coleta não houve corte local de repositórios.

**Não verificado:** cobertura total do índice, comportamento executável dos candidatos, generalização para perguntas novas e jornada remota integrada. O corpus é conhecido e de calibração; não holdout. Nenhum resultado de busca é prova de implementação.

## Validação e custódia
Tekt A/B executado sem atestação física. A produziu domínio/adaptador sem ler testes/corpus; B escreveu testes sem ler produção/patch/corpus; root escreveu contrato, apresentação e checkpoint e integrou. Prompts e contrato antecederam código; hashes congelados antecederam testes. Um helper de admissão obsoleto foi removido antes dos gates, mantendo uma política ativa.

Os controles HTTP verificam que, com MaxResults1, um primeiro candidato insuficiente não elimina um candidato válido posterior; também verificam limites, matchCount ausente/null, divisão entre arquivos, overflow, truncamento Unicode e cap32. Os controles de apresentação verificam10 diagnósticos×3trechos sem mutar JSON; retomada rejeita0055-v2 e aceita live0056-v1. A expectativa antiga de serialização foi explicitamente atualizada para count:200 na fixture com20 candidatos, sem alterar critérios de admissão.

Uma única rodada pública foi executada. Corpus/runner/expectativas0053 permaneceram intactos; reavaliação offline reproduziu as métricas. Tempos são telemetria de Bateia/provedor, não benchmarks dos projetos.

[Verificação final e hashes](../../lab/bancada/history/benchmarks/0056/verification.json), [reprodução](../../lab/bancada/history/benchmarks/0056/README.md). A sintaxe de count segue a [documentação oficial Sourcegraph](https://sourcegraph.com/docs/code-search/queries).

A primeira suíte completa encontrou um identificador incorreto no checkpoint: implementação0056-v1, contrato/teste live0056-v1. Corrigida somente a string de identidade da retomada e preservados os testes. Logs e versão anterior do arquivo foram arquivados. A coleta pública precede essa correção, que não altera consulta, admissão ou análise local; não houve nova busca para escolher resultados. O hash do binário medido está em collected/receipts.json e o do binário final no recibo de verificação.

Verificação final: suíte Go completa, compilação, linteiro sem violações e nove controles do avaliador passaram. Freezes finais mantêm seus hashes; apenas o arquivo de checkpoint mudou após a medição pública, com revisão preservada.
