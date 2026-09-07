# Lint — Passo 0035

Implementação executada sem atestação de isolamento; este relatório não é veredito independente nem reatesta certificados históricos. Fonte externa preservada em `/repos/Antigravity/tekt-linter`, commit `49f46885fd03f753e9f1cf37e271d6baf30ed725`; patch produzido e compilado em cópia `/tmp/bateia-0035-linter`. Nenhuma instalação global foi realizada por esta subtarefa.

## Correção

Baseline: 195 avisos, sendo 186 V12 e 9 V7. O parser Go tratava toda função exportada como `Impl` e todo tipo exportado como `Struct`. Isso rotulava testes e funções de composição como declarações de implementação. O patch distingue métodos de funções livres e identifica structs, interfaces e aliases/tipos definidos, inclusive privados. Métodos permanecem sujeitos a V12.

Arquivos cujo nome termina exatamente em `_test.go` recebem isenção somente de V12. As declarações continuam na IR: imports, tokens e linhagem não são apagados. Em `crystalline.toml`, os arquivos não teste dos harnesses históricos `step0028_gate` e `step0029_experiment` são listados individualmente em `wiring_exceptions.harness_files`, com justificativa. Não há glob, prefixo ou exclusão adicional do walker. Paths absolutos, com `..` ou justificativa vazia não autorizam a isenção. Nove exceções exatas de prompts resolvem V7 para oráculos já excluídos na configuração anterior.

A V12 é uma verificação estrutural de declarações, conforme ADR-0007 da ferramenta. Uma função livre Go não equivale a um `impl` Rust. Portanto **a V12 não demonstra ausência de lógica de negócio em funções livres**; não há exceção por nome `New*`, e o mesmo critério vale para funções públicas/privadas. Os três construtores possuem verificações de configuração e recuperação; o lint limpo não certifica sua pureza comportamental. Análise de corpos com atribuição semântica de decisões fica fora deste patch. Interfaces, aliases, structs (quando configurado) e métodos efetivamente indevidos em produção continuam detectados, inclusive privados anteriormente ignorados.

## Evidência

- Três regressões novas falharam antes da implementação (`lint-tests-red.log`).
- `cargo test --offline --manifest-path /tmp/bateia-0035-linter/Cargo.toml`: 1095 testes aprovados em 71 suites (`lint-tests.log`). Inclui cinco regressões locais novas. Uma primeira rodada revelou alteração na identidade do ratchet de complexidade; o parser foi reorganizado para conservar a estrutura já aceita, sem editar o ratchet ou seu baseline.
- `lint-regressions.py`: seis fixtures CLI adicionais verificam V12 em tipos/métodos de produção, correspondência exata de harness e continuidade de V1/V3 (`lint-regressions.log`).
- `cargo build --offline --manifest-path /tmp/bateia-0035-linter/Cargo.toml`: aprovado.
- Binário corrigido sobre o Bateia: `✓ No violations found` (`lint-after.log`).
- Auto-lint da ferramenta: nenhum erro/aviso acionável; os informativos V19/V20 históricos permanecem. O teste do ratchet completo passou, sem alterar baseline.
- `--fix-hashes` executado na cópia da ferramenta, atualizando somente os cinco prompts/códigos alterados. Nenhum certificado Bateia foi alterado.

## Reprodução e entrega

Execute `bash 00_nucleo/benchmark/0035/reproduce-lint.sh` na raiz do Bateia. O script exige a revisão limpa identificada, copia a ferramenta para `/tmp`, verifica/aplica `crystalline-lint.patch`, executa testes/build offline, seis regressões CLI e lint Bateia. Não instala nada globalmente. Requer Rust/Cargo e dependências já em cache; uma instalação sem esse cache precisa provisioná-las separadamente.

Binário validado e cópia estável para integração: `/tmp/bateia-0035-crystalline-lint` (também em `/tmp/bateia-0035-linter/target/debug/crystalline-lint`).
SHA-256: `e47974fe903b225e96fec580040604d80380edda054a7b2187bbf5b18da9a30b`.
SHA-256 do patch: `ce945cc344bffb4c9067ce74b600baaee640121b4fb915176997b1774dd4235e`.

O patch adiciona um campo público em `WiringConfig`; consumidores Rust que constroem esse tipo literalmente precisam informar `harness_files: vec![]` ou usar `..Default::default()`. Os consumidores internos e testes da ferramenta foram atualizados. Isto não constitui promessa de compatibilidade binária nem análise semântica completa de Go.

Reprodução efetivamente executada em segunda cópia com o runner: teste completo, build, regressões CLI e lint passaram (`lint-reproduce.log`). O diretório de cache Cargo foi compartilhado somente para reduzir recompilação; não constitui isolamento de execução.
