# Gate discriminatório pré-candidato 0068

Materializar um gate L4 autocontido e *binding-free* para o contrato congelado do
compilador finito de consultas técnicas. O gate lê somente a fixture adversarial pública
de 0068 e avalia snapshots declarativos. É proibido importar ou ler a futura
implementação 0068, escrever oráculos produtivos ou adaptar expectativas ao candidato.

O gate deve preservar controles positivos, violar M01–M16 com uma testemunha pública e
manter controles deliberadamente opacos como `Unknown`. `Unknown` em mutação negativa
é mutante sobrevivente. Deve discriminar tradução inglesa exata, literal intacto,
catálogo finito, três propostas ordenadas, estratégia e proveniência, admissão por dois
discriminadores no mesmo arquivo, rejeição de metadados, execução sem parada precoce,
união com duas origens, limites, compatibilidade legada, propagação de erro e separação
entre evidência endereçável e capacidade comprovada.

O gate também verifica as métricas A/B congeladas: recall@5 `0 -> 1`,
contaminação@5 `1 -> 0.5`, evidência endereçável `0 -> 1`, tentativas `1 -> 4` e duas
origens técnicas para o alvo. Deve distinguir `OBSERVED`, `NOT_OBSERVED` e `Unknown`,
repetindo o vetor de classificação em vinte ordens.

O arquivo espelhado é `04_wiring/step0068_gate/discriminatory_test.go`. Ele não integra
produção, não contém lógica consumível pelo compilador e não julga candidato
inexistente. Autoridade: adversário/gate pré-candidato; entradas congeladas: passo,
contrato, baseline, budget, manifesto e APIs públicas anteriores.
