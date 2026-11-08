# 0132 — Diagnóstico de retomada de evidências

**Data:** 2026-10-28. **Contrato:** [passo 0132](../passos/0132-explicar-retomada-de-evidencias.md). **Evidência:** [bancada](../../lab/bancada/confiabilidade/0132/README.md).

## Resultado

`idea recommend --resume` passa a explicar, ao concluir, se reutilizou a evidência anterior e o primeiro impedimento observado em caso de recusa. A avaliação tipada pertence a L1; a CLI somente apresenta estado, código e campo/artefato lógico em português. A função booleana anterior delega à mesma avaliação, preservando a política de elegibilidade, ordem das leituras, short-circuit e duplicatas.

O diagnóstico distingue estado não elegível, schema incompatível, mudança de entrada, catálogo, modelo opt-in, versão de regra, leitor indisponível e artefato ilegível, com hash divergente ou sem registro com hash. A jornada identifica separadamente manifesto que não pôde ser lido/interpretado. O retorno em memória conserva o diagnóstico se uma etapa posterior falhar, sem anunciar conclusão nem promover o resultado.

Nenhum valor de entrada, versão recebida, hash, caminho físico ou erro interno entra no diagnóstico. A apresentação usa apenas códigos e campos conhecidos; valores desconhecidos são omitidos ou recebem mensagem genérica. Sem `--resume`, a saída permanece como antes. A linha original de conclusão e a ordem de publicação/cleanup foram preservadas.

Não há novo artefato, campo de manifesto, regra de aceitação, input causal ou versão. RunID, preimages e formatos mantêm os contratos anteriores. Mudanças apenas de política/decisão ou alvo continuam permitindo reutilizar evidência compatível, com nova avaliação da decisão. Replay que recusa o checkpoint do destino pode continuar pela sua fonte congelada.

## Verificação

Os testes puros cobrem motivos, precedência, todos os campos de versão, imutabilidade e compatibilidade de leituras com o predicado legado. Os testes da jornada exercitam retorno, anúncio, rejeições, manifesto indisponível e falha posterior. Os testes da CLI verificam mensagens/privacidade, replay offline e comportamento efetivo de descoberta/coleta: mudança só de política reaproveita evidências; hash adulterado causa nova coleta. Os oráculos anteriores não foram alterados.

Consultar comandos e resultados da pilha em [final/verification.json](../../lab/bancada/confiabilidade/0132/final/verification.json). A primeira execução restrita bloqueou duas fixtures HTTP locais; elas passaram com loopback autorizado, preservando os testes. Lint e auditorias contemplam 354 arquivos Go produtivos. O gate do módulo raiz mantém seu escopo; gates selados não foram promovidos.

`make verify` passou integralmente na execução final: lint sem violações, auditorias arquiteturais, `go test ./...`, `go vet ./...`, build de `./04_wiring` e `git diff --check`. A linhagem registra sete materializações com prompts proprietários e os oráculos anteriores preservados. Fonte, patch e controles negativos do linter continuam os fixados no projeto.

## Limites

A mensagem relata apenas o primeiro impedimento, após publicação bem-sucedida. A elegibilidade registra os critérios aplicados no momento da checagem; não autentica origem, evita toda alteração concorrente ou aumenta a força da evidência. Não houve avaliação de utilidade com usuários nem consultas reais aos provedores.

O ensaio de comparação encontrou uma oscilação legada entre `null` e `[]` em `requirement_coverage` vazio. A fonte de serialização permaneceu intacta; a observação e a tentativa estão na bancada. A comparação final de bytes usa cobertura não vazia e não certifica idempotência geral dessa serialização. A correção desse formato exige contrato próprio para tratar sua relação com os hashes de artefatos.
