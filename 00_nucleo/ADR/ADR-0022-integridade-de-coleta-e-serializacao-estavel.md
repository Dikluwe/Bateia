# ADR-0022 — Integridade da coleta e serialização estável de bundles

**Status:** ACEITO PARA O RECORTE 0158
**Data:** 2026-11-08

## Contexto

A investigação depende de identidade repo@commit ligada aos bytes. A coleta genérica não confrontava a identidade da resposta de extração e a coleta coordenada não admitia schemas dos planos. O marshaler de coverage vazio alternava array/null, violando o contrato0072 e tornando a identidade causal dependente de quantas vezes o pedido era salvo e reaberto.

## Decisão

L1 admite todos os planos coordenados antes de efeitos e compartilha a validação com a API genérica. Schema1 é obrigatório na coordenação; schema ausente é compatibilidade exclusiva da API legada. Versões desconhecidas são erro. Plano fixado conserva seu pin mesmo quando metadados de branch avançam; identidade divergente de metadados/extração ou extração nil é erro de integridade. Falhas de integridade não são falhas operacionais recuperáveis nem geram contexto.

O domínio emite array de cobertura não nulo em novas serializações, inclusive vazio. Não introduzir outro schema para reparar essa violação do contrato existente. Arrays populados e os demais campos mantêm suas representações. Testes congelam exemplos estáveis antes do patch e testam roundtrips de todas as representações vazias admitidas.

## Compatibilidade e limites

Arquivos históricos, seus IDs e hashes permanecem intocados. Os dois resultados alternantes de uma representação defeituosa não podem ser preservados simultaneamente como um único ID estável. Novas projeções de entradas com cobertura vazia usam a representação corrigida e podem divergir daquele ID antigo; não substituir IDs de registros anteriores nem alegar replay byte a byte dessas entradas sob a árvore nova. Gates selados conservam o commit original de reprodução. Formatos com cobertura populada, já estáveis, são controles de preservação próprios.

Validar identidade declarada e hash não autentica o provedor remoto, não comprova comportamento nem autoriza adoção. Sem mudanças de ownership: formato em domínio L1, regras/admissão em casos de uso L1, transporte L2 e efeitos L3. O [passo0158](../passos/0158-validar-origem-e-estabilizar-bundles.md) declara os observáveis e a verificação.
