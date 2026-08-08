# Especificação L0: Aho-Corasick Automaton

- **Caminho**: `01_core/parser/aho_corasick.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia nos repositórios [`gitleaks/gitleaks`](https://github.com/gitleaks/gitleaks) e [`boyter/lc`](https://github.com/boyter/lc).
- **Objetivo**: Autômato Aho-Corasick para busca e casamento multitermo simultâneo em tempo O(n+m).

## Contratos
```go
type AhoCorasick struct {}
func NewAhoCorasick(patterns []string) *AhoCorasick
func (ac *AhoCorasick) MatchAll(text string) []string
```
