# Protocolo Tekt — Passo 0012

**Regime**: protocolo completo em modo de ensaio  
**Atestação**: executado sem atestação de isolamento  
**Baseline**: `df14ec138261094a9fd33c2e887ef83db228b4a4`  
**Data**: 2026-08-26

## Fragmento observável

Esta execução cobre somente:

1. classificação interna `KEEP | CONSOLIDATE | GAP | UNKNOWN`;
2. exclusão de estados internos opacos do problema de cobertura;
3. gates de commit, falha, licença, operação e manutenção;
4. enumeração determinística da menor composição externa admissível;
5. CLI offline para materializar autofagia e Katamari.

Não cobre descoberta real dos concorrentes, qualidade de seus bundles, recall
competitivo, revisão humana de utilidade ou incorporação de capacidades.

## Entradas e saídas registradas

| Artefato | SHA-256 |
|---|---|
| passo 0012 | `cc12499b512a57bbf1de99b050686f0c8f4bb098f4ee597e2880a81b4c92364b` |
| contrato L1 | `a0ca77deaad5f5d3f7576463ad4daa62cd820332dbda122a999b1d3226866f7a` |
| implementação L1 | `f958d7cd582f6fd998cea01e0fb45541ada9a618cb34d38ee55c65c3cf8f4d54` |
| ataques L1 | `e306447e1124bba194c87b6285b75de5813817f9aa44f2ef0ee9de090c54a827` |
| CLI L2 | `878ae56878bcd1d1c6932a29ee5c5f634508f38f5b7f61053e002d06b6a8fef5` |
| ataques CLI | `fc503d07e64c2948c8c4c56d509a78f0cb29889702cd99f037d1e17d97032d5e` |

## Política de `UNKNOWN`

Ausência de assessments, coleta incompleta, conflito, falha ou evidência parcial
produz `UNKNOWN`; nunca alimenta o Katamari como `GAP`. A ausência só é gap quando
sustentada pelo assessment de um bundle completo. Um candidato externo opaco é
rejeitado quando a política não permite desconhecidos.

## Papéis e capacidades observadas

| Papel | Escrita | Limitação |
|---|---|---|
| intenção | passo 0012 | mesma sessão que solicitou materialização |
| contrato | `01_core/domain/competitive_autophagy.go` | observou a worktree existente |
| adversário | testes L1/L2 | mesma autoridade do implementador |
| implementador | usecase e CLI | pôde ler contrato e ataques |
| verificador | comandos registrados no laudo | não é autoridade independente |

As fusões impedem atestação de isolamento. Os hashes provam identidade e ordem
observada, não independência entre autoridades.

## Gate discriminatório

O teste `TestCompetitiveMutationGateRejectsEverySemanticMutation` rejeita oito
mutações semanticamente negativas:

- alvo `HEAD`;
- `MaxProjects=0`;
- `MaxComponents=0`;
- schema trocado;
- bundle interno em `HEAD`;
- conteúdo adulterado sem hash correspondente;
- item pertencente a outro repositório;
- concorrente em revisão não imutável.

```text
mutation_score = 8 / 8 = 1,0
```

Este score sela apenas o fragmento observável acima, não o experimento de produto.
