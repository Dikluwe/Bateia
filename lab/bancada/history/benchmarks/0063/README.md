# Experimento 0063 — Identidade da fonte e declaração executada

Contrato: ../../passos/0063-vincular-fonte-ao-efeito.md. Manifesto/baseline congelados antes de A/B. Regime A/B sem atestação física: implementação e testes com contextos separados, mesmo filesystem. Hashes identificam artefatos, não provam isolamento.

reference.go.txt e helper-mutated.go.txt são fontes documentais. A função Transform é idêntica em ambas; somente seu helper muda, passando a alterar a entrada. A demonstração usa o mesmo caminho temporário para observar uma mudança entre execuções.

Após build, o script documental run-demo.py.txt recebe o binário e grava executions/ (que precisa não existir). Ele usa /tmp/bateia0063-demo/candidate.go e o oráculo congelado do passo0062, restaurável a partir de ../0062/gate-export.json. Produz recibos de referência, hash antigo, hash atualizado e execução sem vínculo solicitado. Não executar código de terceiros para reproduzir: apenas as fontes controladas acima.

Uso público:

```sh
go build -o /tmp/bateia-0063 ./04_wiring
sha256sum /caminho/candidate.go
/tmp/bateia-0063 effects verify --source /caminho/candidate.go --symbol Transform --expect-sha256 HASH_OBTIDO --output /tmp/efeito.json
```

O hash esperado deve vir da versão previamente examinada. Recalculá-lo cegamente antes de cada execução apenas aceita o estado atual e não protege a ligação com evidência anterior. Conferir source_binding.status, declaration, status global, cases e execution.stop_reason.

matched é identidade de arquivo; a função pode falhar, estar fora do recorte ou não executar. declaration representa localização AST, não prova de efeito. supported_in_scope continua restrito às entradas observadas; o hash não prova autoria, commit, correção geral ou sandbox. Recorte de funções Go autônomas permanece o mesmo do passo0062.

Concluído: 4/4 cenários da CLI e 8/8 regressões0062, testes novos/suíte completa/build/lint aprovados. Uma correção focal de reconhecimento de diretivas, mantendo contrato e testes intactos. Logs e revisão preservados; resultado em verification.json.
