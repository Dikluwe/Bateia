# Prompt: 01_core/usecase/idea_radar_plan

Implementar
`PlanIdeaRadarProbes(domain.IdeaRadarInput) domain.IdeaRadarPlanReport` exatamente segundo
o contrato 0067 congelado. Chamar `AnalyzeIdeaRadar` como autoridade de classificação;
não editar, replicar ou relaxar o 0066.

Gerar somente a fronteira do bloqueador ativo. Combinar `ADD_PERSPECTIVE` com origem
distinta, selecionar o menor subconjunto determinístico de observações de origem
desconhecida, completar apenas déficits residuais de origem e emitir somente
`BIND_EVIDENCE` em convergentes. Ordenar e identificar probes como especificado.

Não executar probes, prometer convergência, fabricar `EvidenceSource`, chamar 0065,
inferir texto ou usar ground truth/baseline produtivamente. Preservar input, arrays não
nulos, limites e operações exatas. Sem I/O, relógio, aleatoriedade ou dependências
externas. O futuro arquivo é `01_core/usecase/idea_radar_plan.go`, com cabeçalho de
linhagem deste prompt e camada L1. O implementador não edita entradas protegidas ou testes.

