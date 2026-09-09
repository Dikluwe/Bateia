# Passo 0044 — Recuperar implementações antes de analisar capacidades

Estado inicial: especificação congelada antes da implementação. Resultado será registrado em relatório separado.
Data: 2026-09-07. Origem: autoanálise pós-0043.

## Escolha

Adotar busca local por símbolos e conteúdo como etapa anterior à seleção de arquivos. Reutilizar `go/parser`/`go/ast` já disponíveis, inspirado na separação busca/navegação do [Sourcegraph](https://sourcegraph.com/docs/deep-search) e nas referências de busca do [Bloop](https://github.com/BloopAI/bloop). Não instalar produtos nem copiar código desses projetos neste passo. Não adicionar LLM ou embeddings antes de recuperar definições exatas de forma verificável.

## Contrato

Porta opcional L1 `ImplementationSearcher.SearchImplementations(ctx, candidate, query, budget) (ImplementationSearchReport,error)`; adaptador local L3 lê bytes com contenção os.Root já existente. Parser L1 puro `RankGoImplementations(query string, files []ExtractedFile, limit int) ImplementationSearchReport` usa AST para definições e tokens para pertinência. Outros provedores continuam funcionando sem a porta, mas informam busca de implementação não avaliada.

`ImplementationSearchBudget`: MaxFiles int, MaxBytes int64, MaxFileBytes int64, MaxResults int. Defaults: 2048 arquivos, 32 MiB, 512 KiB por arquivo, 24 resultados. Limites de varredura são distintos dos limites de extração. Não construir cache persistente neste passo. Registrar arquivos efetivamente lidos, hashes, bytes, limites, erros de leitura/parse e truncamento da busca.

`ImplementationSearchReport`: Status, Method, Reason strings; Hits []ImplementationHit; Examined []ImplementationInput; FilesRead int; BytesRead int64; Budget ImplementationSearchBudget; Limited bool; Issues []string. `ImplementationHit`: Path, Symbol, Kind, SHA256, Reason strings; Line, EndLine, Score int; MatchedTerms []string. `ImplementationInput`: Path,SHA256 strings; Bytes int64. Score é ordenação heurística, não probabilidade nem capacidade comprovada.

Identificadores exatos encontrados como declaração vencem menções em comentários/literais e coincidência de nome de arquivo. Índice cobre somente Go; arquivo inválido ou linguagem sem suporte não gera definição. Separar CamelCase/snake_case, termos comuns e palavras funcionais sem inventar sinonímia ou detectores específicos para Bateia/Coin. Qualquer pequeno glossário PT/EN deve ser explícito e genérico, declarado como expansão lexical. Duplicatas e reordenação de entradas não alteram resultados.

Ignorar por padrão diretórios temporários/vendorizados/gerados e testes de fixture. Oráculos, gates e metadados privados (.git/.agents/.codex) ficam fora da varredura. Testes comuns só entram quando explicitamente buscados, sem abrir corpus protegido. Caso de ausência nunca prova inexistência além do escopo e orçamento.

Integrar a seleção ao funil `analyze` e à coleta de `idea recommend` local. Preservar README/manifesto e limites finais existentes; priorizar caminhos dos hits antes de ranking de caminhos, com distribuição por requisito quando possível. Releitura deve conferir hashes: se o arquivo mudou, declarar obsolescência, sem apresentar hit como confirmado. Registrar a busca no JSON e Markdown como descoberta de implementação estática, nunca prova comportamental. Retomar/replay conserva relatório de busca com seus hashes e valida incompatibilidade de versões.

## Aceitação pré-declarada

1. Consulta exata `SanitizeAssessmentEvidence` extrai `01_core/usecase/evidence_authority.go` e apresenta definição/linha/hash, sem aumentar a nota de comprovação de capacidade automaticamente.
2. Fixture Go em arquivo de nome genérico com função solicitada é recuperada; README ou comentário com o mesmo nome não a substitui. Casos de receita e fila também exercitam recuperação, sem regras com nomes de repositórios.
3. Consulta funcional original é repetida; medir arquivos úteis recuperados e declarar insuficiência remanescente, sem exigir um resultado positivo fabricado.
4. Símbolo inexistente, parser inválido, budget esgotado e alteração entre índice/leitura permanecem explícitos; não geram prova ou ausência global.
5. Testes de contenção local permanecem passando. Nenhum IO em L1; prompts e linhagem precedem Go. Testes A/B separados da implementação; filesystem compartilhado impede atestação física.
6. Suíte completa, lint e build do pacote `./04_wiring` passam. Repetir CLI local e replay, registrar saídas e diferença contra autoanálise anterior.

## Fora do escopo

Não resolver equivalência de programas, fluxo dinâmico, índices semânticos de todas as linguagens, recall público ou integração executável. Não alterar o catálogo de capacidades para forçar a autoanálise a aprovar o Bateia. Não corrigir neste mesmo passo gRPC/Protobuf, Bleve ou todos os problemas de relatório observados.
