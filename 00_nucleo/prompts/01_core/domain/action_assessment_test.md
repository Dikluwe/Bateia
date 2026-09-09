# Testes independentes de ações — passo 0061

Papel B do regime Tekt A/B, ambiente compartilhado sem atestação física.
Raiz normativa: `00_nucleo/passos/0061-discriminar-acoes-estaticas.md`.

JSON snake_case, omitempty no relatório e roundtrip integral dos tipos de ações.

Usar fixtures neutras próprias. Não ler produção candidata, corpus ou relatórios. Não executar antes de congelar os testes. Não alterar expectativas legadas. Nenhum indício sintático prova execução, efeitos ou satisfação.
