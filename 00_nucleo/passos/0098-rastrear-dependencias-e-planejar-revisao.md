# 0098 — Rastrear dependências e planejar revisão incremental

## Resultado

O modo avançado ganha um índice consultivo entre entradas, fontes declaradas,
alegações registradas e decisão histórica. Uma mudança explícita pode produzir
um plano de revisão dos nós alcançáveis, sem recalcular, aprovar ou trocar o
veredito. A pesquisa simples e os comandos sem opt-in permanecem intactos.

## Uso

Adicionar `--trace` a `bateia idea recommend` gera `trace-graph.json` no dossiê
e um estágio `trace` em `run.json`. Não muda RunID, RuleVersions, relatórios
causais ou os critérios de hash dos checkpoints de evidência. Falhar na emissão
solicitada ainda falha o comando e registra a execução como FAILED.

Adicionar `--trace-output grafo.json` a `bateia idea decide` gera um sidecar,
preservando o relatório normal. São duas saídas, não uma transação de arquivos:
uma falha posterior na escrita do relatório pode deixar o sidecar emitido.

Também há duas bocas totalmente offline:

```sh
./bateia idea trace build --input entrada.json --output grafo.json
./bateia idea trace impact --input mudancas.json --output revisao.json
```

O arquivo de build tem `schema_version: "1"`, `report` com o relatório de reúso
existente e, opcionalmente, `idea` e `rule_versions`. Não aceita um relatório nu.
O diretório de saída destes comandos deve existir; recommend cria seu dossiê.

O arquivo de impacto tem este formato; substituir `graph` pelo grafo completo
emitido e `node_id` por um ID realmente presente em `graph.nodes`:

```json
{
  "schema_version": "1",
  "graph": {},
  "changes": [
    { "node_id": "ID_REAL_DA_FONTE", "reason": "EXPIRED" }
  ]
}
```

O exemplo é um molde, não uma entrada válida antes dessas substituições.
Motivos aceitos: CONTENT_CHANGED, EXPIRED, REVOKED, POLICY_CHANGED,
RULE_CHANGED, INPUT_CHANGED e JUDGMENT_CHANGED. `new_fingerprint`, se informado,
é SHA-256 hexadecimal minúsculo do registro declarado do nó, não necessariamente
o hash do arquivo-fonte. CONTENT_CHANGED com fingerprint idêntico é um no-op;
EXPIRED e REVOKED pedem revisão mesmo com bytes iguais. Não há relógio ou monitor
automático: quem chama deve informar vencimento, revogação ou alteração.

## Modelo e semântica

Nós INPUT incluem report/policy e idea/versões quando fornecidos. CLAIMs indexam
matriz de capacidades, razões, bloqueadores, componentes, composições e revisão
técnica. SOURCE indexa EvidenceSource desses registros. DECISION preserva o
veredito. A direção da aresta é dependência para consumidor. O build usa INFORMS
para fontes, sem inventar prova; alegações e entradas alimentam a decisão.
Grafo importado pode declarar SUPPORTS/REFUTES; todos os tipos propagam revisão.

O relatório inteiro tem fingerprint conservador. Campos que não são indexados
individualmente (como ações de investigação) estão cobertos por essa dependência
global, não por uma garantia de rastreabilidade fina. CLAIM armazena label e
fingerprint, não o conteúdo completo: interpretar a alegação exige guardar o
relatório original. Não é ainda um assurance case autossuficiente.

IDs de fonte usam chave estruturada, sem concatenação ambígua, e não incorporam
ContentHash/snippet/strength. Assim uma observação alterada mantém identidade.
Declarações conflitantes para a mesma identidade são recusadas. Ausência de
ContentHash não impede registrar uma fonte; sua validade continua desconhecida.

O selo `graph_hash` usa JSON canônico do grafo com nós/arestas ordenados e o
campo de hash vazio. Fingerprints do report preservam sua representação original:
reordenar vetores do report pode mudar o selo, mesmo sem mudança semântica.
O hash dos bytes do arquivo registrado no manifest é distinto do selo interno.

