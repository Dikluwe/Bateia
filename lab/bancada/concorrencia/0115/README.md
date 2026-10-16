# 0115 — Identidade e estabilidade do canal: auditoria limitada

O código ativo evoluiu para [0116](../0116/README.md); recibos deste diretório
conservam a v7. Usar a bancada atual para publicar resultados posteriores.

Piloto v7 no pacote 0109, isolado do produto. channel_identity acrescenta uma
auditoria por declaração a pares SAME_DECLARED_CHANNEL_FIELD; projeções anteriores
permanecem iguais. Transformação pura, sem execução do candidato.

## Contrato

base_link distingue SAME_BASE_OBJECT, DISTINCT_BASE_OBJECTS e
UNRESOLVED_BASE_EXPRESSION por identidade go/types de bases simples. Mesmo
objeto de base não prova mesmo valor do campo; objetos distintos não provam
instâncias diferentes. State permanece INSTANCE_IDENTITY_UNRESOLVED em todos
os casos, com autoridade DECLARATION_SCOPED_AUDIT.

O inventário nos corpos de funções fornecidos/checados registra FIELD_WRITE,
AGGREGATE_OR_BASE_WRITE, FIELD_ADDRESS_EXPOSED, RECEIVER_ARGUMENT_EXPOSURE e
RECEIVER_CAPTURE. Escritas via selector sobre alias simples também apontam
para a declaração de campo. Atribuição ao agregado/rebinding de base é inventário
amplo, não demonstra substituição daquele canal. Captura rastreia os objetos de
base das pontas do par; não é uma análise completa de escapes.

A cobertura é SUPPLIED_FUNCTION_BODIES_ALL_INSTANCES_OF_FIELD_DECLARATION.
Eventos podem pertencer a outra instância ou trecho inalcançado. São pontos
para investigação, não defeitos. Não resolve aliases transitivos, mutação
externa, unsafe/reflection, ordem de eventos ou estabilidade concorrente.
Ausência de evento não certifica estabilidade. Inicializadores fora dos corpos
de funções e campos promovidos/estruturas indiretas não têm cobertura completa.
Até 10.000 fatos auditados; excesso falha antes da publicação. Posições/hash,
ordem e deduplicação explícitos.

## Fonte congelada

Um par em conc: envio pool/pool.go:56 e defer em :153, bases de objetos distintos
(Pool.Go e worker). Identidade continua inconclusiva. Eventos:

| Evento | Fonte | Significado limitado |
| --- | --- | --- |
| FIELD_WRITE | pool/pool.go:96 | atribuição em WithMaxGoroutines |
| RECEIVER_CAPTURE | pool/pool.go:50 | callback do ramo sem limiter |
| RECEIVER_CAPTURE | pool/pool.go:59 | callback do case de envio ao limiter |

A escrita de configuração não foi classificada como mutação durante execução.
A guarda panicIfInitialized e sua relação com o ciclo de uso ainda precisam
ser estudadas. Os vínculos de callback mostram passagem estrutural; a auditoria
não compõe uma prova de identidade do receiver/canal ou de workers limitados.

## Verificação e uso

13 casos autorais passaram: base igual/distinta/complexa, write direto/via
alias/parentesizado, agregado/rebinding, endereço, argumento, capture, campo
diferente e erro de tipos. Testes adicionais verificam reordenação, hashes e
limite. Os 78 perfis anteriores passaram; as projeções da fonte real permaneceram
iguais. Precisão semântica independente continua não medida.

Na raiz:

```sh
go build -o /tmp/bateia-observer-0115 ./lab/bancada/concorrencia/0109/cmd/observer
python3 lab/bancada/concorrencia/0115/verify.py /tmp/bateia-observer-0115
/tmp/bateia-observer-0115 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json --environment lab/bancada/concorrencia/0111/testdata/environment.json
```

Replay offline com dados congelados; nenhum candidato é executado. CLI mantém
hashes/JSON estrito e publicação exclusiva de novo arquivo. results conserva
recibos e logs; evidência dos passos anteriores não foi reescrita.
Próximo recorte: contexto da escrita de configuração e ligação de receivers,
antes de julgar estabilidade da instância durante reserva/liberação.
