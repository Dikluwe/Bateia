# Prompt proprietário — `03_infra/local/client.go`

Origem L0 exclusiva de `03_infra/local/client.go`. Este contrato foi separado de `00_nucleo/prompts/03_infra/cache/cache_store.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Materializar somente a responsabilidade desta unidade, preservando os contratos observáveis: `NewClient`, `candidateRoot`, `safePath`, `FetchManifest`, `FetchTree`, `ExtractRepoFiles`, `FetchRepoMetadata`, `Client`. Manter a direção de dependências da camada e não incorporar responsabilidades pertencentes aos demais consumidores do prompt de origem.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/03_infra/cache/cache_store.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Revisão 0043 A — confinamento e identidade da leitura

Obrigação congelada: etapa A de `00_nucleo/passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md`. Manifestos, README, árvore e extração abrem uma raiz com `os.OpenRoot` e acessam nomes relativos somente pelas APIs `os.Root`. A contenção precisa vigorar na própria abertura, inclusive durante troca de links/diretórios; validar um caminho e depois abrir por nome absoluto não satisfaz o contrato. Links relativos internos válidos permanecem legíveis; links finais ou intermediários que escapam são rejeitados pelas APIs da raiz. Caminhos absolutos e traversal lexical são recusados.

Cada arquivo é aberto uma vez, inspecionado pelo descritor e lido com limite de 2 MiB + um byte detector. Rejeitar arquivo não regular, maior que o limite ou leitura falha; registrar tamanho dos bytes efetivamente retornados, sem confiar no tamanho anterior à leitura. O limite também vale para manifesto e README. Falhas de fontes opcionais não fabricam conteúdo.

A extração já retorna SHA-256 por arquivo calculado dos bytes retornados, intervalos lidos e identidade de snapshot determinística por caminhos normalizados e hashes ordenados. Identidade de conteúdo independe de HEAD e muda quando mudam bytes extraídos. BaseCommit é metadado Git, não prova da identidade dos bytes nem snapshot atômico do worktree. Não declarar análise semântica nem execução ao registrar leitura. Preservar o formato de snapshot do anotador existente para continuidade das etapas posteriores.

Limites declarados: a raiz escolhida é o diretório aberto; `os.Root` não isola mounts/hardlinks nem cria um snapshot atômico contra escrita concorrente. A implementação depende das garantias de `os.Root` da plataforma Go 1.25. Testes adversariais são de autoria segregada; não devem orientar a implementação.

Abrir descritores com flag não bloqueante para que FIFOs não bloqueiem antes da validação de arquivo regular (plataformas com `syscall.O_NONBLOCK`). A validação continua no descritor aberto, sem depender de `stat` prévio. O subprocesso `git rev-parse HEAD` permanece exclusivamente para metadados de commit-base, incluindo gitdir externo de worktrees legítimas; não executa filtros/hooks nem fornece os bytes analisados.

## Passo 0072 — enumeração completa da superfície local

`FetchTree` não possui corte silencioso por quantidade de entradas. Ele enumera toda a
superfície acessível sob a raiz confinada, em ordem independente do filesystem, ou
retorna erro quando a enumeração falha/cancela. Nunca devolver uma árvore truncada como
sucesso; o chamador transforma erro em `TreeComplete=false`.

### Refinamento A — aliases absolutos internos

Contraexemplo observado pelo verificador: arquivo regular interno + alias interno retornavam somente um arquivo. `os.Root` recusa links absolutos mesmo quando o alvo está sob a raiz. Estender a abertura com fallback limitado a 40 resoluções: inspecionar componentes apenas por `Root.Lstat`/`Root.Readlink`, converter alvo absoluto lexicalmente contido em nome relativo e recomeçar a resolução. Alvos externos são recusados antes de prosseguir. A abertura de conteúdo final continua exclusivamente por `Root.OpenFile`, preservando a contenção contra troca concorrente. Nunca usar `EvalSymlinks` seguido de abertura absoluta. Este fallback aceita aliases absolutos sob o nome absoluto da raiz aberta; aliases para outra grafia externa da mesma raiz podem ser recusados. Preservar o caminho solicitado na evidência e o hash dos bytes retornados.

## Revisão de segurança 2026-09-26

Aplicar `security.go`: árvore em lotes (máximo 20.000 entradas), sem alocar diretório inteiro antes de controlar orçamento. Exceder teto produz erro explícito, não sucesso truncado. Inventário expõe symlink/hardlink/nonregular; leitura rooted mantém alias interno mas recusa hardlinks por nlink do descritor. Extração até 100 arquivos/8 MiB (2 MiB por arquivo). Não executar Git ou herdar ambiente/configuração. BaseCommit consultivo vem de HEAD/refs/packed-refs confinados; gitdir externo fica indisponível explicitamente. Metadata não autentica bytes da worktree e certificados históricos não são atualizados.
