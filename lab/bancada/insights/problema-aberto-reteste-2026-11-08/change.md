# Correção exercitada no reteste 0153

A falha observada em 0152 tinha dois pontos distintos: material do PouchDB que permitia comparar decisões não foi desenvolvido em alternativas separadas pelo agente; outra fonte relevante nem entrou na investigação. Não se atribuiu ambas à incapacidade do mecanismo de recuperação.

O corte 0153 trata a ligação entre coleta e exploração e a revisão de opções. Não altera a relevância semântica dos catálogos, não adiciona um modelo ao produto, não promete descoberta de uma fonte específica e não embute os nomes, mecanismos ou resultados de 0152.

## Comportamento

- `idea evidence collect --question ...`: coleta pela porta existente e devolve envelope com os bundles e a exploração dos bytes fixados. Erros de resolução permanecem nos bundles; identidades sem commit imutável ficam em `unexplored_repositories`. O relatório não transforma falha em leitura.
- `idea insights explore --review-alternatives`: ativa o mesmo modo em bundles existentes. Conserva os cartões da seleção padrão; sob vagas no teto global, acrescenta até duas janelas não sobrepostas por arquivo já representado.
- As tarefas adicionais pedem comparar alternativas dentro da mesma fonte e registrar opções sustentadas, descartes e lacunas. Isso é uma instrução para interpretação pelo agente, não uma avaliação semântica automaticamente concluída.
- O modo padrão e a lista legada de `collect` permanecem. Direções continuam planos; nenhum novo acesso online é implícito.

Regras e coordenação estão em L1 (`insight_collection.go`, `insight_exploration.go`, `insight_passages.go`); L2 adapta flags e publicação. Portas e adaptadores existentes continuam realizando os efeitos. Prompts proprietários e passo foram escritos antes da materialização.

## Verificação

Testes novos verificam ligação de bytes/contexto, referências e tarefas, preservação da primeira passagem, janelas distantes não sobrepostas, teto global, ausência de mutação do pedido, cancelamento, entrada inválida, falhas sem leitura, envelope versus lista legada e proteção contra sobrescrever o plano. Não se verificam marcas ou mecanismos desejados no problema experimental.

`make verify` passou em lint, contratos arquiteturais, testes, vet, build e diff; logs em `verification`. O binário da verificação tem o mesmo SHA-256 do binário congelado para o reteste.

O replay separado em `regression-replay` usa os bytes históricos de 0152 e não foi entregue aos pesquisadores. Com o mesmo pedido, binário antigo e modo padrão novo produziram saída idêntica byte a byte. O modo de revisão passou de 12 para 16 cartões, com quatro janelas adicionais, e de cinco para sete tarefas. Isso verifica o retorno do produto; não mede a qualidade das propostas da IA.

## Desenho do reteste

O problema de negócio de 0152 e a rubrica foram preservados. Dois agentes novos receberam tarefas equivalentes, sem histórico, referências ou mecanismos esperados, baseline antes da primeira chamada, runner auditado e mesmos tetos. Os guias são específicos das interfaces. Ambos receberam exigência de revisar cobertura na entrega; o guia Bateia documenta o novo modo.

Por isso a comparação antes/depois não isola o efeito causal da mudança de código: também mudam a amostra de agentes e a exigência comum de cobertura. O ensaio verifica se a jornada corrigida é exercitada e quais decisões aparecem, com origem das fontes e limites explícitos. Nenhum projeto-fonte ou piloto é executado.
