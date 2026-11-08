# Testes CLI da exploração

**Camada:** L2. **Contrato:** passo 0136.

Testar superfície pública com pergunta/bundles coletados, familiar repetível e ausência de bundles; arquivo/stdout e ausência de conclusões semânticas. JSON estrito, flags/posições inválidas, erro de integridade preservando output, alias, segredo nos trechos, short writes e nenhum input alterado. Sem chamadas remotas.

## 0157 — Controles de unicidade

Testar duplicatas booleanas implícitas/true/false/conflitantes e um/dois hífens, preservando output anterior e stdout vazio. Duplicata deve prevalecer sobre leitura de bundle inexistente ou acesso a catálogo/factory. Testar uma ocorrência true/false e ausência, além de valor inválido, sem enfraquecer consumidores anteriores. Contrato: [passo 0157](../../../passos/0157-fechar-duplicatas-e-limites-de-formulacao.md).
