# Architectural mechanism domain

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/architectural_mechanism.go`

Materialize somente os enums, dados e tags JSON de mecanismo, testemunha,
novidade, baseline convencional, gap e catálogo público definidos pela revisão
ativa de `00_nucleo/benchmark/0028/contract.md`. O catálogo fechado é
`00_nucleo/benchmark/0028/mechanism-rules.json`; não acrescente regra, alias,
heurística ou família por código.

Preserve `Observation != ArchitecturalMechanism != ImprovementHypothesis !=
AdoptionDecision`, identidade `owner/repo@commit`, arrays vazios, estados
`PROVEN|PARTIAL|INCONCLUSIVE` e `CONVENTIONAL|DIVERGENT|UNKNOWN`. O tipo não
contém snippet, arquivo bruto, body, prompt, path local, URL a seguir ou autoridade
de adoção. Zero I/O, MCP, filesystem, rede, env, relógio, goroutine ou mutex.

## Refinamento R2 pré-selo

Defina a chave segura que liga uma `EvidenceSource` já avaliada ao
`CollectedEvidence` do mesmo bundle por `kind,ref,path,content_hash`. Preserve o
`EvidenceStatus` do `CapabilityEvidence` pai em `MechanismWitness`. Rule ID,
symbol, line e kind vêm SOMENTE de `EvidenceSource`; nenhum deles é inferido de
`CollectedEvidence.Content`.
