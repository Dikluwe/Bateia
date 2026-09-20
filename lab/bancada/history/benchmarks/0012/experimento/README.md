# Teste E-0012 — Autofagia competitiva e Katamari real

**Estado**: `EXECUTADO PARCIALMENTE — FAIL_DISCOVERY`  
**Passo avaliado**: [0012](../../../passos/0012-executar-autofagia-competitiva-e-formar-katamari-minimo.md)  
**Pipeline de julgamento**: sem LLM  
**Data do congelamento da intenção**: 2026-08-26

**Execução observada**: [`../run/execution-report.md`](../run/execution-report.md)

## Pergunta

Quando recebe somente a ideia congelada do Bateia, o funil consegue:

1. reconhecer capacidades comprovadas no próprio Bateia;
2. descobrir concorrentes funcionais sem recebê-los como resposta;
3. rejeitar controles lexicalmente atraentes e funcionalmente irrelevantes;
4. identificar lacunas internas reais;
5. propor a menor composição externa admissível que cubra essas lacunas;
6. preservar `UNKNOWN` quando correspondência ou evidência não puder ser provada?

## Hipóteses pré-registradas

### H1 — Autofagia

O Bateia classificará como `KEEP` pelo menos quatro capacidades internas e não
classificará como `GAP` nenhuma capacidade julgada `UNKNOWN` pelo revisor.

### H2 — Descoberta competitiva

Sem receber a lista de sementes durante a busca, o Bateia encontrará pelo menos
70% das sementes que possuam repositório público capturável no instante do selo,
no top 30 agregado. Produtos sem repositório verificável ficam fora do denominador.

### H3 — Discriminação

Os três controles negativos terão decisão de screening `DROP` ou não entrarão no
Katamari. Nenhum poderá aparecer como componente selecionado.

### H4 — Katamari

Se os oráculos comprovarem pelo menos uma lacuna interna e uma cobertura externa
admissível, o Katamari cobrirá mais gaps que a autofagia sem absorver capacidade
`KEEP`, licença proibida ou revisão opaca.

### H5 — Minimalidade

Para a matriz congelada, nenhuma composição admissível com menos projetos terá a
mesma cobertura da proposta selecionada. Em empate de quantidade, nenhuma terá
custo menor segundo a política congelada.

## Artefatos congelados antes da execução

- [`bateia-idea.json`](bateia-idea.json): ideia fornecida à descoberta;
- [`competitive-registry.json`](competitive-registry.json): sementes usadas
  somente para calcular recall depois da descoberta;
- [`katamari-policy.json`](katamari-policy.json): limites e gates;
- [`reuse-policy.json`](reuse-policy.json): política de decisão já existente;
- [`oracle-form.json`](oracle-form.json): formulário vazio de julgamento humano.

O executor da descoberta não pode ler `competitive-registry.json`. O motor não
pode ler `oracle-form.json`; ele recebe apenas os bundles e regras congelados.

## Desenho

### Fase A — descoberta cega

1. Executar `idea plan`, `idea expand` e `idea discover` usando apenas
   `bateia-idea.json`.
2. Fixar estratégia, limite de 30 candidatos e versões das regras antes do run.
3. Salvar snapshot, tentativas, queries efetivas, commit de cada candidato e hash.
4. Só depois do snapshot, abrir `competitive-registry.json` e calcular recall.
5. Ausência de uma semente conta como perda; é proibido inseri-la manualmente.

### Fase B — controles negativos

Antes da peneira, capturar um repositório imutável para cada classe declarada no
registry:

- clone ou ferramenta do jogo Katamari;
- software ou dataset de autofagia biológica;
- projeto sobre mineração/peneiramento sem seleção de software.

A seleção em cada classe deve usar o primeiro resultado público capturável da
query declarada, registrar todos os resultados anteriores e não consultar o
score do Bateia. Se não houver resultado capturável, o controle fica
`UNJUDGEABLE`; não é substituído retrospectivamente por um exemplo conveniente.

### Fase C — bundles correspondentes

1. Capturar o Bateia local em identidade e hash imutáveis, incluindo mudanças da
   worktree necessárias ao Passo 0012.
2. Coletar bundles dos candidatos descobertos e dos controles em `repo@commit`.
3. Não usar `HEAD`, URLs assinadas, tokens ou conteúdo dependente de rede no
   replay.
