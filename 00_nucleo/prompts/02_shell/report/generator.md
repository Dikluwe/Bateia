# Prompt: Gerador de Relatórios Markdown (report/generator)

**Camada**: L₂ — Casca (Apresentação)  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `02_shell/report/generator.go`, `02_shell/report/generator_test.go`  

---

## Contexto

Os resultados do Bateia são entregues como um relatório formatado em Markdown `.md`. Este componente pertence ao catálogo de apresentação da Casca ($L_2$), transformando o `domain.SearchReport` no documento final estruturado de acordo com os 3 cenários do Caso de Uso v2.

---

## Restrições Estruturais

- **L₂ (Casca)**: Função pura de apresentação (transforma struct em string Markdown).
- Pode importar `01_core/domain` e `strings`, `fmt`, `time`.
- Proibido importar $L_3$ ou $L_4$.

---

## Instrução

1. Criar a função `GenerateMarkdownReport(report *domain.SearchReport) string`.
2. Formatar o relatório contendo:
   - Cabeçalho com Query, Linguagem e Data/Hora de geração UTC.
   - Tabela do **Resumo do Funil de Redução** (Descobertos -> Triados -> Extraídos -> Descartados).
   - **Seção 1: Soluções Exatas Encontradas (Com Rastreabilidade repo@commit)**:
     - Lista dos repositórios extraídos com link URL, commit imutável, estrelas, descrição, data de consulta, dependências diretas do manifesto detectado e blocos de código extraídos (`<details><summary>...`).
   - **Seção 2: Ideias Alternativas & Abordagens Encontradas**:
     - Agrupamento das bibliotecas/frameworks identificados.
   - **Seção 3: Descobrimento de Necessidades Ocultas**:
     - Padrões de pastas e configurações observados nos repositórios analisados.

---

## Critérios de Verificação

```
Dado um SearchReport contendo 1 repositório extraído com commit "abc123"
Quando GenerateMarkdownReport for chamada
Então a string resultante contém o título "# Relatório de Análise Arquitetural", a tabela do funil e a referência "repo@abc123".
```

---

## Resultado Esperado

- `02_shell/report/generator.go`: Função `GenerateMarkdownReport`.
- `02_shell/report/generator_test.go`: Teste unitário verificando as seções geradas no Markdown.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | generator.go, generator_test.go |
| 2026-08-03 | Generalização agnóstica de linguagem (Remoção de hardcoding 'go') | generator.go, generator_test.go |

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0035 — Análise semântica C/C++

Cumprir `00_nucleo/passos/0035-corrigir-lint-e-analisar-c-cpp-com-clang.md`. Porta L1 SemanticProvider, execução Clang limitada em L3, injeção L4. Contexto explícito em compile_commands.arguments, sem executar command, plugins ou response files; capturar hashes de configuração, fonte e headers, argumentos, versão, símbolos e chamadas com intervalos. Falta de evidência nunca vira complete. Exposição L2 separa resolução estática e adequação à consulta.

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.

## Passo 0037
Aplicar o contrato de resposta por mecanismo em `00_nucleo/passos/0037-resposta-fundamentada-por-mecanismo.md`. Contexto único de compilação, trechos exatos, limites explícitos e estados desconhecidos preservados.

## Passo 0038
Usar apenas CFG do contexto selecionado e unidade completa. Intervalos ignorados não ancoram seleção lexical. Apresentar estado, localizações, limites e desconhecimento de viabilidade/polaridade.

# 0044 — Recuperação antes da seleção
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md.
L1 somente orquestra porta opcional, sem IO concreto. RetrieveImplementationEvidence usa budget finito padrão e reporta not_evaluated em provedor ausente. SelectImplementationPaths preserva manifesto/README e antepõe hits de caminhos relativos presentes na árvore, com limite de código. ValidateImplementationHits copia relatório, descarta fontes não extraídas ou hashes/intervalos divergentes e marca partial; não apaga identidade do índice.
Funnel/analyze registra relatório; CollectEvidence busca requisitos e prioriza árvore antes do orçamento, sem promoção do status de capacidades. Provedores sem porta preservam contrato antigo. A seleção por conteúdo é distinta da análise de comportamento. Relatório mostra símbolos com linhas/hash e snippet dos bytes revalidados, no máximo20 linhas por hit/8hits apresentados, JSON retém todas as evidências.

