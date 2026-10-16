# 0124 — Pesquisa do Laya com o Bateia atual

Data: 2026-10-07. O binário recompilado da árvore atual executou `idea evidence collect` e `idea insights analyze`. Coleta pública sem tokens, em diretório temporário sem .env; Laya fixado em 3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63. O produto não foi alterado.

## Resultado automático

O coletor preservou 24 itens: 23 alvos e metadata. Os seis alvos explícitos foram coletados; 540 alvos da expansão ficaram SKIPPED_BUDGET. Bundle parcial não implica falha dos seis alvos nem cobertura total do projeto. Hashes de todos os conteúdos conferidos.

Insights produziu INCONCLUSIVE, zero mecanismos e zero hipóteses, com NO_COMPATIBLE_EVIDENCE_OBSERVED. O catálogo atual cobre log, fila, armazenamento por conteúdo e isolamento por portas; não inclui classificação semântica ou calibração. A entrada declarou avaliações vazias, portanto a saída também demonstra dependência de avaliações compatíveis fornecidas, e não incapacidade universal de compreender Python. Não se fabricaram avaliações para induzir uma resposta. Fontes coletadas e proveniência são úteis; descoberta autônoma dessas melhorias ainda não foi demonstrada.

## Controles exploratórios do modelo

Pergunta e pesos de 0123 congelados; seis execuções correlacionadas, sem treino e sem truncamento de estado. Instrução/opções iguais às auditadas em 0123. Nenhum caso é teste independente.

| Entrada | Escolha | answer_confidence |
| --- | --- | --- |
| Original conc | undetermined | 0,6194 |
| Alegação sobre a família With* | undetermined | 0,6555 |
| Fonte sintética nomeando WithMaxGoroutines | undetermined | 0,3913 |
| Original com vínculo familiar explícito do autor | supports | 0,7870 |
| Fonte sintética insuficiente | undetermined | 0,7158 |
| Repetição original | undetermined | 0,6194 |

A referência autoral dos quatro casos documentais conclusivos é contradicts. A fonte sintética e o vínculo familiar são controles de autoria, não citações upstream. A mudança para supports não é melhoria; mostra sensibilidade ao contexto neste caso. A repetição conserva a resposta. A hipótese de falta de vínculo familiar não foi sustentada como correção suficiente. A causa interna permanece desconhecida. Confiança não calibrada no domínio não prova correção. Não escolher prompts ou limiares para decorar estes resultados.

## Achados da análise de bancada

Estes achados vieram de inspeção e controles executados pelo agente usando os bytes coletados, e não da síntese automática do Bateia.

1. **Orçamento da pergunta:** laya/common.py, build_head (linhas 225–276), corta instruções e descrições de opções. build_sequence (linhas 197–217) reporta truncamento do estado. Com instrução autoral longa, 944 tokens viraram 139, mas os três tokens do estado permaneceram inteiros e truncated=false. Comportamento permitido pelo código, não defeito comprovado; oportunidade: expor tokens originais/usados de instruções e opções, com política explícita de rejeição. Os pedidos originais de 0123 não sofreram esse corte.
2. **Separação de calibração:** laya/train.py, split_calibration (474–488), embaralha índices, sem usar família de projeto. Executou-se somente a função extraída por AST dos bytes coletados, com dez itens autorais de uma família: oito treino e dois calibração, com família comum. Não é treino executado nem demonstra vazamento de um dataset upstream. Oportunidade: dividir externamente por família, desativar split interno e reservar calibração independente; ampliar API para grupos apenas com contrato próprio.
3. **Contexto:** vínculo familiar explícito mudou a classificação para suporte no controle. Oportunidade investigativa: ampliar desafios revisados de proibição, condição temporal e exceções, com teste final reservado. Ainda não há intervenção que melhore a decisão comprovadamente.

## Continuação recomendada

Priorizar auditoria de orçamento da entrada e preparação de dados por família, pois são verificáveis deterministicamente. Avaliar qualidade com referências revisadas antes de treino. Para o Bateia, estudar contrato próprio de observações sobre entradas de modelos e protocolos de avaliação; ampliar o catálogo só após exemplos positivos/negativos e consumidores definidos. Preservar coleta, análise manual, diagnóstico e efeito confirmado como escopos distintos.

## Evidência e verificação

[Bancada e reprodução](../../lab/bancada/decisoes/0124/README.md). Resultados brutos, hashes, comandos e recibos ficam no Lab. Lint, testes do módulo raiz, vet e build verificados; build substitui /dev/shm por /tmp/bateia-product-0124. Fixtures usam assinatura Git desativada somente no processo. Gates selados não reproduzidos. Nenhum treino, limiar ajustado, integração produtiva, publicação ou execução de fontes Go de candidatos.
