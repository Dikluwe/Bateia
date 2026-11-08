# 0143 — Vizinhanças sob demanda

Contrato: [passo0143](../../../../00_nucleo/passos/0143-expandir-vizinhancas-a-partir-da-memoria.md). Continuação por relatório e memória, sem nova pergunta nem consultas de catálogo.

```sh
go build -o /tmp/bateia ./04_wiring
/tmp/bateia idea insights expand --input /tmp/investigacao.json \
  --focus https://github.com/org/repo --follow-links 2 \
  --output /tmp/trilha-1.json
/tmp/bateia idea insights expand --input /tmp/trilha-1.json \
  --focus https://github.com/org/repo --output /tmp/trilha-2.json
```

Sem foco, continuar pendências de todas as trilhas; o foco pode ser URL canônica ou entry-id conhecido. Se o foco esgotar, não ler outra trilha. `--memory feedback.json` permite aplicar descartes/reaberturas offline, conservando o relatório que contém a intenção e os documentos congelados. Os overrides disponíveis são follow-links, max-passages, timeout e attempt-timeout. Não há flags de pergunta, termos, fontes ou busca.

A entrada precisa de selo de relatório e memória íntegros. Documentos anteriores são validados contra recibos e fingerprints antes de derivar referências. Não são copiados para cartões atuais. A fronteira é fixada antes de efeitos; um README lido agora só fornece novos destinos na próxima chamada. Cada resultado mantém memória, followups atuais, tarefas e neighborhood com origem/foco/fronteira. Pendências sem citação podem ser tentadas sem inventar lead_id. As quotas e leitores por forge são os mesmos de0142.

`FRONTIER_EXHAUSTED` significa ausência de pendências conhecidas no foco; `FRONTIER_PENDING` indica pendências com leituras desativadas. Falhas/omissões ficam explícitas; falha de todas as leituras publica FAILED e retorna erro. Cancelamento/timeout preservam relatório. Esgotamento é limitado aos caminhos conhecidos, ao primeiro testemunho da memória e às quotas de extração; não demonstra que todas as abordagens foram exploradas. Normalização de aliases como `.git` continua uma limitação de identidade do recorte anterior.

## Controles

Testes L1: raiz→filho→neto entre chamadas, nenhuma consulta de catálogo, imobilidade da entrada, foco esgotado sem vazamento, ciclos, descartes/reabertura, pendência sem lead, namespace GitLab, seleção de fontes, adulteração de selo/conteúdo/fingerprint, geração, quotas e cancelamento.
Testes CLI: cadeia por arquivos, feedback, aliases, JSON/flags antes de factory, ausência de nova pergunta, segredos e falha com publicação. L4 verifica despacho antes do ramo offline. Fixtures desses testes são sintéticas e não comprovam disponibilidade pública.

O relatório público real da pesquisa anterior é usado separadamente para testar a admissão da memória. `real-input-dry-run.json` é o primeiro controle local anterior ao gate final, com follow-links=0: não coletou conteúdo remoto nem comprovou disponibilidade do backend. Novas execuções e recibos posteriores são preservados em arquivos próprios.

`baseline.json` registra arquivos anteriores à mudança. Logs/recibos da pilha ficam em `final/`; o verificador do recorte e seus hashes são registrados nesta pasta após os controles. Nenhuma fonte de candidato é executada. Novidade permanece NOT_ASSESSED e efeito UNVERIFIED.

## Resultado verificado

`final/verification.json`: lint, auditoria arquitetural, go test ./..., go vet ./..., build e diff passaram. A auditoria cobre391 arquivos produtivos. `first-gate/` preserva a execução com duas rejeições V2 do lint por ausência de testes nos arquivos novos de domínio/fronteira; os demais gates passaram nessa execução. Após acrescentar/organizar os testes por módulo, a pilha completa passou. A primeira revisão automática de permissão expirou antes da execução; a repetição foi autorizada.

Seis controles históricos (`replay-0106/0134/0135/0136/0137/0138.json`) passaram. `workflow-0142-replay.json` reproduz exatamente os dois relatórios e o feedback do workflow anterior. Seus filhos continuam dados sintéticos da fixture histórica, separados do controle real abaixo. Os29 arquivos históricos de0142 e todos os testes preexistentes permanecem inalterados.

`live-step-1.json` retoma a memória real da pesquisa anterior, com foco em vitality-vis/rest-api, uma tentativa e nenhuma busca. Em3,151s coletou o README no commit f421ef38c28ed8a4248612e7f6bc6a35e05f4f34. O resultado é CONTEXT_PREPARED, calls=0, com um cartão atual; pergunta herdada e nenhum candidato/catálogo antigo apresentado como coleta nova. O README declara busca por múltiplas sementes, RRF e endpoints de citações/UMAP; essas declarações documentais não foram verificadas no código nem executadas.

`live-step-2.json` usa o resultado anterior como entrada e mantém o mesmo foco: FRONTIER_EXHAUSTED, zero leituras e zero buscas, apesar de ainda haver pendências em outras trilhas. Esse controle demonstra contenção do foco e encerramento observável; não demonstra completude universal da investigação.

`live-receipt.json` registra comando sem credenciais, hashes, duração e status. `candidate-bateia`, usado nos controles, tem o mesmo SHA256 do binário do gate final. `verification.json` e `lineage.json` consolidam os controles e a causalidade. Revalidar sem rede: `python3 lab/bancada/insights/0143/verify.py`.
