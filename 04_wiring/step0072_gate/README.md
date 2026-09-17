# Gate adversarial — passo 0072

Oráculo segregado derivado do passo 0072, dos testes contratuais pré-candidato e do
baseline `03fd5f1`. O autor deste diretório não leu nem corrigiu a implementação
candidata.

## Cobertura dos mutantes

| Mutante | Testemunha principal |
|---|---|
| M01 | `TestStep0072GateLocalTreeIncludesEntryAfterOldTwentyThousandCut` |
| M02 | `TestStep0072GateTargetedReservationBeatsDistractorsAcrossPermutations` |
| M03 | `TestStep0072GateTargetedReservationBeatsDistractorsAcrossPermutations` |
| M04 | `TestStep0072GateTargetedReservationBeatsDistractorsAcrossPermutations` |
| M05 | `TestStep0072GateTargetedReservationBeatsDistractorsAcrossPermutations` |
| M06 | `TestStep0072GateTargetedReservationBeatsDistractorsAcrossPermutations` |
| M07 | `TestStep0072GatePathMatchNeverPromotesCapability` |
| M08 | `TestStep0072GateSkippedPlausiblePathIsExplicitUnknown` |
| M09 | `TestStep0072GateSkippedPlausiblePathIsExplicitUnknown` |
| M10 | `TestStep0072GateSkippedPlausiblePathIsExplicitUnknown` |
| M11 | ambos os testes de vinte permutações |
| M12 | `TestStep0072GateSharedTargetIsDeduplicatedWithoutLosingCoverage` |
| M13 | `TestStep0072GatePlanningIsDeclaredDeterministicAndImmutable` |
| M14 | `TestStep0072GateSmallGoTreeKeepsLegacyCollection` |

Os controles opacos cobrem budget insuficiente e falha de enumeração. Os negativos
cobrem substring interna, extensão documental, arquivo extensionless irrelevante e path
match sem autoridade semântica. Mutation score requerido para selo: `14/14`.

## Execução

```sh
GOCACHE=/tmp/bateia-step0072-gate-cache go test ./04_wiring/step0072_gate -run '^TestStep0072Gate' -count=1
```

Este material registra segregação procedural de autoridade. O filesystem é compartilhado,
portanto não constitui atestação física de isolamento.
