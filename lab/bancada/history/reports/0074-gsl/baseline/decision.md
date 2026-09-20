# Decisão de reúso: INCONCLUSIVE

Fonte canônica: [reuse-decision.json](reuse-decision.json)

## Razões

- `DEC-INCONCLUSIVE-001` : evidência obrigatória ou gate eliminatório permanece inconclusivo
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: capacidade não comprovada: avoid signed-integer overflow when validating a start offset combined with a count against a total length; testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: capacidade não comprovada: clamp an excessive requested count to the remaining length instead of treating it as an error; testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: capacidade não comprovada: express a sub-range using two absolute bounds instead of an offset-plus-count pair; testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: capacidade não comprovada: reject an out-of-range start position unconditionally, not only in debug builds; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: licença não verificada: localizar documento de licença com autoridade no projeto
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: manutenção não verificada: obter metadados datados de atividade e arquivamento
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: operação local não verificada: examinar requisitos e instruções de execução
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: restrição não verificada: C++11 as the minimum supported language standard; restrição sem regra determinística específica
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: restrição não verificada: must not allocate more than the resulting sub-range needs; restrição sem regra determinística específica
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: restrição não verificada: non_goal: changing the container's own element storage; coleta incompleta
- `DEC-EVIDENCE-0040` local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591: restrição não verificada: non_goal: supporting negative indices as counting from the end; coleta incompleta
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: capacidade não comprovada: avoid signed-integer overflow when validating a start offset combined with a count against a total length; testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: capacidade não comprovada: clamp an excessive requested count to the remaining length instead of treating it as an error; testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: capacidade não comprovada: express a sub-range using two absolute bounds instead of an offset-plus-count pair; testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical)
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: capacidade não comprovada: reject an out-of-range start position unconditionally, not only in debug builds; sinais fracos impedem concluir presença ou ausência
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: manutenção não verificada: obter metadados datados de atividade e arquivamento
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: operação local não verificada: examinar requisitos e instruções de execução
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: restrição não verificada: C++11 as the minimum supported language standard; restrição sem regra determinística específica
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: restrição não verificada: must not allocate more than the resulting sub-range needs; restrição sem regra determinística específica
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: restrição não verificada: non_goal: changing the container's own element storage; ausência sem testemunho admissível do escopo examinado; non-goal não observado no escopo coletado
- `DEC-EVIDENCE-0040` local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72: restrição não verificada: non_goal: supporting negative indices as counting from the end; ausência sem testemunho admissível do escopo examinado; non-goal não observado no escopo coletado
- `DEC-CATALOG-0071` : a cobertura do catálogo é ausente ou apenas hipotética para parte dos requisitos; isso não descreve lacuna de nenhum candidato

## Capacidades ainda não comprovadas

- avoid signed-integer overflow when validating a start offset combined with a count against a total length
- clamp an excessive requested count to the remaining length instead of treating it as an error
- express a sub-range using two absolute bounds instead of an offset-plus-count pair
- reject an out-of-range start position unconditionally, not only in debug builds

## Plano de investigação

Somente candidatos, tentativas e evidências registrados; não demonstra busca global completa nem integração verificada.

Encerramento: `RECORDED_SEARCH_FINISHED`. Tempos de busca são telemetria do Bateia, não benchmarks dos candidatos.

