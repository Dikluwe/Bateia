# Prompt: 01_core/domain/idea_radar_plan

Materializar somente os tipos, constantes e tags JSON públicos do planejador de fronteira
definido no passo 0067 e em `00_nucleo/benchmark/0067/contract.md`.

O artefato é L1 puro, aditivo e sombra. Reutilizar `IdeaRadarInput`, `IdeaRadarOutcome`,
`IdeaProjectionKind` e `EvidenceCoalitionRole`; não editar nem duplicar tipos 0065/0066.
Todas as slices de output serializam como arrays não nulos. Não adicionar decisão, I/O,
execução de probes, inferência textual, evidência fabricada ou integração.

O futuro arquivo espelhado é `01_core/domain/idea_radar_plan.go`, com cabeçalho de
linhagem deste prompt e camada L1. Autoridade: implementador segregado após o selo;
proibido editar contrato, baseline, oráculos, mutantes e prompts de teste.

