# Prompt: corpus adversarial local

Camada L3. Gerar `03_infra/local/security_test.go`.

Reproduzir com dados sintéticos hardlink externo que Root permite mas política recusa em extração, manifesto e busca de implementação; inventário rotula symlink/hardlink/FIFO honestamente. Preservar alias interno positivo existente. Testar 100 arquivos, 8 MiB e exato orçamento positivo, árvores gigantes com erro explícito, metadata HEAD detached/loose/packed.

PATH shim e tokens/Git config fictícios não devem executar comandos durante extração ou metadata. `.git` file externo e HEAD em hardlink deixam BaseCommit vazio com explicação explícita e não resolvem candidato. Nenhum segredo/rede real usado; sem atualização de certificados históricos. Contrato 0072 preserva 'não truncar silenciosamente' mas passa a exigir erro quando 20k é excedido.