- **Que evidência responde a: express a sub-range using two absolute bounds instead of an offset-plus-count pair?** (`req-capability-61ba9f8a86bb5068`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: express a sub-range using two absolute bounds instead of an offset-plus-count pair?** (`req-capability-61ba9f8a86bb5068`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: clamp an excessive requested count to the remaining length instead of treating it as an error?** (`req-capability-6de26864b1d0fc5b`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: clamp an excessive requested count to the remaining length instead of treating it as an error?** (`req-capability-6de26864b1d0fc5b`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: avoid signed-integer overflow when validating a start offset combined with a count against a total length?** (`req-capability-c81b145ee0314a93`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: avoid signed-integer overflow when validating a start offset combined with a count against a total length?** (`req-capability-c81b145ee0314a93`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: testemunho sem identidade admissível; símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Localizar implementação, limitações e ensaio pertinente; registrar arquivos e intervalos efetivamente examinados. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: reject an out-of-range start position unconditionally, not only in debug builds?** (`req-capability-f17ebfece51e2855`) — INFERENCE/PARTIAL. Lacuna: símbolo, teste ou indício localizado; comportamento não verificado (avaliação lexical) Próxima ação: Verificar se os indícios atendem ao requisito original e procurar contrapontos antes de recomendar reúso. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:1444` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:1649` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:1921` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:1988` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:2071` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:2181` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:2208` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:2256` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:3165` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:3282` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:3320` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:3321` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:3386` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:3426` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:386` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:569` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:586` (symbol).
  - Fonte: `local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591` em `bits/basic_string.h:604` (symbol).
- **Que evidência responde a: reject an out-of-range start position unconditionally, not only in debug builds?** (`req-capability-f17ebfece51e2855`) — INFERENCE/INCONCLUSIVE. Lacuna: sinais fracos impedem concluir presença ou ausência Próxima ação: Verificar se os indícios atendem ao requisito original e procurar contrapontos antes de recomendar reúso. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
  - Fonte: `local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72` em `README.md:46` (documentation).
- **Que evidência responde a: C++11 as the minimum supported language standard?** (`req-constraint-da9653ee9141cc8c`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: restrição sem regra determinística específica Próxima ação: Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: C++11 as the minimum supported language standard?** (`req-constraint-da9653ee9141cc8c`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: restrição sem regra determinística específica Próxima ação: Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: must not allocate more than the resulting sub-range needs?** (`req-constraint-f0bf6cb3dd3f08c7`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: restrição sem regra determinística específica Próxima ação: Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: must not allocate more than the resulting sub-range needs?** (`req-constraint-f0bf6cb3dd3f08c7`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: restrição sem regra determinística específica Próxima ação: Definir uma verificação delimitada para este requisito e os testemunhos necessários; o catálogo atual não possui detector aplicável. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: supporting negative indices as counting from the end?** (`req-non_goal-8b3e74507d6e1a57`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio. Aceitação: Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento.
- **Que evidência responde a: supporting negative indices as counting from the end?** (`req-non_goal-8b3e74507d6e1a57`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio. Aceitação: Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento.
- **Que evidência responde a: changing the container's own element storage?** (`req-non_goal-e6c41e3d65a6397b`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio. Aceitação: Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento.
- **Que evidência responde a: changing the container's own element storage?** (`req-non_goal-e6c41e3d65a6397b`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Investigar explicitamente o comportamento excluído e sua negação; não inferir ausência por silêncio. Aceitação: Examinar presença e negação do comportamento excluído; ausência de menção não demonstra ausência do comportamento.
- **Que evidência responde a: cmake?** (`req-preference-4682918564fa2db4`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida. Aceitação: Registrar compatibilidade observada; preferência não é gate eliminatório.
- **Que evidência responde a: cmake?** (`req-preference-4682918564fa2db4`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida. Aceitação: Registrar compatibilidade observada; preferência não é gate eliminatório.
- **Que evidência responde a: c++?** (`req-preference-5d050596ec3e3a4f`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida. Aceitação: Registrar compatibilidade observada; preferência não é gate eliminatório.
- **Que evidência responde a: c++?** (`req-preference-5d050596ec3e3a4f`) — NOT_VERIFIED/INCONCLUSIVE. Lacuna: Evidência admissível insuficiente para responder à pergunta. Próxima ação: Comparar a preferência de stack com manifesto e interfaces; não excluir por preferência desconhecida. Aceitação: Registrar compatibilidade observada; preferência não é gate eliminatório.
- **Que evidência responde a: Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.?** (`req-problem-c46287951f7b87df`) — CATALOG_GAP/INCONCLUSIVE. Lacuna: o catálogo não possui regra declarada única e segura para este requisito Próxima ação: Definir uma verificação delimitada e segura; a cobertura do catálogo não descreve ausência no candidato. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.
- **Que evidência responde a: Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.?** (`req-problem-c46287951f7b87df`) — CATALOG_GAP/INCONCLUSIVE. Lacuna: o catálogo não possui regra declarada única e segura para este requisito Próxima ação: Definir uma verificação delimitada e segura; a cobertura do catálogo não descreve ausência no candidato. Aceitação: Localizar implementação e testemunho pertinente ao requisito, com identidade e limites; validar o comportamento sob condições declaradas.

Esforço: INCONCLUSIVE. Método: heuristic_component_gap_count. Não mede custo nem viabilidade; requer ensaio de integração

## Recuperação em local/13-0ea945f7@3b858f322d7ad7f2c92f017e630cd187b5db6770f33070ec0c9ab0bbba62e591

##### Interpretação lexical da consulta

Símbolos explícitos são pedidos de busca; restrições não são afirmações sobre o projeto.

Restrições registradas: 2. Escopo e cumprimento não verificados, inclusive nas linhas não exibidas abaixo.

- Marcador "without", texto original "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.", estado `not_verified`.
- Marcador "not", texto original "reject an out-of-range start position unconditionally, not only in debug builds", estado `not_verified`.

- Original: "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.". Modo: `mixed`. Símbolos (1/1 exibidos): ["SoPath"].
- Original: "avoid signed-integer overflow when validating a start offset combined with a count against a total length". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "clamp an excessive requested count to the remaining length instead of treating it as an error". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "express a sub-range using two absolute bounds instead of an offset-plus-count pair". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "reject an out-of-range start position unconditionally, not only in debug builds". Modo: `natural`. Símbolos (0/0 exibidos): [].

##### Indícios e lacunas por requisito

Correspondência lexical é inferência; requisito satisfeito e comportamento continuam não verificados. Termos sem correspondência não demonstram capacidades ausentes.

- Requisito: "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Restrições preservadas; escopo e cumprimento da negação não foram verificados. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Restrições: escopo e cumprimento da negação não verificados.
  Marcador "without", texto original "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.", estado `not_verified`.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/26 exibidos): ["arithmetic" "copy" "correctly" "count" "empty" "excessive" "handling" "indexed" "integer" "like" "model" "needed"].

- Requisito: "avoid signed-integer overflow when validating a start offset combined with a count against a total length". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/12 exibidos): ["against" "avoid" "combined" "count" "integer" "length" "offset" "overflow" "signed" "start" "total" "validating"].

- Requisito: "clamp an excessive requested count to the remaining length instead of treating it as an error". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (9/9 exibidos): ["clamp" "count" "error" "excessive" "instead" "length" "remaining" "requested" "treating"].

- Requisito: "express a sub-range using two absolute bounds instead of an offset-plus-count pair". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (11/11 exibidos): ["absolute" "bounds" "count" "express" "instead" "offset" "pair" "plus" "range" "sub" "two"].

- Requisito: "reject an out-of-range start position unconditionally, not only in debug builds". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Restrições preservadas; escopo e cumprimento da negação não foram verificados. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Restrições: escopo e cumprimento da negação não verificados.
  Marcador "not", texto original "reject an out-of-range start position unconditionally, not only in debug builds", estado `not_verified`.
  Termos associados (0/0 exibidos): []. Sem correspondência (9/9 exibidos): ["builds" "debug" "only" "out" "position" "range" "reject" "start" "unconditionally"].

#### Hipóteses de ação por candidato

Rótulos e chamadas não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade. A ordem dos candidatos prioriza leitura, não confiança semântica.

- Consulta <code>&#34;Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Modo da consulta não natural; ação não avaliada. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;avoid signed-integer overflow when validating a start offset combined with a count against a total length&#34;</code>; ação <code>&#34;start&#34;</code>; estado <code>&#34;hypothesis&#34;</code>; motivo <code>&#34;Hipótese lexical da ação solicitada: start. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;clamp an excessive requested count to the remaining length instead of treating it as an error&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;express a sub-range using two absolute bounds instead of an offset-plus-count pair&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;reject an out-of-range start position unconditionally, not only in debug builds&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Restrições ou negação presentes; ação não avaliada. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.


## Recuperação em local/gsl-b02e3c12@18ec80b61c87dbdc1472bb3a7169734af0db2d14e435e490ba6eafa137d5aa72

##### Interpretação lexical da consulta

Símbolos explícitos são pedidos de busca; restrições não são afirmações sobre o projeto.

Restrições registradas: 2. Escopo e cumprimento não verificados, inclusive nas linhas não exibidas abaixo.

- Marcador "without", texto original "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.", estado `not_verified`.
- Marcador "not", texto original "reject an out-of-range start position unconditionally, not only in debug builds", estado `not_verified`.

- Original: "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.". Modo: `mixed`. Símbolos (1/1 exibidos): ["SoPath"].
- Original: "avoid signed-integer overflow when validating a start offset combined with a count against a total length". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "clamp an excessive requested count to the remaining length instead of treating it as an error". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "express a sub-range using two absolute bounds instead of an offset-plus-count pair". Modo: `natural`. Símbolos (0/0 exibidos): [].
- Original: "reject an out-of-range start position unconditionally, not only in debug builds". Modo: `natural`. Símbolos (0/0 exibidos): [].

##### Indícios e lacunas por requisito

Correspondência lexical é inferência; requisito satisfeito e comportamento continuam não verificados. Termos sem correspondência não demonstram capacidades ausentes.

- Requisito: "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Restrições preservadas; escopo e cumprimento da negação não foram verificados. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Restrições: escopo e cumprimento da negação não verificados.
  Marcador "without", texto original "Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.", estado `not_verified`.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/26 exibidos): ["arithmetic" "copy" "correctly" "count" "empty" "excessive" "handling" "indexed" "integer" "like" "model" "needed"].

- Requisito: "avoid signed-integer overflow when validating a start offset combined with a count against a total length". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (12/12 exibidos): ["against" "avoid" "combined" "count" "integer" "length" "offset" "overflow" "signed" "start" "total" "validating"].

- Requisito: "clamp an excessive requested count to the remaining length instead of treating it as an error". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (9/9 exibidos): ["clamp" "count" "error" "excessive" "instead" "length" "remaining" "requested" "treating"].

- Requisito: "express a sub-range using two absolute bounds instead of an offset-plus-count pair". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Termos associados (0/0 exibidos): []. Sem correspondência (11/11 exibidos): ["absolute" "bounds" "count" "express" "instead" "offset" "pair" "plus" "range" "sub" "two"].

- Requisito: "reject an out-of-range start position unconditionally, not only in debug builds". Estado: `not_verified`. Nenhuma fonte lexical ou documental confirmada para este requisito no escopo; requisito e comportamento não verificados; não demonstra ausência da capacidade. Restrições preservadas; escopo e cumprimento da negação não foram verificados. Recuperação parcial ou não avaliada; cobertura insuficiente para concluir ausência.
  Restrições: escopo e cumprimento da negação não verificados.
  Marcador "not", texto original "reject an out-of-range start position unconditionally, not only in debug builds", estado `not_verified`.
  Termos associados (0/0 exibidos): []. Sem correspondência (9/9 exibidos): ["builds" "debug" "only" "out" "position" "range" "reject" "start" "unconditionally"].

#### Hipóteses de ação por candidato

Rótulos e chamadas não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade. A ordem dos candidatos prioriza leitura, não confiança semântica.

- Consulta <code>&#34;Model a safe sub-range copy operation (like SoPath::copy(start, count)) over an indexed sequence, correctly handling an empty or out-of-range start, an excessive requested count, and the arithmetic needed to validate start+count without signed-integer overflow.&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Modo da consulta não natural; ação não avaliada. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;avoid signed-integer overflow when validating a start offset combined with a count against a total length&#34;</code>; ação <code>&#34;start&#34;</code>; estado <code>&#34;hypothesis&#34;</code>; motivo <code>&#34;Hipótese lexical da ação solicitada: start. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;clamp an excessive requested count to the remaining length instead of treating it as an error&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;express a sub-range using two absolute bounds instead of an offset-plus-count pair&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Nenhuma família de ação reconhecida na consulta. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.
- Consulta <code>&#34;reject an out-of-range start position unconditionally, not only in debug builds&#34;</code>; ação <code>&#34;&#34;</code>; estado <code>&#34;not&#95;evaluated&#34;</code>; motivo <code>&#34;Restrições ou negação presentes; ação não avaliada. Rótulos e chamadas são indícios sintáticos: não comprovam execução, resolução de chamadas ou efeitos. Objeto e condições não foram verificados; oposição não prova incompatibilidade.&#34;</code>.


## Evidências

Consulte [assessment-report.json](assessment-report.json) e [evidence/bundles.json](evidence/bundles.json).
