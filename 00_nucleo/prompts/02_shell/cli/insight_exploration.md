# CLI de exploração para IA

**Camada:** L2. **Contrato:** passo 0136.

Expor flags diretas question/bundles/familiar repetível/max-passages/output, recusar posição extra, flags desconhecidas/repetidas escalares; ler lista de EvidenceBundle estrita usando helper existente. Não montar avaliações/observações nem classificar fontes. Proteger saída contra alias do input e usar escrita JSON protegida/segredos/short writes. Delegar política a ExploreInsights. Sem L3 concreto, busca ou modelo na casca.

## 0153
Adaptar flag --review-alternatives booleana ao pedido L1; sem flag preservar comportamento.

## 0155 — Ideias de algoritmos
Flag --algorithm-ideas booleana adapta campo L1, preservando --review-alternatives e modo padrão.

## 0157 — Unicidade booleana

Consumir o [passo 0157](../../../passos/0157-fechar-duplicatas-e-limites-de-formulacao.md): a unicidade de flags escalares inclui booleanas, com valores implícitos ou explícitos, mesmo false e mesmo com um hífen. Rejeitar repetição no parse antes de efeitos/publicação; familiar permanece repetível. Uma ocorrência admite as formas booleanas válidas existentes.
