# Passo 0020 — Executar busca autenticada com orçamento observável

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0019 — Executar experimento público da gramática funcional`](0019-executar-experimento-publico-da-gramatica-funcional.md)
**Onde roda**: Bateia + API de busca de repositórios do GitHub autenticada
**Camadas previstas**: L0, L1, L2, L3 e L4
**Natureza**: fronteira externa autenticada, limitada, redigida e reexecutável; sem LLM no pipeline

## Evidência que motiva este passo

O Passo 0019 enviou 36 consultas públicas, mas nenhum braço concluiu 50% delas.
O experimento terminou corretamente como `INCONCLUSIVE_RATE_LIMIT`; precisão e
recall permaneceram `null`. Repetir chamadas anônimas até obter um resultado
favorável confundiria qualidade das queries com exposição ao provedor.

## Incerteza que este passo reduz

Não sabemos se, sob exposição suficiente e igual entre braços, a gramática
funcional encontra candidatos úteis para ideias de busca e avaliação de software.
Também não sabemos se a autenticação pode ser usada sem vazar credenciais,
consultar conteúdo privado ou tornar o replay dependente da rede.

## Hipótese

Uma execução autenticada com orçamento congelado deve concluir ao menos 80% das
queries de cada braço, preservar exposição comparável e produzir snapshots que
possam ser julgados e reexecutados inteiramente offline.

## Contrato de autorização

O executor recebe o token somente por variável de ambiente. O token não entra em
arquivo de configuração, argumento de CLI, hash, log, erro, snapshot ou receipt.

Antes da rede, registrar autorização explícita com:

```json
{
  "provider": "github",
  "purpose": "public_repository_search",
  "authentication": "read_only_token",
  "private_repositories_allowed": false,
  "planned_queries": 36,
  "authorized": true
}
```

O adaptador deve rejeitar resultados marcados como privados e interromper a
persistência se qualquer artefato contiver o token, `Authorization`, cookie ou
valor derivado que permita reconstruir a credencial.

## Orçamento congelado

- três braços: N, F1 e F2;
- 12 queries por braço e 36 requests planejados;
- máximo de 10 resultados por query;
- timeout de 20 segundos por request;
- timeout total de 10 minutos;
- concorrência máxima 2;
- ordem intercalada `N1, F1-1, F2-1 ... N12, F1-12, F2-12`;
- nenhuma repetição de resposta bem-sucedida;
- no máximo uma repetição para falha transitória, contabilizada;
- parada antes de exceder o orçamento ou a margem segura declarada no manifesto.

O executor lê os headers de rate limit quando disponíveis, mas não adapta o
número de queries de um braço isoladamente. Se não houver exposição comparável,
o run permanece `UNKNOWN_PROVIDER_EXPOSURE`.

## Snapshot e receipt

Preservar por tentativa: braço, query ID, request ID, instante, status HTTP,
estado normalizado, rank, URL canônica, título, descrição, visibilidade,
request hash, response hash e contadores de rate limit redigidos.

Estados externos obrigatórios:

```text
SUCCESS · EMPTY · UNKNOWN_RATE_LIMIT · UNKNOWN_TIMEOUT
UNKNOWN_PROVIDER_RESPONSE · REJECTED_PRIVATE · SKIPPED_BUDGET
```

`Unknown` nunca equivale a conjunto vazio. O replay não lê token nem acessa rede.

## Regime Tekt

Aplicar protocolo completo:

1. autor da intenção congela queries, orçamento, política de `Unknown` e baseline 0019;
2. autor do contrato define snapshots e invariantes sem ler a implementação;
3. adversário produz tokens-canário, respostas privadas, 403, 429, timeout e adulterações;
4. mutation score `1,0` sela o contrato;
5. implementador recebe somente prompt e contrato selado;
6. verificador sem escrita nos artefatos produz receipt e certificado.

Sem isolamento demonstrável, registrar `executado sem atestação de isolamento`.

## Gates

- 100% das tentativas possuem query ID, request ID e estado observável;
- zero credencial e zero resultado privado nos artefatos;
- hashes válidos e replay duplo idêntico;
- pelo menos 80% das queries concluídas em cada braço;
- diferença de conclusão entre braços de no máximo 10 pontos percentuais;
- nenhuma estratégia recebe request extra;
- mutações de segredo, privacidade, orçamento, hash e `Unknown` são rejeitadas;
- testes, `go vet`, `staticcheck`, `crystalline-lint` e build permanecem verdes.

## Veredito

- gates de integridade e exposição passam: `READY_FOR_BLINDED_JUDGMENT`;
- cota impede exposição comparável: `INCONCLUSIVE_RATE_LIMIT`;
- token ou conteúdo privado vaza: `FAIL_SECURITY`;
- replay ou hash diverge: `FAIL_REPRODUCIBILITY`;
- orçamento é excedido: `FAIL_BUDGET`.

## Fora de escopo

- repositórios privados;
- alterar queries ou gramática depois da primeira request;
- julgar relevância;
- usar IA, embeddings ou reranking dentro do Bateia;
- trocar de provedor durante o run;
- coletar código dos candidatos;
- decidir reúso.

## Critério de conclusão

O passo termina com três snapshots autenticados, seguros, comparáveis e
reexecutáveis, ou com um veredito explícito que preserve a causa da falha.

## Evidência de execução

Executado em 2026-08-31 com token somente em memória e margem segura 10. A
primeira resposta foi `EMPTY` e informou cota restante dentro da margem; as 35
posições seguintes foram preservadas como `SKIPPED_BUDGET`, sem nova chamada.
Os três snapshots passaram em replay duplo byte a byte e a varredura não encontrou
credencial. O veredito é `INCONCLUSIVE_RATE_LIMIT`: a fronteira segura funciona,
mas não houve exposição suficiente para julgamento. Ver [Laudo 0020](../laudos/0020-executar-busca-autenticada-com-orcamento-observavel.md).
