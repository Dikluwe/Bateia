# 0133 — Foco em insights e primeira investigação do Bateia

**Data:** 2026-10-28. **Contrato anterior ao ensaio:** [passo 0133](../passos/0133-priorizar-insights-de-engenharia.md). **Evidências:** [bancada](../../lab/bancada/insights/0133/README.md).

## Entrega e escopo

O usuário reafirmou insights de engenharia como prioridade e escolheu melhorar o próprio Bateia usando outros projetos como referência. README, visão de produto e mapa passam a explicitar essa direção, preservando as capacidades e políticas de reúso existentes. Todos os arquivos Go permaneceram iguais à baseline deste passo. Nenhuma nova capacidade de geração automática de insights foi declarada como implementada.

A primeira investigação examina se o relatório atual preserva condições e diferenças de uma fonte real para orientar melhoria do Bateia. Usa o corpus conhecido 0106 de `sourcegraph/conc`, no commit `5f936abd7ae87036af1f75c95fb9d0daaf00116b`, e o código atual da fila e do gerador de insights do Bateia. O julgamento é autoral, com observações manuais anteriores; não há independência atestada, holdout ou nova consulta remota.

## Aprendizado 1 — Uma classe de mecanismo não descreve a política da implementação

Na fonte congelada, `Pool.Go` permite concorrência sem teto quando não há limiter; `WithMaxGoroutines` configura o limite explicitamente. Com limite, a submissão pode bloquear; a assinatura de `Go` não recebe contexto. O canal `tasks` é sem buffer. Os testemunhos estão em `pool/pool.go:39–65`, `89–104`, com bytes/hash na bancada.

O relatório atual classifica o mecanismo como PARTIAL, mas sua estratégia é exatamente o texto do catálogo: reserva finita e exposição de rejeição ou cancelamento. Campos de trade-offs e modelo de falha do mecanismo estão vazios. A classificação parcial é preservada corretamente; o texto genérico, sozinho, não comunica a diferença entre teto de workers, bloqueio na submissão e rejeição por capacidade. Isso é uma lacuna de utilidade frente ao foco atual, não uma alegação de que o contrato 0105 tenha executado avaliação semântica automática.

**Insight transferível:** ao investigar controle de sobrecarga, descrever separadamente o que é limitado, como o limite é configurado, o que ocorre na saturação e como se encerra uma espera. Esses detalhes orientam a escolha de estratégia melhor que o rótulo “backpressure”.

## Aprendizado 2 — Escopo do extrator precisa ser distinto de pré-condição operacional

Na saída reproduzida, `preconditions` contém `go`, `java`, `python`, `rust` e `typescript`. A inspeção de `mechanism_extraction.go` mostra que a origem é `LanguageLimits` do catálogo; `insightHypothesis` repassa esses valores. Não são condições operacionais da implementação analisada. A necessidade de configurar `WithMaxGoroutines`, por exemplo, permanece na revisão manual e não ocupa esse campo.

**Melhoria proposta:** separar metadados de cobertura do extrator de condições observadas/declaradas na fonte. Uma condição deve conservar texto, localizador, escopo e origem da avaliação. Campo desconhecido permanece desconhecido. A evolução exige contrato próprio, inclusive para preservar hashes/formatos históricos; não foi materializada neste passo.

## Aprendizado 3 — Comparar com o alvo antes de recomendar uma capacidade

A fila atual do Bateia valida workers/capacidade, reserva capacidade sob mutex e retorna `ErrQueueFull` na saturação. A execução passa um contexto cancelável ao handler; isso é cancelamento cooperativo, não garantia de encerrar código que o ignore. O ensaio executou os três testes existentes de limites/reserva, concorrência/cancelamento e reservas concorrentes; todos passaram.

O pedido de insights atual não fornece esse baseline como evidência do alvo e a saída conserva corretamente `target_state=UNKNOWN`. Uma sugestão genérica de “adicionar fila limitada” perderia uma diferença relevante: esse controle já existe no Bateia, enquanto a fonte investiga outra política de submissão. A observação do alvo foi realizada nesta revisão externa ao comando, não inferida automaticamente por ele.

**Insight transferível:** uma abordagem externa pode ser referência de contraste. A investigação deve poder concluir que uma capacidade já existe e propor apenas um experimento sobre uma diferença concreta, sem recomendar substituição por semelhança de nomes.

## Próximo corte funcional proposto

Representar observações específicas e condicionais das fontes no relatório, distinguindo-as de metadados do catálogo e hipóteses para o alvo. O primeiro aceite deve exigir que os detalhes de configuração e saturação deste caso sejam visíveis com referências, que linguagens não apareçam como pré-condições operacionais e que nenhuma avaliação fornecida seja promovida a conclusão verificada. Este caso é diagnóstico/calibração; avaliar perguntas e fontes novas antes de afirmar generalização.

## Verificação e limites

`make verify` passou: lint, auditorias, suíte raiz, vet, build e diff-check. Os hashes da baseline confirmam preservação dos arquivos Go. O binário atual reproduziu o relatório 0106 byte a byte; os controles existentes verificaram 29 itens, integridade, recusa de adulteração, preservação de saída e ausência de mecanismo sem avaliações. Esses resultados não reescrevem o recibo histórico nem promovem gates selados.

O ensaio não executou testes de `sourcegraph/conc`, não refez sua avaliação semântica de forma independente, não mediu ganho de desempenho do alvo e não implementou as melhorias propostas. O efeito no alvo continua não verificado. Sua entrega é uma investigação rastreável com aprendizados e recorte funcional específico; não uma certificação da utilidade geral do Bateia.
