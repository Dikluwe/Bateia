# Resultado 0061 — Hipóteses da ação por candidato

O Bateia agora separa correspondência lexical de uma comparação limitada da ação pedida. Para início/interrupção, examina nomes de declarações e identificadores de chamadas Go confirmados nos bytes e intervalos físicos. Cada candidato recebe aligned_hint, opposed_hint, conflicting_hints ou unknown, com arquivo, hash, linhas e trecho do identificador. Esses rótulos são hipóteses; execução, resolução de chamadas, objeto e efeitos permanecem não verificados.

## Resultado observado

No caso de interrupção do corpus, a nova análise produz 2 indícios alinhados, 6 opostos e 16 desconhecidos. Remove é aligned_hint; Add é opposed_hint. No controle dirigido que troca apenas iniciar/parar, as relações se invertem corretamente. Fonte: benchmark/0061/action-comparison.json, pair-contract.json e pair-after.json. O controle usa fonte conhecida e não é holdout.

Remove: fsnotify.go:333, nome Remove e chamada Remove. Add: fsnotify.go:314, nome Add e chamada AddWith. SHA256 c721c84f2a2c6734ec0e82c0c9f4133ba3f62761457b52cc2a26aaf511e2c78a. Isso observa rótulos e chamadas, não demonstra que notificações realmente cessem ou comecem.

As outras oito perguntas do corpus ficam com ação not_evaluated: fora das famílias suportadas ou com restrição de negação. Essa limitação é explícita. Documentação, identidade inválida, ausência de bytes, conflito de pistas ou orçamento excedido não são promovidos a comprovação.

Recuperação permanece 6/6 alvos: quatro pelo código, dois pela documentação. Canais e hits são idênticos ao passo 0060; negativos permanecem com zero indicações e zero conclusões indevidas. Fonte: benchmark/0061/retrieval-comparison.json e collected-r1/summary.json. Não houve alteração de scores ou filtros da busca.

## Apresentação

JSON mantém candidatos completos. Markdown prioriza uma cópia por indícios alinhados, conflitantes, opostos e desconhecidos antes de limitar a exibição a 12 candidatos e 4 evidências por candidato. Contagens remetem ao JSON. É prioridade de leitura, não uma nota de confiança. Exemplo final: benchmark/0061/action-report-r1.md.

A revisão inicial ordenava por caminho e escondia Add/Remove no corte. Esse problema foi registrado e corrigido somente no apresentador, preservando entrada e ranking. O teste independente cobre prioridade antes do corte, estabilidade e não mutação.

## Validação e limitações do processo

Tekt A/B sem atestação física, ambiente compartilhado. A leu contrato/produção e B contrato/tipos/testes legados; não houve troca de código candidato com B ou testes com A. Integrador conhece os alvos. Contrato e hashes das revisões estão preservados em benchmark/0061.

Ocorrências preservadas: o primeiro filtro de testes não encontrou nomes e foi corrigido para TestAction61; dois testes de Markdown exigiam texto bruto quando caracteres estavam escapados, sendo corrigidos para comparar conteúdo visível e manter segurança no texto bruto. B informou que um rascunho inicial antecedeu os prompts por falha no comando python; foi removido e rematerializado depois. Não alegamos cronologia estrita para esse rascunho nem conformidade temporal atestada.

Os dez controles novos abrangem pares de ação, nomes enganosos/conflitantes, comentários/literais/funções anônimas, identidade adulterada, linhas físicas, orçamento, integração, JSON e Markdown. Resultados finais de build, testes, lint e replay estão em verification.json e logs do experimento.

## O que continua sem solução

Este passo é um comparador finito de pistas sintáticas, não compreensão geral de programas. Uma função chamada Stop pode mentir no nome; chamadas podem estar em ramos nunca executados ou ter efeitos diferentes. Objeto, condições, fluxo entre funções e efeitos reais não foram verificados. Não executar os projetos ou a busca pública impede alegar validação funcional ou melhora da descoberta remota. Tempos são telemetria do Bateia.

Validação final: build, suíte Go completa r2, dez testes novos, nove controles do runner, crystalline-lint e diff-check passaram. Replay reproduziu o resumo; hashes congelados conferidos.
