# Passo 0077 — Configuração durável em L3

## Objetivo

Retirar de `04_wiring` a leitura de ambiente, filesystem, geração de segredo e
validação de diretórios. Essas operações são infraestrutura e pertencem a
`03_infra/durableconfig`.

## Contrato

`durableconfig.LoadLocal` resolve a raiz de execução, garante diretórios privados,
recusa symlinks e arquivos inseguros, cria ou lê a chave local e valida limites e
capacidade. L4 apenas recebe a configuração pronta e a injeta.

## Aceite

- defaults e variáveis de ambiente continuam compatíveis;
- segredo e diretórios mantêm as restrições de segurança;
- chamadas concorrentes não corrompem a criação inicial;
- L4 não contém implementação de configuração durável.
