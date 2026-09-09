# Testes independentes da interpretação0048

Autor B, regime A/B sem atestação física. Entradas: AGENTS, passo0048 e tipos públicos anteriores; sem leitura de implementação candidata, ranking, gates ou testes anteriores na autoria de casos.

Congelar antes da candidata casos naturais e símbolos explícitos: operadores PT/EN/case, identificador isolado/backticks, CamelCase/snake_case, crases inválidas, apenas operadores, linhas duplicadas/reordenadas. Verificar preservação de texto e status not_verified, lexicalidade dos demais termos e determinismo. Ranking por API pública deve excluir sem natural, preservar sem explícito e manter frequência do corpus. Requisitos devem conservar restrições não verificadas mesmo quando evidência positiva existe; relatório antigo e sem fontes não pode promover cumprimento.

## Revisão mecânica r1

V14 rejeitou import slices fora do allowlist L1. Substituir Equal/Contains por loops locais, Sort por sort.Strings e Clone por append em slice novo. Preservar consultas, expectativas, asserts e semântica nil/empty de Equal. Sem alteração de produção ou configuração lint. Preservar freeze original e registrar nova cadeia antes reexecutar.

## Extensão independente r2 após contraexemplo público

Identificadores Go explícitos s/use não podem desaparecer por filtros de linguagem natural. Sem leitura candidata, adicionar controles isolados e entre backticks: símbolo preservado, token cru presente, função exata recuperada e fonte exata projetada como evidence_located. Extensão posterior à implementação inicial, congelada antes de sua execução e correção subsequente. Preservar todos esperados anteriores e rótulos dos corpora.
