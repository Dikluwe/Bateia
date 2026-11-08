# CLI idea decide

Exponha `bateia idea decide --idea --assessment --policy [--composition] [--output]` como boca de
inspeção totalmente offline. Leia somente artefatos locais, invoque o decisor
determinístico e emita JSON versionado; parâmetros ausentes e JSON inválido devem
falhar antes da decisão.

O argumento composition é opcional. Sem ele, a pesquisa simples permanece
intacta, mas combinações funcionais são hipóteses e não uma decisão COMPOSE.
Conexões, exclusões, dependências e ensaio do conjunto exato são declarativos.

## Passo 0098 — Sidecar opt-in

--trace-output gera grafo de dependências registradas, sem alterar bytes do
report ou recalcular/verificar a decisão. Validar novo destino contra todas as
entradas e saída principal (path, inode, aliases ancestrais). No caminho opt-in,
validar também saída principal contra entradas; preservar CLI sem flag. Construir
grafo antes das saídas; writer trace atômico individual não garante transação
entre report e sidecar. As claims dependem do report original para interpretação.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
