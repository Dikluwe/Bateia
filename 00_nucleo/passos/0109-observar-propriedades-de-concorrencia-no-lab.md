# 0109 Observar propriedades de concorrência no Lab

**Data**: 2026-10-06  
**Estado**: IMPLEMENTADO E VERIFICADO NO RECORTE SINTÁTICO  
**Base**: pesquisas 0107/0108

Materializar observador puro AST, CLI isolada e corpus exploratório em
lab/bancada/concorrencia/0109, sem dependência do produto sobre o Lab.
Autoridade somente sintática: fatos por propriedade/localizador/condições,
com lacunas de tipos, configuração e fluxo. Configuração e identidade da
fonte são declaradas; hashes garantem vínculo, não autenticidade externa.

Não classificar presença de buffer como trabalho total limitado nem Done como
conclusão de cancelamento. Fazer testes de contraste, renomeação, alteração
significativa, parsing incompleto, integridade, escopos e publicação.
Corpus autoral não substitui referência independente. Especificar expectativas
antes da execução; medir conformance sintática sem taxa de qualidade semântica.

Exercitar também bytes reais congelados do conc/0106, sem executar candidato.
Conservar evidências e verificação no Lab e síntese ex-post em laudo L0.
Executar lint, testes focados/race, suíte Go completa, vet e builds do produto
e da CLI experimental. Não alterar catálogo ou resultados históricos.

## Execução

O [laudo 0109](../laudos/0109-observador-de-concorrencia-no-lab.md) registra
implementação, corpus, fonte real e limites. Qualidade semântica independente
e promoção para o produto permanecem pendentes.
