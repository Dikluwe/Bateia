# Gap — documentação Doxygen embutida em headers

## Observação reproduzida

No dossiê `relatorios/coin5-pathassign-dossie`, o Bateia encontrou
`unlink()` e `is_linked()` como indícios de desregistro explícito no candidato
`local/boost-intrusive-ddb60e30`, commit imutável
`81bbd0043760edd5c34a4d3a5b11f12d76e52f97f23fc45809c49e0cfed59631`.
Esses símbolos não respondem, porém, à decisão específica tomada pelo operador
de atribuição.

O mesmo bundle contém evidência documental mais pertinente em
`include/boost/intrusive/list_hook.hpp`, SHA-256
`05a54f6538a093dc51355280ef8cce6ba8e3348ed4af08acd3ad776a2f6cd845`:

- linhas 112–119: o bloco Doxygen declara que a atribuição é uma função vazia,
  ignora o argumento e recomenda `swap` para emular movimento;
- linha 120: `list_base_hook& operator=(const list_base_hook& )`;
- linhas 235–242 e 243: a mesma relação para `list_member_hook::operator=`.

Há ocorrências equivalentes em `include/boost/intrusive/slist_hook.hpp`, linhas
113–121 e 239–247. A leitura manual consegue associar o contrato à declaração;
o extrator atual não materializa essa associação como evidência documental.
Esse formato também é relevante para bibliotecas como Coin, que documentam
grande parte das APIs diretamente em comentários Doxygen de arquivos `.hpp`.

## Causa no Bateia

A classificação documental é atualmente feita no nível do item coletado.
Arquivos `.hpp` entram como `file`, não como `documentation`. Em
`locateRuleContentMatches`, somente itens já classificados como documentação
usam a linha bruta; para código, `lexicalCodeLines` substitui comentários por
espaços. A proteção é correta para evidência de implementação: texto em
comentário não deve satisfazer uma regra de `symbol`. Como efeito colateral, um
bloco Doxygen ligado a uma declaração real desaparece também do canal
`documentation`.

Portanto, o gap não deve ser corrigido aceitando comentários como símbolos. A
correção precisa extrair documentação embutida como um canal separado e manter
a declaração de código associada como âncora.

## Comportamento desejado

Para extensões configuradas como código, o coletor/localizador deve reconhecer
blocos Doxygen (`///`, `//!`, `/**`, `/*!`) e produzir uma ocorrência documental
derivada com:

- caminho, identidade imutável do repositório, hash do arquivo e intervalo de
  linhas do comentário;
- linha e símbolo da declaração imediatamente documentada, quando a associação
  for inequívoca;
- classe `documentation_content`, nunca `strong_content` apenas por estar junto
  de código;
- trecho do comentário que efetivamente casou com a regra, e não o cabeçalho ou
  a primeira linha do arquivo;
- proveniência explícita de que a documentação foi extraída de um arquivo de
  código.

A associação deve parar diante de outra declaração, diretiva ou comentário não
documental incompatível. Comentários comuns, literais, código desabilitado por
pré-processador e documentação sem declaração associável não podem virar prova
forte. Documentação isolada continua sendo indício de intenção, não verificação
do comportamento executado.

## Critérios de aceitação

1. Uma regra documental contendo `Empty function. The argument is ignored`
   localiza `list_hook.hpp:112` e ancora a evidência em `operator=` na linha 120.
2. A mesma frase configurada como regra de símbolo não casa com o comentário.
3. `unlink()` escrito somente em comentário, string ou bloco desabilitado não é
   emitido como implementação.
4. Comentários `//!` consecutivos, blocos `/** ... */`, templates e declarações
   divididas em várias linhas têm testes de regressão.
5. A saída preserva separadamente linha do trecho documental e linha da
   declaração, sem substituir uma pela outra.
6. A confiança não sobe para `PROVEN` apenas pela coocorrência de comentário e
   declaração; ausência de associação segura permanece `INCONCLUSIVE`.

## Impacto esperado

O ganho é precisão explicativa: o Bateia passa de “há operações de ligação e
desligamento nesta biblioteca” para “este operador declara deliberadamente esta
semântica”, mantendo a fronteira entre contrato documentado e efeito verificado.
Isso deve beneficiar especialmente bibliotecas C/C++ com documentação de API
embutida em headers, sem reabrir os falsos positivos em comentários que o
localizador de ocorrência probatória passou a rejeitar.

