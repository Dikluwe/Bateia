# 0136 — Exploração para ampliar o repertório da IA

**Data:** 2026-10-28. **Contrato:** [passo 0136](../passos/0136-preparar-exploracao-criativa-para-ia.md). **Decisão:** [ADR-0017](../ADR/ADR-0017-contexto-de-exploracao-para-ia.md). **Evidências:** [bancada](../../lab/bancada/insights/0136/README.md).

## Direção e entrega

O usuário pediu menos trabalho manual e explicitou que o Bateia deve auxiliar a IA a sair das abordagens habituais. O novo `idea insights explore` recebe uma pergunta e bundles já coletados; familiaridade e limite de trechos são opcionais. Não exige alvo, classes, assessments ou observations preparados à mão.

O comando seleciona cartões de leitura automaticamente, com fonte imutável, hash, linhas e trecho. Prioriza repositórios/tipos ainda não representados, matches literais da pergunta e menos correspondências familiares; variedade de palavras desempata janelas e evita preferir cabeçalhos repetitivos. Não repete trechos de bytes idênticos. Os diagnósticos conservam limites de seleção/prefixo, itens ignorados, fontes e incompletude.

As cinco tarefas para IA propõem alternativas, inversão de premissa, analogia entre domínios, recombinação e contraexemplos. Referenciam cartões existentes e exigem separar observação de hipótese, condições, custos e experimento com critério de rejeição. Recombinação prioriza duas fontes de repositórios distintos quando disponíveis; analogia conserva a lacuna de domínio semântico não classificado. Trechos são dados não confiáveis, nunca instruções.

## Responsabilidades e compatibilidade

`insight_exploration.go` em L1 possui normalização, identidade, tarefas e direções; `insight_passages.go` possui leitura/seleção limitada. L2 adapta flags, lê bundles estritos e publica pelas proteções existentes. Não há chamada de IA, busca, execução de fonte ou novo runtime; as direções de pesquisa são PLANNED_NOT_EXECUTED. Descoberta/coleta continuam disponíveis nos comandos existentes.

Analyze v1/v2 e compare mantêm seus formatos/IDs. A exploração não depende do catálogo fechado de mecanismos e não apresenta seus trechos como classificações. O material do corpus é informado pelo chamador; selecionar diversidade dentro desse conjunto não equivale a descobrir diversidade no ecossistema.

## Ensaio conhecido

O caso reúne bundles congelados de 0106/0122, com 82 itens e 616.356 bytes de conteúdo de quatro repositórios. O pacote selecionou oito cartões, 7.883 bytes de trechos, de três repositórios e quatro tipos de evidência. Nenhuma avaliação ou observação manual foi necessária para essa preparação. A quarta fonte não forneceu candidato selecionado e três bundles estão incompletos, o que permanece explícito. Não houve nova coleta ou execução desses projetos.

Um dos cartões de conc contém o comentário de um teste sobre uma corrida e a espera usada para torná-la improvável. Outro documenta a interface de engines de arbiter. São pistas para perguntar sobre premissas de ordenação, invariantes e diferenças de implementação, em vez de apenas listar ferramentas semelhantes. Essa leitura é autoral; o comando não deduziu essas interpretações nem comprovou benefício de uma ideia no Bateia.

O controle somente com conc funciona sem avaliações. A pergunta em português não encontrou termos no corpus inglês e o pacote registra QUESTION_TERMS_NOT_MATCHED. A seleção não traduz palavras ou reconhece sinônimos. Sem corpus, o estado é SEARCH_NEEDED com tarefas/pesquisas e nenhuma citação inventada. Com cartões, CONTEXT_PREPARED significa contexto pronto para leitura, mantendo novidade NOT_ASSESSED e efeito UNVERIFIED.

## Verificação

`make verify` passou integralmente: lint sem violações, onze controles do linter, auditorias sobre 363 arquivos produtivos, suíte raiz, vet, build e diff-check. O [recibo](../../lab/bancada/insights/0136/final/verification.json) preserva comandos e códigos de saída. As fixtures HTTP usaram loopback, e assinatura Git foi desativada somente no processo do gate.

Os testes focados cobriram integridade/quotas, janelas/bytes/UTF-8/linhas, seleção por repositórios/tipos, revisões do mesmo projeto, clones, familiaridade lexical, prefixos, falta de matches/corpus, referências, determinismo/imutabilidade/causalidade e estados conservadores. A CLI foi exercitada com flags diretas/repetidas, JSON estrito, aliases, segredo nos trechos, short writes e preservação de saída em falha.

O verificador do caso reproduz o pacote e seus controles, confere hashes/localizadores/referências, reordenação/prefixos, adulteração e corpus vazio. Os verificadores anteriores de analyze v1, v2 e compare também passaram com saídas idênticas byte a byte. Comparados à baseline, todos os testes Go anteriores e 23 artefatos históricos inspecionados permaneceram iguais; dos 816 arquivos Go anteriores, somente ajuda/despacho da CLI foram alterados. Os oito arquivos novos e dois alterados possuem [linhagem](../../lab/bancada/insights/0136/lineage.json).

## Limites

A entrega reduz a curadoria para preparar contexto, sem gerar conclusões por modelo ou medir criatividade/utilidade. Diversidade lexical/estrutural pode selecionar material irrelevante; corpo e prefixos são limitados. Associação externa/semântica não é autenticada. As tarefas orientam a IA, mas não garantem que ela encontre uma alternativa útil. O corpus conhecido não é holdout; medir qualidade das propostas e efeito no alvo exige avaliação própria. Gates selados com módulos separados não são promovidos pela suíte raiz.
