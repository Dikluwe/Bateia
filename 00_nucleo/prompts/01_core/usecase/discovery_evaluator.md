# Prompt Spec L0: Avaliador de Descoberta

## Objetivo

Avaliar snapshots congelados por julgamentos explícitos, sem confundir candidato
não julgado com irrelevante.

## Restrições

- L1 puro e determinístico.
- Precisão desconhecida é `null` quando não há julgados positivos/negativos.
- Não calcular recall de universo aberto.
- Rejeitar julgamentos duplicados, conflitantes ou inválidos.
- Medir precisão julgada, incerteza, não julgados, hit rate, yield, overlap e falhas.
