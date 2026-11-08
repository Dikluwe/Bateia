# Parsing determinístico e instante explícito de coleta

**Camada:** L1
**Proprietário:** `01_core/parser/manifest_time_test.go`
**Contrato:** passo 0128 e ADR-0011.

Caracterizar um corpus autoral de todos os tipos de manifesto do dispatcher, congelando o digest da projeção antes da alteração, excluindo apenas ConsultadoEm. Verificar a mesma projeção depois. Testar duas chamadas idênticas, ausência de data em ParseManifest, data explícita em ParseManifestAt, conversão UTC sem alterar o instante, zero desconhecido e propagação dos erros. Não acessar relógio, disco ou rede. Testar parsing de dependências relevantes pelos testes consumidores já existentes.

## Ordem estável dos objetos de dependências

A revisão encontrou iteração não canônica de mapas em package.json, composer.json, box.json, haxelib.json e dub.json. Ordenar as chaves dentro de cada grupo, preservando ordem dos grupos e nomes/versões/tipos/exclusões. Testar múltiplas dependências e grupos repetidamente. A ordem dessas listas passa a ser definida; a caracterização do corpus temporal permanece intacta. Artefatos antigos não são reordenados ou re-hasheados.
