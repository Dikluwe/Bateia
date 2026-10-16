# 0114 — Alcance e ordem estruturados limitados

O código ativo evoluiu para [0115](../0115/README.md); recibos deste diretório
conservam a v6. Publicar novos resultados na bancada atual.

Piloto v6 no pacote 0109, isolado do produto. call_site_reachability e
defer_order são projeções novas; relações e callbacks anteriores não mudam.
Análise pura sobre fontes/ambiente congelados; candidato não é executado.

## Contrato

Cada site de chamada tem localização/hash, escopo, condições sintáticas, modo
DIRECT/DEFERRED/GO e autoridade LIMITED_STRUCTURED_CONTROL. A classificação
refere-se ao site/registro, não à execução posterior da goroutine ou defer.

UNREACHABLE_AFTER_TERMINATOR vincula um blocker físico: return anterior,
chamada direta ao builtin panic ou if cujos dois ramos encerram a sequência.
A expressão de return é analisada antes do terminador. Nome panic sombreado,
go panic e defer panic não encerram a sequência do chamador.

POSSIBLY_REACHABLE é uma sobreaproximação: não garante viabilidade ou execução. Sites lexicais também podem ocorrer em
expressões de tipo constantes, sem invocação runtime.
Não avalia guardas nem prova saída de loops/select. Por exemplo, um site após
for{} conserva possibilidade abstrata neste recorte. Sequências de corpos e
cases permitem distinguir retorno anterior no mesmo bloco/case. Goto, labels
e branch statements (break/continue/fallthrough) tornam todo o escopo
UNSUPPORTED_CONTROL_FLOW. Uma closure possui entrada independente da criação:
seu corpo não recebe inacessibilidade de um return no escopo externo.

PREFIX_DEFER_BEFORE_CALL_SITE liga defer diretamente no corpo do escopo a
site DIRECT em statement posterior, inclusive dentro de if/loop/case. Se o
site for atingido normalmente pelo corpo, o registro desse defer prefixo foi
atravessado. Defer condicional, posterior ou em outro escopo não recebe essa
relação. Não prova execução do callee diferido, retorno do callback, ausência
de bloqueio/panic anterior ou balanceamento de tokens. Escopos com jumps não
emitem ordem. Limite de 10.000 fatos novos, erro antes de publicação.

As projeções lexical/tipada continuam presentes mesmo quando um site é
classificado como inacessível. Isso permite inspecionar o limite entre presença
no código e alcance, preservando lacunas globais de execução/instância.

## Fonte real

No conc congelado, 144 sites classificados POSSIBLY_REACHABLE e dez relações
de ordem. Nenhum desses estados demonstra uma execução. No worker, o registro
do defer de pool/pool.go:153 precede os sites :156 (initialFunc(), sob if) e
:160 (f(), no corpo do range). A ordem se restringe ao escopo de worker.

Na closure de WaitGroup.Go, defer em waitgroup.go:31 precede o site Try em :32.
Em Catcher.Try, defer em panics/panics.go:22 precede a invocação f() em :23.
Condições e escopos permanecem explícitos. O tipo do canal e os callbacks são
rastreáveis, mas alcance global, estabilidade da instância e quantidade de
workers ainda não foram demonstrados.

## Avaliação/reprodução

19 controles autorais passaram: retornos, if, panic/builtin sombreado,
defer/go panic, expressão de return, case, closure, jumps, loop e ordem de defer.
Expectativa errada de alcance foi discriminada. Os 59 perfis anteriores
(sintaxe/tipos/relações/callbacks) passaram; precisão semântica independente
continua não medida. Testes verificam localizadores, reordenação e excesso de
fatos. Evidência anterior não foi reescrita.

Na raiz:

```sh
go build -o /tmp/bateia-observer-0114 ./lab/bancada/concorrencia/0109/cmd/observer
python3 lab/bancada/concorrencia/0114/verify.py /tmp/bateia-observer-0114
/tmp/bateia-observer-0114 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json --environment lab/bancada/concorrencia/0111/testdata/environment.json
```

Replay usa dados congelados, sem rede/execução do candidato. CLI conserva JSON
estrito e saída exclusiva para arquivo novo. results contém recibos e logs.
Próxima obrigação: identidade/estabilidade do canal entre reserva e liberação,
sem converter compartilhamento de declaração em igualdade de instância.
