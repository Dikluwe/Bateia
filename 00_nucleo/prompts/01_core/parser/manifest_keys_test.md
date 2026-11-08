# Verificação das chaves de dependências dos manifestos

**Camada:** L1
**Proprietário:** `01_core/parser/manifest_keys_test.go`
**Contrato:** helper puro manifest_keys.go no passo 0128.

Verificar propriedade de não mutação do mapa/valores e independência da lista devolvida, inclusive após alterar uma lista anterior. Confirmar ordem canônica para chave ASCII e Unicode e comportamento de mapa vazio/nulo. Sem rede, disco ou relógio.
