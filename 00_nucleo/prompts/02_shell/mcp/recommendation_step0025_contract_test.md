# Prompt espelhado — oráculos A/B L2/L4 do Passo 0025

Derive exclusivamente do Passo 0025, do contrato congelado e dos ataques congelados uma
suíte independente para o gateway MCP síncrono `live`. Preserve `NewMCPServer(funnel,
reader, writer)` e permita a injeção compatível da porta `RecommendationRecommender`.

Congele `tools/list` na ordem `bateia_search`, `bateia_benchmark`, `bateia_recommend`, com
nome, descrição e schema exato da terceira ferramenta e sem campos de credencial, caminho,
modo, provider ou estado. O decoder rejeita omissões, valores inválidos, campos extras,
tipos incompatíveis e múltiplos valores antes de chamar o serviço. Sucesso chama a porta
uma vez e devolve um único texto JSON compacto produzido pela projeção pura de L1, sem
refazer ranking. Erros totais usam `-32603` sanitizado; validação usa `-32602`; erros e
schemas antigos permanecem iguais à baseline. Repetições são determinísticas e sem estado.

Audite por AST que `04_wiring/main.go` constrói `NewRecommendationService` e passa a mesma
instância a `NewMCPServer`; construção unitária isolada não satisfaz a fiação L4.

Na revisão R2, omita em runtime cada obrigatório crítico e cada limite, rejeitando-os sem
default; ataque individualmente todos os campos proibidos; congele sucesso e erro legados,
proveniência completa e `excluded`; proíba imports/chamadas de CLI, infraestrutura e
filesystem; e compare sucesso repetido no mesmo servidor com servidor lógico fresco.

Na revisão R3, rejeite qualquer import de `os`, `io/ioutil`, `path/filepath` e seletores
`Create*`, `OpenFile`, `Temp*`, `Mkdir*`, `ReadFile` ou `WriteFile`, além da CLI. Construa
respostas esperadas independentes do candidato e compare bytes do sucesso de search. Para
benchmark, preserve o erro real do baseline: o handler legado não propaga `depth` ao filtro,
portanto não existe sucesso nominal reproduzível sem mudar comportamento antigo. Em decision
e alternatives, afirme independentemente todos os campos públicos de `EvidenceSource`,
removendo somente `snippet`, e preserve `top_choice`, `excluded` e `discovery_locales`.

Na revisão R4 corrigida, M20 usa somente a constante literal com os bytes completos do
sucesso de search capturado no baseline `21c48e3571807777013710bcb6675c95c4130d9e`.
M21 congela o erro público literal de benchmark já exercitado pela suíte. Nenhum expected
chama `SearchFunnel`, gerador de relatório ou outro código de produção em runtime.

Na revisão pós-selo R6, preserve byte a byte o golden literal de sucesso de search, mas
canonicalize antes da comparação somente o valor volátil da linha `Gerado em`, produzido
por relógio não injetável no baseline. A linha deve existir exatamente uma vez e seu
envelope permanece congelado; nenhum outro byte do relatório pode ser normalizado.
O teste legado proprietário de `tools/list` também exige as três ferramentas contratuais,
na ordem congelada, mesmo quando o recommender opcional é `nil`.

No refinamento explícito do Passo 0031, substitua somente os arrays `required` de
policy/search/raiz e acrescente as anotações de defaults congeladas. Campos
operacionais omitidos deixam de ser casos negativos; problema/capacidades,
valores explícitos inválidos, campos desconhecidos, ordem, nomes, projeção e
isolamento continuam protegidos. Preserve esse sucessor como golden independente.

Não leia, adapte nem corrija implementação candidata. `Unknown` não satisfaz caso negativo.
