# Passo 0028 — Extrair mecanismos e formular hipóteses de melhoria

**Estado**: `PLANEJADO`  
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)  
**Depende de**: [`Passo 0027 — Generalizar descoberta multilíngue e multifuente`](0027-generalizar-descoberta-multilingue-e-multifuente.md)  
**Baselines obrigatórios**: certificados `MCP_RECOMMENDATION_GATEWAY_READY`, `MCP_DURABLE_RUNS_READY` e `MULTISOURCE_MARKET_DISCOVERY_READY`  
**Onde roda**: extração determinística em L1, execução durável pelo MCP e síntese opcional por IA externa  
**Camadas previstas**: L0, L1, L2 e L4; L3 somente pelos adaptadores já certificados  
**Natureza**: análise de engenharia orientada por evidência, sem modelo de IA embutido e sem modificar software alvo

## Problema

O Bateia já consegue descobrir repositórios, comprovar capacidades, comparar
alternativas e encontrar produtos em outros idiomas e mercados. Isso responde
“o que existe?”, mas ainda não responde de forma segura à pergunta que motivou o
produto:

> que mecanismo encontrado fora da solução habitual pode melhorar a ideia que
> estamos construindo, sem confundir inspiração, evidência e decisão de adoção?

Uma IA externa consegue ler os resultados e sugerir melhorias, mas tende a cair
em dois extremos. Pode repetir a solução convencional porque ela é frequente no
treino, ou pode produzir uma combinação criativa sem conseguir provar que as
partes existem, resolvem a lacuna ou são compatíveis entre si.

O Bateia também possui conceitos próximos, mas insuficientes para essa resposta:

- `ReusableComponent` identifica partes reaproveitáveis de um candidato;
- `AutophagyReport` separa capacidades internas mantidas, consolidadas, ausentes
  ou desconhecidas;
- `KatamariProposal` procura a menor composição capaz de cobrir lacunas provadas;
- `ImprovementRecommendation` legado usa frequência de adoção entre similares.

Nenhum deles representa explicitamente um mecanismo arquitetural transferível,
uma hipótese causal de melhoria, suas premissas, contraevidências e plano de
validação. Frequência de adoção também não prova qualidade: uma solução rara
pode ser justamente a observação útil, e uma solução comum pode apenas repetir
o estado atual.

## Incerteza que este passo reduz

Não sabemos se é possível transformar resultados congelados do Bateia em
hipóteses de melhoria úteis para engenharia sem:

- declarar arquitetura a partir de nome, descrição ou metadado de loja;
- converter sugestão da IA em fato;
- promover novidade, popularidade ou frequência como qualidade;
- chamar estado desconhecido de lacuna;
- misturar evidência de repositórios diferentes como se fosse uma implementação;
- recomendar cópia quando licença ou compatibilidade não foram verificadas;
- reconstruir o ranking, a autofagia competitiva ou o Katamari;
- expor código de terceiros ou conteúdo não confiável pelo MCP;
- permitir que uma hipótese escreva, instale ou execute algo no software alvo.

## Objetivo

Materializar duas coisas diferentes e causalmente ligadas:

1. um catálogo de **mecanismos arquiteturais observados**, extraídos de evidência
   repo-backed já coletada pelo funil;
2. um conjunto de **hipóteses de melhoria**, propostas por regras ou por uma IA
   externa e avaliadas contra mecanismos, lacunas, premissas, riscos e
   contraevidências.

```text
RECOMMENDATION run com evidência congelada
  → extrator determinístico e versionado
  → mechanism_catalog

MARKET_DISCOVERY run
  → observações de produto sem alegação de arquitetura

IA externa ou regras puras
  → propostas não confiáveis
  → bateia_start_improvement_synthesis
  ← run_id

worker durável do Passo 0026
  → valida hashes e proveniência
  → compara mecanismos com a lacuna alvo
  → procura premissas, riscos e contraevidências
  → improvement_hypotheses + validation_plans

IA externa ou pessoa
  → decide o que investigar, prototipar ou rejeitar
```

O Bateia organiza e qualifica a base factual. A IA externa pode atravessar
domínios, explicar analogias e propor composições; ela não recebe autoridade
para mudar o estado da evidência nem para aprovar adoção.

