# Prompt: Parser de COBOL (parser/cobol)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/cobol.go`, `01_core/parser/cobol_test.go`  

---

## Contexto

Parsing de arquivos de código fonte `*.cbl` / `*.cob` do ecossistema COBOL (Copybooks & Calls).

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.
