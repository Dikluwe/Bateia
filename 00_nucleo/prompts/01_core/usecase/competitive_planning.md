# Planejamento verificável de autofagia e Katamari — 0041

Implementar helpers puros para validar fontes e autoridade, ordenar/deduplicar entradas tratadas como conjuntos, aplicar todos os gates e produzir opções com próxima verificação e fronteira de integração ainda desconhecida. Nenhum helper executa código ou aceita PROVEN sem rastreabilidade declarada. Hash de fontes precisa ser SHA256 hex, referência deve corresponder ao candidato. Canonicalização opera somente sobre inputs de conjuntos deste fluxo, não listas de instruções ordenadas. Rejeitar revisões conhecidamente mutáveis. Não transformar validação estrutural de ponteiro em validação comportamental.

## Revisão 0043 — autoridade e ligação
Aplicar o contrato comum de evidence_authority antes de autorizar capacidade ou operação. Evidência lexical/documental não ganha PROVEN. Negação deve conservar polaridade. Na extração de mecanismos aceitar SHA-256 hexadecimal produzido pela avaliação e formato prefixado legado, verificando conteúdo e projeções de tipo file→symbol e test→test_implementation. Preservar IDs reais: catálogo sem regra compatível retorna lacuna explícita, nunca traduz nomes em prova.
