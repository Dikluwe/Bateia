# Especificacao L0: CLI Manpage Generator

- **Caminho**: `02_shell/cli/man.go`
- **Camada**: $L_2$ (Casca / CLI)
- **Objetivo**: Gerar o manual oficial UNIX bateia(1) em formato manpage.

## Contratos
```go
func GenerateManPage() string
```
