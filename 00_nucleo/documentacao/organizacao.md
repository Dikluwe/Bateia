# Plano de organização do Bateia

A primeira etapa é tornar o funcionamento atual legível por capacidade. O projeto já tem camadas, contratos e acervo de decisões; reorganização física só deve ocorrer quando reduzir um acoplamento identificado e preservar comportamento por testes.

## Entrada organizada — 2026-11-08

O [passo 0148](../passos/0148-organizar-entrada-e-investigar-xerj.md) separa o README de entrada, o [guia de investigação](investigar-software.md) e a [referência completa de comandos](referencia-comandos.md). Os exemplos anteriores foram preservados com seus links ajustados. O caso XERJ fica no Lab e distingue coleta, contexto automático e propostas interpretadas pela IA. Este corte melhora navegação e uso; não declara refatoração integral do código ou encerramento das alterações locais anteriores.

## Organização aplicada

O [mapa arquitetural](../arquitetura/README.md) liga responsabilidades e contratos. A primeira fatia isola análise pura de hipóteses do protocolo de execução. O passo 0104 corrige admissão e projeção do registro de fontes; o passo 0105 acrescenta investigação offline de insights sem decisão de reúso. A linhagem das adições locais de pesquisa e Bitbucket foi regularizada preservando seu comportamento. O [passo 0103](../passos/0103-testar-hierarquia-de-arquiteturas-sob-tekt.md) registra o recorte; trabalhos restantes continuam abaixo.

## Retomada de arquitetura em 2026-10-28

O [passo 0127](../passos/0127-separar-resolucao-de-fontes-da-sintese.md) isola a resolução de fontes da síntese em L1. O resolvedor passa a depender de uma porta de leitura com dois métodos; o protocolo controla admissão e execução, e o analisador conserva conclusões de negócio. A caracterização é congelada antes da extração e reexecutada depois. Consultar o [laudo 0127](../laudos/0127-resolucao-de-fontes-separada.md) para resultados e limites.

As fronteiras inventariadas de efeitos e configuração foram materializadas no [passo 0129](../passos/0129-concluir-fronteiras-arquiteturais.md): relógio/deadlines por porta L1, host L3 para persistência e captura, containers de capacidade separados e testes de conectores em L3. A [bancada 0129](../../lab/bancada/arquitetura/0129/README.md) reúne a reprodução do lint e a auditoria de fronteiras. Essa verificação estrutural não encerra a separação semântica de responsabilidades: ainda havia regras na CLI e composição condicional de serviços fora de L4.

O [passo 0130](../passos/0130-extrair-regras-de-alternativas-e-retomada.md) prioriza regras de alternativas e retomada: L1 passa a decidir remapeamento do próprio alvo, validade semântica do perfil opt-in e elegibilidade/hash dos checkpoints. L2 mantém transporte e adaptação dos leitores. Os critérios vigentes e seus consumidores são preservados. O passo seguinte 0131 conclui os cortes de veredicto/jornada web, composição/jornada do dossiê e investigação/runtime de mercado.

A próxima organização deve seguir capacidade e contratos observáveis; ampliar a árvore ou criar runtimes paralelos não decorre desse corte. Avaliação semântica de insights, utilidade para usuários e evolução das políticas de reúso mantêm seus contratos próprios.

O trabalho 0126 de Gitee/AtomGit permanece uma alteração local anterior, com verificação descrita em seu laudo. A retomada 0127 não transforma essa integração ou as bancadas em validação integral da arquitetura. O tamanho de um arquivo ajuda a navegar; a escolha do corte decorre das responsabilidades e dependências.

## Primeiras correções de divergências — 0128

O [passo 0128](../passos/0128-explicitar-instante-de-coleta-nos-manifestos.md) torna determinísticos os parsers de manifestos: 48 consultas ao relógio são removidas de 43 arquivos, preservando o parsing em corpus autoral de 49 tipos. ParseManifestAt recebe a data explicitamente; as APIs anteriores continuam disponíveis com data desconhecida (zero). GitHub/local adquirem a data em L3; reanálise usa a data da extração. O [ADR-0011](../ADR/ADR-0011-instante-de-coleta-explicito-nos-manifestos.md) registra essa mudança de metadado.

Na verificação 0128, a tabela Go resolveu o carregamento, mas 46 avisos V12 ainda impediam aprovação. O passo 0129 identifica e corrige a classificação do parser numa reprodução isolada, sem desativar regras; também materializa as fronteiras temporais e de arquivos que estavam pendentes. O registro histórico 0128 permanece intacto.

