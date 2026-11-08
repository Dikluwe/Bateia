# Coleta e serialização de evidências v1

## Origem e autoridade

Contrato ex-ante do [passo0158](../passos/0158-validar-origem-e-estabilizar-bundles.md), sob [ADR0022](../ADR/ADR-0022-integridade-de-coleta-e-serializacao-estavel.md) e regras Tekt. O sufixo v1 identifica este contrato, sem introduzir outro schema público. L1 domínio possui serialização; L1 casos de uso possuem admissão e vínculo de origem; os efeitos concretos permanecem em L3.

## Admissão e pin

A coleta coordenada valida todos os planos antes de executar qualquer porta. Exige schema1, identidade de repositório válida e igual à do candidato, quotas não negativas e orçamento de bytes positivo. A API genérica compartilha essas regras e conserva somente a exceção de schema ausente para consumidores legados; versões desconhecidas são rejeitadas. Revisões abreviadas legadas da API genérica não são promovidas a SHA validado pela jornada de insights.

Plano com commit fixado governa o commit das leituras. Metadados de branch mais recente não substituem esse pin. HEAD ou vazio usam a resolução explícita existente. Owner/repo divergente nos metadados, resposta de extração nil ou Ref diferente da identidade e revisão efetivas retornam ErrIntegrity antes de admitir conteúdo. A coordenação propaga esse erro sem bundles/cartões; não o reclassifica como falha operacional de resolução. Falhas operacionais de resolução continuam explícitas, sem promover histórico ou revisão opaca a fonte lida. Cancelamento observado após uma porta prevalece antes da próxima.

## Serialização e histórico

Novas serializações de EvidenceBundle emitem requirement_coverage como array não nulo, inclusive quando nil ou vazio. O marshaler copia o estado necessário e não altera slices do chamador. Após roundtrips, bytes e identidade causal devem estabilizar; null legado continua legível e é projetado como array em nova serialização.

Registros e hashes históricos não são reescritos. Formas vazias antes alternantes convergem nas novas projeções, que podem diferir do ID antigo instável. Não substituir IDs já emitidos nem prometer igualdade com os dois resultados defeituosos. Cobertura populada, já estável, conserva bytes e IDs sob controles congelados antes do patch. Reprodução histórica continua no commit registrado pelo gate; hashes de arquivos preservados são conferidos separadamente da nova execução.

## Aceitação

| ID | Observável | Consumidor |
| --- | --- | --- |
| EVD-01 | Owner/repo/commit divergentes ou resposta nil não geram conteúdo; falha é de integridade | [Regressões de coleta](../../01_core/usecase/evidence_integrity_regression_test.go) |
| EVD-02 | Plano inválido em qualquer posição impede todas as portas; schema ausente é exclusivo da API legada | [Admissão do lote](../../01_core/usecase/evidence_integrity_regression_test.go) |
| EVD-03 | Pin fixado permanece; HEAD resolve explicitamente; falha operacional continua visível | [Metadados e compatibilidade](../../01_core/usecase/evidence_integrity_regression_test.go) |
| EVD-04 | Nil/vazio/null legado estabilizam como array; exemplo populado conserva hash literal | [Regressões de domínio](../../01_core/domain/evidence_integrity_regression_test.go) |
| EVD-05 | Formulação salva/reaberta mantém IDs; controle populado conserva ID literal anterior | [Regressões de formulação](../../01_core/usecase/evidence_integrity_regression_test.go) |

Esses controles verificam integridade declarada e estabilidade local, sem autenticar provedor remoto, comprovar comportamento, novidade ou efeito no alvo.
