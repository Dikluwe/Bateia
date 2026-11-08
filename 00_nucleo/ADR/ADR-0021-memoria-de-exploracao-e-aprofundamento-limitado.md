# ADR-0021 — Memória de exploração e aprofundamento limitado

**Status:** ACEITO PARA O RECORTE0142. **Data:** 2026-10-28. Especializa ADR-0019/0020 e o [passo0142](../passos/0142-aprofundar-trilhas-e-conservar-memoria.md).

L1 possui snapshot canônico, estados e quotas; L2 lê snapshots e publica resultado protegido; L3 implementa leitores; L4 compõe. O resultado inclui exploration_memory. --memory recebe um snapshot standalone ou de relatório anterior; não há memória global oculta nem segunda escrita que possa divergir da publicação atômica do relatório.

A primeira versão registra caminhos de repositório das cinco forges com leitor (GitHub, Gitee, AtomGit/gitcode, Codeberg e GitLab); outros catálogos permanecem fora desta memória. PENDING significa caminho conhecido ainda não encerrado; EXAMINED exige coleta íntegra; DISCARDED exige decisão/motivo. Histórico importado é declarado: hash verifica integridade, não autenticidade. Fingerprints anteriores não viram bundles/cartões atuais. O estado é um resumo, não auditoria completa. Revisões posteriores de um examinado precisam de reabertura explícita, sem revalidação automática por tempo.

Até2 destinos por padrão (0–4 configurável), em um salto, somente nos sítios selecionados, com leitores próprios. Leituras compartilham128KiB,8 documentos e max-passages; falhas conservam pendências/parciais. Destinos não viram candidatos de catálogo. O leitor de aprofundamento GitHub é dedicado, sem fallback raw, com metadata pública, SHA fixado e respostas limitadas. Outros sítios usam seus leitores imutáveis existentes.

Memória até256 identidades/1MiB, com margem de4KiB nas atualizações geradas e omissões explícitas. Validar snapshot, alias e flags antes de rede. Feedback offline descarta/reabre somente alvos conhecidos e exige motivo; não fabrica estado examinado. Identidades pendentes/desconhecidas têm precedência na seleção, e histórico encerrado evita releituras automáticas. Novidade e efeito permanecem não avaliados/verificados.

## Continuação sob demanda —0143

O [passo0143](../passos/0143-expandir-vizinhancas-a-partir-da-memoria.md) acrescenta `insights expand`: herda pergunta e fontes do relatório selado, deriva referências dos documentos congelados da entrada e continua pendências sem catálogo. Foco limita a trilha por citações, com travessia finita e sem releitura automática. Um salto por chamada permanece: descendentes lidos agora só alimentam a próxima chamada. Bytes anteriores não são evidência atual; o resultado opcional neighborhood descreve a expansão sem alterar serialização anterior. Este recorte não muda o comportamento de discover.
