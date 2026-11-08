# insight_neighborhood_sources

**Camada:** L1. **Responsável:** `01_core/usecase/insight_neighborhood_sources`. **Materialização:** [insight_neighborhood_sources.go](../../../../01_core/usecase/insight_neighborhood_sources.go).

## Interface e responsabilidade

Receber relatório copiado com memória original validada e devolver origens documentais verificadas para a admissão. Consumir bundles GitHub selecionados, contextos por sítio e followups da entrada segundo o mesmo contrato de origem.

## Comportamento vigente

Confrontar identidade/namespace, SHA, caminho raiz, item/ref, bytes/hash, recibo COLLECTED/ROOT_README_ONLY e LastRead. Rejeitar duplicidades e pares coleta/bundle incoerentes. Aplicar oito documentos,128 KiB agregados e32 KiB por documento. Não confiar somente no selo externo do relatório: conteúdo e fingerprint precisam ser consistentes entre si.

## Limites e dependências

Sem consulta remota, travessia de grafo ou decisão de seleção. Origens só testemunham links já congelados; não viram cartões atuais. Erro de integridade/contrato deve chegar à admissão antes de qualquer leitor. A memória alternativa de feedback não substitui a verificação contra o histórico original.

As regras compartilhadas estão nas seções [fontes congeladas](../../../contratos/exploracao-de-vizinhancas-v1.md#fontes-congeladas). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [ADM-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [SRC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_sources_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_sources.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
