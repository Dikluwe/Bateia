# 0117 — Cenários sequenciais e concorrentes em modelo autoral

Bancada executável de validação de limites, sem alterar o observador v8 ou o
produto. O modelo é autoral; não é uma implementação de conc e não executa
fontes do candidato. A referência continua não independente.

## Hipótese e abstração

Uma checagem anterior à escrita pode respeitar o contrato sequencial e ainda
não constituir uma operação atômica de check+write. Model substitui o estado
por initialized e limit, variáveis atômicas separadas. Initialize marca o estado;
Configure checa initialized, permite um hook de rendezvous de bancada e grava
limit. O hook pausa o escritor depois do check, antes do store.

Atomicidade por operação evita data race no modelo. Não transforma as duas
operações em transação. Esse modelo omite canal tasks, sync.Once, callbacks,
workers e aliases de receiver; não autoriza transferir resultado para conc.
Snapshot só é lido após rendezvous/joins, sem alegar snapshot conjunto atômico.

## Testemunhos executados

| Cenário | Resultado no modelo | Escopo |
| --- | --- | --- |
| configurar 3, depois inicializar | configuração aceita, limit=3 | ordem sequencial |
| inicializar, depois reconfigurar para 7 | configuração rejeitada, limit=3 preservado | ordem sequencial |
| check, pausa, initialize, retoma store | configuração aceita, limit=7 após initialize | intercalação artificial |
| initialize antes do check em outra goroutine | rejeitada, limit=0, hook não chamado | controle de ordem oposta |

Agendamento usa testing/synctest e canais; sleeps/tempo não determinam vereditos.
A intercalação evidencia uma limitação do padrão no modelo, não execução ou
corrida no candidato. Não foi julgado se a intercalação pertence ao contrato
de uso suportado pelo candidato. Passar no detector de data race não prova atomicidade.

A observação estática congelada de conc permanece a mesma para as hipóteses
sequencial e concorrente, com atomicity=UNRESOLVED. Sua guarda não recebe uma
conclusão mais forte por estes testemunhos. Não foram medidos falsos positivos,
precisão semântica independente, desempenho ou bounded concurrency real.

## Reprodução

Na raiz:

```sh
go test -v ./lab/bancada/concorrencia/0117 -count=1
go test -race ./lab/bancada/concorrencia/0117 -count=1
python3 lab/bancada/concorrencia/0117/verify.py /tmp/bateia-observer-0116
```

Se o binário anterior não estiver disponível, compilar o observador ativo:

```sh
go build -o /tmp/bateia-observer-0116 ./lab/bancada/concorrencia/0109/cmd/observer
```

verify.py executa somente o pacote do modelo duas vezes, conserva logs JSONL e
compara testemunhos normalizados. Logs brutos contêm timestamps/durações e não
precisam ser byte idênticos. O relatório vincula observação estática/input/hash
do ambiente a quatro testemunhos; observação de origem tem replay byte idêntico.
Receipts anteriores preservados; candidato não executado.

Próxima avaliação: corpus de obrigações/alegações revisado separadamente, que
não confunda checagem sequencial, data race e atomicidade. O modelo fornece um
caso de contraste executável para essas alegações, sem substituir o estudo da
instância/ciclo de uso no programa real.
