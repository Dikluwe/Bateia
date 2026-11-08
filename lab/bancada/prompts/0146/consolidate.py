"""Materialize the bounded documentation proposal; no Go or old evidence writes."""
from pathlib import Path
import os
ROOT = Path(__file__).resolve().parents[4]
CONTRACT = ROOT / '00_nucleo/contratos/exploracao-de-vizinhancas-v1.md'
MATRIX = ROOT / '00_nucleo/contratos/exploracao-de-vizinhancas-aceitacao-v1.md'
STEP = ROOT / '00_nucleo/passos/0146-consolidar-prompts-de-exploracao.md'
# owner, sections, criteria, interface, obligations, boundaries, test oracles
specs = [
('01_core/domain/insight_neighborhood', ['interfaces-e-formatos'], ['WIRE-01'],
 'Declarar InsightNeighborhoodRequest, InsightNeighborhood e InsightNeighborhoodNeighbor. O pedido interno transporta relatório, memória opcional, foco, AllTrails e overrides; o relatório projeta origem, foco e fronteira para CLI e consumidores de JSON.',
 'Conservar todos os nomes, tipos e tags existentes. Memory/Focus são opcionais; FollowLinks é ponteiro para distinguir ausência de zero. Focus, SelectionPolicy e Neighbors opcionais preservam a serialização anterior. Distance é ponteiro: valor0 deve ser serializado quando presente. FrontierEntryIDs registra IDs, sem incorporar documentos.',
 'Tipos não validam pedidos nem executam efeitos. AllTrails pertence ao pedido interno; não acrescentar campo público indevido ao relatório. Estados e seleção pertencem aos casos de uso, não aos structs.',
 'Comparar bytes de relatório legado sem neighborhood; serializar expansão com origem/foco/fronteira e decodificar com campos estritos. Comparar ausência de neighbors/selection_policy com presença; distância0 deve permanecer distinta de nil. Verificar roundtrip de related_entry_id e tentativas anteriores sem interpretar esses campos como evidência.'),
('01_core/domain/insight_exploration_memory', ['interfaces-e-formatos','memória'], ['WIRE-02'],
 'Declarar InsightReadObservation, InsightMemoryEntry, InsightExplorationMemory e InsightFollowup. Consumidores são validadores, execução de aprofundamento, relatório multifonte e transporte de feedback.',
 'Snapshot expõe schema_version, id, generation e entries. Entrada expõe id, identity, state, first_generation, last_question e read_attempts; reason, last_status, lead e last_read são opcionais. Observação contém source/path/content_hash. Followup contém entry_id/lead_id/identity/status/cards, com receipt e bundle opcionais. Preservar tags e omissões existentes.',
 'Histórico declarado não é conteúdo atual: cartões/bundle ficam no followup, observação anterior fica na memória. Tipos não autenticam fontes, alteram estado ou fazem I/O. Manter a extensão opcional do relatório sem alterar schemas históricos.',
 'Construir relatório com entrada descartada/motivo e followup coletado/cartão; serializar e comparar roundtrip integral. Controlar campos opcionais ausentes e observação/lead presentes. Não usar uma entrada de memória para satisfazer a expectativa de bundle atual; validação semântica tem consumidor próprio.'),
('01_core/usecase/insight_neighborhood', ['propósito-e-alcance','admissão','contexto-e-resultado'], ['ADM-01','FOC-01','RUN-03','OUT-02'],
 'ExpandInsightNeighborhood recebe contexto, pedido de domínio e bindings; devolve relatório e erro. É a fachada L1 de expansão. ValidateInsightNeighborhood e as origens verificadas pertencem à admissão, não a uma interpretação paralela dentro da fachada.',
 'Coordenar uma única admissão, plano herdado, deadline por runtime, bootstrap de memória, contexto vazio atual, leads, fronteira, executor compartilhado e conclusão. Reutilizar as origens da admissão; propagar erro em vez de ignorar uma segunda validação. Selecionar a fronteira antes de iniciar leituras. Retornar resultado selado para falhas de leitura e interrupções concluídas pelo executor e preservar o pedido do chamador.',
 'Nenhuma consulta de catálogo, pergunta nova ou I/O concreto. Delegar integridade a sources, alcance a frontier, aliases a selection, quotas a followups e estados/selo a result. Não alterar a jornada discover ao reutilizar o executor.',
 'Usar conectores que falham se houver busca e leitores controlados. Encadear relatórios por três chamadas, verificando alvo efetivo, foco, memória e cartões de cada rodada. Adulterar entrada e exigir zero leituras; cancelar com quota de cartões e verificar resultado parcial rastreável. Reabrir alvo conhecido sem Lead e exigir leitura sem testemunho fabricado.'),
('01_core/usecase/insight_neighborhood_plan', ['admissão','foco-e-fronteira','ordenação-e-execução'], ['ADM-01','ADM-02','ADM-03','FOC-01'],
 'ValidateInsightNeighborhood preserva sua API pública e projeta a Request da admissão interna. admitInsightNeighborhood devolve pedido normalizado e Origins verificadas, consumidos pela fachada.',
 'Copiar/validar relatório e memória original; obter origens congeladas antes de aplicar feedback. Feedback conserva identidades e não retrocede geração. Resolver foco explícito/herdado e AllTrails, inclusive aliases descartados. Rejeitar geração efetiva>=999999. Herdar opções e aplicar somente overrides permitidos, revalidando o plano; FollowLinks0 difere de ausência.',
 'Admissão não chama leitores nem relógio/rede concretos. Os limites e classes de erro são os do contrato compartilhado. Cópia profunda inclui dados de fontes e memória; fontes admitidas não podem apontar a armazenamento mutável do chamador. Integridade documental é delegada ao proprietário sources.',
 'Admitir entrada válida e guardar a representação normalizada. Mudar conteúdo e estado no original, depois exigir pedido/origens admitidos intactos e leitura do alvo original. Incluir controles de selo, foco desconhecido, conflito AllTrails/foco, limites e feedback que remove identidade/retrocede geração; todos antes das portas. Omitir override herda; ponteiro0 impede leitura.'),
('01_core/usecase/insight_neighborhood_sources', ['fontes-congeladas'], ['ADM-01','SRC-01'],
 'Receber relatório copiado com memória original validada e devolver origens documentais verificadas para a admissão. Consumir bundles GitHub selecionados, contextos por sítio e followups da entrada segundo o mesmo contrato de origem.',
 'Confrontar identidade/namespace, SHA, caminho raiz, item/ref, bytes/hash, recibo COLLECTED/ROOT_README_ONLY e LastRead. Rejeitar duplicidades e pares coleta/bundle incoerentes. Aplicar oito documentos,128 KiB agregados e32 KiB por documento. Não confiar somente no selo externo do relatório: conteúdo e fingerprint precisam ser consistentes entre si.',
 'Sem consulta remota, travessia de grafo ou decisão de seleção. Origens só testemunham links já congelados; não viram cartões atuais. Erro de integridade/contrato deve chegar à admissão antes de qualquer leitor. A memória alternativa de feedback não substitui a verificação contra o histórico original.',
 'Criar origem válida e variar um aspecto por controle: selo, conteúdo, fingerprint, recibo ausente/incorreto, LastRead divergente, duplicidade, namespace, SHA e tamanho. Quando necessário resselo externo deve ser recalculado para isolar a inconsistência interna. Exigir rejeição e zero chamadas. Exercitar as três formas de origem e limites inclusivos/excedidos; a matriz registra quais controles ainda faltam.'),
('01_core/usecase/insight_neighborhood_frontier', ['foco-e-fronteira'], ['SRC-01','FOC-01','FOC-02','RUN-03'],
 'Receber origens já verificadas, memória, foco normalizado e plano; preparar leads e calcular entradas pendentes alcançáveis com distâncias opcionais para seleção/executor.',
 'Reutilizar derivação de referências documentais e registrar primeiro testemunho sem duplicar leitura. Percorrer citações conhecidas com visitados; foco tem distância0 e só alcança descendentes válidos. Bloquear descartes e conflitos de alias pela projeção proprietária; conservar namespace GitLab. Sem foco, usar fronteira geral sem inventar distância. Permitir pendência reaberta sem Lead em expand.',
 'Integridade de READMEs pertence a sources, aliases a selection e ordenação a followup_selection. Esta fronteira não revalida documentos nem chama provedor. Respeitar os limites de extração/memória do contrato; links obtidos nesta leitura só serão aprofundados em outra chamada. Primeiro testemunho não representa todos os caminhos alternativos.',
 'Construir ciclo e ramo não alcançável a partir de foco conhecido; exigir término, distância mínima e exclusão do ramo externo. Descartar intermediário e verificar bloqueio. Alternar modo geral e foco; conservar identidade GitLab com subgrupo e fonte restrita. Controle de alias não pode liberar caminho conflitante. Ler alvo sem Lead não autoriza fabricar citação.'),
('01_core/usecase/insight_neighborhood_selection', ['aliases-e-vizinhos','ordenação-e-execução'], ['ALS-01','ALS-02','ORD-01'],
 'Resolver relações conservadoras de alias conhecido e projetar InsightNeighborhoodNeighbor a partir da fronteira ordenada. Consumidores são admissão, travessia e relatório de expand.',
 'Reconhecer X/X.git somente com as duas identidades existentes no mesmo connector/host/URL correspondente. Aplicar a tabela compartilhada de supressões, conflitos e reabertura explícita, sem fundir entradas. Snapshot dos vizinhos antecede leituras: registrar tentativas anteriores, distância opcional, motivo e related_entry_id. Fila e projeção usam o mesmo ordenador; supressões vêm depois da fila.',
 'Não remover .git globalmente, mudar memória/contadores nem executar leituras. Conflito bloqueia automação e gera lacuna; exame prévio não comprova mecanismo. A projeção não altera a política sem distância de discover e não mede novidade/relevância.',
 'Combinar estados de base/alias: pendente, examinado, descartado e reaberto com motivo. Exigir código de supressão correto, ausência de leitura/consumo extra e memória anterior intacta. Variar host/forge e base ausente para impedir fusão indevida. Comparar neighbors anterior à execução com followups posterior; testar distância0 e ausência no modo geral.'),
('01_core/usecase/insight_neighborhood_result', ['contexto-e-resultado'], ['OUT-01','OUT-02','RUN-03'],
 'Inicializar InsightSiteReport de expand, preparar contexto vazio para a rodada, projetar tarefas e concluir status/erro/ID. Consumido pela fachada e pela publicação CLI.',
 'Manter coleções de catálogo/coleta antiga vazias, Calls0, origem em neighborhood e memória separada. Tarefas referenciam cartões atuais, memória e leads; catálogo vazio. Aplicar a sequência de estados do contrato, com FAILED/erro quando todas as tentativas falham e CANCELLED/TIMED_OUT prevalecendo. Selar resultado mesmo com erro operacional; nunca promover estado operacional a força de evidência.',
 'Sem provedor, leitura de arquivo, travessia ou mutação de aliases. NoveltyState NOT_ASSESSED e EffectState UNVERIFIED. FRONTIER_EXHAUSTED descreve o escopo conhecido, não inexistência de soluções. Cartões são fornecidos pelo proprietário de contexto, sem reciclagem de histórico.',
 'Falhar a única leitura elegível e exigir FAILED, erro de indisponibilidade, relatório selado e memória pendente com tentativa registrada. Retomar a partir dele e exigir coleta atual. Rodada com follow-links0 mantém fronteira e não chama leitor. Cobrir combinações de coleta/omissão/lacuna e cancelamento/deadline com oráculos da tabela de estados, sem inferir qualidade semântica.'),
('01_core/usecase/insight_exploration_memory', ['memória'], ['MEM-01','MEM-02','MEM-03','MEM-04'],
 'ValidateInsightExplorationMemory devolve cópia canônica ou erro; UpdateInsightExplorationMemory aplica feedback offline. O container possui bootstrap, consultas de estado, registro de leituras/testemunhos e seleção primária sensível à memória.',
 'Validar schema/selo/identidades/estados/observações/leads e limites compartilhados; ordenar entradas e selar com algoritmo histórico. Descarte/reabertura exigem motivo e preservam observações. Coleta marca EXAMINED; falha mantém PENDING com tentativa/status. Preservar histórico ao atingir capacidade e emitir MEMORY_CAPACITY. Na seleção primária de discover, desconhecidas/pendentes precedem examinadas/descartadas; estas só preenchem vagas de metadados e não geram releitura.',
 'Sem I/O, cache global, autenticação remota ou prova de mecanismo. Snapshot devolvido não compartilha estado mutável com entrada. MEM-04 é obrigação de fechamento ainda não satisfeita no feedback em geração999999; não tratar o defeito como regra aprovada. A consolidação não modifica código para ocultar a lacuna.',
 'Validar e modificar a cópia sem alterar original. Adulterar hash/estado/observação/lead/duplicidade, inclusive com resselo externo; exigir rejeição. Descartar e reabrir com motivo, rejeitar EXAMINED declarado e razão vazia. Esgotar256 entradas sem remover história; priorizar candidato novo. Na fronteira999999 exigir rejeição do avanço, registrando MEM-04 como lacuna atual em vez de afirmar cobertura verde inexistente.'),
('01_core/usecase/insight_followup_selection', ['ordenação-e-execução'], ['ORD-01','ORD-02'],
 'Ordenar fila de entradas pendentes recebida, usando plano/bindings, ordem das referências e distâncias opcionais. Mesma política atende executor discover e expand.',
 'Comparar nesta ordem: elegibilidade por fonte selecionada e leitor; distância quando fornecida; tentativas crescentes; FirstGeneration crescente; presença/posição na ordem atual de referências; ID. Ordenar somente a fila, mantendo IDs e memória. Sem distâncias, conservar exatamente a política histórica de discover.',
 'Nenhum efeito, contador incrementado, resolução de alias ou julgamento de relevância. A disponibilidade é inferida da porta configurada; não consultar o provedor para saber se está acessível. O executor possui quotas e a projeção possui motivos/snapshot.',
 'Colocar foco sem leitor e destinos elegíveis perto/longe com contadores distintos: elegibilidade vence distância, distância vence tentativas. Repetir com distância ausente e exigir prioridade por tentativas. Comparar memória/contadores antes/depois e zero leituras. Para empates, variar uma chave por vez e permutar a entrada; ORD-02 registra a tabela de desempates ainda incompleta na suíte.'),
('01_core/usecase/insight_followups', ['ordenação-e-execução','memória'], ['RUN-01','RUN-02','RUN-03','FOC-01'],
 'Executar fila limitada por portas, compartilhada entre discover e expand, produzindo wrappers de followup, recibos, atualizações de memória e pedidos de cartões. Preservar a entrada compatível de discover.',
 'Aplicar a ordem compartilhada; limitar wrappers a16 e tentativas por FollowLinks. Checar cancelamento, fonte, leitor, tentativas, documentos e bytes antes de ler, nessa ordem. Respeitar quotas globais da chamada e deadline por tentativa via runtime. Registar falha sem marcar exame; recusar bundle que exceda agregado. GitHub usa FollowReader e demais forges RepositoryReader. Não seguir descendentes recém-lidos na mesma execução.',
 'Efeitos somente pelas portas injetadas; nenhum cliente/rede/arquivo concreto. Memória pertence ao container, ordem ao seletor e cartões ao contexto. Expand pode ler pendência sem Lead, com lead_id vazio; discover continua exigindo trilha. Erros externos não podem vazar segredos no relatório.',
 'Executar duas rodadas com leitores em memória e verificar avanço sem repetir examinados. Descartar/reabrir; falhar e retentar; omitir fonte/leitor e contar efeitos. Combinar coleta inicial com followups para esgotar orçamento agregado, inclusive retorno excedente após fetch. Cancelar durante leitura e preservar apenas conteúdo coletado. Testar omissão antes do efeito e diferença entre expand sem Lead e discover com trilha.'),
('01_core/usecase/insight_followup_context', ['contexto-e-resultado'], ['RUN-02','RUN-03','OUT-02'],
 'Receber bundle efetivamente coletado e identidade de followup; usar seletor existente para anexar cartões ao followup, dentro do orçamento do contexto da rodada. Consumidores são executor e projeção de tarefas.',
 'Produzir somente cartões de bytes atuais, no máximo um por followup e dentro do max-passages compartilhado. Referenciar forge/namespace/commit/caminho e conservar os IDs/hashes existentes. Manter correspondência entre cartões no followup e referências das tarefas, inclusive quando outra leitura falha ou é cancelada.',
 'Sem chamadas de provedor, selagem de memória ou reutilização de cartões do relatório anterior. Fonte de GitLab conserva subgrupos; não reinterpretar sítios como GitHub. Omissão/bundle recusado não cria cartão. Conteúdo de fonte é dado, não instrução de execução.',
 'Limitar max-passages com várias coletas e comparar número total e refs por destino. Cancelar após uma leitura válida e exigir cartão preservado apenas para a coleta atual. Usar identidade de outro sítio e verificar namespace/commit/path. Controle sem bundle e controle com cartão histórico devem produzir zero cartão novo; a matriz distingue cobertura existente e controles a ampliar.'),
('02_shell/cli/insight_neighborhood', ['transporte-cli'], ['CLI-01','CLI-02','FOC-01','ADM-03'],
 'RunInsightNeighborhoodCommand transporta idea insights expand entre flags/JSON e fachada L1, recebendo factory/host injetados. Input/output são obrigatórios; memory/focus/all-trails e overrides são opcionais.',
 'Rejeitar flags duplicadas, argumentos posicionais, pergunta/fontes novas, valores vazios e limites explicitamente inválidos. FollowLinks0 é válido; max-passages/tempos0 explícitos são inválidos. Ler JSON estrito, proteger output contra input/memory e chamar ValidateInsightNeighborhood antes da factory. Encaminhar AllTrails/foco sem calcular herança localmente. Publicar relatório com ID em falha operacional antes de devolver erro sanitizado.',
 'Política de foco, aliases, memória e quotas pertence a L1. L2 não importa adaptadores L3, não chama rede concreta nem mantém segundo arquivo global de estado. Usar host/publicação protegida existente; barrar vazamento de segredos, sem repetir seu valor na mensagem.',
 'Usar factory contada e leitores controlados. Variar flags/JSON/selo/saída-alias e exigir zero factory/publicação em rejeição. Encadear foco herdado, all-trails e troca explícita, com conflito rejeitado antes da factory. Falha operacional com ID publica relatório útil; falha sem ID não publica. Controlar segredo em conteúdo/erro e exigir mensagem sanitizada.'),
('02_shell/cli/insight_memory', ['transporte-cli','memória'], ['CLI-03','MEM-02','MEM-04'],
 'readInsightMemory admite snapshot autônomo ou relatório multifonte com exploração; runInsightMemory transporta input/output/target/state/reason e publica snapshot atualizado pelo core.',
 'Usar leitura estrita e schema explícito; relatório pode conter neighborhood opcional. Validar memória com L1, sem afirmar autenticação do relatório inteiro. Proteger output contra entrada antes da escrita; delegar transição/selagem a UpdateInsightExplorationMemory. Publicação única contém snapshot novo, preservando arquivo original e recusando segredos pelo host existente.',
 'Comando offline, sem factory de pesquisa, política de estado na casca ou nova persistência global. Formato inválido, memória adulterada e feedback sem alvo/estado/motivo válidos são rejeitados com mensagem sanitizada. MEM-04 permanece uma divergência do core; a CLI não deve duplicar a regra para encobri-la.',
 'Ler as duas formas de entrada e realizar descarte/reabertura com motivo; reutilizar snapshot na próxima investigação. Variar schema, chaves desconhecidas/duplicadas, ID, estado e saída-alias; exigir ausência de publicação indevida e de factory. Verificar proteção de segredo sem ecoar valor. Distinguir leitura de memória de validação integral do relatório; limite de geração consome MEM-04, ainda pendente.'),
]
def link(label, target, base, anchor=''):
    return f'[{label}]({os.path.relpath(target, base)}{anchor})'
