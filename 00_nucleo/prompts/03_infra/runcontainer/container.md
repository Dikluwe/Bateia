# Prompt — container durável L3

Implemente o contrato congelado do passo 0079. O pacote possui privadamente store,
fila e serviço, constrói-os nessa ordem e expõe apenas `New`, `Service`, `Recover`,
`Start` e `Shutdown`. Preserve cleanup e precedência de erros. Não leia ambiente,
não importe L2/L4, não altere componentes existentes e não inicie a fila em `New`.
