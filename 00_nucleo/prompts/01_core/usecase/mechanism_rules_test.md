# Verificação independente do catálogo ativo de mecanismos

Derive os testes somente do contrato do passo 0078. Verifique bytes, tamanho, hash,
JSON v3, chamadas repetidas, clone defensivo, mutação isolada e concorrência sob o
race detector. Confirme também que o accessor não expõe storage compartilhado.

Não adapte expectativas à implementação candidata e não leia seu conteúdo durante
a autoria do oráculo.
