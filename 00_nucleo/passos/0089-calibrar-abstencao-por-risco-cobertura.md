# Passo 0089 — Avaliar abstenção da decisão final

## Resultado inicial

`bateia idea decide evaluate --input selective-decisions.json` mede cobertura
seletiva e risco entre decisões conclusivas adjudicadas. Também separa abstenções
corretas e evitáveis. O modo simples de `search` e `idea decide` não muda.

O campo opcional `reliability` deve ser produzido para a **decisão final** por
um mecanismo futuro e validado em split de teste; não reutilizar o score da
triagem. Com todos os scores presentes, o relatório produz a curva por limiares
observados; sem eles, retorna `SCORE_UNAVAILABLE`, sem curva fictícia. O limite
`risk_upper_95` é Wilson unilateral aproximado, não uma garantia de segurança.

## Entrada mínima

```json
{
  "schema_version": "1",
  "split": "TEST",
  "cases": [
    {"id": "caso-1", "observed": "BUILD", "expected": "BUILD"},
    {"id": "caso-2", "observed": "INCONCLUSIVE", "expected": "COMPOSE"}
  ]
}
```

`expected` omitido marca caso sem adjudicação e o exclui dos denominadores,
mas o relatório o conta. `acceptable` lista outras decisões finais aceitas pelo
julgamento. `reliability` é um número entre 0 e 1 apenas para ordenar conclusões,
não uma probabilidade. Um caso `INCONCLUSIVE` continua sem conclusão, mesmo que
receba score; scores não convertem abstenções em decisões.

## Próximo experimento

O fixture ilustrativo em `lab/bancada/corpora/selective-decision-illustrative.json`
produz 3 conclusões entre 5 casos adjudicados (cobertura 60%), 1 erro nessas
3 conclusões (risco 33,3%) e 2 abstenções (1 correta, 1 evitável). Com limiar
0,8, restam 2 conclusões, nenhum erro observado e cobertura 40%; o limite
superior unilateral de Wilson ainda é cerca de 57,5%, por causa da amostra
minúscula. Esses números **não são uma medição do desempenho do Bateia**.

Congelar um corpus independente de decisões finais com rótulos humanos, grupos
e splits. Rodar o funil Bateia, capturar decisões observadas e avaliar o split
TEST. Só então ajustar uma política de abstenção em VALIDATION e reportar uma
curva de TEST sem escolher limiar nesse mesmo teste. Ainda não existe corpus
adjudicado suficiente para afirmar redução de erros do produto.
