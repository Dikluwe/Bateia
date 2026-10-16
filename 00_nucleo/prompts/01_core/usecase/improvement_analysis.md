# Análise pura de hipóteses de melhoria

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/improvement_analysis.go`  
**Contrato:** passo 0103 e ADR-0008, recorte A1-INV.

Separar do runtime a transformação de um pedido normalizado e de mecanismos previamente resolvidos/validados em comparação, hipóteses, lacunas, planos e proveniência. A base contém somente fatos de mecanismos e sinal de degradação. O resultado nomeia as cinco projeções e seu estado de negócio. Não receber store, fila, serviço de runs ou adaptadores; não usar relógio, contexto operacional, disco, rede ou execução.

Preservar exatamente o algoritmo do contrato 0028: derivação por regras ou proposta externa, IDs causais, ordenação, defaults conservadores, limites de hipóteses, riscos, unknowns e planos. Fonte degradada implica hipótese inconclusiva. Licença desconhecida permanece desconhecida; estado da síntese não autoriza adoção. Preservar as projeções legadas byte a byte no caso de caracterização congelado antes da extração.

Normalização, hash do pedido, validação de referências, resolução e integridade de fontes, admissão, cancelamento e publicação continuam sob o serviço de execução. Esta fatia não introduz nova API pública nem muda schemas.
