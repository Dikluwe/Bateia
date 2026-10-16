# 0116 — Guarda de configuração estruturada

Lab v8, projection configuration_guards independente. Para FIELD_WRITE de
canal candidato, somente assignment direto no corpo de função checada, procurar
call anterior direto nesse corpo a método fornecido sem parâmetros. Destino
suporta padrão exato: corpo com único if, sem init/else, condição campo direto
do receiver != nil (ou nil != campo), corpo com único ExprStmt builtin panic.
Nil e panic por objetos Universe; receiver/campo por go/types, sem promoção de
campos embarcados. Nome não autoriza padrão. Guard não retorna normalmente se
a condição for verdadeira; panic pode ser recuperado fora do guard, sem retomar
as statements seguintes do corpo que sofreu panic.

Base simples do call e escrita devem ser o mesmo objeto. Goto/labels/branches,
rebinding/address-taking/capture dessa base no corpo tornam o recorte não
suportado. Guard condicional/defer/go, depois da escrita ou em closure não
permite projeção. Não compor viabilidade, atomicidade, mutação externa ou
estabilidade até a escrita. Estado PRIOR_GUARD_SAME_BASE_OBJECT, autoridade
LIMITED_GUARD_PATTERN; atomicity UNRESOLVED.

Registrar predicado, panic, definição/call/escrita e writes do campo verificado
em corpos fornecidos, todos físicos/hash. Writes da sentinela são por declaração
e podem ser de outra instância, inclusive em closure. Não provar inicialização
ocorrida. Limite 10.000 fatos; ordem/dedup determinísticos; preservar recibos.
Fonte conc: WithMaxGoroutines chama panicIfInitialized antes de escrever limiter;
guard verifica tasks != nil; init contém atribuição a tasks. Não classificar
proteção concorrente ou workers limitados.
