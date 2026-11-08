# Prompt cristalino — CLI de catálogo

Expor descoberta remota por plano ontológico e peneira offline. Exigir snapshot antes da rede, budgets explícitos e provedor com papel PROJECT_CATALOG. A peneira pode materializar um plano de prova em que toda expansão fica restrita ao repositório `ADVANCE` de origem.

Expor também `discover-surfaces` para README/tópicos e `screen-multiview` offline, preservando snapshots mesmo quando uma superfície fica Unknown.

Segurança: snapshot JSON de screening usa readIdeaJSONInput antes de Unmarshal, conforme security_json_input.md; saídas válidas preservadas.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
