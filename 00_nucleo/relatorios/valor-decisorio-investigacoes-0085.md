# Valor decisório das investigações — passo 0085

## Resultado

O Bateia agora avalia mecanicamente se uma ação de recuperação pode atravessar uma
fronteira real do decisor. A implementação é uma primeira camada de valor da informação:
contrafactual e exata em relação às regras atuais, mas deliberadamente não bayesiana.

Para cada ação, o motor:

1. parte do assessment já admitido e sanitizado;
2. aplica cada resultado declarado somente em memória;
3. reexecuta o mesmo `DecideReuse`;
4. registra as decisões alcançáveis;
5. classifica a ação como pivotal, neutra ou não simulável;
6. compara custos por dominância de Pareto;
7. seleciona no máximo um target executável para a próxima rodada.

Hipóteses nunca aparecem como `EvidenceSource`. O resultado contém estados e decisões
contrafactuais, não evidência observada.

## Superfície

```text
bateia idea recover value \
  --input decision-value-request.json \
  --output decision-value-plan.json
```

O comando recebe `IdeaSpec`, assessment, política decisória, diagnóstico, política de
recuperação, rodada, histórico de targets e autorização de execução. O JSON é estrito.
`idea recover plan` e `idea recover execute` mantêm o comportamento anterior.

## Interpretação

- `DECISION_PIVOTAL`: ao menos um resultado muda `INCONCLUSIVE` para uma decisão;
- `DECISION_NEUTRAL`: todos os resultados declarados mantêm a decisão atual;
- `COUNTERFACTUAL_UNSUPPORTED`: não há transformação local segura do assessment;
- `ACTION_READY`: um target existente pode seguir para a recuperação;
- `AWAITING_CANDIDATE_EXECUTION_AUTHORIZATION`: a melhor opção exige autorização;
- `ACTION_REQUIRES_EXECUTOR`: o plano identificou um ensaio, mas este passo não executa
  código de terceiros;
- `NO_DECISION_PIVOTAL_ACTION`: coletar os targets disponíveis não resolve a decisão.

## Custos

O custo permanece vetorial: chamadas de provedor, documentos, arquivos de código,
testes e execução de candidato. Uma opção só é marcada como dominada quando outra não é
pior em nenhuma dimensão e é estritamente melhor em ao menos uma. Não existe conversão
oculta em dinheiro ou score único.

## Limites

- não há EVPI, EVPPI ou EVSI numérico sem probabilidades e utilidades calibradas;
- revisão imutável não resolvida permanece não simulável porque exige nova coleta;
- teste de compatibilidade é planejado e pede autorização, mas requer um executor
  isolado futuro;
- o valor observado deverá ser calibrado posteriormente com os resultados reais das
  rodadas de recuperação.
