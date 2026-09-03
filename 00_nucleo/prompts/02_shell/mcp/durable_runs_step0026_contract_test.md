# Prompt espelhado — oráculo MCP dos runs duráveis 0026

## Autoridade

Teste A/B independente derivado do Passo 0026, contrato R1 e superfície pública certificada no Passo 0025.

## Obrigações

- Com a dependência durável injetada, anunciar exatamente sete tools; sem ela, preservar byte/deep-equality das três tools 0025.
- Congelar os quatro schemas novos, com objetos fechados e objetos de recomendação deep-equal ao schema síncrono.
- Validar JSON, IDs, chave, artifact, cursor e limites antes da porta L1.
- Rotear uma vez à operação correta e emitir envelopes textuais compactos com arrays `[]`, ponteiros `null` e entries JSON, não strings/base64.
- Mapear categorias públicas para códigos/mensagens exatos, sem detalhes físicos, canários ou erro de adaptador.

## Ambiente

Servidor `stdio` inteiramente em memória, porta fake controlada, zero filesystem e rede.
