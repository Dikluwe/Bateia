# Especificação L0: CLI Manpage Generator

- **Caminho**: `02_shell/cli/man.go`
- **Camada**: $L_2$ (Casca / CLI)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia no repositório [`cpuguy83/go-md2man`](https://github.com/cpuguy83/go-md2man).
- **Objetivo**: Gerar o manual oficial UNIX `bateia(1)` em formato roff para a flag `-man`.

## Contratos
```go
func GenerateManPage() string
func PrintManPage()
```

O manual apresenta `relatorio_bateia.md` como saída padrão, sem direcionar execuções
para `00_nucleo` ou para um diretório estrutural de relatórios.
