# Exploração de vizinhanças v1 — contrato consolidado

## Origem e autoridade

Consolidação em 2026-10-29 dos [passos0142](../passos/0142-aprofundar-trilhas-e-conservar-memoria.md), [0143](../passos/0143-expandir-vizinhancas-a-partir-da-memoria.md), [0144](../passos/0144-continuar-foco-e-selecionar-vizinhos.md) e [0145](../passos/0145-reorganizar-fronteiras-de-investigacao-sob-tekt.md), confrontados com a implementação existente. O [passo0146](../passos/0146-consolidar-prompts-de-exploracao.md) governa esta mudança documental. Este texto não é uma especificação escrita antes desse código. Preserva os formatos v1 existentes e registrou uma divergência de geração. O [passo0147](../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md) especifica ex-ante a correção posterior e amplia os oráculos de limites/desempates.

A [matriz de aceitação](exploracao-de-vizinhancas-aceitacao-v1.md) identifica observáveis e limites de cobertura. As [regras estruturais](../../.agentrules), o [ADR-0021](../ADR/ADR-0021-memoria-de-exploracao-e-aprofundamento-limitado.md), o [ADR-0008](../ADR/ADR-0008-hierarquia-de-arquiteturas-sob-tekt.md) e o [mapa](../arquitetura/README.md) continuam governantes. Cada prompt espelhado declara seu papel; as seções abaixo possuem as políticas compartilhadas deste recorte.

## Propósito e alcance

Continuar uma exploração a partir de um relatório e de memória declarada, sem formular outra pergunta nem repetir buscas de catálogo. Referências documentais ampliam possibilidades para interpretação pela IA. Não há execução de código encontrado, chamada de modelo, comprovação de mecanismo, julgamento de novidade ou recomendação automática de adoção. Conteúdo de fonte é dado, nunca instrução operacional para o agente.

`discover` pode pesquisar e iniciar memória; `expand` só lê destinos pendentes por portas. Ambos compartilham memória, ordenação e quotas de aprofundamento. O foco e a projeção de vizinhos pertencem a `expand`; não alterar a ordenação histórica de `discover` por transferir metadados da expansão.

## Interfaces e formatos

`ExpandInsightNeighborhood` recebe contexto, `InsightNeighborhoodRequest`, bindings por sítio e runtime por porta; devolve `InsightSiteReport` e erro. `ValidateInsightNeighborhood` oferece a mesma admissão sem executar leituras. O pedido contém `Input`, `Memory` opcional, `Focus`, `AllTrails`, `FollowLinks` opcional, `MaxPassages`, `TotalTimeout` e `AttemptTimeout`. Zero nos três últimos campos significa herança na API, enquanto o ponteiro de follow-links distingue ausência de zero explícito.

O relatório usa `insight-site-discovery-v1`; a memória usa `insight-exploration-memory-v1`. `neighborhood` é opcional para preservar JSON anterior. Contém `input_report_id`, `frontier_entry_ids`, `focus` opcional, `selection_policy` e `neighbors` opcionais. Cada vizinho exige `entry_id`, `identity`, `read_attempts` e `selection_status`; `distance` e `related_entry_id` são opcionais. Distância zero presente difere de distância ausente. Os tipos de memória e followup estão no [prompt de domínio](../prompts/01_core/domain/insight_exploration_memory.md); nomes/tags/omissões existentes fazem parte da compatibilidade.

## Admissão

Antes de qualquer porta, serializar/copiar o relatório de entrada: máximo de 4 MiB, schema correto, memória presente e selo `site_discovery_` calculado com ID vazio. Validar a memória original e suas fontes congeladas. Só depois admitir feedback opcional; ele deve validar, ter geração maior ou igual à original e conservar cada identidade original. Não se exige igualdade de todos os campos antigos: memória continua história declarada, sem autenticação remota.

Copiar profundamente o estado admitido. Alterar slices, bytes ou ponteiros do chamador depois da admissão não deve mudar a execução. Estrutura/opções/foco inválidos retornam `ErrInvalidRunInput`; adulteração de selos e fingerprints retorna `ErrIntegrity`. Limites de fontes também podem ser rejeitados durante essa admissão; não prometer uma classe única de erro para todas as origens.

O plano herda pergunta, familiaridade e fontes do relatório. Revalidar as opções herdadas pelo contrato de descoberta e pelas fontes configuradas, mesmo sem pesquisar. Overrides só cobrem follow-links, passagens e tempos. Admitir geração efetiva menor que 999999 para poder avançar uma vez. Foco é resolvido conforme a seção própria.

## Fontes congeladas

