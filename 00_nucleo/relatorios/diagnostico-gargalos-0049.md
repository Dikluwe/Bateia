# Diagnóstico 0049 — onde a qualidade se perde

2026-09-08. Investigação do integrador, sem alteração de produção. A correção 0048 tratou uma ambiguidade real, mas não demonstrou aumento de recuperação. Foi inadequado priorizá-la como resposta ao problema geral de qualidade.

## Resultado principal

Os casos examinados apontam para três gargalos diferentes: a busca remota não recebe a recuperação AST recém-implementada; a recuperação local depende de coincidências lexicais sem uma ponte suficiente entre intenção e implementação; a avaliação de capacidades exige regras de um catálogo que não cobre a maior parte das perguntas da autoanálise. Não há evidência para atribuir todos os problemas do Bateia a um único defeito.

## 1. As melhorias locais não chegam à descoberta remota

**Fato observado no código:** `03_infra/content/router.go`, função `SearchImplementations`, chama o buscador apenas quando a URL começa com `file://`. Para os demais candidatos retorna `not_evaluated`, sem arquivos examinados. A enumeração dos adaptadores confirma que a implementação de busca está em `03_infra/local`.

**Fato histórico:** a busca pública `BM25 code search` do passo 0047 retornou dois candidatos pelo fallback; ambos tinham busca de implementação `not_evaluated`. O primeiro extraiu Makefile, README e go.mod. O [recibo histórico](../../lab/bancada/history/benchmarks/0047/ideas/public-search-observations.json) registra isso.

**Inferência:** ajustar IDF e interpretação no caminho local não melhora esse trecho da jornada remota. O achado explica por que podemos melhorar benchmarks locais e continuar recebendo documentação pouco útil ao procurar projetos externos. Não mede a qualidade atual dos provedores: a consulta remota não foi repetida neste diagnóstico.

## 2. O índice descarta a ponte entre linguagem humana e código

**Fato observado no código:** `implementation_search.go` calcula termos do nome e de `functionalBodyTerms`; esta função visita identificadores AST do corpo, sem comentários e sem a assinatura como campo próprio. O parâmetro `doc` chega à função de coleta mas não é usado. Existe um glossário PT/EN finito; não cobre as formulações de compressão/descompressão examinadas.

**Experimento executado:** mesmo binário 0048 e mesmos arquivos Snappy v1.0.0, variando apenas a consulta. [Comandos, arquivos lidos, hashes, termos e hits](../../lab/bancada/history/benchmarks/0049-diagnostico/comparisons.json).

| Consulta | Observação |
| --- | --- |
| `Descomprimir um bloco e rejeitar dados corrompidos` | Zero hits; `decode.go` foi lido; orçamento não esgotado. |
| `decompress a block and reject corrupt data` | Localiza `decode.go`, mas o primeiro símbolo é `fill`. |
| `decode block corrupt` | `Decode` em primeiro lugar. |
| `Decode` | Definição exata em primeiro lugar. |
| `Comprimir bytes em um bloco e reutilizar o buffer de destino` | Hit lexical em `encodeBlock` por `bytes`; não recupera a API `Encode`. |
| `compress bytes into a block and reuse destination buffer` | `Encode` aparece entre os primeiros cinco hits. |
| `Encode` | Definição exata em primeiro lugar. |
| `decompresses` | Zero hits, embora a palavra esteja na documentação de `NewReader`, em `decode.go`. |

**Inferência sustentada:** a falha de descompressão em português ocorre depois da leitura e antes da seleção útil de símbolos. A mudança de formulação recupera a implementação sem aumentar orçamento ou alterar código. Não é evidência de que tradução sozinha resolva tudo: o resultado em inglês ainda favorece outra função. A documentação vinculada às declarações é um canal de recuperação ausente; deve localizar candidatos, não provar comportamento.

## 3. Ranking não representa cumprimento de requisitos

**Fato observado e reproduzido:** adicionar o termo inexistente `quasarteleport` a `decode block corrupt` preserva integralmente os hits, scores e pesos. `corpus_scoring.go` soma apenas termos presentes; a projeção registra termos faltantes separadamente. Esse comportamento é compatível com recuperação lexical, mas não com interpretar a nota como adequação completa ao pedido.

`requirement_retrieval.go` usa `evidence_located` quando existe qualquer fonte com interseção lexical. Portanto o rótulo significa material localizado, mesmo quando o verbo principal da pergunta não tem correspondência. O JSON avisa sobre isso; a qualidade da recuperação não aumenta por avisar melhor.

**Limite:** o termo artificial é um controle de pontuação, não um requisito funcional real. Este ensaio não demonstra que penalizar toda palavra ausente seja uma correção adequada: sinônimos e palavras de ligação também podem estar ausentes. É necessário distinguir conceito principal, objeto e condições antes de avaliar pertinência conjunta.

## 4. Encontrar mais código não desbloqueia o catálogo de avaliação

**Fato observado:** nove das dez capacidades na recomendação local 0048 têm a explicação `nenhuma regra versionada cobre esta capacidade`. A décima encontra indícios relacionados a offline, mas também permanece inconclusiva. [Causas preservadas](../../lab/bancada/history/benchmarks/0049-diagnostico/assessment-causes.json).

**Causa no código:** em `01_core/usecase/evidence.go`, `findRule` escolhe a primeira regra com correspondência por substring nos termos da capacidade. `assessCapability` retorna imediatamente `INCONCLUSIVE` quando não há regra, antes de examinar os itens coletados. Mais arquivos ou melhor ranking não mudam essa condição. Essa contenção evita inventar uma prova; o problema é a cobertura limitada do avaliador, não o uso de `INCONCLUSIVE` em si.

## Direção de correção

Priorizar uma jornada verificável: pergunta funcional → candidato externo → snapshot imutável → símbolos pertinentes → indícios vinculados ao requisito. A primeira obrigação é fazer a recuperação de implementação alcançar o candidato remoto, com coleta limitada e rastreável, e medir perdas por estágio. Em paralelo conceitual, o índice deve aproveitar documentação vinculada à declaração e registrar expansões de vocabulário como hipóteses, mantendo os termos originais.

Depois, avaliar se os indícios cobrem ação, objeto e condições em vez de usar quantidade de palavras como substituto. Separar capacidade sem avaliador de coleta insuficiente; somente verificações delimitadas podem sustentar conclusões mais fortes. Não promover um score lexical a prova para esconder a lacuna do catálogo.

O próximo contrato de implementação deve exigir melhora em perguntas funcionais novas e nos artefatos da CLI completa, além de controles negativos e equivalência de recuperação local/remota para o mesmo snapshot. Preservar métricas antigas é regressão, não demonstração de ganho. Este diagnóstico não altera os testes congelados nem amplia glossários para memorizar Snappy.

[Identidades e limites da investigação](../../lab/bancada/history/benchmarks/0049-diagnostico/receipt.json). Os testes exploratórios são do integrador, sem alegação de avaliação independente. Código externo foi somente lido. Tempos são telemetria do Bateia.
