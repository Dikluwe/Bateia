# Passo 0019 — Executar experimento público da gramática funcional

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0018 — Versionar gramática funcional e congelar busca web reexecutável`](0018-versionar-gramatica-funcional-e-congelar-busca-web-reexecutavel.md)
**Onde roda**: Bateia + API pública do GitHub
**Camadas previstas**: L0 e artefatos de benchmark; código somente se um defeito for demonstrado antes da execução
**Natureza**: experimento externo autorizado, congelado e reexecutável, sem LLM no pipeline

## Por que este passo é separado

O Passo 0018 refez a gramática, o adaptador observável, os hashes, o replay e as
métricas. Sua primeira execução externa, porém, foi interrompida antes da rede:
o processo poderia carregar um `GITHUB_TOKEN` do `.env` e enviar 36 consultas
derivadas do projeto.

Alterar o Passo 0018 depois de selado invalidaria sua cadeia causal. Este passo
preserva os contratos e executa somente o experimento que ficou pendente.

## Incerteza que este passo reduz

Não sabemos se a gramática funcional v2.1 melhora a descoberta real de
repositórios em relação:

- aos rótulos nominais da ontologia;
- às consultas funcionais v1 do Passo 0017.

Também não sabemos se a API pública do GitHub possui recall suficiente para
esse trabalho ou se a vantagem anteriormente observada dependia de um índice
web mais amplo.

## Hipótese

Mantidos provedor, orçamento, janela e corpus, o braço F2 deve superar N e F1
em precisão e recall repo-backed sem usar nomes esperados na geração.

```text
N  — rótulos nominais
F1 — queries funcionais congeladas no Passo 0017
F2 — queries da gramática v2.1 congeladas no Passo 0018
```

## Autorização externa

Antes da primeira chamada, registrar uma autorização contendo:

```json
{
  "scope": "36 repository-search queries",
  "provider": "GitHub public repository search API",
  "token_policy": "NO_TOKEN",
  "payload": "query text, locale and result limit",
  "result_policy": "public repository metadata only",
  "authorized": true
}
```

Regras:

- não carregar `.env`;
- não enviar `GITHUB_TOKEN`, cookie ou credencial;
- não consultar repositório privado;
- não enviar conteúdo de código ou arquivos locais;
- enviar somente os textos congelados das 36 queries;
- interromper se o processo tentar adicionar autenticação;
- sanitizar receipts antes da persistência.

Sem autorização explícita, o estado permanece `BLOCKED_BY_APPROVAL`.

## Entradas congeladas

Registrar caminho e SHA-256 antes da rede:

```text
lab/bancada/history/benchmarks/0018/queries-n.json
lab/bancada/history/benchmarks/0018/queries-f1.json
lab/bancada/history/benchmarks/0018/queries-f2.json
03_infra/lexical/functional-grammar-v2.json
lab/bancada/history/benchmarks/0018/segregacao/selo.json
```

O executor encerra antes da rede se:

- qualquer hash divergir do manifesto do experimento;
- N, F1 ou F2 não contiver exatamente 12 queries;
- houver ID duplicado;
- uma query estiver vazia;
- F1 divergir do artefato congelado do Passo 0017;
- F2 falhar na auditoria gramatical;
- aparecer nome, owner ou domínio proibido.

## Registry protegido

Antes da rede, o executor pode conhecer somente:

- a quantidade de expected outcomes;
- IDs opacos como `expected-1` a `expected-4`;
- a lista de termos proibidos usada pela auditoria.

Nome, URL e owner dos concorrentes ficam em arquivo separado, ilegível durante:

1. congelamento das queries;
2. execução dos três braços;
3. selagem dos snapshots;
4. deduplicação offline.

O registry é aberto somente para produzir julgamentos depois que os três hashes
de snapshot estiverem registrados.

## Desenho do experimento

### Orçamento

- 12 queries por braço;
- máximo de 5 resultados por query;
- 36 requisições no total;
- uma query por requisição;
- timeout de 20 segundos por requisição;
- timeout total de 15 minutos;
- nenhuma repetição automática de resposta bem-sucedida;
- no máximo uma repetição de falha transitória, registrada separadamente.

### Ordem intercalada

Executar por posição para reduzir efeito temporal:

```text
N1, F1-1, F2-1,
N2, F1-2, F2-2,
...
N12, F1-12, F2-12
```

A ordem interna dos braços é determinada pelo ID canônico. Não reordenar por
qualidade percebida.

### Rate-limit

Sem token, o rate-limit é uma possibilidade esperada.

- HTTP 403 ou 429 vira `UNKNOWN_RATE_LIMIT`;
- preservar request ID, query ID, instante e status HTTP;
- não converter em conjunto vazio;
- não esperar ocultamente até a janela reiniciar;
- não trocar de token ou provedor durante o run;
- não comparar recall de braços com exposição de rede desigual.

Se menos de 50% das queries de qualquer braço concluírem, o experimento é
`INCONCLUSIVE_RATE_LIMIT`.

## Snapshot obrigatório

Cada resposta deve preservar:

```text
run_id
arm
request_id
query_id
query text ou request hash verificável
provider
retrieved_at
status
HTTP status
rank
URL original
URL canônica
título
snippet/descrição
indicador de anúncio
request_hash
response_hash
receipt sanitizado
```

Cada braço produz um snapshot separado. O hash do snapshot é calculado somente
depois de todas as respostas, inclusive falhas e `Unknown`.

## Replay antes do julgamento

Executar duas vezes, sem rede:

```bash
bateia idea web replay \
  --snapshot snapshot-n.json \
  --output replay-n-a.json

