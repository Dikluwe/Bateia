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

Segurança: readJSONFile e modelo/calibração usam readIdeaJSONInput antes de Unmarshal, conforme security_json_input.md; saídas válidas preservadas.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