## Separações epistemológicas obrigatórias

### Observação

Fato limitado a uma superfície congelada: um arquivo, símbolo, teste, manifesto,
documento, item de catálogo ou mercado consultado. Observação não é mecanismo.

### Mecanismo

Estratégia implementada para produzir um efeito, sustentada por evidência de
código e testemunhas independentes dentro do mesmo repositório imutável.
Exemplos de forma, não de catálogo obrigatório: fila com backpressure, log
append-only, cache por conteúdo, sincronização offline com resolução de conflito
ou extensão por adaptadores.

### Hipótese de melhoria

Proposição falsificável de que adaptar um ou mais mecanismos pode melhorar uma
lacuna ou propriedade do alvo, sob premissas explícitas. Mesmo quando
`SUPPORTED`, continua sendo hipótese, não recomendação de adoção.

### Decisão de adoção

Escolha humana ou de processo posterior que considera implementação, custo,
licença, segurança, operação e validação no contexto real. Fica fora deste passo.

O schema e a projeção pública não podem usar esses termos como sinônimos.

## Recorte da primeira materialização

Este passo entrega:

1. vocabulário e schema versionados de mecanismo arquitetural;
2. extrator determinístico baseado no `EvidenceBundle` já coletado;
3. `mechanism_catalog` seguro produzido durante novos runs `RECOMMENDATION`;
4. referências imutáveis a artefatos de runs anteriores;
5. propostas de melhoria `RULE_DERIVED` ou `EXTERNAL_AGENT_PROPOSAL`;
6. síntese como novo job no runtime durável do Passo 0026;
7. estados epistemológicos separados de lifecycle e recomendação;
8. comparação com baseline convencional congelado, sem premiar novidade;
9. plano de validação e lacunas para cada hipótese preservada;
10. benchmark com mecanismos, hipóteses, negativos e replay offline.

Runs `RECOMMENDATION` históricos que não possuam `mechanism_catalog` não são
reinterpretados retroativamente. Devem produzir `SOURCE_BASIS_UNAVAILABLE` ou
ser executados novamente sob o contrato 0028. O Bateia não relê um checkout
móvel para completar um run antigo.

## Contrato de mecanismo arquitetural

L1 define um dado puro equivalente a:

```go
type MechanismEvidenceStatus string

const (
    MechanismProven       MechanismEvidenceStatus = "PROVEN"
    MechanismPartial      MechanismEvidenceStatus = "PARTIAL"
    MechanismInconclusive MechanismEvidenceStatus = "INCONCLUSIVE"
)

type MechanismNovelty string

const (
    NoveltyConventional MechanismNovelty = "CONVENTIONAL"
    NoveltyDivergent    MechanismNovelty = "DIVERGENT"
    NoveltyUnknown      MechanismNovelty = "UNKNOWN"
)

type ArchitecturalMechanism struct {
    ID                string                  `json:"id"`
    Name              string                  `json:"name"`
    Purpose           string                  `json:"purpose"`
    Trigger           string                  `json:"trigger,omitempty"`
    Strategy          string                  `json:"strategy"`
    Components        []ReusableComponent     `json:"components"`
    DataFlow          []string                `json:"data_flow"`
    StateModel        []string                `json:"state_model"`
    FailureModel      []string                `json:"failure_model"`
    Boundaries        []string                `json:"boundaries"`
    Preconditions     []string                `json:"preconditions"`
    Tradeoffs         []string                `json:"tradeoffs"`
    Source            ImmutableCandidateRef   `json:"source"`
    Evidence          []EvidenceSource        `json:"evidence"`
    CounterEvidence   []EvidenceSource        `json:"counter_evidence"`
    Status            MechanismEvidenceStatus `json:"status"`
    Novelty           MechanismNovelty        `json:"novelty"`
    ExtractionRuleIDs []string                `json:"extraction_rule_ids"`
}
```

Nome e forma exatos podem mudar no contrato de L0, mas as distinções semânticas
são obrigatórias.

Um mecanismo `PROVEN` exige, no mesmo `owner/repo@commit`:

