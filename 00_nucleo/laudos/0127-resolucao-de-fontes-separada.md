# Laudo 0127 — Resolução de fontes separada do protocolo de síntese

**Data:** 2026-10-28
**Contrato:** [passo 0127](../passos/0127-separar-resolucao-de-fontes-da-sintese.md)
**Base:** HEAD 1104cad com alterações locais anteriores do passo 0126 e outros artefatos de investigação preservados.

## Resultado arquitetural

A resolução de fontes da síntese passa a pertencer a `01_core/usecase/improvement_sources.go`. O componente recebe uma porta interna com apenas Get e ReadArtifact; não recebe RecommendationRunService e não dispõe de fila, transição ou publicação. Possui checagem de identidade/hash, paginação, interpretação dos mecanismos/lacunas/observações e validação das referências das propostas.

O protocolo em `improvement_synthesis.go` conserva admissão, normalização, hash causal, idempotência, execução, cancelamento e publicação. Chama a mesma resolução antes da fila e durante a execução. `improvement_analysis.go` continua responsável pelos resultados e estado de negócio. Os três permanecem em L1 e no pacote usecase. Não foram alterados construtores públicos, campos públicos, formatos ou regras de hipótese.

Mapa arquitetural, plano de organização e orientação correspondente em .agentrules foram atualizados. O corte limita dependências de uma responsabilidade concreta; não consiste apenas em dividir um arquivo. O serviço durável ainda configura mercado e extração de mecanismos e despacha vários tipos de run.

## Evidência de preservação

Antes da extração, testes autorais caracterizaram o resolvedor original. Foram congelados dois digests: fontes concluídas com sucesso e fonte terminal degradada. A projeção inclui estado resolvido, comparações, hipóteses/IDs, lacunas, planos, proveniência e estado de negócio. Os mesmos digests passaram após a extração.

Cenários adicionais preservam falhas das portas, tipo/lifecycle do run, índice e hash, JSON, mecanismo duplicado, referências inexistentes, fallback de hash, identidade de cada página, cursor inválido, teto de 1000 páginas e cópia dos bytes. O fake final implementa somente os dois métodos de leitura. Uma comparação textual normalizada confirma preservação dos corpos extraídos, com apenas as substituições explícitas de receiver e dependência pela porta.

Registros: [bancada 0127](../../lab/bancada/arquitetura/0127/README.md) e [baseline](../../lab/bancada/arquitetura/0127/baseline.json).

## Verificação e limites

- Caracterização antes e depois: aprovada.
- `go test ./...`: aprovado no módulo raiz, incluindo gates ativos do Lab. Assinatura Git desativada somente no processo das fixtures; cache Go e loopback autorizados.
- `go vet ./...`: aprovado.
- `go build -o /tmp/bateia-arquitetura-0127 ./04_wiring`: aprovado. O destino temporário autorizado substitui /dev/shm.
- `git diff --check`: aprovado.
- `crystalline-lint --fail-on warning .`: continua bloqueado pela incompatibilidade preexistente de languages.go booleano com o lint instalado. Não foi alterada a configuração do projeto.
- Lint com configuração temporária compatível: zero erros e 46 avisos preexistentes, sem novos diagnósticos. O aviso transitório de prompt sem materialização desapareceu após a extração. Essa execução não equivale à aprovação do gate exigido.

A caracterização cobre dois cenários de projeção e falhas declaradas; não prova equivalência universal, autenticidade de fontes externas ou utilidade das hipóteses. Gates selados possuem módulo próprio e não foram reproduzidos. Relógio em L1, persistência na CLI e configuração acumulada do runtime permanecem próximos recortes possíveis.
