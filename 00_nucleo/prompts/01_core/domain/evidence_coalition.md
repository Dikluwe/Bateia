# Prompt: 01_core/domain/evidence_coalition

Materializar somente os tipos, constantes e tags JSON públicos do analisador de
coalizões definidos pelo contrato congelado
`00_nucleo/passos/0065-analisar-coalizoes-de-evidencia-em-modo-sombra.md` e por
`00_nucleo/benchmark/0065/contract.md`.

O artefato é derivado e em modo sombra. `BaselineStatus` deve permanecer um
`EvidenceStatus`; não criar novo veredito nem acrescentar campos a schemas existentes.
Todas as slices do relatório serializam como arrays não nulos. Não incluir lógica de
decisão, I/O, dependências externas ou integração com busca/CLI. O futuro arquivo
espelhado é `01_core/domain/evidence_coalition.go`, com cabeçalho de linhagem deste
prompt e camada L1.

Autoridade: implementador segregado, depois do selo. Não editar contrato, baseline,
oráculos, mutantes ou prompts de teste.

