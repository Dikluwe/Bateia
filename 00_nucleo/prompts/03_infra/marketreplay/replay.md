# Frozen market replay

**Camada:** L3  
**Arquivo proprietário:** `03_infra/marketreplay/replay.go`

Implemente o replay offline da projeção congelada contratada, verificando versão e
SHA-256 antes de entregar dados à porta L1. O adapter não possui transporte HTTP,
DNS, env ou fallback online; mismatch é `INTEGRITY_ERROR`, nunca vazio. Para a
mesma entrada, devolva bytes semanticamente idênticos e preserve cancelamento.
Importe L1, nunca L2/L4.

