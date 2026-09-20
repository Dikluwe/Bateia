# Execução 0052 — Expansão rastreável e admissão por ação

2026-09-08. Implementados o [contrato inicial](../passos/0052-expandir-vocabulario-com-proveniencia.md) e o [adendo r2](../passos/0052-adendo-admissao-por-acao.md). A pergunta portuguesa de descompressão agora localiza `Decode` no Snappy. A primeira candidata ampliava resultados irrelevantes e não foi promovida; a revisão exige ação e contexto para admitir alternativas ambíguas.

## Implementação

Um catálogo finito `0052-v1` acrescenta traduções e alternativas técnicas a tokens naturais. Preserva a consulta original e registra `Source`, `Terms`, `RuleID`, `Kind` e `Status: hypothesis` em `Expansions`. Não expande identificadores explícitos nem linhas com marcadores de restrição, e não encadeia termos gerados como novos gatilhos. O glossário de código e a fórmula IDF não foram alterados.

O adendo filtra cada linha da consulta por suas famílias de ação. Termos secundários, como `data`, não bastam para satisfazer uma ação. Alternativas `encode`/`decode` de compressão/descompressão exigem contexto solicitado (`block`, `blocks`, `corrupt` ou `corrupted`) no mesmo campo. Código usa nome/corpo; documentação usa seu grupo Doc, sem misturar os canais. Em consultas multilinha, os termos de uma linha rejeitada não entram na pontuação daquele candidato nem em suas fontes projetadas. Símbolo exato conserva prioridade.

Esse filtro verifica coocorrência lexical, não fluxo, condições ou equivalência. `decode` não significa necessariamente descompressão. Markdown mostra os recibos como hipóteses, com limites de exibição; JSON preserva tudo. A revalidação reconstrói a interpretação a partir da consulta original externa. Versão `evidence-v13-query-expansions`, checkpoint local/live `0052`. Não houve alteração das consultas enviadas aos provedores, do catálogo de capacidades ou dos limites de coleta.

## Experimento inicial e correção

A candidata r1 colocou `Decode` e `Encode` em primeiro no Snappy, mas fez a mesma pergunta de descompressão retornar doze hits no UUID. Termos secundários e alternativas técnicas amplas admitiam código sem relação com compressão. [Resultado r1 preservado](../../lab/bancada/history/benchmarks/0052/cli-comparison-r1.json), junto aos hashes e arquivos da primeira candidata.

O adendo foi escrito e congelado depois desse contraexemplo público, antes dos novos testes e execução r2. A revisão restringiu a admissão, mantendo o catálogo de pares. Não removemos a consulta negativa nem alteramos suas palavras para aparentar melhora.

## Resultado final pela CLI

Mesmos módulos e perguntas, comparados com o binário 0051:

| Caso | Baseline | R1 irrestrita | R2 final |
| --- | --- | --- | --- |
| Snappy: `Descomprimir um bloco e rejeitar dados corrompidos` | 0 hits | 24 | 8; `Decode` primeiro |
| Snappy: `Comprimir bytes em um bloco e reutilizar o buffer de destino` | 9 hits | 24 | 11; `Encode` primeiro |
| UUID: mesma pergunta de descompressão | 0 | 12 | 0 |
| UUID: `comparar dados` | 0 | 11 | 0 |
| Snappy: `desserializar transação SQL` | 0 | 0 | 0 |
| Snappy: `descomprimir sem alocar` | 0 | 0 | 0; sem novas expansões |
| UUID: `serializar dados` | 0 | não executado | 6; `MarshalJSON` primeiro |

[Recibos r2 com consultas originais, expansões, contribuições, intervalos e hashes](../../lab/bancada/history/benchmarks/0052/cli-comparison-r2.json), [Markdown Snappy](../../lab/bancada/history/benchmarks/0052/snappy.md). O controle de serialização foi acrescentado como ensaio diagnóstico posterior; não é holdout. O número de hits não é precisão: algumas fontes ainda cobrem apenas parte da pergunta, e não houve julgamento exaustivo de todas as fontes.

**Execução remota:** Snappy no commit `43d5d4cd4e0e3390b0b645d5c3ef1187642403d8`, com a pergunta portuguesa original. Os oito hits coincidiram com os locais e foram confirmados após extração/releitura. `Decode` permaneceu primeiro. [Recibo remoto](../../lab/bancada/history/benchmarks/0052/remote-verification.json) e [runner](../../lab/bancada/history/benchmarks/0052/remote-runner.go.txt). Estado `complete` refere-se à recuperação limitada, não à capacidade demonstrada. A rede exigiu permissão após bloqueio no sandbox. Código externo foi somente lido.

## Verificação e custódia

Build de `./04_wiring`, suíte completa `go test ./...`, `crystalline-lint .` e `git diff --check` passaram. Replay reproduziu os cinco artefatos comparados byte a byte; dez capacidades da autoavaliação permaneceram inconclusivas. [Replay](../../lab/bancada/history/benchmarks/0052/replay.json) e [recibo final](../../lab/bancada/history/benchmarks/0052/verification.json).

Tekt A/B executado sem atestação de isolamento físico. Integrador escreveu contratos, apresentação, versões e ensaios; A implementou interpretação/admissão; B escreveu testes sem ler corpos da candidata. Filesystem compartilhado. Busca de assinaturas expôs nomes de testes anteriores a B, sem corpos; isso consta no veredito.

Dois erros de observabilidade dos testes iniciais foram corrigidos com registros preservados: ordem visual da tabela versus ordem lexicográfica dos termos, e tentativa de reconstruir a consulta por uma API que não recebia a consulta original. O teste B1 que exigia admissão irrestrita de `Decode` sem contexto foi arquivado e retirado da suíte ativa por conflito explícito com o novo adendo. Essa supersessão é semântica e não foi apresentada como simples correção mecânica. Demais esperados B1 foram mantidos.

Os testes r2 foram congelados antes da execução e passaram na primeira tentativa: dezoito casos de ação/contexto/famílias/exceção de símbolo, quatro de separação documento/código, além de multilinha/projeção e controles B1 remanescentes. [Veredito B e histórico](../../lab/bancada/history/benchmarks/0052/tests-verdict.md). Não há selo de refinamento ou teste de mutações que prove equivalência geral.

## Limites

O ganho é demonstrado em um recorte conhecido, não em toda consulta portuguesa. O catálogo é uma hipótese manual, finita e escolhida com conhecimento do caso Snappy; não foi aprendido de forma independente. O ensaio UUID acrescenta outro domínio, mas não substitui avaliação externa de generalização.

A revisão elimina o contraexemplo observado, mas código errado que contenha `decode` e `block` ainda pode ser admitido. Requisitos com múltiplas ações distribuídas por funções podem perder fontes por exigir todas no mesmo campo. Linhas com negação continuam sem novas expansões. Aliases contam separadamente no score; calibração por conceito e análise de ação/objeto/condições permanecem pendentes. Nenhum resultado foi promovido a capacidade comprovada. Tempos nos artefatos são telemetria do Bateia.
