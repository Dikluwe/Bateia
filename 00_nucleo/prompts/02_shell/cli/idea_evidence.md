# Idea evidence CLI

Expor `bateia idea evidence plan|collect|assess|evaluate|diagnose`. Planejamento, assessment
e métricas leem e escrevem JSON offline. Coleta recebe o provedor injetado pelo composition root.
Validar entradas com mensagens acionáveis e nunca persistir tokens ou headers.

`plan` exige ideia e peneira; `evaluate` exige assessment, bundles e julgamentos.
`plan --include owner/repo,...` adiciona controles explícitos de calibração sem
alterar scores ou regras.
