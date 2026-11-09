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

## Passo 0088 — composição verificável

Substituir seleção gulosa por candidatas de composição. Cobertura complementar sem conexões admissíveis e ensaio executado do conjunto exato retorna INCONCLUSIVE, nunca COMPOSE. Incompatibilidades e truncamento são explícitos; não transformar ausência de composição encontrada em BUILD quando o limite foi atingido.

## 0166 — Limiares finitos

Aplicar DEC-01 do [contrato de integridade](../../../contratos/integridade-de-entradas-e-replay-v1.md#dec-01--política-finita-e-hash-serializável), conforme [passo0166](../../../passos/0166-validar-bytes-testemunhos-e-politicas.md). Rejeitar NaN/infinitos em ambos os limiares, mantendo (0,1] para finitos; tratar falha ao serializar política, sem relatório útil/hash vazio. Preservar hashes, regras e avaliações válidas e registros históricos.

## Declarações originais de partes — 0171

Transportar por estado privado somente de leitura as partes originais até a guarda CMP-02 da composição, antes de a sanitização apagar diferenças declaradas. Não mutar inputs nem reautorizar evidência/capacidade removida: a enumeração e decisão continuam usando partes sanitizadas e os gates anteriores. Nenhuma API pública, schema, identidade ou hash de controle válido muda; declarações originais válidas/admissíveis contraditórias do mesmo ID retornam erro em vez de seleção por ordem.