- testemunha de implementação ligada a path e símbolo ou estrutura equivalente;
- ao menos uma testemunha independente, como teste, exemplo executável,
  manifesto, schema ou configuração causalmente relacionada;
- regra de extração versionada que explique a transformação;
- evidência suficiente para propósito, estratégia e limite alegados;
- ausência de contraevidência que contradiga a alegação central.

Documentação isolada pode sustentar `PARTIAL`, nunca `PROVEN`. Nome de pacote,
tag, README promocional, estrelas ou quantidade de projetos não bastam. Um
produto sem repositório pode originar observação de produto e inspirar proposta,
mas nunca um `ArchitecturalMechanism` provado.

O extrator não executa código encontrado. Ele usa somente o bundle congelado no
run corrente, com limites de tamanho e tipos de evidência já permitidos.

## Catálogo e regras de extração

O catálogo de regras é fechado antes do candidato e cada regra declara:

- ID e versão;
- classe de mecanismo reconhecida;
- testemunhas mínimas;
- campos que consegue preencher;
- condições de conflito;
- limites de linguagem e ecossistema;
- status máximo que pode conceder.

Regras podem reconhecer topologia, fluxo de dados, coordenação de estado,
persistência, extensibilidade, isolamento, recuperação e políticas operacionais.
O catálogo inicial deve ser pequeno e comprovável. Um mecanismo fora do
vocabulário permanece desconhecido; não é forçado ao rótulo mais próximo.

Extração ocorre dentro do novo run `RECOMMENDATION`, enquanto seu
`EvidenceBundle` imutável ainda está disponível. Somente a representação segura
do mecanismo é publicada. O MCP não publica excerto integral, arquivo bruto,
corpo de página ou conteúdo usado apenas como sinal negativo.

O ID do mecanismo deriva canonicamente de source imutável, regra, alegação e
hashes das testemunhas. Ordem de descoberta, caminho do checkout ou horário não
participam do ID.

## Novidade sem culto à novidade

Cada caso do benchmark possui um `ConventionalBaseline` congelado antes de
observar os candidatos. Ele enumera soluções convencionais esperadas para a
ideia apenas para fins de comparação.

- `CONVENTIONAL`: o mecanismo pertence à família já presente no baseline;
- `DIVERGENT`: mecanismo provado fora dessas famílias;
- `UNKNOWN`: não há base suficiente para classificar.

`DIVERGENT` significa diferença, não superioridade. Não altera evidência, score,
fit, licença ou status da hipótese. O relatório deve mostrar tanto mecanismos
convencionais fortes quanto divergentes relevantes e explicar a diferença.
Frequência no corpus é contexto descritivo, nunca multiplicador de qualidade.

O `MineOntologicalImprovements` e o `ImprovementRecommendation` legado não são
fontes de evidência ou de suporte neste pipeline. Podem permanecer por
compatibilidade, marcados como legado, mas nenhum valor de `AdoptionRate` pode
entrar no cálculo do 0028.

## Contrato de hipótese de melhoria

```go
type ImprovementHypothesisStatus string

const (
    HypothesisProposed     ImprovementHypothesisStatus = "PROPOSED"
    HypothesisSupported    ImprovementHypothesisStatus = "SUPPORTED"
    HypothesisConflicting  ImprovementHypothesisStatus = "CONFLICTING"
    HypothesisInconclusive ImprovementHypothesisStatus = "INCONCLUSIVE"
)

type TransferMode string

const (
    TransferComponentReuse  TransferMode = "COMPONENT_REUSE"
    TransferPatternAdapt    TransferMode = "PATTERN_ADAPTATION"
    TransferReferenceOnly   TransferMode = "REFERENCE_ONLY"
    TransferUnknown         TransferMode = "UNKNOWN"
)

type ImprovementHypothesis struct {
    ID                 string                      `json:"id"`
    Title              string                      `json:"title"`
    TargetProperty     string                      `json:"target_property"`
    TargetGapRefs      []string                    `json:"target_gap_refs"`
    ExpectedEffect     string                      `json:"expected_effect"`
    MechanismRefs      []string                    `json:"mechanism_refs"`
    ProductObservationRefs []string                `json:"product_observation_refs"`
    TransferMode       TransferMode                `json:"transfer_mode"`
    Assumptions        []string                    `json:"assumptions"`
    Preconditions      []string                    `json:"preconditions"`
    IntegrationUnknowns []string                   `json:"integration_unknowns"`
    Risks              []string                    `json:"risks"`
    Tradeoffs          []string                    `json:"tradeoffs"`
    EvidenceRefs       []string                    `json:"evidence_refs"`
    CounterEvidenceRefs []string                   `json:"counter_evidence_refs"`
    ValidationPlan     []ValidationStep            `json:"validation_plan"`
    Derivation         string                      `json:"derivation"`
    Status             ImprovementHypothesisStatus `json:"status"`
}
```

