# Prompt Spec L0: Motor de Ontologia e Similarity Matcher

## Objetivo
Implementar o motor puro em L1 para extração da ontologia de repositórios extraídos, cálculo do Score de Similaridade de Ideias em 4 vetores (Intenção/Semântica, Stack/Ecossistema, Topologia Arquitetural, Cobertura de Capacidades) e mineração de recomendações de melhorias baseadas nas lacunas da ontologia alvo.

## Restrições
- Camada: L1 (Core Puro)
- ZERO I/O (sem HTTP, sem IO de disco, sem banco de dados).
- 100% determinístico e testável unitariamente.
