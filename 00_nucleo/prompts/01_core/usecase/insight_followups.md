# insight_followups

**Camada:** L1. **Responsável:** `01_core/usecase/insight_followups`. **Materialização:** [insight_followups.go](../../../../01_core/usecase/insight_followups.go).

## Interface e responsabilidade

Executar fila limitada por portas, compartilhada entre discover e expand, produzindo wrappers de followup, recibos, atualizações de memória e pedidos de cartões. Preservar a entrada compatível de discover.

## Comportamento vigente

Aplicar a ordem compartilhada; limitar wrappers a16 e tentativas por FollowLinks. Checar cancelamento, fonte, leitor, tentativas, documentos e bytes antes de ler, nessa ordem. Respeitar quotas globais da chamada e deadline por tentativa via runtime. Registar falha sem marcar exame; recusar bundle que exceda agregado. GitHub usa FollowReader e demais forges RepositoryReader. Não seguir descendentes recém-lidos na mesma execução.

## Limites e dependências

Efeitos somente pelas portas injetadas; nenhum cliente/rede/arquivo concreto. Memória pertence ao container, ordem ao seletor e cartões ao contexto. Expand pode ler pendência sem Lead, com lead_id vazio; discover continua exigindo trilha. Erros externos não podem vazar segredos no relatório.

As regras compartilhadas estão nas seções [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução), [memória](../../../contratos/exploracao-de-vizinhancas-v1.md#memória). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [RUN-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_followups_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_followups.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
