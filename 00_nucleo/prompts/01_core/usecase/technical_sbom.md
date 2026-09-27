# Declarações SBOM conservadoras

Aceitar somente declaração top-level name@version de JSON CycloneDX 1.6,
após validar estrutura/chaves sem ambiguidades. Não validar schema inteiro
nem presença/runtime/proveniência. Versão não suportada é UNSUPPORTED.
Limitar dois mil componentes e campos de 8192 bytes; bom-ref não vazio deve
ser único. Nomes iguais com versões diferentes são legítimos quando referências
distintas; não confundir coexistência de versões com contradição. Declaração
sem nome ou versão explícitos não comprova componente versionado.
