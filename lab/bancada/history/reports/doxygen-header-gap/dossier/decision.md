# Decisão de reúso: INCONCLUSIVE

Fonte canônica: [reuse-decision.json](reuse-decision.json)

## Razões

- `DEC-INCONCLUSIVE-001` : evidência obrigatória ou gate eliminatório permanece inconclusivo
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: capacidade não comprovada: document assignment as intentionally leaving container/observer membership untouched, deferring identity transfer to an explicit swap; coleta incompleta ou com falha; ausência não pode ser concluída
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: capacidade não comprovada: explicitly unregister an object's own address from every external observer list it is currently registered on before adopting new state; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: capacidade não comprovada: guard copy assignment against self-assignment before mutating any state; coleta incompleta ou com falha; ausência não pode ser concluída
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: capacidade não comprovada: use the copy-and-swap idiom so the old state's own destructor performs cleanup; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: licença não verificada: localizar documento de licença com autoridade no projeto
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: manutenção não verificada: obter metadados datados de atividade e arquivamento
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: operação local não verificada: examinar requisitos e instruções de execução
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: restrição não verificada: C++11 as the minimum supported language standard; restrição sem regra determinística específica
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: restrição não verificada: must not leave a stale pointer registered on any external list after assignment; restrição sem regra determinística específica
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: restrição não verificada: non_goal: changing the meaning of copy construction; coleta incompleta
- `DEC-EVIDENCE-0040` local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631: restrição não verificada: non_goal: making the type non-copyable; coleta incompleta
- `DEC-CATALOG-0071` : a cobertura do catálogo é ausente ou apenas hipotética para parte dos requisitos; isso não descreve lacuna de nenhum candidato

## Capacidades ainda não comprovadas

- document assignment as intentionally leaving container/observer membership untouched, deferring identity transfer to an explicit swap
- explicitly unregister an object's own address from every external observer list it is currently registered on before adopting new state
- guard copy assignment against self-assignment before mutating any state
- use the copy-and-swap idiom so the old state's own destructor performs cleanup

## Plano de investigação

Somente candidatos, tentativas e evidências registrados; não demonstra busca global completa nem integração verificada.

Encerramento: `RECORDED_SEARCH_FINISHED`. Tempos de busca são telemetria do Bateia, não benchmarks dos candidatos.

- **Que evidência responde a: explicitly unregister an object's own address from every external observer list it is currently registered on before adopting new state?** (`req-capability-2a9d52bc4dc88bda`) — INFERENCE/PARTIAL. Lacuna: símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Verificar se os indícios atendem ao requisito original e procurar contrapontos antes de recomendar reúso. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list_hook.hpp:151` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list_hook.hpp:274` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list_hook.hpp:157` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list_hook.hpp:280` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist_hook.hpp:154` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist_hook.hpp:283` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist_hook.hpp:160` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist_hook.hpp:289` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/nonhook_node.hpp:45` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/nonhook_node.hpp:51` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/nonhook_node.hpp:57` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/nonhook_node.hpp:67` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:437` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:438` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:444` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:445` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:451` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:452` (test_implementation).
- **Que evidência responde a: use the copy-and-swap idiom so the old state's own destructor performs cleanup?** (`req-capability-4007800bc6820bea`) — INFERENCE/PARTIAL. Lacuna: símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Verificar se os indícios atendem ao requisito original e procurar contrapontos antes de recomendar reúso. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:1045` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:1047` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:1053` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:1418` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:1419` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:232` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:238` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:560` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/list.hpp:563` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:1277` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:1472` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:2017` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:2018` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:352` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:358` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:721` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `include/boost/intrusive/slist.hpp:729` (symbol).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:363` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:391` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:399` (test_implementation).
  - Fonte: `local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631` em `test/slist_test.cpp:406` (test_implementation).
