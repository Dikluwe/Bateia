# 0128 — Explicitar instante de coleta nos manifestos

**Data:** 2026-10-28
**Capacidade:** parsing de manifestos e metadado de coleta.
**Responsáveis:** transformação pura em L1; aquisição temporal em adaptadores L3; L4 conserva composição.
**Decisão:** ADR-0011.

## Observável e preservação

Os parsers específicos e o dispatcher devem retornar o mesmo parsing para bytes/ref iguais sem consultar relógio. A API ParseManifestAt aceita instante explícito, UTC ou outro fuso; zero permanece desconhecido. As assinaturas existentes continuam disponíveis, retornando ConsultadoEm zero. Coletores GitHub/local preenchem a hora em L3; a reanálise do funil preserva o instante da extração, inclusive desconhecido.

Congelar antes da mudança a projeção de um corpus autoral que cobre os tipos de manifestos, removendo somente ConsultadoEm. Preservar esse digest depois. Verificar determinismo com metadado explícito, sucesso/erros, hora conhecida/desconhecida e manutenção de data nos adaptadores e na reanálise.

Atualizar o formato languages.go de crystalline.toml para a tabela aceita pelo lint instalado. Essa correção não desabilita regras, níveis ou avisos. Não promover aprovação do gate enquanto existirem diagnósticos.

## Fora de escopo

Relógio de telemetria do runtime/funil, lab_extractor, deadlines, persistência da CLI, refatoração do runtime e migração dos testes de conectores em L4. Preservar alterações locais 0126/0127, artefatos de terceiros e hashes históricos.

## Verificação

Caracterização antes/depois, testes de parsing/coletores/reatribuição do instante no funil, lint padrão, suíte raiz, vet e build. Evidência em lab/bancada/arquitetura/0128; gates selados têm módulo próprio e não estão incluídos na suíte raiz.

## Ordem estável dos objetos de dependências

A revisão encontrou iteração não canônica de mapas em package.json, composer.json, box.json, haxelib.json e dub.json. Ordenar as chaves dentro de cada grupo, preservando ordem dos grupos e nomes/versões/tipos/exclusões. Testar múltiplas dependências e grupos repetidamente. A ordem dessas listas passa a ser definida; a caracterização do corpus temporal permanece intacta. Artefatos antigos não são reordenados ou re-hasheados.
