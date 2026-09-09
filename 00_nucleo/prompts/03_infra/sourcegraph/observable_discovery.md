# Prompt Spec L0: Descoberta Observável no Sourcegraph

## Objetivo

Expor as queries completas e fallbacks realmente tentados pelo Sourcegraph, sem
vazar credenciais e preservando o contrato de busca existente.

## Restrições

- L3, adaptador da porta observável de L1.
- Registrar query, quantidade e erro resumido por tentativa.
- Distinguir resposta vazia de erro.
- Parar no primeiro resultado não vazio.
- Nunca serializar token, cabeçalhos ou corpo sensível.
## Extensão operacional do Passo 0013

Cada fallback deve respeitar o timeout recebido, registrar início, duração,
status e classe segura de erro, limitar tentativas e parar ao cancelar o contexto.
Nunca serialize token e não converta timeout em tentativa vazia.

## Passo 0054
Aplicar o contrato de `00_nucleo/passos/0054-preservar-intencao-na-descoberta.md`: consulta integral única, sem fallback destrutivo, erros explícitos e versão de evidência nova para invalidar cache anterior.

## Passo 0055 — hipótese rastreável

Aplicar o contrato congelado 00_nucleo/passos/0055-planejar-consultas-de-codigo.md, que substitui a tentativa única por original seguido de no máximo uma hipótese apenas após vazio válido. QueryPlan, RawCandidates e RejectedCandidates registram somente a derivada; contagens finais são admitidos. Plano preserva pergunta/restrições não verificadas. Ambas portas usam catálogo puro; linguagem e escopo regex literal persistem. Admitir somente FileMatch com caminho e ação+objeto em snippets efetivamente retidos do mesmo arquivo; reter todos os previews limitados individualmente a 240 runes, sem descrição/nome. Não combinar arquivos. Preservar limites, erros, timeouts, cancelamento e rastros nos candidatos. Validar envelope GraphQL antes de considerá-lo vazio. EvidenceVersion muda para evidence-v15-planned-discovery.

Adendo 0055: aplicar 00_nucleo/passos/0055-adendo-evidencia.md. SearchAttempt.AdmissionEvidence ([]ProviderCandidateMatch, JSON admission_evidence,omitempty) somente na derivada, preserva matches admitidos e os snippets efetivamente examinados da porta Discover, também nos rastros anexados aos candidatos. Não publicar rejeitados como evidência de admissão.

## Revisão de sintaxe 0055-v2
Aplicar 00_nucleo/passos/0055-adendo-sintaxe.md: emitir `patterntype:regexp (a|b) AND (x|y)`, sem aspas/content. O catálogo, filtros, bloqueios e admissão permanecem. Versão do plano 0055-v2 e EvidenceVersion evidence-v16-planned-discovery-syntax. Esta revisão substitui a serialização anterior do contrato.

## Passo 0056 — cobertura observada

Aplicar o contrato congelado `00_nucleo/passos/0056-explicar-cobertura-e-perdas.md`. Somente na derivada: count limitado20..200 sem overflow, coleta de até N repositórios distintos na ordem da resposta e evidências posteriores desses mesmos repositórios, admissão lexical antes do limite final. Cobertura unknown, matchCount nullable, contagens observadas e diagnósticos com motivos, no máximo32 previews retidos de240 runes normalizadas por repositório; omissões e truncamentos explícitos. Diagnósticos usam somente evidência de arquivo retida, exigindo ação e objeto no mesmo arquivo. Selected/insufficient_evidence/candidate_limit distinguem admissão e seleção. Propagar Coverage e CandidateDiagnostics nas tentativas e rastros; Matches/AdmissionEvidence somente selecionados. Literal permanece compatível. EvidenceVersion evidence-v17-discovery-coverage.
