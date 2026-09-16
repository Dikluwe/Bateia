# Prompt L0 — `02_shell/cli/algorithm_history.go`

Implemente a CLI offline `algorithm history analyze --input FILE [--output FILE]` conforme
0070, chamando a autoridade L1. Aceite um JSON estrito, serialize COMPLETE ou UNKNOWN
com dois espaços e newline, use stdout exclusivo sem output e escrita atômica silenciosa
com output. Rejeite alias do input e preserve input/destino em erro. Sem rede, busca,
LLM, Archify ou import de L3; o dispatcher recebe apenas encaminhamento aditivo.
