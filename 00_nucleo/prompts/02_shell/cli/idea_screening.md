# Prompt Spec L0: CLI da Peneira Explicável

## Objetivo

Expor classificação, calibração agrupada por ideia e validação de revisões como
operações JSON reproduzíveis e offline.

## Restrições

- L2 não conhece provedores remotos.
- `screen` aceita snapshot e modelo opcional.
- `screening calibrate` pareia snapshots e julgamentos por ideia.
- `screening audit` preserva revisões, evidência e rótulos anteriores.
- Estratégia, schemas e arquivos inválidos geram erros acionáveis.