Uma hipótese só pode ser `SUPPORTED` quando:

1. o mecanismo referido é `PROVEN` ou a política congelada aceita explicitamente
   `PARTIAL` com teto epistemológico menor;
2. a propriedade ou lacuna do alvo tem testemunha própria;
3. existe relação causal explícita entre mecanismo, lacuna e efeito esperado;
4. premissas, riscos, tradeoffs e contraevidências foram procurados e registrados;
5. referências resolvem para hashes imutáveis do mesmo pedido;
6. não há conflito conhecido que derrube a hipótese;
7. o modo de transferência respeita a informação de licença disponível.

`SUPPORTED` não significa “implementar”, “melhor alternativa” ou “sem risco”. A
projeção deve incluir essa ressalva de forma legível por máquina.

Quando duas fontes são combinadas, cada mecanismo conserva sua proveniência. A
composição não prova que interfaces, estado, performance, licenças ou modelos de
falha sejam compatíveis. Esses itens aparecem em `integration_unknowns` e no
plano de validação. Seleção efetiva de componentes continua pertencendo ao
Katamari do Passo 0012.

## Lacuna e propriedade alvo

Uma lacuna `ABSENT` comprovada pelo `AutophagyReport` pode sustentar uma hipótese
direcionada. `PARTIAL` pode sustentar hipótese de alteração quando existe
baseline observável do comportamento atual. `UNKNOWN` nunca é convertido em
ausência.

Quando não existe repositório alvo ou relatório interno congelado, o agente pode
formular uma hipótese de produto, mas seu estado máximo é `PROPOSED`. O plano de
validação deve começar pela observação do estado atual.

Uma hipótese também pode buscar uma propriedade não ausente, como reduzir
acoplamento, latência ou perda de dados. Nesse caso é obrigatório indicar a
medição basal que falta ou a testemunha que demonstra o problema atual.

## Licença e modo de transferência

Licença desconhecida não elimina mecanismo, candidato ou hipótese. Produz aviso
`LICENSE_NOT_RECOGNIZED` e impede somente a promoção automática para
`COMPONENT_REUSE`.

- licença reconhecida e compatível pode permitir `COMPONENT_REUSE`, sem aprovar
  a adoção por si só;
- licença desconhecida preserva `PATTERN_ADAPTATION`, `REFERENCE_ONLY` ou
  `UNKNOWN`, conforme o contrato, com revisão humana obrigatória;
- incompatibilidade comprovada bloqueia reúso direto e registra o conflito;
- referência abstrata a uma estratégia não autoriza copiar código, texto,
  assets, interface protegida ou implementação derivada;
- o Bateia não emite parecer jurídico.

Essa política refina a triagem sem retornar ao comportamento de excluir
silenciosamente projetos cuja licença não foi reconhecida.

## Referências imutáveis de entrada

Síntese não recebe path local nem texto afirmando o resultado de outro run. Cada
fonte é referida por identidade e hash:

```json
{
  "run_id": "run_opaque_identifier",
  "run_kind": "RECOMMENDATION",
  "artifact": "mechanism_catalog",
  "sha256": "64-hex-chars"
}
```

Para observações de mercado, `artifact` é `market_candidates` ou
`market_snapshot` seguro. O serviço resolve a referência no mesmo store,
verifica lifecycle terminal, kind, integridade e autorização local. Hash
divergente, run ausente, artefato corrompido ou fonte não terminal falha antes da
fila.

O pedido inteiro recebe hash causal. Trocar um source hash com a mesma
`idempotency_key` produz conflito, preservando o contrato do Passo 0026.

