# Prompt — controles documentais independentes O13 do Passo 0028

Materialize exclusivamente
`04_wiring/step0028_gate/o5blackbox/o13_refinement_test.go` como self-tests
binding-free derivados do contrato R4, catálogo v3, protocolo O8 e feedback
público O12. Não ler attacks, records/configs, subjects, candidato, receipts,
seals, certificados ou veredito.

Os controles devem provar genericamente:

- bytes do preimage causal em ordem fixa, sem HTML escaping ou newline;
- normalização local NFC de formas canonicamente equivalentes;
- inclusão, deduplicação e ordenação de todos os hashes de evidence e
  counter_evidence;
- baseline ausente com mecanismo `PROVEN` e `novelty=UNKNOWN` produz `Unknown`,
  sem mascarar divergência decidível;
- licença `UNRECOGNIZED`, warning correspondente e `integration_unknowns`
  explícitos não produzem `Unknown` por si;
- o manifest protegido O13 vincula somente código e prompts oracle
  allowlisted, enquanto o manifest O11 permanece histórico;
- todos os artefatos JSON O13 são documentos completos, UTF-8 e sem chaves
  duplicadas, validados localmente sem ferramenta externa.

Preserve `Evaluate` e `EvaluateSequence`, classificação com precedência
`Violated > Unknown > Preserved`, JSON estrito e independência de IDs/bindings.

Após O13→O14, o manifest O13 é somente histórico e byte-estável; seus hashes não
devem continuar vinculando o harness corrente refinado.
