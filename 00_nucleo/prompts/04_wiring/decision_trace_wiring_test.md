# Despacho real de trace offline — 0098

Testar main em subprocesso do próprio binário de testes, com timeout e cwd isolado.
Uma .env sintética declara GITHUB_TOKEN igual ao RuleID de uma alegação do build.
O writer recusaria a saída se a .env fosse carregada; trace deve emitir JSON válido
sem carregar env, léxico ou provedores. Remover apenas a variável token e chave
do helper do ambiente do processo filho; não reatribuir HOME. Nenhum candidato
é compilado/executado. Fixture e marcador são sintéticos, não credenciais reais.
