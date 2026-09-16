# Testes contratuais: 01_core/usecase/idea_radar

Papel autor independente de oráculos e benchmark. Derivar testes somente do passo,
contrato e baseline 0066 congelados, sem ler a implementação candidata.

Cobrir os seis casos/famílias, os dois limiares independentes, mesma origem, kind
repetido, origem desconhecida, conflito, ambiguidade, duplicação, canonicalização ASCII,
rejeição Unicode/fuzzy, vazio, schema, limites inclusivos/exclusivos, input imutável,
arrays não nulos, operações exatas e 20 permutações/nil-vazio com JSON idêntico.

O oráculo de utilidade permanece externo ao código produtivo: selecionar a melhor
projeção isolada e calcular precisão, recall, F1, falsa convergência, contaminação e ganho
incremental a partir do `baseline.json`. Exigir os estados `OBSERVED`, `NOT_OBSERVED` e
`INCONCLUSIVE` nas suas condições congeladas. Eliminar M01–M14 com score 1.0.

O futuro arquivo é `01_core/usecase/idea_radar_test.go`, preferencialmente em package
externo. Congelar seu hash antes da primeira execução contra o candidato.
