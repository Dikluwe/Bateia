# Prompt L0 — refino O6 de extração e artifact read

Materialize exclusivamente
`04_wiring/step0028_gate/o5blackbox/o6_extraction_artifact.go` como extensão
append-only do ORACLE segregado do Passo 0028. Não editar os predicados O5 nem o
dispatcher de portas; a autoridade integradora fará o wiring separadamente.

Regime: protocolo Tekt completo, com atestação limitada a
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

Entradas autorizadas:

- contrato R3 `00_nucleo/benchmark/0028/contract.md`, SHA-256
  `f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b`;
- catálogo `tekt.mechanism-rules.v2`, SHA-256
  `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5`;
- protocolo público `executable-mutant-protocol-o5.json`, SHA-256
  `455d6beec665765a0cca96b209ed5e651e7103c9a9dc5bebba51faa1b23b1fdb`;
- os artefatos O5 de extração de autoria deste papel;
- diagnóstico sanitizado `coverage-gap-o5.json`, SHA-256
  `bd93b949127bd80495af2840d9bd270c8aa36f6091086cd2c4550b40d2f396bb`;
- diagnóstico sanitizado `discriminatory-failure-o5.json`, SHA-256
  `dc9f22846b8effe85e6fc97ebe0c9f27b50f567977ed5fcf865a58456fd4a96d`;
- skill `tekt-materializacao-segregada` e referências obrigatórias.

Não ler ataques, records, subjects, runner/configuração, manifests, roots,
commitments, candidato, coverage privado, outros gaps/failures, receipts, seals,
verdicts ou artefatos de outros oráculos. Não usar identificador de caso, hash de
binding, expectativa, marcador, descrição de mutação ou branch específico de
fixture.

Fornecer duas APIs genéricas no package `o5blackbox`:

- `evaluateMechanismExtractO6(*evaluation, json.RawMessage, ContractOutput)`;
- `evaluateArtifactReadO6(*evaluation, json.RawMessage, ContractOutput)`.

O helper de extração deve manter catálogo/hash exatos, join uno, recomputação de
content hash, rule-role-kind catalogados, identidade imutável, status, novidade,
ID causal, gaps, proveniência e safety. Deve remover overconstraints meramente
representacionais: `data_flow`, state/failure model, components e witnesses não
ganham ordenação inventada; somente conjuntos e ordens expressamente normativos
são comparados canonicamente. A opacidade só é emitida quando uma ambiguidade R3
real foi preservada e nenhum predicado aplicável falhou.

O helper de artifact read deve:

- validar request/result estritos, allowlist, identidade do artifact, hash, size,
  base64, cursor opaco e avanço de paginação;
- interpretar os bytes publicados na ordem fixa de campos
  `schema_version,artifact,entries`, sem confundi-los com
  `TEKT_CANONICAL_JSON_V1`, que ordena chaves somente para hashes semânticos;
- rejeitar chaves JSON duplicadas e verificar SHA-256 sobre bytes de publicação
  quando a página contém o documento integral;
- preservar ordens de entries explicitamente contratadas, sem impor ordem a
  artefatos herdados ou sequências de domínio;
- permitir reads do store, mas rejeitar rede, writes, processos, mutações,
  admissão durável e vazamentos de conteúdo privado;
- tratar disponibilidade, integridade armazenada, compatibilidade de run e
  autenticidade de cursor não observáveis como `Unknown` somente mediante erro
  controlado aplicável;
- reconhecer Unknown materializado por status/gap contratual, sem convertê-lo em
  sucesso e sem sobrescrever uma violação simultânea.

O resultado cobre somente o fragmento observável R3 e não emite seal, receipt,
certificado nem veredito reservado.

## Suplemento O9

Reconheça também `novelty=UNKNOWN`, `transfer_mode=UNKNOWN` e disponibilidade não
evidenciada como opacidades R4 reais. Não relaxe ordem fixa, hash, size,
identidade, paginação, ordenação de entries ou safety; violação simultânea sempre
precede `Unknown`.

## Suplemento O11

Substitua a calibração O9 ampla. `novelty=UNKNOWN` ou
`transfer_mode=UNKNOWN`, isoladamente, são resultados derivados conservadores e
não provam opacidade do gate. `Unknown` requer uma incerteza epistemológica R4
concreta, como identidade/classe/testemunha insuficiente, baseline ausente,
compatibilidade não medida, licença não reconhecida, budget esgotado ou
disponibilidade não evidenciada. Uma única entrada realmente opaca permanece
opaca mesmo quando o mesmo artifact contém outras entradas resolvidas; nenhuma
opacidade mascara violação decidível.

Na extração, uma lacuna R4 válida ou mecanismo `INCONCLUSIVE` permanece
`Unknown` mesmo ao lado de mecanismos `PROVEN|PARTIAL`. Para o preimage causal
de mecanismo R4 documental, use o RuleID `DOCUMENTARY` efetivamente ligado
quando toda a evidência positiva é documental literal; evidência técnica ou
mista conserva o RuleID pai.

## Suplemento O13

