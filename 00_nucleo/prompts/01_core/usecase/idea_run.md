# Idea run orchestrator

Orquestre os casos de uso existentes, sem duplicar regras, nos modos live, local
e replay. Valide configuração antes de chamar dependências. Replay jamais chama
provedores. Falhas parciais são preservadas e INCONCLUSIVE nunca é promovido a
COMPLETE. Ordene candidatos locais e não execute código de terceiros.
Entre candidatos e etapas externas, respeite cancelamento já observado e não
inicie novas chamadas a provedores.

## Evolução 0043

# Passo 0043 — Plano verificável de investigação

Origem: 00_nucleo/passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md.
L1 puro: transformar IdeaSpec + DiscoverySnapshot + avaliações em plano determinístico, sem IO/relógio. Requisitos recebem IDs por texto normalizado e tipo, ordenados; separar problema, capacidades, restrições, não objetivos e preferências. Covers indexados de expansões referem-se à ideia original do snapshot, sem reinterpretar índice após ordenação. Query alterada tem alteração explícita sem inventar filtros removidos; ausente permanece desconhecida. Propostas de busca não são testemunhos.

BuildInvestigationPlan retorna requisitos, tentativas observadas, ações por candidato/requisito, classes de evidência, método, lacuna, critério de aceitação e encerramento delimitado. Sanitizar avaliações pelo contrato comum antes de atribuir status; sem candidatos produzir NO_CANDIDATES_IN_SCOPE somente sem falhas, nunca inexistência global. Queries e fontes opacas não geram prova. Plano sem snapshot informa DISCOVERY_NOT_RECORDED. Resultado ordenado independe da ordem dos conjuntos de entrada. Sem ensaio de integração sempre preservar condição pendente de validação.

Anexar plano à decisão canônica e à apresentação CLI/MCP. Projetar texto público ao construir o plano com os filtros existentes de privacidade; preservar o gateway certificado. Telemetria identificada como da ferramenta. Compatibilidade aditiva com schema aninhado próprio; não promover resultados históricos.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
