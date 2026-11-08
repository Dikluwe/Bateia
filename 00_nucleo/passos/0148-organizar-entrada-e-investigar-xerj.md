# 0148 — Organizar a entrada e investigar melhorias para o XERJ

**Data:** 2026-11-08. **Escopo:** navegação documental e uso das capacidades existentes.

## Intenção e comportamento esperado

Oferecer uma entrada curta no Bateia, separar o primeiro uso da referência completa de comandos e executar uma investigação cujo alvo de melhoria é o XERJ. O usuário confirmou essa direção. O repositório do Bateia permanece local; nenhuma publicação, envio de mensagem ou alteração no XERJ integra este trabalho.

## Responsáveis e limites

Documentação corrente em `00_nucleo/documentacao`, conforme o mapa arquitetural e o ADR-0008. Coleta e exploração usam as fachadas existentes de L1, a CLI L2, os leitores L3 e o executável composto em L4. Não alterar capacidades, formatos, critérios, organização física do código ou evidência histórica. A worktree contém mudanças anteriores, que devem ser preservadas.

A investigação distingue material do alvo, referências externas, contexto produzido pelo Bateia e interpretação desta sessão pela IA. Seleção lexical não comprova relevância; proposta não comprova novidade nem efeito. Commits, hashes, comandos, falhas, localizadores e limitações devem acompanhar a bancada.

## Aceite

- README curto com navegação por tarefa; exemplos anteriores preservados em referência documental com links corrigidos.
- Guia do primeiro uso diferencia discover, expand, explore, analyze/compare e decisão de reúso.
- Caso XERJ reproduzível no Lab, com coleta real e execução de explore; descoberta online e seus limites registrados quando exercitada.
- Sugestões para o XERJ ligadas às fontes, com condição de aplicação e experimento de rejeição; capacidades já existentes não apresentadas como ausentes.
- Executar lint, testes, vet e build; relatar bloqueios e falhas sem reclassificá-los como aprovação.

Scripts, fontes congeladas, resultados e logs pertencem a `lab/bancada/insights/xerj-2026-11-08`. Esta fatia não demonstra organização integral do brownfield ou utilidade independente das propostas.

## Resultado do recorte

O README de entrada foi reduzido a 64 linhas; os exemplos anteriores ficam na referência de comandos. O guia de investigação e os índices documental e do Lab conectam tarefas às entradas atuais. Nenhum pacote produtivo foi movido ou alterado neste recorte.

A [bancada XERJ](../../lab/bancada/insights/xerj-2026-11-08/README.md) registra descoberta real, coleta limitada de 17 arquivos do alvo e exploração com 16 cartões. A interpretação da IA propõe vínculo da avaliação G7 à versão, recibo dos escopos verificados e contexto de aplicabilidade/contraexemplos. Nenhuma proposta foi implementada ou teve efeito medido no XERJ.

Replay offline idêntico e vínculos de hashes/localizadores verificados. O gate canônico final passou em lint, auditorias, testes, vet, build e diff-check. A primeira tentativa falhou por limitações de armazenamento temporário e loopback; seus logs foram preservados e a execução final usou cache temporário, paralelismo reduzido e loopback permitido. Os [registros de verificação](../../lab/bancada/insights/xerj-2026-11-08/verification/verification.json) mantêm comandos e códigos de saída; isso não valida as hipóteses sobre o XERJ.
