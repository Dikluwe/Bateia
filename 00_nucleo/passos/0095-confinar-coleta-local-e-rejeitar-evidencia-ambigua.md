# Passo 0095 — Confinar coleta local e rejeitar evidência ambígua

Continuação do endurecimento 0094, mantendo pesquisa simples e análise avançada
opcional. Os contraexemplos usam apenas fixtures sintéticas e arquivos temporários;
não há ataque externo, acesso a segredos reais ou execução de candidatos.

## Provedor local

- Remover processos Git. HEAD, refs soltas e packed-refs são lidos dentro de
  `os.Root` com limites; BaseCommit continua informação consultiva, não prova
  que os bytes da worktree pertencem ao commit. Ponteiro gitdir fora da raiz
  não é seguido. Metadata indisponível é explícita.
- Recusar hardlinks pelo descritor, antes e depois da leitura. Symlinks internos
  continuam permitidos quando a resolução permanece confinada. FIFO e outros
  arquivos não regulares não são lidos como conteúdo.
- Inventário distingue file/dir/symlink/hardlink/nonregular. Árvore é percorrida
  em lotes de 128 entradas, teto de vinte mil, caminho de 4096 bytes. Excesso
  retorna erro, nunca uma árvore truncada apresentada como completa.
- Extração: até cem arquivos, 2 MiB por arquivo e 8 MiB total. Busca de código
  também confina descritores, limita varredura e recusa hardlinks.

O contrato 0072 foi supersedido de forma explícita no teste/prompt: permanece a
exigência de não truncar silenciosamente, agora com erro no teto de segurança.
Certificados históricos não foram alterados. Worktrees com gitdir externo e
árvores acima do teto perdem disponibilidade de metadata/coleta nesta etapa.

## Entradas JSON e CLI

`parser.ValidateBoundedJSON` valida um único valor, UTF-8, até 16 MiB,
profundidade 64 e cem mil tokens. Rejeita membros duplicados no mesmo objeto,
inclusive escapes e aliases de casefold Unicode que o decoder Go pode aceitar
como a mesma chave de struct. Objetos independentes podem repetir nomes.
Erros não expõem valores ou nomes de chaves.

Os leitores restantes de `idea_recommend` e as famílias JSON de ideia usam
descritores regulares, não bloqueantes, limitados antes da decodificação.
Rules, modelos e replay também usam a guarda estrutural. Radar mantém suas
proteções contra alias de saída e valida bytes do mesmo descritor aberto.
Hashes continuam sendo calculados sobre os mesmos bytes, sem reinterpretar JSON.

Identificação consultiva de remotos Git também retém um único `os.Root` na
resolução/leitura. Ponteiros gitdir/commondir externos, aliases e hardlinks
não são usados para remapear identidade. Config ausente/insegura dá identidade
desconhecida, sem fallback para outra raiz. Layouts relativos internos continuam
funcionando. Metadados são conferidos antes/depois da leitura; não autenticam
origem remota. Essa proteção não é a quota agregada do hash de árvore.

É uma política mais restritiva que JSON comum: mapas com nomes que diferem apenas
em capitalização são recusados. Caminhos CLI continuam escolhas explícitas do
usuário, não uma sandbox. O cálculo de hash de árvore em recommend limita cada
arquivo, mas ainda não tem orçamento agregado de entradas/bytes; não declarar
cobertura global de recursos para essa rota.

## Coerência dos achados técnicos

Contraexemplos reproduzidos e regressões incorporadas:

- `COPYING` restritivo com Kind=file não pode ser omitido da unicidade da
  licença: a autoridade do inventário depende do caminho, não do rótulo importado.
- Bundles da mesma identidade/revisão não podem dividir LICENSE e COPYING para
  fabricar inventários unívocos. Repetição é erro; schema desconhecido não gera
  fonte verificada.
- Suporte MIT e refutação mit devem conflitar sob a mesma equivalência usada
  pelo reconhecedor. Identidades de predicados e claims usam structs, eliminando
  colisões por delimitadores NUL.
- A API pura limita cem bundles, cem itens por bundle, mil itens totais,
  conteúdo de 2 MiB por item/16 MiB total e campos de 8192 bytes. Achados têm
  até cem limitações de 8192 bytes e datas limitadas a 64 bytes.
- SBOM embutido como string também passa pela validação JSON. Somente declaração
  top-level CycloneDX 1.6 é suportada neste verificador limitado. Versão inventada
  ou sem suporte não vira prova; bom-ref repetido, chave duplicada ou componente
  sem nome/versão explícita são recusados. Até dois mil componentes. Múltiplas
  versões com referências distintas permanecem um controle positivo legítimo.

O verificador SBOM não valida todo o schema nem resolve o grafo de dependências.
Seu resultado permanece DECLARATION_ONLY: presença no artefato, execução,
alcançabilidade e segurança não são inferidas. Integridade dos bytes pode ser
verdadeira mesmo com interpretação UNVERIFIED/UNSUPPORTED. Bundles JSON continuam
sem autenticação de origem; os controles não transformam importação em certificado.

## Verificação

Corpus permanente: parser/bounded_json_test, usecase/technical_semantic_security_test,
local/security_test e cli/security_{advanced_input,git_metadata}_test. Inclui controles positivos,
limites exatos/excedidos, metadados Git sintéticos, aliases, FIFO, hardlinks,
ambiguidade estrutural e conflitos semânticos. A suíte integrada e os resultados
de concorrência/fuzzing são registrados ao concluir a rodada.

Validação concluída em 2026-09-26:

- `go test ./...`, `go vet ./...` e `git diff --check` passaram.
- `go test -race` passou para parser, usecase, CLI e provedor local; o último
  delta de metadata e os campos de enum foram revalidados em corpus específico.
- Fuzzing do validador JSON: 49.391 execuções em cerca de seis segundos, sem falha.
- Harnesses sintéticos bloquearam os bypasses de licença/schema/SBOM/inventário
  dividido e preservaram os controles positivos. Isso não certifica segurança
  de repositórios, ausência de vulnerabilidades ou origem das exportações.
