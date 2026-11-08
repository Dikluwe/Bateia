# Testes de execução de descoberta

**Camada:** L1. **Contrato:** passo 0137.

Portas puras simulam sucesso/empty/falha/rate-limit/payload inválido/cancelamento e runtime explícito. Comprovar que direta bem-sucedida não suprime outras rotas, deadlines são propagados, falhas conservam candidatos, sem retry/novas chamadas após cancelamento e sem mensagens arbitrárias vazando. Conferir bundles/contexto final e IDs. Sem rede/disco/relógio concreto.

Verificar que serializar/reler os bundles coletados e entregá-los a ExploreInsights reproduz a mesma identidade e conteúdo do contexto integrado, inclusive com README contendo HTML. Diferença de representação JSON não deve invalidar esse replay do novo fluxo.
