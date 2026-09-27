# Passo 0043 — Plano verificável de investigação

Origem: 00_nucleo/passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md.
L1 puro: transformar IdeaSpec + DiscoverySnapshot + avaliações em plano determinístico, sem IO/relógio. Requisitos recebem IDs por texto normalizado e tipo, ordenados; separar problema, capacidades, qualidades, condições operacionais, restrições, não objetivos e preferências. Qualidades e condições operacionais não avaliadas geram ações pendentes, não prova. Covers indexados de expansões referem-se à ideia original do snapshot, sem reinterpretar índice após ordenação. Query alterada tem alteração explícita sem inventar filtros removidos; ausente permanece desconhecida. Propostas de busca não são testemunhos.

BuildInvestigationPlan retorna requisitos, tentativas observadas, ações por candidato/requisito, classes de evidência, método, lacuna, critério de aceitação e encerramento delimitado. Sanitizar avaliações pelo contrato comum antes de atribuir status; sem candidatos produzir NO_CANDIDATES_IN_SCOPE somente sem falhas, nunca inexistência global. Queries e fontes opacas não geram prova. Plano sem snapshot e sem candidatos informa NO_CANDIDATES_IN_SCOPE e explicita no escopo que a descoberta não foi registrada; não equivale a busca executada. Resultado ordenado independe da ordem dos conjuntos de entrada. Sem ensaio de integração sempre preservar condição pendente de validação.

Anexar plano à decisão canônica e à apresentação CLI/MCP. Projetar texto público ao construir o plano com os filtros existentes de privacidade; preservar o gateway certificado. Telemetria identificada como da ferramenta. Compatibilidade aditiva com schema aninhado próprio; não promover resultados históricos.

Fontes inadmissíveis não aparecem como evidência de inferência: preservar NOT_VERIFIED e ação para obter testemunho válido. Duplicatas conflitantes da mesma identidade não permitem last-write-wins.

Covers também pode conter texto literal exato do requisito: vincular somente igualdade normalizada, nunca similaridade ou palavras parciais.

Legibilidade 0043: cada ação contém pergunta e lacuna explícita. Sem detector para o requisito, pedir especificação/verificação delimitada, sem fingir que mais coleta resolve ausência de analisador. Markdown mostra texto do requisito além do ID e fontes pertinentes quando presentes.
