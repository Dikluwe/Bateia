# Lab — Arena e Bancada do Bateia

`lab/` é o regime de baixo custo e alta entropia definido pelo Tekt. Nada no sistema
produtivo (`01_core` a `04_wiring`) pode importar deste diretório. O Lab não é um
depósito permanente de toda execução: conserva apenas verificadores ativos e o recorte
histórico ainda ligado ao estado declarado do produto.

## Estrutura

```text
lab/
├── ferramentas/    # gate e ferramentas de desenvolvimento fixadas
└── bancada/
    ├── arquitetura/ # verificadores de fronteiras e ownership
    ├── insights/    # investigações e avaliação de insights
    ├── pesquisa/    # estudos de mecanismos, parsers e integrações
    ├── concorrencia/ # pilotos e observadores de concorrência
    ├── confiabilidade/ # controles de integridade e falhas
    ├── decisoes/    # ensaios de decisão
    ├── forjas/      # integrações com catálogos/forges
    ├── prompts/    # verificação dos contratos de exploração
    ├── corpora/    # corpora de bancada
    ├── culinaria/  # caso de investigação de domínio
    ├── gates/      # gates ativos e selados, conforme seus módulos
    └── history/    # acervo histórico referenciado; preservar hashes/caminhos
```

Artefatos no Lab não definem a arquitetura do produto. Uma promoção para o sistema
exige exegese, prompt proprietário em L0, materialização normalizada e testes.

`arena/` só deve existir enquanto houver experimento ativo. Ao encerrar um experimento,
promova seu resultado, preserve apenas a evidência consumida em `history/` ou elimine-o.
`legado/` só deve existir durante uma convergência brownfield real.

Compilações, caches, mutantes gerados, variações intermediárias e saídas descartáveis
usam `/dev/shm` e não são versionados. Material histórico preserva nomes cronológicos
somente para localização; ele não governa o produto atual.

## Entradas para uso e verificação

| Tarefa | Entrada |
| --- | --- |
| Investigar melhorias para o XERJ com o Bateia | [Caso de 8 de novembro de 2026](bancada/insights/xerj-2026-11-08/README.md) |
| Investigar princípios e decisões com evidência de fonte/alvo | [0155 — Ideias de algoritmos no XERJ](bancada/insights/algoritmos-xerj-2026-11-08/README.md): hipótese condicional e confirmação de capacidade existente, sem ganho de desempenho demonstrado |
| Acompanhar descoberta e recuperação na mesma fonte/commit | [0156 — Bateia e XERJ sobre postings](bancada/insights/bateia-xerj-postings-2026-11-08/README.md): demonstração sequencial assistida, com contexto anterior |
| Preparar contexto para a IA sem avaliações manuais | [Exploração 0136](bancada/insights/0136/README.md) |
| Entender memória e expansão de trilhas | [0142](bancada/insights/0142/README.md), [0143](bancada/insights/0143/README.md) e [0144](bancada/insights/0144/README.md) |
| Executar o gate de desenvolvimento | [Ferramentas](ferramentas/README.md) |
| Reproduzir gates históricos selados | [Regime dos gates selados](bancada/gates/sealed/README.md) |

## Bancada ativa de concorrência

[0109 — Observador sintático](bancada/concorrencia/0109/README.md) contém piloto,
CLI, corpus exploratório e ensaios. Nenhuma camada produtiva importa o pacote;
concordância do corpus não mede qualidade semântica independente.
