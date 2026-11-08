# insight_neighborhood

**Camada:** L2. **Responsável:** `02_shell/cli/insight_neighborhood`. **Materialização:** [insight_neighborhood.go](../../../../02_shell/cli/insight_neighborhood.go).

## Interface e responsabilidade

RunInsightNeighborhoodCommand transporta idea insights expand entre flags/JSON e fachada L1, recebendo factory/host injetados. Input/output são obrigatórios; memory/focus/all-trails e overrides são opcionais.

## Comportamento vigente

Rejeitar flags duplicadas, argumentos posicionais, pergunta/fontes novas, valores vazios e limites explicitamente inválidos. FollowLinks0 é válido; max-passages/tempos0 explícitos são inválidos. Ler JSON estrito, proteger output contra input/memory e chamar ValidateInsightNeighborhood antes da factory. Encaminhar AllTrails/foco sem calcular herança localmente. Publicar relatório com ID em falha operacional antes de devolver erro sanitizado.

## Limites e dependências

Política de foco, aliases, memória e quotas pertence a L1. L2 não importa adaptadores L3, não chama rede concreta nem mantém segundo arquivo global de estado. Usar host/publicação protegida existente; barrar vazamento de segredos, sem repetir seu valor na mensagem.

As regras compartilhadas estão nas seções [transporte cli](../../../contratos/exploracao-de-vizinhancas-v1.md#transporte-cli). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [CLI-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [CLI-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__02_shell__cli__insight_neighborhood.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
