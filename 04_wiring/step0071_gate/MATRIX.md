# Matriz adversarial 0071

| Mutante | Classe | Testemunha do gate |
|---|---|---|
| M01 default substitui catálogo fornecido | negativo | `TestStep0071GateN01ProvidedCatalogCrossesPipeline` e construtor integrado |
| M02 ausência de regra vira ABSENT | opaco | `TestStep0071GateO01CatalogGapNeverDescribesCandidate` |
| M03 gap vira NOT_VERIFIED | opaco | classificação pública `CATALOG_GAP` no mesmo teste |
| M04 gap entra em CandidateAssessment.Gaps | negativo | `TestStep0071GateO01CatalogGapNeverDescribesCandidate` |
| M05 substring reversa casa regra | opaco | `TestStep0071GateO02UnsafeAndAmbiguousMatchesStayUnknown` |
| M06 primeira regra vence ambiguidade | opaco | duas regras normalizadas aplicáveis no mesmo teste |
| M07 hipótese serializada como declarada | negativo | `TestStep0071GateN02HypothesisHasNoDeclaredAuthority` |
| M08 hipótese lexical vira PROVEN | negativo | símbolo lexical no mesmo teste |
| M09 hipótese sem fonte vira ABSENT | opaco | bundle completo sem testemunho no mesmo teste |
| M10 artifact/hash usa default | negativo | `TestStep0071GateN03CLICatalogArtifactIsCausal` |
| M11 resume ignora catálogo | negativo | duas execuções com catálogos divergentes no mesmo teste |
| M12 reordenação altera saída | determinismo | `TestStep0071GateD01TwentyPermutationsAreByteStable` |
| M13 nil serializa em vez de [] | positivo/schema | `TestStep0071GateP01PublicSchemaAndArrays` |
| M14 cobertura inferida após coleta | opaco/ordem | conteúdo forjado não muda cobertura + `TestStep0071GateN04CoveragePrecedesCollectionStatically` |

Política: somente os casos `O*` podem satisfazer Unknown. `P*` exige preservação,
`N*` exige rejeição do comportamento mutante e `D*` exige igualdade byte a byte.

