# Bateia — Investigação de software para insights de engenharia

**O foco atual do Bateia é produzir insights de engenharia a partir da investigação de software.** Uma pergunta, ideia ou projeto orienta a busca por estratégias, mecanismos, condições e trade-offs que ajudem a conceber ou melhorar uma solução. O [passo 0133](00_nucleo/passos/0133-priorizar-insights-de-engenharia.md) registra essa prioridade e distingue o resultado desejado das capacidades já implementadas.

O Bateia deve ampliar o repertório da IA: oferecer exemplos e perguntas que ajudem a explorar alternativas ao caminho habitual, inversões, analogias, combinações e contraexemplos. O [modo de exploração](00_nucleo/passos/0136-preparar-exploracao-criativa-para-ia.md) prepara esse contexto automaticamente a partir de uma pergunta e bundles já coletados, sem exigir avaliações ou observações escritas à mão.

A decisão entre `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE`, `BUILD` ou `INCONCLUSIVE` é uma saída. O projeto também possui extração de mecanismos e síntese de hipóteses de melhoria com premissas, lacunas e planos de validação. Uma hipótese não comprova efeito no alvo nem autoriza adoção; essas capacidades ainda têm dependências e limites descritos no mapa arquitetural.

Um insight útil explica o que foi observado numa fonte, em quais condições funciona, que custo ou limitação traz e como investigar sua pertinência à pergunta. Aprender com uma implementação não exige decidir por seu fork ou pela adoção de seu código. A versão 2 da investigação conserva observações específicas declaradas com arquivo, hash, linhas e trecho, separando condições operacionais da cobertura do extrator. A interpretação continua fornecida e o catálogo de mecanismos é limitado.

