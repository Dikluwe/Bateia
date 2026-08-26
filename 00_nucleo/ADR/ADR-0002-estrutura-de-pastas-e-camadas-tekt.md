# ADR-0002: Nomenclatura, Estrutura de Subpastas e Mapeamento de Camadas Tekt no Bateia

**Status**: ACEITO  
**Data**: 2026-07-25  
**Documentos Relacionados**: ADR-0001 (revisada), Caso de Uso v2, Manifesto Tekt v1.3  

---

## Contexto

A primeira versão da implementação do projeto Bateia gerou código de forma amorfa, introduzindo pastas não-canônicas (como `02_usecase`) e um aninhamento desnecessário na fiação (`04_wiring/cmd/bateia/main.go`). 

Além disso, houve acoplamento direto entre as regras de orquestração do funil e os drivers de infraestrutura I/O ($L_3$), violando o Princípio da Gravidade e o isolamento de camadas da **Arquitetura Cristalina (Tekt)**.

Esta ADR formaliza a nomenclatura exata, o propósito de cada subpasta, as regras de contenção estrutural e a organização espelhada dos prompts em `00_nucleo/prompts/`.

---

## Decisão

1. Adotar a estrutura canônica Tekt de 5 estratos ($L_0$ a $L_4$), eliminando a subpasta redundante `cmd/bateia/` em favor de uma fiação direta em `04_wiring/main.go`.
2. Organizar a pasta de sementes `00_nucleo/prompts/` em **subpastas por camada**, espelhando exatamente a localização e o nome do código executável resultante (substituindo a extensão `.go` por `.md`).

### 1. Estrutura de Pastas e Mapeamento de Camadas

```text
Bateia/
├── 00_nucleo/
│   ├── ADR/
│   │   ├── ADR-0001-revisada.md
│   │   └── ADR-0002-estrutura-de-pastas-e-camadas-tekt.md
│   ├── caso-de-uso-motor-de-busca-v2.md
│   └── prompts/                 # Prompts de semente L₀ organizados em subpastas espelhadas
│       ├── 01_core/
│       │   ├── domain/
│       │   │   └── types.md
│       │   ├── parser/
│       │   │   ├── gomod.md
│       │   │   └── structure.md
│       │   └── usecase/
│       │       └── funnel.md
│       ├── 02_shell/
│       │   ├── cli/
│       │   │   └── app.md
│       │   └── report/
│       │       └── generator.md
│       ├── 03_infra/
│       │   ├── cache/
│       │   │   └── cache_store.md
│       │   ├── sourcegraph/
│       │   │   └── client.md
│       │   └── github/
│       │       └── client.md
│       └── 04_wiring/
│           └── main.md
├── 01_core/                     # L₁ — Núcleo (O Cristal: Lógica Pura & Contratos)
│   ├── domain/                  # Entidades e interfaces/contratos
│   ├── parser/                  # Parsers puros (gomod.go, structure.go)
│   └── usecase/                 # SearchFunnel (Orquestração pura dependendo apenas de L₁)
├── 02_shell/                    # L₂ — Casca (Adaptadores Primários & Apresentação)
│   ├── cli/                     # Handlers de linha de comando
│   └── report/                  # Formatação do relatório Markdown
├── 03_infra/                    # L₃ — Infraestrutura (Adaptadores Secundários & I/O)
│   ├── cache/                   # DiskCacheStore em disco
│   ├── sourcegraph/             # Client GraphQL Sourcegraph
│   └── github/                  # Client Contents API / Raw / Git Clone
├── 04_wiring/                   # L₄ — Fiação (Composição & Ponto de Entrada)
│   └── main.go                  # main() principal
└── go.mod
```

---

## Justificativa da Estrutura Espelhada de Prompts em $L_0$

A organização em subpastas espelhadas dentro de `00_nucleo/prompts/` traz os seguintes benefícios:

1. **Rastreabilidade Fisiológica Unívoca**: O prompt `00_nucleo/prompts/01_core/parser/gomod.md` é a causa biunívoca e direta de `01_core/parser/gomod.go`; `01_core/parser/gomod_test.go` possui sua própria causa em `00_nucleo/prompts/01_core/parser/gomod_test.md`. Cada prompt proprietário possui exatamente um consumidor.
2. **Navegação Intuitiva**: Desenvolvedores e agentes estatísticos encontram a especificação $L_0$ navegando no mesmo caminho de pasta da implementação.
3. **Escalabilidade**: Impede a poluição da pasta `00_nucleo/prompts/` à medida que o sistema cresce.

---

## Regras de Gravidade e Contenção

- **$L_0$ (`00_nucleo/`)**: Apenas Markdown. Zero código executável.
- **$L_1$ (`01_core/`)**: Lógica pura, algoritmos e contratos. Zero I/O. Não pode importar $L_2$, $L_3$ ou $L_4$.
- **$L_2$ (`02_shell/`)**: CLI e Apresentação. Pode importar $L_1$. Proibido importar $L_3$.
- **$L_3$ (`03_infra/`)**: Drivers e I/O. Importa $L_1$ para implementar as interfaces. Proibido importar $L_2$ ou $L_4$.
- **$L_4$ (`04_wiring/`)**: Instancia $L_3$ e injeta em $L_1$ e $L_2$. Zero lógica de negócio.