## Entrega documental inicial

Esta pasta fornece uma visão de produto, um mapa de arquitetura e um guia de engenharia. O [diagnóstico inicial](../relatorios/diagnostico-organizacao-2026-10-06.md) distingue condições verificadas e pendências. Os próximos itens são propostas de trabalho, não ADRs aceitos nem implementação concluída.

## Sequência proposta

| Ordem | Trabalho | Por que priorizar | Critério de aceite |
| --- | --- | --- | --- |
| 1 | Linhagem e prompts das adições locais regularizados | Contratos e cabeçalhos adicionados | Especificações coerentes, cabeçalhos e lint sem erros; testes dos conectores passam |
| 2 | Ampliar o mapa inicial por capacidade | O recorte atual cobre os fluxos centrais | Cada capacidade principal liga contrato, fachada, portas, testes, evidência e limites |
| 3 | Uniformizar estados de entrega e validação | Implementação e qualidade demonstrada possuem critérios diferentes | Cada capacidade informa disponibilidade, verificação e avaliação de produto separadamente |
| 4 | Tornar a verificação repetível em ambiente limpo | Testes dependem de loopback e configuração Git da máquina | Procedimento reproduzível para testes, lint, vet e build, com execução de CI se adotada |
| 5 | Delimitar efeitos e determinismo em L1 | Uso de relógio conflita com a intenção de pureza | Exegese e contrato para uma primeira fatia, seguido de testes de preservação |
| 6 | Avaliar persistência do dossiê e fronteiras de apresentação | A CLI concentra arquivos, checkpoints e projeções | Decisão de ownership registrada e migração pequena com compatibilidade comprovada |
| 7 | Avaliar e reduzir despacho/composição em L4 | Configuração e features ampliam o ponto de entrada | Mapa de responsabilidades e corte justificado sem mover política para adaptadores |
| 8 | Ampliar avaliação da decisão de reúso | A promessa depende da utilidade para ideias reais | Corpus e julgamentos declarados, resultados por etapa e abstenções explicadas |

## Inventário por capacidade

Para cada capacidade, preencher:

- Resultado para o usuário e interface de entrada.
- Contrato proprietário e decisão arquitetural aplicável.
- Fachada, tipos, portas e adaptadores efetivamente usados.
- Artefatos produzidos, identidade e versões dos formatos.
- Testes ativos e evidência de avaliação com recorte temporal.
- Limites conhecidos, falhas parciais e dúvidas ainda abertas.

Começar pela recomendação ponta a ponta. Depois registrar descoberta de mercado/pesquisa, coleta e comprovação, rastreabilidade/revisão e runtime durável. Essa ordem acompanha a promessa central do produto.

## Política de mudança

Preservar nomes históricos, hashes e gates selados. Antes de mover um pacote ou arquivo, registrar o problema concreto, dependências afetadas, comportamento preservado e critério de reversão. Novos agrupamentos por capacidade podem inicialmente ser documentais; a estrutura física só muda depois da exegese e do contrato.

As regras estruturais de agentes estão centralizadas em `.agentrules`. Após autorização explícita do usuário, `.agents/AGENTS.md` passou a exigir sua leitura integral e incorporar todas as regras por referência. A comparação anterior à consolidação confirmou a preservação literal das seções históricas de regras na fonte canônica. O `AGENTS.md` da raiz mantém a orientação complementar de navegação e domínio.

O mesmo recorte estabiliza a ordem dos objetos JSON de dependências em cinco parsers. Preserva grupos e conteúdo das dependências, evitando variação causada pela iteração de mapas Go.

## Pontos da revisão concluídos — 0131

O [passo 0131](../passos/0131-concluir-divergencias-da-revisao-tekt.md) e o [laudo](../laudos/0131-conclusao-da-revisao-tekt.md) registram o fechamento dos oito pontos da revisão após o primeiro recorte 0130. Regras e jornadas ficam em L1, transporte/renderização em L2, efeitos em L3 e composição/lifecycle em L4. Resultados brutos e auditores gerais ficam no Lab; o gate canônico é `make verify` com ferramenta fixada no projeto.

A sequência proposta inicial acima permanece como contexto do planejamento. A conclusão atual diz respeito aos cortes arquiteturais inventariados e à sua preservação observável. Avaliação de relevância de buscas, qualidade das hipóteses, utilidade para usuários e a validade geral das organizações candidatas A1 mantêm seus próprios contratos.
