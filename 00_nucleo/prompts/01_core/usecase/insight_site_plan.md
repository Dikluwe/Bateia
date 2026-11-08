# Plano de fontes e orçamento global

**Camada:** L1. **Responsável:** opções e admissão do plano multifonte. **Materialização:** [insight_site_plan.go](../../../../01_core/usecase/insight_site_plan.go).

## Interface e responsabilidade

ValidateInsightSiteOptions normaliza o pedido sem consultar conectores. PlanInsightSites recebe pedido e bindings disponíveis, valida capacidades e produz plano selado com fontes canônicas. O plano é consumido por discover e pela fachada de expand. L4 injeta as portas, sem decidir elegibilidade de memória ou limites de negócio.

## Comportamento vigente

Delegar a normalização lexical e seus limites ao [plano de descoberta](insight_discovery_plan.md). FollowLinks ausente recebe2; copiar o ponteiro fornecido e admitir0..4. Validar memória por cópia profunda; pedir ao container sua guarda de avanço e verificar que o próximo snapshot selado cabe em1 MiB antes de admitir. Histórico999999 é inspecionável, mas não pode iniciar rodada;999998 admite avanço para999999, inclusive com follow-links0. Não chamar capabilities, busca ou leitura antes de rejeitar opções inválidas.

Admitir de1 a12 bindings com conectores não nulos, IDs únicos e capacidades canônicas de até4096 bytes. Reader legado e FollowReader são exclusivos de GitHub; RepositoryReader opcional só pode atender Gitee/AtomGit/Codeberg/GitLab e não coexiste com Reader. Leitores ausentes não impedem catálogo. Não confundir namespace GitLab com owner/repo de outro sítio.

Selecionar de1 a12 fontes pelo ID ou aliases conhecidos: github, sourcegraph, gitee, atomgit, codeberg, gitlab, bitbucket, joss, zbmath, itunes, playstore e fdroid. `all` só é válido sozinho e seleciona os bindings disponíveis. Rejeitar desconhecidos, repetidos após resolução de alias ou IDs com espaços externos; ordenar IDs canônicos e capacidades para obter plano determinístico. Não impor autenticação aqui.

MaxCalls0 recebe20; limite1..60. Locale vazio recebe en-US; demais valores exigem BCP47 válido. Selar site_plan_ pelo algoritmo histórico. Registrar a política compartilhada de [agenda](insight_site_schedule.md); este proprietário não executa essa agenda nem possui seleção de cartões, leitura ou derivação de leads.

## Limites e dependências

Sem I/O ou fallback concreto. O [contrato de exploração](../../../contratos/exploracao-de-vizinhancas-v1.md#admissão) e o [container de memória](insight_exploration_memory.md) governam a retomada. Orçamento herdado de catálogo não autoriza buscas em expand. Preservar schemas, IDs e o ramo GitHub legado sem --sources; coleção e runtime têm outros proprietários.

## Aceitação

Consumir [o prompt de testes](insight_site_plan_test.md), ADM-03 e MEM-04 da [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Exigir equivalência do plano ao permutar bindings/aliases e rejeição antes dos efeitos quando a geração não puder avançar.

## Origem e preservação

Consolidação posterior ao código das regras0138–0142 no [passo0147](../../../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md); a guarda de avanço é intenção ex-ante da correção desse passo. O [texto anterior](../../../../lab/bancada/prompts/0147/before/00_nucleo__prompts__01_core__usecase__insight_site_plan.md.txt) preserva a cronologia e suas precedências. Não reclassificar a consolidação como origem ex-ante do comportamento legado.
