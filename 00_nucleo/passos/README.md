# Passos do Bateia

Passos são unidades executáveis e verificáveis que promovem desejos em software
sem tentar materializar o desejo inteiro de uma vez.

Cada passo deve declarar:

- qual incerteza reduz;
- qual observável produz;
- como pode falhar;
- como será verificado;
- o que permanece explicitamente fora de escopo.

## Sequência

| Passo | Objetivo | Estado |
| :---- | :------- | :----- |
| [0001](0001-materializar-ideaspec-e-plano-de-expansao.md) | Materializar `IdeaSpec → ExpansionPlan` como dado puro e auditável | `IMPLEMENTADO` |
| [0002](0002-expandir-vocabulario-com-baseline-mensuravel.md) | Expandir vocabulário com baseline lexical offline e avaliação comparável | `IMPLEMENTADO` |
| [0003](0003-executar-expansoes-e-medir-precisao-da-descoberta.md) | Executar expansões e medir precisão real da descoberta | `IMPLEMENTADO` |
| [0004](0004-construir-peneira-probabilistica-explicavel.md) | Rastrear o ruído e construir uma peneira probabilística explicável | `IMPLEMENTADO` |
| [0005](0005-triar-candidatos-e-comprovar-capacidades.md) | Triar candidatos e comprovar capacidades sem LLM | `IMPLEMENTADO` |
| [0006](0006-endurecer-e-calibrar-evidencias.md) | Endurecer a coleta e calibrar evidências reais antes da decisão | `IMPLEMENTADO` |
| [0007](0007-diagnosticar-e-ampliar-cobertura-de-evidencias.md) | Diagnosticar perdas de cobertura e atingir o gate sem afrouxar evidências | `IMPLEMENTADO` |
| [0008](0008-comparar-candidatos-e-decidir-reuso.md) | Comparar candidatos e decidir o ponto de partida sem LLM | `IMPLEMENTADO` |
| [0009](0009-orquestrar-funil-e-entregar-dossie-de-reuso.md) | Orquestrar o funil e entregar um dossiê reproduzível de reúso | `EM VALIDAÇÃO` |
| [0010](0010-diagnosticar-inconclusoes-e-coletar-adaptativamente.md) | Diagnosticar inconclusões e coletar somente a evidência faltante | `EM VALIDAÇÃO` |
| [0011](0011-congelar-benchmark-ponta-a-ponta-e-auditar-utilidade.md) | Congelar corpus correspondente e auditar evidência, decisão, recuperação e utilidade | `EM VALIDAÇÃO` |
| [0012](0012-executar-autofagia-competitiva-e-formar-katamari-minimo.md) | Auditar concorrentes e propor a menor composição de capacidades que cubra lacunas reais | `EM VALIDAÇÃO` |
| [0013](0013-limitar-descoberta-e-preservar-resultados-parciais.md) | Impor deadlines, registrar progresso e preservar snapshots parciais da descoberta | `EM VALIDAÇÃO` |
| [0014](0014-construir-ontologia-lexical-para-descoberta-de-software.md) | Traduzir intenção de reúso em queries conceituais discriminatórias e medir recall | `IMPLEMENTADO` |
| [0015](0015-separar-catalogo-de-projetos-e-comprovacao-em-codigo.md) | Restaurar GitHub como catálogo e restringir Sourcegraph à comprovação em código | `EM VALIDAÇÃO` |
| [0016](0016-comparar-superficies-publicas-de-catalogo-e-recuperar-recall.md) | Comparar README, tópicos, pacotes e outro catálogo por ganho marginal de recall | `EM VALIDAÇÃO` |
| [0017](0017-gerar-consultas-por-frames-funcionais-e-descobrir-produtos-na-web.md) | Descrever o trabalho do produto com frames funcionais e comparar busca nominal versus comportamental | `EM VALIDAÇÃO` |
| [0018](0018-versionar-gramatica-funcional-e-congelar-busca-web-reexecutavel.md) | Corrigir a gramática funcional e tornar a busca web congelável, reexecutável e causalmente comparável | `EM VALIDAÇÃO` |
| [0019](0019-executar-experimento-publico-da-gramatica-funcional.md) | Executar N/F1/F2 sem token, congelar resultados públicos e decidir a gramática por replay e julgamento cego | `IMPLEMENTADO` |
| [0020](0020-executar-busca-autenticada-com-orcamento-observavel.md) | Repetir a descoberta com token somente leitura, orçamento fixo, segurança e replay offline | `IMPLEMENTADO` |
| [0021](0021-julgar-candidatos-com-ia-e-escalonamento-humano.md) | Obter primeiro veredito da IA e escalar dúvidas e casos de borda ao usuário | `IMPLEMENTADO` |
| [0022](0022-calibrar-regras-de-evidencia-por-oraculos-independentes.md) | Ampliar a comprovação de capacidades com oráculos segregados e teste oculto | `IMPLEMENTADO` |
| [0023](0023-validar-funil-completo-para-decisoes-reais-de-reuso.md) | Validar descoberta, evidência, decisão e utilidade em dez ideias reais do domínio escolhido | `DESCOBERTA APROVADA — DOWNSTREAM PENDENTE` |
| [0024](0024-ranquear-alternativas-funcionais.md) | Ranquear alternativas funcionais para um software alvo com evidência e filtros | `IMPLEMENTADO` |
| [0025](0025-expor-recomendacao-completa-como-ferramenta-mcp.md) | Expor o funil completo como ferramenta MCP estruturada para uma IA externa | `IMPLEMENTADO` |
| [0026](0026-persistir-execucoes-mcp-e-consultar-por-run-id.md) | Persistir execuções MCP, acompanhar estado e ler artefatos seguros por `run_id` | `IMPLEMENTADO` |
| [0027](0027-generalizar-descoberta-multilingue-e-multifuente.md) | Descobrir produtos em múltiplos idiomas, mercados e fontes com proveniência | `IMPLEMENTADO` |
| [0028](0028-extrair-mecanismos-e-formular-hipoteses-de-melhoria.md) | Extrair mecanismos comprovados e formular hipóteses de melhoria rastreáveis | `IMPLEMENTADO` |
| [0029](0029-medir-ganho-da-ia-usando-bateia.md) | Medir em Luna, Terra e Sol se usar o Bateia melhora decisões de engenharia | `IMPLEMENTADO` |
| [0030](0030-corrigir-teste-de-catalise-da-busca.md) | Corrigir pareamento, limites de saída e telemetria do teste IA+web versus IA+Bateia+web | `IMPLEMENTADO` |
| [0031](0031-autopreencher-parametros-da-recomendacao.md) | Autopreencher política, alvo e busca omitidos nas ferramentas MCP de recomendação | `IMPLEMENTADO` |
| [0032](0032-executar-teste-end-to-end-de-validacao-do-bateia.md) | Executar teste end-to-end de validação do Bateia na CLI e servidor MCP | `IMPLEMENTADO` |
| [0033](0033-expandir-descoberta-global-federada-de-mercados-e-locales.md) | Expandir descoberta global federada de mercados e locales no catálogo e MCP | `IMPLEMENTADO` |
| [0034](0034-integrar-conectores-android-e-forjas-federadas-de-codigo.md) | Integrar conectores Android (Play Store & F-Droid) e forjas federadas (Gitee, Codeberg, GitLab) | `IMPLEMENTADO` |
| [0035](0035-corrigir-lint-e-analisar-c-cpp-com-clang.md) | Corrigir precisão do lint e integrar análise estática C/C++ com Clang e cobertura explícita | `IMPLEMENTADO` |

