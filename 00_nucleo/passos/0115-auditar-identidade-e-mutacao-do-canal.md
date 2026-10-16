# 0115 — Auditar identidade e mutação do canal

Capacidade investigativa no Lab sob Tekt A0. Prompts de auditoria, observador,
relações e testes precederam código. v7 acrescenta auditoria por declaração,
com bases go/types, eventos de escrita/exposição/captura e localizadores/hash.

Na fonte conc, bases de Pool.Go/worker são objetos diferentes; escrita em
WithMaxGoroutines e duas capturas localizadas. Nenhum evento demonstra que o
canal é substituído durante execução. Identidade continua inconclusiva;
compartilhamento de declaração não vira igualdade de instância.

13 casos novos e 78 perfis anteriores passaram. Qualidade semântica independente
continua pendente. [Bancada](../../lab/bancada/concorrencia/0115/README.md).
Próximo recorte: guarda da escrita de configuração e ligação entre receivers.