Impacto valida o selo, percorre somente dependências declaradas e emite IDs
alterados, impactados e não alcançados, tarefas e ancestrais causadores.
REVIEW_REQUIRED pede revisão; NO_RECORDED_IMPACT significa somente que nenhuma
mudança ativa foi registrada. `decision_before == decision_after` e
`decision_impact: "NONE"` significam ausência de mutação, não autorização para
continuar confiando na decisão após uma revogação. O grafo histórico não é editado.

## Segurança e limites

O núcleo não faz chamadas de rede, não lê repositórios, não executa candidatos
e não chama o decisor.
Valida schema, enums, IDs únicos, fingerprints, fontes ambíguas, exatamente uma
decisão, DAG sem ciclos/autorreferências/arestas pendentes ou duplicadas.
Limites máximos: 1000 nós, 4000 arestas, 100 mudanças, 8192 bytes por string,
2 MiB de JSON/texto, profundidade 64 e orçamento estrutural conservador de 20000
valores. Este último pode recusar entradas antes dos outros tetos. O plano de
impacto também passa pela guarda de saída antes da serialização integral:
expansão de muitos ancestrais é recusada com erro, não com JSON parcial.

CLI trace lê JSON estrito limitado, recusa chaves duplicadas/desconhecidas e
entradas especiais. O novo destino não pode sobrescrever entradas por caminho,
hardlink ou alias ancestral; saídas simbólicas são recusadas. O writer usa pai
confinado, temporário exclusivo, sync e rename, verifica identidade do diretório
e rejeita segredos configurados antes de emitir bytes. Não é uma sandbox contra
outro processo do mesmo usuário que altera continuamente o namespace.
No decide opt-in, a saída principal também não pode identificar uma entrada.

Sem --trace na retomada, um grafo antigo é preservado, mas não é consumido nem
incluído no novo manifest; um aviso explicita isso. Verificar estágio/hashes do
manifest atual, não apenas a existência de trace-graph.json.

Hash não autentica fatos, autoria, fonte ou completude. Reconstruir o selo de
um grafo inventado não o transforma em prova. O plano depende de mudanças
declaradas e não descobre dependências ocultas.

## Verificação

Testes do núcleo cobrem selagem, propagação seletiva/refutação, conflitos,
entradas inválidas/limites, no-op e preservação do veredito. Testes CLI cobrem
JSON estrito, aliases, escrita segura, comandos offline e ausência de bytes em
erros de validação. Integração compara execução com/sem trace, igualdade de
artefatos causais/identidade, hash do estágio, impacto de expiração, preservação
histórica e retomada sem opt-in. Implementação, integração e revisão foram
divididas entre agentes com responsabilidades delimitadas.

Próximo incremento possível: localizadores explícitos de cada claim no relatório
e execução controlada das tarefas de revisão, sempre gerando uma nova decisão
separada da decisão histórica.


## Adendo 0098 — Execução do legado após extração do transporte

O recorte antigo de quatro arquivos MCP deixou de compilar depois da extração
anterior dos helpers de dispatch/segurança. Executar o pacote completo
./02_shell/mcp, selecionando explicitamente os mesmos goldens R3 e certificados,
mais o teste real TestStep0025LegacySuccessfulSearchUsesEvidenceContract
(sucessor já existente do nome MatchesFrozenBaselineBytes). Não alterar hashes,
goldens, expectativas, mutantes ou classificação dos casos. Essa correção é de
invocação: recompõe helpers ausentes e impede seletor legado vazio; não credita
qualidade do produto ou nova prova aos testes de infraestrutura.


## Validação final observada

`go test -p 1 ./...`, `go vet ./...`, `git diff --check` e testes -race
direcionados do núcleo, CLI e wiring passaram. A execução paralela da suíte
completa também expôs um timeout intermitente no gate legado R5: Start não
retomou dentro do prazo de um segundo do fixture de recuperação. Esse teste
passou em rodadas isoladas e na suíte sequencial. Não alterar seu prazo ou
expectativa para ocultar a intermitência; tratar robustez do harness em trabalho
separado. Isso não converte testes de infraestrutura em prova de decisões.
