# 0143 — Expandir vizinhanças a partir da memória

**Data:** 2026-10-28. **Contrato ex-ante.** A1-INV sob Tekt. L1 possui admissão, fronteira, foco e execução limitada; L2 transporta relatório/flags e publica; leitores L3 existentes; L4 injeta a factory existente. Preservar descoberta, schemas e hashes anteriores.

## Jornada

`idea insights expand --input anterior.json [--focus URL-ou-entry-id] [--memory feedback.json] [--follow-links 0..4] [--max-passages 1..16] [--timeout duração] [--attempt-timeout duração] --output proximo.json`.
Não recebe pergunta, fontes ou termos novos. Herdar pergunta, familiaridade, sítios e limites do relatório anterior; overrides só de orçamento. Zero consultas de catálogo. Sem estado global, modelo, clone ou execução de candidatos. Memória standalone opcional aplica feedback validado; não substitui o relatório que fornece pergunta e documentos congelados.

Admitir relatório multifonte com memória e selo site_discovery_ íntegros antes de factory/efeitos; cópia profunda. Derivar referências apenas dos READMEs congelados e íntegros da entrada, incluindo followups, com correspondência de identidade/commit/path/hash na memória original. Até8 documentos/128KiB; validar recibos/bundles antes de derivação. Bytes anteriores não aparecem como leitura atual, não consomem orçamento atual e não promovem estado.

Foco opcional identifica entrada conhecida por URL canônica ou entry-id. Rejeitar desconhecido/descartado. Sem foco, continuar pendências conhecidas; com foco, ler apenas o próprio foco pendente ou descendentes alcançáveis pelas citações guardadas/recém-derivadas. Não consumir outras trilhas quando a escolhida estiver vazia. Travessia iterativa limitada às256 identidades; ciclos não causam releitura; descartados interrompem a travessia. Memória conserva apenas primeiro testemunho por identidade: caminhos alternativos não recuperáveis pelo resumo são limitação explícita.

Selecionar snapshot da fronteira antes de ler: um salto por chamada; links dos documentos lidos agora só serão extraídos da entrada da próxima chamada. EXAMINED/DISCARDED não geram leitura automática. Pendências sem lead (falha primária) podem ser tentadas; lead_id vazio não fabrica citação. Reutilizar ordenação, leitores por sítio, deadlines e quotas de0142 (até4 tentativas,32KiB/documento,128KiB/8 documentos, max-passages global). Falhas continuam pendentes; omissões são explícitas.

## Saída

Mesmo insight-site-discovery-v1, extensão opcional neighborhood com ID da entrada, foco e IDs da fronteira. Plano declara agenda memory-neighborhood-v1; calls=0, attempts/candidates/readmes/bundles/site_contexts vazios; followups somente leituras atuais. Referências derivadas testemunham bytes da entrada anterior. Tarefas separam cartões atuais, referências e memória. Geração avança uma vez; snapshot único publicado atomicamente sem sobrescrever input/memory/aliases. Integridade não autentica história importada.

Estados: CONTEXT_PREPARED com leitura; PARTIAL com falhas/omissões; FAILED com erro se todas as tentativas efetivas falharam; FRONTIER_PENDING com pendências sem tentativa; FRONTIER_EXHAUSTED sem pendências no foco. Timeout/cancelamento preservam relatório. Novidade NOT_ASSESSED e efeito UNVERIFIED.

## Verificação

Testes puros de três passos raiz→filho→neto, retomada sem catálogo/nova pergunta, foco sem vazamento, esgotamento, ciclos/descartes, retry, quotas, adulteração e cancelamento. CLI: JSON/flags antes da factory, feedback, saída/alias/segredos e retomada encadeada. Composição do despacho e replay dos formatos anteriores. Gates completos; evidências em lab/bancada/insights/0143.
