# Execução 0044 — Recuperação local de implementações Go

2026-09-07. [Contrato congelado](../passos/0044-recuperar-implementacoes-por-simbolos.md). Estado: implementado e validado no recorte declarado; recuperação funcional em linguagem natural ainda insuficiente.

## Escolha e implementação

Antes de ampliar inferências, recuperar código com uma operação verificável. Reutilizamos `go/parser` e `go/ast`, sem dependência nova. O adaptador local lê arquivos com limites e contenção; o núcleo puro identifica declarações, separa identificadores e aplica expansão lexical PT/EN explícita. Declarações com identificador exato têm prioridade. Consultas com vários requisitos distribuem os primeiros resultados entre requisitos.

A seleção em `analyze` e `idea recommend` local agora prioriza os caminhos recuperados. README e manifesto continuam disponíveis; o limite de código permanece separado desses arquivos de contexto. A releitura confere SHA-256 antes de apresentar uma localização. Erro, limite e mudança dos bytes produzem cobertura parcial. Provedores sem a porta opcional registram não avaliado.

O relatório mostra símbolo, arquivo, intervalo da declaração, hash, termos e motivo da seleção; limita a exibição a um trecho identificado. O JSON conserva o inventário da varredura. Pontuação lexical não altera os detectores de capacidades. A versão de evidência passou a `evidence-v6-go-retrieval`; checkpoints locais registram `implementation_search: 0044`.

## Resultado observado

| Ensaio | Antes | Depois |
| --- | --- | --- |
| `SanitizeAssessmentEvidence` | Apenas README e go.mod; definição perdida | Recuperou `01_core/usecase/evidence_authority.go:125–198`, com hash confirmado e trecho 125–144 exibido |
| Pergunta funcional original | Não chegou às três implementações centrais observadas no diagnóstico | Ainda não chegou a `idea_discovery.go`, `reuse_decision.go` e `investigation_plan.go`; predominam declarações lexicamente relacionadas |
| Recomendação da ideia do Bateia | 10 capacidades inconclusivas; seleção pouco dirigida | 10 continuam inconclusivas; hits incluem `investigation_plan.go`, `evidence_recovery.go` e `software_alternative.go`, mas isso não comprova os requisitos |
| Replay | — | Assessment, decisão e alternativas idênticos byte a byte ao dossiê local |

A busca das consultas individuais examinou 223 arquivos Go, 1.443.663 bytes. A recomendação examinou 456 arquivos: seus requisitos mencionam testes, habilitando testes comuns. Esses números são cobertura/telemetria da ferramenta, não desempenho do projeto analisado. O limite de 24 hits tornou a busca parcial; a saída não declara ausência global. Na recomendação permaneceram 43 alvos coletados, 849 ignorados por orçamento e quatro não encontrados: aumentar a cobertura útil não se confunde com coletar mais alvos.

**Inferência limitada:** o controle exato demonstra correção daquele defeito de recuperação. As fixtures independentes também exercitam outros nomes e caminhos genéricos. Isso não demonstra boa compreensão de perguntas naturais nem qualidade global de busca. O ranking ainda pode preferir tipos ou comentários associados a palavras comuns. A próxima melhoria deve ser avaliada com consultas funcionais e implementações esperadas congeladas, incluindo negativos, antes de mudar pesos ou ampliar glossários.

## Verificação e rastreabilidade

- `GOCACHE=/tmp/bateia-0044-cache go test ./...`: passou, incluindo gates existentes; permissão ampliada apenas para servidores HTTP locais dos testes.
- `crystalline-lint .`: `✓ No violations found`.
- `GOCACHE=/tmp/bateia-0044-cache go build -o /tmp/bateia-0044 ./04_wiring`: passou. Build do pacote, pois o wiring contém vários arquivos.
- CLI: duas consultas em Markdown e JSON; recomendação local e replay; comandos, versão, commit-base, hashes dos arquivos examinados e comparação em [observations.json](../../lab/bancada/history/benchmarks/0044/observations.json).
- Especificação e baseline prévio em [baseline.json](../../lab/bancada/history/benchmarks/0044/baseline.json); testes e duas revisões explícitas de fixtures/contrato em [test-manifest.json](../../lab/bancada/history/benchmarks/0044/test-manifest.json). Nenhum ajuste de certificado de gates.

Regime Tekt A/B: autores de implementação e testes separados; filesystem compartilhado, portanto **sem atestação de isolamento físico ou certificação independente**. O coordenador integra produção e consulta resultados públicos. A especificação congelada não foi reescrita para acomodar resultados.

## Limites preservados

Varredura nova somente local/Go; não há novo índice persistente nem prova de equivalência, comportamento ou integração. A exclusão de oráculos/gates se refere à nova varredura; este passo não redesenha todos os coletores legados. Não corrigimos a classificação gRPC/Protobuf, o ruído Clang em projetos Go ou o fallback Bleve. Os relatórios anexos preservam esses defeitos como observados. A busca natural e os planos de investigação ainda precisam melhorar.

Evidências legíveis: [controle exato](../../lab/bancada/history/benchmarks/0044/symbol.md), [pergunta natural](../../lab/bancada/history/benchmarks/0044/natural.md). Os inventários não duplicam corpos de testes/oráculos. Mudanças posteriores no worktree devem ser comparadas aos hashes registrados, não tratadas como a mesma análise.