- [0036 — Projetar evidências da consulta](0036-projetar-evidencias-da-consulta.md)

- [0037 — Resposta fundamentada por mecanismo](0037-resposta-fundamentada-por-mecanismo.md)

- [0038 — Protótipo Clang AST + CFG](0038-prototipo-clang-ast-cfg.md)

- [0039 — Analisar a jornada principal de reúso](0039-analisar-jornada-principal-de-reuso.md)

- [0040 — Corrigir autoridade da evidência](0040-corrigir-autoridade-da-evidencia.md)

- [0041 — Reestruturar autofagia e Katamari](0041-reestruturar-autofagia-e-katamari.md)

- [0042 — Auditar promessas e limites do projeto](0042-auditar-promessas-e-limites-do-projeto.md)

- [0043 — Delimitar promessas e planejar buscas verificáveis](0043-delimitar-promessas-e-planejar-buscas-verificaveis.md) — [Execução A–F e ensaio G; utilidade externa não avaliada](../relatorios/execucao-promessas-buscas-0043.md)

- [0044 — Recuperar implementações por símbolos](0044-recuperar-implementacoes-por-simbolos.md) — [Execução e limites](../relatorios/recuperacao-implementacoes-0044.md)

- [0045 — Recuperação por perguntas funcionais](0045-recuperar-perguntas-funcionais.md) — [Execução e limites](../relatorios/recuperacao-funcional-0045.md)

