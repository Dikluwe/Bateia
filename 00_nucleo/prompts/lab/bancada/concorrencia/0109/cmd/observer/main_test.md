# Prompt 0109 — cmd/observer/main_test

**Regime**: Lab, sem promoção para L1–L4  
**Data**: 2026-10-06

Testar comandos observe/evaluate por fronteira pública CLI, entrada desconhecida, duplicada, excesso de bytes, múltiplos JSON, saída existente/alias preservada, output 0600, equivalência stdout/arquivo, short writer e evaluate com expectativa errada publica resultado e falha.

0111: contrastar observe com/sem ambiente fonte, JSON extra/duplicado inválido,
ambiente inválido sem publicação, evaluate com ambiente rejeitado.

0112: no congelamento conc, verificar uma ligação do case de limiter ao
lançador, duas chamadas tipadas desse lançador no Pool.Go (default e case),
ligação do defer ao receive e compartilhamento de declaração do campo. Validar
hash/posição em ambas extremidades; lacunas de callback/instância permanecem.

0113: verificar no conc congelado dois bindings no parâmetro de WaitGroup.Go,
forwarding capturado para o parâmetro de Catcher.Try e invocação desse parâmetro.
Comparar localizadores de declaração para compor vínculo estrutural; escopo
físico da goroutine explícito e CALLBACK_EXECUTION_UNRESOLVED preservado.
