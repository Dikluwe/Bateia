# Especificacao L0: Rule DSL Parser

- **Caminho**: `01_core/parser/rule_dsl.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Parser declarativo de regras de linter em $L_1$.

## Contratos
```go
type DSLRule struct {
    Name     string
    Severity string
    Pattern  string
}

func EvaluateDSL(rules []DSLRule, text string) []string
```
