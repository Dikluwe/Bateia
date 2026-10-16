# 0112 — Relações estruturais com identidade de tipos

Lab, transformação pura. Indexar declarações de funções por objeto go/types
somente em grupos checados sem erros. Acrescentar relations separadas dos
perfis sintáticos existentes; versão v4, recibos anteriores preservados.
Tipos de relação: STATIC_CALL_TARGET (incluindo modo DIRECT/DEFERRED/GO),
FUNCTION_CONTAINS_GO, CALL_ARGUMENT_FUNCTION, SEND_CASE_CONTAINS_GO_CALL,
DEFERRED_RECEIVER_RECEIVE e SAME_DECLARED_CHANNEL_FIELD.

STATIC_CALL_TARGET requer objeto *types.Func e definição fornecida; chamadas
indiretas por variáveis/interface não autorizam destino. Seleções usam identidade
de objeto, não nome. SEND_CASE_CONTAINS_GO_CALL vincula envio tipado de canal
na comunicação select a chamada no corpo do mesmo case cujo destino contém
um GoStmt no próprio corpo (excluindo closures). É contenção lexical, não
alcance/dominância ou comprovação de reserva. Condições if/else/loops/select
permanecem sintáticas. Closures reiniciam condições/contexto do case.

DEFERRED_RECEIVER_RECEIVE requer chamada defer de método de receptor canal,
cuja definição contém receive do objeto receptor exato, com condições do alvo
separadas das condições do chamador. Não implica que o defer é sempre registrado.
SAME_DECLARED_CHANNEL_FIELD liga envio e receptor de chamada defer pelo mesmo
objeto de campo *types.Var canal. Não prova mesma instância/valor: receivers de
funções distintas, aliases e mutação não resolvidos.

Localizadores físicos/hash em ambas extremidades; autoridade
STATIC_TYPES_AND_SYNTAX, call_mode explícito. Não chamar grafo de prova de fluxo,
limite global, callback executado ou ciclo balanceado. Ordenação determinística;
até 10.000 relações, atingir limite é erro explícito sem publicação parcial.
Emitir lacunas FLOW_PATH_FEASIBILITY_UNRESOLVED,
CALLBACK_EXECUTION_UNRESOLVED e CHANNEL_INSTANCE_IDENTITY_UNRESOLVED quando
houver relações correspondentes. Nenhuma execução nem novo acesso IO.

0113: coletar sites de chamada tipados com escopo/condições/modo para analisador
de callback separado. Reportar callback_relations sem alterar relations v4;
limite total compartilhado, erro sem publicação parcial.

0114: usar funções e calls já tipados para projeção de alcance e ordem limitada
independente, com limite próprio de 10.000 fatos. Não alterar relações antigas.

0115: conservar bases simples dos usos de campo e pares existentes para auditor
separado por declaração; não modificar relações históricas nem promover instância.

0116: aplicar observador separado de guarda a canais dos pares existentes,
funções e chamadas resolvidas; não alterar relações anteriores.
