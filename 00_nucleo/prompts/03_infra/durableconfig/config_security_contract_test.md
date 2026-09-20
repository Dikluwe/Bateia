# Contrato estrutural da configuração durável local

Derive este oráculo exclusivamente do passo 0077. Audite as fontes Go do pacote por
AST, sem depender de nomes privados, e exija os seguintes controles observáveis na
estrutura do código:

1. a abertura do cursor combina uma operação de abertura com
   `syscall.O_NOFOLLOW`;
2. todo `Chmod(0600)` aplicável ao cursor aberto é lexicalmente precedido, na mesma
   função, por validação de `Nlink` exatamente igual ou diferente de um;
3. uma chamada ordinária a `Read` valida tanto a contagem quanto o erro retornados;
   primitivas que já garantem leitura completa, como `io.ReadFull`, permanecem
   alternativas válidas.

Inclua controles sintéticos que preservem a forma válida e rejeitem isoladamente a
remoção de `O_NOFOLLOW`, a inversão entre `Chmod` e a validação de `Nlink`, e o
descarte da contagem retornada por `Read`. Ausência de fonte produtiva, construção
AST não reconhecida ou capacidade POSIX indisponível não pode ser convertida
silenciosamente em sucesso.
