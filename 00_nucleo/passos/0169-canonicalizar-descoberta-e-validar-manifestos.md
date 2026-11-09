# 0169 — Canonicalizar descoberta e validar manifestos

## Intenção anterior à materialização

Corrigir os quatro achados de 2abef74 em L1: idea_discovery possui união/origens, regional_discovery admite catálogo, recommendation_dossier_artifacts possui pré-imagens/estágios, recommendation_dossier admite e publica. Preservar schemas, assinaturas públicas, perfil sem locales, respostas sem variantes e todos os registros/hashes anteriores.

## Aceitação

Permutar variantes e matches sem mudar snapshot nem perder origens; testar desempates completos e inputs intactos. Falhas de hash em entradas e saídas impedem COMPLETE e publicação. Locales equivalentes possuem identidade/retomada equivalentes pela API; inválidos falham antes de efeitos. Catálogo com chave não canônica é erro, enquanto tradução ausente válida é Unknown. Controles congelam hashes anteriores do caminho válido.

Executar make verify com V0–V27 ativas e comparar o acervo histórico por SHA256. Evidências de execução ficam em /tmp, sem promover gates selados.

## Materialização e regressões

DSC-01/DOS-01/DOS-02/REG-01 foram acrescentados ao [contrato de integridade](../contratos/integridade-de-entradas-e-replay-v1.md). Os quatro prompts proprietários e os dois prompts de regressão explicitam admissão, observáveis e limites. União escolhe representante inteiro por ordem total e conserva todos os matches; cópias incluem Attempts antes do enriquecimento regional. A montagem de estágios mantém assinatura pública e representa falha de hash como FAILED sem saídas; a fachada persiste a falha e impede publicação. Locales são admitidos antes da identidade e efeitos; chave não canônica é rejeitada sem lacuna artificial.

As regressões reproduziram os quatro achados antes do patch (/tmp/bateia-0169-before.log), e passaram depois (/tmp/bateia-0169-focused.log). Cobrem permutações de variantes/matches/evidências, desempates completos, origens idênticas, slices retidos pela porta, nil/vazio, catálogo inconsistente e ausência legítima. Cada artefato obrigatório foi exercitado com erro de hash e hash vazio. A API completa do dossiê confirmou equivalência de identidade e retomada para ordem, repetição e all, além de rejeição de locale inválido antes de Prepare. Os testes congelam hashes anteriores do snapshot único, pré-imagem sem locales, RunID e estágios válidos; não se alteraram as assertivas anteriores. A suíte L1 completa passou sem cache de resultados em /tmp/bateia-0169-core.log.

A verificação canônica usa /tmp/bateia-0169-verification e a comparação SHA256 do acervo usa /tmp/bateia-0169-historical.json. Não há mudança de schema, política de decisão, registros, selos ou hashes históricos; os módulos/gates selados mantêm seus commits e recibos próprios.

## Resultado da verificação

make verify final passou em lint, arquitetura, go test ./..., go vet ./..., build de ./04_wiring e diff-check. O alcance efetivo confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true, 1.192 fontes analisadas e zero erros de leitura/sintaxe/infraestrutura; passaram os 11 controles de parser e os dois controles V3 ligado/desligado. As três exclusões históricas de oráculos permaneceram configuradas.

A comparação SHA256 confirmou 4.883 arquivos do acervo Lab/passos/laudos intactos (o índice ativo de passos é excluído da comparação por ser atualizado nesta entrega). Os probes da revisão reaplicados em /tmp/bateia-0169-after-{discovery,dossier,catalog}.log confirmaram snapshots estáveis com origens completas, estágios FAILED para hashes ausentes, identidade equivalente dos locales e erro de catálogo inconsistente. Nenhum registro ou hash histórico foi recalculado. Evidências completas: /tmp/bateia-0169-verification/verification.json e /tmp/bateia-0169-history-check.json; são diagnósticos locais, não novos recibos selados.