- [0046 — Projetos reais e lacunas por requisito](0046-validar-recuperacao-em-projetos-reais.md) — [Execução: sem ganho demonstrado de recall](../relatorios/recuperacao-real-0046.md)

- [0047 — Ponderar termos pelo corpus](0047-ponderar-termos-pelo-corpus.md) — [Execução e ganho delimitado](../relatorios/ranking-corpus-0047.md)

- [0048 — Distinguir consulta natural e símbolos](0048-distinguir-consulta-e-simbolos.md) — [Execução e restrições não verificadas](../relatorios/consulta-simbolos-0048.md)

- [0050 — Recuperar implementações remotas](0050-recuperar-implementacoes-remotas.md) — [Execução pública e limites](../relatorios/recuperacao-remota-0050.md)

- [0051 — Recuperar documentação de declarações](0051-recuperar-documentacao-de-declaracoes.md) — [Execução e limites](../relatorios/documentacao-declaracoes-0051.md)

- [0052 — Expansão lexical com proveniência](0052-expandir-vocabulario-com-proveniencia.md), com [adendo de admissão por ação](0052-adendo-admissao-por-acao.md).

- [0053 — Avaliação por estágio em perguntas novas](0053-avaliar-jornada-em-perguntas-novas.md) — [Resultados e próximo gargalo](../relatorios/avaliacao-jornada-0053.md)

- [0054 — Preservar intenção na descoberta](0054-preservar-intencao-na-descoberta.md) — Remover fallback destrutivo e medir no corpus congelado.

- [0055 — Consultas de código como hipóteses rastreáveis](0055-planejar-consultas-de-codigo.md).

- [0056 — Cobertura observada e perdas na admissão](0056-explicar-cobertura-e-perdas.md).

- [0057 — Observações de comportamento por função](0057-observar-comportamento-por-funcao.md).

- [0062 — Verificar efeito executado](0062-verificar-efeito-executado.md) — [Piloto e limites](../relatorios/efeitos-executados-0062.md).

- [0063 — Vincular fonte ao efeito](0063-vincular-fonte-ao-efeito.md) — [Execução e rastreabilidade](../relatorios/vinculo-fonte-efeito-0063.md).

- [0064 — Identificar analisador e parâmetros](0064-identificar-analisador-e-parametros.md) — [Recibos e limites](../relatorios/rastreabilidade-analisador-0064.md).

## Organização arquitetural em preparação

- [0103 — Testar a hierarquia de arquiteturas sob Tekt](0103-testar-hierarquia-de-arquiteturas-sob-tekt.md) — mapa inicial e primeira separação materializados; avaliação da hierarquia em andamento.

- [0104 — Compor descoberta por registro validado](0104-compor-descoberta-por-registro-validado.md) — implementado e verificado localmente: admissão em L1, composição em L4 e schema do registro em L2.

- [0105 — Investigar insights sem decisão de reúso](0105-investigar-insights-sem-decisao-de-reuso.md) — implementado e verificado no recorte offline; utilidade real ainda não avaliada.

- [0106 — Corrigir esperas e investigar GitHub](0106-corrigir-esperas-e-investigar-github.md) — suíte concorrente aprovada e fonte real coletada; aplicabilidade ainda pendente.

- [0107 — Pesquisar qualidade do classificador com Bateia](0107-pesquisar-qualidade-do-classificador-com-bateia.md) — pesquisa exploratória concluída; avaliação semântica independente proposta.

