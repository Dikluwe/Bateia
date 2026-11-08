# 0130 — Regras de alternativas e retomada em L1

**Data:** 2026-10-28  
**Contrato:** [passo 0130](../passos/0130-extrair-regras-de-alternativas-e-retomada.md) e [ADR-0013](../ADR/ADR-0013-regras-de-alternativas-e-retomada-em-l1.md).  
**Evidência:** [bancada 0130](../../lab/bancada/arquitetura/0130/README.md).

## Resultado do recorte

Três decisões que estavam na CLI passam a proprietários L1: remapeamento consultivo para exclusão do próprio alvo local, admissão semântica do perfil de screening opt-in e elegibilidade/verificação de checkpoint de evidência. A CLI conserva transporte, JSON estrito, hashes dos bytes, erros sanitizados e adaptação dos leitores às portas físicas existentes.

O núcleo decide quais artefatos lógicos de checkpoint podem ser consultados e em que ordem. Metadados inelegíveis não provocam leitura; erro ou hash divergente recusa reaproveitamento, sem tentar uma duplicata posterior. Os mesmos estados, oito versões, presença/hash do modelo e dois artefatos obrigatórios continuam exigidos. Nenhuma política de evidência ou aprovação de reúso foi ampliada.

Hints Git permanecem consultivos. Ausência mantém a classificação Stale; somente correspondência ao alvo remapeia a projeção de candidato, assessment e fontes já tratadas. Entradas e identidade dos bytes são preservadas. O perfil opt-in continua separado das regras que admitem modelos calibrados em outras jornadas.

## Preservação e verificação

A mesma seleção de contratos ativos CLI foi executada antes e depois. Os fontes desses consumidores, da recomendação, do ranking e do screening geral mantêm os hashes do início do recorte; seus esperados não foram alterados. Formatos públicos, RunID, hashes e versões permanecem os atuais. Mudanças somente de política/alvo continuam permitindo reutilizar evidência causalmente compatível para recalcular a decisão.

Testes L1 usam leitores em memória para verificar invariantes, guardas, estados, erros, duplicatas, ordem e ausência de mutação. A pilha final e seus códigos de saída constam de [checks.json](../../lab/bancada/arquitetura/0130/checks.json): lint com parser validado de 0129, auditor AST, suíte raiz, vet, build e diff-check. A falha inicial de assinatura SSH de uma fixture e sua resolução restrita ao processo são preservadas na bancada. Não há promoção dos gates selados ou reescrita de seus recibos.

## Limites

Este é o primeiro corte priorizado da revisão. A construção do manifesto/identidade da jornada e a composição condicional de serviços permanecem em L2; veredicto experimental web, runtime de mercado e outras divergências inventariadas exigem recortes próprios. Os checks estruturais não demonstram conformidade semântica integral ou utilidade de produto.
