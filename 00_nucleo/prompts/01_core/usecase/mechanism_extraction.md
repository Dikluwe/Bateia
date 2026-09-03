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
