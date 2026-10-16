# Laudo 0116 — Guarda anterior à configuração

O observador v8 acrescenta configuration_guards, separado dos fatos históricos.
Para escrita direta de canal candidato, um padrão limitado exige guard direto
anterior no mesmo corpo e base, com predicado de campo != nil e panic builtin
no método fornecido. Estado PRIOR_GUARD_SAME_BASE_OBJECT, autoridade
LIMITED_GUARD_PATTERN; atomicidade sempre inconclusiva.

Em conc congelado, write de limiter em pool/pool.go:96 é precedido pelo call
em :92, destino panicIfInitialized em :111, predicado p.tasks != nil.
Assignment à sentinela tasks localizado em init, :104, dentro de callback.
Quando o predicado é verdadeiro, esse guard não retorna normalmente. Isso não
prova o valor de tasks na escrita, execução de init na mesma instância ou
proteção contra uso concorrente. Identidade do canal/workers limitados pendentes.

Guard ausente/condicional/adiado/posterior, outra base, rebind/endereço/capture,
jumps, nomes sombreados e recover ficam fora do recorte. Efeito intermediário
não apaga a existência de um check anterior e continua explicitamente não
resolvido. Writes da sentinela são inventário por declaração, não efeitos de
inicialização comprovados. Nenhum nome sozinho autoriza o padrão.

20 controles autorais passaram; expectativa errada foi discriminada. 91 perfis
anteriores passaram e projeções reais conservaram igualdade. Replays byte
idênticos, localizadores/hash, ordem e limite explícito verificados. Sem ambiente,
nenhuma guarda em conc. Referência semântica independente não medida; candidato
não executado. Evidência histórica preservada.

crystalline-lint sem violações, go test ./..., race do observador, go vet ./...,
builds produto/piloto e diff check passaram. Fixtures usaram loopback e
assinatura Git desativada apenas no processo. Produto não importa Lab; gates
selados com módulos próprios não reproduzidos.

[Recibo com hashes](../../lab/bancada/concorrencia/0116/results/verification.json)
e [contrato/reprodução](../../lab/bancada/concorrencia/0116/README.md).
Próxima avaliação: distinguir configuração sequencial de garantias concorrentes,
com ligação dos receivers e revisão independente das conclusões.