for owner, sections, cases, interface, behavior, bounds, tests in specs:
    layer = 'L2' if owner.startswith('02_shell') else 'L1'
    for test in (False, True):
        suffix = '_test' if test else ''
        path = ROOT / '00_nucleo/prompts' / (owner + suffix + '.md')
        snapshot = ROOT / 'lab/bancada/prompts/0146/before' / (str(path.relative_to(ROOT)).replace('/', '__') + '.txt')
        assert snapshot.exists(), path
        base = path.parent
        deps = ', '.join(link(s.replace('-', ' '), CONTRACT, base, '#' + s) for s in sections)
        criterion_links = ', '.join(link(c, MATRIX, base) for c in cases)
        # Row IDs are text in tables; deliberately link to the matrix, not fictitious anchors.
        material = ROOT / (owner + suffix + '.go')
        production = ROOT / '00_nucleo/prompts' / (owner + '.md')
        if test:
            body = f'''# {path.stem}

**Camada:** {layer}. **Responsável:** aceitação de `{owner}`. **Materialização consumidora:** {link(material.name, material, base)}.

## Interface e contrato

Consumir {link('o prompt proprietário', production, base)} e as seções {deps}. Oráculos desta fatia: {criterion_links}. Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

{tests}

## Limites e dependências

{'Testes do core usam apenas dados e portas controladas, sem I/O concreto.' if layer == 'L1' else 'Usar host de teste e portas/factory controladas; não depender de rede pública.'} Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A {link('matriz', MATRIX, base)} distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob {link('passo0146', STEP, base)}. O {link('texto anterior', snapshot, base)} preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
'''
        else:
            body = f'''# {path.stem}

**Camada:** {layer}. **Responsável:** `{owner}`. **Materialização:** {link(material.name, material, base)}.

## Interface e responsabilidade

{interface}

## Comportamento vigente

{behavior}

## Limites e dependências

{bounds}

As regras compartilhadas estão nas seções {deps}. Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos {criterion_links} e {link('o prompt de testes', path.with_name(path.stem + '_test.md'), base)}. A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob {link('passo0146', STEP, base)}; não é autoria ex-ante desse código. O {link('texto anterior', snapshot, base)} conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
'''
        path.write_text(body)
print(f'{len(specs) * 2} prompts consolidados')
