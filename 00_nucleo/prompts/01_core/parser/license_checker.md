# Especificação L0: License Checker

- **Caminho**: `01_core/parser/license_checker.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na análise de governança e conformidade de projetos open-source minerados.
- **Objetivo**: Identificador estático de licenças de código (Permissiva vs Copyleft/GPL) e avaliação de riscos.

## Contratos
```go
type LicenseCategory string
const (
    Permissive LicenseCategory = "PERMISSIVE"
    Copyleft   LicenseCategory = "COPYLEFT"
    Unknown    LicenseCategory = "UNKNOWN"
)

func CategorizeLicense(licenseText string) LicenseCategory
```
