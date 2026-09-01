# Passo 0022 — Calibrar regras de evidência por oráculos independentes

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0021 — Julgar candidatos com IA e escalonamento humano`](0021-julgar-candidatos-com-ia-e-escalonamento-humano.md)
**Onde roda**: Bateia, inteiramente offline sobre bundles congelados
**Camadas previstas**: L0, L1 e L2
**Natureza**: ampliação mensurável do catálogo de evidências sem aprender nomes de concorrentes

## Evidência que motiva este passo

Na execução E-0012, um bundle local completo do Bateia continha 40 itens, quase
96 KB e nenhuma falha, mas as dez capacidades terminaram `INCONCLUSIVE`. O
catálogo atual cobre poucas famílias e não reconhece formulações de produto como
decisão de partida, composição mínima e autofagia competitiva.

## Incerteza que este passo reduz

Não sabemos quanto das inconclusões vem de regra ausente, alvo não planejado,
padrão não reconhecido, vínculo fraco com símbolo ou falta real de evidência.

## Corpus e divisão

Usar ideias reais da família “busca e avaliação de software” e candidatos aceitos
no Passo 0021. Congelar antes da implementação:

- 4 ideias `CONTRACT`, visíveis ao autor das obrigações;
- 3 ideias `VALIDATION`, usadas pelo verificador;
- 3 ideias `TEST`, ocultas do implementador até o selo da implementação.

Cada julgamento de capacidade deve apontar `repo@commit`, arquivo, trecho ou
símbolo e estado esperado `PROVEN`, `PARTIAL`, `ABSENT`, `CONFLICT` ou
`INCONCLUSIVE`. Casos sem prova suficiente permanecem opacos.

## Política de evidência

- `PROVEN`: exige ao menos uma fonte forte em código, teste ou exemplo ligada a
  símbolo, mais uma fonte independente quando a regra declarar necessidade;
- `PARTIAL`: documentação ou sinal estrutural sustenta parte da capacidade;
- `ABSENT`: somente ausência observável em árvore completa e alvo planejado;
- `CONFLICT`: evidência negativa explícita ou incompatibilidade eliminatória;
- `INCONCLUSIVE`: regra ausente, coleta incompleta, identidade opaca, budget
  esgotado ou independência insuficiente.

É proibido usar nome, owner, domínio, slogan ou caminho exclusivo de concorrente
como regra. Regras nascem de classes de capacidade e devem ter controles
negativos semanticamente próximos.

## Implementação prevista

- tornar o catálogo versionado carregável como dado validado, mantendo L1 puro;
- mapear capacidades da IdeaSpec para famílias de regra com proveniência;
- produzir diagnóstico causal para toda evidência esperada não emitida;
- ampliar o planejamento de alvos somente quando uma regra explicar por que o
  arquivo, símbolo, teste ou exemplo pode provar a capacidade;
- preservar versões antigas para replay de dossiês existentes;
- rejeitar regra duplicada, excessivamente genérica ou dependente de identidade.

## Regime Tekt

O autor dos oráculos não lê as regras candidatas. O implementador não lê o split
`TEST`. O adversário produz mutações que removem símbolos, trocam evidência entre
repositórios, duplicam documentação, inserem README enganoso e substituem commit
por `HEAD`. O verificador não corrige regra nem oráculo.

Contrato e oráculos só são selados com mutation score `1,0`. Qualquer mudança nos
julgamentos protegidos invalida o selo e reinicia a fase afetada.

## Métricas e gates

- precisão de evidência `PROVEN` ≥ 90%;
- recall de evidência `PROVEN` ≥ 85%;
- zero promoção de caso opaco;
- 100% das evidências apontam fonte existente e commit imutável;
- duplicação de conteúdo não aumenta força ou independência;
- controles negativos não recebem `PROVEN`;
- diagnóstico classifica 100% das perdas esperadas;
- desempenho no split `TEST` satisfaz os gates sem edição posterior das regras;
- decisões antigas em replay só mudam quando a nova evidência explica a mudança.

## Veredito

- todos os gates passam nos três splits: `PROMOTE_EVIDENCE_CATALOG`;
- contrato passa e teste oculto falha: `FAIL_GENERALIZATION`;
- precisão falha: `FAIL_FALSE_PROOF`;
- recall falha: `KEEP_INCONCLUSIVE_AND_DIAGNOSE`;
- oráculo ou isolamento é violado: `INVALID_EXPERIMENT`.

## Fora de escopo

- gerar regras automaticamente por LLM;
- aprender com nomes dos candidatos;
- afrouxar `PROVEN` para elevar cobertura;
- executar binários ou testes dos repositórios;
- aconselhamento jurídico;
- decidir ainda a utilidade ponta a ponta.

## Critério de conclusão

O passo termina quando capacidades da família escolhida podem ser comprovadas com
precisão e recall mensurados em casos não vistos, preservando `INCONCLUSIVE` onde
a evidência não autoriza conclusão.
