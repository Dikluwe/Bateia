# Prompt: Parser de BASIC / VB (parser/basic)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/basic.go`, `01_core/parser/basic_test.go`  

---

## Contexto

Parsing de manifestos `.vbp` do ecossistema BASIC/Visual Basic.

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.
