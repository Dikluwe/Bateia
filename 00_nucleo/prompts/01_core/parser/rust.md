# Prompt: Parser de Manifestos Rust Cargo.toml (parser/rust)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/rust.go`, `01_core/parser/rust_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para o ecossistema Rust, o conteúdo de texto do arquivo `Cargo.toml` precisa ser analisado para extrair o nome da crate/pacote, edição Rust e dependências (produção, dev-dependencies e build-dependencies).

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar a função:
`func ParseCargoToml(ref domain.CommitRef, content string) (*domain.Manifest, error)`

Regras do parsing:
1. Extrair o nome do pacote `name = "..."` sob o bloco `[package]`.
2. Extrair a edição `edition = "..."` sob o bloco `[package]`.
3. Extrair dependências diretas de produção sob `[dependencies]`.
4. Extrair dependências de desenvolvimento sob `[dev-dependencies]` ou `[target.*.dev-dependencies]`.
5. Extrair dependências de build sob `[build-dependencies]`.
6. Suportar sintaxe simples (`tokio = "1.0"`) e sintaxe inline table (`serde = { version = "1.0", features = ["derive"] }`).
7. Preencher `Language: "rust"`.
8. Preencher `ConsultadoEm` com a data/hora UTC atual.

---

## Critérios de Verificação

```
Dado um Cargo.toml com name "axum-app", tokio = "1.35.0" em [dependencies] e criterion = "0.5" em [dev-dependencies]
Quando ParseCargoToml for executada
Então retorna um Manifest com ModuleName "axum-app", Language "rust", 1 dependência "direct" (tokio) e 1 dependência "dev" (criterion).
```

---

## Resultado Esperado

- `01_core/parser/rust.go`: Implementação do parser de Cargo.toml.
- `01_core/parser/rust_test.go`: Suite de testes cobrindo Cargo.toml com tabelas simples e inline.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Rust | rust.go, rust_test.go |
