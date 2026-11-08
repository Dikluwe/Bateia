# Testes CLI da descoberta

**Camada:** L2. **Contrato:** passo 0137.

Portas falsas pela API pública injetada: perguntar/buscar/coletar/explorar sem input JSON, output obrigatório, limites antes de efeitos, flags duplicadas/desconhecidas/posicionais, publicação parcial, segredos e short writes. Não modificar testes antigos.

## 0157 — Controles de unicidade

Testar duplicatas booleanas implícitas/true/false/conflitantes e um/dois hífens, preservando output anterior e stdout vazio. Duplicata deve prevalecer sobre leitura de bundle inexistente ou acesso a catálogo/factory. Testar uma ocorrência true/false e ausência, além de valor inválido, sem enfraquecer consumidores anteriores. Contrato: [passo 0157](../../../passos/0157-fechar-duplicatas-e-limites-de-formulacao.md).
