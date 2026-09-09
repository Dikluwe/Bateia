# Passo 0060 — Tradução local explícita de caminho

## Contrato antes do candidato
Na interpretação de consulta local, adicionar tradução finita caminho/caminhos → path. Somente Mode=natural, sem marcadores de negação existentes; somente palavras naturais originais elegíveis para expansão, nunca símbolos/backticks ou termos derivados. Preservar Original, Symbols, Constraints, expansões existentes e todos os termos anteriores.

Cada palavra original distinta reconhecida gera uma ImplementationQueryExpansion com Source em minúsculas (caminho ou caminhos), Terms=[path], RuleID=0060-v1:path:translation, Kind=translation, Status=hypothesis. Terms globais deduplicados e ordenados. Não inventar tradução parar→Remove, nomes de funções, capacidades ou mecanismos. Não alterar glossário usado para indexação de código; a tradução é da consulta, com provenance explícita. Não alterar busca pública, planejamento0055/0058, pesos, orçamento, seleção de arquivos ou afirmações comportamentais.

EvidenceVersion=evidence-v22-path-query; checkpoint implementation_search=0060 local/live; discovery_query inalterado. Novo termo não autoriza satisfação de requisitos.

## Regime e escopo
Tekt A/B sem atestação física. A lê contrato e produção; modifica somente query_interpretation.go, evidence.go, idea_recommend.go e seus prompts antes do código. B lê contrato/tipos/testes legados, sem produção candidata/corpus/relatórios, cria controles independentes de tradução, provenance, singular/plural, determinismo, símbolos, mixed, negação, não recursão e recuperação em fixtures neutras. Congelar testes antes da execução. Integrador conhece alvos e executa comparação no corpus0053 congelado, separado por canais, sem alterar métricas. Trata-se de calibração, não holdout. Registrar qualquer falha e revisão.

## Critério de aceitação
Tradução observável em consulta natural elegível, sem expansão nos modos excluídos. Recuperar implementação em fixture cujo corpo usa path, preservando evidência como lexical e comportamento não verificado. Suíte Go completa, lint, build e replay do corpus. Relatar também candidatos adicionais irrelevantes e limites da amostra: uma palavra traduzida não resolve entendimento semântico.
