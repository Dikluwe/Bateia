# Prompt proprietário — `02_shell/report/idea_radar_archify_test.go`

Teste independentemente o contrato 0069: golden e schema, topologia, deduplicação,
UNKNOWN, limites, IDs adversariais collision-free, origem desconhecida, cards legíveis,
ausência de probe/relação inventada, 20 permutações, nil/vazio e imutabilidade profunda.
Não leia nem adapte os casos à implementação candidata; zero rede, render, relógio ou
aleatoriedade.
Cubra 5/5/5 inclusivo, cada sexto exclusivo, rows 0..4, labels <=10 bytes e preservação da
identidade completa. O harness externo, não o teste Go, valida e entrega o exemplar no
commit Archify fixado; não exija `showcase` universal.
