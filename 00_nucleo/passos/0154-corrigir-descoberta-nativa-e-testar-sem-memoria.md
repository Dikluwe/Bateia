# 0154 — Corrigir descoberta nativa e testar sem memória de fontes

## Contrato antes da implementação

Capacidade: descoberta de referências para investigação. L1 possui política, consultas e seleção; L3 realiza catálogo/resolução; L2 adapta flags. Problemas observados: escopo README causa menções incidentais; rotas reduzem contexto a token com sufixos genéricos; variedade precede correspondência; resolução antecipada custa chamadas por candidato e fontes sem pin não têm nova tentativa limitada antes da leitura.

Adicionar política context-v2 explícita no pedido, padrão da CLI, com lexical-v1 disponível para compatibilidade. Pedido L1 sem política preserva contrato histórico e hashes. Derivar cinco rotas apenas de termos da pergunta/familiaridade, usando conjunções de termos distintos quando disponíveis, sem nomes de projetos, glossário de soluções, tradução ou expansão semântica fingida. Default de resultados na política nova: cinco por rota, teto público preservado.

Na política nova, GitHub busca nome/descrição/tópicos, sem README, e adia resolução de revisão até referências escolhidas. Selecionar por quantidade de termos literais em metadados antes da diversidade; não selecionar correspondência zero, conservando os candidatos e o motivo. Aplicar também à projeção multifonte e sua memória, sem alterar regra de exame de pendências. Familiaridade e popularidade não vencem correspondência. Seleção lexical não garante relevância semântica.

Antes da leitura de referência sem pin, usar porta de resolução se disponível, exigir mesma identidade/URL e SHA imutável, respeitar deadline e conflitos observados. Sem resolver ou falha continua lacuna explícita; não relaxar integridade para obter resultado. Coleta continua de README raiz e não comprova capacidade. Campo/flags omitidos conservam JSON antigo.

## Critérios e ensaio

Testes sem rede devem provar: termos preservados nas rotas, nenhuma resposta/repositório embutido, política histórica preservada, correspondência antes de cobertura/autor/popularidade, metadados sem README, resolução só escolhidos e pin antes de árvore/conteúdo, falhas/identidades divergentes sem leitura. Testar outros domínios com nomes fictícios. Rodar make verify.

Novo ensaio: mesma pergunta de negócio de 0152/0153, dois agentes novos (Bateia e XERJ), sem histórico ou fontes prontas. Não permitir aquisições originadas na memória nem busca por nomes lembrados. Consultas por problemas/capacidades permitidas. Cada aquisição tem de indicar testemunho de descoberta ou link efetivamente lido pelo produto. Baseline e ações auditados. O runner limita URLs elegíveis a retornos de catálogo e links dos bytes lidos; não admite URL ecoada em plano/pergunta como descoberta. Mesmos limites, cobertura e rubrica de relevância/apoio/decisão; parecer anônimo após congelar resultados. Ausência não prova incapacidade. Se um braço não encontra material, relatar sem resgate por memória.

Bancada em lab/bancada/insights/descoberta-nativa-2026-11-08. Nenhum ajuste a partir dos resultados dos pesquisadores neste ensaio; melhorias posteriores exigem outro registro. Não executar fontes nem publicar.

## Resultado

Correção e gate canônico concluídos. O reteste com origem nativa obrigatória encontrou fontes úteis nos dois braços, sem aquisições de memória. Bateia retornou material mais próximo de implementação cliente/inventário nesta execução; o parecer registrou famílias operacionais compartilhadas e diferença complementar de merge, sem solução completa exclusiva ou família nova. Conferências de proveniência e literais, limites e parecer na [bancada 0154](../../lab/bancada/insights/descoberta-nativa-2026-11-08/README.md).
