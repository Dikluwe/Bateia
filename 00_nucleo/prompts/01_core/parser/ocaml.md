# Prompt: Parser de OCaml (parser/ocaml)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/ocaml.go`, `01_core/parser/ocaml_test.go`  

---

## Contexto

Parsing de manifestos `dune-project` e `*.opam` do ecossistema OCaml.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura.

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.
