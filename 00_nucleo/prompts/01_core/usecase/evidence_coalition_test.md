# Testes contratuais: 01_core/usecase/evidence_coalition

Papel autor independente dos oráculos. Derivar os testes apenas da intenção, contrato e
baseline 0065 congelados, sem ler implementação candidata ou adaptar expectativas ao
patch.

Cobrir o corpus mínimo do passo, precedência de conflito, papéis exatos, limiares 1 e 2,
origem desconhecida, mesma origem, deduplicação idêntica e contraditória, input inválido,
schema desconhecido, vazio e todos os limites inclusivos/exclusivos. Verificar input
imutável, slices não nulas, contagem exata de operações e JSON idêntico em pelo menos 20
permutações determinísticas e entre `nil`/vazio.

Construir oráculo de classificação e serialização independentemente. Os testes devem
eliminar M01–M12 do contrato; mutation score exigido 1.0. O futuro arquivo espelhado é
`01_core/usecase/evidence_coalition_test.go`, preferivelmente em package externo.
Congelar testes antes de qualquer execução contra o candidato.

