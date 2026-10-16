# 0115 — Auditoria limitada de identidade de canal

Lab puro v7. channel_identity separado de fatos anteriores. Para pares
SAME_DECLARED_CHANNEL_FIELD, registrar base de acesso por objeto go/types:
SAME_BASE_OBJECT / DISTINCT_BASE_OBJECTS / UNRESOLVED_BASE_EXPRESSION. Somente
selector sobre Ident simples, com parênteses removidos. Igualdade de base não
certifica valor do campo, instância ou estabilidade concorrente.

Auditar fontes checadas do tipo de campo: FIELD_WRITE em lhs selector (inclusive
alias simples do receiver), AGGREGATE_OR_BASE_WRITE quando lhs tem tipo struct/
ponteiro cuja declaração contém campo candidato, FIELD_ADDRESS_EXPOSED em
&campo, RECEIVER_ARGUMENT_EXPOSURE em argumentos de tipo struct/ponteiro que
contém o campo e RECEIVER_CAPTURE em closure usando o objeto de uma das bases
do par candidato. Inventário por declaração cobre várias instâncias e trechos
não alcançados; não classificar hazards como defeitos ou mutação daquele canal.
Não resolver ordem, aliases transitivos, unsafe/reflection ou mutação externa.

State sempre INSTANCE_IDENTITY_UNRESOLVED, autoridade DECLARATION_SCOPED_AUDIT.
Localizar cada evento e pontas com posição/hash, ordenar/deduplicar, até 10.000
fatos auditados, falha sem publicação. Ausência de evento não certifica estabilidade.
Preservar lacunas globais/recibos e nenhuma execução. conc deve registrar write
WithMaxGoroutines e bases diferentes entre Pool.Go e worker; captures explícitos.
