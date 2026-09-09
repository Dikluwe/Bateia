# 0044 — Recuperação antes da seleção
Origem: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md.
L1 somente orquestra porta opcional, sem IO concreto. RetrieveImplementationEvidence usa budget finito padrão e reporta not_evaluated em provedor ausente. SelectImplementationPaths preserva manifesto/README e antepõe hits de caminhos relativos presentes na árvore, com limite de código. ValidateImplementationHits copia relatório, descarta fontes não extraídas ou hashes/intervalos divergentes e marca partial; não apaga identidade do índice.
Funnel/analyze registra relatório; CollectEvidence busca requisitos e prioriza árvore antes do orçamento, sem promoção do status de capacidades. Provedores sem porta preservam contrato antigo. A seleção por conteúdo é distinta da análise de comportamento. Relatório mostra símbolos com linhas/hash e snippet dos bytes revalidados, no máximo20 linhas por hit/8hits apresentados, JSON retém todas as evidências.

## Passo 0045 — Identidade de relações estáticas

Conferir também RelatedTo por SHA-256 e linha da chamada no arquivo chamador extraído. Referência não extraída, alterada ou fora de intervalo é removida e torna relatório parcial; hit exclusivamente direct_call sem referência válida é descartado. Não mutar slices da entrada. A confirmação identifica bytes de uma relação AST, não execução nem satisfação funcional.

Revisão 0045 após ataque público: intervalo válido deve apontar uma CallExpr por identificador ao símbolo alvo dentro do corpo da função chamadora declarada, não apenas uma linha existente no arquivo. Reparse puro Go dos bytes conferidos; ignorar closures e chamadas sombreadas; erro de parse não confirma referência.

A validação calcula hashes e índice de chamadas uma vez por arquivo em cada releitura, evitando reparse por referência.

0051: ValidateImplementationHits valida hits documentation contra declaração AST e grupo Doc anexado, incluindo símbolo, intervalos da declaração e documentação, após conferência SHA256. Parse por arquivo com cache local; rejeitar arquivo gerado, sem Doc ou intervalo forjado. Não admitir fonte documental como código lexical.
