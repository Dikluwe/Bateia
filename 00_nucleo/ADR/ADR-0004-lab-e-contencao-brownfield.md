# ADR-0004: Lab canônico e contenção brownfield da bancada

**Status**: ACEITO
**Data**: 2026-09-18
**Baseline**: `d02eafefddef0c8fadf6b6751dfed121105963ac`
**Relacionado**: ADR-0002, ADR-0003, passo 0076, Manifesto Tekt v1.4

## Contexto

O Bateia preservou a gravidade principal entre L1, L2 e L3, mas passou a armazenar
executores de benchmark em L0 e gates, mutantes e experimentos em L4. O linteiro
aceitava esse estado por exclusões e exceções explícitas. A topologia compilava, porém
as pastas deixaram de expressar o regime real dos artefatos.

Tekt distingue o regime do sistema (`00_nucleo` a `04_wiring`) do regime de bancada
(`lab/`). L0 contém causa e evidência declarativa, nunca executores. L4 contém somente
entrada, composição, injeção e lifecycle. Experimentos, mutantes, gates e corpora
executáveis pertencem ao Lab e nunca podem ser importados pelo sistema produtivo.

## Decisão

1. Adotar `lab/` como Workbench canônico do Bateia.
2. Organizar o Lab em `arena/` para experimentos e `bancada/` para gates,
   benchmarks e medições reproduzíveis.
3. Retirar de L4 todos os diretórios cujo papel seja gate, mutante, fixture ou
   experimento histórico.
4. Retirar de L0 toda bancada, inclusive benchmarks, fixtures e resultados passivos.
   L0 pode descrevê-los e referenciá-los; o acervo pertence a `lab/bancada/history`.
5. Preservar nomes cronológicos dentro do Lab: ali eles identificam a execução
   histórica e não definem a arquitetura do produto.
6. Não alterar neste marco domínio, CLI, MCP, persistência, formatos públicos,
   mensagens, defaults nem política de `Unknown`.
7. Atualizar o linteiro para reconhecer `lab/` e remover exceções de wiring que só
   existiam porque a bancada estava em L4.

## Consequências

- L0 volta a ser exclusivamente declarativo, sem armazenar o acervo medido que descreve.
- L4 deixa de ser depósito de verificação histórica.
- A suíte de produto continua funcionalmente idêntica; material histórico útil fica
  isolado no Lab.
- A decomposição por capacidade, a remoção de relógio de L1 e o enxugamento semântico
  de L4 permanecem trabalhos posteriores, cada um precedido por exegese e prompt.

## Critérios de aceite

- `crystalline-lint --fail-on warning .` sem violações;
- `go build -o /dev/shm/bateia-cleanup ./04_wiring` passa;
- `go test ./...` passa;
- nenhum fonte executável, inclusive código mascarado por sufixo `.txt`, permanece em `00_nucleo`;
- nenhum diretório de gate, mutante ou experimento permanece em `04_wiring`;
- nenhuma camada produtiva importa `bateia/lab`;
- o diff não altera comportamento produtivo.

## Política de incerteza

Arquivo de natureza ambígua, import não resolvido, teste não executado ou necessidade
de mudança semântica resulta em `Unknown` e bloqueia este marco. `Unknown` nunca é
convertido em preservação por padrão.
