# Bounded local recommendation worker queue

Implemente em L3 a fila local e reservas congeladas no contrato 0026. Aceite entre
1 e 4 workers e capacidade entre 1 e 32; rejeite configuração inválida. `Reserve`
é não bloqueante, `Commit` de reserva válida não falha por capacidade, `Release`
é idempotente e start algum cria goroutine própria.

Workers chamam somente o handler L1 injetado. Mantenha um contexto cancelável por
run ativo; `Cancel` retira job queued ou sinaliza ativo de forma idempotente.
Shutdown para reservas, cancela ativos e aguarda limitadamente. Nunca invente
estado, publique artifact, chame provider diretamente ou importe L2/L4. Não logue
pedido, chave, cursor, conteúdo ou erros externos.

