# Laudo 0113 — Callbacks até a invocação

O observador v5 acrescentou callback_relations por identidade de parâmetros,
separadas das relações anteriores. No conc congelado, 18 vínculos: sete
bindings, quatro forwardings e sete invocações de parâmetro. As 80 relações
v4 permaneceram iguais. Dados do ambiente 0111 e fontes não foram executados.

Literals de Pool.Go em pool/pool.go:50 e :59 ligam-se ao parâmetro f de
WaitGroup.Go em waitgroup.go:28. O uso capturado na goroutine em :32 aponta
para o parâmetro de Catcher.Try em panics/panics.go:21, invocado em :23.
Localizadores das declarações coincidem nos pontos de junção. Escopo do
forwarding é a closure de waitgroup.go:30; condições dos ramos de Pool.Go
permanecem distintas. Essa cadeia é estrutural, não execução comprovada.

Aliases locais, variádicos, indireção e transformações estão fora do recorte.
Parâmetro sobrescrito ou com endereço tomado perde crédito conservadoramente,
inclusive em closures. Sombra por outro objeto não herda identidade. Relações
mantêm autoridade STATIC_TYPES_AND_SYNTAX e lacunas de execução/alcance;
identidade do canal e balanceamento/limite de workers continuam pendentes.

15 casos autorais em 14 famílias passaram: remoção de forwarding/invocação,
nomeação, alias, sobrescrita direta/parentesizada/capturada, endereço, sombra,
indireção, condições, variádico e erro de tipos. Perfil esperado errado foi
discriminado. 13 perfis anteriores de relações e 31 de sintaxe/tipos passaram;
replays byte idênticos e ausência de vínculos em conc sem seu ambiente.
Testes adicionais verificaram escopos, localizadores, reordenação e limite
combinado de 10.000 relações. Precisão semântica independente não medida.

crystalline-lint sem violações, go test ./..., race do observador, go vet ./...,
builds produto/piloto e diff check passaram. Fixtures usaram loopback e
assinatura Git desativada apenas no processo. Produto não importa Lab. Gates
selados com módulo próprio não foram reproduzidos. Recibos históricos intactos.

[Recibos](../../lab/bancada/concorrencia/0113/results/verification.json) e
[contrato/reprodução](../../lab/bancada/concorrencia/0113/README.md).
Próximo passo: alcance/ordem dos pontos relevantes, seguido de estabilidade da
instância do canal; manter referências independentes como avaliação pendente.
