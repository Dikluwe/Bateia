# 0116 — Guarda anterior à escrita de configuração

Piloto v8 no pacote 0109, isolado do produto. configuration_guards é uma
projeção independente; fatos anteriores permanecem iguais. Fontes e ambiente
congelados; candidato não é executado.

## Contrato

Para escrita direta de campo de canal candidato no corpo de função checada,
localiza chamada direta anterior no mesmo corpo, sobre a mesma base simples.
O destino deve ser método fornecido, sem parâmetros/resultados, com padrão
exato: único if sem init/else, campo direto do receiver != nil (ou comparação
invertida), corpo com único panic builtin. Identidades de nil/panic/receiver/
campo vêm de go/types; nomes não autorizam a classificação.

Nesse padrão, guard não retorna normalmente quando a condição é verdadeira.
Panic pode ser recuperado fora desse método, sem retomar statements seguintes
do corpo que sofreu panic. O vínculo registrado é PRIOR_GUARD_SAME_BASE_OBJECT,
com autoridade LIMITED_GUARD_PATTERN e atomicity=UNRESOLVED. Registra checagem
anterior; não prova o valor do campo no momento da escrita ou proteção atômica.

Guard condicional/adiado/go/posterior, receiver diferente, jump, rebinding,
address-taking ou capture da base não autorizam esse recorte. Closure não é
composição do corpo externo. Guard com defer/recover ou método fora do padrão
fica fora da projeção. Um efeito intermediário pode mudar a sentinela sem
remover a existência da checagem anterior; permanece uma limitação explícita.

checked_field_writes inventaria assignments à declaração do campo verificado,
inclusive em closures de corpos fornecidos. Podem ser outra instância ou caminho
inalcançado. Não prova inicialização ocorrida, ligação de receivers ou ciclo de
uso. Posições/hash, ordenação e limite de 10.000 fatos. Ausência de projeção não
significa ausência de guarda implementada por outro padrão.

## Fonte congelada

Em conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b:

| Ponto | Fonte |
| --- | --- |
| call anterior à escrita | pool/pool.go:92 |
| escrita de limiter | pool/pool.go:96 |
| definição panicIfInitialized | pool/pool.go:111 |
| predicado | p.tasks != nil |
| assignment à sentinela tasks, em init | pool/pool.go:104 |

Isso localiza a relação configuração/checagem no código. A atribuição em init
ocorre dentro de callback de initOnce.Do; sua execução e identidade com a base
da configuração não foram comprovadas. A guarda tampouco estabelece exclusão
mútua ou ausência de corrida entre checagem e escrita. Limite global de workers
e estabilidade do canal continuam pendentes.

## Avaliação e reprodução

20 controles autorais passaram: renomeação/inversão da comparação, remoção,
condição, defer/go, ordem, base diferente, rebind/endereço/capture, nomes
sombreados, recover, predicado errado, closure/jump, efeito intermediário e
pacote inválido. Testes adicionais verificam sentinelas de outra função,
reordenação, posições/hash e excesso de fatos. Referência autoral; precisão
semântica independente não medida. 91 perfis anteriores passaram, sem reescrever
recibos históricos.

```sh
go build -o /tmp/bateia-observer-0116 ./lab/bancada/concorrencia/0109/cmd/observer
python3 lab/bancada/concorrencia/0116/verify.py /tmp/bateia-observer-0116
/tmp/bateia-observer-0116 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json --environment lab/bancada/concorrencia/0111/testdata/environment.json
```

Replay offline, JSON/hashes estritos e publicação exclusiva de arquivo novo.
results contém recibos e logs. Próxima avaliação precisa distinguir o contrato
de configuração sequencial das garantias de uso concorrente, com ligação de
receivers e referência revisada independentemente.
