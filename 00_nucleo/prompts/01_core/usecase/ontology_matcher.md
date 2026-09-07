# Prompt Spec L0: Motor de Ontologia e Similarity Matcher

## Objetivo
Implementar o motor puro em L1 para extração da ontologia de repositórios extraídos, cálculo do Score de Similaridade de Ideias em 4 vetores (Intenção/Semântica, Stack/Ecossistema, Topologia Arquitetural, Cobertura de Capacidades) e mineração de recomendações de melhorias baseadas nas lacunas da ontologia alvo.

## Restrições
- Camada: L1 (Core Puro)
- ZERO I/O (sem HTTP, sem IO de disco, sem banco de dados).
- 100% determinístico e testável unitariamente.

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.
