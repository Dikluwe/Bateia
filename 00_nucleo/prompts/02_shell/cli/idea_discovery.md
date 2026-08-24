# Prompt Spec L0: CLI de Descoberta de Ideias

## Objetivo

Expor captura remota por `idea discover` e avaliação offline por
`idea discovery evaluate`.

## Restrições

- L2 recebe portas prontas por injeção e não importa L3.
- Captura exige arquivo de saída e valida limites antes da rede.
- Snapshot e relatório são JSON; logs não contaminam stdout.
- Replay não constrói provedor remoto.
