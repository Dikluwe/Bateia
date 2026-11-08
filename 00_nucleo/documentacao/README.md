# Documentação do Bateia

O Bateia investiga software para produzir oportunidades de reúso e insights de engenharia. Este mapa conecta propósito, arquitetura e trabalho de engenharia sem exigir a leitura cronológica de todos os passos.

## Arquitetura antes da documentação

O [mapa de responsabilidades](../arquitetura/README.md) e o [ADR-0008](../ADR/ADR-0008-hierarquia-de-arquiteturas-sob-tekt.md) orientam a organização atual. Os textos desta pasta começaram como rascunho documental e não constituem a especificação completa da arquitetura. Reúso é uma saída do produto; insights possuem contratos e limites próprios.

## Por onde começar

| Necessidade | Leitura |
| --- | --- |
| Entender o problema e o resultado esperado | [Visão do produto](produto.md) |
| Fazer a primeira investigação | [Investigar software](investigar-software.md) |
| Localizar responsabilidades, fluxos e interfaces | [Arquitetura atual](arquitetura.md) |
| Escrever ou revisar especificações de componentes | [Guia de prompts](escrever-prompts.md) e [contratos consolidados](../contratos/README.md) |
| Desenvolver, verificar e preservar contratos | [Guia de engenharia](engenharia.md) |
| Planejar a organização e reduzir riscos conhecidos | [Plano de organização](organizacao.md) |
| Consultar comandos e exemplos | [Referência de comandos](referencia-comandos.md) |
| Consultar o diagnóstico desta leitura | [Diagnóstico de 6 de outubro de 2026](../relatorios/diagnostico-organizacao-2026-10-06.md) |

## Onde cada informação pertence

| Informação | Local | Como interpretar |
| --- | --- | --- |
| Propósito e intenção | [Desejos](../desejos/README.md) | O resultado desejado; não atesta implementação |
| Contrato de uma evolução | [Passos](../passos/README.md) | Escopo, observável e critérios de aceite |
| Especificação de um componente | `00_nucleo/prompts/` | Causa declarada da materialização |
| Decisão transversal | `00_nucleo/ADR/` | Contexto, decisão e consequências |
| Síntese de execução e limites | `00_nucleo/laudos/` e [relatórios](../relatorios/README.md) | Resultado relativo ao recorte e à versão descritos |
| Funcionamento corrente | Código em `01_core` a `04_wiring` | Confirmar com testes e contratos aplicáveis |
| Corpus, fixtures, gates e resultados brutos | [Lab](../../lab/README.md) | Evidência de bancada com identidade histórica |
| Orientação transversal de leitura | Esta pasta | Síntese mantida junto às mudanças do produto |

As regras estruturais estão em [.agentrules](../../.agentrules) e [.agents/AGENTS.md](../../.agents/AGENTS.md). Esta documentação explica sua aplicação e as diferenças observadas; não substitui essas regras nem os ADRs aceitos.

## Manutenção

Ao mudar uma capacidade, revise a seção correspondente da arquitetura e os comandos afetados. Registre decisões novas em ADR e contratos em passos/prompts. Preserve laudos, hashes e recibos históricos; acrescente uma síntese atual quando sua interpretação mudar.
