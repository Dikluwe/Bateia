# 0102 — orçamento agregado da projeção de saída

Testar boundedReviewContextProjection, admissão final privada usada pelo facade.
Construir projeção tipada sintética de 100 tarefas com 100 ancestrais cada e notas
no limite de 8192 bytes. Ledger e entries devem caber isolados; contexto agregado
excede e retorna erro controlado com resultado zero. Reduzir notas para vizinho
positivo, sem truncar ancestrais. Este teste verifica orçamento, não a coerência
semântica do join (validada nos testes do facade público). Sem I/O/rede.
