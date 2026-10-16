# Laudo 0112 — Relações do ciclo de tokens

O observador v4 acrescentou um grafo de relações por identidade de tipos e
contenção sintática, com localizadores físicos/hash e autoridade
STATIC_TYPES_AND_SYNTAX. Permanece no Lab, puro e com dados do ambiente 0111.
Não há promoção de mecanismos produtivos ou inferência de limite global.

Na fonte conc congelada, 80 relações: 68 destinos estáticos, oito argumentos
função, um FUNCTION_CONTAINS_GO, um SEND_CASE_CONTAINS_GO_CALL, um
DEFERRED_RECEIVER_RECEIVE e um SAME_DECLARED_CHANNEL_FIELD.

O envio ao limiter em pool/pool.go:56 está no mesmo case da chamada em :59,
cujo alvo é WaitGroup.Go em waitgroup.go:28; essa definição contém go em :30.
O argumento literal possui chamada a worker em pool/pool.go:60, alvo em :150.
O defer em :153 liga-se ao receive do receptor em :172, sob if(l != nil).
Campo do envio e receptor do defer compartilham objeto de declaração; isso
não comprova instância/valor do canal. A chamada no default em :50 é distinta
e não recebe relação de case com envio ao limiter.

Estes vínculos localizam o ciclo candidato para investigação. Não demonstram
alcance, execução do callback, ausência de mutação, registro de defer em todos
os caminhos, balanceamento de tokens ou workers limitados. Lacunas explícitas
de caminho, callback e instância acompanham o grafo. Condições de origem e
alvo são separadas; chamadas adiadas conservam modo DEFERRED.

13 casos autorais em 12 famílias passaram, incluindo remoção de go,
receive/campo/case diferente, closure, defer condicional, chamada adiada ou
após return, variável de função, interface e erro de tipos. Perfil esperado
errado foi discriminado. Reordenação e localizadores conferiram; excesso de
10.000 relações foi rejeitado. Os 15 e 16 perfis anteriores também passaram.
Referências autorais: precisão semântica independente continua não medida.
Replays byte idênticos e nenhuma relação em conc sem seu ambiente fornecido.

crystalline-lint sem violações, go test ./..., race do observador, go vet ./...,
builds produto/piloto e diff check passaram. Fixtures da suíte usaram loopback
e assinatura Git desativada apenas no processo. Produto não importa Lab;
gates selados com módulo próprio não foram reproduzidos. Candidato não foi
executado. Evidência histórica preservada.

[Recibos](../../lab/bancada/concorrencia/0112/results/verification.json) e
[contrato/reprodução](../../lab/bancada/concorrencia/0112/README.md).
Próxima fatia: passagem de argumentos de função até invocação, com alcance e
identidade de instâncias mantidos como obrigações distintas.
