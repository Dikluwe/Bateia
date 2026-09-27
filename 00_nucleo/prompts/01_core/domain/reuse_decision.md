# Reuse decision domain

Defina contratos JSON versionados para decidir entre FORK, EXTEND, COMPOSE,
INTEGRATE, REFERENCE, BUILD e INCONCLUSIVE. Preserve cobertura funcional,
licença, manutenção, operação, esforço e conflitos como dimensões separadas.
Toda decisão deve expor matriz, política e hash, razões, bloqueios e evidências.

## Evolução 0043

# Passo 0043 — Plano verificável de investigação

Origem: 00_nucleo/passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md.
L1 puro: transformar IdeaSpec + DiscoverySnapshot + avaliações em plano determinístico, sem IO/relógio. Requisitos recebem IDs por texto normalizado e tipo, ordenados; separar problema, capacidades, restrições, não objetivos e preferências. Covers indexados de expansões referem-se à ideia original do snapshot, sem reinterpretar índice após ordenação. Query alterada tem alteração explícita sem inventar filtros removidos; ausente permanece desconhecida. Propostas de busca não são testemunhos.

BuildInvestigationPlan retorna requisitos, tentativas observadas, ações por candidato/requisito, classes de evidência, método, lacuna, critério de aceitação e encerramento delimitado. Sanitizar avaliações pelo contrato comum antes de atribuir status; sem candidatos produzir NO_CANDIDATES_IN_SCOPE somente sem falhas, nunca inexistência global. Queries e fontes opacas não geram prova. Plano sem snapshot informa DISCOVERY_NOT_RECORDED. Resultado ordenado independe da ordem dos conjuntos de entrada. Sem ensaio de integração sempre preservar condição pendente de validação.

Anexar plano à decisão canônica e à apresentação CLI/MCP. Projetar texto público ao construir o plano com os filtros existentes de privacidade; preservar o gateway certificado. Telemetria identificada como da ferramenta. Compatibilidade aditiva com schema aninhado próprio; não promover resultados históricos.

## Passo 0088

Adicionar CompositionInput opcional à requisição e CompositionReport ao relatório, preservando schema de decisão e o modo simples. IDs e lacunas de composição devem ser rastreáveis.
