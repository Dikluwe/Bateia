# Bancada da composição de descoberta

`startup-checks.json` registra duas execuções do binário compilado com `mcp`: registro padrão de cinco fontes e registro padrão com Bitbucket de seis fontes. As duas usaram ambiente isolado, diretório temporário, cache/home/raiz temporários e nenhuma credencial. Foram enviados apenas tools/list e list_connectors; nenhuma pesquisa remota foi solicitada. Ambas anunciaram onze tools, com enum correspondente ao registro, e encerraram com código zero.

Os testes `04_wiring/market_composition_test.go` verificam ambos os construtores com 1, 3, 5, 6, 10 e 11 fontes, rejeição antes de criar recursos e shutdown/reabertura. `01_core/usecase/market_registry_test.go` verifica admissão pura sem Discover. Reprodução: `go test ./01_core/usecase ./04_wiring -run 'TestMarketRegistry|TestConfiguredMarketComposition' -count=1`.

Antes da correção, os testes mostraram rejeição do registro padrão pela síntese e enum basal incorreto nos registros ampliados. As capturas congeladas anteriores não foram alteradas. Contrato: `00_nucleo/passos/0104-compor-descoberta-por-registro-validado.md`.