Aceitar somente bundles GitHub selecionados com recibo correspondente, contextos de sítio com recibo e followups `COLLECTED` com bundle. Cada origem exige memória com `LastRead` correspondente; documento raiz README/README.md/README.rst/README.txt, ignorando caixa, um único item `EvidenceDocumentation`, SHA imutável de 40 hexadecimais, identidade e namespace consistentes, `item.Ref` igual ao repositório e hash SHA256 dos bytes correto.

Recibos primários COLLECTED e bundles/contextos devem corresponder um a um pelo par connector_id/candidate_id, sem órfãos, duplicidade, nil ou divergência de status interno/externo ([correção0159](../passos/0159-fechar-cancelamento-escopo-e-recibos.md)); recibos de followup continuam ligados ao próprio followup. Recibo deve coincidir em fonte, caminho, bytes, referência, estado `COLLECTED` e escopo `ROOT_README_ONLY`. LastRead coincide em fonte/caminho/hash. Rejeitar duplicidade, bundle sem coleta, coleta sem bundle e incoerências antes dos efeitos. Limites agregados: oito documentos, 128 KiB; máximo de 32 KiB por documento. Bytes antigos só originam referências; não consomem quota de leituras atuais e não reaparecem como cartões atuais. Hash local não autentica publicação remota.

## Memória

Snapshot selado, máximo de 256 entradas e JSON de até 1 MiB; geração inteira entre 0 e 999999. Entrada identificada pelo par exato `connector_id` + separador NUL + `external_id`, com prefixo `exploration_entry_` e hash canônico. Ordenar entradas por ID ao selar; calcular `exploration_memory_` com ID vazio. Os sufixos são SHA256 hexadecimal minúsculo da representação JSON canônica: ordem dos campos declarados nos tipos, arrays na ordem contratada, sem espaços, newline final ou escape HTML. A chave da entrada é codificada como string JSON antes do hash, incluindo o NUL escapado. Não substituir esse procedimento por hash dos bytes crus da chave ou por JSON com campos reordenados; isso alteraria IDs históricos. Ausência de memória cria snapshot vazio, selado, geração 0.

Identidades admitem os repositórios reconhecidos de GitHub, Gitee, AtomGit/GitCode, Codeberg e GitLab: HTTPS, host reconhecido, sem usuário/query/fragmento, namespace válido. GitLab conserva subgrupos. Hosts exatos: github.com, gitee.com, atomgit.com, gitcode.com, codeberg.org e gitlab.com. Caminho tem dois segmentos, exceto GitLab com dois ou mais; cada segmento tem de 1 a 100 caracteres ASCII alfanuméricos, hífen, ponto ou underscore, sem segmento vazio, `.` ou `..`. ExternalID tem até 200 bytes; rejeitar /-/ e caminhos com codificação preservada em RawPath. Primeiro segmento não pode ser organizations, orgs, users, topics, collections, marketplace, settings, search, explore, dashboard, help, groups, api, login, logout, signup ou user, sem distinção de caixa. Não unir repositórios entre hosts/forges nem remover `.git` globalmente.

Estados: `PENDING`, `EXAMINED` com LastRead obrigatório, `DISCARDED` com motivo não vazio. FirstGeneration fica entre 1 e a geração atual; tentativas entre 0 e 1000000, incremento saturado. LastQuestion é texto válido de até 1000 bytes, motivo de até 500. LastRead requer tentativa positiva, SHA/caminho/hash/identidade válidos. Lead opcional deve apontar para a identidade da entrada e ter selo, identidade de origem e observação válidos; é `REFERENCE_NOT_READ`, nunca leitura atual. Textos válidos exigem UTF-8, não vazios e sem controles U+0000..001F, U+007F..009F; campos com motivo obrigatório também rejeitam apenas espaços. Status aceitos, além de vazio: COLLECTED, AUTH_REQUIRED, RATE_LIMIT, NOT_FOUND, INTEGRITY_ERROR, CONTENT_LIMIT, TREE_LIMIT, READ_FAILED, CANCELLED, TIMED_OUT, README_NOT_FOUND, OPAQUE_REVISION, REVISION_CONFLICT e CONTENT_BUDGET. Lead exige SourceCandidateID de até 150 bytes, ID reference_lead_ de até 100, linha inicial>=1 e intervalo de até duas linhas, trecho UTF-8 não branco de até 1024 bytes, até 24 termos e oito motivos de seleção (cada texto até 120 bytes). Seu selo é calculado com ID vazio e motivos de seleção como array vazio. Referência conserva caminho e hash da observação de origem.

Na seleção primária de discover, identidades desconhecidas/pendentes precedem examinadas/descartadas; estas podem preencher vagas restantes de metadados, sem provocar releitura.

