# Passo0047 — pontuação lexical pelo corpus AST

Origem exclusiva de `01_core/parser/corpus_scoring.go`, subordinada ao contrato
`00_nucleo/passos/0047-ponderar-termos-pelo-corpus.md`. Funções puras sem IO.
Coletar cada declaração nomeada elegível uma única vez no parser; cada nome em
ValueSpec constitui um documento como no universo de hits existente. N inclui
não matches, tipos e valores. DF conta uma vez por documento a união de termos do
nome e corpo. Repetir entrada idêntica não altera corpus. Não alterar vocabulário,
exclusões, escopo AST ou identidade da coleta existente.

IDF = log1p((N-df+0.5)/(df+0.5)). Contribuição nome = 2*IDF; contribuição
corpo = 3*IDF/(0.25+0.75*comprimento/média). Comprimento é cardinalidade do
conjunto lexical de identificadores AST únicos do corpo (após expansão vigente),
nunca bytes/comentários/literais. Média considera funções com corpo do recorte,
inclusive não matches; piso 1 evita divisão por zero. Registrar ambos campos se
ambos casarem. Guardar contribuições finitas em ordem termo/campo e somar nessa
ordem. Score lexical = ceil(900*soma/(8+soma)); bônus 1000 função com corpo
pertinente; bônus 100000 identificador exato. A saturação preserva faixas de
prioridade absoluta. Nenhuma contribuição existe para termos ausentes.

Reason declara fórmula, quantização, comprimento, média e limites do recorte. N/DF
não estimam ecossistema nem pesos são probabilidades. Scorer de requisitos usa a
mesma estatística e função. Expansão de chamadas v2 permanece; helper sem termos
não recebe contribuições artificiais. Método go_ast_corpus_v1.
