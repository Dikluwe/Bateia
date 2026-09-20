# Lab — Arena e Bancada do Bateia

`lab/` é o regime de baixo custo e alta entropia definido pelo Tekt. Nada no sistema
produtivo (`01_core` a `04_wiring`) pode importar deste diretório. O Lab não é um
depósito permanente de toda execução: conserva apenas verificadores ativos e o recorte
histórico ainda ligado ao estado declarado do produto.

## Estrutura

```text
lab/
└── bancada/
    ├── gates/       # verificadores executáveis ainda relevantes
    └── history/
        ├── benchmarks/ # evidências citadas por contratos, testes ou sínteses atuais
        ├── oracles/    # oráculos supersedidos que explicam gates selados
        ├── prompts/    # prompts supersedidos por decisão arquitetural explícita
        └── reports/    # dossiês ainda ligados ao estado atual
```

Artefatos no Lab não definem a arquitetura do produto. Uma promoção para o sistema
exige exegese, prompt proprietário em L0, materialização normalizada e testes.

`arena/` só deve existir enquanto houver experimento ativo. Ao encerrar um experimento,
promova seu resultado, preserve apenas a evidência consumida em `history/` ou elimine-o.
`legado/` só deve existir durante uma convergência brownfield real.

Compilações, caches, mutantes gerados, variações intermediárias e saídas descartáveis
usam `/dev/shm` e não são versionados. Material histórico preserva nomes cronológicos
somente para localização; ele não governa o produto atual.
