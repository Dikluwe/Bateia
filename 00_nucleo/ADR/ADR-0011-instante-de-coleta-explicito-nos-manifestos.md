# ADR-0011 — Instante de coleta explícito nos manifestos

**Status:** ACEITO PARA O RECORTE 0128
**Data:** 2026-10-28
**Relacionado:** ADR-0002, ADR-0008 e passo 0128.

## Contexto

Os parsers de manifestos são declarados puros, mas 43 arquivos consultam o relógio para preencher ConsultadoEm. O instante representa a coleta e não pode ser deduzido do texto. Essa ambiguidade impede determinismo mesmo com bytes e referência iguais. O requisito histórico de preencher a hora atual nos parsers conflita com a pureza de L1.

## Decisão

1. Os parsers específicos e ParseManifest conservam assinaturas e parsing; não consultam relógio e deixam ConsultadoEm zerado, indicando instante não fornecido.
2. Adicionar ParseManifestAt(filename, ref, content, collectedAt). O instante é entrada explícita, convertido a UTC; zero permanece desconhecido. A função não consulta relógio nem cria datas substitutas.
3. GitHub e coleta local em L3 registram o instante da leitura e o fornecem à função. O funil, ao reanalisar bytes extraídos, utiliza ExtractedRepository.ConsultadoEm já recebido da porta.
4. Preservar nomes/tags JSON, conteúdo, dependências, erros e arquivos históricos. Chamadores diretos dos parsers deixam de receber uma hora fabricada; devem usar a API explícita quando possuem informação de coleta. O campo JSON permanece presente, incluindo valor zero quando desconhecido.

## Limites

Esse recorte elimina o relógio dos parsers de manifestos. Não elimina a telemetria temporal do funil, de idea_run ou de lab_extractor, nem altera deadlines ou efeitos da CLI. Não redefine elegibilidade de reúso, integridade de bytes ou força de evidência. Relatórios anteriores não são recalculados.

## Ordem estável dos objetos de dependências

A revisão encontrou iteração não canônica de mapas em package.json, composer.json, box.json, haxelib.json e dub.json. Ordenar as chaves dentro de cada grupo, preservando ordem dos grupos e nomes/versões/tipos/exclusões. Testar múltiplas dependências e grupos repetidamente. A ordem dessas listas passa a ser definida; a caracterização do corpus temporal permanece intacta. Artefatos antigos não são reordenados ou re-hasheados.