bateia idea web replay \
  --snapshot snapshot-n.json \
  --output replay-n-b.json
```

Repetir para F1 e F2. Os pares devem ser byte a byte idênticos.

O replay deve verificar:

- hashes de snapshot e resposta;
- associação request/query;
- rank e campos materiais;
- deduplicação por URL canônica;
- melhor posição preservada;
- separação de anúncios;
- `Unknown` e rate-limit preservados.

Falha de replay invalida o braço inteiro.

## Julgamento cego

Produzir uma união das URLs canônicas dos três braços, embaralhada por hash.
O julgador não recebe:

- braço de origem;
- query que encontrou a URL;
- posição original;
- métrica parcial.

Para cada URL, registrar:

```json
{
  "canonical_url": "...",
  "judgment": "RELEVANT | IRRELEVANT | UNKNOWN",
  "is_software": true,
  "is_repo_backed": true,
  "expected_id": "expected-2",
  "reasons": ["..."],
  "evidence": ["title", "description", "repository URL"]
}
```

`expected_id` aparece somente quando a identidade corresponde ao registry
aberto após os snapshots. Repositório relevante novo permanece relevante sem
receber expected ID.

## Métricas

Calcular por braço:

```text
query_completion_rate
unknown_rate
unique_organic_hits
judged_precision_at_10
relevant_repository_yield
repo_backed_competitive_recall
expected_ids_found
functional_over_nominal_recall_gain
f2_over_f1_recall_gain
control_contamination
```

Regras de denominador:

- expected outcomes vêm dos IDs explícitos do registry;
- denominador vazio permanece `null`;
- `UNKNOWN` não conta como relevante nem irrelevante;
- anúncios não entram no top 10 orgânico;
- a mesma expected ID conta uma vez por braço;
- ausência por rate-limit não vira falso negativo silencioso.

## Gates

### Integridade

- 36 requests planejados e rastreáveis;
- zero credenciais em request ou receipt;
- zero alteração das queries depois da rede;
- 100% das respostas com request ID e query ID;
- hashes válidos em 100% dos snapshots;
- replay duplo byte a byte idêntico;
- registry aberto somente depois dos snapshots.

### Cobertura externa

- pelo menos 50% das queries concluídas em cada braço;
- nenhuma falha de um braço é registrada como sucesso;
- diferença de conclusão entre braços ≤ 20 pontos percentuais;
- caso contrário: `INCONCLUSIVE_PROVIDER_EXPOSURE`.

### Produto

- precisão julgada em 10 de F2 ≥ 60%;
- recall repo-backed de F2 ≥ 70%;
- F2 encontra ao menos 3/4 expected IDs;
- ganho de F2 sobre N ≥ 30 pontos percentuais;
- precisão de F2 não piora mais de 5 pontos contra F1;
- nenhum controle negativo recebe `RELEVANT`.

## Decisão

- todos os gates passam: `PROMOTE_GRAMMAR_V2`;
- F2 melhora F1, mas não passa recall: `KEEP_VALIDATING_GRAMMAR_V2`;
- N, F1 e F2 equivalentes: `INDEX_DOMINATES_QUERY_FORMULATION`;
- F2 piora F1: `REJECT_GRAMMAR_V2`;
- exposição desigual, rate-limit ou replay inválido: `INCONCLUSIVE`;
- GitHub falha e outro índice amplo havia funcionado: escrever novo passo para
  comparar provedores; não trocar provedor durante este run.

Nenhuma decisão altera automaticamente a gramática. Promoção exige atualizar o
estado e o laudo em revisão posterior, preservando os snapshots.

## Comandos previstos

Executar por um orquestrador que intercale os braços. Os comandos abaixo
representam cada braço, mas não autorizam execução sequencial N → F1 → F2:

```bash
env -u GITHUB_TOKEN -u GH_TOKEN \
  bateia idea web discover-observable \
  --queries lab/bancada/history/benchmarks/0018/queries-n.json \
  --arm N --max-per-query 5 \
  --output lab/bancada/history/benchmarks/0019/snapshot-n.json

