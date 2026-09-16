# Contrato canônico 0069

Autoridade normativa: `../../passos/0069-exportar-radar-para-archify.md`.

Preserved exige um Dataflow v1 determinístico com três stages, topologia exclusivamente
origem→perspectiva→target, observações canônicas como flows e outcomes/probes apenas como
texto descritivo. IDs são tipados e hexadecimais, válidos e collision-free. Bateia chama
0066/0067, rejeita qualquer `UNKNOWN`, não muta input e não executa renderização.

R1 exige no máximo cinco origens únicas, cinco perspectivas e cinco targets, rows 0..4 e
labels de node de até 10 bytes pela função `short` normativa do passo. Acima disso é
`render_budget_exceeded`, sem parcial. Identidades completas ficam nos IDs hexadecimais,
labels de flow e cards. O exemplar congelado deve passar `archify validate` e
`archify deliver` no commit `6db72a9aea3d0f67a6a034e41f8a5491476a11c1`; nenhuma
exportação comum é obrigada a usar perfil `showcase`.

Violated inclui relação inventada, probe convertido em flow, colapso de origem desconhecida,
sanitização colidente, label maior que 10 bytes, row fora de 0..4, mais de cinco nodes em
qualquer stage, dependência de ordem, parcial em UNKNOWN/budget, overwrite do input
ou execução de Node/Archify. Entrada/schema opaco e budget excedido permanecem Unknown no
harness, nunca sucesso. O contrato cobre somente esta exportação observável.
