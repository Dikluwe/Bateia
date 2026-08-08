# Especificação L0: Struct Flags CLI Parser

- **Caminho**: `02_shell/cli/struct_flags.go`
- **Camada**: $L_2$ (Casca / CLI)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) no repositório [`alecthomas/kong`](https://github.com/alecthomas/kong).
- **Objetivo**: Mapeamento declarativo de opções e subcomandos da CLI via Go struct tags em L2.

## Contratos
```go
type CLIOptions struct {
    Query string `flag:"query" help:"Query de busca"`
    Lang  string `flag:"lang" help:"Linguagem alvo"`
    Depth int    `flag:"depth" help:"Profundidade recursiva"`
}

func ParseStructFlags(target interface{}, args []string) error
```
