# Testes independentes da identidade e trace 0064

Papel B, regime A/B em filesystem compartilhado, sem atestação física. Entradas: contrato0064 congelado, contratos0062/0063 e API domain baseline. Não consultar produção candidata, testes legados, fixtures ou relatórios.

Criar testes externos de effects.NewRunner().Verify. Fontes autônomas controladas definidas no próprio teste. Comparar identidade com SHA-256 independente do executável corrente (/proc/self/exe em Linux); conferir metadados runtime/debug quando disponíveis, horários UTC e fingerprint por serialização independente. Exercitar efeito supported, repetição em caminho diferente/expectativa diferente, mismatch unknown sem casos ou execução, rejeição de recorte unknown rastreável. Conferir parâmetros pedidos e limites efetivos. Não executar antes do congelamento dos artefatos.

## Revisão focal B r1 — casos planejados em preflight

A execução integrada identificou exigência extra de lista vazia em dois preflights. O contrato0063 exige zero casos completos, permitindo conservar casos planejados incompletos do0062; o0064 preserva esse comportamento. Antes da revisão de código, ajustar exclusivamente essas duas condições para rejeitar qualquer Observation.Completed, mantendo status, razão, identidade, declaração, vínculo e demais expectativas. Não consultar implementação candidata nem logs completos. Congelar novamente sem executar testes; r0 preservado em /tmp/bateia0064-b-r0.sha256.
