# Recusa de efeitos sem host

**Camada:** L2
**Proprietário:** `02_shell/cli/host_test.go`
**Contrato:** 0129. Verificar que a implementação indisponível retorna erro para leitura/escrita/abertura sem executar efeitos, e não inventa data ou identidade de arquivos. Não alterar configuração global de outros testes.
