# Efeitos compartilhados da pesquisa acadêmica

**Camada:** L3  
**Arquivo proprietário:** `03_infra/research/common.go`

Preservar transporte HTTP limitado a 4 MiB, timeout finito, bloqueio de redirects, status explícito, saneamento de texto/URLs e hash de resposta. Timestamp operacional pertence a L3. Manter distinção entre resposta vazia observada e falha do provedor.

## Regularização de linhagem

Este contrato formaliza a adição local já existente, após leitura do código e testes. A regularização acrescenta especificação e cabeçalho; não altera o comportamento do componente nem representa aprovação de sua utilidade em ambiente remoto.
