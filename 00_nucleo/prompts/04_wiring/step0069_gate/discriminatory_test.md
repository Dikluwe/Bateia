# Gate discriminatório binding-free — passo 0069

## Autoridade e isolamento

Este prompt espelha `04_wiring/step0069_gate/discriminatory_test.go`. A autoridade é o
passo congelado `00_nucleo/passos/0069-exportar-radar-para-archify.md`, o contrato e o
baseline em `00_nucleo/benchmark/0069`. O gate é anterior ao candidato, usa somente a
biblioteca padrão e não importa, descobre ou executa código de produção.

O papel é adversário. Ele pode escrever somente as fixtures adversariais, este prompt e o
gate. Não pode corrigir implementação, alterar contrato, criar selo ou emitir veredito de
materialização.

## Receipt observável

Construa em memória um controle completo com origens registradas e duas origens
desconhecidas, perspectivas, targets, projeções, observações com stance, cards e probes.
IDs devem ser prefixo tipado mais `hex(bytes UTF-8 da identidade canônica)`. O receipt
registra estados públicos das autoridades, digest do input antes/depois, limites, destino
de escrita, atomic rename, documento e 20 resultados de serializações semanticamente
equivalentes.

Classifique sem bindings:

- controle íntegro como `Preserved`;
- entrada/schema opaco, autoridade `UNKNOWN` sem parcial ou budget esgotado sem parcial
  como `Unknown`;
- toda emissão parcial nesses estados e toda divergência observável como `Violated`.

Valide forma Dataflow v1, três stages em ordem, nodes por stage/identidade, rows canônicas,
flows de projeção e observação exatos e ordenados, cards por target e probes na ordem do
plano. Rejeite relações extras, probe-flow, origem desconhecida coalescida ou promovida,
IDs inválidos/colidentes, stance trocada, perdas, target errado, duplicatas, mutação do
input, excesso com parcial e receipt de escrita não atômica.

## Poder discriminatório

Carregue `00_nucleo/benchmark/0069/adversarial-fixtures.json`; exija pelo menos M01–M16,
classifique cada caso duas vezes, execute o corpus também em 20 ordens determinísticas e
exija o mesmo mapa de resultados. `Unknown` só satisfaz fixtures opacas explícitas.

Calcule `mutation_score = mutações M* classificadas Violated / mutações M*` e exija o
valor definido na fixture, atualmente `1.0`. Registre o score no output de `go test`.
O gate atesta apenas que o contrato observável discrimina o corpus pré-candidato; não sela
o contrato e não afirma equivalência funcional de qualquer implementação.
