# Gate de desenvolvimento do Bateia

`make verify` é a entrada canônica de verificação: executa lint, auditorias estruturais/de ownership, suíte raiz, vet, build e diff-check, registrando comandos/códigos de saída em diretório temporário. `BATEIA_VERIFY_OUTPUT` seleciona o destino dos logs e do binário. `make lint` executa apenas a ferramenta arquitetural fixada.

A fonte upstream do linter está congelada em `crystalline-lint/source.tar.gz`, com revisão, origem, licença declarada pelo upstream e SHA-256 em `manifest.json`. `go-parser.patch` corrige classificação Go sem mudar regras/severidades. O runner valida hashes, compila/testa em cache temporário fora do produto e executa 11 controles positivos/negativos; os controles são repetidos também ao reutilizar o binário validado. Nenhum clone externo ou instalação global é necessário.

Pré-requisitos: Go conforme go.mod, Python com tarfile data filter, Git, Make, Rust/Cargo e crates da revisão fixada disponíveis no cache Cargo para a compilação offline inicial. As dependências Rust estão no Cargo.lock do arquivo de fonte. A ausência de uma dependência falha explicitamente. O produto compilado não depende dessas ferramentas nem importa o Lab. `BATEIA_TOOL_CACHE` pode selecionar o cache do linter.

O runner da pilha usa `crystalline-lint --fail-on warning .` com o binário fixado no PATH do processo. A assinatura de commits de fixture é desativada somente no processo, acrescentando a opção à configuração Git de ambiente existente. Testes httptest precisam de loopback. `GOCACHE` pode apontar para um diretório autorizado. Gates selados têm módulo e reprodução próprios; a suíte raiz não os promove.

Os checks de ownership cobrem os pontos da revisão 0131 e possuem controles de regressão. Eles não constituem análise semântica universal de arquitetura ou avaliação da utilidade das conclusões do produto. Evoluir as regras exige preservar controles negativos e declarar o novo contrato.