Ao iniciar uma rodada admitida, clonar e avançar geração uma vez, inclusive sem leituras. Coleta válida marca `EXAMINED`, registra LastRead, status e tentativa. Falha tentada registra status/tentativa e mantém pendência; não descarta automaticamente. Sem espaço, conservar histórico e emitir `MEMORY_CAPACITY`. A inserção/atualização limitada reserva 4096 bytes sob 1 MiB; não remover entradas antigas para abrir espaço.

Feedback offline só altera alvo conhecido para `PENDING` ou `DISCARDED`, com estado sem distinção de caixa e motivo válido obrigatório. Preserva observações/tentativas, incrementa geração e ressela; não aceita `EXAMINED` por declaração. **Obrigação de fechamento:** saída bem-sucedida deve voltar a validar. Feedback, discover e expand devem usar a mesma política L1 para rejeitar avanço de geração999999 com ErrInvalidRunInput antes de efeitos ou mutação de entrada; a validação do histórico999999 continua permitida. Geração999998 admite o último avanço para999999. Não saturar nem reiniciar a geração. MEM-04 e o passo0147 governam essa correção, conservando o contraexemplo anterior na bancada0146.

## Foco e fronteira

Sem foco explícito, herdar o foco do relatório anterior; `AllTrails` retorna à fronteira geral. `AllTrails` com foco explícito é inválido. Foco aceita ID de entrada ou URL de identidade conhecida; rejeitar desconhecida, descartada ou alias cuja base esteja descartada. Alias conhecido é resolvido à base canônica conforme a seção seguinte.

Derivar links de origens congeladas ativas, conservando o primeiro testemunho por destino. Examinar no máximo 128 URLs por documento, 64 referências antes da seleção e 12 leads; testemunho de até duas linhas/1024 bytes. Não reintroduzir destinos examinados/descartados nem origens atuais. Links dos novos READMEs só podem ser seguidos na próxima chamada: uma chamada não faz recursão ilimitada.

Percorrer citações conhecidas com visitados, limitado à memória de 256 identidades. Com foco, incluir apenas sua região alcançável, distância mínima em arestas e o próprio foco com distância 0; descartes e conflitos de alias bloqueiam travessia. Sem foco, distância ausente. O grafo resume o primeiro testemunho; não recupera todos os caminhos alternativos. Uma entrada conhecida reaberta pode ser lida por expand sem inventar Lead; discover exige trilha para seus followups.

## Aliases e vizinhos

Reconhecer X/X.git apenas quando ambos já existem, com mesmo connector e URL completa/host correspondente. Não inferir base ausente, fusão global, renomeação ou equivalência semântica. Não reescrever identidade, estado, motivo ou contador da memória.

| Situação conhecida | Seleção |
| --- | --- |
| Alias pendente, base pendente | Base entra; alias `DUPLICATE_ALIAS` |
| Alias pendente, base examinada | Alias `ALIAS_ALREADY_EXAMINED` |
| Alias pendente, base descartada | Alias `ALIAS_DISCARDED` |
| Base pendente, alias examinado | Base `ALIAS_ALREADY_EXAMINED`; motivo explícito na base permite reabertura |
| Base pendente, alias descartado | Base `ALIAS_STATE_CONFLICT`, lacuna e bloqueio da travessia |

`neighbors` é snapshot anterior às leituras: fila ordenada seguida das supressões. `frontier_entry_ids` enumera a fila admitida, inclusive destinos que serão omitidos por fonte/leitor. `selection_policy` é `eligible-distance-attempts-age-reference-v1`. Resultado atual fica em followups; posição, distância e seleção não medem novidade ou relevância.

## Ordenação e execução

Ordenação determinística: fonte selecionada com leitor disponível; distância se houver foco; menos tentativas; FirstGeneration mais antiga; presença e ordem na lista atual de referências; ID como desempate. Discover usa o mesmo ordenador sem distâncias.

Planejar no máximo 16 wrappers de followup; excesso permanece pendente com `FOLLOWUP_PLAN_LIMITED`. Verificar, nesta ordem: cancelamento (`NOT_STARTED`), fonte (`SOURCE_NOT_SELECTED`), leitor (`READER_UNAVAILABLE`), tentativas (`FOLLOWUP_BUDGET`), documentos (`DOCUMENT_BUDGET`), bytes (`CONTENT_BUDGET`), e só então ler. GitHub usa FollowReader; demais forges usam RepositoryReader. Runtime/deadline são portas explícitas; rede/relógio concretos pertencem a L3.

