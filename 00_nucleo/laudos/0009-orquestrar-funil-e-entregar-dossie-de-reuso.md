# Laudo 0009 — Orquestrar o funil e entregar um dossiê de reúso

**Data**: 2026-08-26  
**Veredito técnico**: `IMPLEMENTAÇÃO MATERIALIZADA`  
**Veredito de produto**: `GATE PENDENTE`  
**Estado do passo**: `EM VALIDAÇÃO`  
**Regime**: A/B limitado, executado sem atestação de isolamento

## Resultado

O comando `bateia idea recommend` encadeia expansão, descoberta, screening,
planejamento e coleta de evidências, assessment e decisão. Ele opera nos modos
`live`, `local` e `replay`, sem LLM ou embeddings, e entrega:

```text
run.json
idea.json
reuse-policy.json
expansion-plan.json
discovery-snapshot.json
screening-report.json
evidence/plans.json
evidence/bundles.json
assessment-report.json
reuse-decision.json
decision.md
```

O manifesto registra hashes de entrada e saída, versões de regras, contagens,
duração por estágio, falhas e caminho da decisão. Arquivos finais usam escrita
temporária seguida de rename. Valores configurados em `GITHUB_TOKEN` e
`SOURCEGRAPH_TOKEN` são rejeitados antes da serialização.

## Achados dos ataques

Dois defeitos foram encontrados durante runs reais e congelados em regressão:

1. `--resume` aceitava um run `FAILED` que só possuía `run.json`; agora exige
   snapshot e bundles existentes antes de reutilizar evidências.
2. falha ao resolver `HEAD` gerava um bundle não imutável que derrubava o
   assessment; agora o bundle e a falha permanecem no dossiê e produzem
   assessment e decisão `INCONCLUSIVE`, sem inventar commit.

Também foram atacados: replay tocando provedor, divergência entre replays,
segredo serializado, saída inválida após chamada externa, identidade pública
inventada para worktree local, alteração local não refletida no hash e cadeia de
hashes incompleta.

## Validação executada

### Live

| Ideia | Limite | Resultado |
|---|---:|---:|
| processamento de eventos | 1 candidato, 90 s | `INCONCLUSIVE` |
| análise estática | 1 candidato, 90 s | `INCONCLUSIVE` |

As duas execuções consultaram Sourcegraph/GitHub e terminaram com dossiê. Falhas
de resolução foram preservadas; nenhuma virou `BUILD`.

### Local

| Alvo | Remoto exigido | Resultado |
|---|---:|---:|
| worktree do Bateia com alterações locais | não | `INCONCLUSIVE` |
| fixture `/tmp/bateia-local-no-remote` sem `.git` | não | `INCONCLUSIVE` |

O hash da segunda worktree mudou depois de alterar um arquivo não publicado.
Ambas receberam identidade `local/...` e URL `file://`, nunca identidade GitHub.

### Replay

Foram gerados cinco dossiês para ideias distintas. Um replay adicional sobre o
run local reproduziu byte a byte:

- `reuse-decision.json`:
  `5ae2d9e6ac4f5599bc139b2f5804ea7d27faa02d1319189bfaa3bbc6f58395c8`;
- `assessment-report.json`:
  `c7f4509a8521a5afab0aad674b6d99169d7643386179641824fe560c44c8a33a`.

## Gates técnicos

```bash
GOCACHE=/tmp/bateia-go-cache go test ./...
GOCACHE=/tmp/bateia-go-cache go vet ./...
/repos/Antigravity/tekt-linter/target/debug/crystalline-lint .
git diff --check
```

Resultados: suíte integral verde, `go vet` verde, V15 sem violações e diff sem
erros de whitespace.

## Gate de produto não satisfeito

Todos os nove dossiês reais/replay terminaram `INCONCLUSIVE`. O comportamento é
mais seguro que fabricar recomendação, mas não demonstra que 80% dos revisores
conseguiriam escolher um ponto de partida. A causa observada é cobertura de
evidência ponta a ponta insuficiente ou não correspondente à ideia do replay.

O passo não deve ser marcado `IMPLEMENTADO` até um corpus congelado conter
snapshots e bundles correspondentes às ideias e produzir ao menos quatro decisões
úteis em cinco revisões. O próximo trabalho deve ampliar esse corpus, não
afrouxar `INCONCLUSIVE` nem ajustar regra por candidato.

## Papéis e limitação de independência

- intenção e obrigação: Passo 0009 previamente escrito;
- contrato e ataques: sessão atual, congelados em arquivos próprios;
- implementação: sessão atual;
- executor e veredito: sessão atual.

Houve ordem obrigação → contrato/ataques → implementação → runs → gates, mas não
isolamento de autoridades. Este laudo é executado sem atestação de isolamento e
não alega verificação independente.
