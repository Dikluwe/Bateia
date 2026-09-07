# Prompt espelhado — Oráculo End-to-End L4 do Passo 0032

Exercite de ponta a ponta todas as vias de consumo operacionais do Bateia de forma determinística e offline:

1. CLI Local:
   - Invoque o funil da CLI em worktree local mockada ou existente;
   - Assegure que as três etapas (descoberta, triagem AST Go e extração) ocorram com sucesso;
   - Comprove que o relatório Markdown é emitido com estrutura íntegra e sem erros.

2. Servidor MCP via stdio:
   - Instancie a composição L4 com suporte a persistência durável de execuções;
   - Exercite o handshake `initialize` e comprove a conformidade JSON-RPC 2.0;
   - Exercite `tools/list` e valide o anúncio estrito das ferramentas duráveis (`bateia_recommend`, `bateia_start_recommendation`, `bateia_get_run`, `bateia_get_artifact`, `bateia_cancel_run`);
   - Exercite `bateia_recommend` fornecendo somente `idea` (comprovando a atuação dos defaults V1 consolidados no Passo 0031);
   - Exercite `bateia_start_recommendation` com `idempotency_key` mínima, capture o `run_id` emitido em estado `QUEUED`;
   - Exercite `bateia_get_run` com o `run_id`, comprovando a integridade da persistência durável em disco até resolução em `SUCCEEDED`;
   - Exercite `bateia_get_artifact` com hash seguro SHA256 e `bateia_cancel_run`.

3. Invariantes de Isolamento:
   - Proibido qualquer acesso de rede externa durante a suíte E2E;
   - Zero I/O em L1.
