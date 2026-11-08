# Regressões de fronteiras da investigação

**Camada:** L1. **Origem ex-ante:** [passo0159](../../../passos/0159-fechar-cancelamento-escopo-e-recibos.md).

Usar APIs públicas e portas contadas em memória. Recheck: cancelamento e deadline após árvore ou extração devolvem contexto e relatório vazio; nenhuma extração após árvore cancelada; contexto nil/precancelado também sem efeitos; sem interrupção mantém VERIFIED. Recovery: target explícito LICENSE não lê outro arquivo da árvore; target derivado só lê caminhos casados; resposta extra é rejeitada; modo geral continua podendo selecionar árvore. Neighborhood: relatório válido, depois resselar alterações de recibo órfão, nil, duplicado ou status divergente, exigindo ErrIntegrity e zero leitores. Históricos não coletados e pares válidos continuam aceitos; assertivas independentes da implementação. Reusar fixtures existentes sem modificar seus consumidores históricos.
