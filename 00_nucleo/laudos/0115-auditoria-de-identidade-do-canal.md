# Laudo 0115 — Auditoria de identidade do canal

O observador v7 acrescenta channel_identity aos pares de campo já encontrados,
sem modificar relações/callbacks/controle históricos. Auditoria pura no Lab,
com bases por objetos go/types e inventário localizado de writes/exposições/
captures. Estado de identidade permanece inconclusivo em todos os casos.

Em conc congelado, o par do envio pool/pool.go:56 e defer :153 tem bases
DISTINCT_BASE_OBJECTS: receivers de Pool.Go e worker. Isso não demonstra
instâncias diferentes. Eventos localizados: FIELD_WRITE em :96, configuração
WithMaxGoroutines, e RECEIVER_CAPTURE em :50 e :59. Não se demonstrou mutação
da mesma instância durante execução, nem estabilidade do canal. A guarda
panicIfInitialized e ligação de receivers continuam obrigações separadas.

Inventário cobre corpos fornecidos/checados e várias instâncias da declaração;
eventos podem ser inalcançados ou pertencer a outro objeto. Não são defeitos.
Aliases transitivos, mutação externa/concorrente, unsafe/reflection e cobertura
de estruturas indiretas não resolvidos. Ausência de evento não prova estabilidade.
Autoridade DECLARATION_SCOPED_AUDIT, posições/hash, deduplicação e limite explícito.

13 controles autorais passaram, com base igual/distinta/complexa, substituição
direta/via alias/parentesizada, agregado/rebinding, endereço, argumento, captura,
campo distinto e erro de tipos. Testes adicionais verificaram reordenação e
excesso de eventos. 78 perfis anteriores passaram; projeções reais antigas
conservaram igualdade. Replays byte idênticos e nenhum audit em conc sem
ambiente. Precisão semântica independente não medida; candidato não executado.

crystalline-lint sem violações, go test ./..., race do observador, go vet ./...,
builds produto/piloto e diff check passaram. Fixtures usaram loopback e
assinatura Git desativada apenas no processo. Produto não importa Lab; gates
selados de módulos separados não reproduzidos. Evidência histórica preservada.

[Recibo com hashes](../../lab/bancada/concorrencia/0115/results/verification.json)
e [contrato/reprodução](../../lab/bancada/concorrencia/0115/README.md).
Próximo recorte: contexto da escrita de configuração e ligação dos receivers,
antes de julgar identidade/estabilidade e limite global de concorrência.
