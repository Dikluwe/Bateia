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

## Passo 0090

Roteie `idea landscape analyze` para a casca offline, sem alterar `idea plan` ou a busca simples.

Segurança: IdeaSpec usa leitura limitada seguida de validateIdeaJSONStructure antes de Unmarshal; preservar mensagens existentes de JSON inválido e saídas válidas, conforme security_json_input.md.


## Passo 0098

Roteie idea trace build/impact para a casca offline própria, sem alterar plan
ou a busca simples. Erros de subcomandos e uso incluem trace. Writer trace não
reutiliza o writer legado não atômico; o contrato está em idea_trace.md.


## Passo 0105

Roteie idea insights analyze para a casca offline específica, preservando demais comandos e mensagens legadas.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
