# Prompt Spec L0: Expansor Lexical Local

## Objetivo

Implementar a porta de expansão semântica com regras locais, versionadas,
determinísticas e auditáveis.

## Restrições

- Camada L3, pacote `lexical`.
- Vocabulário JSON validado, com IDs únicos e expansões não vazias.
- Correspondência de frase normalizada em limites de palavras.
- Sem rede, relógio, aleatoriedade ou regras ocultas.
- Limite de expansões por query de origem.
- Cada proposta informa regra, estratégia e expansão de origem.

## Critérios de Verificação

- Vocabulário inválido é rejeitado com erro acionável.
- Substrings parciais não ativam regras.
- Ordem, IDs de regra e resultados são estáveis.

## Passo 0068 — catálogo técnico finito

Aplicar `00_nucleo/passos/0068-compilar-capacidade-em-consultas-tecnicas.md`. Preservar o
schema legado e aceitar opcionalmente `strategy` e dois `discriminators` por regra. Sem
strategy, usar `lexical`; regra técnica usa somente `technical-lexicon-v1`. Materializar
a tradução inglesa canônica e duas queries curtas para a frase portuguesa integral,
conservando facet, covers, RuleID e SourceID. Frase parcial ou entrada fora do catálogo
não ativa tradução. Validar strings, limites e duplicatas; não inventar tradução geral,
nome de repositório ou prova de capacidade.