## Ferramenta MCP

As nove ferramentas certificadas no Passo 0027 permanecem nas mesmas posições.
Uma ferramenta é acrescentada ao final:

10. `bateia_start_improvement_synthesis`.

Alterar ordem, remover ferramenta ou inserir uma nova antes dela invalida a
compatibilidade observável e exige refinamento explícito dos certificados 0025,
0026 e 0027.

### `bateia_start_improvement_synthesis`

Entrada representativa:

```json
{
  "idempotency_key": "improvement-task-2026-09-02-0001",
  "idea": {
    "problem": "reduzir desperdício alimentar doméstico",
    "users": ["famílias"],
    "capabilities": ["inventário", "previsão de consumo"]
  },
  "source_artifacts": [
    {
      "run_id": "run_recommendation",
      "run_kind": "RECOMMENDATION",
      "artifact": "mechanism_catalog",
      "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    },
    {
      "run_id": "run_market",
      "run_kind": "MARKET_DISCOVERY",
      "artifact": "market_candidates",
      "sha256": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
    }
  ],
  "proposals": [
    {
      "proposal_id": "agent-proposal-1",
      "title": "usar consumo previsto para priorizar alertas",
      "target_property": "diminuir alertas irrelevantes",
      "expected_effect": "ordenar lembretes pelo risco estimado de perda",
      "mechanism_refs": ["mechanism-id"],
      "product_observation_refs": [],
      "assumptions": ["há histórico local suficiente"],
      "derivation": "EXTERNAL_AGENT_PROPOSAL",
      "rationale": "mecanismo observado trata prioridade como estado derivado"
    }
  ],
  "policy": {
    "max_hypotheses": 12,
    "max_mechanisms_per_hypothesis": 4,
    "require_counterevidence_search": true
  }
}
```

Regras:

- ideia reutiliza a normalização canônica existente;
- ao menos uma fonte imutável é obrigatória;
- no máximo 12 fontes, 20 propostas e os limites menores configurados pela
  política congelada;
- IDs referidos precisam existir nos artefatos declarados;
- título, rationale, premissas e efeito têm limites e rejeitam controles;
- proposta externa é conteúdo não confiável e recebe essa proveniência;
- `mechanism_refs` não pode criar, alterar ou elevar um mecanismo;
- produto sem repo entra somente em `product_observation_refs`;
- ausência de proposta externa permite hipóteses geradas por regras puras;
- política não aceita pesos de popularidade, prompt, URL, path ou expressão de
  ferramenta;
- nenhum campo solicita escrita, execução, instalação ou acesso privado.

Resposta:

```json
{
  "schema_version": "1",
  "run_id": "run_opaque_identifier",
  "run_kind": "IMPROVEMENT_SYNTHESIS",
  "lifecycle": "QUEUED",
  "version": 1
}
```

Fila, lease, retry, idempotência, cancelamento, store, quota e publicação atômica
são os mesmos do Passo 0026. Não existe worker ou banco paralelo para síntese.

## Artefatos duráveis

Novos runs `RECOMMENDATION` acrescentam à allowlist segura:

- `mechanism_catalog`;
- `mechanism_gaps`.

Runs `IMPROVEMENT_SYNTHESIS` publicam:

- `architecture_comparison`;
- `improvement_hypotheses`;
- `hypothesis_gaps`;
- `validation_plans`;
- `synthesis_provenance`.

`bateia_get_run` passa a aceitar `run_kind=IMPROVEMENT_SYNTHESIS` e o campo
`improvement_status` com `COMPLETE`, `PARTIAL` ou `INCONCLUSIVE`. Lifecycle
`SUCCEEDED` indica somente que o job terminou e publicou artefatos íntegros; não
significa que uma hipótese foi suportada.

`architecture_comparison` associa propriedades do alvo a mecanismos sem achatar
suas fontes. `hypothesis_gaps` distingue falta de evidência, desconhecimento do
alvo, licença, compatibilidade, orçamento e proposta inválida. Artefatos seguem
paginação, cursor HMAC, hash, quota e projeção segura já certificados.

## Geração por regras e propostas externas