- [0108 — Refinar pesquisa antes do avaliador](0108-refinar-pesquisa-antes-do-avaliador.md) — propriedades, configuração, referência e escopo do piloto refinados.

- [0109 — Observar propriedades de concorrência no Lab](0109-observar-propriedades-de-concorrencia-no-lab.md) — piloto sintático, CLI e corpus verificados; qualidade semântica pendente.

- [0110 — Resolver canais nomeados no Lab](0110-resolver-canais-nomeados-no-lab.md) — checagem restrita em memória; dependências de conc ainda incompletas.

- [0111 — Fornecer dependências rastreáveis no Lab](0111-fornecer-dependencias-rastreaveis-no-lab.md) — limiter resolvido com ambiente congelado; fluxo global pendente.

- [0112 — Relacionar reserva, lançamento e liberação](0112-relacionar-reserva-lancamento-e-liberacao.md) — relações tipadas e contrastes; alcance/callback/instância pendentes.

- [0113 — Rastrear parâmetros de callback](0113-rastrear-parametros-de-callback.md) — cadeia estrutural até Try/f(); execução e alcance pendentes.

- [0114 — Verificar alcance e ordem limitados](0114-verificar-alcance-e-ordem-limitados.md) — retorno/prefixo de defer distinguidos; alcance global e instância pendentes.

- [0115 — Auditar identidade e mutação do canal](0115-auditar-identidade-e-mutacao-do-canal.md) — writes/captures localizados; identidade de instância ainda inconclusiva.

- [0116 — Observar guarda da configuração](0116-observar-guarda-da-configuracao.md) — checagem anterior localizada; atomicidade/instância pendentes.

- [0117 — Testar ordem sequencial e intercalação](0117-testar-ordem-sequencial-e-intercalacao.md) — quatro testemunhos autorais; resultados restritos ao modelo.

- [0118 — Confrontar alegações com contrato](0118-confrontar-alegacoes-com-contrato.md) — revisão autoral rastreável; reutilização e reconfiguração distinguidas, precisão independente pendente.

- [0119 — Executar referência original do candidato](0119-executar-referencia-original-do-candidato.md) — testes originais de conc executados; quatro comportamentos confirmados no recorte, precisão semântica geral pendente.

- [0120 — Avaliar configuração em fontes externas](0120-avaliar-configuracao-em-fontes-externas.md) — observador congelado em dois projetos; referência estrutural localizada, cobertura contratual ausente e precisão independente pendente.

- [0121 — Ligar documentação a evidência estática](0121-ligar-documentacao-a-evidencia-estatica.md) — vista documental no Lab, com donos AST e associações fornecidas explícitas; interpretação contratual pendente.

- [0122 — Pesquisar modelos de decisão tipados](0122-pesquisar-modelos-de-decisao-tipados.md) — coleta com Bateia, alternativas e quatro payloads diagnósticos; inferência/calibração pendentes.

- [0123 — Testar Laya e planejar treinamento](0123-testar-laya-e-planejar-treinamento.md) — inferência real isolada em CPU; três concordâncias diagnósticas, plano de treino sem dataset/ajuste executado.

- [0124 — Investigar melhorias do Laya com Bateia](0124-investigar-melhorias-do-laya-com-bateia.md)

- [0127 — Separar resolução de fontes da síntese](0127-separar-resolucao-de-fontes-da-sintese.md) — porta de leitura em L1 e caracterização antes/depois; testes, vet e build aprovados, lint padrão bloqueado por incompatibilidade anterior.

- [0128 — Explicitar instante de coleta nos manifestos](0128-explicitar-instante-de-coleta-nos-manifestos.md) — parsers determinísticos, coleta temporal em L3 e configuração do lint compatível; avisos arquiteturais anteriores permanecem.

- [0129 — Concluir fronteiras arquiteturais](0129-concluir-fronteiras-arquiteturais.md) — portas temporais e de arquivos, host L3, containers de capacidade e verificação reproduzível.

- [0130 — Extrair regras de alternativas e retomada](0130-extrair-regras-de-alternativas-e-retomada.md) — alvo local, perfil opt-in e política de checkpoint em L1, com leitores restritos e preservação dos consumidores CLI.

