# Oráculo independente de domínio do Passo 0028

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/architectural_mechanism_step0028_oracle_test.go`

Congele antes do candidato as assinaturas, constantes, strings e tags JSON do
contrato R3. Verifique por compilação somente as três interfaces públicas e não
imponha construtor ou tipo concreto. Preserve o struct basal `LicenseAssessment`
e exija o enum distinto `ImprovementLicenseAssessment`.

Serializações de mecanismo e hipótese devem ter arrays explícitos, provenance,
status epistemológico, licença, target state, unknowns e disclaimer de não
adoção. Rejeite superfícies com snippet, conteúdo bruto, README, checkout ou URL.
O teste é L1 puro e não usa rede, filesystem, relógio, env, goroutine ou mutex.
