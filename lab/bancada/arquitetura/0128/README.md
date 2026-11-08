# Bancada 0128 — Parsers sem relógio

O corpus autoral de `manifest_time_test.go` cobre 49 tipos do dispatcher. `baseline.json` registra fontes anteriores e digest de parsing, com exclusão apenas de ConsultadoEm. `before.log` comprova a caracterização antes da alteração. `focused-tests.log` registra determinismo, hora explícita/desconhecida e testes dos consumidores e coletores após a alteração.

`lint-compatible-before.log` usa configuração temporária compatível antes da correção; `lint-standard.log` usa diretamente crystalline.toml corrigido. Ambos preservam regras, limiar warning e 46 avisos V12 anteriores. Configuração corrigida não significa gate aprovado.

`go-test.log`, `go-vet.log` e `go-build.log` registram a verificação do módulo raiz. A suíte usa cache/loopback autorizados e assinatura Git desativada somente no processo das fixtures. Gates selados têm módulo próprio e não foram reproduzidos.

O instante zero indica metadado desconhecido. Relógio de telemetria e extrator do Lab continuam fora do recorte. Corpus e testes não comprovam cobertura semântica universal dos ecossistemas. Nenhuma evidência histórica é reescrita.

Os objetos JSON de dependências em package.json, composer.json, box.json, haxelib.json e dub.json agora usam chaves ordenadas dentro de cada grupo. `dependency-order-before.log` registra a falha anterior nos cinco formatos, e `after.log` registra sua correção. Nomes, versões, grupos e exclusões permanecem; a ordem antes indeterminada passa a ser estável.
