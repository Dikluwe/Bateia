# Execução 0046 — Projetos reais e lacunas por requisito

2026-09-07, horário local. [Passo congelado](../passos/0046-validar-recuperacao-em-projetos-reais.md). Implementado e executado. **O ensaio real não demonstrou melhora de recuperação; a entrega acrescenta avaliação reproduzível e comunicação de evidência por requisito.**

## O que mudou no produto

`analyze` e `idea recommend` agora projetam indícios por requisito somente depois de revalidar os hashes dos arquivos extraídos. Para cada linha da consulta, o JSON registra:

- `evidence_located`: declaração associada por termos lexicais, sem comprovar o requisito;
- `not_verified`: nenhuma fonte lexical confirmada naquele escopo, sem concluir inexistência;
- termos associados e termos sem correspondência, com fontes por arquivo, símbolo, linha e SHA-256.

Chamadas diretas sem correspondência lexical própria não viram fontes independentes do requisito. Hash, caminho ou linha inválidos impedem fonte; projeções antigas são substituídas após a releitura. Quando arquivos mudam ou não são extraídos, a projeção perde essas fontes. `Reason` identifica a inferência e os limites; não há promoção de score ou capacidade.

Markdown apresenta no máximo oito requisitos, três fontes por requisito e doze termos por lista, explicitando os totais e remetendo ao JSON. A apresentação aparece também em `decision.md`, mantendo decisão canônica separada. Versão de contrato `evidence-v8-requirement-retrieval`; checkpoint local `implementation_search: 0046`. Ranking continua `go_ast_functional_v2`.

## Resultado nos projetos reais

Autor B escolheu snapshots de módulos Go reais disponíveis no cache, com commit de origem, ZIP/checksum, licenças e hashes por arquivo: **bbolt v1.4.0, Snappy v1.0.0 e OAuth2 v0.36.0**. São domínios de armazenamento, compressão e autorização. Corpus de conveniência com seis perguntas funcionais em português e três negativas, congelado antes do baseline. Não representa todo o ecossistema.

| Medida do ranking final | Resultado |
| --- | --- |
| Arquivos-alvo rotulados recuperados nos primeiros 3 caminhos distintos | 0/6 |
| Arquivos-alvo rotulados recuperados nos primeiros 12 caminhos distintos | 1/6 |
| Precisão em relação aos arquivos rotulados, corte 12 | 1/43 |
| Negativas com algum indício lexical | 1/3 |

Precisão rotulada não significa que os outros 42 caminhos foram julgados irrelevantes. Métrica por arquivo também não demonstra recuperação do símbolo esperado. O runner aplica orçamento de 24 hits e simula cortes de caminhos; não equivale à execução da coleta. Manifestos e [avaliação B detalhada](../../lab/bancada/history/benchmarks/0046/evaluation-real.md) documentam essas distinções, fontes e perdas por requisito.

A varredura não esgotou orçamento de leitura neste corpus. Houve perdas por formulação sem correspondência útil, termos genéricos compartilhados e corte de resultados. A pergunta de descompressão Snappy, por exemplo, não retornou definição; em OAuth2, a negativa sobre compressão retornou indícios pelo termo HTTP.

## Experimento descartado

Testamos distribuir hits em rodadas de arquivos, para impedir que vários símbolos do mesmo caminho consumissem o orçamento. O arquivo de leitura do bbolt passou a aparecer, mas em 14º caminho, fora do corte 12. O recall aceito não melhorou; os retornos não rotulados aumentaram. A alteração foi retirada.

O ranking final tem hashes idênticos ao 0045 nos dois arquivos centrais. O ensaio final reproduziu exatamente os nove relatórios do baseline. A tentativa descartada e seu teste exploratório foram preservados como histórico, fora da implementação e da suíte ativa. Não alteramos rótulos, glossário ou pesos para aparentar ganho. A nova projeção explica a lacuna; não corrige o recall.

## CLI e replay final

Executamos também as **nove consultas pela CLI completa**, nos diretórios dos mesmos módulos. Todos os hashes dos arquivos Go varridos coincidiram com o manifesto congelado. Como são distribuições de módulos sem checkout Git, a CLI não fornece commit-base Git desses diretórios; o vínculo de versão/commit da origem vem do manifesto, separado dos hashes efetivamente lidos. Não se presume que o módulo seja uma árvore Git completa.

Na negativa OAuth2, a projeção apresenta apenas `http` entre os termos associados e deixa os termos de compressão sem correspondência. Continua sendo inferência lexical, nunca compressão demonstrada. Snappy/descompressão fica `not_verified`. Esse é o ganho entregue na comunicação, sem esconder a recuperação insuficiente.

No controle multilinha `SanitizeAssessmentEvidence` + `quasar teleportation`, o primeiro requisito tem fontes lexicais conferidas e o segundo não tem fonte. A recomendação da ideia do próprio Bateia mantém dez capacidades INCONCLUSIVE. Replay preservou byte a byte assessment, decisão, alternativas, bundles e `decision.md`.

[Recibos, comandos, identidades e projeções da CLI](../../lab/bancada/history/benchmarks/0046/cli-observations.json), [controle legível](../../lab/bancada/history/benchmarks/0046/requirements.md) e [decisão com requisitos](../../lab/bancada/history/benchmarks/0046/decision.md). Tempos das execuções são telemetria da ferramenta, não benchmarks dos projetos. Nenhum código externo foi executado; fontes foram lidas para análise.

## Validação e Tekt

- Suíte completa `GOCACHE=/tmp/bateia-0046-cache go test ./...`: passou na versão final; testes HTTP locais exigiram execução com permissão correspondente.
- `crystalline-lint .`: `✓ No violations found`.
- `GOCACHE=/tmp/bateia-0046-cache go build -o /tmp/bateia-0046 ./04_wiring`: passou.
- `git diff --check`: passou.
- Testes independentes da projeção: passaram, inclusive releitura alterada/ausente e não reutilização de fontes obsoletas. [Recibo final](../../lab/bancada/history/benchmarks/0046/tests/execution-final.json).
- Regressão 0045: corpus/testes originais preservados; recall e precisão sintéticos continuam 100%, sem generalizar para o corpus real. [Avaliação da regressão](../../lab/bancada/history/benchmarks/0046/regression0045/final/evaluation.md).

Regime **A/B executado sem atestação física de isolamento**. Autor de corpus e autor de testes não leram a implementação candidata para formular os casos. Coordenador escreveu contrato e produção e recebeu resultados públicos para decidir o experimento. O corpus real foi congelado antes do baseline; testes exploratórios da diversificação foram congelados após a implementação, com essa limitação registrada. Não há novo selo de refinamento nem alteração de gates/certificados existentes.

## Limites que permanecem

A recuperação funcional em português ainda é fraca nesta amostra real. Vocabulário, termos genéricos e a relação entre pergunta e mecanismo precisam de investigação própria. A seleção continua Go/lexical/AST, sem execução ou equivalência semântica. Os limites anteriores de gRPC/Protobuf, Bleve e linguagens não cobertas não foram corrigidos neste passo. A utilidade com usuários e a qualidade global permanecem não avaliadas.
