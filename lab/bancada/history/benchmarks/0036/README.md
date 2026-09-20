# Validação do passo 0036

Implementado em 2026-09-07. Coin permanece apenas corpus externo de análise; não foi incorporado como dependência ou fonte do Bateia.

## Resultado observado
- Pipeline completo da CLI: 16 arquivos extraídos, incluindo `src/actions/SoGLRenderAction.cpp`; os headers XML irrelevantes da amostra anterior deixaram de preencher a cota.
- Markdown original: 38.140 bytes contra 476.063 no passo 0035 (redução de aproximadamente 92%). Volume de documento, não benchmark do Coin.
- As seis definições centrais identificadas na auditoria aparecem no relatório. A cadeia `doSortedLayersBlendRendering` → inicialização / camada / composição FP e NV foi conferida.
- 24 definições e 64 chamadas na projeção, com limites explícitos. Trechos e hashes conferidos contra o conteúdo integral extraído.
- Score de descoberta permanece 16; similaridade geral fica não avaliada (`overall_score=-1`, `evaluation_status=insufficient`). Dimensões ausentes também usam -1, separadas do sinal lexical conhecido. Saúde continua não avaliada.
- Clang: 12/898 unidades; cobertura parcial. Enumeração com enumeradores é reconhecida como definição. Nenhuma execução gráfica foi realizada.
- `go test ./...`, build de `./04_wiring`, `git diff --check` e lint sem violações passaram.

## Artefatos e reprodução
- `coin-analysis.md`: apresentação gerada; apenas link de recibo atualizado após arquivamento.
- `coin-analysis.json.gz`: saída JSON integral, preservando contextos por unidade, argumentos, hashes e bytes extraídos.
- `validation.json`: comando, hashes do executável e fontes do analisador, métricas e resultados.
- `go-test.log`, `lint.log`, `coin-cli.log`: resultados de execução. Tempos são telemetria do Bateia.
- `check_projection.py`: executar com os caminhos Markdown e JSON/gzip; confere as seis definições, chamadas, ausência de sobreafirmação, tamanho e correspondência dos trechos.

O corpus e a configuração são os mesmos de `../0035/coin-inputs.json`, preparados por `../0035/prepare-coin-snapshot.py`. O snapshot não contém `.git`; o commit-base externo permanece no recibo do corpus. Os hashes do relatório identificam os bytes efetivamente examinados.

## Limites
Melhoria demonstrada neste caso e nos testes sintéticos de outro domínio; não uma avaliação universal ou independente. A seleção continua lexical e por chamadas estáticas; inclui funções auxiliares e pode perder relações indiretas ou além dos dois saltos. Chamadas estáticas não provam execução conjunta ou correção visual. A AST completa mantém duplicações de headers por contexto de compilação; a projeção deduplica identidades e chamadas, não elimina a informação original. O JSON lateral automático foi integrado à saída Markdown da CLI; outros consumidores recebem a projeção e devem solicitar exportação JSON para recibos integrais.
