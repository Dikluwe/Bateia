# Public end-to-end calibration O20b

This is a public, sanitized protocol projection. It is not an attack, protected fixture, test vector, oracle expectation, or candidate input.

Each source exchange preserves the public envelope and artifact-document structure observed at the `Evaluate` boundary. Neutral anonymous values replace all protected content. The replay derives mandatory transport-integrity fields in memory and sends serialized `ContractInput` as the first argument and serialized `ContractOutput` as the second. The canonical artifact document is carried in the second argument at `/return_value/canonical_json`.

Run from the repository root:

```text
go run ./00_nucleo/benchmark/0028/public-end-to-end-calibration-o20b-replay.go
```

Against oracle O20, the seven exchanges classify as `Preserved` (7), `Unknown` (0), and `Violated` (0). This observation is replay output, not an expectation embedded in the corpus.
