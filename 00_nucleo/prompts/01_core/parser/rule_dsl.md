# Especificação L0: Rule DSL Parser

- **Caminho**: `01_core/parser/rule_dsl.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia no repositório [`bradleyjkemp/sigma-go`](https://github.com/bradleyjkemp/sigma-go).
- **Objetivo**: Parser declarativo e avaliador de regras de linter em L1.

## Contratos
```go
type DSLRule struct {
    Name     string
    Severity string
    Pattern  string
}

func EvaluateDSL(rules []DSLRule, text string) []string
```
