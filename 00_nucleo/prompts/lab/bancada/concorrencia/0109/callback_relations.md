# 0113 — Passagem estrutural de callbacks

Lab puro, versão v5; callback_relations separado das relações v4 para preservar
perfis históricos. Indexar parâmetros função por objetos go/types/declarações.
CALLBACK_ARGUMENT_BINDING liga literal função ao parâmetro de destino estático.
CALLBACK_PARAMETER_FORWARDING liga uso de parâmetro função ao parâmetro de
outro destino estático. CALLBACK_PARAMETER_INVOCATION liga declaração de
parâmetro à chamada via seu objeto exato. ExecutionScope localiza função/closure
em que ocorre o uso, e condições/modo permanecem locais desse escopo.

Apenas argumentos diretos, aridade exata, funções não variádicas e destinos
fornecidos sem erro. Não rastrear aliases locais, interface, método via variável,
chamadas genéricas, transformação/armazenamento, mutação ou passagem por ponteiro.
Qualquer atribuição ou address-taking do parâmetro no corpo, inclusive closures,
suprime seu crédito em bindings/forwarding/invocation. Sombra por outro objeto
não recebe crédito. Forwarding parte da declaração/uso identificado, não texto.
Não resolver alcançabilidade ou afirmar que uma execução do argumento ocorreu.

Localizadores físicos/hash em ambas pontas e escopo, autoridade
STATIC_TYPES_AND_SYNTAX. Lacuna CALLBACK_EXECUTION_UNRESOLVED permanece;
CALLBACK_VALUE_FLOW_UNRESOLVED marca parâmetros mutados/address-taken.
Limite combinado de 10.000 relações com falha explícita; ordenação estável.
Fonte conc: dois bindings de callbacks de Pool.Go em WaitGroup.Go, forwarding
do f capturado na goroutine para Catcher.Try e invocação do parâmetro em Try.
Relações locais não certificam um caminho de execução conectado ou workers.
