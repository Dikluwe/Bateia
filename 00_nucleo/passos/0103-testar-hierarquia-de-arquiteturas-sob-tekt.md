# 0103 Testar a hierarquia de arquiteturas sob Tekt

**Estado**: EM VALIDAÇÃO — mapa inicial e primeira fatia materializados; avaliação integral da hierarquia pendente  
**Data**: 2026-10-06  
**Decisão relacionada**: [ADR-0008](../ADR/ADR-0008-hierarquia-de-arquiteturas-sob-tekt.md)

## Incerteza

É possível separar investigação, execução durável e rastreabilidade como organizações locais, preservando Tekt como arquitetura governante e permitindo reúso e insights com critérios próprios?

## Observável

Um mapa verificável de ownership e colaboração identifica, para cada responsabilidade do recorte, seu contrato A1, camada Tekt, portas usadas, invariantes, testes consumidores e conflitos atuais. O mapa deve revelar uma primeira fatia pequena de reorganização com comportamento preservável.

## Recorte inicial

Investigar a relação entre `recommendation_run.go`, `improvement_synthesis.go`, `mechanism_extraction.go`, `recommendation.go`, `idea_run.go`, container, fila, store e composição MCP. Usar rastreabilidade como consumidor de referência para verificar a distinção entre artefato, afirmação e progresso declarado.

## Procedimento

1. Congelar o commit base e registrar alterações locais relevantes. Preservar os arquivos de terceiros e a evidência histórica.
2. Mapear funções e estados ao contrato A1 e à camada L. Registrar responsabilidades ambíguas como conflito, sem forçar classificação para completar a tabela.
3. Desenhar separadamente a árvore de autoridade, o grafo de colaboração e as dependências de código. Verificar cada relação contra a gravidade Tekt.
4. Registrar entradas e saídas entre investigação e runtime: pedido, resultado, estado operacional, artefato e proveniência. Enumerar os consumidores e testes de compatibilidade.
5. Examinar pelo menos três cenários: execução concluída com conclusão inconclusiva; fonte inelegível para reúso com mecanismo útil para investigar; interrupção ou cancelamento sem promoção da hipótese.
6. Identificar o menor corte que retire uma responsabilidade de negócio do runtime ou isole síntese de seu lifecycle, com riscos e critérios de reversão explícitos.
7. Antes de implementar o corte, criar ou atualizar os prompts proprietários e seu contrato de preservação. A implementação e os testes constituem um marco posterior, não consequência automática da classificação.

## Critérios de aceite do mapa

- Tekt ocupa A0 e nenhuma regra local concede exceções às camadas.
- Os contratos A1 têm escopo, ownership e invariantes próprios; capacidades não são promovidas a arquiteturas apenas por nome.
- As relações de autoridade não são confundidas com imports ou ordem de execução.
- Estado de execução, força de evidência e estado de hipótese são distinguíveis nos três cenários.
- O mapa mostra tanto código compatível quanto mistura de responsabilidades e desconhecidos.
- A primeira fatia proposta declara comportamento preservado, consumidores, testes e reversão.

## Critérios do marco de implementação

Para o corte que vier a ser contratado: prompts e linhagem causal válidos, testes comportamentais dos consumidores, `go test ./...`, `go vet ./...`, build e `crystalline-lint --fail-on warning .`. Pendências anteriores devem ser discriminadas; não declarar conformidade completa enquanto permanecerem violações.

## Como pode falhar

A hierarquia falha se apenas renomear pacotes, duplicar contratos, criar ciclos de dependência, mover política para infraestrutura, apagar distinções probatórias ou aumentar o custo de mudança sem reduzir uma mistura concreta de responsabilidades. Se nenhum corte justificável aparecer, registrar o resultado e revisar os limites A1 antes de materializar.

## Fora de escopo

Reorganização completa de diretórios, microserviços, migração de formatos públicos, novos provedores, geração automática de insights e revisão do Tekt. Nenhum teste remoto, execução de candidatos ou benchmark de utilidade é implicitamente autorizado por este contrato.

## Evidência

Mapas e resultados brutos de bancada ficam em `lab/bancada`, quando o experimento for executado. L0 conserva o contrato e a síntese de conclusão. Gates selados e seus hashes não são reescritos.


## Materialização inicial

O [mapa de responsabilidades](../arquitetura/README.md) registra ownership, colaboração e conflitos. A primeira fatia separa análise pura em `01_core/usecase/improvement_analysis.go` do protocolo de execução em `improvement_synthesis.go`, mantendo pacote, contratos públicos e ordem de publicação. Um teste de caracterização congelado antes da extração preserva a projeção legada; testes adicionais verificam abstenção, determinismo e distinção de estado de negócio.

Esta materialização não conclui a hipótese geral do experimento, não demonstra ganho da IA e não resolve a independência de insights em relação à jornada de recomendação. A [síntese de execução](../laudos/0103-organizar-responsabilidades-sob-tekt.md) registra as verificações e os limites do corte.
