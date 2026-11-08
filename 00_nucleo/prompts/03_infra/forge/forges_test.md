# Testes dos catálogos de forjas

**Camada:** L3
**Proprietário:** `03_infra/forge/forges_test.go`
**Contrato:** passo 0126 e prompt `forges.md`.

Verificar HTTP de Gitee e AtomGit por transporte falso: endpoint, consulta Unicode, token no mecanismo de cada provedor, limite e uma única página. Conferir identidade com namespace aninhado, URLs AtomGit e legado GitCode, visibilidade booleana/inteira e descarte de privados ou identidade incompatível. Conferir hash SHA-256 prefixado e instante de coleta; nunca promover branch ou metadados a commit verificado. Cobrir falta de token sem rede, cancelamento, redirect sem propagação de credenciais, corpo excedido, JSON inválido, 401/403/404/429/503 e array vazio. Preservar cliente injetado e evitar merge entre plataformas por nome ou branch. Testes não usam credenciais reais ou rede.
