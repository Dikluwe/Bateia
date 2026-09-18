# Teste do passo 0075 para `01_core/usecase/evidence.go`

Derivar testes da obrigação em
`00_nucleo/passos/0075-localizar-expressoes-probatorias-compostas.md`.

Exigir que padrões fortes compostos por identificadores e operadores, como
`&__ht == this`, sejam localizados na linha exata de código e gerem somente evidência
lexical parcial. Cobrir variação de espaços, atribuição a membro e módulo. Rejeitar
comentários, strings, substrings de identificador, texto natural sem operador e
expressões divididas entre linhas. Confirmar que padrões probatórios não entram nos
termos de recuperação de caminho, que continuam derivados apenas de `TargetHints`.

Não alterar catálogo, pesos, confiança, orçamento ou política de Unknown para fazer o
teste passar.

