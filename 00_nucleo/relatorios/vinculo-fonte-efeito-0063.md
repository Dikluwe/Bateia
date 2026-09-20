# Passo 0063 — Vínculo entre fonte e efeito

Implementado `--expect-sha256` em `effects verify` e `source_binding` no recibo JSON. O hash esperado identifica o arquivo completo previamente examinado. Se os bytes mudaram, o Bateia retorna unknown/source_identity_mismatch antes de gerar wrapper, compilar ou executar. Sem hash solicitado, o piloto continua disponível e declara not_requested.

O recibo agora localiza por AST a função selecionada: offsets de bytes [início,fim), linhas físicas inclusivas, trecho bruto e hash do trecho. O mesmo snapshot lido fornece o hash, a análise AST, o conteúdo registrado e os bytes compilados. Identidade de arquivo e resultado da propriedade são estados separados: matched não significa sucesso da função.

## Evidência observada

Demonstração pela CLI pública, com o mesmo caminho de arquivo entre execuções:

| Situação | Vínculo | Efeito | Chamadas completas |
| --- | --- | --- | --- |
| Referência e hash correspondente | matched | supported_in_scope | 39 |
| Helper modificado, hash antigo | mismatch | unknown | 0 |
| Helper modificado, hash atualizado | matched | counterexample | 39 |
| Helper modificado, hash omitido | not_requested | counterexample | 39 |

A declaração Transform permaneceu idêntica nas versões: bytes [77,135), linhas 6–8, SHA-256 aae12cbfa9011c0c7c88be6eeb9342438f3d4aa28359b1d893fa0d2fc2220cd3. Alterar apenas o helper mudou o efeito: a entrada [2,2,1,2] passou a [1,2,2,2], violando preservação. O hash de arquivo mudou e o hash do trecho não. Isso demonstra por que vincular somente o nome ou corpo da função não basta.

Os quatro recibos e verificações estão em `../../lab/bancada/history/benchmarks/0063/executions/`. A conferência dos efeitos completos reutilizou o oráculo Python independente congelado no passo0062. As oito fixtures anteriores mantiveram os resultados esperados: referência supported_in_scope, cinco defeitos counterexample e dois casos unknown. Evidência em `../../lab/bancada/history/benchmarks/0063/regression0062/verification.json`.

## Revisão encontrada pelos testes

A primeira execução dos testes independentes rejeitou três fontes válidas. O filtro herdado procurava a substring go: em qualquer comentário, confundindo //line inventado.go:700 com diretiva //go:. A revisão focal passou a reconhecer o prefixo real //go: do comentário AST. O contrato e as expectativas permaneceram congelados; todos os novos testes passaram depois da correção. Logs, hashes e snapshot r0 estão em `../../lab/bancada/history/benchmarks/0063/revisions/r0/`. Esse snapshot foi reconstruído revertendo a alteração focal e confirmado pelos hashes congelados antes dos testes, pois a primeira cópia ocorreu após a revisão.

Testes novos cobrem hash omitido, correto e maiúsculo, argumentos malformados, mudança só no helper, recorte rejeitado mesmo com hash correto, símbolos ausentes/duplicados, Unicode, CRLF e coordenadas físicas apesar de //line. O teste com stub de go no PATH só observa comandos resolvidos por esse PATH; não atesta ausência de todo subprocesso. A conclusão de preflight sem execução também se apoia na auditoria do fluxo (retorno anterior à criação de temporários/build) e nos recibos sem wrapper/build/binário e sem chamadas completas.

## Interpretação e limites

- Fato observado: hashes dos bytes, declaração localizada e resultados das chamadas instrumentadas.
- Inferência limitada: os bytes executados correspondem ao hash esperado fornecido; a propriedade vale somente nas entradas exercitadas quando supported_in_scope.
- Hipótese ainda aberta: esses bytes implementam a intenção da consulta em outros contextos e entradas.
- Não verificado: vínculo com repositório/commit, resolução de dependências externas, integração automática com busca, correção geral, sandbox e efeitos além do contrato0062.

Obter o hash da evidência previamente examinada; recalculá-lo cegamente antes de cada execução aceita a versão atual e não protege a ligação com a versão anterior. Declaration é localização, não prova de execução. Tempos continuam sendo telemetria da ferramenta.

Autoria A/B: integrador escreveu contrato e verificou; A escreveu produção/prompts; B escreveu testes/prompts sem ler candidato. Contextos separados, filesystem compartilhado: ensaio sem atestação física. Manifesto/baseline e freezes em `../../lab/bancada/history/benchmarks/0063/`. Nenhum score ou mecanismo de busca foi alterado. O próximo vínculo com resultados de busca ainda precisa preservar a identidade e o contexto completo de compilação; este passo fornece uma precondição verificável, não essa integração.

## Uso

```sh
go build -o /tmp/bateia-0063 ./04_wiring
/tmp/bateia-0063 effects verify --source /caminho/candidate.go --symbol Transform --expect-sha256 HASH_DA_VERSAO_EXAMINADA --output /tmp/efeito.json
```

Hash deve ter 64 dígitos hexadecimais; maiúsculas são normalizadas. Exit code zero também ocorre em unknown/counterexample; ler status global e source_binding separadamente.

Validação final: testes novos sem cache, suíte Go completa, build, crystalline-lint e git diff --check passaram. Contrato e testes conservaram seus hashes; uma revisão focal da implementação foi registrada. Comandos e resultados em `../../lab/bancada/history/benchmarks/0063/verification.json`. Sem commit.
