# Conclusão dos pontos da revisão Tekt — 0131

Contrato anterior à implementação: passo 0131 e ADR-0014. Continua o recorte 0130 e preserva as alterações anteriores da worktree. A conclusão refere-se aos oito pontos inventariados na revisão; não é prova universal de semântica, segurança ou utilidade do produto.

## Fechamento por responsabilidade

| Ponto | Materialização e evidência |
| --- | --- |
| Regras de alternativas e perfil de screening em L2 | Extraídas em 0130; a jornada L1 também aplica a admissão do perfil opt-in antes de publicar. Oráculos CLI existentes preservados. |
| Política de retomada/jornada do dossiê em L2 | `recommendation_dossier.go` e `recommendation_dossier_artifacts.go` possuem identidade, versões, manifesto, estágios, retomada e execução. Portas lógicas adaptadas pela CLI; mensagens e cleanup preservam a ordem anterior. |
| Veredicto e execução experimental em L2 | `web_experiment_run.go` possui admissão, queries, intercalação, budgets/deadlines, retomada e veredicto. A descoberta simples também delega SearchOne a L1. |
| Composição de recomendação em L2 | A CLI não instancia RecommendationService nem escolhe pipeline/ranking. L4 fornece o serviço e as dependências; a API compatível usa a mesma fachada L1, sem segundo pipeline. Configuração dinâmica de regras permanece no núcleo. |
| Investigação de mercado misturada com runtime | `market_execution.go` implementa a porta observável. O runtime chama o executor configurado, observa progresso, preserva CAS/cancelamento e publica os artefatos. L4 compõe o executor explicitamente. |
| Resultado bruto em L0 | Dumps/contagens Clang movidos byte a byte para `lab/bancada/pesquisa/ast-cfg`; hashes em `moved-artifacts.json`. L0 conserva o relatório e os links. |
| Testes puros e auditor geral em L4 | Validadores em `global_market_validation_test.go` L1; auditor AST em `contracts/legacy_architecture_test.go`. Asserções preservadas em `oracle-preservation.json`; oráculos MCP duplicados já presentes em L2 foram consolidados. |
| Gate dependente de clone externo/patch manual | `make verify` e `make lint` usam fonte e patch fixados em `lab/ferramentas/crystalline-lint`, com checksums, Cargo.lock preservado, --locked/--offline e controles negativos. Sem alteração da instalação global. |

Também foi retirada de L2 a projeção dos targets originais de recuperação e sua associação a RuleIDs, agora em `evidence_recovery.go`.

## Preservação

`baseline.json` registra HEAD e fontes antes do corte. Os formatos, versões, hashes causais, policies, defaults, ordem dos artefatos, captura na mesma leitura do hash, estrito JSON, filtros de segredo e proteção de sidecars permanecem. A assinatura Git de fixtures é desativada somente no processo do gate; configuração global permanece intacta. Os artefatos históricos e módulos selados não foram reescritos nem promovidos.

`oracle-preservation.json` compara a transferência dos oráculos de domínio e AST. A diferença do auditor genérico é somente o diretório de trabalho delimitado pelo teste para conservar as mesmas raízes relativas. O teste de composição MCP anterior não reconhecia construtor encadeado com WithExecutionRuntime; a composição foi escrita em duas instruções, sem alterar o teste. `attempt-1` preserva a falha da primeira pilha e `mcp-composition-test.log` sua verificação focada.

O catálogo de mecanismos, recommendation.go e os consumidores CLI/MCP anteriores permanecem protegidos pelos próprios goldens/hashes e pela suíte. A fábrica validada de extrator substitui a extração fictícia usada para admitir configuração em L4.

## Verificação

Os logs focados por capacidade estão nesta pasta. `final/verification.json` contém os comandos e códigos de saída do gate completo: lint, auditorias, suíte raiz, vet, build e diff-check. O auditor produtivo inclui 353 arquivos Go. Os testes de ownership contêm controles de regressão para composição na casca, execução de busca web em L2, política de mercado no runtime, testes puros/gates em L4 e material não declarado em L0.

A fonte/patch do linter foi compilada com os 671 testes da biblioteca e 11 controles positivos/negativos do parser; `lint-tool-build.log` registra essa construção. `tool-lock.json` comprova a preservação do Cargo.lock congelado. Os controles do parser são repetidos no uso do cache validado. O código do produto não importa o Lab.

## Reprodução

```bash
make verify
```

Pré-requisitos e opções em `lab/ferramentas/README.md`: Go, Python, Git, Make, Rust/Cargo e crates fixados em cache para a primeira compilação offline. `BATEIA_VERIFY_OUTPUT` seleciona a evidência; `GOCACHE` e `BATEIA_TOOL_CACHE` podem usar diretórios autorizados. Os testes httptest precisam de loopback. O executável de verificação fica no diretório de evidência, sem deploy.

## Limites

As auditorias conhecem construções e contratos declarados, não toda a semântica possível. O gate inspeciona o grafo produtivo e os pontos de ownership desta revisão; testes de composição legítimos permanecem em L4. Dois arquivos JSON declarativos permanecem em L0: parâmetros do perfil opt-in e recibo causal de sucessão de contrato, ambos com papel distinto de resultado medido. Avaliação de produto e validade geral das organizações candidatas A1 conservam contratos próprios.

O controle adicional `configured-market-result-test.log` exige preservar também as conclusões do executor configurado: uma abstenção deliberada não é substituída por candidatos rederivados das observações. A última pilha completa passou no ambiente restrito, após timeout de revisão automática de permissões de uma tentativa que não chegou a iniciar; esse timeout não foi um veredito de insegurança nem exigiu mudança de testes.
