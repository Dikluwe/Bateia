# Passo 0035 — Resultado e reprodução

Implementado: precisão do lint Go/V12 e análise estática C/C++ via Clang no funil local. Regime A/B para os testes semânticos; execução **sem atestação de isolamento de capacidades**. Os certificados históricos permanecem históricos.

## Lint

195 avisos → 0. `crystalline-lint . --fail-on warning` retorna `✓ No violations found`. O patch diferencia funções Go, métodos, tipos, testes e arquivos exatos de harness; V1/V3 continuam ativos. Detalhes, regressões e reprodução em [lint-report.md](lint-report.md) e `reproduce-lint.sh`.

O binário corrigido foi instalado em `/home/dikluwe/.cargo/bin/crystalline-lint`, com backup em `/tmp/bateia-0035-crystalline-lint.before`. Hashes e localização em `lint-installation.json`. A fonte externa não foi modificada; o patch versionado permite reproduzir a instalação. V12 é estrutural e **não prova ausência de regras de negócio dentro de funções livres**.

## Integração semântica

- Porta L1: `SemanticProvider.Analyze`; sem I/O direto no núcleo.
- Adaptador L3: `03_infra/clang`. Usa `compile_commands.json` na raiz ou `build/`, campo `arguments`. Ignora o executável indicado pelo repositório e invoca apenas o Clang configurado.
- O campo `command` isolado é reportado como não avaliado. Não há shell, build, plugin ou expansão de response files. Opções desconhecidas também impedem completar a unidade.
- Dependências de compilação, incluindo headers de macros e headers de sistema, são enumeradas pelo Clang e têm bytes/hashes capturados antes da AST e comparados depois. Isso detecta mudanças entre essas leituras, mas não constitui snapshot atômico contra alterações transitórias.
- Declarações, definições, tipos e chamadas têm arquivo/linha; identidades de chamadas distinguem sobrecargas quando o Clang fornece nome mangled. Chamadas virtuais/indiretas não demonstram destino em execução.
- Contextos ausentes, erros, cancelamento, mudança de inputs e limites preservam `partial` ou `not_evaluated`; AST completa se refere somente às unidades analisadas.
- Limites: 12 unidades, 20 segundos por invocação, 32 MiB de AST, 128 KiB de diagnóstico, 128 MiB de entradas por unidade; seleção lexical inicial lê até 64 MiB de fontes. Linguagem/flags/headers são os do contexto registrado, sem inventar includes.
- O relatório Markdown apresenta a evidência semântica separada da relevância. JSON conserva símbolos, chamadas, argumentos e hashes completos. A cobertura global e a saúde continuam sem promoção automática.

Exemplo de uso com banco de compilação em formato `arguments`:

```sh
./bateia analyze -path /caminho/do/projeto -query 'order independent transparency opengl' -lang cpp -format json -output /tmp/analysis.json
```

## Coin real

A validação usou cópia dos bytes de `/mnt/Laranja/Git/externos/coin`, com os parâmetros já gerados em `/tmp/coin-contribution-review-build`. O script reconstruiu `arguments` de `flags.make`/`build.make`, sem executar CMake ou recompilar o projeto. O commit-base original, caminhos e hashes dos parâmetros constam de `coin-inputs.json`. A cópia não possui `.git`; sua origem está nesse recibo, não em um commit inventado no relatório.

Resultado pela CLI: **12 de 898 unidades**, status global **partial**. As 12 unidades selecionadas terminaram com AST válida, incluindo `SoGLRenderAction.cpp`, `SoTransparencyType.cpp`, `SoTransparencyElement.cpp`, `gl.cpp` e `SoOffscreenRenderer.cpp`.

Fatos estáticos em `src/actions/SoGLRenderAction.cpp`:

| Implementação | Intervalo | Evidência relevante |
| --- | --- | --- |
| `handleTransparency` | 1205–1311 | Seleção `SORTED_LAYERS_BLEND`, configuração de fragment programs/register combiners e manipulação do estado de transparência. |
| `doSortedLayersBlendRendering` | 2042–2075 | Inicialização, textura, iteração por camadas e chamadas `renderOneBlendLayer`, `renderSortedLayersFP` ou `renderSortedLayersNV`. |
| `renderOneBlendLayer` | 2117–2200 | Implementação de uma camada. |
| `renderSortedLayersFP` | 2525–2580 | Caminho de composição por fragment program. |
| `renderSortedLayersNV` | 2583–2678 | Caminho alternativo de composição. |

Inferência: há implementação pertinente de transparência por camadas para investigar a consulta OIT. Não verificado: correção visual, compatibilidade da GPU, ordem independente em todos os casos, cobertura de todo o Coin, execução ou benchmark. A enumeração acima resulta de símbolos e leitura dos intervalos registrados, não do score de busca.

Artefatos:

- [coin-summary.json](coin-summary.json): cobertura e contagens.
- [coin-analysis.md](coin-analysis.md): apresentação gerada do mesmo JSON da CLI.
- `coin-analysis.json.gz`: relatório JSON integral comprimido; hash dos bytes descomprimidos em `coin-summary.json`.
- `coin-cli.log`: telemetria do Bateia, sem benchmark do Coin.
- `coin-inputs.json`: origem do snapshot e parâmetros existentes.
- `prepare-coin-snapshot.py`: reprodução da cópia/configuração a partir dos caminhos registrados.

## Verificação

- Oráculo A/B congelado antes da implementação: `03_infra/clang/semantic_contract_test.go`, SHA-256 `a63e06e6ea8ed345d7b918426dec13fee0542383a9a5cae89e8096bb8461ce7d`. Não foi adaptado ao patch.
- `go test ./...` aprovado; regressões complementares verificam a injeção no funil, apresentação e limite de unidades. Logs em `go-tests.log` e `semantic-tests.log`.
- Ferramenta lint: 1095 testes Rust e seis fixtures CLI aprovados, incluindo casos negativos de produção e preservação de imports/linhagem.
- Build do pacote inteiro: `go build -o bateia ./04_wiring`.
