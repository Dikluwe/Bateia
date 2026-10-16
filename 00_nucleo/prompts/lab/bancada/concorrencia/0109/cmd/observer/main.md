# Prompt 0109 — cmd/observer/main

**Regime**: Lab, sem promoção para L1–L4  
**Data**: 2026-10-06

CLI Lab observe --input --output opcional e evaluate --corpus --output. Entrada JSON até 2MiB, campos desconhecidos/múltiplos documentos/duplicados recusados. Saída JSON determinística; --output somente arquivo novo, publicar completo por arquivo temporário e link exclusivo 0600, nunca sobrescrever entrada/saída prévia. Retornar erro em short write e entrada inválida. Não ligar ao produto nem executar fontes. Evaluation com divergências deve emitir relatório e retornar exit code não zero.

0111: observe aceita --environment JSON opcional, leitura regular até 64MiB
com as mesmas verificações de chaves/estrutura. evaluate rejeita esse flag.
Entrada de fontes permanece limitada a 2MiB; ambiente valida limites próprios.
CLI não busca nem gera export data; somente recebe dados congelados.