O caminho `RULE_DERIVED` pode combinar lacuna comprovada com mecanismos cujo
propósito declarado cobre a mesma capacidade ou propriedade, usando taxonomia e
regras congeladas. Ele não usa modelo, embedding online ou similaridade opaca.

O caminho `EXTERNAL_AGENT_PROPOSAL` permite analogias que as regras não previram.
O Bateia valida referências e reclassifica o resultado pela evidência disponível.
Uma proposta pode ser preservada como `PROPOSED`, rebaixada a `INCONCLUSIVE`,
marcada `CONFLICTING` ou elevada a `SUPPORTED`; sua autoridade externa nunca é
copiada para a evidência.

Texto de produto pode sugerir uma pergunta, nunca responder como o mecanismo foi
implementado. Se a proposta depender apenas de produto não repo-backed, o plano
de validação deve pedir uma fonte técnica, experimento ou observação verificável.

## Contraevidência e plano de validação

Para cada hipótese, o serviço procura ao menos:

- evidência de que a lacuna não existe ou já foi resolvida;
- limites declarados do mecanismo fonte;
- falhas e testes negativos relacionados;
- incompatibilidades de estado, interface ou operação conhecidas;
- restrições de licença disponíveis;
- premissas sem testemunha.

Ausência de contraevidência encontrada é registrada como busca limitada, não
como prova de ausência.

Cada `ValidationStep` possui pergunta falsificável, observável, procedimento
seguro, critério de aprovação, critério de rejeição, custo relativo e evidência
que falta. Neste passo o plano é apenas materializado: não cria branch, não gera
código, não roda experimento no alvo e não abre PR.

## Corpus congelado e oráculos

Congelar antes da implementação ao menos oito casos distribuídos por quatro
famílias de software:

- dois mecanismos provados por código + teste;
- um mecanismo documentado sem testemunha suficiente;
- um produto não repo-backed que inspira uma proposta;
- uma composição de dois mecanismos com compatibilidade desconhecida;
- uma lacuna interna `UNKNOWN`, que não pode virar ausência;
- uma licença não reconhecida, que deve gerar aviso sem excluir;
- uma contraevidência que deve impedir `SUPPORTED`.

Cada caso possui:

- `IdeaSpec` e, quando aplicável, estado alvo congelado;
- bundles e artefatos com hashes imutáveis;
- `ConventionalBaseline` escrito antes de revelar candidatos;
- catálogo esperado de mecanismos e teto epistemológico;
- hipóteses admissíveis, conflitantes e inválidas;
- plano mínimo de validação;
- split `CONTRACT`, `VALIDATION` ou `TEST`.

O split `TEST`, as mutações e os rótulos finais permanecem ocultos do
implementador até o selo. Julgadores de utilidade veem relatórios sem nomes dos
projetos e respondem se conseguem entender: o que adaptar, de onde veio, por que
poderia ajudar, o que pode falhar e como testar.

## Métricas e gates de produto

- precisão de mecanismos `PROVEN` = 100%; qualquer falso provado bloqueia;
- recall de mecanismos esperados ≥ 80% no conjunto `VALIDATION + TEST`;
- 100% dos mecanismos preservam source imutável, regras e testemunhas;
- 100% das hipóteses preservam derivação e referências resolvíveis;
- taxa de hipótese `SUPPORTED` sem lacuna/propriedade alvo testemunhada = 0%;
- taxa de produto-only promovido a mecanismo arquitetural = 0%;
- taxa de `UNKNOWN` promovido a ausência = 0%;
- taxa de composição promovida como compatível sem testemunha = 0%;
- 100% das hipóteses preservadas possuem premissas, riscos e plano de validação;
- 100% das buscas de contraevidência declaram escopo e resultado;
- licença não reconhecida preserva candidato e produz aviso em 100% dos casos;
- nenhum `AdoptionRate`, estrela, download ou ranking altera status;
- recall de mecanismos divergentes provados ≥ 80% quando o oráculo os contém;
- falso rótulo `DIVERGENT` = 0%; divergência não muda o ranking;
- utilidade cega ≥ 80% nas cinco perguntas do relatório;
- replay offline reproduz todos os artefatos byte a byte;
- mutation score `1,0`.

