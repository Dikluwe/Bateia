# 0131 — Concluir as divergências restantes da revisão Tekt

**Data:** 2026-10-28. **Governante:** Tekt e ADR-0008, preservando ADR-0012/0013.

## Contrato anterior à materialização

Concluir os pontos inventariados na revisão, após 0130: mover jornada, identidade e manifesto de recomendação para L1; compor executores/fábricas em L4; mover admissão, scheduling e veredicto de experimento web para L1; delegar investigação de mercado a executor L1 observável sem store/fila; manter lifecycle/publicação no serviço durável; retirar resultados brutos de L0 e testes puros/gates gerais de L4; padronizar o gate com versão reproduzível do linter e auditoria do regime dos artefatos.

Preservar contratos, formatos, hashes históricos, versões, defaults, ordem de efeitos, rejeição prévia à publicação, retomada e cancelamento. Nenhum novo pipeline, runtime ou promoção de evidência. Os consumidores CLI continuam públicos e seus testes exercitam o fluxo composto, com adaptação em L2 e efeitos L3. Código legado de compatibilidade pode delegar a uma fachada L1, sem manter uma segunda regra. Dados físicos permanecem atrás das portas.

O executor de mercado possui planejamento, budget, coleta e interpretação; aceita observação dos mesmos estágios e respeita cancelamento antes de novas chamadas. O serviço durável publica os mesmos documentos em sua ordem e preserva CAS/versionamento/cancelamento. Um executor configurado deve ser utilizado; a compatibilidade sem executor explícito utiliza a mesma implementação de investigação, sem duplicação.

Resultados Clang preservam bytes ao mover para o Lab. Testes de validadores vão para L1 e auditores gerais para Lab; testes de composição/lifecycle ficam em L4. Gates selados e recibos históricos não são reescritos.

O gate canônico do projeto deve construir a ferramenta de fonte e patch fixados, sem exigir checkout externo ou modificar a instalação global. Nenhuma regra/severidade será desativada. Deve verificar fontes produtivos, ownership conhecido desta revisão, localização dos artefatos e ausência de dependência produtiva do Lab.

## Aceite

Cada corte passa nos contratos focados antes de seguir. Na conclusão: auditorias, lint, suíte raiz, vet e build passam. Baseline e resultados em `lab/bancada/arquitetura/0131`; laudo registra cobertura e limitações. O aceite fecha os pontos declarados; não constitui prova universal de semântica ou utilidade de produto.

## Complemento de ownership
A projeção dos targets originais de recuperação, incluindo a associação de tipo de evidência a RuleID, passa de idea_recover.go para o proprietário L1 evidence_recovery.go. Preservar ordem, primeiro capability, IDs e duplicatas atuais; a casca apenas lê os bundles e delega.

## Registro da execução
Consultar o [laudo 0131](../laudos/0131-conclusao-da-revisao-tekt.md), a [bancada](../../lab/bancada/arquitetura/0131/README.md) e o gate canônico `make verify`. Os registros históricos 0129/0130 e os módulos selados não foram promovidos por este resultado.
