# ADR-0013 — Regras de alternativas e retomada em L1

**Status:** ACEITO PARA O RECORTE 0130  
**Data:** 2026-10-28  
**Governante:** Tekt e ADR-0008; preserva as portas do ADR-0012.

## Decisão

`alternative_target.go` possui a interpretação dos hints de identidade para a projeção de alternativas locais. Recebe somente um leitor restrito de identidades; o comando liga o lookup congelado da captura, ou o adaptador de metadata existente. Nenhum caminho físico ou cliente concreto entra no núcleo. Hints mantêm sua autoridade consultiva.

`screening_model_policy.go` possui a admissão semântica do perfil opt-in, separada da desserialização estrita em L2 e da admissão de outros perfis de screening. Mover regras não altera pesos, thresholds ou autoridade de evidência.

`evidence_checkpoint.go` possui a política de elegibilidade e verificação de checkpoint por uma porta restrita que retorna SHA-256 de um artefato lógico. O adaptador CLI resolve o caminho do dossiê e usa as leituras já implementadas pelo host L3. A política controla quais artefatos podem ser consultados e a ordem; não recebe raiz física nem faz I/O concreto.

## Consequências

L2 delega as três decisões a L1 e conserva transporte, apresentação e integração de arquivos. As funções privadas históricas necessárias aos consumidores podem permanecer como delegações. Formatos, IDs, hashes, versões, códigos de estado, critérios e efeitos de leitura permanecem os atuais. Policy/target podem mudar sem invalidar evidência causalmente compatível; a decisão é recalculada.

As restrições antigas do prompt CLI contra alteração de L1 no marco do modelo opt-in são substituídas apenas quanto ao ownership dessas regras neste recorte. Não há autorização para alterar a política de calibração, ranking, retomada ou evidência.
