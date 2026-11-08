# insight_exploration_memory

**Camada:** L1. **Responsável:** `01_core/domain/insight_exploration_memory`. **Materialização:** [insight_exploration_memory.go](../../../../01_core/domain/insight_exploration_memory.go).

## Interface e responsabilidade

Declarar InsightReadObservation, InsightMemoryEntry, InsightExplorationMemory e InsightFollowup. Consumidores são validadores, execução de aprofundamento, relatório multifonte e transporte de feedback.

## Comportamento vigente

Snapshot expõe schema_version, id, generation e entries. Entrada expõe id, identity, state, first_generation, last_question e read_attempts; reason, last_status, lead e last_read são opcionais. Observação contém source/path/content_hash. Followup contém entry_id/lead_id/identity/status/cards, com receipt e bundle opcionais. Preservar tags e omissões existentes.

## Limites e dependências

Histórico declarado não é conteúdo atual: cartões/bundle ficam no followup, observação anterior fica na memória. Tipos não autenticam fontes, alteram estado ou fazem I/O. Manter a extensão opcional do relatório sem alterar schemas históricos.

As regras compartilhadas estão nas seções [interfaces e formatos](../../../contratos/exploracao-de-vizinhancas-v1.md#interfaces-e-formatos), [memória](../../../contratos/exploracao-de-vizinhancas-v1.md#memória). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [WIRE-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_exploration_memory_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__domain__insight_exploration_memory.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
