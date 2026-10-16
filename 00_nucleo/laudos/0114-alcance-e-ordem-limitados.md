# Laudo 0114 — Alcance e ordem limitados

O observador v6 acrescenta call_site_reachability e defer_order, mantendo
relações e callbacks históricos. Return, builtin panic direto e if com dois
ramos terminantes distinguem sites após terminadores. Goto/labels/branches
ficam inconclusivos por escopo. Closure tem entrada própria. A expressão de
return é analisada antes da interrupção. Guardas e saída de loops/select não
têm viabilidade provada; POSSIBLY_REACHABLE é sobreaproximação, não execução.

Na fonte conc congelada, 144 sites POSSIBLY_REACHABLE e dez vínculos de ordem.
O registro de defer em pool/pool.go:153 precede sites :156 e :160 no corpo de
worker, condicionado ao alcance normal desses sites. Defer condicional ou em
outro escopo não vira prefixo obrigatório. Registrar defer não prova a execução
da liberação ou a conclusão do callback. A identidade/estabilidade do canal
permanece pendente, assim como o limite global de concorrência.

19 controles autorais passaram: retorno/if, panic/sombra/defer/go, case,
closure, jumps, loop e registro de defer. Expectativa errada foi discriminada.
Os 59 perfis anteriores passaram e as 80 relações/18 vínculos de callback
reais permaneceram iguais. Replays byte idênticos; nenhum fato novo em conc
sem ambiente. Testes adicionais conferiram posições/hash, reordenação, excesso
de fatos e expressões de tipo constantes sem alegação de invocação runtime.
Precisão semântica independente não medida. Candidato não foi executado.

[Contrato e reprodução](../../lab/bancada/concorrencia/0114/README.md).
Próximo recorte: identidade e estabilidade da instância do canal, sem promover
compartilhamento de declaração a igualdade de valor. Verificação final e hashes
constam no recibo de bancada.

Verificação final concluída: crystalline-lint sem violações, go test ./...,
race do observador, go vet ./..., builds produto/piloto e diff check passaram.
Fixtures usaram loopback e assinatura Git desativada somente no processo.
Produto não importa Lab; gates selados de módulos separados não reproduzidos.
[Recibo com hashes](../../lab/bancada/concorrencia/0114/results/verification.json).
