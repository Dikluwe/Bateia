# Bridges dos provedores de repositórios

**Camada:** L1. **Contrato:** passo 0138.

Implementar bridges MarketDiscoveryConnector sobre ProjectCatalog GitHub e ObservableDiscoveryProvider Sourcegraph. Delegar efeitos; validar respostas/cardinalidade e converter metadados em ProductCandidate com identidade própria, sem inventar URL/revisão. Derivar hash de resposta sem publicar erros arbitrários. Sourcegraph máximo de uma tentativa solicitada por operação; não chamar cliente concreto. GitHub ainda usa leitor injetado separado.

## 0154 — Descoberta contextual
Bridge GitHub traduz QueryPolicy context-v2 em metadados apenas e resolução adiada; sem política preserva comportamento. Código externo continua desconhecido.
