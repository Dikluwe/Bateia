# Passo 0023 — Validar funil completo para decisões reais de reúso

**Estado**: `DESCOBERTA APROVADA — DOWNSTREAM PENDENTE`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0022 — Calibrar regras de evidência por oráculos independentes`](0022-calibrar-regras-de-evidencia-por-oraculos-independentes.md)
**Onde roda**: descoberta autenticada uma vez; replay, evidência e decisão offline
**Camadas previstas**: L0, L1, L2, L3 e L4
**Natureza**: benchmark vertical de produto para um domínio estreito

## Objetivo

Demonstrar que o Bateia ajuda o usuário a decidir um ponto de partida para ideias
reais de busca e avaliação de software, em vez de apenas produzir candidatos ou
componentes tecnicamente válidos.

## Corpus congelado

Usar as mesmas 10 ideias dos splits `CONTRACT`, `VALIDATION` e `TEST` do Passo
0022. Para cada ideia, congelar antes do run:

- IdeaSpec e restrições eliminatórias;
- política de reúso;
- candidatos esperados quando conhecidos, sem os expor à geração;
- julgamentos de capacidade e decisões aceitáveis/proibidas;
- pergunta de utilidade que o relatório deve permitir responder.

O split `TEST` permanece oculto do implementador até o hash da implementação e
dos contratos ser selado.

## Fluxo observado

```text
IdeaSpec
  → expansão funcional
  → busca autenticada com orçamento
  → candidatos e identidades imutáveis
  → coleta de bundles
  → assessment de capacidades
  → recuperação adaptativa, no máximo duas rodadas
  → FORK · EXTEND · COMPOSE · INTEGRATE · REFERENCE · BUILD · INCONCLUSIVE
  → dossiê revisável
```

A avaliação da IA serve apenas para formar e auditar o oráculo de relevância. A
decisão operacional continua sendo calculada por regras determinísticas do
Bateia e deve citar as evidências que a sustentam.

## Decisão útil

Uma decisão é útil quando o usuário consegue responder, a partir do dossiê:

1. qual projeto ou componente deve ser considerado primeiro;
2. quais capacidades estão comprovadas;
3. quais lacunas ainda precisam ser construídas;
4. quais gates de licença, manutenção ou operação bloqueiam o reúso;
5. por que as alternativas foram preteridas;
6. qual incerteza impede decisão, quando o resultado é `INCONCLUSIVE`.

`INCONCLUSIVE` conta como útil somente em caso pré-registrado como opaco e se
indicar a evidência específica necessária para resolver a dúvida.

## Regime Tekt

- intenção, casos e decisões proibidas são congelados antes da implementação;
- julgadores não veem a decisão do motor antes de registrar o oráculo;
- implementador não lê o split `TEST`;
- adversário troca ideias, bundles, commits, evidências, políticas e identidades;
- verificador executa replay duplo e auditoria de utilidade sem escrever saídas;
- qualquer refinamento posterior cria nova versão e não reclassifica este run.

## Gates de produto

- conclusão de queries ≥ 80% em cada braço usado;
- precisão julgada no top 10 ≥ 70%;
- recall de candidatos adequados ≥ 70%;
- precisão de evidência `PROVEN` ≥ 90%;
- recall de evidência `PROVEN` ≥ 85%;
- acurácia de decisão ≥ 80%;
- taxa de decisão proibida = 0%;
- pelo menos 8 das 10 ideias produzem decisão útil;
- concordância IA/humano ≥ 80% na amostra revisada;
- replay offline byte a byte idêntico;
- zero segredo, conteúdo privado ou identidade mutável nos dossiês;
- mutation score `1,0` e arquitetura sem violações.

## Auditoria de utilidade

Após selar as decisões, o usuário responde um formulário sem ver os expected
outcomes: escolha possível, confiança, evidência suficiente, lacunas claras e
ação seguinte. O gate de utilidade exige pelo menos 80% de respostas positivas.
Casos substituídos ou reescritos depois do resultado não entram no denominador.

## Veredito

- todos os gates passam: `VALIDATED_FOR_SOFTWARE_SEARCH_IDEAS`;
- busca falha, demais fases não são avaliadas: `FAIL_DISCOVERY`;
- candidatos são bons, mas evidência falha: `FAIL_EVIDENCE`;
- evidência passa, mas decisão falha: `FAIL_DECISION`;
- decisão correta não ajuda o usuário: `FAIL_UTILITY`;
- split, cegamento ou correspondência é violado: `INVALID_BENCHMARK`.

O veredito vale somente para a família “busca e avaliação de software”. Não
autoriza afirmar que o Bateia funciona para qualquer ideia.

## Fora de escopo

- generalização para todos os domínios de software;
- integração, fork ou modificação automática de candidatos;
- repositórios privados;
- aconselhamento jurídico;
- ajuste de gates depois de observar o split `TEST`;
- terceira rodada de recuperação;
- transformar toda inconclusão em decisão.

## Critério de conclusão

O passo conclui com certificado reproduzível e veredito de produto delimitado.
Se falhar, a causa deve apontar para descoberta, evidência, decisão ou utilidade,
sem criar mais uma camada genérica de ranking.
