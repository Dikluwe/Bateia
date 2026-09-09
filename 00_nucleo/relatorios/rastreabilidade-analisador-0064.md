# Passo0064 — Identidade do analisador e parâmetros

Cada recibo de effects verify agora inclui trace: identidade do executável Bateia, informações disponíveis de build, intervalo UTC e fingerprint dos parâmetros efetivos. O hash da fonte e a localização da declaração do passo0063 permanecem separados. Nenhum score, consulta ou propriedade foi alterado.

## O que foi observado

A demonstração pública gerou quatro recibos, conferidos por um script Python conforme o protocolo congelado:

| Cenário | Efeito | Trace | Comparação |
| --- | --- | --- | --- |
| Referência | supported_in_scope | recorded | Hash do produtor igual ao SHA-256 independente do binário Bateia |
| Mesma fonte em outro caminho | supported_in_scope | recorded | Mesmo request_sha256 |
| Mesma fonte, seed 2 em vez de 1 | supported_in_scope | recorded | request_sha256 diferente |
| Fonte modificada, hash esperado antigo | unknown | recorded | Fonte identifica o novo snapshot; nenhuma candidata executada |

Nos três casos executados, 39 observações completaram cada avaliação. Hashes de produtor e candidata são distintos e aparecem em campos separados. O trace da versão construída registrou módulo bateia, Go go1.26.5, revisão de build 9056cd259e40d6360a06e6e21a157aa2e5b9eac7 e vcs_modified=true. Essa revisão pertence ao analisador; não foi atribuída à fonte candidata. Dados disponíveis foram obtidos do executável em execução, sem consultar .git durante a análise.

Recibos e conferência: ../benchmark/0064/executions/verification.json. Comparações: executions/comparison.json. O protocolo canônico do fingerprint está no contrato ../passos/0064-identificar-analisador-e-parametros.md e foi conferido independentemente em Go e Python.

## Significado dos campos

- trace.analyzer.binary_sha256: bytes do Bateia produtor. Em Linux, leitura stream de /proc/self/exe, sem cache global ou fallback para um pathname potencialmente substituído.
- trace.analyzer.module_version, go_version, vcs_revision, vcs_modified: informações de build quando disponíveis. Ausência de VCS fica vazia/null; não significa árvore limpa.
- trace.inputs: contrato, produtor, fonte, símbolo, seed, timeout de execução, timeout de build e limites stdout/stderr.
- trace.request_sha256: fingerprint canônico desses inputs. Ignora caminho, horários, hash esperado e resultado. Alterar apenas a exigência de identidade da fonte não muda o trabalho identificado.
- trace.started_at/finished_at: intervalo UTC da ferramenta, incluindo preparação. Não é benchmark da candidata. Se o relógio de parede recuar, a implementação limita finished_at ao início; não há atestação de relógio nem precisão temporal garantida.

recorded não é um novo nível de aprovação da função. Um pedido interrompido por mismatch pode estar bem identificado. Se não for possível identificar o produtor, trace fica unknown e o fingerprint vazio, sem promover nem invalidar por si só o veredito de efeito.

## Autoria e correção dos testes

Contrato/baseline congelados antes de A/B; produção e testes em contextos separados, filesystem compartilhado. Ensaio Tekt sem atestação física. Nenhuma alegação de equivalência geral ou selo discriminatório completo.

Na primeira execução, dois testes exigiram Cases vazio no preflight. Essa exigência excedia o contrato0063, que exige zero casos completos e permite conservar os casos planejados. O autor B corrigiu apenas essa condição para inspecionar Observation.Completed; status unknown, binding, razão e demais exigências permaneceram. A implementação e o contrato não mudaram. Falha, snapshots e freezes B r0 preservados em ../benchmark/0064/revisions/b-r0/. Os testes novos passaram após a correção. Também conferem hash independente do processo de teste, metadados runtime/debug, coordenadas temporais UTC, inputs, sensibilidade, seed máxima e entradas inválidas.

## Limites da conclusão

Fatos observados são os bytes/metadata lidos, parâmetros registrados e chamadas completas do piloto. A inferência é que o recibo identifica esse produtor e esse pedido no recorte observado. Generalização da propriedade continua hipótese. A identidade de um repositório candidato/commit, autenticidade criptográfica do recibo, integralidade do worktree e isolamento físico continuam não verificados.

request_sha256 não é checksum do JSON inteiro nem assinatura digital. Um terceiro pode editar um recibo ou recalcular hashes; não há prova de autenticidade. Mesmo fingerprint não garante determinismo, mesmos resultados nem identidade de todo ambiente externo. O compilador da candidata continua identificado separadamente nos metadados execution do piloto. Não houve execução de projetos externos, somente fontes controladas dos passos0062/0063.

## Uso

O comando permanece compatível:

```sh
go build -o /tmp/bateia-0064 ./04_wiring
/tmp/bateia-0064 effects verify --source /caminho/candidate.go --symbol Transform --expect-sha256 HASH_PREVIAMENTE_EXAMINADO --seed 1 --timeout 1s --output /tmp/efeito.json
```

Comparar trace.inputs para explicar diferenças entre pedidos; source_binding para identidade previamente esperada; cases/status para o resultado limitado da propriedade. Um hash diferente sozinho não explica a causa da diferença.

Validação final: testes novos sem cache, suíte Go completa, build, crystalline-lint e git diff --check passaram. Quatro cenários novos, quatro regressões0063 e oito regressões0062 aprovados. Contrato e produção mantiveram freezes; houve uma revisão contratualmente justificada dos testes. Comandos/logs em ../benchmark/0064/verification.json. Falha de leitura de identidade e recuo de relógio foram auditados no código, mas não forçados no ambiente. Sem commit.
