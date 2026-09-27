# Prompt Spec L0: CLI de Expansão e Avaliação de Ideias

## Objetivo

Expor `idea expand` e `idea evaluate` como comandos JSON offline, recebendo a
estratégia lexical por injeção.

## Restrições

- Camada L2, pacote `cli`.
- L2 não importa o adaptador lexical de L3.
- `expand` aceita `--input`, `--strategy` e `--output` opcional.
- `evaluate` aceita `--corpus`, `--strategy` e `--output` opcional.
- Somente JSON em stdout; erros devem ser acionáveis.
- Estratégias suportadas: `literal` e `lexical`.

## Critérios de Verificação

- Os comandos funcionam sem infraestrutura remota.
- `--output` grava arquivo e mantém stdout vazio.
- Corpus, estratégia e JSON inválidos são rejeitados.

Segurança: corpus JSON usa readIdeaJSONInput antes de Unmarshal, conforme security_json_input.md; saídas válidas preservadas.
