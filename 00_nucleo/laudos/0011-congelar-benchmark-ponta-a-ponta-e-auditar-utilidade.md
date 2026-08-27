# Laudo 0011 — Benchmark ponta a ponta e auditoria de utilidade

**Veredito**: `EM VALIDAÇÃO`  
**Regime Tekt**: executado sem atestação de isolamento  
**Data**: 2026-08-26

## Resultado

Foi materializada a infraestrutura offline do benchmark:

- contratos versionados para cases, manifesto, proveniência, correspondência,
  mutações, julgamentos, métricas e resultados;
- canonicalização e hash estável independentemente da ordem de candidatos,
  bundles, itens e evidências esperadas;
- gate `PRESERVED | VIOLATED | UNKNOWN`, sem contar identidade opaca como acerto;
- mutation gate com bloqueio de certificação abaixo de `1,0`;
- replay do screening, assessment e decisão reais, sem LLM;
- métricas com denominadores vazios serializados como `null`;
- comandos `bateia benchmark validate`, `mutate` e `run` despachados antes de
  `.env`, caches e clientes remotos.

Os ataques cobrem case válido, troca de ideia, bundle/candidato, commit, conteúdo,
`HEAD`, reordenação e alteração posterior de julgamento. O conjunto selado de
doze mutações usado no teste obteve `mutation_score = 1,0`.

## Evidência executada

```text
GOCACHE=/tmp/bateia-go-cache go test ./...
PASS em todos os pacotes (fora do sandbox apenas para httptest loopback)

/repos/Antigravity/tekt-linter/target/debug/crystalline-lint .
✓ No violations found
```

## Limites encontrados

O Passo 0011 ainda não satisfaz o critério de produto:

1. não existe corpus correspondente com quinze casos e a distribuição exigida;
2. não existem oráculos humanos congelados antes de observar o motor;
3. nenhum split `TEST` foi ocultado do implementador;
4. a recuperação sobre bytes congelados ainda não integra o runner;
5. métricas de recuperação e revisão humana continuam `null`;
6. portanto, os gates de precisão, decisão e utilidade não foram avaliados.

Gerar casos sintéticos e rotulá-los pela mesma sessão produziria circularidade:
validaria o formato, mas não responderia se o Bateia ajuda alguém a escolher um
projeto real. O protocolo aplicou o critério de parada em vez de promover um
resultado artificial.

## Veredito de produto

Os Passos 0009 e 0010 permanecem `EM VALIDAÇÃO`. O motor do benchmark e sua
fronteira offline estão prontos, mas a alegação de utilidade aguarda corpus e
julgamento independentes. O próximo incremento deve materializar primeiro os
oráculos correspondentes de `CONTRACT` e `VALIDATION`; o split `TEST` só pode ser
selado por uma execução que o implementador não tenha lido.
