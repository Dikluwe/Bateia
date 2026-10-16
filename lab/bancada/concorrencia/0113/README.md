# 0113 — Parâmetros de callback até a invocação

O código ativo evoluiu para [0114](../0114/README.md); os recibos aqui conservam
a v5. Usar a bancada atual para publicar resultados de versões posteriores.

Piloto v5, no pacote 0109, ainda isolado do produto. O relatório acrescenta
callback_relations; as 80 relações do ensaio real v4 permanecem iguais.
Código puro, com fontes/export data congelados do 0111, sem rede ou execução
de candidato.

## Contrato e limites

CALLBACK_ARGUMENT_BINDING liga literal função ao parâmetro de um destino
estático fornecido. CALLBACK_PARAMETER_FORWARDING liga uso de parâmetro função
ao parâmetro do próximo destino; parameter_declaration identifica a declaração
do parâmetro de origem. CALLBACK_PARAMETER_INVOCATION liga declaração à chamada
por esse objeto exato. Os vínculos podem ser unidos por localizadores iguais
sem depender do texto dos nomes.

execution_scope localiza a função ou closure do uso; conditions e call_mode
pertencem a esse escopo. Uma closure não herda condições de sua criação.
Autoridade STATIC_TYPES_AND_SYNTAX e hash/posição em ambas extremidades. Tipos
com erro não autorizam vínculos. Limite combinado de 10.000 relações com falha
explícita, sem relatório parcial publicado.

Nesta fatia: argumentos literais e referências diretas a parâmetros, aridade
exata, funções não variádicas. Não rastreia aliases locais, chamadas por
interface/variável, argumentos transformados/armazenados ou chamadas genéricas.
Parâmetro sobrescrito ou com endereço tomado em qualquer trecho do corpo,
inclusive closures, perde crédito: análise conservadora, sem resolver ordem
ou alcance dessas operações. Sombra por outro objeto não é parâmetro original.
Funções variádicas ficam fora do recorte.

As relações não afirmam execução, caminho viável, quantidade de invocações ou
registro/liberação de tokens. CALLBACK_EXECUTION_UNRESOLVED permanece;
CALLBACK_VALUE_FLOW_UNRESOLVED registra mutação/address-taking. Condições de
chamadores e de invocações não são compostas como garantia global.

## Fonte congelada

Em conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b, foram localizadas 18
relações novas: sete bindings, quatro forwardings e sete invocações de
parâmetros. A cadeia relevante para o pool é:

| Relação | Origem | Alvo | Escopo/condição |
| --- | --- | --- | --- |
| literal no ramo default -> parâmetro f | pool/pool.go:50 | waitgroup.go:28 | if(p.limiter == nil), select.case:1 |
| literal no case do limiter -> mesmo parâmetro f | pool/pool.go:59 | waitgroup.go:28 | else(p.limiter == nil), select.case:0 |
| uso do parâmetro f -> parâmetro de Try | waitgroup.go:32 | panics/panics.go:21 | closure da goroutine em waitgroup.go:30 |
| parâmetro de Try -> chamada f() | panics/panics.go:21 | panics/panics.go:23 | corpo de Catcher.Try |

O forwarding referencia exatamente a declaração que recebe ambos os literals;
seu destino coincide com a declaração da invocação. Isso identifica uma cadeia
estrutural consistente. Ainda não prova que uma execução a percorreu ou que o
callback chamará worker em todo caminho. As obrigações de alcance, identidade
do canal e balanceamento dos tokens continuam separadas.

## Verificação e uso

15 casos autorais em 14 famílias: nomes, remoção de forwarding/invocação,
alias, sobrescrita direta/parentesizada/capturada, endereço, sombra, indireção,
condições, variádico e erro de tipos. Todos os perfis passaram; expectativa
errada foi discriminada. Precisão semântica independente não medida.
Os 13 perfis de relações anteriores e 31 perfis de sintaxe/tipos passaram.
Testes adicionais verificam reordenação, localizadores, escopos e limite total.

Na raiz:

```sh
go build -o /tmp/bateia-observer-0113 ./lab/bancada/concorrencia/0109/cmd/observer
python3 lab/bancada/concorrencia/0113/verify.py /tmp/bateia-observer-0113
/tmp/bateia-observer-0113 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json --environment lab/bancada/concorrencia/0111/testdata/environment.json
```

Replay usa arquivos congelados e não substitui recibos dos passos anteriores.
results contém relatórios e hashes. CLI mantém JSON estrito e publicação de
novo arquivo. Nenhum vínculo de callback em conc sem o ambiente fornecido.

Próxima fatia: alcance/ordem dos pontos relevantes, incluindo retorno antes de
invocação e registro condicional de defer; depois identidade/estabilidade do
canal. Um vínculo estrutural sozinho permanece insuficiente para bounded
concurrency ou qualidade semântica do classificador.
