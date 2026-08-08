# Especificacao L0: Aho-Corasick Automaton

- **Caminho**: `01_core/parser/aho_corasick.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Autômato Aho-Corasick para busca multitermo simultânea em tempo O(n+m).

## Contratos
```go
type AhoCorasick struct {}
func NewAhoCorasick(patterns []string) *AhoCorasick
func (ac *AhoCorasick) MatchAll(text string) []string
```
