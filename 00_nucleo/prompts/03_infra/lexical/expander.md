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
