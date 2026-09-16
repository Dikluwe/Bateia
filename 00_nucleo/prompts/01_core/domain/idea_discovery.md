# Prompt proprietário — `01_core/domain/idea_discovery.go`

Origem L0 exclusiva de `01_core/domain/idea_discovery.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/idea_discovery.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Materializar somente a responsabilidade desta unidade, preservando os contratos observáveis: `DiscoveryRequest`, `ProviderAttempt`, `ProviderDiscovery`, `DiscoveredCandidate`, `ExpansionDiscovery`, `DiscoverySnapshot`, `RelevanceLabel`, `CandidateJudgment`, `DiscoveryJudgments`, `DiscoveryMetrics`, `ExpansionDiscoveryEvaluation`, `DiscoveryEvaluation`. Manter a direção de dependências da camada e não incorporar responsabilidades pertencentes aos demais consumidores do prompt de origem.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/idea_discovery.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
## Extensão operacional do Passo 0013

Modele política de deadline total e por tentativa, estados `COMPLETE`, `PARTIAL`
e `FAILED`, causas de término, receipts seguros e eventos de progresso. Campos
novos do snapshot devem ser opcionais para preservar replays legados. Timeout,
cancelamento, falha, vazio e `NOT_STARTED` nunca podem ser equivalentes.

## Contrato 0043 — consultas ligadas à intenção

Antes da chamada ao provedor, cada expansão recebe RequirementIDs derivados dos Covers da IdeaSpec original, com IDs estáveis por conteúdo e tipo. Enviar os IDs na DiscoveryRequest e preservá-los na expansão do snapshot. Sem vínculo válido não fabricar associação: lista vazia exige investigação. Estes IDs descrevem a pergunta, não evidência do candidato.

## Passo 0055 — hipótese rastreável

Aplicar o contrato congelado 00_nucleo/passos/0055-planejar-consultas-de-codigo.md, que substitui a tentativa única por original seguido de no máximo uma hipótese apenas após vazio válido. QueryPlan, RawCandidates e RejectedCandidates registram somente a derivada; contagens finais são admitidos. Plano preserva pergunta/restrições não verificadas. Ambas portas usam catálogo puro; linguagem e escopo regex literal persistem. Admitir somente FileMatch com caminho e ação+objeto em snippets efetivamente retidos do mesmo arquivo; reter todos os previews limitados individualmente a 240 runes, sem descrição/nome. Não combinar arquivos. Preservar limites, erros, timeouts, cancelamento e rastros nos candidatos. Validar envelope GraphQL antes de considerá-lo vazio. EvidenceVersion muda para evidence-v15-planned-discovery.

## Passo 0056 — cobertura observada

Aplicar o contrato congelado `00_nucleo/passos/0056-explicar-cobertura-e-perdas.md`. Somente na derivada: count limitado20..200 sem overflow, coleta de até N repositórios distintos na ordem da resposta e evidências posteriores desses mesmos repositórios, admissão lexical antes do limite final. Cobertura unknown, matchCount nullable, contagens observadas e diagnósticos com motivos, no máximo32 previews retidos de240 runes normalizadas por repositório; omissões e truncamentos explícitos. Diagnósticos usam somente evidência de arquivo retida, exigindo ação e objeto no mesmo arquivo. Selected/insufficient_evidence/candidate_limit distinguem admissão e seleção. Propagar Coverage e CandidateDiagnostics nas tentativas e rastros; Matches/AdmissionEvidence somente selecionados. Literal permanece compatível. EvidenceVersion evidence-v17-discovery-coverage.

## Passo 0068 — proveniência técnica

Adicionar somente a `DiscoveryRequest` o campo opcional `Derivation
*ExpansionDerivation` com JSON `derivation,omitempty`, para que o adaptador diferencie
propostas técnicas finitas de consultas literais e hipóteses 0055. Não duplicar ground
truth, não tornar derivação obrigatória para chamadas legadas e não alterar schemas de
relatório persistido.
