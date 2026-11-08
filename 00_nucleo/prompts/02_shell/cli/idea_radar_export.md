# Prompt L0 — `02_shell/cli/idea_radar_export.go`

Implemente o handler de `bateia idea radar export --input <json> [--output <json>]` conforme
0069. Aceite um único JSON estrito; chame o apresentador L2. Sem output, stdout contém
somente JSON; com output, escrita atômica e stdout vazio. Em qualquer erro não emita
parcial nem altere destino/input; rejeite quando output identifica input. Não execute
Node, Archify, rede ou renderização. O dispatcher recebe apenas encaminhamento aditivo.


## Segurança de entrada

Abrir com O_NONBLOCK e ler bytes limitados do mesmo descritor usado nas verificações de aliases; recusar não regulares/FIFO/>16MiB. Aplicar validateIdeaJSONStructure antes do Decoder estrito. Preservar escrita atômica, identidade de entrada, saída válida e certificados históricos.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
