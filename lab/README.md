# Lab — Arena e Bancada do Bateia

`lab/` é o regime de baixo custo e alta entropia definido pelo Tekt. Nada no sistema
produtivo (`01_core` a `04_wiring`) pode importar deste diretório.

## Estrutura

```text
lab/
├── arena/       # experimentos e protótipos ainda não promovidos
├── bancada/     # gates, mutantes, benchmarks, corpora e medições
│   └── history/ # resultados passivos e acervo histórico da bancada
└── legado/      # substrato brownfield read-only, quando existir
```

Artefatos no Lab não definem a arquitetura do produto. Uma promoção para o sistema
exige exegese, prompt proprietário em L0, materialização normalizada e testes.

Material histórico movido de L4 preserva nomes cronológicos somente para localização.
Ele não governa o produto atual.
