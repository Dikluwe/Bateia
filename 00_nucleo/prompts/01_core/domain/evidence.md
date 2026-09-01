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
