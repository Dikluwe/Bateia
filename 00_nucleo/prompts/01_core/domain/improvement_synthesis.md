# Improvement synthesis domain

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/improvement_synthesis.go`

Materialize os tipos, enums, tags JSON, source refs, propostas, hipóteses, passos
de validação, gaps, provenance e projeções duráveis exatos do contrato 0028.
Preserve lifecycle separado de `ImprovementSynthesisStatus`, hipótese separada de
decisão, `UNKNOWN` explícito, licença separada de evidência e proveniência por
mecanismo em composições.

IDs e hashes são valores causais canônicos; arrays públicos são `[]`. Nenhum tipo
aceita path, URL de fetch, comando, prompt, token, conteúdo bruto ou efeito sobre
software alvo. Zero I/O, MCP, filesystem, rede, env, relógio, goroutine ou mutex.

## Refinamento R3 pré-selo

O enum escalar de licença desta síntese chama-se
`ImprovementLicenseAssessment`, com identificadores Go prefixados
`ImprovementLicenseCompatible|ImprovementLicenseIncompatible|ImprovementLicenseUnrecognized`
e strings JSON preservadas `COMPATIBLE|INCOMPATIBLE|UNRECOGNIZED`. Não redeclare
nem substitua o struct basal `LicenseAssessment` do domínio de evidência.
