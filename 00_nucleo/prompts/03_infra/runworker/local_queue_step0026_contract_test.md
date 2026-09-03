# Prompt espelhado — oráculo da fila local 0026

## Autoridade

Teste independente L3 baseado somente no Passo 0026 e contrato R1.

## Obrigações

- Rejeitar workers fora de `[1,4]` e capacidade fora de `[1,32]`.
- Reserva é atômica; commit de reserva válida não falha por capacidade.
- Não iniciar workers antes de `Start`, não exceder concorrência configurada e não criar worker por request.
- Cancelar contexto da execução ativa e liberar capacidade sem reiniciar o mesmo run.
- Liberar reservas perdedoras/abortadas sem vazamento, inclusive sob concorrência e `go test -race`.

## Limites

Fila inteiramente local, handler fake bloqueável e timeouts limitados; nenhum acesso a rede ou store real.
