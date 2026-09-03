# Environment file adapter

**Camada:** L3  
**Arquivo proprietário:** `03_infra/envfile/load.go`

Encapsule a leitura opcional de um arquivo `.env` e a publicação somente de
variáveis ainda ausentes. Falhas de leitura são não fatais. O parser aceita
linhas `KEY=VALUE`, prefixo `export`, comentários e aspas externas simples ou
duplas. Não registre nem retorne valores, tokens ou conteúdo do arquivo.
