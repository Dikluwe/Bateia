# Prompt proprietário — `01_core/domain/interfaces.go`

Origem L0 exclusiva de `01_core/domain/interfaces.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/domain/types.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Materializar somente a responsabilidade desta unidade, preservando os contratos observáveis: `SearchProvider`, `ObservableDiscoveryProvider`, `ContentProvider`, `CacheStore`, `ProgressReporter`. Manter a direção de dependências da camada e não incorporar responsabilidades pertencentes aos demais consumidores do prompt de origem.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/domain/types.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Passo 0035 — Análise semântica C/C++

Cumprir `00_nucleo/passos/0035-corrigir-lint-e-analisar-c-cpp-com-clang.md`. Porta L1 SemanticProvider, execução Clang limitada em L3, injeção L4. Contexto explícito em compile_commands.arguments, sem executar command, plugins ou response files; capturar hashes de configuração, fonte e headers, argumentos, versão, símbolos e chamadas com intervalos. Falta de evidência nunca vira complete. Exposição L2 separa resolução estática e adequação à consulta.
