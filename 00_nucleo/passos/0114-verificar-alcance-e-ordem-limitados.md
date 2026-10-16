# 0114 — Verificar alcance e ordem limitados

Capacidade investigativa no Lab sob Tekt A0. Contratos de alcance, relações,
observador e testes precederam código. Projeções v6 independentes preservam
relações históricas. Return/builtin panic/if terminante distinguem sites
bloqueados; jumps ficam inconclusivos; guardas/loops não têm viabilidade provada.

Na fonte conc, defer de worker precede os sites de tarefa no seu próprio
escopo, de modo condicional ao alcance desses sites. Registro de defer não
certifica liberação concluída, estabilidade do canal ou workers limitados.

19 casos autorais e 59 perfis anteriores passaram. Avaliação semântica
independente continua pendente. [Bancada](../../lab/bancada/concorrencia/0114/README.md).
Próximo recorte: identidade e estabilidade da instância do canal.
