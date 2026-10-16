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
