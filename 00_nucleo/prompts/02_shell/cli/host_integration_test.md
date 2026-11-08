# Composição do host nos testes de integração CLI

**Camada:** L2
**Proprietário:** `02_shell/cli/host_integration_test.go`
**Contrato:** 0129. O TestMain externo compõe o driver L3 antes das suítes de integração; nenhum import de driver é introduzido no produto L2. O teste de integração verifica a CLI com efeitos reais confinados às fixtures existentes.
