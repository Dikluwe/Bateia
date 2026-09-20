# Contrato canônico 0070

Autoridade normativa: `../../passos/0070-reconstruir-historia-de-algoritmos.md`.

`Preserved` exige análise pura e determinística de alegações históricas tipadas, cada uma
com evidência citada, proveniência rastreável e inferências explicitamente separadas.
Outcomes seguem exclusivamente as contagens públicas de evidência e origem: `CONTESTED`,
`UNKNOWN`, `CONTRADICTED`, `CORROBORATED` ou `SUPPORTED`, na precedência normativa.
Inferências permanecem visíveis, mas nunca contam como evidência ou corroboração.

`Violated` inclui inferência promovida a evidência, fonte duplicada contada como origem
independente, relação inventada/transitiva, `RESEMBLES` promovido a causalidade, perda de
controvérsia, UNKNOWN convertido em sucesso, referência órfã, dedup dependente de ordem,
mutação do input, parcial em falha, I/O em L1, busca/rede/LLM, Archify ou overwrite do
input/destino na CLI.

Entrada inválida, ambígua, opaca ou acima do budget produz relatório global `UNKNOWN`,
assessments vazios e reason pública. UNKNOWN epistemológico por alegação permanece em
relatório `COMPLETE` e nunca satisfaz um caso positivo que exija suporte. O contrato cobre
somente o fragmento observável do passo 0070, não verdade histórica geral.
