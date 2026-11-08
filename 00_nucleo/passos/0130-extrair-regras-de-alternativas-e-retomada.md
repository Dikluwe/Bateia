# 0130 — Extrair regras de alternativas e retomada da CLI

**Data:** 2026-10-28  
**Governante:** Tekt, ADR-0008 e ADR-0013.

## Comportamento e limites anteriores à materialização

Priorizar as regras ainda em `02_shell/cli/idea_recommend.go`: remapeamento/exclusão do próprio alvo em worktrees locais, admissão semântica do perfil de screening opt-in e elegibilidade de retomada da evidência. Os três proprietários passam a L1; CLI conserva flags, JSON estrito, mensagens sanitizadas, leitura por host, projeção e publicação.

Preservar comportamento: remotes locais são hints consultivos e não autenticação; identidade ausente marca Stale; somente identidade igual ao alvo remapeia candidato/assessment/fontes; outros modos e ausência de alvo não consultam hints. Preservar ordem e quantidade das consultas e não mutar entradas. Ranking e políticas de evidência permanecem os atuais.

O perfil opt-in aceita somente schema 1, uncalibrated sem calibrador, pesos presentes, thresholds e taxa nos limites vigentes e números finitos. O validador específico não altera ScreenDiscovery, defaults, calibradores ou os caminhos que admitem modelos calibrados. Bytes, hash, RunID, erros de transporte e rejeição antes da saída permanecem idênticos.

A política de retomada L1 recebe manifesto, critérios causais e porta restrita de SHA-256 por caminho lógico. Somente COMPLETE/PARTIAL/INCONCLUSIVE podem reutilizar evidência; comparar schema, hashes causais, presença/hash do modelo e as oito versões atuais. Exigir discovery-snapshot.json e evidence/bundles.json do manifesto, usando a primeira ocorrência com hash não vazio e ignorando duplicatas posteriores já verificadas. Divergência ou erro de leitura recusa reaproveitamento; nenhuma leitura antes de admitir metadados. Preservar ordem das leituras, short-circuit, hashes históricos e critérios de mudança somente de política/alvo. Não acrescentar campos aos formatos públicos existentes, versões ou endurecimento de política nesta extração.

## Verificação

Congelar hashes dos fontes e executar contratos CLI antes/depois; novos testes L1 com leitores em memória para modos/hints/fontes, mutação das entradas, perfil opt-in e matriz de checkpoints, incluindo leituras e duplicatas. Reexecutar contratos de consumidor, auditor AST, lint reproduzível, suíte raiz, vet e build. Resultados brutos em `lab/bancada/arquitetura/0130`. Não reescrever testes/recibos selados nem alterações anteriores da worktree.

Este recorte não encerra a jornada inteira da CLI, o veredicto de experimentos web, a composição de serviços em L2 ou a separação restante do runtime de mercado.

## Registro da execução

Consultar o [laudo 0130](../laudos/0130-regras-de-alternativas-e-retomada.md), a [bancada](../../lab/bancada/arquitetura/0130/README.md) e os códigos de saída em checks.json. Consumidores e critérios foram preservados; o recorte não declara conformidade integral da arquitetura.
