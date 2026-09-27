# 0100 — despacho offline da comparação

Executar o main real em subprocesso do binário de teste com idea trace compare.
Criar .env com segredo sintético que coincida com uma declaração visível na saída.
Retirar a variável correspondente do ambiente herdado. A comparação deve emitir
JSON válido e preservar a declaração, demonstrando que .env não foi carregado.
Não utilizar credenciais reais nem chamadas de rede. Prazo limitado a dez segundos.
