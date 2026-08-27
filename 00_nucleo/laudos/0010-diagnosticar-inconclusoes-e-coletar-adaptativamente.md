# Laudo 0010 — Diagnosticar inconclusões e coletar adaptativamente

**Data**: 2026-08-26  
**Veredito técnico**: `IMPLEMENTAÇÃO MATERIALIZADA`  
**Veredito de produto**: `GATE AGREGADO PENDENTE`  
**Estado do passo**: `EM VALIDAÇÃO`  
**Regime**: A/B limitado, executado sem atestação de isolamento

## Resultado

O Bateia agora oferece:

```bash
bateia idea recover diagnose --run dossie
bateia idea recover plan --run dossie --diagnosis diagnosis.json --policy policy.json
bateia idea recover execute --run dossie --policy policy.json [--resume]
```

`diagnose` e `plan` são offline. `execute` recebe o provedor em L4, coleta apenas
targets congelados, une conteúdo por hash e reexecuta `AssessEvidence` e
`DecideReuse` sem alterar seus critérios.

Cada rodada é gravada em `recovery/round-NN` e o índice mantém hashes de
diagnóstico, plano, bundles, assessment, decisão e métricas. `run.json` e os
artefatos originais permanecem byte a byte inalterados.

## Ataques que alteraram a implementação

1. A ordem dos findings alterava o hash e a seleção sob `max_candidates`. O
   diagnóstico agora é canonicalizado antes de hash, limite e ordenação.
2. O histórico considerava somente rodadas anteriores e repetia `LICENSE` e
   `go.mod` do run original. A CLI agora converte os targets originais em
   histórico e escolhe apenas alternativas novas.
3. Findings perdiam o commit porque guardavam somente `owner/repo`. Agora usam
   `repo@commit`; revisão opaca permanece explicitamente `HEAD`.
4. Uma recuperação incompleta degradava um bundle original conclusivo e criava
   novos findings funcionais. O merge agora nunca reduz a completude já
   comprovada.
5. `useful_evidence_rate` podia exceder 1 ao dividir findings resolvidos por
   bundles agrupados. A métrica agora usa targets novos e possui regressão que
   exige intervalo `[0,1]`.

## Diagnóstico da amostra

Foram encontrados 28 findings nos nove dossiês:

- revisões não resolvidas;
- falhas de provedor;
- licenças desconhecidas;
- operação desconhecida;
- coleta incompleta;
- lacuna de cobertura de regra;
- conflito eliminatório.

Lacuna de regra e conflito foram mantidos terminais: nenhum deles criou regra ou
target automaticamente.

## Recuperação executada

| Grupo | Findings resolvidos | Decisão final |
|---|---:|---:|
| live — processamento de eventos | 4/8 | `INCONCLUSIVE` |
| live — análise estática | 4/8 | `INCONCLUSIVE` |
| dois runs locais | 0 | `INCONCLUSIVE` |
| cinco replays | 0 | `INCONCLUSIVE` |

Nos runs live, commits e licenças passaram a ser observáveis, mas operação e
falhas remanescentes impediram decisão conclusiva. Nenhuma falha virou `ABSENT`
ou `BUILD`.

O run local sem remoto completou duas rodadas e uma tentativa de terceira foi
rejeitada por `max_rounds = 2`.

## Métricas

- resolução live: `8/16 = 50%`;
- resolução agregada: `8/28 = 28,6%`;
- precisão de decisões alteradas: não aplicável — nenhuma decisão mudou;
- falsos `BUILD`: `0`;
- segredos serializados observados: `0`;
- repetição deliberada de target concluído: `0` após a correção do histórico.

O gate exigia resolução agregada ≥ 40% e `useful_evidence_rate` agregado ≥ 30%.
Ambos permanecem pendentes na amostra integral, embora os dois runs live tenham
atingido 50% individualmente.

## Gates técnicos

```bash
GOCACHE=/tmp/bateia-go-cache go test ./...
GOCACHE=/tmp/bateia-go-cache go vet ./...
/repos/Antigravity/tekt-linter/target/debug/crystalline-lint .
git diff --check
```

Os testes incluem decisão conclusiva sem recuperação, separação de gates,
canonicalização, deduplicação histórica, limite de rodadas, lacuna de regra
terminal, não regressão de completude, taxa limitada, CLI offline, preservação do
run e execução sem provedor rejeitada antes de mutação.

## Pendência de produto

A recuperação mostrou valor onde o problema era revisão/licença remota, mas não
resolveu operação local nem evidência funcional nos replays. A próxima iteração
deve atacar esses findings específicos:

- reconhecer evidência positiva de operação local sem confundir manifesto com
  prova de offline;
- usar bundles correspondentes à ideia no corpus replay;
- produzir targets de implementação/teste com maior novidade útil;
- medir transições de evidência, não somente redução da quantidade de findings.

Não se recomenda afrouxar `PROVEN`, `ABSENT`, `BUILD` ou os gates de licença e
operação para elevar a métrica.

## Papéis e limitação de independência

- obrigação: Passo 0010 previamente congelado;
- contratos e ataques: sessão atual, em arquivos próprios;
- implementação, execução e veredito: sessão atual.

Houve ordem obrigação → contrato/ataques → implementação → experimento → gates,
mas não houve isolamento de autoridades. Este resultado foi executado sem
atestação de isolamento e não alega verificação independente.
