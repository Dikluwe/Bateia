# JSON decisório limitado e sem membros ambíguos

Validar um único valor JSON, UTF-8 válido, até 16 MiB, profundidade 64 e cem
mil tokens antes da materialização tipada. Rejeitar membros duplicados no mesmo
objeto, inclusive escapes e aliases de casefold Unicode reconhecidos pelo
decodificador Go. Objetos independentes podem repetir nomes. Não materializar
árvore genérica nem executar IO. Erros não incluem conteúdo ou nomes de campos.
É uma política deliberadamente mais restritiva que JSON, não autenticação.
