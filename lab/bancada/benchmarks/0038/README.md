# Validação do passo 0038 — AST + CFG + pré-processamento

## Entrega
Auxiliar C++ LibTooling com contrato `bateia-flow-v1`, integrado de forma opcional ao provedor Clang. Exporta definições/localizações AST, blocos CFG, sucessores, terminadores, chamadas diretas ou não resolvidas, localização de expansão/origem inicial de macros e intervalos ignorados pelo pré-processador. O núcleo Go permanece sem I/O.

O auxiliar deve ser explicitamente configurado com `BATEIA_CLANG_FLOW_TOOL`. Instruções de compilação em `../../../03_infra/clang/flowtool/README.md` (a partir da raiz do projeto, consultar `03_infra/clang/flowtool/README.md`). Não foi instalado globalmente; binário validado em `/tmp/bateia-clang-flow-final`. Dependências LLVM/Clang 18 já estavam disponíveis. Usar compilador e bibliotecas auxiliares compatíveis; outras versões não foram validadas.

## Evidência obtida
- Suíte Go completa aprovada com `BATEIA_TEST_FLOW_TOOL=/tmp/bateia-clang-flow-final`.
- Clang real: if/else, for/while, break/continue, curto-circuito, retorno antecipado, chamada indireta, macros e #if 0 dentro de função ativa.
- Conferência de arestas: o bloco com fallback e retorno não alcança o bloco com done no caso sintético.
- Seleção lexical rejeita termos no intervalo ignorado de função ativa.
- Auxiliar ausente resulta em fluxo não avaliado; unidade inválida não exibe CFG válido.
- Coin: 12/898 unidades analisadas. Unidade do renderizador contém CFGs de 71 funções e um intervalo ignorado. Laço de passagens e chamadas de inicialização/camada/composição conferidos no JSON.
- Markdown original: 52.564 bytes. Arestas são fatos estruturais; não são prova de execução ou viabilidade.
- Build Go e compilação C++ aprovados, lint sem violações, diff sem erros de whitespace.

## Artefatos
`coin-analysis.md`, `coin-analysis.json.gz`, `validation.json`, `go-test.log`, `flow-tests.log`, `lint.log`, `coin-cli.log` e `check_flow.py`. O script recebe JSON/gzip e Markdown, verifica o grafo e as chamadas esperadas. JSON preserva argumentos efetivos, hash do auxiliar e hashes das entradas. Coin continua corpus externo, sem alteração no projeto original; recibo de preparação em `../0035/coin-inputs.json`.

## Limites explícitos
Protótipo intraprocedural restrito ao arquivo principal de cada unidade, CFG sem poda de arestas trivialmente falsas. Sem análise de dominância, prova de viabilidade, interpretação geral do algoritmo, resolução de ponteiros ou execução gráfica. Sucessores não recebem rótulos de verdadeiro/falso por suposição de ordem. A apresentação limita blocos; o JSON mantém o grafo integral. Origens de macros registram spelling e expansão, sem histórico completo de expansões aninhadas. Hash do executável não substitui reprodução do ambiente de bibliotecas dinâmicas. Não é auditoria independente.
