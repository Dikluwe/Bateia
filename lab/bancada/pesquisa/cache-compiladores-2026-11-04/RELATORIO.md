# Invalidação de caches e reúso de recursos em compiladores

Pesquisa Bateia, 2026-11-04. Resultado: **hipóteses de transferência, benefício não verificado**. Nenhuma otimização, alteração de L0 do typst-crystalline, branch, worktree, commit ou push. A prioridade recomendada é medir e tentar reduzir trabalho repetido dentro da prova vigente; depois separar produtos imutáveis e dependências; eventos e catálogo nativo ficam condicionados a provas adicionais.

## Pergunta, método e estado

Como diminuir aquisição/validação de fontes e trabalho repetido de compilação preservando reconhecimento na próxima recompilação, seleção, conteúdo, diagnósticos e renderização?

Foram lidos `.agentrules`, ADR-0008, mapa arquitetural e guia de engenharia do Bateia, contratos/consumidores de `insights discover` e `explore`. A capacidade utilizada pertence à investigação A1-INV: planejamento e seleção em L1, CLI em L2, leitores em L3 e composição em L4. Não houve evolução dessa capacidade. A bancada conserva artefatos; o relatório é análise de pesquisa, sem autorização de adoção.

O binário foi compilado da árvore atual com `GOCACHE=/tmp/bateia-research-gocache`, saída `/tmp/bateia-cache-research-20261104`. O primeiro build tentou o cache padrão e falhou por filesystem somente leitura; o seguinte passou. Uma tentativa inicial de discover incluiu `--no-token`, flag que pertence ao fluxo web e não a insights: foi rejeitada antes de produzir pesquisa (`discovery-queries.log`). A tentativa válida usou:

```sh
/tmp/bateia-cache-research-20261104 idea insights discover \
  --question 'incremental compiler cache invalidation tracked dependencies' \
  --sources github --max-calls 5 --max-references 4 --follow-links 2 \
  --timeout 70s --attempt-timeout 8s --output discovery.json
```

`discovery.json` e a repetição `discovery-network.json` preservam plano, consultas DIRECT/BROADENED/ALTERNATIVES/TRADEOFFS/COUNTEREXAMPLES, cinco tentativas por rodada e `FAILED/PROVIDER`, zero candidatos admitidos. A repetição com permissão ampliada também não completou consultas. Teste sem credenciais a api.github.com/raw.githubusercontent.com retornou falha de resolução de nomes. Não se atribui causa remota específica apenas ao erro protegido do Bateia e não se interpreta falha como busca vazia.

O conector GitHub conseguiu resolver commits e ler READMEs de Salsa e comemo. `github-meta-{0,1}.json` registra origem e revisão; `github-file-{0,1}.json` conserva conteúdo e blob SHA. Esses dados alimentaram bundles reais, sem assessments inventados. O primeiro bundle usou kind `README`, não reconhecido pelo seletor, e gerou `SEARCH_NEEDED`; foi preservado. O segundo usa `documentation`, conforme domínio, em `bundles-v2.json`. Executou-se:

```sh
/tmp/bateia-cache-research-20261104 idea insights explore \
  --question 'incremental compiler cache invalidation tracked dependencies' \
  --familiar 'full metadata validation' --bundles bundles-v2.json \
  --output exploration-v2.json
```

O resultado final é `CONTEXT_PREPARED`, com dois cartões (comemo README:49–63 e Salsa README:25–40), cinco perspectivas e a lacuna `INCOMPLETE_SOURCE_BUNDLES`; `NOT_ASSESSED/UNVERIFIED` não significa insight validado. Bundles continuam incompletos porque só READMEs foram coletados. Não se forçou `insights analyze` nas quatro classes do extrator, que não modelam esta pergunta. Complementou-se a leitura com documentação primária de mecanismos, APIs e código local das dependências fixadas. `web1..3.json` registra consultas; `read1..6.json` registra respostas/recortes e falhas. As páginas web móveis foram consultadas nesta data, não representam versões fixadas do alvo. São recibos de leitura parcial, não arquivos integrais das páginas.

### Interpretação das tarefas produzidas pelo Bateia

