# Avaliação0046 — recuperação em fontes reais

Autor B, regime A/B lógico executado sem atestação física de isolamento. Perguntas, alvos por inspeção manual e snapshots foram congelados em `real/freeze.json` antes da primeira medição. Não foram lidos corpos do parser de produção nem testes existentes para construir os rótulos; somente API pública de Rank e tipos domain. A leitura das fontes externas foi estática. Nenhum código dos projetos externos foi executado, importado ou adicionado como dependência.

## Corpus e custódia

Amostra de conveniência de três módulos disponíveis integralmente no cache local; não é amostra representativa. Cobra foi considerada, mas somente seus metadados estavam no cache; a substituição por OAuth2 ocorreu antes do congelamento. Cada projeto tem duas perguntas funcionais em português e uma negativa. As perguntas não contêm o símbolo esperado. O conjunto de alvos rotulados é intencionalmente pequeno: arquivos auxiliares retornados podem ser úteis, e não foram rotulados como irrelevantes.

| Projeto | Domínio | Versão / commit da origem | Licença | Go sem testes / bytes |
|---|---|---|---|---:|
| bbolt | Banco local chave/valor | v1.4.0 / `0d516858c4bad9ad7c7ccb637adc32bcd78ca8cb` | MIT | 64 / 313556 |
| golang/snappy | Compressão de blocos e streams | v1.0.0 / `43d5d4cd4e0e3390b0b645d5c3ef1187642403d8` | BSD-3-Clause | 8 / 33163 |
| golang.org/x/oauth2 | Cliente de autorização OAuth2 | v0.36.0 / `4d954e69a88d9e1ccb8439f8d5b6cbef230c4ef9` | BSD-3-Clause | 62 / 227503 |

`real/manifest.json` registra URL de origem, Origin.Hash/Ref, SHA256 do ZIP, checksum Go ziphash, licença e SHA256/bytes de cada arquivo. As fontes são snapshots de módulos distribuídos, não necessariamente árvores Git completas: testes e arquivos não Go não entram no ranking; variantes de arquitetura, pacotes internos e comandos Go entram, sem filtragem por build tags. Assembly de Snappy, por exemplo, não é analisado. Os ZIPs permanecem no cache, sem cópia dos projetos para produção.

`real/questions.json` guarda texto exato, justificativa e path/símbolo/intervalo esperado. Para negativas, a ausência é rótulo manual dos fluxos e recortes declarados; não prova de inexistência em todo o repositório. `real/freeze.json` registra hashes e horário UTC. O relógio efetivo retornou 2026-09-08 UTC (a sessão foi iniciada em 2026-09-07); usamos horário observado, sem retrodatação. Hashes identificam os artefatos, não atestam custódia independente ou autenticidade da origem remota.

## Método e orçamento

Execute, na raiz ou por caminho absoluto:

```sh
python3 00_nucleo/benchmark/0046/real/run.py baseline
python3 00_nucleo/benchmark/0046/real/run.py candidate
python3 00_nucleo/benchmark/0046/real/run.py final
```

O runner recusa sobrescrever uma medição existente. Para reprodução preservando os recibos originais, copie os cinco arquivos `run.py`, `manifest.json`, `questions.json`, `freeze.json` e o recibo desejado para custódia antes de remover apenas o resultado local que será regenerado. É necessário o mesmo checkout Bateia e os ZIPs nos paths declarados; rede fica desligada (`GOPROXY=off`, `GOSUMDB=off`). O adaptador Go temporário fica em `/tmp` e importa apenas API pública Bateia. Não há Go em L0.

Varredura simulada sobre arquivos Go em ordem de path: 2048 arquivos, 32 MiB total, 512 KiB por arquivo; nenhum arquivo foi excluído por orçamento nesta amostra. Rank recebe limite de 24 hits. @3 e @12 são os primeiros caminhos **distintos** desses hits: simulam slots de extração, não executam provedor, leitura final, revalidação de hash ou CLI. Não há orçamento adicional de bytes de extração medido aqui. Os JSONs registram a varredura efetiva, relatório cru, contagens, perdas por pergunta e hashes das fontes Bateia compiladas.

