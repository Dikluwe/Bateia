# Laudo 0035 — Configuração durável em L3

**Resultado**: aprovado

A configuração durável foi movida de L4 para `03_infra/durableconfig`. O pacote passou a
ser responsável por ambiente, diretórios, segredo local e limites. `04_wiring` apenas
carrega a configuração e a injeta.

Os testes cobrem defaults, permissões, recusa de symlinks e arquivos inseguros,
concorrência na criação inicial e mensagens de erro. O comportamento público foi
preservado e duas condições de corrida da implementação anterior deixaram de existir.
