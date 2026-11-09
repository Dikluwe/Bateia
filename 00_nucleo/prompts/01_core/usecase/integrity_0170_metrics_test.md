# Métricas comprovadas e ordem total das alternativas

Testar `EvaluateEvidenceCalibration` com Expected e Observed PROVEN em cada split
CONTRACT/VALIDATION/TEST. Fontes ausentes, evidência ausente, commit divergente e
referência mutável devem conservar denominadores ExpectedProven/ObservedProven e
InvalidSources, porém TrueProven, precisão e recall ficam zero no split e no
agregado. Controle imutável mantém taxas 1; corpus misto válido/inválido produz
1/2. As observações não são modificadas.

Testar `RankSoftwareAlternatives` com snapshots duplicados que diferem somente em
Stale, Archived, branch padrão, data atualizada ou tentativas de busca. Permutar
os snapshots, inclusive repetindo um deles; relatório, representante e motivos
permanecem idênticos, inputs intactos e duplicatas continuam INCONCLUSIVE sem
TopChoice. No empate dos campos prioritários, a ordem pelo JSON completo seleciona
Stale=false antes de Stale=true; não fabricar diagnóstico de identidade opaca a
partir do representante preterido. O controle com candidato único comprovado
mantém RECOMMENDED, cobertura 1 e primeira escolha. Não combinar escalares nem
promover evidência por metadados.
