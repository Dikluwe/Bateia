# Laudo 0002 — Expansão lexical com baseline mensurável

**Data**: 2026-08-24  
**Passo verificado**: [0002 — Expandir vocabulário com baseline mensurável](../passos/0002-expandir-vocabulario-com-baseline-mensuravel.md)  
**Resultado**: `APROVADO`

## O que foi implementado

- proveniência semântica com estratégia, regras e expansões de origem;
- porta `SemanticExpander` em L1 e orquestração pura de enriquecimento;
- deduplicação com mesclagem de `Covers`, regras e origens;
- barreira de `NonGoals` antes da promoção de queries;
- avaliador puro para cobertura, contaminação, novidade e duplicação;
- adaptador lexical L3 com vocabulário JSON embutido, oito regras auditáveis e
  limite de duas expansões por origem;
- corpus versionado com oito ideias, incluindo um caso fora de ferramentas de
  desenvolvimento;
- comandos offline `idea expand` e `idea evaluate`;
- teste de regressão que exige que o baseline lexical supere o literal.

## Resultado medido

As duas estratégias foram executadas sobre o mesmo corpus:

| Métrica | Literal | Lexical |
| :------ | ------: | ------: |
| Cobertura | 0% | 100% |
| Contaminação | 0% | 0% |
| Novidade | 0% | 66,67% |
| Duplicação | 0% | 0% |
| Casos com query nova | 0/8 | 8/8 |

O baseline lexical cumpriu os critérios: cobertura estritamente maior, nenhuma
contaminação, novidade em mais da metade dos casos, preservação das expansões
literais e determinismo.

## Prova de mordida

Para `mostrar o impacto estrutural de alterar código`, a estratégia lexical
preservou a query literal e acrescentou:

```text
code impact analysis
change blast radius
```

Cada expansão aponta para a regra `structural-impact` e para o ID da expansão
literal que a originou. `runtime monitoring` e `application observability`
permaneceram ausentes das queries positivas.

O golden do Passo 0001 continuou passando porque `derivation` usa `omitempty` e
não altera expansões literais.

## Execução offline

O binário foi executado em `/tmp` com ambiente vazio para `expand` e para as duas
avaliações:

```bash
env -i /tmp/bateia-step-0002 idea expand \
  --input /home/dikluwe/Documentos/Antigravity/Bateia/02_shell/cli/testdata/idea_spec.json \
  --strategy lexical

env -i /tmp/bateia-step-0002 idea evaluate \
  --corpus /home/dikluwe/Documentos/Antigravity/Bateia/02_shell/cli/testdata/evaluation_corpus.json \
  --strategy lexical
```

Ambos retornaram código zero e somente JSON em `stdout`, sem `.env`, tokens,
cache, GitHub, Sourcegraph ou rede.

## Verificações

```text
go test -race ./...    PASS
go vet ./...           PASS
go build ...           PASS
crystalline-lint .     PASS — No violations found
git diff --check       PASS
```

Os testes cobrem ainda vocabulário inválido, limite por origem, falso casamento
por substring, proveniência ausente, não mutação, deduplicação, barreiras negativas,
listas métricas vazias, saída em arquivo e estratégia desconhecida.

## Limites e interpretação

A cobertura de 100% não significa generalização semântica ampla. O vocabulário foi
deliberadamente limitado aos conceitos visíveis no corpus; portanto, o resultado
prova o contrato, a medição e a auditabilidade, mas pode superestimar desempenho
fora dessas formulações.

O casamento lexical reconhece frases normalizadas e variações explicitamente
listadas. Não realiza stemming, tradução geral, similaridade vetorial ou inferência.
Essas limitações são úteis neste baseline porque tornam falsos positivos e regras
causais inspecionáveis.

## Próxima decisão

O baseline está apto a ser executado num provedor de descoberta. O próximo passo
deve medir se as queries novas encontram projetos relevantes e diversos, mantendo
separados:

1. qualidade conceitual da expansão;
2. precisão do provedor de busca;
3. relevância e cobertura de cada candidato.

Ainda não há evidência para escolher embeddings ou LLM. Essa comparação só passa
a ser necessária quando o corpus lexical revelar baixa cobertura em ideias não
contempladas pelas regras.
