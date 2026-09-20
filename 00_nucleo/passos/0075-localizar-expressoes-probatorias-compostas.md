# Passo 0075 — Localizar expressões probatórias compostas

## Estado e regime

- Estado: obrigação e contrato observável anteriores à implementação focal.
- Regime: Tekt A/B reduzido, executado como ensaio sem atestação de isolamento; o mesmo
  ambiente e autoridade integram contrato, teste e candidato.
- Baseline: worktree posterior aos passos 0072–0074, com alterações ainda não
  consolidadas. A correção deve ser focal e preservar esses artefatos.
- Budget: uma implementação e até duas revisões focais; nenhum ajuste de catálogo para
  fabricar o resultado.

## Diagnóstico congelado

No dossiê `lab/bancada/history/reports/0075-coin-sbhash/dossier`, a recuperação dirigida funcionou:

- `bits/hashtable.h` foi candidato e arquivo examinado para
  `self-assignment-guard-hashtable`; o guard está na linha 1257;
- `src/include/robin_hood.h` também foi candidato e examinado; o guard está na linha
  1614.

Logo, o defeito não é a exclusão de padrões com espaço das sementes de caminho. Os
termos dirigidos permanecem `hashtable.h` e `robin_hood.h`, derivados de `TargetHints`.
O defeito ocorre na atribuição de autoridade intrarquivo: o padrão literal
`&__ht == this` casa com a linha limpa, mas `supportedSymbolOccurrence` aceita apenas
declarações ou identificadores seguidos de `(`/`<`, rejeitando expressões com operador.

## Obrigação

Quando um padrão forte declarado descreve uma expressão de código composta e sua forma
literal ocorre integralmente numa linha de código admissível, o localizador deve emitir
a ocorrência exata como `STRONG_CONTENT`. A expressão precisa conter operador de código
explícito e identificadores; uma sequência natural de palavras não se torna expressão.

Continuam válidas as fronteiras do passo 0073:

- comentários e literais são removidos antes do casamento;
- substring de identificador não casa;
- macros, token pasting e expressão dividida em linhas permanecem opacos;
- o achado lexical produz no máximo evidência parcial, nunca prova comportamental;
- padrões probatórios não vazam para `EvidenceRequirementQuery.Terms`; recuperação de
  caminho continua derivada exclusivamente de `TargetHints`.

## Casos observáveis

Positivos:

1. `if (&__ht == this)` localiza o padrão `&__ht == this` na própria linha;
2. `if (&o == this)` localiza `&o == this` apesar de espaços variáveis;
3. `__p->_M_nxt = __next` e `__num % __den` são expressões compostas admissíveis;
4. a capacidade recebe `PARTIAL`, fonte `STRONG_CONTENT`, linha e snippet coerentes.

Negativos e opacos:

1. os mesmos textos somente em comentário ou string não geram fonte;
2. `&__htable == this` não casa `&__ht == this`;
3. texto natural multiword sem operador não recebe autoridade de símbolo;
4. expressão repartida em duas linhas não é reconstruída;
5. os termos de busca dirigida permanecem somente `hashtable.h`/`robin_hood.h`.

## Verificação

Gate focal:

```sh
go test ./01_core/usecase -run '^TestStep0075' -count=1
```

Depois do gate focal: pacote afetado, suíte completa, race do pacote, build,
`crystalline-lint .` e `git diff --check`. O veredito cobre somente localização lexical
de expressões declaradas; não comprova a ordem de mutações nem a segurança geral do
operador de atribuição.

## Resultado materializado

O gate pré-correção compilou e rejeitou os quatro positivos com
`LEXICAL_MATCH_OPAQUE`; negativos e separação entre recuperação/prova permaneceram
aceitos. Depois da alteração focal, o gate passou duas vezes, o pacote afetado passou
com e sem race, e a suíte completa, build, lint cristalino e diff-check passaram.

Replay do dossiê congelado, sem nova coleta:

- libstdc++: `PARTIAL (0.60)`, `bits/hashtable.h:1257`, padrão
  `&__ht == this`, classe `STRONG_CONTENT`;
- robin-hood-hashing: `PARTIAL (0.60)`, `src/include/robin_hood.h:1614`, padrão
  `&o == this`, classe `STRONG_CONTENT`.

O dossiê geral continua `INCONCLUSIVE` por outras capacidades; a correção não transforma
localização lexical em prova comportamental. Regime final: executado sem atestação de
isolamento e sem mutation score independente.
