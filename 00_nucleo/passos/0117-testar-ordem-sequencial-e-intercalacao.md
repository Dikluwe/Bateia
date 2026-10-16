# 0117 — Testar ordem sequencial e intercalação

Bancada autoral no Lab sob Tekt A0. Contratos do modelo/testes precederam código.
Sem alteração do produto/observador v8. Marcador e limit atômicos abstratos;
rendezvous permite intercalação check -> initialize -> write sem data race.

Quatro cenários verificam ordens sequenciais e controle concorrente. Modelo
aceita escrita quando initialization ocorre depois do check e antes do store.
Resultado restrito ao modelo, sem diagnóstico de corrida ou execução de conc.
Observação estática conserva atomicity=UNRESOLVED. Referência independente
continua pendente. [Bancada](../../lab/bancada/concorrencia/0117/README.md).

Próximo recorte: revisar obrigações/alegações de configuração e concorrência
separadamente das implementações, usando os testemunhos como controles.
