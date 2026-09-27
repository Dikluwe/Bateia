# Passo 0085 — Priorizar investigações por valor decisório

## Estado e propósito

- Estado: implementado e verificado em 2026-09-24.
- Depende dos passos 0010, 0023, 0083 e 0084.
- Objetivo: transformar `INCONCLUSIVE` em uma escolha operacional da próxima
  investigação, sem inventar probabilidades, utilidades ou evidências.

## Entregas

1. enumerar ações recuperáveis sem alterar o plano legado;
2. projetar os resultados admissíveis de cada ação sobre uma cópia do assessment;
3. reexecutar o mesmo `DecideReuse` para cada hipótese;
4. classificar ações como `DECISION_PIVOTAL`, `DECISION_NEUTRAL` ou
   `COUNTERFACTUAL_UNSUPPORTED`;
5. expor custo vetorial, sem colapsá-lo em dinheiro ou score arbitrário;
6. selecionar a ação pivotal de menor custo observável;
7. abster e perguntar quando a ação exigir execução de código não autorizada.

## Invariantes

- uma hipótese contrafactual nunca é serializada como evidência observada;
- `EvidenceStatus`, certeza e probabilidade continuam conceitos separados;
- a ordenação não usa estrelas, popularidade ou score de screening;
- ação neutra pode continuar registrada como gate, mas não vence uma ação pivotal;
- ação não simulável não recebe valor fabricado;
- o resultado não executa código de terceiros nem concede autorização;
- o decisor produtivo e seus gates permanecem a única fronteira de decisão.

## Superfície

`idea recover value --input request.json` recebe JSON estrito e produz um plano
contrafactual. O comando legado `idea recover plan` mantém o comportamento anterior.

## Gate

Testar ação neutra, ação pivotal, ordem invariável, custo vetorial, cenário não
simulável, protocolo de autorização e ausência de evidência contrafactual no resultado.

## Resultado

- `PlanDecisionValue` reexecuta o decisor sob resultados declarados;
- custos são comparados por dominância de Pareto;
- o plano seleciona no máximo uma ação pronta por rodada;
- teste de compatibilidade sem autorização produz pergunta bloqueante;
- revisão opaca e transformações inseguras permanecem sem valor fabricado;
- a CLI está disponível em `idea recover value --input request.json`.
