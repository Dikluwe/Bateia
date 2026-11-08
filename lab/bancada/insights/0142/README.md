# 0142 — Aprofundamento automático e memória de exploração

O modo multifonte aprofunda até2 trilhas por padrão. followups separa identidade, origem, recibo, bundle e cartões. exploration_memory registra PENDING/EXAMINED/DISCARDED, motivo e fingerprint da última leitura. Tarefas distinguem cartões atuais de memory_entry_refs.

## Uso

```sh
/tmp/bateia idea insights discover \
  --question 'serendipitous software discovery' --sources all \
  --follow-links 2 --output /tmp/primeira.json
/tmp/bateia idea insights discover \
  --question 'serendipitous software discovery' --sources all \
  --memory /tmp/primeira.json --output /tmp/segunda.json
```

O primeiro resultado já inclui memória. A segunda rodada pula examinados/descartados e prioriza outros pendentes ou desconhecidos. O usuário/IA escolhe o snapshot com --memory: não há arquivo global oculto nem segunda escrita separada do relatório. Nesta versão, memória cobre repositórios das cinco forges com leitor; outros catálogos permanecem fora desse histórico. --memory aceita snapshot standalone ou relatório novo que contenha exploration_memory.

```sh
/tmp/bateia idea insights memory --input /tmp/segunda.json \
  --target https://github.com/org/repo --state discarded \
  --reason 'Este caminho não atende à restrição atual' \
  --output /tmp/memoria.json
```

O alvo precisa existir na memória. Use state pending para reabrir, também com motivo. O comando é offline, protege a entrada e publica saída atomicamente com a proteção de segredos existente. EXAMINED só decorre de coleta íntegra; feedback não simula leitura. Para continuar, use --memory /tmp/memoria.json em uma descoberta posterior e salve em novo caminho.

## Limites e estados

--follow-links aceita0–4;0 desativa. Somente um salto e fontes selecionadas; demais ficam pendentes com SOURCE_NOT_SELECTED/READER_UNAVAILABLE. Identidades são deduplicadas por sítio/external_id. GitHub usa adaptador dedicado sem fallback raw; outras forges reutilizam leitores existentes. Quatro GETs máximos por destino, com revisão fixada e metadata pública. max-calls continua quota de catálogo; leituras adicionais seguem deadlines total/por leitura.

32KiB/documento,128KiB agregados,8 documentos e max-passages são compartilhados com a coleta primária. Até16 destinos entram no plano de followups; demais ficam na memória com lacuna de limitação. Falha mantém PENDING e registra tentativa/status; cancelamento conserva resultado parcial. Catálogo indisponível não impede progresso de pendências quando o leitor retorna conteúdo válido.

Até256 identidades/1MiB, com margem de4KiB nas atualizações geradas. MEMORY_CAPACITY informa atualizações que não couberam, conservando estado anterior. Memória é resumo do estado/última justificativa, não cache de README ou auditoria completa. Hash verifica integridade, não autoria do histórico importado. Leituras anteriores não geram cartões atuais. EXAMINED significa documento coletado, não mecanismo/benefício comprovado. Não há atualização automática por tempo: para rever nova versão de um examinado, reabra o caminho com motivo.

## Controle reproduzível offline

```sh
GOCACHE=/tmp/bateia-review-go-cache go run ./lab/bancada/insights/0142/replay > /tmp/workflow-fixture.json
```

workflow-fixture.json usa dois READMEs públicos de origem congelados do ensaio anterior e documentos de destino AUTORAIS/SINTÉTICOS. Não são observações públicas atuais dos destinos; não use memória de fixture como histórico real. A primeira rodada aprofunda dois destinos, incluindo Serendip-o-matic; a segunda pula pais/destinos examinados e lê outros dois. Estado:14 entradas, primeiro4 examinadas/10pendentes, depois6examinadas/8pendentes. Um feedback descarta outro caminho com motivo. Esse ensaio comprova fluxo de memória e origem, não insights dos destinos reais.

Testes de transporte GitHub usam RoundTripper controlado para metadata pública, identidade/SHA, headers, quatro GETs, erros e blob/tamanho/symlink/redirect. Testes puros cobrem validação/imutabilidade, duas rodadas sem repetição, seleção orientada pela memória, fallback de catálogo, quotas e cancelamento. CLI testa relatório/standalone, alias, validação antes de factory e segredo em motivo. Gate: make verify com loopback para fixtures gerais. Evidências antigas permanecem intactas; novos logs/replays/recibos estão nesta pasta.

`verification.json` consolida os controles e `lineage.json` registra hashes/prompts. Conferência offline: `python3 lab/bancada/insights/0142/verify.py`. O binário final exercitou descarte/reabertura em `cli-workflow-check.json`, usando apenas dados sintéticos da fixture. A primeira revisão automática de permissão para o gate expirou antes de execução; a repetição com loopback foi aprovada e passou.
