# Passo 0035 — Precisão do lint e análise C/C++ com Clang

**Estado:** IMPLEMENTADO — análise estática limitada às unidades/configurações registradas
**Origem:** pedido do usuário de 2026-09-07: escrever o passo e implementar a correção dos 195 avisos e a integração semântica inicial proposta.
**Baseline:** worktree contendo as correções de evidência da auditoria de 2026-09-07. Não alterar ou reatestar certificados históricos.

## Obrigação e limites

Reduzir avisos falsos sem desativar regras de gravidade/linhagem e obter evidência semântica de unidades C/C++ usando Clang com configuração explícita. Uma AST válida comprova resolução estática no contexto de compilação registrado; não comprova adequação à pergunta, execução, saúde arquitetural ou solução exata.

Regime A/B: testes do contrato semântico escritos a partir desta especificação por executor separado, sem leitura do patch candidato. Execução sem atestação de isolamento de capacidades. Sem selo de refinamento ou extensão de certificados históricos.

## Contrato observável inicial

- L1 define `SemanticRequest{Root string, Query string, Paths []string}` e `SemanticReport{Status string, Reason string, Analyzer string, Version string, Units []SemanticUnit}`. Cada unidade registra Path, Status, Diagnostics, Arguments, Symbols e Calls, mais hashes dos arquivos e configuração examinados.
- L3 expõe `clang.NewClient(binary string)` e método `Analyze(context.Context, domain.SemanticRequest) (domain.SemanticReport,error)`.
- Localizar `compile_commands.json` na raiz ou em `build/`; usar `arguments` quando disponível. Entradas apenas com `command` não são executadas via shell e ficam não avaliadas nesta primeira versão.
- Configuração ausente, incompleta, opção não suportada, include ausente, timeout ou orçamento excedido nunca viram análise completa. Estados: `complete` (somente unidades selecionadas), `partial`, `not_evaluated`. Cobertura global permanece limitada.
- Registrar executável/versão, argumentos efetivos, configuração e hashes de arquivos lidos; distinguir declarações e definições; registrar chamadas com destino resolvido e intervalos do código. Não promover busca a solução.
- Não executar comando arbitrário, compilador indicado pelo repositório, plugin, resposta `@file`, build ou teste do projeto. Invocar somente o Clang configurado com flags permitidas e limites de tempo/saída/unidades.
- Integrar no funil local por uma porta L1 e apresentar resultados no Markdown; chamadas remotas ou sem provedor explícito mantêm `not_evaluated`.
- Identificar símbolos pertinentes à consulta e suas chamadas, conservando trechos/arquivos referenciados na evidência exibida.

## Lint

Inventário inicial: 186 V12 (136 em `_test.go`, 50 em outros arquivos) e 9 V7. Corrigir a identificação de testes, funções e composição pelo linteiro; preservar detecção de tipos/lógica efetivamente indevidos. Arquivos de gates/experimentos devem ser classificados de forma explícita, sem exclusão ampla de produção. Resolver prompts de oráculos excluídos por vínculo/exceção específica documentada. Publicar patch, testes discriminatórios e comando de verificação reproduzível quando a ferramenta estiver fora do repositório.

## Aceite

1. C/C++ pequeno com include e macro configurados: definição e chamada resolvidas com caminho/linha.
2. Include ausente e compilação sem contexto: estado não completo com diagnóstico.
3. Mudança em fonte/header/configuração altera hash; versão e argumentos presentes.
4. `command`, plugin e resposta arbitrária não executam código do repositório.
5. Cancelamento e limites retornam estado não completo.
6. USalign/BinDiff não viram solução por pontuação; relatório conserva distinção semântica estática/comportamento.
7. Validar no checkout Coin disponível, registrando cobertura e impedimentos reais.
8. Lint sem avisos injustificados; regressões positivas e negativas da regra; `go test ./...`, build do pacote `./04_wiring` e `git diff --check`.

## Entrega

Registrar implementação, comandos, resultados e limites em `lab/bancada/history/benchmarks/0035/`. Não declarar cobertura de todo o Coin a partir de unidades selecionadas.

## Resultado

195 → 0 avisos com correção da ferramenta e classificação explícita dos harnesses. Clang integrado à CLI/funil local; testes A/B passaram sem alteração. Coin validado em cópia dos fontes reais com configuração existente: 12/898 unidades, status parcial. Evidências, limitações e reprodução em [`../../lab/bancada/history/benchmarks/0035/README.md`](../../lab/bancada/history/benchmarks/0035/README.md). V12 estrutural não certifica pureza de corpos de funções; Clang não certifica comportamento gráfico.
