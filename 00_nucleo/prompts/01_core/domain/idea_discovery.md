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
