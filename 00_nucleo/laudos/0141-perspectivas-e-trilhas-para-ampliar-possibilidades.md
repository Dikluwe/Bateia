# Laudo 0141 — Perspectivas e trilhas para ampliar possibilidades

**Data:** 2026-10-28. **Resultado:** materialização e verificação aprovadas para o [contrato0141](../passos/0141-ampliar-possibilidades-com-perspectivas-e-trilhas.md).

## Entrega

O modo multifonte agora alterna perspectivas e fontes, com agenda determinística coprima que percorre cada par uma vez. O plano declara scheduling_policy; recibos permanecem em ordem query×fonte e recebem execution_order para as operações iniciadas. opportunity_coverage distingue tentativas, respostas válidas e fontes alcançadas. Lacunas de perspectivas não iniciadas ficam explícitas. O modo GitHub legado e o normalizador de explore não foram alterados.

READMEs admitidos geram reference_leads: candidato/identidade de origem, commit, caminho/hash, linhas/trecho literal e identidade do repositório citado. URLs são interpretadas como dados; rotas conhecidas de navegação, destinos já selecionados, self-links e URLs fora do contrato são excluídos. Limites128URLs/documento,64relações/12trilhas e1024bytes/trecho são explícitos. Tarefas recebem reference_lead_refs separado de catálogo/cartões. O destino permanece REFERENCE_NOT_READ; não houve rede adicional, classificação automática de mecanismo, promoção de evidência ou execução de candidatos.

As regras e o estado ficam em L1. L2/L3/L4 não mudaram. Foram acrescentados três arquivos produtivos e três de testes; três arquivos produtivos existentes foram alterados, sem modificar testes anteriores. Há também um replay de bancada no Lab, que não é dependência do produto.

## Verificação

make verify passou com lint/11 controles negativos, auditoria de382 arquivos produtivos, testes do módulo raiz, vet, build e diff-check. O gate usou loopback autorizado para fixtures HTTP gerais. Os seis verificadores históricos0106 e0134–0138 passaram com o binário final;265 arquivos de evidência anteriores permanecem idênticos ao baseline. Prompts e hashes atuais estão na linhagem nova.

Os controles puros verificam permutação de agenda para1–12 fontes, orçamento e cobertura. Para11 fontes elegíveis/20operações, a nova política fez quatro tentativas por perspectiva e alcançou todas as fontes. A política anterior query-major usava as20 células em apenas duas perspectivas. Isso é cobertura operacional demonstrada em fixture, não garantia de sucesso remoto.

O replay offline conserva os dois READMEs do ensaio comparacao-ia-2026-10-28, usando catálogo/leitor de fixture em memória. Produziu12trilhas, incluindo chnm/serendipomatic nas linhas106–107 da coleção MuseCompNet/muse-tech-central, com hash e trecho exatos. Cartões/bundles continuam apenas os dois documentos lidos; tarefas recebem as12 referências indiretas. O controle CRLF preserva bytes literais, e organizações/rotas de navegação não são apresentadas como raiz de repositório.

## Limites

A agenda amplia oportunidades de consulta; timeouts, auth, locale e quotas ainda podem impedir perspectivas. Consultas permanecem lexicais. Trilhas expõem caminhos adicionais, sem confirmar existência de destino, relação útil ou mecanismo; links relativos não são resolvidos. A seleção é limitada e lexical/estrutural, com omissões registradas. Leitura automática de destinos é uma etapa futura com contrato próprio. Não foi medido ganho de novidade/utilidade, e o critério desta fatia é ampliar possibilidades consideradas pela IA. Novidade NOT_ASSESSED e efeito UNVERIFIED permanecem.

Consultar [bancada0141](../../lab/bancada/insights/0141/README.md), [recibo consolidado](../../lab/bancada/insights/0141/verification.json) e [linhagem](../../lab/bancada/insights/0141/lineage.json). Os resultados anteriores foram preservados como históricos.