Refine `REAL_OPACITY_CLASSIFICATION` somente pelas relações públicas do feedback
O12. Licença `UNRECOGNIZED`, warning `LICENSE_NOT_RECOGNIZED` e
`integration_unknowns` explícitos conformes são resultados conservadores
obrigatórios e não tornam o gate `Unknown` isoladamente. Em
`MECHANISM_EXTRACT`, baseline convencional ausente/inválido torna `Unknown` um
mecanismo `PROVEN` cuja `novelty=UNKNOWN`, pois nesse caso a falta do baseline é
o limite causal observável; para mecanismo `PARTIAL`, a novidade já seria
`UNKNOWN` independentemente do baseline. Violação decidível continua tendo
precedência absoluta.

## Suplemento O14

Reconheça relacionalmente, sem IDs de fixture, os estados opacos literais do
protocolo O8: witness mínimo semanticamente insuficiente; target `UNKNOWN` com
baseline ausente; composição não medida; licença `UNRECOGNIZED` cujo
transfer/status permanece incapaz; counterevidence com budget esgotado;
identidade imutável ausente/ambígua; disponibilidade de provider não evidenciada
por capture congelado; e efeito de utility/security sem baseline. Exija os
status, plano e contexto correspondentes. `UNRECOGNIZED` com
`PATTERN_ADAPTATION` ou `REFERENCE_ONLY` explícito e hipótese apenas `PROPOSED`
permanece `Preserved`. `Unknown` só é emitido quando nenhum predicado simultâneo
falhou.
## Refinamento O15

Calibre `REAL_OPACITY_CLASSIFICATION` pelas relações estruturais públicas do
protocolo O8 e pelo tipo declarado do artefato. Um grupo técnico `PARTIAL` que
não alcança o witness mínimo é epistemicamente limitado; documentary-only
`PARTIAL` permanece decidível. Em `hypothesis_gaps` e `mechanism_gaps`, códigos
controlados podem representar baseline, composição, budget, identidade,
provider ou witness indisponível, mas texto livre jamais controla o resultado.
Licença `UNRECOGNIZED` só é opaca com transfer `UNKNOWN` ou status
`INCONCLUSIVE`; `PATTERN_ADAPTATION` e `REFERENCE_ONLY` propostos permanecem
preservados. Toda violação simultânea precede `Unknown`.

## Refinamento O16

Não promova limites decidíveis do contrato a proxies de opacidade. Um witness
técnico `PROVEN` que produz agregado `PARTIAL`, target `UNKNOWN` acompanhado de
plano e `integration_unknowns` obrigatórios em proposta são estados válidos e
decidíveis. A insuficiência semântica exige assessment técnico `PARTIAL` ou gap
controlado coerente. Licença `UNRECOGNIZED` na hipótese continua dependente de
transfer/status; o code `LICENSE_NOT_RECOGNIZED` somente é opaco quando aparece
no artifact tipado `hypothesis_gaps`. Texto livre não participa da decisão.

## Refinamento O17

Consolide a política total de `Unknown` em uma tabela tipada. Relações completas
de hipótese exigem os enums de incapacidade previstos por R4; gaps somente usam
o campo controlado `code` dentro de `mechanism_gaps|hypothesis_gaps`. Reconheça
as formas públicas missing/absent/unknown/unavailable/unmeasured/insufficient,
ambiguous/unbound/unresolved, without, not-evidenced/not-measured/
not-established/not-available para baseline, compatibility, witness, immutable
identity/binding, provider/capture e utility/security. Budget exhausted e
`LICENSE_NOT_RECOGNIZED` são literais. Texto humano nunca decide.

## Refinamento O18

Substitua aproximações lexicais pela tabela pública necessária e suficiente do
coverage gap O17. Composição e counterevidence exigem todas as relações tipadas,
incluindo `INCONCLUSIVE` e `validation_plan` não vazio. Códigos de gap usam
somente as conjunções publicadas para composition/compatibility, license,
counterevidence/budget, identity/binding, provider/capture e
utility|cost|security com baseline. Preserve O01–O03 e a precedência de falhas.

## Refinamento O19

Na fronteira serializada pública, O04–O10 são documentos schema `1` do artifact
`hypothesis_gaps`, com entry tipada por `hypothesis_id`, `code`, `source_ref`.
Reconheça por igualdade exata somente os sete códigos normativos publicados na
projeção sanitizada O18. Ausência de qualquer campo obrigatório não ativa esse
predicado. Preserve O01–O03 e faça toda violação preceder `Unknown`.

## Refinamento O20

Reconheça as sete relações no shape end-to-end público: hipótese `PROPOSED` sem
validation plan e (a) composição múltipla com integration unknowns, (b) licença
`UNRECOGNIZED` em `REFERENCE_ONLY`, (c) counterevidence
`NONE_FOUND_IN_SCOPE` com escopo não vazio e resultado vazio, ou (d) target
`UNKNOWN` com target gap refs vazio. A presença de plano não vazio torna o
vizinho decidível. Não use semantic_state nem texto humano.

## Refinamento O23

O protocolo público O22 substitui os sete estados antes indistinguíveis por
documentos `hypothesis_gaps` binding-free. Para identidade de hipótese no
formato canônico, aceite `Unknown` somente quando o documento integral possui
uma entrada, ordem fixa `hypothesis_id,code,source_ref`, nenhuma propriedade
adicional, um dos sete pares normativos exatos e exposição do mesmo artifact e
SHA-256 na observação de publicação. Valide apenas o formato da identidade e as
relações de integridade; seus valores não controlam a classificação. Um par,
shape ou relação de publicação vizinho permanece `Preserved`, e qualquer
violação contratual continua precedendo a opacidade.
