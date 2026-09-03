# Prompt — controles relacionais independentes O14 do Passo 0028

Materialize exclusivamente
`04_wiring/step0028_gate/o5blackbox/o14_refinement_test.go` a partir de R4,
catálogo v3, protocolo público O8, artefatos oracle O13 e diagnósticos públicos
O13. Não ler attacks, mutants, records/configs, subjects, candidato, receipts,
seals, certificados ou veredito.

Teste sem IDs de caso:

- `mechanism_id.rule_id` vem do registro da classe no catálogo, enquanto
  `extraction_rule_ids` preserva os RuleIDs das sources;
- witness mínimo insuficiente com gap/teto coerentes é `Unknown`;
- target `UNKNOWN` com baseline ausente, composição não medida, licença incapaz,
  budget esgotado, identidade ausente/ambígua, provider sem capture e efeitos de
  utility/security sem baseline são `Unknown` somente com contexto coerente;
- licença `UNRECOGNIZED` com `PATTERN_ADAPTATION` ou `REFERENCE_ONLY` explícito
  e status `PROPOSED` permanece `Preserved`;
- qualquer falha simultânea continua `Violated`;
- manifest O14 e JSON públicos são estritos, duplicate-free e vinculam somente
  os 28 arquivos oracle allowlisted.

Preserve as APIs `Evaluate`/`EvaluateSequence` e todos os predicados negativos
anteriores. Não emita receipt, seal ou veredito.
