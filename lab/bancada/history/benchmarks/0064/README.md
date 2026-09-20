# Experimento0064 — Analisador e parâmetros

Contrato ../../passos/0064-identificar-analisador-e-parametros.md. A/B com contextos novos, filesystem compartilhado; sem atestação de isolamento ou assinatura de autenticidade. Manifesto e baseline congelados antes de ambos os autores. Os freezes de implementação/testes precedem execução.

O script documental run-demo.py.txt recebe o binário Bateia e requer executions/ ausente. Usa a fonte controlada ../0063/reference.go.txt, sem executar projetos externos. Grava recibos de mesma fonte em caminhos distintos, seed diferente e fonte divergente do hash esperado. Confere o fingerprint em Python conforme protocolo público e confronta a identidade do analisador com hash independente do binário CLI.

trace.analyzer identifica o Bateia; execution.binary_sha256 identifica a candidata instrumentada. O commit de build no trace pertence ao analisador e não à candidata. (devel), versão vazia e vcs_modified null não devem ser substituídos por dados inventados. A identidade via /proc/self/exe em Linux evita confundir arquivo substituído no disco com processo em execução.

trace.request_sha256 identifica contrato, produtor, fonte e parâmetros computacionais. Exclui caminhos/horários/resultado/hash esperado e não é checksum do recibo inteiro. Mesmo fingerprint não garante resultados determinísticos nem prova execução autêntica. Hash diferente não explica sozinho qual parâmetro mudou: conferir trace.inputs. recorded rastreia pedido, incluindo pedido que terminou unknown; não promove veredito.

Regressões anteriores mantêm expectativas congeladas. Comandos/logs/recibos finais serão registrados em verification.json. Tempos são telemetria da ferramenta.

Concluído: 4/4 cenários novos + 4/4 regressões0063 + 8/8 regressões0062, testes completos/build/lint aprovados. Produção sem revisões; teste B revisado uma vez para distinguir casos planejados de executados, conforme contrato anterior. Falha e snapshots preservados.
