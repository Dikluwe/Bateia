# Prompt L0 — `02_shell/cli/idea_radar_export.go`

Implemente o handler de `bateia idea radar export --input <json> [--output <json>]` conforme
0069. Aceite um único JSON estrito; chame o apresentador L2. Sem output, stdout contém
somente JSON; com output, escrita atômica e stdout vazio. Em qualquer erro não emita
parcial nem altere destino/input; rejeite quando output identifica input. Não execute
Node, Archify, rede ou renderização. O dispatcher recebe apenas encaminhamento aditivo.

