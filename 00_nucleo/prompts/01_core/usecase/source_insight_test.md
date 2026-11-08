# Testes da análise localizada de insights

**Camada:** L1. **Contrato:** passo 0134.

Testar referências exatas e hashes com/sem prefixo, intervalos inválidos/limites e classes/tipos/textos inválidos, duplicatas e budgets. Demonstrar determinismo por reordenação de fontes/observações/referências e imutabilidade; IDs mudam por alteração de declaração/tipo/referência e nenhuma identidade v1 é promovida. Cobertura não aparece como pré-condição; declaração sem mecanismo não gera hipótese; PARTIAL/conflito/incompleto não são promovidos por anotações. Duas fontes/classes não misturam hipóteses ou referências. Sem observações produzir lacunas e precondições vazias. Cancelamento e esquema incorreto rejeitados. Reutilizar fixture pura anterior sem alterar testes congelados.
