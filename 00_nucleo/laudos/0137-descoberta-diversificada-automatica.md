# 0137 — Descoberta diversificada automática

**Data:** 2026-10-28. **Contrato:** [passo 0137](../passos/0137-descobrir-referencias-diversificadas-automaticamente.md). **Decisão:** [ADR-0018](../ADR/ADR-0018-descoberta-diversificada-e-coleta-limitada-para-insights.md). **Evidências:** [bancada](../../lab/bancada/insights/0137/README.md).

## Entrega

`idea insights discover` conecta pergunta e familiaridade a cinco rotas independentes de pesquisa, seleção de referências, coleta de README fixado e contexto para IA. O executável usa o catálogo e o cliente de conteúdo GitHub existentes. O usuário não precisa montar IdeaSpec, avaliações, observações ou plano de coleta para essa primeira investigação.

As rotas direta, ampliada, alternativas, trade-offs e falhas são tentadas mesmo após sucesso inicial, dentro dos budgets. Candidatos repetidos conservam todas as query IDs; revisões conflitantes são registradas e não lidas. A seleção prioriza ganho de consultas representadas, autores novos e menos matches familiares, sem usar estrelas. Isso diversifica rotas/autores; não classifica relevância, domínio semântico ou novidade.

O fluxo lê somente um README regular de raiz de cada selecionado, em SHA imutável, conferindo ref/path/cardinalidade/UTF-8 e admissão de até 32 KiB. Não trunca, clona, instala ou executa candidatos. O bundle permanece incompleto pelo escopo ROOT_README_ONLY e registra somente a cobertura documental examinada. Explore produz os cartões e tarefas no próprio relatório, com novidade NOT_ASSESSED e efeito UNVERIFIED.

## Fronteiras

L1 possui contratos e owners separados de plano, execução, união/seleção e política de coleta. Recebe ProjectCatalog, porta estreita de árvore/extração e runtime explícito no contexto. L3 mantém HTTP e controles do cliente existente. L2 possui flags/sinais/publicação protegida de resultado completo ou parcial. L4 injeta as portas no ramo online específico, preservando analyze/compare/explore offline. Não há novo runtime durável ou alteração da política de reúso.

## Ensaio público

O binário real executou uma pergunta curta de backpressure/queue/cancellation, sem token GitHub e em diretório sem `.env`. O budget foi reduzido a um resultado por consulta, duas referências, quatro trechos, 45 s totais e 6 s por operação. As cinco buscas concluíram e produziram quatro candidatos únicos. Foram selecionados autores distintos:

- AvdLee/Swift-Concurrency-Agent-Skill, SHA `d5770817d2622e1585b1f7eaebc791a9cb0959c8`: README de 9.949 bytes admitido, com cartão documental sobre orientações de concorrência Swift.
- ReactiveX/RxJava, SHA `b0e8737b956f05cebca961899193c80104c94941`: README rejeitado por exceder o teto de admissão.

O resultado ficou PARTIAL com um cartão e cinco tarefas para IA; preservou candidatos, origens e limite de conteúdo. O documento de skill é tratado como dado não confiável, sem instalação ou seguimento de suas instruções. A coleta comprova o vínculo aos bytes recebidos, não utilidade ou novidade das ideias que a IA possa formular.

## Reprodução e correção encontrada

Na conferência offline, a cobertura vazia do bundle legado alternava sua representação entre lista vazia e null após round-trip JSON, mudando a identidade do contexto apesar de cartões/tarefas iguais. O novo produtor passou a registrar a cobertura documental real do README examinado, mantendo bundle incompleto e sem avaliar capacidade. O marshaller e os contratos antigos não foram alterados.

Um teste de regressão com conteúdo HTML confere IDs após serializar/reler bundles. O ensaio público inicial e o gate anterior foram preservados em pastas próprias; o ensaio final foi repetido após a correção. O [verificador offline](../../lab/bancada/insights/0137/verify.py) confere identidades do plano, origens, hashes/localizadores e igualdade integral do contexto reconstruído a partir dos bundles publicados. Não afirma replay byte a byte de uma nova pesquisa online, cujos resultados/metadados/durações podem variar.

## Verificação

Os testes focados passaram: normalização/budgets, rotas independentes e deadlines, cobertura/autores/popularidade, duplicação/conflito de revisão, falhas sanitizadas, cancelamento preservando parcial, seleção/limites/integridade de README e replay do contexto. A CLI verifica output obrigatório e flags antes da rede, configuração explícita, saída protegida em falha total e bloqueio de segredos. A integração HTTP L3 usa cliente real contra servidor local e quota parcial; o teste L4 exercita o despacho real em subprocesso sem rede.

Os replays históricos de 0106/0134/0135/0136 passaram. Comparada à baseline, a mudança preservou todos os testes Go preexistentes e 36 artefatos históricos inspecionados. Dos 824 arquivos Go anteriores, apenas ajuda/despacho da CLI e composição main foram alterados. Quatorze arquivos novos e três alterados têm prompts espelhados e [linhagem](../../lab/bancada/insights/0137/lineage.json).

`make verify` passou integralmente após a correção: lint sem violações, onze controles do linter, auditorias sobre 369 arquivos produtivos, suíte raiz, vet, build e diff-check. O [recibo final](../../lab/bancada/insights/0137/final/verification.json) registra os comandos/códigos. O binário final passou no verificador do snapshot público e nos quatro verificadores históricos. As fixtures HTTP usaram loopback e assinatura Git foi desativada somente no processo do gate.

## Limites

Este primeiro fluxo usa GitHub, sem integrar implicitamente os outros conectores. As consultas são heurísticas lexicais sem tradução ou descoberta semântica de sinônimos. Cobertura estrutural, presença no catálogo e leitura de README não comprovam novidade, relevância ou comportamento. Documentos grandes, revisões opacas, quotas e prazos podem deixar poucas referências; o resultado conserva as lacunas. Não foi medida utilidade dos insights, e gates selados com módulos próprios não são promovidos pela suíte raiz.
