# 0112 — Relações tipadas do ciclo de tokens

O código ativo evoluiu para [0113](../0113/README.md). Os recibos deste
diretório conservam o ensaio v4; não executar seu verificador para publicar
resultados de uma versão posterior sobre estes recibos.

Piloto v4, código em lab/bancada/concorrencia/0109. Relações são campos novos
do relatório, separados dos perfis de observação anteriores. A transformação
permanece pura, com o ambiente congelado do 0111. Nenhum import produtivo ou
nova dependência. Resultados v1/v2/v3 preservados.

## Contrato

Somente pacotes sem erros de tipos fornecem funções/campos resolvidos.
STATIC_CALL_TARGET liga chamada a uma definição fornecida pelo objeto
*types.Func; interface e variável de função não autorizam destino concreto.
FUNCTION_CONTAINS_GO liga definição à sintaxe go no próprio corpo; closures
possuem escopo separado. CALL_ARGUMENT_FUNCTION registra argumento literal,
sem afirmar sua execução.

SEND_CASE_CONTAINS_GO_CALL liga um envio de canal na comunicação de select
à chamada no corpo do mesmo case, quando o destino contém go no seu próprio
corpo. Essa relação é contenção lexical. Um return anterior pode impedir o
alcance, e uma chamada defer conserva call_mode=DEFERRED. Não é dominância
nem prova de reserva antes de execução em todos os caminhos.

DEFERRED_RECEIVER_RECEIVE liga defer de método de receptor canal a receive do
objeto exato desse receptor, com condições do alvo separadas das condições do
chamador. SAME_DECLARED_CHANNEL_FIELD vincula envio e receptor do defer pelo
mesmo objeto de declaração de campo; não comprova mesma instância, valor ou
ciclo balanceado. Receivers diferentes e mudanças do campo permanecem lacunas.

Autoridade STATIC_TYPES_AND_SYNTAX, posições físicas/hash em ambas pontas,
ordem determinística e limite de 10.000 relações com erro explícito se excedido.
Condições continuam sintáticas. Relações exigem tipos válidos; nenhuma relação
na fonte conc quando se remove seu ambiente. Lacunas de alcance, execução de
callback e identidade de instâncias acompanham as relações correspondentes.
Limites de workers/admissão/cancelamento não são promovidos.

## Fonte real

Em conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b, o relatório contém 80
relações: 68 destinos estáticos, oito argumentos função e uma relação de cada
uma das quatro categorias restantes. Os pontos relevantes são:

| Relação | Origem | Alvo | Condição/limite |
| --- | --- | --- | --- |
| envio no case e chamada ao lançador | pool/pool.go:56 | pool/pool.go:59 | else(p.limiter == nil), select.case:0 |
| alvo tipado do lançador | pool/pool.go:59 | waitgroup.go:28 | definição WaitGroup.Go fornecida |
| presença de go | waitgroup.go:28 | waitgroup.go:30 | execução do callback ainda não resolvida |
| chamada a worker no argumento literal | pool/pool.go:60 | pool/pool.go:150 | definição fornecida; vínculo com execução não provado |
| defer e receive no receptor | pool/pool.go:153 | pool/pool.go:172 | alvo sob if(l != nil) |
| mesma declaração de campo | pool/pool.go:56 | pool/pool.go:153 | mesma declaração, identidade de instância pendente |

Há também chamada ao mesmo lançador no ramo default, pool/pool.go:50, quando
p.limiter == nil. Ela não recebe SEND_CASE_CONTAINS_GO_CALL. Isso permite
inspecionar separadamente a via sem limite e a via com token sem classificar
o pool inteiro como bounded.

A existência desses pontos não demonstra uma trajetória executável que os
atravessa. A execução de callback passa por Catcher.Try e ainda não foi ligada
por análise de argumentos; aliases, mutação do limiter, uso concorrente da
configuração, panic antes do registro do defer e saída do worker não foram
verificados globalmente.

## Avaliação e reprodução

13 casos autorais em 12 famílias, com perfis completos de relações:
renomeação, ausência de go, outro receive/campo/case, closure, defer condicional,
launch adiado/inalcançável, indireção, interface e erro de tipos. Todos passaram;
perfil esperado errado foi discriminado. Precisão semântica independente não
medida. Os 15 casos históricos e 16 contrastes de tipos também passaram.
Testes adicionais verificam reordenação, localizadores e excesso de relações.

Na raiz:

```sh
go build -o /tmp/bateia-observer-0112 ./lab/bancada/concorrencia/0109/cmd/observer
python3 lab/bancada/concorrencia/0112/verify.py /tmp/bateia-observer-0112
/tmp/bateia-observer-0112 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json --environment lab/bancada/concorrencia/0111/testdata/environment.json
```

Replay usa dados congelados, sem rede/execução do candidato. CLI preserva JSON
estrito, hashes e publicação exclusiva de novo arquivo. verify.py checa perfis
e recibos em results. A checagem de fluxo exige ainda alcance/ordem, passagem
de callbacks e identidade dos canais; o grafo atual serve para localizar essas
obrigações e testar classificações futuras.