Recall = arquivos esperados encontrados / arquivos esperados rotulados. Precisão rotulada = arquivos esperados encontrados / caminhos retornados no corte. Negativas não entram nessas duas métricas: medimos presença de indícios retornados separadamente. Agregado usa macro recall por requisito positivo e micro precisão sobre seus caminhos. Um indício negativo não é chamado de capacidade implementada nem de falso positivo semântico demonstrado.

## Baseline0045

| Requisito | Alvo manual | Hits | Recall @3 | Recall @12 | Paths @3 / @12 | Diagnóstico observável |
|---|---|---:|---:|---:|---:|---|
| bbolt-write | db.go:891–916, DB.Update | 24 | 0 | 1 | 3 / 6 | db.go aparece em 5º caminho; três slots perdem o alvo |
| bbolt-read | db.go:922–948, DB.View | 24 | 0 | 0 | 3 / 12 | alvo fora do relatório; variantes de mmap ocupam o início com `read` |
| snappy-compress | encode.go:20–43, Encode | 4 | 0 | 0 | 1 / 1 | encode_other.go localizado por `bytes`, seguido por três direct_call no mesmo arquivo |
| snappy-decompress | decode.go:57–76, Decode | 0 | 0 | 0 | 0 / 0 | nenhuma correspondência retornada para a formulação funcional |
| oauth2-header | transport.go:32–56, Transport.RoundTrip | 24 | 0 | 0 | 3 / 12 | alvos concorrentes com auth/http/token; alvo manual fora dos hits |
| oauth2-exchange | oauth2.go:222–235, Config.Exchange | 24 | 0 | 0 | 3 / 12 | vários fluxos OAuth com auth/token; alvo manual fora dos hits |

Baseline: recall macro **0/6 = 0% @3** e **1/6 = 16,67% @12**. Precisão micro rotulada **0/13 = 0% @3** e **1/43 = 2,33% @12**. Quatro dos seis requisitos positivos esgotam 24 hits e o relatório sinaliza partial; os dois de Snappy não saturam. Todos os alvos estavam disponíveis para varredura: a perda não foi causada pelos tetos de leitura do adaptador. A presença de um path não garante recuperação do símbolo esperado; estas métricas são estritamente por arquivo.

| Negativa | Indícios retornados | Hits | Observação |
|---|---:|---:|---|
| bbolt: enviar email quando transação falha | não | 0 | não prova ausência universal |
| Snappy: persistir blocos em transação SQL com rollback | não | 0 | não prova ausência universal |
| OAuth2: comprimir corpo HTTP em blocos Snappy | sim | 24 | hits lexicais `http`; a capacidade de compressão não está demonstrada |

Há indícios em **1/3 negativas**. O último caso evidencia por que a projeção por requisito deve explicitar fundamento lexical e termos observados, evitando converter correspondência parcial em satisfação funcional.

## Diagnóstico antes de mudar ranking

Os resultados públicos mostram duas classes distintas: formulação sem correspondência útil (descompressão), e termos muito compartilhados levando a fontes concorrentes (read/auth/token/http). O ranking limitado ainda introduz perda de arquivo no corte de extração, demonstrada concretamente por DB.Update em 5º caminho. Um arquivo auxiliar encontrado para compressão pode conter parte real do trabalho, embora não seja o alvo manual: isto limita a interpretação da precisão.

Os dados não isolam uma correção de pesos que resolva as classes sem regressão. Aumentar um limite ou privilegiar os nomes dos arquivos esperados seria ajuste à amostra. Não se solicitou alteração dos rótulos para acomodar o resultado. A recomendação prévia deste ensaio foi conservar o ranking na ausência de uma falha generalizável isolada e registrar as limitações, melhorando a explicação por requisito sem apresentar ganho artificial.

## Tentativa experimental descartada

A tentativa de diversificação foi executada em `real/candidate.json` após a liberação pelo autor A. Ela foi posteriormente descartada e não integra a entrega final. Este autor comunicou diversificação por rodadas de caminhos antes da expansão; o avaliador B não inspecionou a implementação. `real/comparison.json` confirma igualdade de freeze, runner, adaptador, orçamento e varredura. O diagnóstico anterior permanece registrado como avaliação prévia, sem ajuste de perguntas ou alvos.

