# Adaptador local de busca de implementações — passo 0044

Origem exclusiva de `03_infra/local/implementation_search.go`. Implementar porta
ImplementationSearcher com varredura ordenada via os.Root; não seguir links simbólicos.
Aplicar escopo puro antes de abrir arquivos, mesmo quando a consulta cita caminho
protegido. Defaults 2048 arquivos, 32 MiB totais, 512 KiB por arquivo, 24 hits.
Normalizar limites não positivos, limitar a leitura real pelo orçamento restante,
registrar hashes e bytes de arquivos lidos; não indexar arquivos lidos parcialmente.
Erros de leitura, parse, cancelamento e limites produzem relatório parcial explícito.
Ranking e AST pertencem ao parser puro. Sem cache persistente ou aumento de capacidade.

0045: falhas anteriores ao parsing identificam método go_ast_functional_v2; manter orçamento e contenção sem alterações. Parsing funcional continua no núcleo puro.

0046: método go_ast_functional_v3 também nas falhas antes de parsing; mesmos limites e contenção.

0046 — Veredito experimental: a tentativa v3/diversificação acima foi descartada após comparação congelada sem ganho de recall@12 e com mais retornos não rotulados. Contrato ativo de ranking permanece v2/0045; histórico experimental preservado, não implementação ativa.

0047: candidato usa método go_ast_corpus_v1 também em falhas antes do parser; não mudar leitura/orçamento/contenção.

0048: método go_ast_query_v1 nas falhas anteriores ao parser; mesma contenção/orçamento.

## Revisão de segurança 2026-09-26

Aplicar `security.go`: árvore em lotes (máximo 20.000 entradas), sem alocar diretório inteiro antes de controlar orçamento. Exceder teto produz erro explícito, não sucesso truncado. Inventário expõe symlink/hardlink/nonregular; leitura rooted mantém alias interno mas recusa hardlinks por nlink do descritor. Extração até 100 arquivos/8 MiB (2 MiB por arquivo). Não executar Git ou herdar ambiente/configuração. BaseCommit consultivo vem de HEAD/refs/packed-refs confinados; gitdir externo fica indisponível explicitamente. Metadata não autentica bytes da worktree e certificados históricos não são atualizados.
