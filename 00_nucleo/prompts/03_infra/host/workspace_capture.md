# Captura física e lifecycle local

**Camada:** L3
**Proprietário:** `03_infra/host/workspace_capture.go`
**Contrato:** 0129 e 0097. Possuir container temporário, descritores, cópia na mesma leitura do hash, modos, quotas, retenção/cleanup e identidade física da captura. Projeção/Close por porta; CLI não implementa outro provider físico. Preservar tree-sha256-v2, frozen-local-v1 e todos os limites e erros de segurança.
