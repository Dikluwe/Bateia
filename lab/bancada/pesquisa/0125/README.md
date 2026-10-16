# Bancada 0125 — Mesa e referências externas

Data: 2026-10-09. Nenhum draw GPU executado. Fontes Mesa somente lidas;
nenhum driver ou biblioteca Mesa instalado/recompilado nesta bancada.

Resultado legível: [pesquisa](../../../../00_nucleo/relatorios/pesquisa-mesa-0125.md).
Contrato: [passo](../../../../00_nucleo/passos/0125-investigar-clipping-e-transicao-gs.md).

## Artefatos

- `mesa-clipping.summary.json`, `mesa-gs-transition.summary.json`: projeções
  das duas consultas originais. Mantêm diagnóstico, argumentos e inputs com
  hashes. Nenhuma AST final concluiu: todas excederam 32 MiB. Não são novos vereditos.
- `mesa-clipping.json.gz`, `mesa-gs-transition.json.gz`: relatórios CLI originais
  completos, comprimidos sem alteração. Hash dos bytes descomprimidos na projeção.
- `mesa-clipping-focused.json.gz`, `mesa-gs-focused.json.gz`: rodada complementar
  da porta SemanticProvider com Paths fornecidos, preservando pedido e resposta;
  projeções correspondentes em `*.summary.json`. A seleção foi autoral.
- `mesa-sources.zip`, `mesa-manifest.json`: oito fontes centrais, versão, release
  note anterior e banco de compilação de 6 MB. Origem local é uma árvore extraída
  de Mesa 25.2.8; hashes identificam os bytes, não um commit autenticado upstream.
- `external/sources.zip`, `external/manifest.json`: fontes e manifestos originais
  de SwiftShader e ANGLE em commits fixados, com URLs, tamanhos e SHA256.
  Recortes não contêm fechamento de headers/dependências nem contexto de build.
- `swiftshader-clipping.json`: consulta Bateia no recorte, sem análise Clang.
- `angle-gs-transition.json`: consulta Bateia descartada pela triagem; manifesto
  GN não suportado. Isso não indica irrelevância da fonte. Leitura manual separada.
- `oracle-sources.zip`, `oracle-manifest.json`: referência e logs da campanha
  anterior, copiados sem executar oracle/patch.
- `receipts.json`, `go-test.log`: binário usado, consultas, comandos e verificação.
- `provenance-check.json`: conferência dos hashes e vínculos entre relatório e
  arquivos de evidência. Confere integridade; não demonstra semântica.
- `before-output-bound/`: resultados superados anteriores à correção de ReadFrom,
  que permitia contornar o orçamento. Não usar sua contagem de complete como
  cobertura sob o contrato. `focused-gs-interrupted.log` conserva a pilha da
  decodificação interrompida; a rodada final dirigida terminou sob limite real.

## Reprodução das consultas

Na raiz Bateia:

```sh
GOCACHE=/tmp/bateia-go-cache-0125 go build -o /tmp/bateia-mesa-0125 ./04_wiring
```

Em um diretório sem `.env`, definir `mesa_source_root` para o fonte local Mesa
25.2.8. Ele possui `../build/compile_commands.json` original, que referencia o
build e seus headers gerados. Rodar:

```sh
/tmp/bateia-mesa-0125 analyze -path "$mesa_source_root" \
  -query 'clip distance edge flag polygon mode GL_POLYGON' \
  -lang cpp -format json -output /tmp/mesa-clipping-replay.json
/tmp/bateia-mesa-0125 analyze -path "$mesa_source_root" \
  -query 'si_bind_gs_shader si_update_shaders dirty_shaders_mask VGT_GS_MODE' \
  -lang cpp -format json -output /tmp/mesa-gs-replay.json
```

O banco congelado identifica os argumentos originais; o ZIP não contém todos os
headers do build ou do sistema. Replay Clang requer os mesmos bytes de inputs
listados, versão do compilador e caminhos válidos. Sem eles, o resultado pode
ser não avaliado; não reconstruir um contexto falso para reproduzir complete.

Para inspecionar somente os resultados existentes, descomprimir os JSONs e
conferir hashes. `request` dos relatórios focados registra exatamente Root,
Query e Paths passados a `clang.NewClient("clang").Analyze` pela porta existente.
Não houve mudança na CLI para exigir conhecimento prévio dos arquivos.

Os testes discriminantes no relatório são propostas de fixtures CPU/mocks ou
comparação de capturas anteriores. Não foram executados aqui. Nenhuma conclusão
de causalidade do reset, nenhum efeito de patch e nenhum teste GPU é alegado.
