# Prompt: Parser de Julia (parser/julia)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/julia.go`, `01_core/parser/julia_test.go`  

---

## Contexto

Parsing de manifestos `Project.toml` do ecossistema Julia.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Recebe string de conteúdo e retorna `*domain.Manifest`.
- Proibido I/O ou dependências externas.