- [0131 — Concluir divergências da revisão Tekt](0131-concluir-divergencias-da-revisao-tekt.md) — jornadas em L1, execução de mercado delegada, composição em L4, bancada no Lab e gate reproduzível do projeto.

- [0132 — Explicar retomada de evidências](0132-explicar-retomada-de-evidencias.md) — diagnóstico tipado em L1 e mensagem na CLI, preservando elegibilidade, privacidade e identidade causal.

- [0133 — Priorizar insights de engenharia](0133-priorizar-insights-de-engenharia.md) — foco de produto, critérios de utilidade e lacuna prioritária registrados; sem alteração funcional ou promoção de evidências históricas.
- [0134 — Conservar observações específicas nos insights](0134-conservar-observacoes-especificas-nos-insights.md) — contrato v2 com observações declaradas/localizadas, condições operacionais separadas da cobertura e replay v1 preservado.
- [0135 — Comparar insights com baseline do alvo](0135-comparar-insights-com-baseline-do-alvo.md) — snapshot explícito do alvo, declarações por fonte/classe/tipo e perguntas para validar diferenças, sem inferir equivalência ou efeito.
- [0136 — Preparar exploração criativa para IA](0136-preparar-exploracao-criativa-para-ia.md) — pergunta direta, cartões automáticos variados e tarefas de alternativas/inversão/analogia/recombinação/contraexemplos, sem curadoria prévia de avaliações.
- [0137 — Descobrir referências diversificadas automaticamente](0137-descobrir-referencias-diversificadas-automaticamente.md) — buscas independentes, origens, seleção estrutural e coleta de READMEs fixados alimentando o contexto para IA.
- [0138 — Ampliar descoberta de insights para outros sítios](0138-ampliar-descoberta-de-insights-para-outros-sitios.md) — fontes federadas, incluindo Gitee/AtomGit, com orçamento global, identidade por sítio e contexto de catálogo separado das leituras.

- [0139 — Leitura automática de README em Gitee e AtomGit](0139-ler-readmes-de-gitee-e-atomgit.md): SHA fixado, contextos por sítio, cartões para IA e budgets compartilhados.

- [0140 — Leitura automática de README em Codeberg e GitLab](0140-ler-readmes-de-codeberg-e-gitlab.md) — APIs públicas, revisão fixa, subgrupos GitLab e contexto para IA.

- [0141 — Ampliar possibilidades com perspectivas e trilhas](0141-ampliar-possibilidades-com-perspectivas-e-trilhas.md) — agenda equilibrada e referências documentais para novas investigações da IA.

- [0142 — Aprofundar trilhas e conservar memória](0142-aprofundar-trilhas-e-conservar-memoria.md) — leitura de um salto e snapshots com feedback.

- [0143 — Expandir vizinhanças a partir da memória](0143-expandir-vizinhancas-a-partir-da-memoria.md) — continuação focada sem nova pergunta ou catálogo, com fronteira e um salto por chamada.

- [0144 — Continuar o foco e selecionar vizinhos](0144-continuar-foco-e-selecionar-vizinhos.md) — foco herdado, fronteira ordenada por proximidade e aliases conhecidos sem tentativas duplicadas.

- [0145 — Reorganizar fronteiras de investigação sob Tekt](0145-reorganizar-fronteiras-de-investigacao-sob-tekt.md) — admissão, fontes, fronteira, execução, memória e resultado com ownership explícito; interpretação de argumentos em L2.

- [0146 — Consolidar prompts de exploração](0146-consolidar-prompts-de-exploracao.md) — contrato vigente, matriz de aceitação e obrigações locais; distinguir consolidação documental de autoria ex-ante.

- [0147 — Corrigir limites e verificar contratos de exploração](0147-corrigir-limites-e-verificar-contratos-de-exploracao.md) — fechar o avanço de geração e ampliar controles de preservação, fontes congeladas e desempates.

- [0148 — Organizar a entrada e investigar melhorias para o XERJ](0148-organizar-entrada-e-investigar-xerj.md) — guia de primeiro uso, referência de comandos e caso reproduzível de investigação, com propostas e limites explícitos.

- [0149 — Reproduzir herança G7 no XERJ](0149-reproduzir-heranca-g7-no-xerj.md) — reprodução local, rastreamento dos gates de publicação e teste de um agente usando somente a recuperação do XERJ.

