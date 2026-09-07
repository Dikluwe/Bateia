# Testes independentes locais 0043
Autoria: agente testador separado do implementador, acesso declarado ao contrato congelado 0043, AGENTS e tipos públicos. Isolamento físico não atestado: filesystem compartilhado.

Antes do código: criar diretórios temporários independentes e testemunhos sentinela. Manifesto e README finais externos e arquivo sob diretório symlink externo não podem retornar bytes da sentinela. Erro explícito ou omissão conservadora são admissíveis. Um arquivo interno real e link interno legível são controles positivos. Mudança de bytes no worktree deve alterar identidade de extração, preservando conteúdo novo; repetição sem alteração deve manter identidade. Não exige uma política específica de erros ou formato específico de hash, nem alteração do commit Git original.
