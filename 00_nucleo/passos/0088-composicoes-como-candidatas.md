# Passo 0088 — Composições como candidatas

## Objetivo

Substituir a seleção gulosa por uma busca declarativa e limitada de conjuntos de componentes. O modo simples de pesquisa permanece disponível; a entrada avançada opcional é `bateia idea decide --composition composition.json`.

## Contrato

- Cada membro tem ID estável associado a repositório, commit e identidade da parte; o primeiro relatório fornece os IDs a usar no JSON posterior.
- A parte contribui apenas com capacidades PROVEN no assessment do mesmo repositório e fontes admissíveis.
- `connections` liga duas partes e pode ter intervalo de custo não negativo. Uma conexão PROVEN exige fonte de teste de implementação ancorada em um dos membros.
- `exclusions` proíbe co-seleção; `requirements` exige um membro quando outro é selecionado.
- `verifications` se refere ao conjunto exato de IDs e declara `status: PROVEN`, `executed: true`, `passed: true` e fonte de teste. Bateia verifica a forma, a procedência e o hash do snippet fornecido, mas **não executa o teste** nem prova a verdade material da declaração externa.
- Ausência de ligação ou ensaio deixa a composição PROPOSED e a decisão INCONCLUSIVE; conflito a torna INCOMPATIBLE. Custo ausente é desconhecido, não zero.
- A busca enumera no máximo 16 partes, conjuntos de até 6, 30.000 subconjuntos e 100 resultados. Truncamento é sinalizado e não autoriza conclusão negativa.
- A fronteira multicritério inclui composições apenas quando VERIFIED. Medições de critérios continuam explícitas e ancoradas em evidências dos membros; nada é convertido automaticamente em score.

## Exemplo mínimo de `composition.json`

```json
{
  "connections": [
    {"from": "<id-A>", "to": "<id-B>", "status": "PROVEN", "sources": ["<EvidenceSource de teste>"], "cost": {"unit": "hours", "low": 2, "high": 4}}
  ],
  "exclusions": [],
  "requirements": [],
  "verifications": [
    {"members": ["<id-A>", "<id-B>"], "status": "PROVEN", "executed": true, "passed": true, "sources": ["<EvidenceSource do ensaio>"]}
  ]
}
```

O exemplo é esquemático: cada `sources` deve conter objetos EvidenceSource, não strings. O fluxo de busca/triagem simples não requer este arquivo.

A atribuição de capacidade também exige que a parte compartilhe um testemunho com a evidência PROVEN dessa capacidade no assessment; não basta o repositório ter a capacidade em outro módulo.
