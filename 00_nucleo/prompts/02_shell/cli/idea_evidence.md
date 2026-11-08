# Idea evidence CLI

Expor `bateia idea evidence plan|collect|assess|evaluate|diagnose`. Planejamento, assessment
e métricas leem e escrevem JSON offline. Coleta recebe o provedor injetado pelo composition root.
Validar entradas com mensagens acionáveis e nunca persistir tokens ou headers.

`plan` exige ideia e peneira; `evaluate` exige assessment, bundles e julgamentos.
`plan --include owner/repo,...` adiciona controles explícitos de calibração sem
alterar scores ou regras.

Segurança: planos e bundles JSON usam readIdeaJSONInput antes de Unmarshal (incluindo fallback objeto/lista), conforme security_json_input.md; saídas válidas preservadas.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

## 0153 — Coleta com pergunta
Adicionar --question, --familiar repetível e --max-passages. Se question foi informado, delegar nova coordenação L1 e emitir envelope bundles+exploration via publicação protegida, recusando alias do plano e posicionais. Flags de contexto sem question são erro. Sem flags novas manter array legado. Não executar exploração na casca.

## 0155 — Ideias de algoritmos
Flag --algorithm-ideas booleana exige --question e é passada ao pedido L1; sem flags novas preservar lista e saídas.