Para investigar ideias de algoritmos e outros caminhos, `--algorithm-ideas` prepara a cadeia observação → princípio/passos → condições → alvo → decisão → experimento. O comando offline `idea insights formulate` confere as citações e registra hipóteses declaradas com evidência de fonte e alvo, sem catálogo fechado de algoritmos. O [guia](00_nucleo/documentacao/investigar-software.md#investigar-princípios-de-algoritmos-e-decisões) descreve o uso; a interpretação e o efeito continuam não verificados automaticamente.

O funil **Descoberta → Triagem → Extração**, os parsers de manifestos, os coletores de código e os adaptadores AST/CFG são meios de obter evidência para essa jornada. Um repositório analisado é um candidato ou corpus de validação; não define a identidade do Bateia. Evidências de dependências sugerem aspectos a investigar, sem comprovar automaticamente capacidades ou necessidades.

A qualidade do produto deve ser avaliada pela utilidade e fundamentação das decisões em diferentes ideias, projetos e ecossistemas. Quantidade de arquivos, símbolos ou grafos extraídos e aprovação de testes de infraestrutura não substituem essa avaliação.

Intenção de origem: [D-001 — Da ideia à decisão de reúso](00_nucleo/desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [D-002 — Alternativas de software com evidência](00_nucleo/desejos/D-002-alternativas-de-software-com-evidencia.md).

## Organização arquitetural

Tekt governa o programa. O [mapa de responsabilidades](00_nucleo/arquitetura/README.md) liga capacidade, código, contrato e verificação. O [ADR-0008](00_nucleo/ADR/ADR-0008-hierarquia-de-arquiteturas-sob-tekt.md) delimita organizações locais candidatas; a primeira separação distingue análise de insights de seu protocolo de execução durável. A composição de descoberta usa um registro validado em L1 e o MCP anuncia suas fontes configuradas ([ADR-0009](00_nucleo/ADR/ADR-0009-admissao-de-conectores-e-projecao-mcp.md)). Para agentes, comece por [AGENTS.md](AGENTS.md).

## Documentação do projeto

Comece pelo [mapa de documentação](00_nucleo/documentacao/README.md), que reúne a [visão do produto](00_nucleo/documentacao/produto.md), a [arquitetura atual](00_nucleo/documentacao/arquitetura.md), o [guia de engenharia](00_nucleo/documentacao/engenharia.md) e o [plano de organização](00_nucleo/documentacao/organizacao.md).

## Comece por uma tarefa

| Quero… | Entrada |
| --- | --- |
| Usar o Bateia pela primeira vez | [Guia de investigação](00_nucleo/documentacao/investigar-software.md) |
| Encontrar referências para uma pergunta | `idea insights discover` — [primeiro uso](00_nucleo/documentacao/investigar-software.md#encontrar-referências) |
| Continuar uma trilha já encontrada | `idea insights expand` — [continuação](00_nucleo/documentacao/investigar-software.md#continuar-uma-trilha) |
| Explorar código/documentos já coletados | `idea insights explore` — [contexto para IA](00_nucleo/documentacao/investigar-software.md#explorar-material-coletado) |
| Investigar princípios de algoritmos e decisões | `--algorithm-ideas` e `idea insights formulate` — [guia](00_nucleo/documentacao/investigar-software.md#investigar-princípios-de-algoritmos-e-decisões) |
| Consultar todas as opções e exemplos | [Referência de comandos](00_nucleo/documentacao/referencia-comandos.md) |
| Desenvolver ou verificar uma mudança | [Engenharia](00_nucleo/documentacao/engenharia.md), `make verify` e [mapa de responsabilidades](00_nucleo/arquitetura/README.md) |
| Ver uma investigação aplicada | [Melhorias para o XERJ](lab/bancada/insights/xerj-2026-11-08/README.md) |

## Compilar e experimentar

Requer Go 1.25 ou superior. Na raiz:

```sh
GOCACHE=/tmp/bateia-go-cache go build -o /tmp/bateia ./04_wiring
/tmp/bateia idea insights discover \
  --question 'backpressure queue cancellation' \
  --output /tmp/descoberta-insights.json
```

O comando consulta o GitHub e tenta coletar READMEs em commits fixados. As consultas não são traduzidas: use termos na língua das fontes. O JSON contém material para a IA interpretar, tentativas e lacunas; uma execução concluída pode entregar conteúdo parcial. Veja os [limites e o passo seguinte](00_nucleo/documentacao/investigar-software.md).

## O que está disponível

| Jornada | O que entrega hoje | Limite principal |
| --- | --- | --- |
| Descoberta e expansão | Referências, leituras limitadas, memória e fronteira de investigação | Catálogo encontrado, destino lido e novidade são estados diferentes |
| Exploração | Trechos com localizadores e tarefas de raciocínio para IA | Seleção lexical/estrutural; não chama modelo nem produz avaliação semântica automática |
| Formulação de ideias | Princípio, passos, evidências de fonte/alvo, condições, decisão e teste ligados aos bytes | Interpretação declarada; não certifica novidade ou benefício |
| Analyze/compare | Organiza mecanismos e observações declaradas, hipóteses e perguntas de validação | Depende de entradas fornecidas e catálogo limitado; não demonstra efeito no alvo |
| Decisão de reúso | Dossiê e alternativas sob política explícita | Evidência insuficiente permanece inconclusiva |
| CLI, MCP e execução durável | Interfaces e lifecycle conforme os contratos de cada jornada | A existência de CLI não implica paridade completa de ferramentas MCP |

## Onde trabalhar

`00_nucleo` contém intenção e contratos; `01_core` domínio e regras; `02_shell` entrada e apresentação; `03_infra` efeitos; `04_wiring` composição. O [Lab](lab/README.md) conserva experimentos e evidências, sem ser dependência do produto.

Mudanças, caches e resultados locais não são uma publicação. Preservar segredos, formatos e evidência histórica faz parte do [ciclo de engenharia](00_nucleo/documentacao/engenharia.md). A organização documental deste recorte está registrada no [passo 0148](00_nucleo/passos/0148-organizar-entrada-e-investigar-xerj.md).
