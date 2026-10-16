# Prompt — Declarações e evidência, 0121

Regime Lab; data 2026-10-07. Pacote puro declarations. Inspect recebe fonte
Request v8, Environment opcional, documentos path/content/hash e associações
document_path/first_line/last_line/owner Location. Primeiro executar observador
v8 puro; não modificar sua versão, campos ou decisões. Validar arquivos
documentais sem parsing Go, com limites equivalentes a v8 e sem exigir .go, incluindo conflito
de caminho com fonte. Vazio permitido para documentos/associações. Até 1000
associações, 10.000 declarações/links totais.

Grupos de comentários preservam bytes exatos, incluindo diretiva line que
integre o grupo anexado. Diretiva não altera posições físicas nem vira contrato.

ParseComments/AllErrors; trabalhar apenas arquivos sem erro de parsing para
declarações, sem depender de sucesso de tipos. Identificar FuncDecl cujo nome
é exportado, mesmo com receiver privado, e campos de método exportado com
FuncType em interface nomeada. Doc apenas anexado por AST; comentários soltos,
trailing, internos ou pacote não são documentação da operação. Interfaces e
implementações são donos distintos, sem casamento nominal.

Posição física com PositionFor(false), offsets de bytes, texto exato do grupo
de comentários. Dono inclui path/hash, intervalo inteiro de declaração, nome
sintático e kind. ID determinístico por dono+localizador documental+associação.
Associação externa aponta exatamente para dono AST fornecido; rejeitar hash,
posição, linha/coluna ou dono discrepante, spans fora de limites e duplicatas.
Texto externo é intervalo de linhas exato sem quebra final adicionada.

Declaration Authority DOCUMENTATION_TEXT_ONLY; Association
AST_ATTACHED_GO_DOC ou CALLER_PROVIDED_BINDING. Fatos vinculados por contenção
lexical no dono: observações, from das relações e guarda/write de configuração;
limitado ao mesmo path/hash e offsets. Interface não possui corpo: sem vínculos
por mera igualdade de nome. Lexical containment inclui closures, sem executar
ou provar que comentário se aplica a cada caminho. Ordenar resultados por
path/offset documental e ID, links por localização/kind; invariância à ordem.
Report inclui observer_report original, schema api-declarations-v1,
declarations, scope e limites explícitos. Nenhum veredito semântico.

Report inclui api_owners mesmo quando dono não tem Doc, permitindo associação
externa exata sem inventar identidade por nome. Até 10.000 donos, ordenados por
path/offset; inventário sintático não prova exposição externa ou implementação.
