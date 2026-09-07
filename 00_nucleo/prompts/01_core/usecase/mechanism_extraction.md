# Deterministic mechanism extraction

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/mechanism_extraction.go`

Implemente o extrator puro e determinístico R4 sobre o `EvidenceBundle` congelado
e `[]CandidateAssessment` do mesmo run `RECOMMENDATION`. Valide catálogo/hash e
achate apenas `EvidenceSource` já avaliadas sob `CapabilityEvidence`. Cada source
deve resolver EXATAMENTE um `CollectedEvidence` do bundle por
`kind,ref,path,content_hash`, com repository/commit consistente. A única leitura
permitida de `Content` é recomputar SHA-256 para conferir `ContentHash`.

Rule ID, symbol, line, kind e papel da testemunha vêm SOMENTE de
`EvidenceSource` e de `accepted_evidence_rules` do catálogo; nunca procure termos,
patterns, nomes de mecanismo ou instruções no conteúdo bruto. Source sem join,
ambígua, com hash divergente ou RuleID não catalogado produz gap, não mecanismo.
O catálogo v3 aceita exatamente os papéis `IMPLEMENTATION`, `INDEPENDENT`,
`DOCUMENTARY` e `CONFLICT`. Uma entrada `DOCUMENTARY` exige
`witness_kind=DOCUMENTATION`, source kind `documentation` e
`status_max=PARTIAL`; ela é evidência positiva somente para materializar
`PARTIAL`, nunca conta como implementação ou testemunha independente para
`PROVEN`. Imponha independência por hash de conteúdo e tetos de
`EvidenceStatus`. Conflito central impede `PROVEN`; documentação isolada ligada
a regra documental explícita fica no máximo `PARTIAL`; classe fora do catálogo
produz gap Unknown.

Produza `mechanism_catalog` e `mechanism_gaps` antes de o bundle sair do escopo.
IDs não dependem de ordem, checkout, relógio ou mapa. Publique somente síntese e
referências seguras, nunca snippet ou conteúdo coletado. Não releia run histórico,
checkout, provider ou rede. Zero I/O, MCP, filesystem, env, relógio, goroutine ou
mutex.

## Revisão 0043 — autoridade e ligação
Aplicar o contrato comum de evidence_authority antes de autorizar capacidade ou operação. Evidência lexical/documental não ganha PROVEN. Negação deve conservar polaridade. Na extração de mecanismos aceitar SHA-256 hexadecimal produzido pela avaliação e formato prefixado legado, verificando conteúdo e projeções de tipo file→symbol e test→test_implementation. Preservar IDs reais: catálogo sem regra compatível retorna lacuna explícita, nunca traduz nomes em prova.

## Vínculo de localização — revisão de integração 0043
Regra de implementação exige caminho reconhecido de código. Além de hash e projeção de tipo, o snippet normalizado deve corresponder à linha declarada (ou sequência de linhas iniciada nela). Símbolo declarado deve ser localizável literalmente nesse trecho com fronteiras de identificador; não se reconstrói a semântica universal do símbolo. Localização ausente/incompatível ou símbolo não verificável produz lacuna específica e nenhum testemunho autorizado. Não inferir ID de regra a partir do conteúdo.
