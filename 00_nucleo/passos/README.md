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
