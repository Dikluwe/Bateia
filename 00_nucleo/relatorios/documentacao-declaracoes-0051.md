# Execução 0051 — Documentação vinculada a declarações

2026-09-08. O [passo 0051](../passos/0051-recuperar-documentacao-de-declaracoes.md) torna pesquisáveis os comentários Doc de funções e métodos Go. O ganho é localizar descrições funcionais antes ignoradas, com autoridade documental explícita. Não é tradução nem comprovação da implementação descrita.

## Implementação

O parser associa termos dos grupos Doc às declarações pelo AST. Comentários de pacote, soltos, internos ao corpo, literais, tipos e arquivos gerados não originam esse canal. Consultas exclusivamente por símbolo permanecem no canal de código. O glossário e o IDF do código não mudaram.

Os hits documentais ocupam somente slots restantes depois dos hits de código e chamadas. Têm `retrieval: documentation`, score pela quantidade de termos distintos (teto 99), intervalos da declaração e do comentário, símbolo e SHA-256. Omissão pelo orçamento fica explícita. Não deslocam os hits de código nem criam pesos de IDF artificiais.

A revalidação após extração confere hash, símbolo, intervalo da declaração e grupo Doc exato pelo AST. Um intervalo forjado, mesmo dentro do arquivo, não basta. A projeção separa `documentation_sources` de `sources`; quando só há documentação, usa `documentation_located`. O Markdown do hit mostra até vinte linhas do comentário, com aviso de que se trata de descrição do autor. A seção geral de arquivos extraídos continua exibindo o conteúdo disponível, separadamente.

Versão de evidência `evidence-v12-declaration-documentation`; checkpoint local/live `implementation_search: 0051`. O canal é compartilhado entre os provedores local e GitHub. Replay conserva as fontes congeladas.

## Resultados observados

Ensaios exploratórios do integrador, nos mesmos arquivos antes/depois e sem executar código externo:

| Consulta | Baseline 0050 | Candidata 0051 |
| --- | --- | --- |
| Snappy: `decompresses framing` | Zero hits | `NewReader` primeiro, seguido de duas correspondências documentais parciais. |
| bbolt: `multiple concurrently` | Zero hits | `Begin` primeiro; também retorna outras declarações com correspondência parcial. |
| Snappy: `Descomprimir um bloco e rejeitar dados corrompidos` | Zero hits | Zero hits; continua `not_verified`. |

[Comparação com consultas, termos, arquivos e hashes](../../lab/bancada/history/benchmarks/0051/cli-comparison.json). A escolha das duas primeiras consultas veio da inspeção dos comentários; são controles diagnósticos, não um holdout de perguntas de usuários. O ganho demonstra acesso ao canal documental, sem estimar precisão geral ou ganho em português. A terceira consulta preserva visível a limitação de vocabulário.

No Snappy, `NewReader` corresponde a `decompresses` e `framing`; `NewBufferedWriter` e `crc` correspondem apenas a `framing`. Essa dispersão mostra que soma lexical ainda não verifica a relação funcional completa entre os termos. `documentation_located` não significa que todos os resultados resolvem a pergunta.

**Ensaio remoto público:** `golang/snappy@43d5d4cd4e0e3390b0b645d5c3ef1187642403d8`, mesma consulta `decompresses framing`. Oito arquivos lidos, três hits documentais, idênticos aos locais; extração e revalidação confirmaram os três. Estado `complete` da recuperação significa término no orçamento desse recorte, não solução comprovada. [Recibo remoto](../../lab/bancada/history/benchmarks/0051/remote-verification.json) e [runner](../../lab/bancada/history/benchmarks/0051/remote-runner.go.txt). Rede exigiu a permissão correspondente do ambiente.

## Verificação e autoria

Build de `./04_wiring`, `go test ./...`, `crystalline-lint .` e `git diff --check` passaram. A suíte completa incluiu servidores HTTP locais com permissão. Replay reproduziu byte a byte assessment, decisão, alternativas, Markdown e bundles. A autoavaliação manteve as dez capacidades inconclusivas. [Replay](../../lab/bancada/history/benchmarks/0051/replay.json) e [recibo final](../../lab/bancada/history/benchmarks/0051/verification.json).

Tekt A/B sem atestação de isolamento físico: integrador escreveu contrato, revalidação, apresentação e versões; A escreveu domínio/parser/projeção; B escreveu casos e veredito sem ler corpos da candidata. Filesystem compartilhado. O [veredito B](../../lab/bancada/history/benchmarks/0051/independent-verdict.md) registra as entradas e limites.

Na primeira execução B, núcleo passou e dois testes de Markdown falharam por observarem o relatório inteiro, incluindo arquivos gerais. O contrato restringia o trecho do hit documental. A revisão corrigiu esse escopo, preservou valores esperados e arquivou primeira falha/teste/congelamento. Os testes do núcleo foram movidos mecanicamente para o arquivo companheiro exigido pelo linteiro. A segunda execução passou após novo congelamento. Não se apresenta essa revisão de observabilidade como defeito corrigido em produção ou aprovação na primeira tentativa.

## Limites e continuidade

Documentação pode estar incorreta, desatualizada ou descrever uma negação. Continua sendo uma declaração do autor. O canal só usa espaço restante; orçamento ocupado por código impede novas fontes documentais. Não foram modificados o catálogo de capacidades, tradução, análise de ação/objeto/condições ou política conjunta de ranking. Os recibos não são uma nova avaliação exaustiva dos corpora antigos; preservação de código é verificada pelo contrato e controles, sem alegar que o conjunto total de hits permaneça igual.

O próximo refinamento deve tratar a ligação entre termos da pergunta e mecanismos, preservando estas fontes como evidência de origem. A pergunta portuguesa ainda falha e os matches parciais continuam visíveis: esses são critérios mais úteis para o próximo experimento do que apenas aumentar o número de resultados. Durações são telemetria do Bateia, não benchmarks dos projetos.
