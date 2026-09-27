# Prompt Spec L0: Expansor Determinístico de Ideia

## Objetivo

Transformar uma `IdeaSpec` já estruturada num `ExpansionPlan` auditável, com uma
expansão para o problema e uma por capacidade, além de expansões de restrição
quando presentes.

## Restrições

- Camada L1, pacote `usecase`.
- Zero I/O, rede, filesystem, relógio ou aleatoriedade.
- Normalização e deduplicação determinísticas.
- Preservar qualidades e condições operacionais normalizadas sem convertê-las em buscas não autorizadas.
- IDs estáveis derivados do conteúdo normalizado.
- Cada expansão declara faceta, justificativa e campos cobertos.
- Não inventar sinônimos e não transformar `NonGoals` em buscas positivas.
- Não modificar a entrada recebida.

## Critérios de Verificação

- Mesma entrada sempre produz saída idêntica.
- Ideia sem problema e sem capacidades é rejeitada.
- Duplicatas por caixa e espaço são eliminadas.
- Ordem da saída é estável.
- Prova de mordida preserva alvo e consultas em estruturas distintas.