- [0150 — Comparar jornadas completas Bateia e XERJ](0150-comparar-jornadas-completas-bateia-xerj.md) — dois agentes sem referências prontas, com descoberta nativa, aquisição registrada e propostas independentes.

- [0151 — Ampliar referências por linguagem](0151-ampliar-referencias-por-linguagem.md) — investigação suplementar de mecanismos em sete linguagens, preservando a comparação original.

- [0152 — Comparar alternativas em problema aberto](0152-comparar-alternativas-em-problema-aberto.md) — problema concreto sem mecanismos sugeridos, baselines congelados e avaliação anônima das decisões propostas.

- [0153 — Ligar coleta à revisão de alternativas](0153-ligar-coleta-a-revisao-de-alternativas.md) — contexto atualizado após coleta, passagens adicionais e revisão de cobertura; reteste independente registrado no Lab.

- [0154 — Corrigir descoberta nativa e testar sem memória](0154-corrigir-descoberta-nativa-e-testar-sem-memoria.md) — consultas contextuais, correspondência em metadados, pin nas leituras escolhidas e ensaio com origem nativa obrigatória.

- [0155 — Descobrir princípios e decisões no XERJ](0155-descobrir-principios-e-decisoes-no-xerj.md) — exploração de algoritmos e registro de evidências/condições/decisões; ensaio no XERJ com referências nativas.

- [0156 — Demonstrar Bateia e XERJ sobre postings](0156-demonstrar-bateia-e-xerj-sobre-postings.md) — descoberta nativa, mesmo commit, indexação e recuperação progressiva com evidência.

- [0157 — Fechar duplicatas de flags e limites de formulação](0157-fechar-duplicatas-e-limites-de-formulacao.md) — unicidade booleana em explore/discover e ampliação dos consumidores de limites e identidade da formulação.

- [0158 — Validar origem da coleta e estabilizar bundles](0158-validar-origem-e-estabilizar-bundles.md) — admissão antecipada, vínculo repo@commit e serialização estável, sem reescrever evidência histórica.

## Complemento de navegação histórica

As entradas abaixo completam a navegação para passos, adendos e verificações já existentes. Os títulos são os dos documentos originais; inclusão no índice não altera seu estado, suas evidências ou seus critérios de aceite. A sequência acima conserva as descrições e os estados anteriormente registrados.

