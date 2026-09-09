# Evidence domain

Contratos imutáveis para planejar, congelar e avaliar evidências de capacidades em
`repo@commit`. Os estados funcionais, licença, manutenção e operação permanecem
separados. Ausência só pode ser declarada após busca completa; falhas produzem
inconclusão. Nenhum contrato admite saída de LLM, embedding ou inferência remota.

Cada alvo registra `COLLECTED`, `NOT_FOUND`, `FAILED`, `SKIPPED_BUDGET` ou
`UNSUPPORTED`. Identidade `owner/repo` e revisão SHA são dimensões separadas.

Na versão 2, cada regra declara família, proveniência, alvos justificáveis e se
`PROVEN` exige fontes independentes. Versões antigas permanecem legíveis para
replay, enquanto catálogos v2 são validados estritamente.

# 0044 — Porta e evidência de recuperação de implementação
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md (contrato congelado).
Tipos de domínio para orçamento, relatório, inputs hashados e hits de AST; porta opcional de recuperação sem IO concreto em L1. Busca pode ser partial/not_evaluated; hits não são comprovação de capacidade. Preservar hashes da varredura e conferir identidade antes de exibir definição. Campos opcionais em resultados de extração/bundles mantêm compatibilidade com consumidores existentes.
