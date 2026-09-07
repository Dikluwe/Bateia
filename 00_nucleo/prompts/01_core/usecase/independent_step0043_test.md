# Testes independentes puros 0043
Autoria: agente testador separado do implementador; contexto permitido: AGENTS, spec congelada 0043, tipos/assinaturas públicos e dados públicos do catálogo. Sem leitura de implementações/diffs, testes históricos ou oráculos. Isolamento físico não atestado: filesystem compartilhado.

Antes do código: frases `cannot run locally` e `never uploads code to the cloud` não provam os respectivos contrários. README positivo permite declaração parcial, nunca comprovação de execução local. Avaliação inteiramente sem fontes, incluindo gates favoráveis, não autoriza FORK nem cobertura completa; candidato opaco e ausência de candidatos mantêm limites. Mesmas entradas produzem saída idêntica.

Mecanismos: bytes hash SHA256 hexadecimal canônico recebidos pelo assessor real e consumidor real não podem gerar HASH_MISMATCH; regra desconhecida não pode produzir mecanismo. Catálogo público embutido como literal, sem I/O em L1; não fabricar EvidenceSource.RuleID. Ausência de mecanismo é admissível quando não há testemunho pertinente. Exigir ao menos uma fonte real do produtor para evitar controle vacuamente aprovado.

Plano: BuildInvestigationPlan produz IDs por texto/tipo estáveis sob permutação, zero candidatos encerra NO_CANDIDATES_IN_SCOPE, desconhecido tem ação e aceitação. Status sem fontes não vira fato observado comprovado. Consultas efetivas e originais/fallback permanecem rastreáveis quando há tentativas.

## Correção de fixture pública antes de reexecutar
A primeira fixture de mecanismo usou “code analysis”, que o catálogo público DefaultEvidenceRules não reconhece, portanto gerou zero fontes. Trocar somente entrada por “static analysis”, termo declarado publicamente na regra static-analysis; manter todas as assertions. Hash anterior e conteúdo arquivados em benchmark0043. Não é adaptação ao patch: o erro é construção sem testemunho, conhecido via saída pública do catálogo.

## Correção da fixture de catálogo fechado
O catálogo público é fechado sobre bytes distribuídos: a cópia compactada por JSON conserva dados mas não identidade do artefato. Substituir somente literal pelos bytes públicos exatos de 04_wiring/mechanism-rules.json e manter SHA256 calculado sobre eles. Sem I/O runtime L1. Arquivo anterior arquivado, assertions intactas; ContentHash hexadecimal do assessor continua sob teste.
