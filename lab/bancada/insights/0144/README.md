# 0144 — Foco contínuo e próximos vizinhos

Contrato: [passo0144](../../../../00_nucleo/passos/0144-continuar-foco-e-selecionar-vizinhos.md).

```sh
go build -o /tmp/bateia ./04_wiring
/tmp/bateia idea insights expand --input pesquisa.json \
  --focus https://github.com/org/repo --follow-links 2 --output trilha-1.json
/tmp/bateia idea insights expand --input trilha-1.json --output trilha-2.json
/tmp/bateia idea insights expand --input trilha-2.json --all-trails --output geral.json
```

A segunda chamada herda o foco, pergunta e limites. --all-trails volta explicitamente à fronteira geral; um foco novo troca a trilha. A combinação all-trails/foco é inválida antes da factory. Cada chamada continua limitada a um salto e nunca pesquisa catálogos.

neighbors é uma projeção anterior às leituras: entry-id, identidade, distância opcional, tentativas anteriores, selection_status e relação com outra entrada quando há supressão. FRONTIER significa inclusão no plano, sem afirmar sucesso ou disponibilidade do leitor. followups registra o resultado atual. Distância só existe com foco e mede citações conhecidas; não é nota de relevância/novidade. A fila prioriza leitor elegível, distância, tentativas, geração, referência e ID.

Somente pares X/X.git já conhecidos no mesmo host/conector são relacionados. Não há mesclagem de memória ou inferência sobre renomes, maiúsculas, sites ou um destino X ainda desconhecido. DUPLICATE_ALIAS evita uma segunda tentativa; ALIAS_ALREADY_EXAMINED/ALIAS_DISCARDED preservam decisões. ALIAS_STATE_CONFLICT exige resolver decisões explícitas antes de ler. Reabertura canônica com motivo é respeitada quando a outra entrada só informa exame anterior. Os testemunhos, estados e contadores suprimidos permanecem intactos.

baseline.json conserva hashes anteriores. before/ guarda os testes anteriores que recebem apenas novos casos. Evidências novas ficam nesta pasta; as pastas históricas0142/0143 são preservadas. Fixtures unitárias são sintéticas; controles com o relatório público anterior são identificados separadamente. Nenhum código de candidato é executado. Novidade/efeito permanecem não avaliados/verificados.

## Resultado verificado

Todos os gates em final/verification.json passaram: lint, auditoria (392 arquivos produtivos), testes completos, vet, build e diff. Os corpos dos testes anteriores foram preservados; somente novos casos foram acrescentados nos arquivos de domínio/CLI. O novo módulo de seleção possui seu próprio teste e prompt. verification.json/lineage.json registram hashes e controles; revalidar sem rede com `python3 lab/bancada/insights/0144/verify.py`.

Oito controles históricos passaram: seis replays anteriores, igualdade dos dois relatórios/feedback da fixture0142 e admissão do relatório selado0143 pelo novo executável. Os78 arquivos de evidência anteriores no baseline permanecem intactos. A fixture0142 continua sintética; não confundir seus destinos com o controle público a seguir.

live-step-1.json reutiliza a memória real da pesquisa de28/10, com foco em vitality-vis/frontend e orçamento de uma tentativa. neighbors lista rest-api como FRONTIER, distância1, e frontend.git como ALIAS_ALREADY_EXAMINED, distância0; somente o backend foi consultado. O README foi coletado em1,596s no commit f421ef38c28ed8a4248612e7f6bc6a35e05f4f34, com zero buscas. Essa duração é uma observação desta execução, sem comparação controlada de desempenho.

live-step-2.json foi produzido sem --focus: herdou o frontend, não fez leituras/buscas e retornou FRONTIER_EXHAUSTED. A memória conserva pendências de outras trilhas. live-general-preview.json usa --all-trails e follow-links=0: mostra11 destinos na fila geral, nenhum com distância inventada, mais o alias suprimido. Nenhuma pendência externa foi consumida ao encerrar o foco.

O alias preserva exatamente o estado PENDING, contadores e testemunho anteriores; sua supressão é uma decisão operacional da expansão, sem forjar leitura ou descarte. Os recibos e os hashes do conteúdo/cartão lido são verificados. candidate-bateia usado nos controles tem o mesmo SHA256 do binário final.
