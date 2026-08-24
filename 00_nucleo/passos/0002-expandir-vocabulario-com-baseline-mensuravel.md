# Passo 0002 — Expandir vocabulário com baseline mensurável

**Estado**: `IMPLEMENTADO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0001 — Materializar IdeaSpec → ExpansionPlan`](0001-materializar-ideaspec-e-plano-de-expansao.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2 e L3  
**Natureza**: experimento comparável e primeiro incremento semântico offline

## Intenção

Fazer o Bateia formular buscas que não repetem literalmente a ideia, sem perder a
origem de cada termo e sem depender inicialmente de um modelo remoto.

O passo deve responder a uma pergunta concreta:

> Um vocabulário local, pequeno e auditável aumenta a cobertura conceitual das
> expansões sem introduzir ruído proibido?

Não queremos ainda descobrir repositórios. Queremos comparar o plano literal do
Passo 0001 com um plano semanticamente ampliado antes de misturar a qualidade da
expansão com a qualidade de Sourcegraph, GitHub ou ranking.

## Incerteza reduzida

O Passo 0001 provou que uma ideia pode ser decomposta de modo estável, mas suas
queries continuam limitadas às palavras fornecidas pelo usuário. Ainda não sabemos:

- quanto um vocabulário alternativo aumenta a capacidade de descoberta;
- quanto ruído aparece quando termos relacionados são adicionados;
- qual proveniência é suficiente para auditar uma expansão semântica;
- quais medidas permitirão comparar futuramente regras, embeddings e modelos de
  linguagem sob o mesmo contrato.

## Hipótese

Para um corpus pequeno e versionado, regras locais explícitas devem:

1. aumentar a cobertura média de conceitos esperados em relação ao plano literal;
2. manter contaminação por conceitos proibidos em zero;
3. produzir resultado determinístico e completamente offline;
4. explicar qual regra e qual expansão de origem produziram cada nova query.

Se a cobertura não aumentar, o baseline é rejeitado. Se aumentar apenas à custa de
contaminação, o baseline também é rejeitado.

## Entrega

### 1. Proveniência semântica no domínio L1

Estender o contrato sem alterar o JSON já fixado para expansões literais. Uma forma
de referência é:

```go
type ExpansionDerivation struct {
    Strategy  string   `json:"strategy"`
    RuleID    string   `json:"rule_id,omitempty"`
    SourceIDs []string `json:"source_ids"`
}

type SearchExpansion struct {
    ID         string               `json:"id"`
    Facet      ExpansionFacet       `json:"facet"`
    Query      string               `json:"query"`
    Rationale  string               `json:"rationale"`
    Covers     []string             `json:"covers"`
    Derivation *ExpansionDerivation `json:"derivation,omitempty"`
}
```

Os nomes podem mudar, mas toda expansão não literal deve declarar:

- a estratégia que a produziu;
- uma regra ou mecanismo identificável;
- uma ou mais expansões de origem;
- os campos da `IdeaSpec` que continuam cobertos.

Expansões do Passo 0001 devem permanecer byte a byte iguais no golden existente.

### 2. Porta de estratégia em L1

Criar uma abstração pequena para enriquecer um plano já materializado:

```go
type SemanticExpander interface {
    Expand(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error)
}
```

O caso de uso orquestrador deve:

- manter as expansões literais;
- pedir novas expansões à estratégia selecionada;
- normalizar e deduplicar contra o plano inteiro;
- mesclar `Covers` e origens quando estratégias gerarem a mesma query;
- preservar ordem estável;
- rejeitar expansão sem proveniência;
- impedir que `NonGoals` sejam emitidos como queries positivas;
- não modificar o plano recebido.

O núcleo não conhece arquivos, variáveis de ambiente, rede, relógio nem o formato
concreto do vocabulário.

### 3. Baseline lexical auditável em L3

Implementar uma estratégia local alimentada por um vocabulário versionado. Cada
regra precisa possuir ID estável, termos de ativação e expansões explícitas. Exemplo
ilustrativo, não definitivo:

```json
{
  "id": "structural-impact",
  "when_any": ["impacto estrutural", "impact analysis"],
  "expand_to": ["code impact analysis", "change blast radius"]
}
```

Requisitos:

- vocabulário em JSON ou YAML legível por pessoas;
- regras pequenas, revisáveis e sem chamadas externas;
- correspondência por termos normalizados, não por substring acidental;
- nenhuma regra genérica como `code → software`;
- IDs e ordem estáveis;
- limite configurável de expansões por query de origem;
- erro acionável para vocabulário inválido;
- nenhuma regra embutida silenciosamente no caso de uso L1.

O vocabulário inicial deve cobrir apenas os conceitos necessários ao corpus de
avaliação. Não tentar criar uma ontologia universal de software neste passo.

### 4. Corpus de avaliação versionado

Criar ao menos oito casos diversos, incluindo:

- análise de impacto estrutural;
- ferramenta de linha de comando;
- persistência ou banco de dados;
- interface web;
- processamento de eventos;
- análise estática;
- automação para agentes;
- um caso fora do domínio de ferramentas de desenvolvimento.

Cada caso deve conter:

```json
{
  "id": "structural-impact",
  "idea": {},
  "expected_concepts": ["code impact analysis", "change blast radius"],
  "forbidden_concepts": ["runtime monitoring", "observability platform"]
}
```

Os conceitos esperados não precisam todos aparecer literalmente na ideia. Os
conceitos proibidos devem representar vizinhos lexicalmente atraentes, mas
funcionalmente divergentes. O corpus é teste de regressão, não dado de treinamento
oculto: regras e expectativas ficam visíveis e versionadas.

### 5. Avaliador puro e relatório comparativo

Implementar métricas determinísticas para cada caso e para o corpus:

```text
coverage       = expected_concepts encontrados / expected_concepts
contamination  = forbidden_concepts encontrados / forbidden_concepts
novelty        = queries semânticas novas / queries totais
duplication    = duplicatas removidas / expansões propostas
```

Quando uma lista estiver vazia, a métrica deve ter semântica explícita e teste,
sem divisão por zero ou valor enganoso.

Adicionar uma boca de inspeção offline:

```bash
bateia idea expand --input ideia.json --strategy lexical
bateia idea evaluate --corpus caminho/do/corpus --strategy literal
bateia idea evaluate --corpus caminho/do/corpus --strategy lexical
```

`idea expand` emite o plano enriquecido. `idea evaluate` emite JSON contendo
resultados por caso, agregados e identificação da estratégia. Ambos devem funcionar
sem `.env`, tokens, cache e rede.

## Comparação obrigatória

O laudo deve comparar pelo menos:

| Estratégia | Definição |
| :--------- | :-------- |
| `literal` | comportamento implementado no Passo 0001 |
| `lexical` | plano literal acrescido do vocabulário local |

Critério inicial para aprovar o baseline lexical:

- cobertura média estritamente maior que a literal;
- contaminação total igual a zero;
- ao menos metade dos casos recebe uma query nova;
- nenhuma perda das expansões literais;
- resultado idêntico em duas execuções consecutivas.

Os números devem aparecer no laudo. Não declarar melhoria apenas por inspeção de
exemplos favoráveis.

## Prova de mordida

Usar a ideia de impacto estrutural como controle:

```text
entrada literal: mostrar a consequência estrutural de alterar código
conceitos esperados: code impact analysis, change blast radius
conceitos proibidos: runtime monitoring, application observability
```

O teste deve falhar se:

- o plano lexical não acrescentar vocabulário alternativo;
- uma expansão semântica não apontar sua regra e expansão de origem;
- um conceito proibido aparecer como query positiva;
- a ideia-alvo for modificada;
- o golden literal mudar devido à nova estratégia.

Isso evita uma implementação cosmética que apenas renomeia o expansor literal e
evita também uma expansão ampla demais que confunde impacto estrutural com
monitoramento em runtime.

## Testes obrigatórios

### Unidade — L1

1. expansões literais sobrevivem sem alteração;
2. plano de entrada não é modificado;
3. query semântica repetida é deduplicada contra literal e semântica;
4. cobertura e origens são mescladas na deduplicação;
5. expansão sem estratégia, regra ou origem é rejeitada;
6. `NonGoals` são usados como barreira e nunca como promoção positiva;
7. ordem e IDs permanecem determinísticos;
8. métricas tratam listas vazias explicitamente;
9. avaliação repetida gera JSON idêntico;
10. o caso de uso não realiza I/O.

### Integração — L2/L3

1. vocabulário válido é carregado e aplicado;
2. vocabulário inválido produz mensagem acionável;
3. limite de expansões é respeitado;
4. correspondência não ocorre por substring parcial;
5. `idea expand` funciona com a fixture do Passo 0001;
6. `idea evaluate` compara `literal` e `lexical` sobre todo o corpus;
7. os comandos funcionam com ambiente vazio e diretório de trabalho sem `.env`;
8. o golden do Passo 0001 continua passando sem alterações.

## Critérios de verificação

- [x] Proveniência de expansões semânticas é explícita e serializável.
- [x] Uma porta L1 permite trocar a estratégia sem acoplar o núcleo à L3.
- [x] O baseline lexical é local, limitado, determinístico e auditável.
- [x] O corpus possui ao menos oito ideias e vizinhos proibidos relevantes.
- [x] As métricas são reproduzíveis e possuem testes.
- [x] `idea expand` e `idea evaluate` funcionam completamente offline.
- [x] A estratégia lexical supera a literal pelos critérios declarados.
- [x] A prova de mordida preserva alvo, origem e barreiras negativas.
- [x] O contrato literal e seu golden permanecem estáveis.
- [x] Novos componentes possuem prompts L0 e linhagem válida.
- [x] `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.
- [x] Um laudo registra medidas, regras que ajudaram, ruído e lacunas.

## Como este passo pode falhar

- regras apenas memorizam o corpus e não generalizam nem para variações simples;
- cobertura sobe com termos genéricos que futuramente produzirão muitos falsos
  positivos;
- `NonGoals` são ignorados na expansão;
- proveniência existe no JSON, mas não permite reconstruir a derivação;
- a estratégia lexical altera silenciosamente o plano literal;
- o avaliador favorece a própria implementação por correspondência permissiva.

Essas falhas devem ser registradas, não escondidas ajustando o corpus depois da
execução. Mudanças nas expectativas exigem justificativa no laudo.

## Fora de escopo

- embeddings, banco vetorial ou modelo de linguagem;
- provedores remotos de expansão;
- tradução automática geral entre idiomas;
- execução das queries no Sourcegraph ou GitHub;
- medição de precisão sobre repositórios encontrados;
- agregação e ranking de candidatos;
- extração de evidência dos projetos;
- score de compatibilidade ou decisão de reuso;
- interpretação irrestrita de texto livre em `IdeaSpec`.

## Próximo passo condicionado

Se o baseline lexical for aprovado, executar suas expansões num provedor de
descoberta e medir **precisão de candidatos**, não apenas cobertura conceitual.

Se for rejeitado por baixa cobertura, manter o mesmo corpus e contrato para
comparar embeddings ou um modelo de linguagem. Se for rejeitado por contaminação,
melhorar barreiras e proveniência antes de adicionar uma estratégia mais opaca.

Nenhuma estratégia futura deve ser aceita apenas porque produz mais queries.
