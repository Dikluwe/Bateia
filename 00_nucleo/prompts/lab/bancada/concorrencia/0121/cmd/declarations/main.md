# Prompt — CLI de declarações, 0121

Regime Lab, 2026-10-07. --input JSON obrigatório, --output arquivo novo opcional.
Entrada até 64 MiB (ambiente exportado incluído), JSON estrito: desconhecidos,
duplicatas, múltiplos documentos e nesting >64 rejeitados. Chamar Inspect puro.
Sem rede, execução de fontes ou carregamento de dependências implícitas. Saída
JSON determinística com newline. Publicar por hardlink de temporário completo
0600, recusando substituir alvo/alias existente. Erros no stderr e exit 1.
