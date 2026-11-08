# 0145 — Reorganizar fronteiras de investigação sob Tekt

**Data:** 2026-10-29. **Contrato ex-ante.** Reorganização após0142–0144 sob ADR-0008/A0. Inspeção inicial não encontrou imports concretos entre camadas proibidas; os conflitos observados são ownership interno da expansão e interpretação de argumento em L4. Não promover esta revisão a prova universal da semântica do programa.

## Recorte e responsáveis

- L1/admissão: insight_neighborhood_plan.go recebe relatório/memória/limites, valida/normaliza e produz admissão interna com request e origens congeladas já verificadas. ValidateInsightNeighborhood mantém sua assinatura como projeção pública. A fachada usa essa mesma admissão, sem ignorar um segundo erro de interpretação das fontes.
- L1/fontes: insight_neighborhood_sources.go verifica bundles, recibos e fingerprints; não possui foco, grafo ou chamadas remotas.
- L1/fronteira: insight_neighborhood_frontier.go deriva referências admissíveis, travessia e pendências; aliases permanecem em insight_neighborhood_selection.go.
- L1/resultado: insight_neighborhood_result.go inicializa/projeta relatório/tarefas e decide estados operacionais; não lê arquivos ou provedores. Fachada insight_neighborhood.go apenas coordena admissão, deadline, plano, memória, fronteira, execução e resultado.
- L1/memória: bootstrap, estado e registro de testemunhos pertencem ao container insight_exploration_memory.go. Ordenação compartilhada pertence a insight_followup_selection.go; projeção de cartões a insight_followup_context.go; insight_followups.go conserva execução limitada por portas.
- L2/entrada: web_input.go interpreta a presença literal histórica de --no-token. L4 só consulta essa projeção e compõe ambiente/cliente. Preservar inclusive comportamento anterior das variantes com =; o parser completo continua no comando existente.
- L3 mantém os efeitos e L4 mantém os construtores existentes. Sem novo subsistema, runtime, rede real, movimento cosmético da árvore ou expansão de produto.

## Invariantes

Preservar APIs públicas, JSON, schemas, IDs, hashes, ordenação, budgets, erros, cancelamento, semântica de foco/aliases e estado de evidência. Nenhuma alteração nos artefatos históricos ou em corpos dos testes consumidores existentes; realocar testes por responsável com seus nomes/asserts conservados. Prompts espelhados atualizados antes do código; regras canônicas e arquivos de orientação não mudam.

## Verificação

Auditoria global de imports/efeitos e regime L0 por gates existentes; novo controle de ownership no Lab, inclusive controles negativos. Testes existentes e testes de isolamento da admissão, prioridade e probe de transporte. Comparar expansão offline do binário anterior e novo em relatórios reais congelados (foco, geral, feedback, aliases) byte a byte, sem ler destinos. Reproduzir descoberta/memória de fixture e replays históricos. Pilha completa: lint, testes, vet e build. Evidência em lab/bancada/arquitetura/0145; L0 somente síntese e referências.
