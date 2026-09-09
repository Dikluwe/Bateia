# Testes independentes de ações — passo 0061

Papel B do regime Tekt A/B, ambiente compartilhado sem atestação física.
Raiz normativa: `00_nucleo/passos/0061-discriminar-acoes-estaticas.md`.

Projeção Markdown segura, ação/status/razões/identidade, limites 12 candidatos e 4 evidências e preservação integral no JSON, integração no gerador.

Usar fixtures neutras próprias. Não ler produção candidata, corpus ou relatórios. Não executar antes de congelar os testes. Não alterar expectativas legadas. Nenhum indício sintático prova execução, efeitos ou satisfação.

## Revisão B r1
Normalizar entidades HTML com `html.UnescapeString` apenas nas verificações de conteúdo visível do Markdown. Escaping seguro não altera o conteúdo contratual; underscores codificados não devem causar falso negativo. Preservar os limites 12/4, JSON completo e todas as verificações de segurança sobre o texto bruto. Revisão motivada pelo diagnóstico textual do integrador, sem acesso à produção candidata ou execução por B.

## Revisão B r2 — Adendo de apresentação
Acrescentar fixture neutra em ordem deliberadamente adversa, com relações intercaladas e candidatos alinhados após a posição 12. Exigir prioridade aligned_hint, conflicting_hints, opposed_hint, unknown/outros antes do corte 12, estabilidade entre empates inclusive valores não reconhecidos e unknown. Comparar serialização antes/depois das chamadas direta e integrada para detectar ordenação in-place ou qualquer mutação da lista, evidências ou JSON. Manter controles anteriores integralmente.
