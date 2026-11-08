# 0144 — Continuar o foco e selecionar vizinhos

**Contrato ex-ante. Data:** 2026-10-28. A1-INV sob Tekt. L1 possui foco, aliases e seleção; L2 transporta --all-trails; leitores/composição não mudam. Preservar snapshots/IDs históricos, comportamento de discover e evidência anterior.

## Continuação

Expand herda o foco de input.neighborhood.focus quando não houver foco novo. --all-trails solicita explicitamente a fronteira geral e não pode ser combinado com foco novo. Entrada discover sem neighborhood continua geral como antes. Foco desconhecido/descartado continua inválido. As leituras continuam em um salto, com deadlines/quotas de0143 e zero catálogo.

## Aliases conservadores

Somente na expansão, reconhecer a relação X.git→X quando ambas as identidades canônicas já existem na memória, com mesmo conector e host exato, e X não termina também em .git. Não normalizar nomes, maiúsculas, hosts, renomes ou sufixos arbitrários; não inferir existência remota. Não alterar identidade, testemunhos, selos ou estados armazenados, nem regras globais de leitura/descoberta.

Se ambas estão pendentes, somente X entra na fila; X.git aparece como DUPLICATE_ALIAS. Se X está examinado/descartado, o alias pendente é bloqueado como ALIAS_ALREADY_EXAMINED/ALIAS_DISCARDED. Se X está pendente mas o alias está examinado, bloquear X salvo reabertura explícita de X com motivo. Se X está pendente e alias descartado, bloquear com ALIAS_STATE_CONFLICT: não resolver decisões incompatíveis silenciosamente. Reabertura/decisões continuam no feedback existente. Foco em alias conhecido passa ao membro X; rejeitar X descartado. Estados e counters de destinos suprimidos permanecem inalterados.

A travessia usa a relação de alias para alcançar citações do mesmo destino sem gerar ciclos/releituras. Descartes interrompem a travessia. Permanecem as limitações do primeiro testemunho e das quotas de memória/extração.

## Seleção observável

No expand, ordenar por leitor elegível no sítio selecionado, menor distância em citações ao foco, menos tentativas anteriores, geração mais antiga, ordem das referências atuais e ID. Sem foco, não atribuir distância: preservar os desempates de0143. Foco pendente elegível tem distância0. Distância mede saltos conhecidos, não relevância, novidade ou benefício. Discover conserva sua ordem anterior.

Extensões opcionais em neighborhood: selection_policy e neighbors. Cada vizinho registra entry-id, identidade, distância opcional, tentativas anteriores, selection_status e related_entry_id quando suprimido por alias. Projeção é o snapshot anterior às leituras; estado efetivo atual continua em followups/memória. Listar fila na ordem de execução e suprimidos depois, deterministicamente; frontier_entry_ids contém somente fila admitida, incluindo omissões de leitor/fonte com diagnóstico em followups. Vizinhos suprimidos não consomem tentativas. Ausência desses campos em relatórios antigos preserva sua serialização e selo.

## Verificação

Testar herança/override/all-trails e conflitos antes de efeitos, alias examinado/descartado/pendente/conflitante, reabertura explícita, foco em alias, ciclos e separação de host/forge/subgrupo. Testar distância0/1/2 com tentativas/idade adversariais, ordem sem foco e ausência de catálogo. JSON antigo deve conservar bytes; testes CLI e todos os gates. Controle usando a memória pública anterior demonstra supressão de frontend.git sem fingir leitura. Replays históricos e hashes em lab/bancada/insights/0144; nenhuma execução de candidatos ou promoção de hipóteses.
