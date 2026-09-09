# Prompt proprietário — `03_infra/content/router.go`

Origem L0 exclusiva de `03_infra/content/router.go`. Este contrato foi separado de `00_nucleo/prompts/03_infra/cache/cache_store.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Materializar somente a responsabilidade desta unidade, preservando os contratos observáveis: `NewRouter`, `provider`, `FetchManifest`, `FetchTree`, `ExtractRepoFiles`, `FetchRepoMetadata`, `Router`. Manter a direção de dependências da camada e não incorporar responsabilidades pertencentes aos demais consumidores do prompt de origem.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/03_infra/cache/cache_store.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Passo 0044

Adicionar porta opcional SearchImplementations: candidatos file:// delegam somente se
provedor local implementa ImplementationSearcher; demais casos retornam not_evaluated
com razão explícita, sem simular busca remota nem alterar ContentProvider obrigatório.

0045: provedor sem implementação registra not_evaluated para método go_ast_functional_v2, sem inferir suporte.

0046: método go_ast_functional_v3 para provedor sem suporte; permanece not_evaluated.

0046 — Veredito experimental: a tentativa v3/diversificação acima foi descartada após comparação congelada sem ganho de recall@12 e com mais retornos não rotulados. Contrato ativo de ranking permanece v2/0045; histórico experimental preservado, não implementação ativa.

0047: candidato usa método go_ast_corpus_v1 no relatório not_evaluated do provedor sem suporte.

0048: método go_ast_query_v1 no provedor sem suporte; interpretação será apresentada pela projeção da query original, sem promover suporte.

0050: SearchImplementations delega ao provedor selecionado por provider(candidate), local ou remoto, quando implementa a porta. Sem suporte permanece not_evaluated; nunca trocar origem.
