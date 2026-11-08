# Parecer cego

Os dois pacotes sustentam as mesmas duas políticas para compromissos sobre peças: registro livre com pedidos provisórios e conciliação, ou custódia local com um único ponto autorizado a confirmar. X também apresenta uma diferença relevante de responsabilidade de implementação: adaptar um aplicativo Android de formulários, em vez de desenvolver o cliente. Essa opção é exclusiva observada nesta execução, em dimensão parcial; já constava da intenção de seu baseline. Não há família exclusiva nova, acrescentada pela pesquisa, demonstrada pelos pacotes.

Foram usados somente os quatro arquivos autorizados. Trechos e linhas foram aceitos como verificados em Git fixado. Nenhum projeto foi executado. Notas são ordinais e não foram somadas.

## Propostas

| Proposta | Relevância 0–4 | Evidência 0–3 | Especificidade 0–2 | Candidata relevante |
|---|---:|---:|---:|---|
| X.A1 | 4 | 3 | 2 | Sim |
| X.A2 | 4 | 2 | 2 | Sim |
| X.A3 | 4 | 3 | 2 | Sim |
| Y.A | 4 | 3 | 2 | Sim |
| Y.B | 4 | 2 | 2 | Sim |

### X.A1

Captura por aplicativo Android de formulários existente, ocorrências independentes e conciliação própria; confirmação física antes de viagem. Define quem revisa, o que fica pendente, o que desenvolver e como rejeitar o piloto. Documentação específica apoia captura offline, encontro por entidade, auditoria opcional, CSV e limites de desligamento/manutenção.

Trechos usados:

- X E1: docs/collect-intro.rst:13–17, designed to work well offline
- X E2: docs/entities-intro.rst:151–153, Entity List that represents encounters
- X E3: docs/form-audit-log.rst:91, old answers and new answers
- X E9: docs/central-install.rst:48–51, infrastructure expenses and ongoing maintenance
- X E17: docs/form-audit-log.rst:266–268, ausência de entrada de tela no desligamento

Restrições não resolvidas: Projeção de estoque, fila de revisão e protocolo de confirmação próprios; Exportação completa de originais/fotos e restauração não demonstradas pelos recortes; 72h, Android barato, 45 mil peças, custos e carga humana não medidos.

Baseline: **mesma família** (B1, B4). Auditoria opcional e suporte de formulários tornam a implementação mais concreta. B1 já exigia ocorrências append-only e confirmação; B4 já propunha adaptar aplicativo de campo. A exclusividade de responsabilidade abaixo não foi introduzida pela pesquisa.

### X.A2

Um ponto local de custódia confirma reservas/uso; demais aparelhos só observam/pedem; trânsito e entrega transferem autoridade sem timeout. Muda concorrência permitida e responsabilidade física, mantém captura livre e explicita gargalo e perda de aparelho. Fontes dão primitivas/limites, não demonstram exclusividade ou handover seguro.

Trechos usados:

- X E8: docs/entities-intro.rst:187–192, control which Entities are delivered
- X E12: src/docs/src/replication/conflicts.rst:111–115, single node e 409
- X A2 work_change/software_change: custódio, transferência e não restaurar autoridade por timeout

Restrições não resolvidas: Protocolo de autoridade/recuperação sem implementação ou prova; Disciplina de retirada e disponibilidade do telefone; Mesmas lacunas de captura/exportação de A1.

Baseline: **mesma família** (B2). B2 já continha custódio/dispositivo único e handoff explícito. Limites de filtro/409 e cenário de falha refinam esse desenho; não criam família.

### X.A3

Cliente e projeções próprios, réplicas colaborativas com divergências visíveis, pedidos provisórios ou custódia para compromisso firme. Explicita custo de construir cliente e UI de conflito e separa histórico de negócio das revisões descartáveis; tem piloto de integração e carga.

Trechos usados:

- X E10: src/docs/src/replication/conflicts.rst:90–105, arbitrary revision as winner
- X E11: src/docs/src/replication/conflicts.rst:185–189, non-leaf bodies discarded
- X E13: rust/automerge/src/lib.rs:58–74, cannot merge sensibly
- X E14: javascript/src/implementation.ts:1005–1012, concurrent put operations

Restrições não resolvidas: Android, autorização, transporte de fotos e exportador próprios; Memória, armazenamento, custo e prazo de dois generalistas não medidos; Convergência não garante reserva única.

Baseline: **revisão material** (B3). A fonte de compactação exige manter eventos vivos separados das revisões/histórico do motor. A família colaborativa, preservação de conflitos e separação de reservas já constavam em B3; não são descobertas novas.

