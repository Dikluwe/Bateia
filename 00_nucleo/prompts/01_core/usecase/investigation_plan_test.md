# Plano independente 0043
Exigir ligações de cada consulta a IDs reais dos requisitos declarados, preservando consulta original, efetiva e requisito sob fallback. Covers explicitamente identifica o requisito. Busca vazia produz lacuna acionável, nunca inventa candidato. Testes somente APIs públicas e dados literais sem I/O.

## Controle da emissão ao provedor
Spy público em memória captura DiscoveryRequest antes da descoberta; plano com Covers capabilities[0] deve enviar RequirementIDs correspondentes ao texto original, preservar ideia de entrada e guardar tentativa efetiva de fallback. A captura é sincronizada para comportar concorrência sem I/O nem data race.
