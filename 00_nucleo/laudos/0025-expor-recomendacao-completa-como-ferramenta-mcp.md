# Laudo 0025 — Gateway de recomendação completa por MCP

**Data**: 2026-09-01  
**Estado**: `IMPLEMENTADO`  
**Veredito**: `MCP_RECOMMENDATION_GATEWAY_READY`  
**Regime Tekt**: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

## Resultado

O corte síncrono `live` de `bateia_recommend` está aprovado para os hashes
registrados no certificado. O MCP anuncia exatamente as três ferramentas
esperadas, rejeita entrada ausente, desconhecida ou fora dos limites antes das
portas externas e delega a execução a `RecommendationService` em L1. A CLI live
e o MCP usam essa mesma orquestração de domínio; o adaptador MCP não executa a
CLI, não cria dossiê ou arquivo temporário e não refaz ranking em L2.

A resposta preserva decisão, alternativas, exclusões, evidências rastreáveis e
locales, removendo snippets e conteúdo bruto. Licença desconhecida continua
visível como aviso não bloqueante. Falhas parciais e estados opacos não são
promovidos a `COMPLETE`, e as ferramentas legadas mantêm os observáveis
congelados pelo baseline.

## Reavaliação das correções finais

O primeiro recibo final, identificado por
`e295ffdec609515379bae62283ce865cc2a5adb72476ee2a25bdc471bd68ae2e`,
bloqueou a candidatura por duas violações. As mesmas testemunhas independentes,
sem alteração, foram repetidas contra a correção:

- M18 agora passa: cancelamento observado durante uma porta de conteúdo impede
  `ExtractRepoFiles` e impede o início do candidato seguinte;
- M11 agora passa: mensagens externas são projetadas como categoria pública
  estável, e falha de coleta embutida em razão não ecoa o texto do adaptador.

O receipt sucessor preserva o predecessor por hash, a fixture independente e os
resultados anteriores e atuais. O mutation score permanece `26/26 = 1,0` e
O01–O05 permanecem `Unknown`.

## Gates finais

- os dois overlays independentes M18/M11 passaram;
- `go test -count=1 ./...` passou duas vezes em execuções separadas pelo
  verificador, além das duas passagens informadas pelo coordenador;
- `go build`, `go vet` e `git diff --check` passaram;
- `crystalline-lint .` terminou com exit 0; os avisos V12 restantes pertencem
  somente aos fakes da suíte L4 protegida;
- a auditoria de imports confirmou L1 sem dependência de shell/infra e L2 sem L3;
- o teste offline L4 anunciou e chamou `bateia_recommend` pela composição real.

O primeiro ensaio da suíte completa dentro do sandbox falhou apenas porque três
testes `httptest` não puderam abrir listeners localhost. As duas repetições com
essa capacidade permitida passaram integralmente; a falha ambiental não foi
tratada como falha do produto.

## Cadeia e limites

- manifesto final: `9f4a06917139098b2a7a63f36cb438b936de9bdbb95d049dd9277b7f1296bc5f`;
- selo R6: `9ee11e374a114cf760b1bdac51fa17fa4dd49ab7717671555a5b2820c090f5af`;
- receipt final: `a49e2733cbf47f03a7f537b8ec4efd548333302d36ad877a4422344f95f4bf4a`;
- certificado: `bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d`.

A execução ocorreu em worktree compartilhada, e o refinamento R6 e as correções
finais ocorreram depois de já existir candidato. Portanto o laudo registra
reprodutibilidade e segregação de autoridade, mas não alega isolamento forte nem
independência temporal plena.

O veredito cobre somente o gateway MCP síncrono `live` e os hashes registrados.
Não implica equivalência funcional geral, persistência MCP, `run_id`, replay,
resume, execução assíncrona, cobertura geral de mercados ou qualidade de
tradução multilíngue.