### Y.A

Aplicativo próprio com eventos locais replicados, revisão explícita e confirmação física antes de promessa; Automerge é opção interna. Define identidade/idempotência, visibilidade de conflitos, papel do steward e entrega de saída de dados; distingue primitivas de exportação de exportador e mede carga em Android.

Trechos usados:

- Y E1: docs/guides/conflicts.md:142–148, only create new documents
- Y E2: docs/guides/conflicts.md:74–90, arbitrary winner
- Y E4: docs/_includes/api/compaction.html:3–7, removing non-leaf revisions and attachments
- Y E6: docs/_includes/api/replication.html:29–33, low amounts of memory
- Y E8/E10: get_all e teste de ambos os valores; teste lido, não executado
- Y E9: rust/automerge/src/sync.rs:1–10, reliable in-order stream

Restrições não resolvidas: Retenção Android, app shell, fotos, autenticação e custos sem medição; Exportação/restauração integral a implementar; Não garante reserva única entre escritores isolados.

Baseline: **revisão material** (b1, b4). A regra de compactação torna explícito que eventos e fotos devem permanecer documentos vivos, sem depender de revisões antigas; é revisão de retenção/armazenamento, não nova família. b1 já dizia append-only e revisão de conflitos, b4 já tinha cache/fila/reconciliação.

### Y.B

Telefone compartilhado da oficina é único escritor de compromissos, retirada física controlada e passagem durável de autoridade; observações pessoais livres. Distingue disponibilidade de registrar de disponibilidade de comprometer. Explicita perda de confirmação, bloqueio/recontagem e ausência de failover por timeout. Primitivas citadas não validam o protocolo.

Trechos usados:

- Y E5: docs/guides/conflicts.md:22–24, 409
- Y E2: docs/guides/conflicts.md:74–90, eventual deterministic winner
- Y E1: docs/guides/conflicts.md:142–148, immutable documents
- Y relatório B: aparelho único, trânsito indisponível e handover durável

Restrições não resolvidas: Ownership/handover não implementados nem formalmente verificados; Disciplina física e falha do aparelho podem bloquear compromissos; Sem orçamento ou resultados de campo.

Baseline: **revisão material** (b2, b3). b2 já tinha direitos transferíveis por oficina. A substituição da configuração de b3 (servidores locais mantidos) por telefone de autoridade portátil muda responsabilidade operacional e falha; a relação causal exclusiva com uma fonte é incerta, pois 409 não determina escolha de hardware.

A nota 2 de evidência das custódias não significa que sejam vagas ou irrelevantes. Elas mudam explicitamente quem pode confirmar e o que acontece numa falha. Entretanto, as fontes apoiam captura/eventos e a diferença entre conflito local e eventual; não oferecem o protocolo de autoridade que daria a garantia proposta. A nota 3 das demais usa documentação específica de mecanismos e limites; não certifica a aplicação completa.

## Equivalência por decisão e garantia

### G1 — Observações livres offline, preservação de divergências, pedidos provisórios e revisão/confirmacão física antes de promessa

X.A1, X.A3, Y.A. Equivalentes quanto à concorrência e garantia de compromisso; não são três garantias diferentes.

- X.A1 reaproveita aplicativo de coleta Android e deixa projeção/conciliação próprias; X.A3/Y.A exigem construir cliente e fluxo de conflito. Isso altera responsabilidade de manutenção, portanto merece comparação como configuração interna relevante.
- X.A3 explicita descrições lado a lado; Y.A também exige que conflitos não desapareçam atrás de vencedor e inclui Automerge internamente. Não há exclusiva de colaboração demonstrada.
- Y.A especifica transporte confiável, lotes/memória e eventos vivos; X.A3 também trata compactação e eventos imutáveis. São diferenças de profundidade/integração, sem nova garantia de reserva.

### G2 — Custódia local restringe quem confirma compromissos; captura livre continua; transferência não duplica autoridade

X.A2, Y.B. Equivalentes quanto ao escritor único, gargalo, controle físico, trânsito e rejeição de failover por timeout.

- Y detalha geração/transfer ID, aceite durável, recontagem e limiar de 5%/5 minutos; X começa por simulação antes de construir e usa limite de um turno. São diferenças verificáveis de avaliação e detalhamento, sem outra garantia demonstrada.
- O telefone de autoridade em lugar de servidores locais é revisão material do baseline de Y, mas X já oferece o mesmo telefone de balcão.

