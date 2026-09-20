# Laudo 0038 — Contenção de descartes não calibrados

**Resultado**: aprovado

`ScreenDiscovery` agora separa ranking de autoridade decisória. Modelo explicitamente
não calibrado continua ordenando e avançando candidatos, mas não pode emitir `DROP` nem
probabilidade. Estado contraditório, desconhecido ou com parâmetros não finitos resulta
em `AUDIT`. Calibração válida preserva o comportamento calibrado anterior.

Sinais, pesos, limiares, ordenação e schemas públicos não foram alterados. Esta mudança
reduz risco de falso negativo; não demonstra que a classificação seja boa.
