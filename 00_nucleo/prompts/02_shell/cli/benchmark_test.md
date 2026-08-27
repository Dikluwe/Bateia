# Prompt — testes da CLI offline do benchmark

Ataque a fronteira `bateia benchmark`: prove que validação e execução funcionam
somente com arquivos locais, que denominadores vazios permanecem `null`, e que
entradas ausentes e subcomandos desconhecidos falham de forma explícita. Os testes
não podem depender de tokens, rede ou relógio.
