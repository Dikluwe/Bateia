# Testes independentes da revalidação anterior à projeção — 0046

Autor B, regime A/B sem atestação física. Usar apenas o contrato do passo 0046 e API pública. Gerar conteúdo sintético em memória e SHA256 real. Releitura idêntica preserva a fonte lexical; conteúdo alterado ou não extraído elimina a fonte e mantém requisito not_verified, mesmo quando o relatório traz projeção anterior. Não executar IO no teste L1. Congelar hashes antes de executar. Não ler candidata, testes anteriores, diffs ou gates.
