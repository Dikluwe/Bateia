# 0129 — Concluir as fronteiras arquiteturais declaradas

**Data:** 2026-10-28
**Contrato:** ADR-0012, sob Tekt.

## Recorte e aceite

Concluir as pendências inventariadas: efeito temporal restante em L1, deadlines, efeitos diretos de arquivos/entropia/ambiente na CLI, configuração acumulada no runtime e classificação dos avisos de L4. Preservar a integração 0126 e correções 0127/0128 já presentes na worktree. Não reescrever evidência histórica.

Declarar portas de runtime e arquivos em L1, implementá-las em L3 e compor em L4. L1/L2 não contêm chamadas concretas ao sistema operacional para esses efeitos. L2 pode operar sobre handles abstratos e fluxos de entrada/saída fornecidos, preservando validação de transporte e apresentação. Não usar interfaces que exponham os.File/os.Root/syscall.

Manter funcionalidade do executável, segurança de entrada/saída/captura, hashes dos bytes, retomada, cancelamento, schema MCP e único lifecycle. Configuração e estado de capacidade têm containers internos próprios. Identificar limitações do lint por análise da implementação do verificador e testar a correção de forma reproduzível; não alterar níveis para silenciar avisos.

## Verificação

Inventário de efeitos/imports e linhagem; testes controlados das portas; contratos ativos de arquivos, captura, deadlines, duração e runtime; suíte raiz, vet, build e lint. Evidência bruta em lab/bancada/arquitetura/0129. Gates selados têm módulo separado e só são reproduzidos conforme seus contratos.

## Resultado

Consultar o [laudo 0129](../laudos/0129-fronteiras-arquiteturais.md) e a [bancada reproduzível](../../lab/bancada/arquitetura/0129/README.md). A entrada protegida de instruções históricas foi preservada após rejeição da revisão automática; não é alterada por este recorte.