O cartão `card_f23eb1331f055ac23a055aa9ed8afb3bfde11d303f8c629fc92ad173e163f64a` localiza o mecanismo comemo; `card_9aec801e93e907607e5b1e61ee358f96bf5456f0b9ee327edac42439711e09d0` localiza inputs/funções memoizadas Salsa. As cinco perspectivas do resultado foram respondidas assim: ALTERNATIVES gera A/C/D; INVERSION pergunta se a unidade de invalidação deve ser uma projeção usada em vez do catálogo inteiro (B/C); CROSS_DOMAIN usa Skyframe como analogia explícita (C/E); RECOMBINATION combina prova externa atual com rastreamento de consultas (C/D), sem confundir observação com memoização; COUNTEREXAMPLE produz a matriz adversarial e rejeita H no contrato vigente. A interpretação e a seleção das alternativas são desta análise, não classificações semânticas emitidas pela CLI.

## Baseline e contrato conferidos no alvo

Os diagnósticos `bottleneck-first-round-2026-11-04.md` e `memoization-gap-2026-11-04.md` foram lidos. O estado informado é branch Tekt, HEAD `3c6b0ce7b7632cc74b004cea6681e4746e63d28d`, árvore alterada. O diagnóstico mede texto pequeno CLI em 95,15 ms e CPU de ciclo watch em 19,09 ms; suas amostras inclusivas apontam validação de fontes. Esses números são contexto recebido, não medições refeitas aqui. `ProbeKey::unchanged`, `scan_stamp` e `GenerationKey::unchanged` aparecem sobrepostos; não se somam as amostras nem se deduz economia máxima delas. Outline tem seis finalizações de candidatos, não seis iterações de um único loop. Source/CST já é retido. Contagem de atributos ou eviction não mede hit.

FontCatalogSession já mantém uma geração e o par slots/book por Arc; o World é novo e demanda o catálogo preguiçosamente. O hot path em `03_infra/src/world/font_session.rs:61` verifica FastKey e depois GenerationKey. Em `fontdb/generation.rs:68`, scan_stamp faz symlink_metadata e metadata; a prova de raízes e a lista de faces podem observar caminhos relacionados. Isso sugere sobreposição, não comprova redundância removível: link e alvo têm semânticas diferentes e uma face pode corresponder a índice diferente na mesma coleção.

Os L0 `prompts/infra/fontdb/generation.md` e `world/font_session.md` exigem:

- prova rápida apenas em Linux fontconfig padrão sem fontes de projeto; outros casos recusam e usam aquisição fresca;
- bytes/configuração efetiva, ambiente/cwd, diretórios inclusive ausentes, candidatos por extensão inclusive inválidos e symlinks/alvos; ordem e duplicatas dos loaders;
- identidade Unix dev/ino/size/mode/mtime/ctime com nanossegundos; rechecagens depois da preparação e antes de reúso/publicação; mudança/erro/instabilidade recusa retenção;
- uma geração futura, Arc antigos válidos para Worlds ativos, sem reciclar índices ou negativos antigos; 16.384 descritores e 1 MiB de caminhos para elegibilidade, sem truncar catálogo fresco;
- sem watcher novo, fc-list, cache em disco, fingerprint só por mtime ou snapshot só de diretório; a fronteira é demanda na recompilação, não snapshot atômico de todo o filesystem ou evento por mudança apenas de fonte.

Há ainda instrução explícita em font_session sobre preservar ordem, short-circuit e postcheck e **não memoizar o resultado de verify**. Logo, compartilhar observações entre fases ou remover verify exige proposta própria de L0, mesmo se a intenção for preservar garantia observável. Esta pesquisa não concede essa autorização.

No vanilla, conferiram-se por `git cat-file` os três blobs contra `64c5e8c02fe0394d913ead36e42ef8c1334d2107` (`vanilla-pin-check.json`): watch cria World fora do loop, reset limpa arquivos/relógio, não fontes; FontStore mantém OnceLock por slot. Isso é precedente de retenção e carga por demanda, não prova equivalente de atualização de catálogo. Não foi medido ganho contra vanilla.

## Evidências primárias e o que estabelecem

