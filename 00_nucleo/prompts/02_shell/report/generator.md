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