- [Adendo 0055 — Identidade de retomada](0055-adendo-checkpoint.md)
- [Adendo 0055 — Recibo da admissão na busca direta](0055-adendo-evidencia.md)
- [Revisão de contrato 0055 — Sintaxe validada no provedor](0055-adendo-sintaxe.md)
- [Adendo0057 — Não promover origem documental](0057-adendo-origem-documental.md)
- [Passo 0058 — Diagnosticar a perda de uma implementação conhecida](0058-diagnosticar-perda-do-alvo.md)
- [Passo 0059 — Preservar o canal documental sob limite de resultados](0059-preservar-canal-documental.md)
- [Passo 0060 — Tradução local explícita de caminho](0060-traducao-local-de-caminho.md)
- [Passo 0061 — Discriminar hipóteses de ação por candidato](0061-discriminar-acoes-estaticas.md)
- [Passo 0065 — Analisar coalizões de evidência em modo sombra](0065-analisar-coalizoes-de-evidencia-em-modo-sombra.md)
- [Passo 0066 — Prototipar radar de convergência de ideias em modo sombra](0066-prototipar-radar-de-convergencia-de-ideias.md)
- [Passo 0067 — Planejar a fronteira mínima do radar de ideias](0067-planejar-fronteira-minima-do-radar-de-ideias.md)
- [Passo 0068 — Compilar uma capacidade em consultas técnicas discriminatórias](0068-compilar-capacidade-em-consultas-tecnicas.md)
- [Passo 0069 — Exportar radar para Archify](0069-exportar-radar-para-archify.md)
- [Passo 0070 — Reconstruir história de algoritmos](0070-reconstruir-historia-de-algoritmos.md)
- [Passo 0071 — Expor cobertura do catálogo sem atribuir lacunas ao candidato](0071-expor-cobertura-do-catalogo-sem-atribuir-lacunas-ao-candidato.md)
- [Passo 0072 — Recuperar evidência local por requisito antes da amostra geral](0072-recuperar-evidencia-local-por-requisito-antes-da-amostra-geral.md)
- [Passo 0073 — Localizar a ocorrência probatória dentro do arquivo](0073-localizar-a-ocorrencia-probatoria-dentro-do-arquivo.md)
- [Passo 0074 — Alinhar regra de precondição ao vocabulário do GSL](0074-alinhar-regra-de-precondicao-gsl.md)
- [Passo 0075 — Localizar expressões probatórias compostas](0075-localizar-expressoes-probatorias-compostas.md)
- [Passo 0076 — Separar produto e bancada](0076-quarentenar-bancada-e-limpar-l0-l4.md)
- [Passo 0077 — Configuração durável em L3](0077-extrair-configuracao-duravel-para-l3.md)
- [Passo 0078 — Catálogo de mecanismos em L1](0078-mover-catalogo-de-mecanismos-para-l1.md)
- [Passo 0079 — Container durável em L3](0079-extrair-container-duravel-para-l3.md)
- [Passo 0080 — Conter descartes não calibrados](0080-conter-descartes-nao-calibrados.md)
- [Passo 0081 — Medir a qualidade da triagem](0081-medir-qualidade-da-triagem-em-holdout-cego.md)
- [Passo 0082 — Modelo operacional opt-in](0082-validar-modelo-operacional-opt-in-em-holdout-novo.md)
- [Passo 0083 — Estruturar entradas e expor abstenção interativa verificável](0083-estruturar-entradas-e-abstencao-interativa.md)
- [Passo 0084 — Revisão sistemática viva do funil](0084-revisao-sistematica-viva-do-funil.md)
- [Passo 0085 — Priorizar investigações por valor decisório](0085-priorizar-investigacoes-por-valor-decisorio.md)
- [Passo 0086 — Fronteira multicritério robusta](0086-fronteira-multicriterio-robusta.md)
- [Passo 0088 — Composições como candidatas](0088-composicoes-como-candidatas.md)
- [Passo 0089 — Avaliar abstenção da decisão final](0089-calibrar-abstencao-por-risco-cobertura.md)
- [Passo 0090 — Landscape de soluções e tecnologias](0090-landscape-de-solucoes-e-tecnologias.md)
- [Passo 0091 — Evidências técnicas especializadas](0091-evidencias-tecnicas-especializadas.md)
- [Passo 0092 — Verificar achados técnicos antes de usá-los](0092-verificar-achados-tecnicos.md)
- [Passo 0093 — Reler fontes técnicas no provedor](0093-reler-fontes-tecnicas-no-provedor.md)
- [Passo 0094 — Endurecer segurança e garantias](0094-endurecer-seguranca-e-garantias.md)
- [Passo 0095 — Confinar coleta local e rejeitar evidência ambígua](0095-confinar-coleta-local-e-rejeitar-evidencia-ambigua.md)
- [Passo 0096 — Orçamento agregado e identidade local versionada](0096-orcamento-agregado-e-identidade-local-versionada.md)
- [Passo 0097 — Evidências consistentes e ADR sem aprovação automática](0097-capturar-evidencias-consistentes-e-propor-adr.md)
- [0098 — Rastrear dependências e planejar revisão incremental](0098-rastrear-dependencias-e-planejar-revisao.md)
- [0099 — Localizar declarações no relatório original](0099-localizar-declaracoes-no-relatorio-original.md)
- [0100 — Comparar grafos e planejar revisão incremental](0100-comparar-grafos-e-planejar-revisao.md)
- [Verificação do incremento 0100](0100-verificacao.md)
- [0101 — Registrar o progresso consultivo das tarefas de revisão](0101-registrar-progresso-de-revisao.md)
- [0102 — Ligar tarefas de revisão às declarações originais](0102-ligar-revisao-as-declaracoes-originais.md)
- [0125 — Contexto de compilação e investigação Mesa](0125-investigar-clipping-e-transicao-gs.md)
- [0126 — Integrar Gitee e AtomGit à descoberta](0126-integrar-gitee-e-atomgit.md)
