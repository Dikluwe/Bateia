# Reuse decision engine

Implemente um decisor offline, determinístico e sem LLM sobre IdeaSpec,
CandidateAssessment e política versionada. Avalie na ordem FORK, EXTEND,
COMPOSE, INTEGRATE, INCONCLUSIVE, REFERENCE e BUILD. Falha operacional e
ambiguidade nunca podem virar BUILD. Ordene entradas antes de decidir, preserve
os assessments e exponha regras, evidências, bloqueios, lacunas e esforço apenas
por fatores observáveis.
