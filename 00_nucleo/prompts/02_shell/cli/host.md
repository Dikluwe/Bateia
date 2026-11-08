# 02_shell/cli/host

**Camada:** L2
**Proprietário:** `02_shell/cli/host.go`
**Contrato:** passo 0129 e ADR-0012.

Receber CommandHost por ConfigureCommandHost, sem fallback concreto. Propagar runtime por contexto. A configuração é de startup; não alterar durante execução. Handles de I/O vêm apenas da porta.
