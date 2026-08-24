# Prompt Spec L0: Comando Offline `idea plan`

## Objetivo

Expor em L2 a materialização inspecionável de `IdeaSpec → ExpansionPlan` por JSON,
sem inicializar qualquer infraestrutura remota ou cache.

## Restrições

- Camada L2, pacote `cli`.
- Entrada por `--input` em JSON.
- Saída JSON em stdout ou `--output`.
- Stdout contém somente JSON quando não houver arquivo de saída.
- Erros de leitura, JSON e validação são acionáveis.
- O wiring deve despachar o comando antes de `.env`, Bleve, GitHub e Sourcegraph.

## Critérios de Verificação

- Golden test fixa a forma pública.
- Funciona sem token, rede e cache gravável.
- JSON inválido e ideia vazia retornam erro.
- Com `--output`, grava o plano esperado.
