# Prompt espelhado — gate discriminatório offline R3 do Passo 0026

Autoria independente pré-candidato. Materialize um pacote de teste Go isolado em L4 que
compile e execute sem qualquer implementação candidata. Leia apenas os casos adversariais
públicos congelados. Para P01–P05 execute fixtures positivas concretas. Para M01–M60
execute modelos de referência e mutações/faults observáveis concretos; nunca derive o
resultado da presença de bindings. Preserve O01–O05 como `Unknown`.

O gate deve repetir as 70 avaliações três vezes em cada uma de três ordens (lexical,
reversa e intercalada), exigir os resultados congelados e comparar o digest normalizado
das nove rodadas. Inclua concorrência real para criação única/quota, HMAC e comparação
constante para cursores, cursores re-assinados com versão/hash/offset inválidos, injeção
real de falhas nos prefixos de persistência, reabertura lógica de quota cheia, decoding do
texto MCP antes do JSON interno, AST via parser Go, golden literais legados e canários.

Não importe pacotes da implementação, não leia código candidato e não escolha nomes de
arquivos, índices, segredos ou helpers privados. O pacote existe para o verificador
executá-lo pré-candidato com `go test ./04_wiring/r3_gate`.
