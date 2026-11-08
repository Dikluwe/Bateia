# 0126 — Integrar Gitee e AtomGit à descoberta

**Data:** 2026-10-18
**Capacidade:** descoberta de repositórios, A1-INV sob Tekt.
**Responsáveis:** contratos e ID em L0/L1; HTTP em L3; ambiente e composição em L4; projeção MCP existente em L2.

## Intenção e observável

Gitee já possui adaptador, mas está ausente do conjunto padrão. AtomGit não possui adaptador. Ambos devem aparecer em bateia_list_connectors e no schema da descoberta do executável, respeitando seleção máxima de dez. Busca pública retorna candidatos com identidade da plataforma e proveniência; ausência de token, falha remota e resposta vazia permanecem distintas.

## Limites e incertezas

A API AtomGit documenta token obrigatório e flags de visibilidade em dois formatos. Validar ambos em fixtures autorais. Corrigir branch usada como commit e capacidades excessivas do adaptador compartilhado. Não adicionar clone, inspeção de código, tradução, promoção de evidência, regras de licença ou deduplicação heurística. Não mudar contratos históricos ou gates selados.

## Verificação prevista

Testes L3 de protocolo e falhas; composição real offline do MCP de mercado e síntese; crystalline-lint, go test ./..., go vet ./... e build de ./04_wiring em destino temporário. Consultas públicas reais, quando viáveis, são verificação separada e não substituem testes determinísticos. Síntese ex-post no laudo 0126.

## Dependência encontrada antes da materialização L1

O executor pula tentativas MARKET_NOT_FILTERED, inclusive com orçamento reservado. Isso impede o uso dos catálogos globais. A menor correção permite sua execução quando PlannedCost é positivo e mantém uma lacuna geográfica independente do sucesso operacional. A projeção não atribui ao candidato o país solicitado sem filtro aplicado. Ownership continua em L1, sem novo enum ou alteração de hash da requisição.
