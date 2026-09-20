# Execução 0050 — Recuperação remota de implementações

2026-09-08. O [passo 0050](../passos/0050-recuperar-implementacoes-remotas.md) implementa o primeiro gargalo do diagnóstico 0049: candidatos GitHub passam pela recuperação de definições Go antes da seleção dos arquivos. Antes, o Router devolvia `not_evaluated` para candidatos remotos.

## Mudança entregue

O adaptador GitHub implementa `ImplementationSearcher`; o Router delega ao provedor selecionado quando existe suporte. A coleta exige commit SHA completo, lê a árvore nesse commit e busca blobs pelos identificadores da árvore. Verifica tamanho, encoding e SHA-1 do objeto Git antes de enviar os bytes ao parser. Cada arquivo admitido recebe também SHA-256 no relatório. Não faz clone, executa código ou troca a revisão por HEAD.

Tetos remotos: 64 tentativas de blob, 4 MiB de fonte, 512 KiB por arquivo, 24 hits e resposta de árvore limitada a 8 MiB. Orçamentos menores prevalecem. Respostas HTTP e envelopes JSON têm limites próprios; `BytesRead` contabiliza fonte decodificada, inclusive descartada, e não representa tráfego HTTP total. Um servidor que envia conteúdo além do declarado pode fazer esse contador ultrapassar o orçamento; o excesso não é admitido e a coleta é limitada. Falhas, opacidade e esgotamento não produzem conclusão de ausência.

O ranking, glossário e interpretação permanecem iguais. A revalidação posterior continua removendo hits sem arquivo extraído ou com bytes divergentes. EvidenceVersion passa a `evidence-v11-remote-retrieval`; checkpoints local/live usam `implementation_search: 0050`, impedindo reaproveitar silenciosamente coleta anterior sem esta etapa. Replay permanece baseado nos bundles da origem.

## Evidência executada

**Ensaio público no commit fixo** `golang/snappy@43d5d4cd4e0e3390b0b645d5c3ef1187642403d8`, com consulta `decode block corrupt`: oito arquivos íntegros, 33.163 bytes de fonte, `Decode` em primeiro lugar. Os 24 hits, pesos, linhas e hashes são idênticos à execução local anterior com os mesmos arquivos. A extração e releitura confirmaram os 24 hits. Status `partial` por corte de resultados, não por leitura incompleta dos oito arquivos.

[Coleta, seleção e revalidação](../../lab/bancada/history/benchmarks/0050/live-verification.json); [comparação local/remota](../../lab/bancada/history/benchmarks/0050/local-remote-comparison.json); [runner do integrador](../../lab/bancada/history/benchmarks/0050/live-runner.go.txt). A primeira tentativa foi bloqueada pelo sandbox; o ensaio foi repetido com permissão de rede e completou. Código externo somente lido.

**CLI completa:** `search -repos golang/snappy -query 'decode block corrupt' -lang go -max 1 -depth 1 -format json` resolveu o commit `9ae09f520e93143c782aa6088d18a5e300c69dd0`. Recuperou oito arquivos na etapa AST e reteve 22 hits após extração; `Decode` permaneceu primeiro. Dois hits de `cmd/snappytool/main.go` foram removidos porque o arquivo não foi extraído. O relatório registra explicitamente esses descartes e permanece parcial. Esse ensaio usa outro commit e não é a comparação de igualdade local/remota.

[Resultado da CLI](../../lab/bancada/history/benchmarks/0050/cli-observations.json). As cópias dos recibos substituem conteúdos integrais de arquivos por hashes/tamanhos; mantêm relatórios de busca e intervalos. Tempos registrados são telemetria do Bateia, não benchmarks do Snappy.

## Verificação

Autor B escreveu testes sem ler a candidata: seis testes raiz e 21 subcasos passaram na primeira execução, sem ajuste de esperados. Cobrem consulta funcional em caminho neutro, igualdade local/remota, encaminhamento pelo Router, seleção/revalidação, identidade e conteúdo divergentes, árvore truncada, symlinks, budgets, 403/429 e cancelamento prévio. [Avaliação independente](../../lab/bancada/history/benchmarks/0050/independent-evaluation.md) e [congelamento](../../lab/bancada/history/benchmarks/0050/tests-freeze.json).

Build de `./04_wiring`, suíte `go test ./...`, linteiro cristalino e `git diff --check` passaram. Os testes HTTP foram executados com permissão para servidores locais. [Comandos e hashes](../../lab/bancada/history/benchmarks/0050/verification.json).

Regime Tekt A/B executado sem atestação de isolamento físico. Integrador escreveu o contrato, Router e versões; A escreveu adaptador/prompt; B escreveu testes e veredito. Filesystem compartilhado. Não há selo de refinamento ou teste de mutações. A avaliação B declara limites: cancelamento durante transferência não foi exercitado; não saturou os tetos de 4 MiB e 24 hits distintos, embora tenha verificado seus valores e orçamentos menores.

## O que melhora e o que continua aberto

O ganho observado é acesso a implementações remotas verificáveis, antes inexistente nessa etapa. Isso não demonstra ganho universal de pertinência: a consulta técnica foi escolhida para isolar a disponibilidade da coleta. As perguntas naturais de compressão/descompressão, o índice de documentação e a cobertura do catálogo de capacidades continuam pendentes. `partial` também não deve ser lido como solução demonstrada.

Há custo de até uma árvore e 64 blobs por chamada da nova etapa, além de metadados, triagem e releitura do funil existente. O adaptador não acrescenta cache de blobs; consumo agregado de cota, seleção quando muitos arquivos não cabem e os descartes na extração são limites operacionais para o próximo refinamento. Não se presume que habilitar a etapa remota torne todo repositório analisável integralmente.