Recall de novidade só é julgado nos casos em que o oráculo contém mecanismo
divergente provado. Ausência de mecanismo divergente não reprova a execução e
jamais autoriza inventar um.

## Política de `Unknown`

- run fonte não terminal: referência inválida antes da fila;
- artefato histórico sem catálogo: `SOURCE_BASIS_UNAVAILABLE`;
- hash divergente ou artefato corrompido: erro de integridade;
- regra sem testemunhas mínimas: mecanismo `INCONCLUSIVE` ou gap;
- documento sem implementação: no máximo `PARTIAL`;
- produto sem repo verificado: observação, nunca mecanismo;
- lacuna alvo desconhecida: hipótese no máximo `PROPOSED`;
- baseline de propriedade ausente: efeito não mensurável ainda;
- compatibilidade entre mecanismos desconhecida: não afirmar integração;
- licença não reconhecida: aviso e bloqueio de reúso direto, sem exclusão;
- contraevidência não procurada por budget: hipótese `INCONCLUSIVE`;
- proposta da IA sem referência válida: rejeitada como entrada, não alucinada;
- ausência de hipótese suportada: resultado válido e possivelmente
  `INCONCLUSIVE`, nunca erro fabricado;
- mecanismo não coberto pela taxonomia: desconhecido, sem classificação forçada.

## Segurança e conteúdo não confiável

- nenhum conteúdo de repositório, produto ou proposta controla tools, políticas,
  prompts, paths ou próximas chamadas;
- mecanismo publica síntese estrutural e referências, não arquivo ou excerto
  integral de terceiro;
- nomes, descrições, rationales e mensagens são limitados, sanitizados e tratados
  como dados;
- URLs não são seguidas durante síntese;
- source refs resolvem somente dentro do store e da allowlist locais;
- o serviço não clona, baixa, instala, compra ou executa candidatos;
- o serviço não escreve no repositório alvo;
- artefatos não contêm segredo, token, review, dado pessoal ou prompt privado;
- erros públicos não ecoam conteúdo bruto, canários nem configuração;
- licenças, avisos e conflitos não podem ser removidos por proposta externa;
- conteúdo contraditório não é apagado para tornar a hipótese mais limpa.

## Observáveis a congelar antes da implementação

1. as nove tools certificadas preservam ordem, schema e comportamento;
2. `tools/list` acrescenta somente `bateia_start_improvement_synthesis` ao final;
3. novo run `RECOMMENDATION` produz catálogo a partir do bundle corrente sem
   reler checkout móvel;
4. run histórico sem catálogo permanece indisponível para extração;
5. documentação isolada nunca produz mecanismo `PROVEN`;
6. produto-only nunca produz mecanismo arquitetural;
7. path e símbolo de uma fonte não se combinam com teste de outra;
8. ordem do bundle não altera IDs nem catálogo;
9. regra e testemunhas aparecem na proveniência de todo mecanismo;
10. proposta externa não eleva o status do mecanismo referido;
11. `UNKNOWN` no alvo não vira gap;
12. ausência de target permite proposta, mas não `SUPPORTED`;
13. composição preserva incompatibilidades desconhecidas;
14. licença desconhecida mantém a hipótese e bloqueia somente reúso direto;
15. licença incompatível não desaparece do relatório;
16. popularidade e taxa de adoção não alteram nenhum status;
17. divergência é visível, mas não recebe bônus de score;
18. contraevidência rebaixa ou conflita a hipótese conforme o oráculo;
19. ausência de contraevidência é descrita como busca limitada;
20. source hash divergente falha antes da fila;
21. mesma idempotency key e mesmas fontes retorna o mesmo run;
22. mesma key com fonte ou proposta diferente produz conflito;
23. cancelamento impede novos passos e preserva artefatos já publicados;
24. lifecycle `SUCCEEDED` não é projetado como hipótese `SUPPORTED`;
25. replay offline não chama rede e é byte a byte idêntico;
26. canários de código, produto e proposta não atravessam artefatos ou logs;
27. nenhuma saída inicia implementação ou altera o alvo;
28. certificados 0025, 0026 e 0027 recebem recibos explícitos de refinamento.

## Mutações mínimas a rejeitar