- **Que evidência responde a: document assignment as intentionally leaving container/observer membership untouched, deferring identity transfer to an explicit swap?** (`req-capability-6b9c235339c79256`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: coleta incompleta ou com falha; ausência não pode ser concluída Próxima ação: Completar a coleta indicada nas falhas do assessment antes de avaliar a capacidade; não inferir ausência. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: guard copy assignment against self-assignment before mutating any state?** (`req-capability-a638d285b0fe11f0`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: coleta incompleta ou com falha; ausência não pode ser concluída Próxima ação: Completar a coleta indicada nas falhas do assessment antes de avaliar a capacidade; não inferir ausência. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: C++11 as the minimum supported language standard?** (`req-constraint-da9653ee9141cc8c`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: restrição sem regra determinística específica Próxima ação: Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: must not leave a stale pointer registered on any external list after assignment?** (`req-constraint-e2c530f308726f35`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: restrição sem regra determinística específica Próxima ação: Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: changing the meaning of copy construction?** (`req-non_goal-03ee38507131a5cb`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio. Aceitação: Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento.
- **Que evidência responde a: making the type non-copyable?** (`req-non_goal-8ad5701813765558`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio. Aceitação: Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento.
- **Que evidência responde a: cmake?** (`req-preference-4682918564fa2db4`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida. Aceitação: Registrar compatibilidade observada; preferência não é gate eliminatório.
- **Que evidência responde a: c++?** (`req-preference-5d050596ec3e3a4f`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida. Aceitação: Registrar compatibilidade observada; preferência não é gate eliminatório.
- **Que evidência responde a: Correctly implement copy assignment for a type that registers itself, by its own address, as an observer/auditor on external containers reachable from its current state (like SoPath registering itself on each node's SoChildList so scene-graph edits notify it), so that assigning a new route unregisters the old observation exactly once, registers the new one exactly once, and self-assignment is a safe no-op.?** (`req-problem-eab51dc68a7bb30c`) — CATALOG_GAP/INCONCLUSIVE. Lacuna: o catálogo não possui regra declarada única e segura para este requisito Próxima ação: Definir uma verificação delimitada e segura; a cobertura do catálogo não descreve ausência no candidato. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.

Esforço: INCONCLUSIVE. Método: heuristic_component_gap_count. Não mede custo nem viabilidade; requer ensaio de integração

## Recuperação em local/boost-intrusive-ddb60e30@81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631

##### Interpretação lexical da consulta

Símbolos explícitos são pedidos de busca; restrições não são afirmações sobre o projeto.

- Original: "Correctly implement copy assignment for a type that registers itself, by its own address, as an observer/auditor on external containers reachable from its current state (like SoPath registering itself on each node's SoChildList so scene-graph edits notify it), so that assigning a new route unregisters the old observation exactly once, registers the new one exactly once, and self-assignment is a safe no-op.". Modo: `mixed`. Símbolos (2/2 exibidos): ["SoChildList" "SoPath"].
- Original: "document assignment as intentionally leaving container/observer membership untouched, deferring identity transfer to an explicit swap". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "explicitly unregister an object's own address from every external observer list it is currently registered on before adopting new state". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "guard copy assignment against self-assignment before mutating any state". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "use the copy-and-swap idiom so the old state's own destructor performs cleanup". Modo: `natural`. Símbolos (0/0 exibidos): [].

##### Indícios e lacunas por requisito

Correspondência lexical é inferência; requisito satisfeito e comportamento continuam não verificados. Termos sem correspondência não demonstram capacidades ausentes.

- Requisito: "Correctly implement copy assignment for a type that registers itself, by its own address, as an observer/auditor on external containers reachable from its current state (like SoPath registering itself on each node's SoChildList so scene-graph edits notify it), so that assigning a new route unregisters the old observation exactly once, registers the new one exactly once, and self-assignment is a safe no-op.". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/40 exibidos): ["address" "assigning" "assignment" "auditor" "child" "containers" "copy" "correctly" "current" "each" "edits" "exactly"].

- Requisito: "document assignment as intentionally leaving container/observer membership untouched, deferring identity transfer to an explicit swap". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/13 exibidos): ["assignment" "container" "deferring" "document" "explicit" "identity" "intentionally" "leaving" "membership" "observer" "swap" "transfer"].

- Requisito: "explicitly unregister an object's own address from every external observer list it is currently registered on before adopting new state". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/15 exibidos): ["address" "adopting" "before" "currently" "every" "explicitly" "external" "list" "new" "object" "observer" "own"].

- Requisito: "guard copy assignment against self-assignment before mutating any state". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (9/9 exibidos): ["against" "any" "assignment" "before" "copy" "guard" "mutating" "self" "state"].

- Requisito: "use the copy-and-swap idiom so the old state's own destructor performs cleanup". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (10/10 exibidos): ["cleanup" "copy" "destructor" "idiom" "old" "own" "performs" "so" "state" "swap"].

#### Hipóteses de ação por candidato

Rótulos e chamadas não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade. A ordem dos candidatos prioriza leitura, não confiança semântica.

- Consulta <code>&#34;Correctly implement copy assignment for a type that registers itself, by its own address, as an observer/auditor on external containers reachable from its current state (like SoPath registering itself on each node&#39;s SoChildList so scene-graph edits notify it), so that assigning a new route unregisters the old observation exactly once, registers the new one exactly once, and self-assignment is a safe no-op.&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Modo da consulta não natural; ação não avaliada. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;document assignment as intentionally leaving container/observer membership untouched, deferring identity transfer to an explicit swap&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;explicitly unregister an object&#39;s own address from every external observer list it is currently registered on before adopting new state&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;guard copy assignment against self-assignment before mutating any state&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;use the copy-and-swap idiom so the old state&#39;s own destructor performs cleanup&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.


## Evidências

Consulte [assessment-report.json](assessment-report.json) e [evidence/bundles.json](evidence/bundles.json).
