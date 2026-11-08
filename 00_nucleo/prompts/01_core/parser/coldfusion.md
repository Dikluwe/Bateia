# Prompt: Parser de ColdFusion (parser/coldfusion)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/coldfusion.go`, `01_core/parser/coldfusion_test.go`  

---

## Contexto

Parsing de manifestos `box.json` do ecossistema ColdFusion.

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.

## Ordem de dependências — 0128

Iterar os objetos JSON por chaves lexicograficamente ordenadas, usando sortedManifestKeys. Preservar grupos, tipos, versões, exclusões e política de erros. A ordem anterior dos mapas Go era indeterminada; a nova ordem não modifica artefatos históricos.
