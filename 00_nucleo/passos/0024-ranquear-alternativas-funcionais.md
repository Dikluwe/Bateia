# Passo 0024 — Ranquear alternativas funcionais para um software alvo

**Estado**: `EM IMPLEMENTAÇÃO`
**Desejo de origem**: [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Depende de**: Passos 0005, 0008 e 0009
**Onde roda**: offline sobre snapshot e assessments já coletados

## Objetivo

Transformar o resultado do funil existente numa experiência “alternativas para
X”, sem duplicar descoberta, coleta ou regras de evidência.

## Contrato observável

- alvo possui nome e, quando conhecido, identidade `owner/repo`;
- o alvo é excluído independentemente do commit;
- cobertura conta somente capacidades `PROVEN`;
- licença incompatível, projeto arquivado proibido e operação local ausente são
  bloqueios explícitos;
- qualquer dependência decisiva `INCONCLUSIVE`, falha de coleta ou identidade
  mutável impede destaque;
- ordenação: classe de fit, cobertura, alinhamento do problema, popularidade e
  identidade, nessa ordem;
- popularidade é somente desempate posterior à evidência;
- reordenação dos inputs produz o mesmo relatório;
- primeira escolha só existe para `RECOMMENDED` ou `PARTIAL` elegível.

## Normalização e join

- `owner/repo` é comparado sem diferença de caixa; a CLI também aceita URL GitHub
  HTTP(S), forma SSH e sufixo `.git`, convertendo-os para identidade canônica;
- somente `target-name`, sem repositório, não exclui candidatos por nome;
- assessments dirigem o ranking; metadados vêm do snapshot por `owner/repo` e
  commit exatos;
- commit vazio ou `HEAD`, commit divergente, assessment ausente, snapshot ausente
  ou duplicata em qualquer lado tornam o candidato `INCONCLUSIVE`;
- órfãos do snapshot permanecem visíveis em `excluded`, nunca desaparecem;
- o próprio alvo pode aparecer apenas em `excluded` como `NOT_ALTERNATIVE`, nunca
  em `alternatives` ou `top_choice`.

## Cobertura

O denominador é o conjunto de capacidades obrigatórias não vazias da IdeaSpec,
normalizadas por caixa/espaço e deduplicadas. Somente `PROVEN` conta no numerador.
Capacidades extras são ignoradas. Capacidade obrigatória omitida ou duplicada no
assessment é `INCONCLUSIVE`; conflito explícito é bloqueio. Cobertura nunca excede
`1.0`.

## Gates normativos

- licença `INCONCLUSIVE` torna o candidato inconclusivo; `ABSENT`/`CONFLICT` ou
  identificador fora da allowlist não vazia bloqueiam; allowlist vazia desativa
  somente o filtro de identificador, não a exigência de prova da licença;
- manutenção `INCONCLUSIVE` torna o candidato inconclusivo; status conflitante ou
  arquivado quando proibido bloqueia;
- operação local só é decisiva quando `RequireLocal=true`: `INCONCLUSIVE` torna o
  candidato inconclusivo e ausência/conflito bloqueiam;
- qualquer `Failure` impede elegibilidade e produz `INCONCLUSIVE` salvo quando já
  existe bloqueio conclusivo;
- bloqueio conclusivo tem precedência sobre opacidade, pois o candidato já é
  inelegível independentemente de coleta adicional.

Todos os estados são totais:

| Campo decisivo | `PROVEN` | `PARTIAL` | `ABSENT` | `INCONCLUSIVE` | `CONFLICT` |
|---|---|---|---|---|---|
| capacidade obrigatória | conta cobertura | sinal fraco, não conta | gap conclusivo | candidato inconclusivo | bloqueado |
| alinhamento do problema | permite ranking | no máximo `REFERENCE` | `NOT_ALTERNATIVE` | inconclusivo | bloqueado |
| licença | aplica allowlist | inconclusivo | bloqueado | inconclusivo | bloqueado |
| manutenção | aplica archived | inconclusivo | bloqueado | inconclusivo | bloqueado |
| operação quando local é exigido | permite ranking | inconclusivo | bloqueado | inconclusivo | bloqueado |
| operação quando local não é exigido | ignorado | ignorado | ignorado | ignorado | ignorado |

Capacidade omitida ou duplicada equivale a `INCONCLUSIVE`. Conflito conclusivo
sempre produz `BLOCKED`; opacidade sempre produz `INCONCLUSIVE`; ambos ficam em
`excluded`.

## Derivação da classe

Aplicar na ordem abaixo:

| Condição | Classe | Lista |
|---|---|---|
| é o próprio alvo | `NOT_ALTERNATIVE` | `excluded` |
| possui bloqueio conclusivo | `BLOCKED` | `excluded` |
| identidade/join/falha/gate decisivo opaco | `INCONCLUSIVE` | `excluded` |
| alinhamento do problema é `ABSENT` | `NOT_ALTERNATIVE` | `excluded` |
| problema `PROVEN` e cobertura ≥ `ExtendMinCoverage` | `RECOMMENDED` | `alternatives` |
| problema `PROVEN` e cobertura ≥ `IntegrationMinValue` | `PARTIAL` | `alternatives` |
| problema `PROVEN`/`PARTIAL` e há ao menos uma prova ou sinal parcial | `REFERENCE` | `alternatives` |
| demais casos conclusivos | `NOT_ALTERNATIVE` | `excluded` |

`top_choice` é a primeira alternativa `RECOMMENDED` ou `PARTIAL`; fica ausente se
não houver nenhuma. `REFERENCE` nunca vira destaque.

## Ordem canônica

Em `alternatives`: `RECOMMENDED → PARTIAL → REFERENCE`, depois cobertura
decrescente, alinhamento (`PROVEN → PARTIAL`), estrelas decrescentes e identidade
canônica crescente. Em `excluded`: `BLOCKED → INCONCLUSIVE → NOT_ALTERNATIVE`,
seguido dos mesmos desempates. Estrelas nunca participam da classe ou cobertura.

Capacidades, lacunas, razões, bloqueios e fontes também são deduplicados e
ordenados. Permutar snapshot, assessments ou qualquer coleção aninhada deve gerar
JSON canônico byte a byte idêntico.

## Duplicatas e órfãos

Agrupe a união de snapshot e assessments por `owner/repo` normalizado e emita
exatamente um registro por grupo. Mais de um assessment ou mais de um candidato
de snapshot no grupo torna o registro `INCONCLUSIVE`. Assessment sem snapshot e
snapshot sem assessment também produzem um único `INCONCLUSIVE` em `excluded`.

Para representar o grupo opaco, escolha deterministicamente o menor CommitRef
lexicográfico entre os registros disponíveis e a menor metadata pela tupla
`owner, repo, commit_id, url, description, stars`; essa escolha nunca torna o
grupo elegível. Commit ou metadata divergente é explicado em `reasons`.

## Validação

O caso de uso retorna erro quando:

- target.name é vazio;
- IdeaSpec não contém capacidade obrigatória normalizada;
- policy.schema_version é vazio;
- `ExtendMinCoverage` ou `IntegrationMinValue` não são finitos ou não pertencem
  a `(0,1]`;
- `IntegrationMinValue > ExtendMinCoverage`.

`SoftwareAlternativeFilters` não é uma segunda entrada. É a projeção canônica de
`ReuseDecisionPolicy`: allowlist normalizada/deduplicada/ordenada,
`RequireLocal`, `AllowArchived`, `RecommendedMinCoverage=ExtendMinCoverage` e
`PartialMinCoverage=IntegrationMinValue`.

## Política de Unknown

Ausência comprovada pode excluir. Ausência de coleta, regra, commit ou prova de
licença/manutenção e, quando exigida, operação permanece `INCONCLUSIVE` e não é
convertida em recomendação.

## Entrega

O comando `idea recommend` aceita `--target-name` e `--target-repo` e grava:

- `alternatives.json`, fonte canônica;
- `alternatives.md`, projeção legível;
- estágio `alternatives` no manifesto do dossiê.

## API pública congelada

```go
type SoftwareAlternativeFit string

const (
    AlternativeRecommended    SoftwareAlternativeFit = "RECOMMENDED"
    AlternativePartial        SoftwareAlternativeFit = "PARTIAL"
    AlternativeReference      SoftwareAlternativeFit = "REFERENCE"
    AlternativeBlocked        SoftwareAlternativeFit = "BLOCKED"
    AlternativeInconclusive   SoftwareAlternativeFit = "INCONCLUSIVE"
    AlternativeNotAlternative SoftwareAlternativeFit = "NOT_ALTERNATIVE"
)

type SoftwareAlternativeTarget struct {
    Name       string
    Repository *RepositoryIdentity
}

type SoftwareAlternativeFilters struct {
    AllowedLicenses       []string
    RequireLocal          bool
    AllowArchived         bool
    RecommendedMinCoverage float64
    PartialMinCoverage     float64
}

func RankSoftwareAlternatives(
    target domain.SoftwareAlternativeTarget,
    idea domain.IdeaSpec,
    discovery domain.DiscoverySnapshot,
    assessments []domain.CandidateAssessment,
    policy domain.ReuseDecisionPolicy,
) (domain.SoftwareAlternativeReport, error)
```

`SoftwareAlternative` e `SoftwareAlternativeReport` preservam exatamente os
campos declarados no prompt de domínio: identidade/metadata, fit, cobertura,
alinhamento, matriz de capacidades, gaps, gates, razões, bloqueios, evidências,
listas `alternatives`/`excluded` e `top_choice`.
