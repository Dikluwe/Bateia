# Matriz discriminatória — passo 0073

| Mutante | Testemunha principal |
|---|---|
| M01 | `TestStep0073GateCitesExactCppDeclarationAndExpression` |
| M02 | `TestStep0073GateCitesExactCppDeclarationAndExpression` |
| M03 | `TestStep0073GateRejectsPathCommentsStringsSubstringsAndOpaqueSyntax` |
| M04 | `TestStep0073GateRecognizesCppMemberAndCallWithoutLanguageKeyword` |
| M05 | `TestStep0073GateCitesExactCppDeclarationAndExpression` |
| M06 | `TestStep0073GateCitesExactCppDeclarationAndExpression` |
| M07 | `TestStep0073GateRejectsPathCommentsStringsSubstringsAndOpaqueSyntax` |
| M08 | `TestStep0073GateRejectsPathCommentsStringsSubstringsAndOpaqueSyntax` |
| M09 | `TestStep0073GateRejectsPathCommentsStringsSubstringsAndOpaqueSyntax` |
| M10 | `TestStep0073GatePrecedenceAndDistinctOccurrencesAreCanonical` |
| M11 | `TestStep0073GateAssessmentIsDeterministicAndInputsImmutable` |
| M12 | `TestStep0073GatePrecedenceAndDistinctOccurrencesAreCanonical` |
| M13 | `TestStep0073GateRejectsPathCommentsStringsSubstringsAndOpaqueSyntax` |
| M14 | `TestStep0073GateAssessmentIsDeterministicAndInputsImmutable` |

Controles opacos preservam `COLLECTION_INCOMPLETE`, `CATALOG_GAP` e
`HYPOTHESIS_ONLY`. Chamada e declaração C++ são cobertas dentro do recorte lexical, sem
alegar resolução semântica. Mutation score exigido: `14/14 = 1.0`.

