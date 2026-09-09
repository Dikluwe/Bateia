# Testes independentes da CLI de efeitos

Papel C, contrato congelado: `00_nucleo/passos/0062-verificar-efeito-executado.md`.
Criar `package cli_test` com verificador fake pela interface pública de domínio.
Argumentos inválidos e ajuda não chamam Verify; defaults e limites são exercitados.
Verificar JSON público em stdout e arquivo, inclusive status unknown, propagação
de erro e falha de escrita. Não consultar produção ou fixtures de outros papéis.
Não executar testes; integrador executa artefatos congelados. Prompt precede código.
