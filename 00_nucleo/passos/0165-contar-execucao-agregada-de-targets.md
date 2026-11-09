# 0165 — Contar execução agregada de targets

## Intenção anterior à materialização

Corrigir o achado confirmado na revisão de 7809644: um target derivado lê integralmente um arquivo, mas TargetsExecuted fica zero quando o path seguinte termina SKIPPED_BUDGET/Bytes0. Responsável L1: evidence_recovery. A contagem deve reconhecer qualquer tentativa correspondente da coleta atual, uma vez por target do plano, sem alterar a política de retomada que usa a última ocorrência.

## Comportamento e limites

Agregar resultados COLLECTED/NOT_FOUND/FAILED e fragmentos SKIPPED_BUDGET com Bytes>0 para a métrica. Omissão posterior não apaga execução anterior. Omissão integral e UNSUPPORTED isolados continuam sem execução; conservar fallback legado sem resultado correspondente. Não contar itens dos bundles originais como execução atual, duplicar targets por path, mudar bytes/chamadas ou promover evidência. Histórico, schemas e registros anteriores permanecem intactos.

## Aceitação e preservação

Regressão pela API ExecuteEvidenceRecovery com um target derivado e dois paths, saldo exato para o primeiro e zero para o segundo. Controles com um byte a menos/a mais, saldo para ambos e omissão integral por bytes/quantidade. Conferir status/bytes/conteúdo/hash, TargetsPlanned/Executed, chamadas reais e inputs intactos. Na omissão posterior, ExecutedRecoveryTargets e PlanEvidenceRecovery continuam permitindo retomada do target. Executar regressão antes/depois, gates canônicos com V0–V27 ativas e comparação SHA256 do acervo. Logs/binários em /tmp; não reproduzir gates selados como resultados atuais nem recalcular hashes históricos.

## Resultado e verificação

A regressão falhou antes do patch somente no saldo exato para o primeiro arquivo: TargetsExecuted era zero. Após a correção, os seis cenários passaram. A métrica agrega resultados da coleta atual e conta cada target uma vez; o predicado de histórico anterior permanece intacto. As verificações confirmam bytes, hashes, chamadas e inputs preservados; conteúdo original correspondente não fabrica execução nova, e a omissão posterior continua permitindo novo plano.

O probe público da revisão confirmou a correção no orçamento267, com metadados236, primeiro arquivo31 e segundo omitido: TargetsExecuted agora é1. Os controles266/268 e leitura completa mantiveram uma execução. Logs locais em /tmp/bateia-0165-before.log, /tmp/bateia-0165-focused.log e /tmp/bateia-0165-after-probe.log; não são recibos selados.

make verify passou em lint, arquitetura, go test ./..., go vet ./..., build de ./04_wiring e diff-check. Evidências/binário em /tmp/bateia-0165-verification. O alcance do linter confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true, 1.179 fontes lidas e zero violações/erros de leitura/sintaxe. Os 11 controles do parser e dois controles de V3 ligado/desligado passaram; as três exclusões históricas de oráculos permaneceram declaradas.

A comparação SHA256 confirmou 4.692 arquivos históricos de Lab, passos e laudos intactos; somente o índice ativo de passos foi excluído da comparação por integrar esta atualização. Testes anteriores não foram editados; links, linhagem e diff foram conferidos. Nenhum registro ou hash histórico foi recalculado, e módulos selados mantêm seus próprios commits e recibos.
