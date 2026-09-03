# Prompt L0 — predicados black-box de extração O5

Materialize exclusivamente `04_wiring/step0028_gate/o5blackbox/extraction.go`, no
papel segregado de autoridade ORACLE do Passo 0028.

Entradas autorizadas e congeladas:

- `00_nucleo/benchmark/0028/contract.md`, revisão R3;
- `00_nucleo/benchmark/0028/mechanism-rules.json`, schema
  `tekt.mechanism-rules.v2`, versão 2 e SHA-256
  `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5`;
- `00_nucleo/benchmark/0028/executable-mutant-protocol-o5.json`;
- a skill `tekt-materializacao-segregada` e suas referências obrigatórias.

Não ler implementação candidata, descrições adversariais, records/configuração do
runner, bindings privados, compromissos de mutação, manifestos ou raízes de hashes
de entradas, gaps/falhas de cobertura, artefatos de oráculos anteriores, receipts,
seals ou verdicts. O trabalho ocorre em workspace compartilhado e, portanto, deve
ser reportado como `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

Implementar
`evaluateMechanismExtract(state *evaluation, request json.RawMessage, output ContractOutput)`
como predicados externos, determinísticos, genéricos e binding-free. A função deve
depender somente do request contratual e da saída contratual, nunca de identificador
de registro, hash de binding, expectativa pré-carregada, marcador de classificação
ou branch específico de fixture.

Os predicados devem cobrir:

- decode estrito do request e do result, campos obrigatórios, ausência de campos
  extras e slices públicas não nulas;
- catálogo fechado exato, versão e SHA-256 congelados; catálogo divergente deve ser
  recusado por erro controlado de integridade e nunca produzir resultado;
- identidade imutável `owner/repo@commit`, join exatamente uno entre source avaliada
  e item do mesmo bundle por `(kind,ref,path,content_hash)`, e recomputação do
  `sha256:` dos bytes de `content_utf8` somente para integridade;
- `RuleID` como única identidade de regra e `role`/`witness_kind` derivados somente
  da entrada literal do catálogo, com compatibilidade estrutural de source kind;
- agrupamento estrito por repositório, commit e regra de mecanismo, deduplicação de
  witnesses e proibição de composição cross-repo ou cross-commit;
- precedência de `INCONCLUSIVE`, teto `PARTIAL` e requisitos de `PROVEN`, inclusive
  implementação e testemunha independente `PROVEN`, hashes distintos, campos
  centrais completos e ausência de contraevidência central;
- novidade calculada somente depois do status contra o baseline convencional:
  `CONVENTIONAL`, `DIVERGENT` apenas para mecanismo `PROVEN`, ou `UNKNOWN`;
- gaps reais para regra fora do catálogo, binding zero/ambíguo, hash divergente,
  identidade/estrutura/status inconsistente e witness insuficiente; `Unknown` não
  pode ser transformado em prova, ausência ou sucesso;
- campos de catálogo, witnesses, rule IDs, source, status, novidade e `mechanism_id`
  causal verificáveis; arrays conjuntistas deduplicados e ordenados e mecanismos e
  gaps na ordem contratual;
- segurança e pureza: zero rede, filesystem, processos, mutações externas, admissão
  durável ou publicação; nenhum conteúdo bruto, snippet, path físico, segredo,
  canário ou campo privado na saída.

O oráculo deve emitir findings apenas pelos métodos do `evaluation` fornecido pelo
package e não deve criar executor, consultar filesystem/rede, importar expectativas
ou executar o software avaliado. Ele atesta somente o fragmento observável de
extração definido por R3.

## Suplemento O8 / R4

O binding ativo passa a ser o catálogo literal `tekt.mechanism-rules.v3`, SHA-256
`52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142`.
Aceite `DOCUMENTARY` somente com um dos quatro RuleIDs documentais explícitos,
`witness_kind=DOCUMENTATION`, `status_max=PARTIAL` e source kind
`documentation`. Essa evidence pode satisfazer `PARTIAL`, mas nunca os mínimos
de implementação, independência ou hashes distintos de `PROVEN`.

## Suplemento O13

O predicado `CANONICAL_CAUSAL_MECHANISM_ID` deve reconstruir literalmente o
preimage fixo da seção 6/R4: campos na ordem normativa, JSON compacto UTF-8,
HTML escaping desabilitado e sem newline. Normalize localmente para NFC e
whitespace canônico `source`, RuleID, classe, propósito e estratégia, sem
consultar dependência, filesystem ou ambiente. `witness_hashes` contém todos e
somente os hashes de evidence e counter_evidence, normalizados, deduplicados e
ordenados por bytes. A implementação local pode ser conservadora para sequências
que não consegue provar, mas jamais pode inventar uma composição Unicode.

## Suplemento O14

O campo `rule_id` do preimage causal é sempre `MechanismRule.id`, resolvido pela
classe/source contra o registro do catálogo v3. Não o derive de
`ExtractionRuleIDs`: esses IDs preservam as regras das EvidenceSources e podem
ser documentais, técnicas ou de conflito sem substituir a identidade da regra
de mecanismo.
