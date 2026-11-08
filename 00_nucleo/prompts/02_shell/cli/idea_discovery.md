# Prompt Spec L0: CLI de Descoberta de Ideias

## Objetivo

Expor captura remota por `idea discover` e avaliação offline por
`idea discovery evaluate`.

## Restrições

- L2 recebe portas prontas por injeção e não importa L3.
- Captura exige arquivo de saída e valida limites antes da rede.
- Snapshot e relatório são JSON; logs não contaminam stdout.
- Replay não constrói provedor remoto.
## Extensão operacional do Passo 0013

Aceite `--timeout` e `--attempt-timeout`, valide-os antes da rede, trate sinais do
sistema como cancelamento estruturado, mostre progresso em stderr e grave
atomicamente snapshots completos, parciais ou falhos antes de retornar o erro.
Registre `NO_SEMANTIC_EXPANSION` quando a estratégia lexical não propuser nada.

Segurança: planos/corpus de descoberta usam readIdeaJSONInput antes de Unmarshal, conforme security_json_input.md; validar antes de acessar provedores e preservar saídas válidas.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
