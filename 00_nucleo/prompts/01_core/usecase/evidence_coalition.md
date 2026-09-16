# Prompt: 01_core/usecase/evidence_coalition

Implementar a função pura
`AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput) domain.EvidenceCoalitionReport`
exatamente segundo o contrato 0065 congelado. Normalizar, validar, deduplicar por
identidade física completa, classificar com a precedência pública e ordenar todo output
canonicamente. Preservar `Unknown`, conflito e `BaselineStatus`; impacto sempre `NONE` e
zero mudanças de veredito.

Respeitar os limites 128 regras, 128 claims, 64 testemunhas por claim, 4096 testemunhas
totais, quatro papéis por regra e 8192 operações contratuais. Usar índices e passagens
limitadas; não enumerar combinações/subconjuntos. Nenhum I/O, relógio, aleatoriedade,
estado global, biblioteca externa ou integração com avaliação existente.

O futuro arquivo espelhado é `01_core/usecase/evidence_coalition.go`, com cabeçalho de
linhagem deste prompt e camada L1. Autoridade do implementador após selo; proibido ler ou
editar oráculos, mutantes, baseline e prompts de teste.

