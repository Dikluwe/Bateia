# Prompt R2 — oráculos MCP ampliados

## Compatibilidade não autorreferente

Hashes do certificado, prompts e testes imutáveis do Passo 0025 são verificados antes da comparação. Os schemas legados completos são goldens literais, e `search`, `benchmark` e `recommend` são realmente invocados com o serviço durável presente; nenhuma porta assíncrona pode ser tocada. O comportamento alcançável certificado de benchmark (inclusive seu erro legado para o argumento declarado pelo schema) e a recomendação usam bytes literais independentes, nunca o serializer candidato.

## Fechamento e projeção

- campos extras e ausentes são exercitados na raiz e em todos os objetos aninhados;
- todos os lifecycles preservam chaves nullable e arrays `[]` explícitos depois de decodificar `content[0].text`, não o envelope re-escapado;
- lookup ausente, três falhas terminais controladas e bytes de artifact malformados permanecem categorias distintas e sanitizadas.

O transporte usa apenas buffers e fakes offline.

## Sucessor histórico O3

O hash predecessor de `server.md` permanece registrado em
`historical-hash-invalidation-o2-to-o3.json`; o teste vincula o hash sucessor
exato e os hashes da invalidação e do sucessor. A aceitação continua dependendo
da execução real das três tools legadas, de seus schemas/goldens literais e da
ausência de aquisição de comportamento assíncrono. A extensão de budget não
autoriza regenerar expected bytes nem remover qualquer propriedade legada.
## Sucessor explícito 0031

Vincule os novos hashes do schema/prompt síncrono ao artefato
`historical-contract-supersession.json`. Continue exigindo execução real das três
tools, os byte-goldens e toda a superfície durável não substituída.
