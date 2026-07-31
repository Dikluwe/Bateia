# Prompt: Parser de Manifestos Python (parser/python)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/python.go`, `01_core/parser/python_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para o ecossistema Python, o conteúdo de texto de manifestos como `pyproject.toml` (Poetry/Flit/Hatch/Setuptools) ou `requirements.txt` precisa ser analisado para extrair o nome do projeto e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParsePyProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name` sob blocos `[tool.poetry]`, `[project]` ou `[tool.flit.metadata]`.
   - Extrair dependências de produção sob `[tool.poetry.dependencies]`, `dependencies = [...]` ou `requires = [...]`.
   - Extrair dependências de dev sob `[tool.poetry.group.dev.dependencies]` ou `[project.optional-dependencies]`.
   - Marcar `Language: "python"`.

2. `func ParseRequirementsTxt(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas linha a linha (ex: `fastapi>=0.100.0`, `sqlalchemy==2.0.0`).
   - Ignorar linhas de comentário (`#`) ou flags do pip (`-r`, `-e`, `--index-url`).
   - Marcar `Language: "python"`.

---

## Critérios de Verificação

```
Dado um pyproject.toml com name "my-api" e dependencies "fastapi", "pytest"
Quando ParsePyProjectToml for executada
Então retorna um Manifest com ModuleName "my-api", Language "python" e 2 dependências.

Dado um requirements.txt com "django>=4.2" e "psycopg2-binary"
Quando ParseRequirementsTxt for executada
Então retorna um Manifest com Language "python" e 2 dependências.
```

---

## Resultado Esperado

- `01_core/parser/python.go`: Implementação dos parsers de Python.
- `01_core/parser/python_test.go`: Suite de testes cobrindo `pyproject.toml` e `requirements.txt`.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Python | python.go, python_test.go |
