# evidence_checkpoint_test

**Camada:** L1. **Contrato:** passo 0130 e ADR-0013.

Matriz de checkpoints em memória: estados, schema, hashes, oito versões, presença/ausência/mudança de modelo, policy/target alterados, dois artefatos obrigatórios, erro/hash divergente, outputs incompletos/vazios/duplicados, ordem e ausência de leituras antes de admitir metadados. Preservar a regra da primeira ocorrência verificável; leitor nil recusa quando uma leitura seria necessária.