4. Registrar falhas e budget; não convertê-los em ausência.

### Fase D — oráculos humanos cegos

Um revisor recebe a ideia, regras e bytes congelados, mas não recebe screening,
assessment, autofagia ou Katamari produzidos pelo motor. Para cada par
capacidade–repositório, preenche uma cópia do `oracle-form.json` com:

- `PROVEN`, `PARTIAL`, `ABSENT`, `CONFLICT` ou `INCONCLUSIVE`;
- paths e linhas que sustentam o julgamento;
- gate de licença, operação e manutenção;
- justificativa;
- decisão sobre utilidade do componente.

O formulário é hashado antes da execução candidata. Alteração posterior invalida
o teste.

### Fase E — replay offline

Executar primeiro autofagia e depois Katamari somente com os artefatos congelados:

```bash
go run ./04_wiring idea autofagia \
  --idea 00_nucleo/benchmark/0012/experimento/bateia-idea.json \
  --bundles 00_nucleo/benchmark/0012/run/bateia-internal.json \
  --rules 00_nucleo/benchmark/0012/run/evidence-rules.json \
  --target local/bateia@<WORKTREE_HASH> \
  --output 00_nucleo/benchmark/0012/run/autofagia.json

go run ./04_wiring idea katamari \
  --autofagia 00_nucleo/benchmark/0012/run/autofagia.json \
  --assessments 00_nucleo/benchmark/0012/run/competitors-assessment.json \
  --policy 00_nucleo/benchmark/0012/experimento/katamari-policy.json \
  --reuse-policy 00_nucleo/benchmark/0012/experimento/reuse-policy.json \
  --output 00_nucleo/benchmark/0012/run/katamari.json
```

`run/` só deve ser criado durante a execução. O replay é repetido com rede e
tokens indisponíveis e os hashes das saídas precisam coincidir.

## Métricas

| Métrica | Numerador | Denominador |
|---|---|---|
| precisão interna | findings internos iguais ao oráculo | findings internos julgáveis |
| recall interno | capacidades internas `PROVEN` recuperadas | capacidades internas `PROVEN` no oráculo |
| recall competitivo | sementes descobertas no top 30 | sementes correspondentes pré-registradas |
| rejeição irrelevante | controles rejeitados | controles negativos julgáveis |
| cobertura antes | capacidades `KEEP` | capacidades julgáveis da ideia |
| cobertura depois | `KEEP` + gaps cobertos | capacidades julgáveis da ideia |
| absorção redundante | componentes que só cobrem `KEEP` | componentes selecionados |
| minimalidade | custo da proposta menos custo ótimo | custo ótimo enumerável |

Denominador zero é `null`. `INCONCLUSIVE` e `UNJUDGEABLE` nunca são sucesso.

## Gates

- precisão interna ≥ 90%;
- recall interno ≥ 85%;
- recall competitivo ≥ 70%;
- rejeição irrelevante = 100%;
- cobertura depois > cobertura antes quando houver gap coberto admissível;
- absorção redundante = 0%;
- absorções proibidas = 0;
- arrependimento de minimalidade = 0;
- duas execuções offline produzem JSON canônico equivalente;
- toda divergência aponta capacidade, regra, fonte ou gate.

## Vereditos possíveis

- `PASS`: todos os gates aplicáveis passam;
- `FAIL_DISCOVERY`: concorrentes pré-registrados não foram recuperados;
- `FAIL_CORRESPONDENCE`: ideia, candidato e bundle não formam cadeia causal;
- `FAIL_EVIDENCE`: estados divergem dos oráculos;
- `FAIL_AUTOPHAGY`: estado interno ou gap incorreto;
- `FAIL_KATAMARI`: composição redundante, proibida ou não mínima;
- `INCONCLUSIVE`: identidade, captura ou oráculo insuficiente.

Falha não autoriza alterar expected outcomes. Ela origina um finding e um passo
posterior específico.

## Critério de parada

Interromper sem calcular utilidade se ocorrer qualquer uma destas condições:

- o executor da descoberta leu o registry antes do snapshot;
- algum candidato foi inserido manualmente;
- bundle usa `HEAD` ou não corresponde ao candidato descoberto;
- formulário de oráculo mudou depois de hashado;
- replay acessou rede, cache externo ou token;
- um controle negativo foi escolhido após observar o score;
- `UNKNOWN` foi contado como gap, cobertura ou acerto.
