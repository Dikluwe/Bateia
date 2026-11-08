# Prompt L0 — `02_shell/cli/algorithm_history.go`

Implemente a CLI offline `algorithm history analyze --input FILE [--output FILE]` conforme
0070, chamando a autoridade L1. Aceite um JSON estrito, serialize COMPLETE ou UNKNOWN
com dois espaços e newline, use stdout exclusivo sem output e escrita atômica silenciosa
com output. Rejeite alias do input e preserve input/destino em erro. Sem rede, busca,
LLM, Archify ou import de L3; o dispatcher recebe apenas encaminhamento aditivo.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
