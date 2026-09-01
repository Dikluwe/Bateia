# Laudo 0023 — Validação vertical de decisões reais de reúso

**Data**: 2026-08-31
**Estado**: `EXECUTADO`
**Veredito**: `FAIL_DISCOVERY`
**Domínio**: busca e avaliação de software
**Regime Tekt**: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

## Resultado

O benchmark reutilizou os dez IDs congelados no passo 0022 e o run autenticado
selado do passo 0020. Nenhum braço atingiu o mínimo de 80% de queries concluídas:

| Braço | Concluídas | Planejadas | Taxa |
|---|---:|---:|---:|
| N | 4 | 12 | 33,33% |
| F1 | 6 | 12 | 50% |
| F2 | 6 | 12 | 50% |

Conforme o protocolo, o benchmark parou na primeira fase. Evidência, decisão,
utilidade e mutation score do produto ficaram `NOT_EVALUATED`; não foram
convertidos em zero nem em aprovação. Assim, este run não autoriza afirmar que o
Bateia já funciona ponta a ponta para pesquisar ideias.

## Materializado

- gate offline por braço, resistente a média agregada enganosa;
- validação de dez IDs únicos e contagens coerentes;
- short-circuit causal `FAIL_DISCOVERY`;
- diagnóstico reproduzível e replay byte a byte;
- comando `bateia benchmark product-gate` sem rede ou credenciais.

## Verificação técnica

Passaram `go test ./...`, `go vet ./...`, `staticcheck 2026.2.1`, Crystalline
Lint e `git diff --check`. Esses gates comprovam a execução do avaliador, não o
gate de produto, que falhou legitimamente.

## Próxima condição

Executar uma nova descoberta autenticada em janela de cota capaz de concluir ao
menos 10 de 12 queries em cada braço. Somente então evidência, decisão e utilidade
podem ser avaliadas no mesmo corpus congelado.

## Atualização após retomada

O `run-5` posterior concluiu 12/12 queries em N, F1 e F2. Portanto a causa
`FAIL_DISCOVERY` deste run histórico foi resolvida. O passo volta a avançar, mas
evidência, decisão e utilidade ainda precisam ser recalculadas sobre os snapshots
novos; não são promovidas automaticamente a partir do relatório antigo.
