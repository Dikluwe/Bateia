# 0146 — Consolidar prompts de exploração

**Data:** 2026-10-29. **Contrato ex-ante desta mudança documental.**

## Intenção e recorte

Tornar os prompts de expansão, memória e aprofundamento suficientes para orientar mudanças e testes, sem exigir a reconstrução cronológica de0142–0145. Consolidar contratos vigentes com entradas, saídas, limites, erros, responsáveis e critérios observáveis. A consolidação usa especificações anteriores e inspeção do código existente: não constitui autoria ex-ante desse código.

## Entrega e fronteiras

Criar contrato compartilhado e matriz de aceitação em L0, substituir os 28 prompts proprietários deste recorte por obrigações locais e referências precisas, e publicar um guia de escrita. Preservar snapshots dos textos anteriores no Lab. Não modificar código Go, schemas, regras canônicas, testes existentes ou evidência histórica. Não transformar defeitos observados em regras desejadas; registrar divergências e lacunas de cobertura explicitamente.

O contrato compartilhado possui as regras transversais da capacidade; cada prompt espelhado possui sua realização e fronteira. A0 continua em .agentrules e nos ADRs; contratos funcionais não a substituem. Nenhuma regra nova de produto é materializada nesta fatia.

## Aceitação

Cada prompt declara proprietário, camada, interface, comportamento, limites, dependências, casos observáveis e origem temporal. Cada caso da matriz distingue obrigação de evidência existente; teste citado não significa cobertura completa. Referências locais e âncoras devem resolver. Evidências e executáveis ficam em lab/bancada/prompts/0146. Verificar preservação de todos os Go existentes e das bancadas0142–0145 em relação ao início desta mudança; executar lint, testes, vet e build pelo gate canônico. A aprovação estrutural não atesta utilidade dos insights nem completude semântica dos prompts.