A distinção entre G1 e G2 muda efetivamente a decisão: em G1, qualquer voluntário registra pedidos concorrentes e a confirmação chega depois; em G2, apenas a custódia pode confirmar, inclusive offline, com perda dessa capacidade se o aparelho/ponto de liberação falhar. O registro de observações permanece disponível nos dois. Essa diferença não se reduz a uma troca de biblioteca.

## Exclusividade e contribuição da pesquisa

**Exclusiva qualificada, parcial: X.A1 como divisão de responsabilidade.** X propõe reutilizar captura Android, formulários e auditoria de um aplicativo mantido, deixando projeção de inventário e conciliação aos dois desenvolvedores. Y.A exige construir cliente e app shell, e reconhece não ter investigado uma suíte de formulários Android offline. Isso pode mudar a escolha sob a restrição de manutenção. O apoio vem de X E1/E2/E3, com custo operacional e limite de auditoria em E9/E17. Não é exclusividade de preservação de observações ou de reserva segura, ambas compartilhadas. Se a unidade for exclusivamente política de compromisso, não há exclusiva.

**Contribuição da pesquisa para essa exclusividade: não demonstrada.** X B4 já propunha adaptar aplicativo de campo offline/exportável e B1 já previa ocorrências, conciliação e confirmação. Encontrar mecanismo concreto e seus limites melhora a fundamentação, mas não acrescenta por si essa escolha de responsabilidade. Nenhuma família nova é demonstrada em X ou Y.

Há revisões materiais internas: X.A3 e Y.A tornam a retenção independente de revisões compactáveis; Y.B troca a configuração de servidores de b3 por telefone de autoridade. No primeiro caso, append-only já estava no baseline, portanto a novidade é o contrato explícito de armazenamento, não a ideia geral de preservar relatos. No segundo, a consequência operacional é real, mas sua atribuição a uma fonte é incerta: 409 não implica telefone em vez de servidor. Nenhuma dessas revisões fica sem equivalente no outro pacote.

Exclusivas apenas aparentes:

- Três títulos em X contra dois em Y: X.A3 corresponde ao cliente replicado de Y.A. Não constitui ganho automático de cobertura.
- Colaboração/CRDT só em X: Y descreve Automerge no corpo, inclui leitura de valores concorrentes, teste estático e limite de transporte; a opção interna deve ser considerada.
- Telefone de autoridade ou proteção contra compactação só em Y: X.A2 e X.A3 também os descrevem.

Descartes são coerentes: exigir visita diária/online viola o problema; LWW/vencedor técnico não decide disponibilidade física; revisões compactáveis não são auditoria permanente; trocar biblioteca não cria nova política. Catálogos de inventário sem fonte útil em X e servidores locais não investigados em Y ficam como lacunas, não como candidatas sustentadas. Não se atribui alternativa adicional a nomes comerciais ou a uma opção só mencionada.

## Incertezas e limites

- Nenhum projeto, teste de fonte ou piloto foi executado; notas avaliam proposta e evidência estática, não utilidade comprovada.
- Literalidade/linhas fixadas foram aceitas como verificadas conforme instrução; não houve consulta a outras fontes.
- Recortes X E5, E6, E7 e E18 terminam antes de parte das propriedades atribuídas. CSV básico, conflito e estado local estão apoiados; anexos/audit integral, LWW, espera de cinco dias e ausência de mídia em Entity Lists não ficam integralmente demonstrados nesses trechos. Não usei essas extensões como garantia.
- Recortes X E13/E14 não exibem completamente o acesso a valores/histórico alegado; Y E8/E10 fornece apoio mais explícito ao acesso a conflitos. Isso não cria alternativa nova.
- Ambos deixam autenticação/revogação offline, restauração, 45 mil peças/fotos, custos recorrentes e disciplina real com voluntários sem validação.
- Atribuição causal à pesquisa é limitada pelos baselines curtos: refinamentos de retenção são materiais no desenho explícito, mas append-only já era intenção de ambos. Não há evidência suficiente para atribuir uma nova família exclusiva à pesquisa.
- Ausência de alternativa em Y ou X não demonstra incapacidade de um produto; nenhuma inferência de identidade ou comparação de desempenho dos produtos foi feita.

Os pilotos são concretos e podem rejeitar os desenhos: 72 horas, oficinas/voluntários reais, reservas rivais, reinício, entrega duplicada/fora de ordem, fotos e exportação/restauração. X começa a custódia com simulação de mesa e mede integração; Y detalha falha de autoridade e separa disponibilidade de registro e compromisso. Os limiares propostos são hipóteses a pactuar. Nenhum resultado de campo permite escolher vencedor, estimar redução de viagens ou afirmar que dois generalistas conseguirão manter o sistema.
