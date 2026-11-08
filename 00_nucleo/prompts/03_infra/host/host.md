# 03_infra/host/host

**Camada:** L3
**Proprietário:** `03_infra/host/host.go`
**Contrato:** passo 0129 e ADR-0012.

Implementar portas de runtime e arquivos sobre os.Root/os.File. Traduzir flags não bloqueantes/nofollow, reter handles e preservar erros. Relógio, sinais, ambiente e entropia pertencem aqui. Não importar L2/L4.
