# Bancada de organização 0127

Caracterização autoral antes/depois da extração da resolução de fontes de `RecommendationRunService` para porta de leitura em L1. Contrato e síntese em `00_nucleo/passos/0127-separar-resolucao-de-fontes-da-sintese.md` e `00_nucleo/laudos/0127-resolucao-de-fontes-separada.md`.

- `baseline.json`: commit, hash do proprietário antes da extração e digests congelados de sucesso/degradação.
- `before.log`: mesmos cenários no resolvedor original, antes da extração.
- `after.log`: cenários e digests preservados depois da extração.
- `go-test.log`, `go-vet.log`, `go-build.log`: verificação do módulo raiz.
- `lint-standard.log`: resultado do comando exigido.
- `lint-compatible-before.log` e `lint-compatible-after.log`: diagnóstico com cópia temporária da configuração, substituindo a declaração Go booleana pela tabela aceita pelo lint instalado. Não modifica `crystalline.toml` do projeto. O diagnóstico anterior inclui o aviso transitório de prompt ainda sem materialização; os demais avisos já existiam.

Os testes usam fontes em memória, sem I/O ou credenciais. Cobrem identidade/hash, paginação, propagação de erros e a projeção de dois cenários. Não comprovam toda a arquitetura nem disponibilidade de fontes externas. Gates selados têm módulo separado e não foram reproduzidos.
