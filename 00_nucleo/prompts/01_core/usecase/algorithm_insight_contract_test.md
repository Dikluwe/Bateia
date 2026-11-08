# Consumidores adicionais de integridade e limites da formulação

**Camada:** L1. **Proprietário:** algorithm_insight_contract_test.go. **Origem:** intenção ex-ante no [passo 0157](../../../passos/0157-fechar-duplicatas-e-limites-de-formulacao.md), ampliando os consumidores do [contrato de formulação](algorithm_insight.md) sem reescrever os testes 0155.

Testar FormulateAlgorithmInsights com bundles fictícios de fonte/alvo distintos, referências literais independentes dos helpers internos e hashes SHA256 calculados nos testes. Permutações de referências/bundles/ideias conservam o relatório inteiro e a entrada; mudança nos bytes efetivamente citados, com hashes/citações válidos, muda IDs e InputSHA256. Rejeitar referências duplicadas em ambos os papéis.

Exigir aceitação no limite e rejeição no primeiro excesso para títulos200, princípio1500, campos de decisão/observável/rejeição1000, pergunta1000 e cada passo/condição600 bytes; limites são de bytes, incluindo UTF-8. Steps2..12, condições1..12, fontes/ideias/referências1..8. Cobrir campos vazios/controles e schema inválido. Pedido até4 MiB, cada grupo de exploração até2 MiB e512 itens. Faixa de citação até20 linhas e4096 bytes, excerpt literal até2048 bytes, sem referências vazias, fora da faixa ou com duplicatas. Rejeições devolvem erro e relatório sem ID; sucessos mantêm DECLARED_NOT_REASSESSED/UNVERIFIED.

Usar tamanhos, bytes e resultados esperados explícitos. Não calcular o ID esperado chamando a função de produção que o gera. Não executar I/O ou fontes, não introduzir nomes/algoritmos do ensaio XERJ e não alterar políticas produtivas de formulação para fazer um teste passar.
