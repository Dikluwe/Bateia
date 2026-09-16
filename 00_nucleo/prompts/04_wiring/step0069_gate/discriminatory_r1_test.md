# Gate discriminatório binding-free r1 — passo 0069

## Delta contratual

Este prompt espelha `04_wiring/step0069_gate/discriminatory_r1_test.go` sem substituir o
gate r0. A entrada nova é a revisão contratual r1 causada pelo grid fixo do renderer
Archify: rows 0..4, no máximo cinco nodes por stage e labels normativas de até 10 bytes.

O executor permanece apenas adversário. O gate não importa nem lê produção, não executa
Archify e não altera oráculos, manifesto ou selo. Ele classifica receipts observáveis em
memória e não escolhe uma implementação.

## Controle r1

Construa um documento 5/5/5 com identidades longas e carregadas de pontuação. Cada stage
deve ter rows exatamente 0..4. A função `short(s)` opera sobre bytes: cada byte fora de
`[A-Za-z0-9]` vira `_`; até 8 bytes o resultado é conservado, acima disso ficam os sete
primeiros bytes e `~`.

Labels de origem são `O ` + `short(OriginID)`, de perspectiva são `P ` + o código de kind
`PROB|BEHV|ARCH|CONS|OPER`, e de target são `T ` + `short(TargetID)`. Exija a label exata,
ASCII e com no máximo 10 bytes. O ID tipado hexadecimal continua codificando a identidade
integral. Labels de flow preservam projection/origin/observation integrais; card title e
items preservam TargetID e provenance integrais.

## Casos e classificação

Carregue `adversarial-fixtures-r1.json`. P03/P04 são `Preserved`. U04/U05/U06 representam
a sexta origem, perspectiva ou target corretamente recusada com
`render_budget_exceeded`, documento nil e zero bytes, portanto permanecem `Unknown` pela
política do harness. M18–M20 emitem parcial na mesma condição e são `Violated`.

M21–M27 atacam row 5, limite e regra exata de label, identidade abreviada no ID, perdas de
identidade em flows/cards e contagem por bytes UTF-8. Classifique duas vezes, repita o
corpus em 20 ordens determinísticas e exija score r1 `1.0`. Depois do recorte focal passar,
o corpus completo r0+r1 pode ser reexecutado. O resultado é evidência discriminatória
pré-candidato, não selo nem veredito de implementação.
