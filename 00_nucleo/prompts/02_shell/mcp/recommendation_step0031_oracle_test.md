# Oráculo independente — defaults de recomendação MCP do Passo 0031

## Linhagem

- Passo autoritativo: `00_nucleo/passos/0031-autopreencher-parametros-da-recomendacao.md`
- Contrato candidato: `00_nucleo/benchmark/0031/contract-candidate.json`
- Baseline: commit `6aae9ea29b65bb942d14aff2fb6d659668f01bd5`
- Papel: autor segregado de oráculo e ataques, sem leitura da implementação candidata.

## Observáveis

Exercitar somente as superfícies MCP públicas e ports in-memory injetados. O schema de
`bateia_recommend` deve exigir apenas `idea`; o de `bateia_start_recommendation`, apenas
`idempotency_key` e `idea`. Ambos devem manter `problem` e `capabilities` obrigatórios,
objetos estritos, schemas operacionais profundamente equivalentes e anotações dos defaults
estáticos V1. `target.name` deve indicar inequivocamente que a omissão deriva de
`idea.problem`.

Uma intenção mínima, assim como objetos operacionais vazios, deve chegar aos dois ports com
o perfil `BATEIA_RECOMMEND_DEFAULTS_V1` exato. Cada chamada deve receber slices próprios.
Valores explícitos válidos prevalecem, incluindo `false`, string vazia e listas vazias;
uma requisição completa previamente válida deve chegar semanticamente inalterada.

A allowlist V1 é somente política de aceitação. Um candidato projetado com licença
`INCONCLUSIVE`, sem identificador nem fonte, deve permanecer opaco mesmo quando a mesma
resposta expõe a allowlist efetiva. Positivos e negativos são repetidos em ordens
determinísticas direta e inversa, alternando também o primeiro caminho observado.

Zeros inválidos, `null`, strings de domínio inválidas, campos desconhecidos, intenção de
negócio ausente e relação inválida entre thresholds devem produzir JSON-RPC `-32602` antes
de qualquer chamada downstream, igualmente nos caminhos síncrono e assíncrono.

## Limites da alegação

O oráculo observa schema, materialização e encaminhamento. Não prova qualidade de busca,
completude jurídica da allowlist, proveniência interna de valor explícito igual ao default
nem equivalência funcional geral.
