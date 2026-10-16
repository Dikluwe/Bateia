# 0113 — Rastrear parâmetros de callback

Capacidade investigativa no Lab, governada por Tekt A0. Prompts de callback,
relações, observador e testes precederam a materialização. v5 acrescenta relações
separadas por identidade de parâmetros, com escopo/condições e hashes.

No conc congelado, literals de Pool.Go ligam-se ao parâmetro de WaitGroup.Go;
o uso capturado na goroutine é repassado a Catcher.Try e seu parâmetro é
invocado. As declarações coincidem nos pontos de junção. O relatório não
promove esses vínculos a execução verificada ou limite global.

Mutação/endereço/sombra/indireção/aliases/variádicos possuem limites explícitos.
15 casos autorais e 44 perfis anteriores passaram; referência independente
continua pendente. [Bancada](../../lab/bancada/concorrencia/0113/README.md).
Próximo recorte: alcance e ordem, depois estabilidade de instância do canal.
