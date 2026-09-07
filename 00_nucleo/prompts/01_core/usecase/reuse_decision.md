# Reuse decision engine

Implemente um decisor offline, determinístico e sem LLM sobre IdeaSpec,
CandidateAssessment e política versionada. Avalie na ordem FORK, EXTEND,
COMPOSE, INTEGRATE, INCONCLUSIVE, REFERENCE e BUILD. Falha operacional e
ambiguidade nunca podem virar BUILD. Ordene entradas antes de decidir, preserve
os assessments e exponha regras, evidências, bloqueios, lacunas e esforço apenas
por fatores observáveis.

## Passo 0040

Aplicar o contrato de autoridade, sinais lexicais sem comprovação comportamental, procedência local, seleção de produção, negação conservadora e lacunas acionáveis de `00_nucleo/passos/0040-corrigir-autoridade-da-evidencia.md`. Preservar determinismo e ausência de I/O no núcleo.

## Revisão 0043 — autoridade e ligação
Aplicar o contrato comum de evidence_authority antes de autorizar capacidade ou operação. Evidência lexical/documental não ganha PROVEN. Negação deve conservar polaridade. Na extração de mecanismos aceitar SHA-256 hexadecimal produzido pela avaliação e formato prefixado legado, verificando conteúdo e projeções de tipo file→symbol e test→test_implementation. Preservar IDs reais: catálogo sem regra compatível retorna lacuna explícita, nunca traduz nomes em prova.