env -u GITHUB_TOKEN -u GH_TOKEN \
  bateia idea web discover-observable \
  --queries lab/bancada/history/benchmarks/0018/queries-f1.json \
  --arm F1 --max-per-query 5 \
  --output lab/bancada/history/benchmarks/0019/snapshot-f1.json

env -u GITHUB_TOKEN -u GH_TOKEN \
  bateia idea web discover-observable \
  --queries lab/bancada/history/benchmarks/0018/queries-f2.json \
  --arm F2 --max-per-query 5 \
  --output lab/bancada/history/benchmarks/0019/snapshot-f2.json
```

## Testes obrigatórios antes da rede

1. `.env` não é carregado no modo `--no-token`;
2. headers não contêm `Authorization`;
3. 12 queries existem em cada braço;
4. IDs são únicos dentro e entre requests;
5. F1 corresponde ao hash do Passo 0017;
6. F2 passa auditoria gramatical;
7. termos proibidos bloqueiam antes da rede;
8. execução intercala N/F1/F2;
9. timeout vira `UNKNOWN_TIMEOUT`;
10. 403/429 vira `UNKNOWN_RATE_LIMIT`;
11. erro de parser vira `UNKNOWN_PROVIDER_RESPONSE`;
12. resposta vazia bem-sucedida permanece distinguível de falha;
13. request e response hashes detectam adulteração;
14. token e cookies são removidos do receipt;
15. replay não acessa rede;
16. replay duplo é idêntico;
17. anúncio não entra no orgânico;
18. URL canônica deduplica sem fundir owners distintos;
19. julgamento não contém braço ou query;
20. denominador vazio gera `null`;
21. expected ID repetida conta uma vez;
22. snapshot incompleto não pode promover;
23. exposição desigual produz `INCONCLUSIVE`;
24. todos os gates são necessários para promoção;
25. testes completos, build, linter V15 e `git diff --check` permanecem verdes.

## Artefatos esperados

```text
lab/bancada/history/benchmarks/0019/authorization.json
lab/bancada/history/benchmarks/0019/manifest.json
lab/bancada/history/benchmarks/0019/snapshot-n.json
lab/bancada/history/benchmarks/0019/snapshot-f1.json
lab/bancada/history/benchmarks/0019/snapshot-f2.json
lab/bancada/history/benchmarks/0019/replay-n.json
lab/bancada/history/benchmarks/0019/replay-f1.json
lab/bancada/history/benchmarks/0019/replay-f2.json
lab/bancada/history/benchmarks/0019/blinded-candidates.json
lab/bancada/history/benchmarks/0019/blinded-judgments.json
lab/bancada/history/benchmarks/0019/experiment.json
lab/bancada/history/benchmarks/0019/receipt.json
00_nucleo/laudos/0019-executar-experimento-publico-da-gramatica-funcional.md
```

## Fora de escopo

- usar token, cookie ou repositório privado;
- usar LLM, embeddings ou reranker neural no Bateia;
- alterar queries, gramática ou expected outcomes durante o run;
- trocar de provedor para melhorar o resultado;
- repetir chamadas até obter ranking favorável;
- interpretar rate-limit como irrelevância;
- promover a gramática com snapshots sintéticos;
- fabricar julgamento ou identidade repo-backed;
- executar Sourcegraph antes de resolver identidades;
- corrigir código depois da primeira request e continuar o mesmo run.

## Critério de conclusão

O passo termina quando os três braços reais estão congelados e reexecutáveis e
o veredito pode ser recalculado offline. Resultado ruim conclui o experimento;
resultado não observável ou exposição desigual permanece inconclusivo.
