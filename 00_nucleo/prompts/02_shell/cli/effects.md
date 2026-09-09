# Prompt: 02_shell/cli/effects

RunEffectsCommand(args,out,verifier) offline: verify --source --symbol --output opcional --timeout default5s positivo até30s --seed uint64 default1. Validar todos argumentos antes de executar. Help sem execução, limites explícitos. JSON para avaliação concluída inclusive unknown; erros CLI/I/O retornam erro. Dependência apenas interface L1.

Contrato autoritativo: `00_nucleo/passos/0062-verificar-efeito-executado.md`. Implementação A; testes e veredito independentes pelo protocolo segregado.

Passo 0063: aceitar --expect-sha256 HEX opcional, exatamente 64 dígitos hexadecimais ASCII; normalizar para minúsculas e transmitir ExpectedSourceSHA256. Flag explicitamente vazia ou malformada é erro antes de chamar verifier; ausência mantém compatibilidade. Help explica hash de arquivo completo, mismatch unknown sem execução e ausência de vínculo automático com busca/repositório/commit ou dependências externas. Contrato adicional autoritativo: `00_nucleo/passos/0063-vincular-fonte-ao-efeito.md`.