- README promocional sozinho gera mecanismo `PROVEN`;
- nome de produto vira arquitetura implementada;
- proposta da IA cria evidência ou testemunha inexistente;
- teste de um repo comprova código de outro;
- `owner/repo` sem commit imutável entra no catálogo;
- estado `UNKNOWN` do alvo vira `ABSENT`;
- mecanismo frequente ganha suporte adicional;
- mecanismo divergente é automaticamente preferido;
- `AdoptionRate` entra no score ou status;
- ausência de contraevidência vira prova de segurança;
- duas partes provadas viram composição compatível sem evidência;
- licença desconhecida exclui candidato ou hipótese;
- licença desconhecida permite `COMPONENT_REUSE` automático;
- licença incompatível é omitida da projeção;
- source hash informado pelo cliente não é verificado;
- run histórico relê HEAD atual para inventar catálogo;
- lifecycle `SUCCEEDED` vira `SUPPORTED`;
- artifact MCP inclui arquivo, trecho integral ou canário;
- proposta contém URL/path e provoca acesso externo;
- síntese usa fila ou store paralelos;
- cancelamento inicia novo trabalho;
- replay offline acessa rede;
- ferramenta nova é inserida antes das nove existentes;
- certificados anteriores são tratados como preservados sem recibo.

Mutation score exigido para selar o contrato: `1,0`.

## Gates de aceitação

- spec, baselines, schemas, regras, corpus, oráculos, mutações e splits congelados
  antes do candidato;
- `crystalline-lint .` sem violações de gravidade;
- `go test -count=1 ./...` passa duas vezes;
- `go test -race` passa nos pacotes novos e durable runs afetados;
- `go build -o bateia ./04_wiring/main.go` passa;
- `go vet ./...` passa;
- nenhum teste unitário exige rede ou modelo de IA;
- suites protegidas cobrem falso mecanismo, cross-source contamination,
  unknown, licença, popularidade, composição, injection e integridade;
- replay repetido produz artefatos byte a byte idênticos;
- scan de projeções não encontra segredo, conteúdo integral ou canário;
- métricas de mecanismo, hipótese, divergência e utilidade passam no split oculto;
- mutation score `1,0` sem transformar caso opaco em sucesso;
- receipts demonstram refinamento observavelmente compatível de 0025, 0026 e
  0027.

## Regime Tekt previsto

Aplicar protocolo completo. O usuário ocupa a autoria da intenção. Esta spec e
os três certificados formam a raiz causal. Autor do contrato recebe intenção,
spec e schemas públicos, mas não o patch. Autor dos oráculos recebe corpus sem
implementação. Implementador não lê split `TEST`, rótulos finais, mutações ou
suites protegidas. Adversário pode cruzar evidências, adulterar hashes, promover
metadados, esconder licenças e injetar propostas, mas não corrigir o candidato.
Verificador não escreve produção nem entradas seladas.

O laudo final precisa registrar hashes de spec, contratos, corpus, catálogo de
regras, baselines, mutações, candidato e evidências dos gates. Trocar taxonomia,
regra, baseline convencional, política de licença ou corpus invalida o selo a
partir do primeiro artefato afetado.

Sem isolamento material forte, declarar
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`, mesmo com todos os gates aprovados.

## Fora de escopo

- incorporar LLM, embedding ou agente autônomo ao binário do Bateia;
- provar arquitetura de aplicação fechada por descrição, review ou marketing;
- decidir automaticamente adoção, compra, fork, instalação ou substituição;
- gerar ou copiar implementação, abrir branch, commit, issue ou PR;
- executar plano de validação no repositório alvo;
- substituir `AutophagyReport`, `KatamariProposal` ou ranking de alternativas;
- transformar divergência em superioridade;
- emitir parecer jurídico ou autorizar cópia com licença desconhecida;
- resolver compatibilidade de múltiplos mecanismos sem evidência própria;
- recuperar `EvidenceBundle` de runs históricos relendo fontes mutáveis;
- expor código-fonte integral pelo MCP;
- adicionar novos conectores de mercado;
- transporte MCP remoto, autenticação multiusuário ou controle por tenant;
- retenção, deleção ou garbage collection de runs;
- aprender ou alterar regras automaticamente a partir do split oculto.