## Passo 0045 — Explicação funcional delimitada

Ao apresentar recuperação, separar modo lexical de expansão por chamada direta estática. Mostrar referências chamadoras com caminho, símbolo, linha e hash, já validadas pela releitura. Não chamar relação sintática de comportamento demonstrado. A ausência de relações ou definições não demonstra inexistência da capacidade. Preservar limites de exibição e intervalos.

Exibir no máximo três referências por hit, informando referências adicionais disponíveis no JSON.

0046: antes de hits, exibir seção de indícios por requisito (máximo8 requisitos; máximo3 fontes por requisito; termos até12). Status evidence_located explicita correspondência lexical; not_verified não demonstra ausência. Listar termos correspondentes e sem correspondência, caminhos/linhas/hash das fontes e limitar exibição com quantidade restante no JSON. Não apresentar requisito como satisfeito. Nenhum uso de score para promoção. Não omitir limites de cobertura.

0046: exportar WriteRequirementRetrieval para compartilhar a apresentação limitada entre analyze e decision.md local/replay.

0047: quando houver TermWeights no hit, mostrar até cinco contribuições com termo, campo, df/N e Weight; declarar contribuição de ranking heurístico no corpus lido, não probabilidade. Restante no JSON. Sem pesosinventados para expansão direct_call semtermos. Referências/hashes e fonte de requisito seguem0046.

0048: apresentação mostra interpretação por linha, modo e símbolos, até8 linhas/8 símbolos por linha, com total restante no JSON. Requisitos exibem até3 restrições (marcador, texto original, status not_verified) e aviso de que escopo/cumprimento não foram verificados. Renderização analyze inclui QueryInterpretation; decisão já usa projeção Requirements e deve mostrar Constraints também. Pesos para símbolos explícitos permanecem possíveis; nenhuma restrição vira ausência/satisfação. Preservar limites e hashes.

0048: exportar WriteQueryInterpretation para reutilização da apresentação em decision.md.

0048: antes do limite de oito linhas, apresentar resumo global de restrições: até três marcadores/textos originais/status e quantidade total, inclusive quando o requisito estiver fora das primeiras linhas exibidas. Assim truncamento visual não esconde a existência de negações; restante no JSON.

0051: hits documentation exibem rótulo documental e somente trecho do grupo Doc (máximo20 linhas), mantendo localização da declaração e hash. Projeção separa DocumentationSources de Sources; até3 fontes documentais e contagem restante. Não afirmar implementação analisada pelo texto documental.

0052: WriteQueryInterpretation exibe Expansions da consulta, até6 registros/linha e8 termos/registro, contagens omitidas, Source/RuleID/Kind/Status. Aviso hipótese lexical; tradução/alternativa técnica não provam pertinência ou capacidade.

## Passo 0055 — Hipótese de consulta
Ao exibir tentativa com QueryPlan, apresentar Status, Version, Reason, RuleIDs, ActionTerms, ObjectTerms, Unverified, RawCandidates e RejectedCandidates. Explicitar que admissão lexical é indício e não comprova solução ou restrições. Preservar OriginalQuery e Query efetiva.

Adendo de evidência 0055: exibir AdmissionEvidence como indício lexical observado, com repositório, caminho e snippet retido. Nunca converter esse rastro em prova de solução.

## Passo 0056
Aplicar 00_nucleo/passos/0056-explicar-cobertura-e-perdas.md: report exibe cobertura desconhecida, motivos e limites sem inferir ausência de capacidade; checkpoint live discovery_query passa0056-v1.

## Passo0057
Exibir Behavior como observação sintática por função, com referência/hash/linhas e limitações, inclusive not_evaluated/zero fatos. Máximo10funções/8fatos, contagens omitidas e JSON integral. Não afirmar comportamento comprovado.


## Passo 0061
Chamar WriteActionAssessments junto às projeções existentes da busca de implementação.