| Requisito | Recall @3 antes → depois | Recall @12 antes → depois | Caminhos no relatório antes → depois | Efeito no alvo |
|---|---:|---:|---:|---|
| bbolt-write | 0 → 0 | 1 → 1 | 6 → 13 | db.go continua em 5º |
| bbolt-read | 0 → 0 | 0 → 0 | 13 → 20 | db.go passa de ausente para 14º, fora da extração de 12 |
| snappy-compress | 0 → 0 | 0 → 0 | 1 → 1 | encode.go continua ausente |
| snappy-decompress | 0 → 0 | 0 → 0 | 0 → 0 | decode.go continua ausente |
| oauth2-header | 0 → 0 | 0 → 0 | 12 → 17 | transport.go continua ausente |
| oauth2-exchange | 0 → 0 | 0 → 0 | 12 → 16 | oauth2.go continua ausente |

Todos os hits por pergunta mantêm a mesma contagem; varia sua composição e diversidade. O recall macro permanece **0% @3 e 16,67% @12**. A precisão micro rotulada permanece **0% @3** e passa de **1/43 = 2,33%** para **1/49 = 2,04% @12**, porque mais caminhos não rotulados entram no corte de bbolt-write. Isso não demonstra queda de relevância real: os novos arquivos não foram julgados irrelevantes. As três negativas mantêm exatamente sua contagem de hits e seus caminhos, com indícios em **1/3**.

O único ganho sobre um alvo manual é intermediário: a candidata passa a incluir DB.View pelo arquivo db.go no relatório limitado de bbolt-read, mas em 14º caminho; portanto não melhora os cortes aceitos de 3 ou 12. Mesmo um relatório com mais arquivos ainda pode perder uma função esperada. A avaliação pública de casos mínimos e preservação de símbolos exatos0045 pertence aos controles separados, não foi inferida deste corpus.

**Veredito limitado:** comparação reproduzível concluída; nenhum ganho de recall rotulado nos cortes de extração, nenhuma mudança na taxa de negativas com indícios. As perdas funcionais em português e a concorrência entre termos genéricos permanecem. O resultado não sustenta afirmar recuperação funcional geral, capacidade satisfeita, ou ausência universal de capacidades. A projeção por requisito não é exercitada por Rank diretamente neste runner; sua validação é separada. Não alteramos o contrato nem fabricamos melhora.


## Entrega final0046

O autor A descartou a diversificação após o resultado real: nenhum ganho nos cortes aceitos, sem base para promovê-la como melhoria. O ranking final mantém os hashes de `implementation_search.go` e `functional_retrieval.go` da baseline0045. A mudança entregue0046 fica na projeção e explicação por requisito, cuja validação funcional pertence ao conjunto separado.

Executamos `real/final.json` sobre os mesmos snapshots e budgets. `real/final-comparison.json` verifica igualdade exata dos **nove relatórios crus**, de todas as linhas de métricas por requisito, dos agregados, do adaptador, da varredura e do freeze. Resultado final: recall rotulado **0/6 @3, 1/6 @12**, precisão rotulada **0/13 @3, 1/43 @12**, indícios em **1/3 negativas**. Não houve ganho ou regressão neste corpus. A lacuna de cobertura permanece explicitamente reconhecida.

Para preservar os resultados experimentais sem sobrescrever `candidate.json`, `run.py` recebeu somente a opção de rótulo `final` na lista do argparse. Os bytes exatos do runner usado na baseline e tentativa estão preservados em `real/run-v1.py`; o recibo final confirma essa única alteração e o mesmo SHA do adaptador. A diferença no hash do arquivo runner final é de orquestração, não de métrica.

**Veredito da entrega:** corpus congelado e comparação final concluídos, com preservação observada do ranking0045 neste recorte. Não foi demonstrada melhora de recuperação em projetos reais; a utilidade da projeção por requisito está em comunicar evidência e incerteza, não em alterar o recall. Os controles separados devem julgar sua correção. Todos os recibos anteriores permanecem como histórico da tentativa rejeitada.