Follow-links limita tentativas a 0..4. Leitura atual compartilha teto de oito documentos/128 KiB/32 KiB por arquivo com a coleta primária de discover; expand começa sem conteúdo atual. Falha do provedor consome tentativa; cancelamento e omissão antes da leitura não viram sucesso. Se bytes retornados ultrapassarem o agregado, descartar o bundle, zerar bytes no recibo, marcar `CONTENT_BUDGET` e manter pendência. Não vazar erro bruto externo.

| Opção | Normalização inicial de discover | Limite/herança em expand |
| --- | --- | --- |
| follow-links |2 se ausente |0..4; ausência herda, ponteiro 0 desabilita |
| max-passages |8 se omitido |1..16; zero na API herda |
| timeout total |90 s se omitido |1 s..5 min; zero na API herda |
| attempt-timeout |10 s se omitido |1 s..30 s e não maior que total; zero na API herda |

Opções herdadas de catálogo continuam válidas pelo [plano de sítios](../prompts/01_core/usecase/insight_site_plan.md), mas não autorizam consultas em expand. A publicação CLI rejeita zeros explícitos para passagens/tempos.

## Contexto e resultado

Cartões só podem vir de bundles válidos coletados na chamada, sob max-passages global e no máximo um cartão por followup. Conservar forge, namespace, commit, caminho e trecho na referência; não fabricar cartão para metadado, memória ou fonte omitida. IDs e hashes seguem o algoritmo existente. Tarefas apontam cartões atuais, entradas de memória e leads; referências de catálogo ficam vazias na expansão.

Expand publica Calls=0, attempts/candidates/readmes/bundles vazios e nenhum contexto antigo como coleta atual. Memória aparece no próprio relatório; não existe cache global novo. NoveltyState permanece `NOT_ASSESSED`, EffectState `UNVERIFIED`.

Estado final, antes de cancelamento:

1. Sem leituras, começar por `FRONTIER_EXHAUSTED`, ou `FRONTIER_PENDING` quando houver pendências.
2. Com coleta, `CONTEXT_PREPARED`; falha, wrapper omitido ou lacuna tornam `PARTIAL`.
3. Sem coleta, com tentativa e falha, `FAILED` e `ErrInsightDiscoveryUnavailable`.
4. Sem coleta/falha tentada, com pendência e follow-links positivo, `PARTIAL`.
5. Qualquer lacuna torna `PARTIAL`, exceto se já for `FAILED`.

Na conclusão do executor, erro de contexto prevalece: `CANCELLED` ou `TIMED_OUT`. Selar esse relatório resultante, inclusive quando retorna erro de leitura. Contexto nulo/já cancelado, erro de admissão, plano ou preparação não produzem relatório publicável; não prometer relatório selado para todo erro possível. Fronteira esgotada significa somente ausência de pendentes conhecidos naquele escopo, sem concluir inexistência de outras abordagens.

## Transporte CLI

Expand exige input e output; aceita memory, focus, all-trails, follow-links, max-passages, timeout e attempt-timeout. Não aceita pergunta/fontes novas, argumentos posicionais ou flags duplicadas. JSON é estrito: campos desconhecidos, chaves duplicadas e múltiplos valores são rejeitados pelo transporte. Proteger saída contra alias dos arquivos input e memory antes de compor clientes; validar semântica em L1 antes da factory. O core revalida ao executar.

Memory aceita snapshot autônomo ou relatório contendo memória, valida pelo core e publica somente o snapshot atualizado. Target/state/reason são obrigatórios semanticamente. Não busca rede nem precisa de factory de pesquisa. Extrair memória de relatório não autentica o relatório inteiro; expand verifica o selo e as fontes da entrada original.

Publicação usa o host e proteção de escrita existentes, sem ecoar segredos. Relatório selado de falha parcial deve ser publicado antes de devolver erro sanitizado. Erro de admissão sem ID não publica. Um único relatório contém estado e conteúdo da chamada; persistência, sinais e efeitos concretos permanecem em L3.

## Perfil textual compartilhado — 0162

Fontes congeladas ROOT_README_ONLY devem cumprir o perfil da coleta atual: texto UTF-8 não branco, sem controles exceto tab/CR/LF. A guarda de fontes usa o mesmo predicado puro da coleta antes de leitores, além de recibos/hash/LastRead e limites existentes. Vazio, espaços e NUL não se tornam origem válida por recalcular selos. Esta é coerência dos dados declarados, sem autenticação remota. O corpo histórico0145 de insightNeighborhoodOrigins permanece intacto. Origem: [passo0162](../passos/0162-aplicar-quotas-e-coerencia-de-fontes.md); consumidor: [regressões de perfil](../../01_core/usecase/insight_frozen_text_regression_test.go). SRC-01 ganha estes controles sem declarar cobertura integral das demais combinações.
