# Laudo 0117 — Modelo de checagem e intercalação

A bancada autoral executou quatro cenários com marcador de inicialização e
limite representados por variáveis atômicas separadas. Não implementa conc;
omite canais, sync.Once, workers, callbacks e aliases. Hook artificial entre
checagem e escrita permite rendezvous determinístico com testing/synctest.

Configurar antes de inicializar foi aceito; configurar depois foi rejeitado,
conservando o limite. Pausar depois do check, inicializar e retomar o store
permitiu escrita após inicialização no modelo. O controle oposto, inicializar
antes do check em outra goroutine, rejeitou e não chamou o hook. Snapshot lido
somente após sincronização. Detector race passou no modelo; isso não estabelece
atomicidade da sequência check+write.

Testemunhos normalizados foram iguais em duas execuções; logs brutos têm tempo
e não requerem byte-identidade. Observação estática de conc foi repetida sem
executá-lo e conserva atomicity=UNRESOLVED. Não se demonstrou corrida, defeito,
estabilidade ou limite global no candidato, nem que a intercalação pertença
a seu contrato de uso suportado. Referência autoral, não independente.

A distinção observável é entre contrato sequencial, data race e atomicidade.
O resultado fornece controle executável para avaliar futuras alegações do
classificador, mas não mede sua precisão semântica ou autentica hipótese sobre
o programa real. Observador v8 e produto não foram alterados.

crystalline-lint sem violações, go test ./..., race do modelo, go vet ./...,
build produto e diff check passaram. Fixtures usaram loopback e assinatura Git
desativada apenas no processo. Produto não importa Lab; módulos selados não
reproduzidos. Recibos anteriores preservados.

[Resultados](../../lab/bancada/concorrencia/0117/results/scenario-evaluation.json),
[verificação/hashes](../../lab/bancada/concorrencia/0117/results/verification.json)
e [contrato/reprodução](../../lab/bancada/concorrencia/0117/README.md).
Próxima avaliação: obrigações e alegações contra o contrato permitido de uso,
revisadas separadamente da implementação e destes modelos autorais.
