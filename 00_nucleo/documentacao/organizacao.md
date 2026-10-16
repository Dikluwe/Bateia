# Plano de organização do Bateia

A primeira etapa é tornar o funcionamento atual legível por capacidade. O projeto já tem camadas, contratos e acervo de decisões; reorganização física só deve ocorrer quando reduzir um acoplamento identificado e preservar comportamento por testes.

## Organização aplicada

O [mapa arquitetural](../arquitetura/README.md) liga responsabilidades e contratos. A primeira fatia isola análise pura de hipóteses do protocolo de execução. O passo 0104 corrige admissão e projeção do registro de fontes; o passo 0105 acrescenta investigação offline de insights sem decisão de reúso. A linhagem das adições locais de pesquisa e Bitbucket foi regularizada preservando seu comportamento. O [passo 0103](../passos/0103-testar-hierarquia-de-arquiteturas-sob-tekt.md) registra o recorte; trabalhos restantes continuam abaixo.

## Entrega documental inicial

Esta pasta fornece uma visão de produto, um mapa de arquitetura e um guia de engenharia. O [diagnóstico inicial](../relatorios/diagnostico-organizacao-2026-10-06.md) distingue condições verificadas e pendências. Os próximos itens são propostas de trabalho, não ADRs aceitos nem implementação concluída.

## Sequência proposta

| Ordem | Trabalho | Por que priorizar | Critério de aceite |
| --- | --- | --- | --- |
| 1 | Linhagem e prompts das adições locais regularizados | Contratos e cabeçalhos adicionados | Especificações coerentes, cabeçalhos e lint sem erros; testes dos conectores passam |
| 2 | Ampliar o mapa inicial por capacidade | O recorte atual cobre os fluxos centrais | Cada capacidade principal liga contrato, fachada, portas, testes, evidência e limites |
| 3 | Uniformizar estados de entrega e validação | Implementação e qualidade demonstrada possuem critérios diferentes | Cada capacidade informa disponibilidade, verificação e avaliação de produto separadamente |
| 4 | Tornar a verificação repetível em ambiente limpo | Testes dependem de loopback e configuração Git da máquina | Procedimento reproduzível para testes, lint, vet e build, com execução de CI se adotada |
| 5 | Delimitar efeitos e determinismo em L1 | Uso de relógio conflita com a intenção de pureza | Exegese e contrato para uma primeira fatia, seguido de testes de preservação |
| 6 | Avaliar persistência do dossiê e fronteiras de apresentação | A CLI concentra arquivos, checkpoints e projeções | Decisão de ownership registrada e migração pequena com compatibilidade comprovada |
| 7 | Avaliar e reduzir despacho/composição em L4 | Configuração e features ampliam o ponto de entrada | Mapa de responsabilidades e corte justificado sem mover política para adaptadores |
| 8 | Ampliar avaliação da decisão de reúso | A promessa depende da utilidade para ideias reais | Corpus e julgamentos declarados, resultados por etapa e abstenções explicadas |

## Inventário por capacidade

Para cada capacidade, preencher:

- Resultado para o usuário e interface de entrada.
- Contrato proprietário e decisão arquitetural aplicável.
- Fachada, tipos, portas e adaptadores efetivamente usados.
- Artefatos produzidos, identidade e versões dos formatos.
- Testes ativos e evidência de avaliação com recorte temporal.
- Limites conhecidos, falhas parciais e dúvidas ainda abertas.

Começar pela recomendação ponta a ponta. Depois registrar descoberta de mercado/pesquisa, coleta e comprovação, rastreabilidade/revisão e runtime durável. Essa ordem acompanha a promessa central do produto.

## Política de mudança

Preservar nomes históricos, hashes e gates selados. Antes de mover um pacote ou arquivo, registrar o problema concreto, dependências afetadas, comportamento preservado e critério de reversão. Novos agrupamentos por capacidade podem inicialmente ser documentais; a estrutura física só muda depois da exegese e do contrato.

As regras de agentes existem em dois arquivos. Avaliar uma referência canônica com compatibilidade para as ferramentas consumidoras, sem apagar uma das cópias antes de conhecer quem a utiliza.
