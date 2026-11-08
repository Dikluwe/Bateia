# Gate de desenvolvimento do Bateia

`make verify` é a entrada canônica de verificação: executa lint, auditorias estruturais/de ownership, suíte raiz, vet, build e diff-check, registrando comandos/códigos de saída em diretório temporário. `BATEIA_VERIFY_OUTPUT` seleciona o destino dos logs e do binário. `make lint` executa apenas a ferramenta arquitetural fixada.

A fonte upstream do linter está congelada no arquivo indicado por `crystalline-lint/manifest.json`, que registra revisão, origem, licença declarada pelo upstream e SHA-256. As versões anteriores permanecem preservadas conforme a seção de seleção efetiva abaixo. `go-parser.patch` corrige classificação Go sem mudar regras/severidades. O runner valida hashes, compila/testa em cache temporário fora do produto e executa 11 controles positivos/negativos; os controles são repetidos também ao reutilizar o binário validado. Nenhum clone externo ou instalação global é necessário.

Pré-requisitos: Go conforme go.mod, Python com tarfile data filter, Git, Make, Rust/Cargo e crates da revisão fixada disponíveis no cache Cargo para a compilação offline inicial. As dependências Rust estão no Cargo.lock do arquivo de fonte. A ausência de uma dependência falha explicitamente. O produto compilado não depende dessas ferramentas nem importa o Lab. `BATEIA_TOOL_CACHE` pode selecionar o cache do linter.

O runner da pilha usa `crystalline-lint --fail-on warning .` com o binário fixado no PATH do processo. A assinatura de commits de fixture é desativada somente no processo, acrescentando a opção à configuração Git de ambiente existente. Testes httptest precisam de loopback. `GOCACHE` pode apontar para um diretório autorizado. Gates selados têm módulo e reprodução próprios; a suíte raiz não os promove.

Os checks de ownership cobrem os pontos da revisão 0131 e possuem controles de regressão. Eles não constituem análise semântica universal de arquitetura ou avaliação da utilidade das conclusões do produto. Evoluir as regras exige preservar controles negativos e declarar o novo contrato.

## Seleção efetiva de regras — 0159

O gate ativo fixa a revisão 20eeb9b do tekt-linter, incluindo a configuração enabled por regra introduzida em 6c4bf23. O manifest ativo aponta para source-20eeb9b.tar.gz; source.tar.gz e manifest-e53de6f.json preservam a fonte/hash anteriores. O patch Go permanece idêntico. Nenhum checkout externo ou binário global governa make lint.

crystalline.toml declara V0–V27 com enabled=true, sem redução de severidades ou flags no-drift/no-stale. A saída informa selectedRules, disabledRules e analysisComplete. O gate repete os 11 controles anteriores do parser e dois controles reais de V3: com enabled=true a importação proibida é diagnosticada e falha; com enabled=false ela é omitida, inclusive com --checks all. Os controles isolados não desligam nenhuma regra no Bateia. Regras selecionadas podem não ter aplicação a todas as linguagens; alcance sintático não comprova comportamento semântico.

Ao mudar a seleção por projeto, revisar a política e seu alcance junto ao código. Não usar fix-hashes ou update-snapshot para fazer o gate passar: registros/recibos históricos são preservados, e reparos de linhagem exigem seu contrato próprio.
