# Prompt: 01_core/domain/idea_radar

Materializar somente os tipos, constantes e tags JSON públicos do radar de convergência
definidos no passo 0066 e em `00_nucleo/benchmark/0066/contract.md`.

O artefato é L1 puro e sombra. Não editar `idea_spec.go`, `evidence_coalition.go` ou
qualquer schema existente. Reutilizar os tipos públicos `EvidenceOriginStatus` e
`EvidenceCoalitionRole`; não criar equivalentes concorrentes. Todas as slices do
relatório devem serializar como arrays não nulos. Não incluir decisão, I/O, inferência,
LLM, embeddings ou integração.

O futuro arquivo espelhado é `01_core/domain/idea_radar.go`, com cabeçalho de linhagem
deste prompt e camada L1. Autoridade: implementador segregado após o selo; proibido
editar contrato, baseline, oráculos, mutantes e prompts de teste.
