# Laudo 0110 — Resolução limitada de canais

A versão go-concurrency-syntax-0110-v2 acrescenta observação por go/types a
canais nomeados e aliases em pacotes fornecidos sem erros. Transformação
permanece pura no Lab e não carrega imports. Identidade de tipo tem autoridade
GO_TYPES_SUPPLIED_PACKAGE; capacidade segue SYNTAX_ONLY, sem promoção a
limite global de concorrência.

Os 15 perfis históricos passaram na cópia baseline v2. Os 16 contrastes novos
em 15 famílias também passaram: nomes, aliases, tipos locais, capacidades,
slice/map, sombras, ciclo e erros. Testes adicionais verificam declaração entre
arquivos, grupos separados, ordem e parsing parcial. Expectativa de capacidade
errada foi rejeitada e saída existente foi preservada. Replays byte a byte
conferiram. Referência autoral; precisão semântica independente não medida.

Nos dois arquivos de conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b,
conservaram-se sete observações sintáticas e zero atribuições de tipos.
Dependências ausentes/pacote incompleto geraram lacunas explícitas e impediram
crédito indevido. Não se demonstrou inexistência do canal limiter, nem limites
de admissão/workers ou comportamento de p.handle.Go.

Evidência atual em [0110/results](../../lab/bancada/concorrencia/0110/results/),
contrato/uso em [README](../../lab/bancada/concorrencia/0110/README.md).
Recibos e corpus v1 não foram reescritos. Próximo passo: contrato para
fornecimento rastreável de dependências e seleção de build, mantendo IO fora
da transformação pura; depois estudar a relação entre reserva e trabalho.

Verificação concluída: crystalline-lint sem violações, go test ./...,
race do observador (incluindo contrastes finais), go vet ./..., builds do
produto/piloto e diff check passaram. A suíte completa usou fixtures loopback
e assinatura Git desativada apenas no processo. Gates de módulos selados
separados não foram reproduzidos. O recibo verification.json vincula artefatos
por SHA256 e registra o runtime Go utilizado.