**S1 — rustc.** O [guia de incrementalidade](https://rustc-dev-guide.rust-lang.org/queries/incremental-compilation-in-detail.html) descreve red-green, fingerprints estáveis, recuperação seletiva e queries de projeção. Reexecutar uma dependência pode produzir o mesmo resultado e evitar propagar mudança. A [descrição do DAG](https://github.com/rust-lang/rustc-dev-guide/blob/main/src/queries/incremental-compilation.md) explica por que validar dependências na ordem observada importa quando o caminho de controle muda. Isso documenta validade de consultas sob entradas rastreadas, não descoberta automática de toda mudança externa.

**S2 — Salsa.** [Algorithm](https://salsa-rs.github.io/salsa/reference/algorithm.html), [Memo](https://salsa-rs.github.io/salsa/plumbing/terminology/memo.html) e [LRU](https://salsa-rs.github.io/salsa/plumbing/terminology/LRU.html) separam revisão verificada, revisão de mudança, dependências/durabilidade e retenção de valor. Durabilidade reduz verificações quando a revisão pertinente não mudou; não é previsão temporal de frescor. O [exemplo lazy](https://salsa-rs.github.io/salsa/common_patterns/on_demand_inputs.html) lê arquivos por demanda e usa watcher para atualizar inputs; não é demonstração de recuperação de toda perda/race. README fixado em `30c80938bcc32cdf19ed33dd5d58e6ddbccd34d2` foi processado pelo Bateia.

**S3 — comemo.** O [README fixado](https://github.com/typst/comemo/blob/5944487f5e9d11949038662df4b307ec56f153b3/README.md) mostra memoize + track + Tracked. A [documentação 0.5.1](https://docs.rs/comemo/0.5.1/comemo/attr.memoize.html) explicita Hash completo, restrições de impureza e replay de mutações tracked. O alvo fixa **0.4.0**: foram lidos localmente src/lib.rs, src/cache.rs e comemo-macros/src/lib.rs dessa versão. Não se transplantam APIs/restrições de 0.5.1 sem auditoria. O fonte 0.4 também exige entradas Hash completas e restringe impureza; `track` isolado não armazena resultado de eval/layout. Revalidar consultas externas pode continuar custando I/O.

**S4 — LLVM/Clang.** [PCH internals](https://clang.llvm.org/docs/PCHInternals.html) documenta desserialização por demanda e validação de metadados de linguagem/alvo/predefines. [DependencyScanningFilesystem.h](https://clang.llvm.org/doxygen/DependencyScanningFilesystem_8h_source.html), versão exibida 24.0.0git, documenta cache compartilhado por filename/UID, espera de produtor em andamento e cache local por worker. Seu detector de entradas obsoletas descreve negativos que passam a existir e tamanho alterado; não é prova suficiente para substituição de mesmo tamanho/mtime restaurado. Retenção e coalescência são precedentes; o checker não substitui a identidade Tekt.

**S5 — sistemas de build.** [Skyframe](https://bazel.build/reference/skyframe) modela valores imutáveis e dependências explícitas, inclusive resolução de symlink em passos; leitura externa sem dependência registrada causa incrementalidade incorreta. O [texto no repositório](https://bazel.googlesource.com/bazel/+/d167757e358b239f261afa41c0725ae74e86aee5/site/en/reference/skyframe.md) discute invalidação bottom-up/top-down, stat dos inputs e eventos como otimização. A entrada no grafo continua precisando de um observador correto.

**S6 — Fontconfig e fontdb.** [Referência Fontconfig 2.18.3](https://fontconfig.pages.freedesktop.org/fontconfig/fontconfig-devel/) define FcDirCacheRead/Rescan, FcConfigUptoDate por arquivos/timestamps e FcInitBringUptoDate sujeito ao intervalo de rescan. O [manual de configuração](https://fontconfig.pages.freedesktop.org/fontconfig/fontconfig-user.html) descreve cache por arquitetura/versão, configuração e SOURCE_DATE_EPOCH. No fontdb **0.23.0** do Cargo.lock, `src/lib.rs:503–575` usa fontconfig-parser, extrai diretórios/aliases e faz load_fonts_dir; `:309–384` canonicaliza links e filtra extensões. Não chama o cache binário nativo Fontconfig nesse caminho. A [API fixada](https://docs.rs/fontdb/0.23.0/fontdb/struct.Database.html) é contexto, mas a evidência decisiva é o fonte local. Não há equivalência demonstrada entre catálogos, ordem, metadados, aliases ou seleção.

**S7 — eventos Linux.** O [manual inotify](https://www.man7.org/linux/man-pages/man7/inotify.7.html) documenta overflow, observação por inode, recursão manual, races de diretório/rename, ausência de notificação por mount sobreposto, alterações remotas em rede e por mmap. Eventos podem ser coalescidos. Isso refuta tratar fila vazia como prova universal de ausência de alteração.

**S8/S9 — plataformas.** [Apple FSEvents](https://developer.apple.com/library/archive/documentation/Darwin/Conceptual/FSEvents_ProgGuide/UsingtheFSEventsFramework/UsingtheFSEventsFramework.html) exige rescan para MustScanSubDirs, perda e mudança da raiz; histórico é consultivo. [ReadDirectoryChangesW](https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-readdirectorychangesw) prevê overflow com descarte/zero bytes, enumeração para recuperar e filesystems que não suportam a operação. São APIs diferentes, sem prova portátil única.

## Alternativas de transferência, ordenadas pela preservação do contrato

As propostas abaixo são **inferências desta pesquisa**. Documentação do mecanismo não demonstra benefício no alvo. “Compatível em intenção” não significa permitido pelo L0 atual nem correto sem experimento. A matriz adversarial posterior integra os critérios de cada proposta.

### A — Reduzir observações duplicadas dentro de uma fase de prova

Precedente: cache local/compartilhado e produtor único de Clang (S4), revisão verificada de Salsa (S2). Evitaria stats repetidos para várias faces de uma coleção e consumidores da mesma observação. Não evita a cobertura de entradas exigida nem a segunda amostragem que detecta concorrência.

Hipótese: inventariar `(operação, caminho, follow/nofollow, fase, configuração)` e deduplicar somente leituras que podem compartilhar a mesma observação. Uma chave só por caminho é insuficiente: lstat do link, stat do alvo, caminho canônico e ausência têm contratos distintos. Não unir precheck e postcheck, não persistir erro/ausência em outro ciclo, não tratar Mutex da sessão como lock de escritores externos. A memória proposta é limitada à fase/chamada, com teto e fallback original. Mantém fallback não Linux e configuração customizada.

Compatibilidade: intenção preservadora, mas remover observações pode mudar a janela de corrida; compartilhar resultado de verify entre fases contraria texto atual. Primeiro experimento é apenas contagem em Lab, sem patch produtivo. Refutação mínima: substituição/retarget entre duas observações que o candidato uniria resultar em reúso quando o oracle fresco recusa ou muda seleção. Sem redução material de chamadas/CPU, rejeitar também a hipótese de benefício.

### B — Separar catálogo atual de recursos efetivamente usados

Precedentes: carga lazy PCH (S4), FontStore/OnceLock vanilla e sessão já existente. Evitaria reconstruir/ler/decodificar objetos de fontes demandadas que continuam idênticas quando algum outro item do catálogo muda; pode beneficiar aquisição CLI somente com outro desenho de persistência, que fica fora desta opção em RAM.

Manter a prova completa do catálogo: uma fonte nova pode vencer uma busca por família/estilo/fallback mesmo se a fonte anterior continua válida. Em uma geração nova, refazer índices/associações no par book/slots; só transportar payload imutável quando identidade de arquivo/índice, bytes e contexto relevante de parser garantirem equivalência. Não transportar um índice cru de geração antiga. Erro em OnceLock não pode se tornar negativo permanente para arquivo reparado. Retenção de métricas deve carregar identidade da fonte, não apenas família/nome.

Compatibilidade: preserva intenção de frescor se o catálogo atual continuar sendo a autoridade; atravessar gerações amplia a retenção além da sessão vigente e pede L0 próprio. Limitar por bytes/entradas e medir gerações antigas ainda presas por Arc, inclusive muitos Worlds concorrentes. Refutação: adicionar fonte homônima preferida, retirar/substituir a usada e reparar fonte antes inválida; qualquer seleção/diagnóstico diferente do fresco rejeita. Medir reconstruções e RSS com catálogo mudando e mesmas fontes usadas.

### C — Banco de consultas com projeções e validação preguiçosa

Precedentes S1/S2/S5. Evitaria propagar toda mudança do catálogo para todo eval/layout quando consultas observadas continuam iguais; não evita por si só descobrir/validar filesystem. Separar conceitualmente observação externa, catálogo ordenado, lookup de seleção e recurso identificado. Uma revisão só avança ou permanece válida depois da prova externa; contador de ciclos não certifica que o mundo não mudou.

Uma query de seleção depende também de ausência e da ordem das alternativas, aliases, idioma, cobertura e candidatos antes inválidos. Entradas devem representar essa competição, não apenas o arquivo vencedor. Validar dependências na ordem de controle observada; não executar ramo antigo depois de uma condição mudar. Regras de negócio/identidades em L1, I/O concreto em L3, lifecycle em L4, mediante contrato específico. Grafo e memos precisam quotas; durabilidade não pode esconder modificação externa não ingerida.

Compatibilidade: preservadora se a prova vigente alimenta o banco; instalação de Salsa é hipótese arquitetural maior, não primeiro patch. Refutação mínima: um lookup com candidato antes ausente/ruim que se torna preferido deve mudar o resultado; um arquivo não pertinente pode mudar sem reexecutar a projeção dependente. Comparar hits reais, custo do grafo/validação e memória com cálculo direto. Se custo de validação domina ou falta dependência, não adotar.

### D — Fronteira mínima de memoização de imports/eval; layout separado

Precedente comemo S3 e queries S1. Evitaria avaliação de imports/derivações que não mudaram; para outline, talvez trabalho repetido entre candidatos com dependências idênticas. Não há evidência de que seis candidatos sejam equivalentes, nem de que eval/layout explique o pequeno documento mais que fontes.

Condição: explicitar World e dependências transitivas, Source, inputs, features/target, rotas, contexto, fontes, introspecção e diagnósticos/efeitos. O hash público de Content não é prova de identidade suficiente para layout/origem. Reprodução de warnings/sink deve ser especificada; cache hit não pode suprimir erro nem vazar resultado entre Worlds/hosts. Introspecção/rollback e estilos/regiões podem diferir por candidato. Limites de retenção/eviction não constituem prova de hits. Verificar implementação 0.4.0 e tensão ADR-0030/L1/globals antes de qualquer escolha.

Compatibilidade: não autorizada pela existência de track/plugins; precisa fronteira e L0 próprios. Primeiro candidato experimental: import caro puro sob World congelado em Lab; layout tem experimento independente depois de decompor tentativas. Refutar com alteração transitiva, plugin host, diagnóstico repetido e nova fonte preferida; na sequência comparar saída, contagem de execuções efetivamente evitadas, CPU/latência e RSS. Mesmo hit correto perde se hashing/replay/validação custa mais.

### E — Snapshot/epoch de compilação com entradas imutáveis

Precedentes de valores imutáveis Skyframe, revisões Salsa e Arc já presente. Evitaria que consumidores repitam validação de uma geração já capturada e permitiria compartilhar recursos entre tarefas concorrentes. Uma epoch nomeia o conjunto capturado; não garante captura completa nem verdade do estado externo.

Dentro de World atual, a aquisição OnceLock já evita parte desse custo. O ganho adicional deve ser contado, não presumido. Um snapshot completo de bytes pode fixar o conteúdo observado; FD aberto sozinho não congela escrita in-place e mmap pode refletir alterações. Para transportar catálogo para a próxima recompilação, nova prova segue necessária. Não chamar amostragem sequencial de snapshot atômico do filesystem. Captura com concorrência exige checagens antes/depois, recusa e fallback sem loop ilimitado; reter bytes pode elevar RSS e antecipar carga hoje lazy.

Compatibilidade: snapshot lógico dentro da fronteira já existente é candidato preservador; fonte hermética/copiar todos bytes ou nova fronteira de atualização altera contrato/custo e precisa decisão própria. Refutação: modificar arquivo/config/link durante captura e demandar fonte posteriormente; nenhuma mistura book/payload entre gerações pode aparecer. Medir ganho marginal sobre OnceLock e memória com Worlds simultâneos; rejeitar se trabalho só foi antecipado.

### F — Eventos como sinal de sujeira e recuperação conservadora

Precedentes S5/S7/S8/S9. Eventos podem antecipar invalidação/coalescer mudanças; a versão conservadora mantém a prova atual na demanda e pode evitar tentativas inúteis de hit, mas não promete retirar a varredura de validação quente. A versão que elimina stats depende de cobertura e barreira de observação ainda não demonstradas.

Plano experimental, não implementação: instalar observadores antes da enumeração; observar pais de raízes ausentes e configuração/includes, links e destinos pertinentes; enumerar, consumir eventos e revalidar até uma fronteira estável ou recusar. Diretório novo precisa registrar watch e enumerar o que surgiu antes do registro. Eventos invalidam a geração, não editam índices de book ativo. Em overflow, perda, falha de watch/limites, mount, mudança de raiz ou backend sem garantias, marcar estado desconhecido e fazer aquisição/prova fresca. Não publicar enquanto uma recuperação não atingiu seu critério. Um contador de eventos só ajuda se todos os eventos pertinentes entram antes do teste; fila drenada não elimina entrega tardia nem mudança não reportada.

Compatibilidade: L0 hoje proíbe watcher novo. Usar como dica preserva garantia apenas mantendo a prova; substituir a prova está fora do contrato vigente e sem demonstração de equivalência. Linux tem mmap/rede/mount como contraexemplos; outras plataformas exigem provas próprias. Memória inclui watches, fila, mapa de caminhos e recuperação sob tempestade. Refutação mínima: perder artificialmente um evento e forçar overflow, rename de árvore e alteração no intervalo scan→watch; na próxima compilação o resultado deve igualar aquisição fresca. Repetir com escrita mmap/rede elegível. Qualquer stale rejeita remoção de stats, mesmo se eventos comuns passam.

### G — Cache nativo Fontconfig como acelerador de metadados

Precedente S6. Pode evitar parse de faces e obtenção de metadados em inicialização CLI ou nova geração. Não é substituto demonstrado de descoberta fontdb, nem Fonte de Verdade automática de frescor. O alvo usa parser Rust e travessia; integrar biblioteca nativa não é apenas ativar uma opção já presente.

Variante conservadora: usar dados nativos como pista e validar contra a descoberta/identidade atual, comparando candidatos e derivando os campos consumidos exatamente como hoje. Se ainda precisar ler todas as faces, a economia pode desaparecer. Variante que substitui catálogo depende de paridade de diretórios, ordem/duplicatas, aliases, índices de coleções, extensões, erros/negativos, cobertura, configurações customizadas e fontes de projeto. FcConfigUptoDate/rescan interval não prova o caso de arquivo substituído com timestamp restaurado; FcDirCacheValid documenta validade no seu próprio regime, sem promessa equivalente ao Tekt.

Compatibilidade: uso nativo/cache em disco/fc-list é excluído do L0 atual; exige contrato e dependência novos. Plataforma/versionamento/ABI, inicialização global e concorrência precisam auditoria; não chamar FcFini enquanto outros usam objetos. Memória de mapas e caches externos entra no custo. Refutação mínima: comparar catálogos ordenados/seleção contra fontdb com cache nativo quente/desatualizado, fonte reparada, symlink retarget, config incluída, fonte de projeto e replacement com mesmo tamanho/mtime. Uma divergência rejeita equivalência; um benchmark separado pode ainda justificar pista sob validação.

### H — Catálogo congelado por sessão ou ambiente hermético explícito

Precedente vanilla para retenção; opção de produto, não correção transparente. Evitaria quase toda validação entre ciclos e permitiria queries sobre entradas imutáveis. Adições/remoções/substituições/configurações e reparos depois da captura ficariam invisíveis até recarregar/reiniciar, salvo mecanismo novo. Cache por TTL tem a mesma janela de staleness sem criar garantia melhor.

Compatibilidade: **incompatível com reconhecimento exigido na próxima recompilação** no modo atual. Hermeticidade só vale se toda entrada relevante for realmente fechada e identificada, inclusive recursos externos, links e configuração; copiar uma lista de caminhos não basta. Conservar Worlds concorrentes exige isolamento e orçamento de memória. Refutação imediata: instalar fonte melhor e editar o documento; se o resultado conserva seleção antiga, o modo viola a garantia vigente. Só estudar mediante decisão explícita de modo/contrato; nenhuma autorização concedida aqui.

## Matriz adversarial obrigatória

Cada alternativa A–G deve rodar contra aquisição fresca e contra a semântica de seleção/renderização, não apenas igualdade de hashes ou endereços Arc. Em A, C e E, manter barreiras da prova externa; em B e D, invalidar produtos após ingestão correta; em F e G, a cobertura é a própria hipótese em disputa. H falha por desenho nas mutações posteriores e permanece rejeitada para o modo vigente.

| Caso | O que precisa continuar verdadeiro / critério de rejeição |
| --- | --- |
| Adição/remoção/nested rename | Nova opção pode alterar seleção; remover vencedor precisa fallback atual. Validar namespace, inclusive raiz ausente criada e diretório renomeado para dentro. Observar apenas arquivos usados falha. |
| Replacement/mtime restaurado | Mesmo tamanho e mtime não bastam. Preservar dev/ino/ctime/mode exigidos e associação índice/bytes; chave mais fraca ou payload antigo rejeita. Não se promete proteção contra adversário capaz de forjar toda identidade além do contrato atual. |
| Symlinks/hardlinks | Link lexical, alvo, aliases e ciclos importam. Retarget com nome sem extensão, alvo externo e substituição por arquivo regular precisam oracle. UID não elimina dependência do nome; hardlink pode modificar bytes por outro nome. |
| Configuração/ambiente/cwd | Incluir includes, aliases/diretórios efetivos, config ausente que aparece, HOME/XDG/FONTCONFIG_FILE e cwd relevantes. Não classificar config como “alta durabilidade” sem ingestão de mudanças. |
| Inválida/inacessível → válida | Vigiar candidatos e erros que não entraram no book; corrigir bytes/permissões sem alteração do diretório deve permitir descoberta. Negativo por sessão é rejeitado. |
| Concorrência | Writers externos não usam Mutex da sessão. Intercalar mudança antes/durante/depois da preparação/verificação; nunca publicar par instável ou recolocar geração velha após recusa. Worlds ativos conservam seu par sem reciclar slots. |
| Eventos perdidos/tardios | Injetar perda, overflow e atraso, falha ao registrar watch, restart e limite de recursos. Estado desconhecido usa prova fresca. Ausência de evento/TTL não é prova. Para A–E com prova mantida, eventos não participam da validade. |
| Plataformas/filesystems | Prova rápida atual é Linux elegível; recusa é comportamento correto nos demais. Testar config customizada/projeto e manter full fallback. Backend de eventos Linux não autoriza macOS/Windows/rede. |
| Memória | Startup/frio/quente, muitas fontes/coleções, várias revisões e Worlds ativos. Medir RSS pico/retido, bytes dos memos e watches; manter limites e sem truncar resultado fresco para caber no cache. |

## Ordem proposta de experimentos e decisões

1. **E0 — mapa de observações, somente Lab.** Sobre o mesmo estado do diagnóstico, contar syscalls por fase/caminho/operação/ciclo e multiplicidade por face; separar aquisição CLI, startup watch, prova quente e rebuild. Instrumentação separada dos tempos, custo do instrumento controlado. Isso refuta a premissa de duplicação antes de desenhar A. Não derivar exclusividade de amostras inclusivas.
2. **E1 — modelo adversarial de A.** Replay determinístico de precheck/postcheck com mudanças entre leituras, sem alterar produto. Provar qual união preserva evidência; se exigiria memoizar verify ou suprimir barreira, registrar mudança contratual necessária. Só depois medir variante experimental isolada, com oracle fresco. Qualquer divergência elimina variante.
3. **E2 — catálogo versus recursos B/C.** Registrar quais lookups/fontes/métricas cada fixture usa e o que muda quando uma fonte não usada é instalada. Protótipo em Lab conserva prova inteira e compara projeção de seleção/identidade. Medir bytes/objetos reutilizáveis e custos do banco. Rejeitar se a mudança afeta seleção ou não evita custo material.
4. **E3 — observabilidade eval/layout, depois D.** Completar duração exclusiva de parse e cada tentativa outline. Contar consultas, execuções, validações e hits reais. Testar import caro em fronteira mínima; outline separadamente, incluindo mudanças de contexto/estilos/regiões/introspecção. Não transportar resultado entre seis candidatos sem prova de dependências iguais.
5. **E4 — snapshots E.** Medir ganho marginal sobre retenção já existente e testar captura concorrente; comparar RSS e coerência de payload/book. Não transformar epoch em certificado de disco.
6. **E5 — F e G em trilhas independentes.** Primeiro eventos como dicas mantendo oracle; em outra bancada comparar Fontconfig nativo contra fontdb, com caches obsoletos e mesma matriz. Remoção de validação e troca do catálogo só avançam após contrato específico e demonstração. Não combiná-las no primeiro experimento: isso esconderia causa de divergência.

Para cada bancada futura: repetir as cinco fixtures recebidas e controles sem demanda de fonte, CLI processo novo e watch persistente. Usar amostras pareadas/intercaladas, mesmas fontes/config/afinidade, três rodadas como contrato atual pede, startup separado, CPU/latência até publicação, RSS pico/retido e saídas esperadas. Fixar HEAD+dirty hashes, binário e ambiente; não comparar medianas de estados diferentes nem subtrair medianas para inventar fases. Incluir sucesso, fallback, erros e tempestade de mudanças. Critério funcional é zero divergência nos controles, com limites explicitados; ganho só existe após medição repetível e custo de chave/validação incluído. Não há percentual de ganho previsto nesta pesquisa. Builds Rust futuros seguem o diagnóstico: alvo exclusivo /dev/shm, um job, teto de 12 GiB e zero swap, sem concorrer com benchmarks.

## Lacunas e caminhos rejeitados

- Descoberta online Bateia não retornou resultados; corpus de README foi curado via conector. A coleta alternativa não é apresentada como sucesso do provedor Bateia. Falha DNS/provedor e ruído lexical não medem abrangência da literatura.
- Processaram-se dois READMEs no explore; o restante é leitura primária complementar registrada. Não há execução dos compiladores/bibliotecas pesquisados, prova formal de invalidação ou teste novo de performance Tekt.
- Não se auditou toda implementação Fontconfig nem se executou fc-cache/fc-list. A garantia nativa para as mutações adversariais não está estabelecida. Não se certifica a equivalência de metadados fontdb/Fontconfig.
- Páginas móveis Salsa/Clang/Fontconfig não fixam as dependências do alvo. comemo 0.4.0 foi conferido localmente; APIs modernas não são autorização de atualização. Integração de banco de consultas e cancelamento/snapshots concorrentes exigem estudo próprio.
- Não se demonstrou que sobreposições de stat possam ser removidas sem afetar detecção de corrida. A cláusula de não memoizar verify permanece vinculante. Não se alterou L0 para contorná-la.
- Rejeitados no contrato vigente: TTL como prova; “fila vazia implica fresco”; vigiar só fontes usadas; directory-only/mtime-only; UID/filename sem contexto; cache negativo permanente; reutilizar índice de geração velha; transplante do checker de tamanho Clang; trocar descoberta por fc-list; copiar retenção vanilla como se fosse paridade de atualização; usar track/eviction como hit; escolher memoização ampla com base em contagem de atributos.
- Não se propõe revisão de políticas de reúso do Bateia nem aprovação de licença/adoção; as fontes fornecem precedentes investigativos.

## Artefatos e verificação

`state-before.json` registra HEAD, status e SHA-256 dos arquivos rastreados/não ignorados dos dois repositórios antes das escritas desta pesquisa. `local-sources.json` e `local-sources.tar.gz` congelam o recorte lido com hashes; `vanilla-pin-check.json` comprova os três blobs. `preservation.json` confere os arquivos preexistentes no encerramento e diferencia mudanças observadas de causalidade não provada. `SHA256SUMS.json` indexa os artefatos finais.

O gate canônico `lab/ferramentas/verify.py` foi executado com cache/saída em /tmp; logs e resultados são copiados em `verification/`. Ler os códigos exatos, incluindo falhas ambientais/preexistentes: sucesso de lint/build ou do explore não promove as hipóteses. Não se modificou código para fazer gates passarem. Gates selados de módulo próprio não foram promovidos pela suíte raiz.

A repetição de `go test ./...` com loopback passou (exit 0); a primeira revisão automática de permissão expirou antes de executar, e a única repetição foi aprovada. A tentativa inicial restrita fica registrada. Os snapshots externos inicialmente descompactados causaram falhas de lint por seus `@prompt` apontarem para outro repositório; foram arquivados, sem modificar bytes ou suprimir regras. `local-sources.json` fornece os membros e hashes para consulta. O lint foi repetido após o empacotamento.

Resultado final: lint, auditorias arquiteturais, suíte raiz com loopback, vet, build de `./04_wiring` e diff-check passaram. `verification/final-summary.json` consolida os resultados finais sem apagar tentativas anteriores. Integridade dos membros do arquivo de fontes e dos bundles foi revalidada. Nenhum efeito de performance foi medido.
