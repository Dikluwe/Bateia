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
