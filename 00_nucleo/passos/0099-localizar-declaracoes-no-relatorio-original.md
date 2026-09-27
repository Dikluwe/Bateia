# 0099 — Localizar declarações no relatório original

## Resultado e uso

O grafo 0098 continua decision-trace-v1, com os mesmos IDs, bytes e hashes.
Um índice separado acrescenta localizadores e resumos das declarações de cada
nó. Não é uma aprovação, um novo decisor ou uma prova de evidência correta.

`idea recommend --trace` também emite trace-index.json. O estágio trace em
run.json identifica os hashes do grafo (primeiro) e índice (segundo), com contrato
de projeção decision-trace-projection-v2. Isso não muda o schema do grafo. Sem a flag,
a pesquisa simples e suas saídas ficam intactas. O índice não participa de
RunID, RuleVersions ou hashes causais de checkpoint.

`idea decide --trace-output grafo.json` continua emitindo somente o grafo.
Para explicá-lo, há um comando offline, com diretório de saída já existente:

```sh
./bateia idea trace explain --input pedido.json --output indice.json
```

Sem --output, escreve JSON em stdout. O pedido tem schema_version "1", graph
com o grafo completo congelado, build_input com a entrada original da construção
e node_ids opcional. build_input contém seu próprio schema_version "1", report
e idea/rule_versions quando usados. node_ids omitido/vazio explica todos os nós;
uma lista seleciona IDs existentes, sem dispensar validação do relatório completo.
IDs desconhecidos ou repetidos falham. Não há inferência de versões históricas.

O núcleo valida o grafo, reconstrói um com build_input e exige igualdade de
graph_hash. Relatório/política/ideia/versões alterados falham antes de emitir uma
explicação. A igualdade liga representações declaradas, não autentica sua autoria
ou veracidade. Um grafo manual diferente do construtor não é explicado como se
tivesse vindo desses documentos.

## Localizadores e resumos

A saída decision-trace-explanation-v1 tem graph_hash, report_fingerprint, scope
RECORDED_DEPENDENCIES_ONLY, entries e limitations. Cada entrada contém node_id,
kind, label, fingerprint, statement, statement_truncated e locations.

Cada location tem documento lógico report, idea ou rule_versions e JSON Pointer
RFC6901 relativo a ele. Ponteiro vazio significa a raiz. Exemplos:

- report /decision: veredito histórico.
- report /matrix/0/candidates/1: declaração para uma capacidade.
- report /reasons/0/evidence/0: fonte da primeira razão.
- report /composition/candidates/0/members/1/component/evidence/0: fonte de membro.

Esses nomes não são arquivos a abrir. O consumidor resolve ponteiros nos
documentos fornecidos em build_input. Fontes compartilhadas mantêm suas
ocorrências; localizadores são ordenados e não duplicados. Fingerprint igual
não identifica sozinho uma alegação: capacidades distintas podem ter payloads
idênticos. A indexação reconstrói a chave do ID e confere o fingerprint do
registro. Deve cobrir todos os nós construídos antes de aplicar seleção.

statement projeta o estado/mensagem existente, não um novo juízo. Limita-se a
1024 bytes sem cortar caracteres UTF-8; corte tem statement_truncated true e
localizador do registro integral. Guardar os documentos originais continua
necessário; o índice não duplica payloads completos nem reinterpreta PROVEN.

Fingerprint é do JSON tipado do registro, não dos bytes do arquivo-fonte, e não
é substituível por ContentHash. Para conferir um ponteiro no JSON original,
compactar espaços preservando a ordem das chaves. Remontar mapas genéricos com
outra ordem pode gerar outro hash. JSON externo pode ter números/campos/escapes
representados diferentemente: normalizar pelos tipos do contrato antes da
comparação. O vínculo é com a representação tipada, não com bytes arbitrários.

## Segurança e compatibilidade

Guardas de 0098 continuam: JSON estrito, hash/DAG coerentes, UTF-8 e números
finitos, limites antes da serialização integral, proteção de segredos configurados
e recusa de saídas que identificam entradas por path/hardlink/alias ancestral.
trace-index.json também passa por preflight antes de criar o dossiê. Escritas
de grafo/índice são atômicas individualmente, não uma transação conjunta. Falha
da emissão solicitada marca execução FAILED; consultar o manifest final para
identificar pertencimento dos artefatos, não apenas sua existência.

Entrada/saída compartilham limites conservadores de 0098: 2 MiB texto/JSON,
profundidade 64 e orçamento estrutural de 100000 valores. Esse teto já estava no
código ao retomar este incremento; corrigi a referência documental antiga de
20000 em 0098, sem alterar a guarda do núcleo. Como o pedido inclui grafo e build_input,
pode exceder o limite agregado mesmo se cada parte couber isoladamente.
Retomada sem --trace preserva arquivos antigos, mas não os consome/anuncia no
novo estágio; o aviso declara que não pertencem à execução atual. O grafo
histórico permanece utilizável por idea trace impact.

Não há consultas a fontes, execução de candidatos, vencimento automático,
autenticação ou aprovação. Campos não indexados individualmente continuam sob
a dependência global do relatório. O índice não certifica completude.

## Verificação

Testes conferem ponteiros/fingerprints no JSON original, capacidades com payload
igual, fontes compartilhadas, composições/verificações, seleção, adulteração,
truncamento UTF-8, limites, não mutação e grafo histórico intacto. Integração
compara bytes causais/RunID com/sem --trace, hashes, saídas protegidas, erros sem
stdout/arquivos e despacho real offline sem carregar .env.

Próximo incremento possível: executar uma tarefa delimitada de revisão e
registrar a nova decisão sem substituir a decisão histórica.


Validação final: `go test -p 1 ./...`, `go vet ./...`, `git diff --check`
e testes de corrida direcionados de domínio, núcleo, CLI e wiring passaram.
A suíte sequencial mantém as expectativas dos gates legados, sem alterar prazos
ou goldens para contornar a instabilidade paralela registrada em 0098.
