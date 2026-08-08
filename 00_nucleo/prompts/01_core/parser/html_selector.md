# Especificação L0: HTML Documentation Selector

- **Caminho**: `01_core/parser/html_selector.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) no repositório [`PuerkitoBio/goquery`](https://github.com/PuerkitoBio/goquery) (via [`sw33tLie/bbscope`](https://github.com/sw33tLie/bbscope)).
- **Objetivo**: Seletor sintático de tags de documentação, headers e elementos HTML/XML em L1.

## Contratos
```go
func ExtractHTMLTags(htmlContent string, tagName string) []string
```
