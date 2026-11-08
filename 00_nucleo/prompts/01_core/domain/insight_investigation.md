# Contrato da investigação de insights

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/insight_investigation.go`  
**Contrato:** passo 0105

Declarar pedido independente de IdeaSpec, política/decisão de reúso e run_id: pergunta, alvo com contexto/propriedade/efeito desejado/restrições, classes de mecanismo, fontes (bundle e avaliações declaradas) e teto de hipóteses. Declarar relatório com identidade/hash, proveniência, mecanismos, hipóteses, lacunas e limites. Distinguir escopo de evidências congeladas, avaliações fornecidas e aplicabilidade ainda desconhecida. Sem I/O.

0134: acrescentar `observations,omitempty` por fonte usando tipos proprietários de source_insight.go. Ausência/vazio preserva a serialização causal v1. A nova projeção tem tipos e schema próprios; não mudar campos da saída v1.
