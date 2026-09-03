# Prompt espelhado — oráculo de domínio dos runs duráveis 0026

## Autoridade

Teste A/B independente derivado exclusivamente do Passo 0026 e do contrato observável R1. Não adaptar à implementação candidata.

## Obrigações

- Congelar valores de lifecycle, estágios, nomes e estados de artefatos.
- Congelar JSON público de recibo e visão de run, incluindo `null` explícito e arrays `[]`.
- Distinguir lifecycle operacional de `recommendation_status`; `SUCCEEDED` pode carregar `INCONCLUSIVE`.
- Impedir que a visão pública exponha pedido, hash ou digest de idempotência.

## Limites

O teste é puro, determinístico e sem I/O, rede, relógio ou dependência de ordem acidental.
