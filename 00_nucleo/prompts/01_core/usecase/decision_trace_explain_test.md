# Testes do caso de uso de explicação de rastreabilidade

Validar `ExplainDecisionTrace` com entrada congelada e todas as famílias de registros: report/policy/idea/rule_versions, decisão, capacidades, razões, bloqueadores, componentes, composições, membros, revisões e verificações técnicas. Resolver cada JSON Pointer RFC 6901 no JSON original produzido com MarshalIndent; compactar o registro sem reordenar chaves e verificar SHA-256 igual à fingerprint do nó. Tratar documentos como identificadores lógicos, nunca caminhos para coleta.

Uma mesma fonte deve aparecer uma vez no índice com todos os localizadores, inclusive posições repetidas em arrays, componentes repetidos e membros. Duas capacidades com o mesmo registro de candidato devem manter IDs distintos mesmo com fingerprints iguais. Entradas/localizadores têm ordem canônica e não são duplicados; explicações repetidas são idênticas. A explicação não muda bytes, schema v1 ou hash do grafo, build input ou decisão histórica.

Recusar esquema incorreto, hash ausente/tamper, alterações de declarações/versões, seletores desconhecidos ou duplicados e excesso de orçamento estrutural, textual, vetorial ou de seleção. Toda falha devolve resultado vazio, nunca uma explicação parcial. Verificar seleção opcional independente da ordem dos IDs.

Statements são projeções explicitamente declarativas, limitadas a 1024 bytes UTF-8 válidos, com truncamento marcado. Cobrir texto multibyte no limite; localizadores e fingerprints continuam vinculados ao registro original não truncado. Não promover fontes sem digest/autenticação nem afirmar prova de capacidade ou licenciamento. Não executar regras decisórias nem realizar I/O.
