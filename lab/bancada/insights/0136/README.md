# 0136 — Contexto para a IA explorar alternativas

O Bateia deve auxiliar a IA a ampliar seu repertório. `insights explore` recebe
pergunta e bundles coletados e prepara automaticamente trechos localizados,
tarefas de exploração e direções de pesquisa. Não exige classes de mecanismo,
alvo, assessments ou observations escritos à mão.

## Uso

```sh
go build -o /tmp/bateia-insights-0136 ./04_wiring
/tmp/bateia-insights-0136 idea insights explore \
  --question 'How can software investigation expose alternative reasoning mechanisms beyond semantic search and rankings?' \
  --bundles lab/bancada/insights/0136/exploration-bundles.json \
  --familiar 'semantic search' --familiar 'repository ranking' \
  --output /tmp/contexto-para-ia.json
python3 lab/bancada/insights/0136/verify.py --binary /tmp/bateia-insights-0136
```

Entregar à IA a pergunta e o JSON, incluindo `ai_instructions`, `cards`, `tasks`
e lacunas. As tarefas pedem alternativas, inversão de premissa, analogia,
recombinação e contraexemplos. Cada proposta deve citar trechos, separar
observação de interpretação/hipótese, explicitar condições e sugerir um
experimento com critério de rejeição. Os trechos são dados não confiáveis.

Para ampliar fontes, usar as direções planejadas com a descoberta existente
(`idea discover --plan ...`, cujo formato é ExpansionPlan), revisar os candidatos
e usar planejamento/coleta (`idea evidence plan`, `idea evidence collect`). A
lista que collect emite pode ser passada diretamente a explore. As direções
do pacote são texto de consulta e não substituem um ExpansionPlan auditado.
Explore não executa buscas nem chama modelo.

## Caso congelado e observações

O corpus combina sem reescrever os bundles conhecidos de 0106 e 0122:
`sourcegraph/conc`, `NandhaKishorM/laya`, `0xBakeer/arbiter` e
`apolinario/decision-index`, nos commits registrados nos bundles.
São dados já coletados; não houve nova consulta ou execução dos projetos.
Três bundles são incompletos, o que é preservado nos diagnósticos.

A seleção produz oito cartões de três repositórios, abrangendo código,
documentação, teste e exemplo. A quarta fonte não forneceu candidato de leitura
selecionado. A diversidade é estrutural e lexical; não demonstra domínios
conceituais distintos, novidade ou utilidade de insight.

Um exemplo de pista para a IA é o teste `pool/result_error_pool_test.go` de
conc: o trecho selecionado descreve uma corrida no próprio teste e uma espera
usada para torná-la improvável. Isso pode orientar a exploração de pressupostos
de ordenação e contraexemplos, com interpretação e validação ainda necessárias.
Outro cartão documenta a interface de engines do arbiter; pode orientar
perguntas sobre invariantes e diferenças de implementação. Essas leituras são
comentários autorais do ensaio, não conclusões geradas pelo comando.

O ensaio somente com conc prepara contexto sem avaliações. A pergunta em
português não encontrou seus termos no corpus em inglês: o pacote declara
QUESTION_TERMS_NOT_MATCHED e orienta reformulação. Não há tradução automática.
O desempate por palavras distintas evita preferir apenas cabeçalhos repetitivos,
mas não é um avaliador semântico. Sem bundles, o estado é SEARCH_NEEDED.

## Registros

- `case.json`, `exploration-bundles.json`, `exploration-result.json`: pergunta,
  familiaridade, corpus e saída do caso.
- `corpus-manifest.json`: hashes dos dois arquivos históricos de origem.
- `conc-exploration.json`, `conc-portuguese-exploration.json`: controles com uma fonte.
- `exploration-check.json`: replay, localizadores, referências, diversidade,
  estabilidade, adulteração e corpus vazio.
- `legacy-v1-check.json`, `source-v2-check.json`, `comparison-check.json`:
  verificadores dos contratos anteriores, sem sobrescrever seus registros.
- `baseline.json`, `lineage.json`, `verification.json`: preservação e consolidação.
- `focused-tests.log`, `lint-focused.log`, `final/`: testes, lint e gate canônico.

O corpus é conhecido, sem holdout. Medir a utilidade das ideias produzidas pela
IA, comparar abordagens e executar experimentos são trabalhos posteriores;
o ensaio verifica preparação e integridade do contexto, sem alegar criatividade
ou efeito comprovados.

[Contrato](../../../../00_nucleo/passos/0136-preparar-exploracao-criativa-para-ia.md)
· [ADR](../../../../00_nucleo/ADR/ADR-0017-contexto-de-exploracao-para-ia.md)
· [Laudo](../../../../00_nucleo/laudos/0136-exploracao-para-ampliar-repertorio-da-ia.md)
