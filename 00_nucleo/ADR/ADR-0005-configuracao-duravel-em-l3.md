# ADR-0005: Configuração durável local pertence a L3

**Status**: ACEITO
**Data**: 2026-09-18
**Relacionado**: ADR-0004, passo 0077, laudo 0034

## Contexto

`04_wiring/durable_config.go` lê ambiente, cria e audita diretórios, inspeciona
symlinks/hardlinks, gera segredo criptográfico e usa syscalls POSIX. Esses efeitos
são infraestrutura; L4 deve apenas solicitar a configuração e injetar seus valores.

O prompt histórico atribuiu essa responsabilidade a L4 antes da explicitação da
fronteira atual. Ele permanece evidência causal do baseline, mas não governa a nova
materialização.

## Decisão

1. Criar `03_infra/durableconfig` como adaptador proprietário da configuração local.
2. Expor somente `LoadLocal() (Config, error)` e o valor imutável necessário à
   composição: raiz, segredo, limites, workers e capacidade.
3. Manter em L3 os helpers de filesystem, segurança e geração do segredo.
4. L4 chama o adaptador nos mesmos dois pontos e preserva ordem, mensagens e
   lifecycle existentes.
5. Preservar os prompts L4 anteriores, byte a byte, no histórico do Lab; criar
   prompts L3 novos para a obrigação refinada.

## Consequências

L4 perde imports e operações de baixo nível sem alterar o comportamento público. O
novo pacote permanece Linux/POSIX nos mesmos limites do baseline. Injeção de
filesystem ou RNG, parsing/dispatch de `main.go` e refatoração de `runstore` ficam
fora desta decisão.
