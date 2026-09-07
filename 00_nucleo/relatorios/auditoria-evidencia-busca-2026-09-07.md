# Auditoria do caminho busca → afirmação — 2026-09-07

Escopo: fluxo `SearchFunnel` e relatório Markdown de busca/benchmark. Diagnóstico por leitura de código e regressões offline; não houve repetição de buscas públicas ou análise do checkout real do Coin. Os relatórios específicos USalign/BinDiff/Coin citados no pedido não foram localizados nos arquivos versionados pesquisados. Os nomes e sintomas do pedido foram usados em fixtures; não se atribuem termos de README históricos que não foram fornecidos.

| Tema | Causa observada no código anterior | Correção / resultado |
| --- | --- | --- |
| Relevância | `EvaluateRelevance` tinha `EqualFold(manifest.ModuleName, cand.Repo)` como condição suficiente para 85. Parsers C/C++ podem copiar `ref.Repo` ao módulo quando o manifesto não informa nome. | Removida a condição circular. Nome de módulo é comparado à consulta. USalign e bindiff com módulo homônimo e consulta OIT ficam com 0 quando não há termos pertinentes. |
| README | Prefixos de cinco caracteres e um único termo autorizavam 65. `Manifest.Content` também continha manifestos, sem distinguir sua origem. | `ContentPath` identifica README; termos literais únicos, score proporcional à fração coberta, campo e termos na justificativa. Fixture com `OpenGL transparency` para quatro termos dá 32/100; não prova implementação. |
| Consulta | Sourcegraph tentava consulta completa com linguagem, dois primeiros termos com linguagem e esses termos sem linguagem, sem relatar tentativas. OIT podia virar `order independent`. | Tentativas, erros, consulta original, consulta enviada e relaxamentos registrados; preservados nos candidatos em cache. Provedores sem rastreio ficam explicitamente não verificados. Resultados recursivos de dependências são reavaliados pela consulta da chamada pai. |
| Arquivos | Primeiros três arquivos Go/Rust/TS/Python/Java; C/C++ e shaders excluídos. Sem relação entre caminho e pergunta. | Ranking determinístico por termos no caminho, até 12 fontes mais manifestos/README; C/C++ e shaders incluídos. Penalidades para tmp, vendor, generated etc. Sem prova semântica, cobertura fica `insufficient`. |
| Cache | Extração era recuperada só por commit, embora a pergunta devesse mudar a seleção. | Reuso desse cache desativado neste fluxo; mantém cache de descoberta/manifesto. Maior custo de extração é uma consequência deliberada. Manifesto final é reparsed dos bytes extraídos. |
| Truncamento | `TruncateSmart`: 1.500 runes / 40 linhas, na apresentação. Parsers recebiam `ExtractedFile.Content`; a IA via MCP recebia Markdown truncado. | SHA-256, método, linhas examinadas, termos e linhas de ocorrência, intervalo exibido e arquivos solicitados/não extraídos. Lexical integral não equivale a leitura semântica integral pela IA. |
| Scores | Relevância não controla classificação “solução exata”. Similaridade usa nome/descrição, dependências e caminhos; soma intenção 35%, stack 30%, capacidades 20%, topologia 15%. Capacidades genéricas eram injetadas na ideia-alvo. | Todos os resultados são candidatos. Relatório reconcilia métricas por candidato, evidencia `Divergent`/intenção baixa e não promete solução. Capacidades genéricas do alvo removidas; classificação de “ideia idêntica” substituída por similaridade heurística não verificada. |
| Saúde | 100 menos 20/10/5 por gravidade; zero achados virava 100, inclusive com manifesto ausente. | Relatório usa “Não avaliado”; `health_score=-1` neste fluxo. `CalculateHealthScore(nil)=-1`. Não existe cobertura validada que sustente uma nota de saúde global. |
| Linguagens | Default de ORM era Go e substring `ent` confundia dependências; recomendação default usava `*_test.go`; insights selecionavam catálogo pela primeira linguagem de todo o lote. | Go explicitamente condicionado; identidades/segmentos de dependências substituem substring; recomendações genéricas fora das linguagens reconhecidas. Catálogos aplicados por grupo de linguagem. Dependência não prova vazamento arquitetural; achado passa a hipótese, sem linha fictícia de acesso ao ORM. |
| Apresentação | “Soluções Exatas Encontradas”, “Descobrimento de Necessidades Ocultas”, “Commit Rastreável” e fallback com padrões Go mesmo sem resultados. | Candidatos, indícios/hipóteses, referência sem imutabilidade presumida e ausência de evidência explícita. Contrato distingue fato observado, inferência, hipótese e não verificado. Tempos identificados como telemetria do Bateia, não benchmarks do projeto. |
| Worktree | `WORKTREE` não era resolvido por `FetchRepoMetadata`. Nenhum hash de arquivos examinados. | Commit-base via Git, hash de cada conteúdo e hash ordenado do conjunto de arquivos. Teste altera fonte rastreada e acrescenta não rastreada mantendo mesmo commit-base: hashes mudam. Clone remoto de fallback deve conseguir fazer checkout do commit solicitado e usa os caminhos selecionados. |

## Evidências e limites

- Relevância: `01_core/parser/relevance.go`, `01_core/parser/evidence.go`; regressões em `evidence_test.go`.
- Orquestração, cache e intenção original: `01_core/usecase/funnel.go`.
- Consulta efetiva: `03_infra/sourcegraph/client.go`; teste `TestSearchTraceRetainsAllRelaxations` inclui resposta vazia nas duas primeiras tentativas.
- Worktree: `03_infra/local/client.go`; teste `TestWorktreeRecordsBaseAndChangedUntrackedBytes` usa repositório Git temporário com commit real.
- Apresentação: `02_shell/report/generator.go`; `TestEvidenceLanguageDoesNotPromoteSearch` verifica zero resultados, score 85 com Divergent e saúde legada 100.
- Não há alegação de testes segregados ou veredito independente. A sucessão do golden MCP está documentada em `evidence-contract-supersession.json`; não altera certificados históricos.

O ranking é lexical por caminhos, não por entendimento de implementações. A fixture Coin verifica a seleção de fontes do renderizador quando os nomes contêm os termos; não comprova cobertura do Coin real. O dispatcher de símbolos deste fluxo continua limitado a Go (AST) e Rust (heurísticas). C/C++ têm extração textual, sem prova semântica. A ausência de achados nunca promove cobertura insuficiente a sucesso.

O leitor local limita árvore a 20.000 entradas, ignora alguns diretórios e recusa arquivos maiores que 2 MiB. Logo, caminhos já omitidos pela infraestrutura não podem ser recuperados pelo ranking, mesmo quando um diretório vendorizado é explicitamente citado. Arquivos solicitados que não foram lidos aparecem como ausentes; as razões específicas de cada falha não são fornecidas pela interface atual. O snapshot identifica os bytes lidos, mas não é captura atômica do worktree, backup dos arquivos ou prova de que o checkout não mudou durante a leitura. Reproduzir bytes modificados exige preservá-los além dos hashes.

As classificações de dependências continuam heurísticas e os catálogos não são completos. A matriz de protocolos também é uma inferência por dependências, não uma observação de comunicação em execução. Relatórios históricos não foram reescritos: refletem o comportamento anterior.

## Validação final

- `go test ./...`: exit 0, incluindo regressões de evidência e contratos MCP atualizados.
- `crystalline-lint .`: exit 0, nenhum erro; 195 avisos de L4/prompts órfãos.
- `go build -o bateia ./04_wiring`: exit 0. A forma documentada `./04_wiring/main.go` falha por não incluir funções definidas nos demais arquivos do pacote.
- `git diff --check`: exit 0.

Testes HTTP foram executados com permissão para portas locais; nenhuma busca pública foi usada para demonstrar as regressões. Saídas sintéticas não são reproduções dos relatórios históricos nem benchmarks dos projetos citados.
