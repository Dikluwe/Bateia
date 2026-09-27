# CLI idea recommend

Exponha uma boca de inspeção única para live, local e replay. Valide saída e
parâmetros antes de provedores, grave JSON atomicamente, rejeite sobrescrita sem
resume, preserve file:// e hash de conteúdo local, nunca serialize segredos e
gere manifesto, artefatos intermediários e projeção Markdown rastreável.

## Passo 0040

Aplicar o contrato de autoridade, sinais lexicais sem comprovação comportamental, procedência local, seleção de produção, negação conservadora e lacunas acionáveis de `00_nucleo/passos/0040-corrigir-autoridade-da-evidencia.md`. Preservar determinismo e ausência de I/O no núcleo.

## Evolução 0043

# Passo 0043 — Plano verificável de investigação

Origem: 00_nucleo/passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md.
L1 puro: transformar IdeaSpec + DiscoverySnapshot + avaliações em plano determinístico, sem IO/relógio. Requisitos recebem IDs por texto normalizado e tipo, ordenados; separar problema, capacidades, restrições, não objetivos e preferências. Covers indexados de expansões referem-se à ideia original do snapshot, sem reinterpretar índice após ordenação. Query alterada tem alteração explícita sem inventar filtros removidos; ausente permanece desconhecida. Propostas de busca não são testemunhos.

BuildInvestigationPlan retorna requisitos, tentativas observadas, ações por candidato/requisito, classes de evidência, método, lacuna, critério de aceitação e encerramento delimitado. Sanitizar avaliações pelo contrato comum antes de atribuir status; sem candidatos produzir NO_CANDIDATES_IN_SCOPE somente sem falhas, nunca inexistência global. Queries e fontes opacas não geram prova. Plano sem snapshot informa DISCOVERY_NOT_RECORDED. Resultado ordenado independe da ordem dos conjuntos de entrada. Sem ensaio de integração sempre preservar condição pendente de validação.

Anexar plano à decisão canônica e à apresentação CLI/MCP. Projetar texto público ao construir o plano com os filtros existentes de privacidade; preservar o gateway certificado. Telemetria identificada como da ferramenta. Compatibilidade aditiva com schema aninhado próprio; não promover resultados históricos.

Identidade 0043: registrar versão do analisador/promessas e invalidar checkpoints de evidência sem essa versão; código de coleta e autoridade mudaram mesmo quando catálogo lexical não mudou.

Legibilidade 0043: cada ação contém pergunta e lacuna explícita. Sem detector para o requisito, pedir especificação/verificação delimitada, sem fingir que mais coleta resolve ausência de analisador. Markdown mostra texto do requisito além do ID e fontes pertinentes quando presentes.

## 0044
Adicionar versão de recuperação implementation_search=0044 à identidade/retomada somente em modo local, que ganhou nova coleta opcional. Replay usa bundle congelado; provedores remotos sem porta preservam identidade.

0045: versão implementation_search do checkpoint local passa a 0045; incompatibilidade impede retomada silenciosa do ranking anterior. Replay mantém evidência da origem.

0046: checkpoint local implementation_search passa0046 pela adição de lacunas por requisito; replay conserva bundle.

0046: decision.md também apresenta projeções dos bundles, após razões/capacidades/planos e antes links de evidência. Reutilizar renderizador L2 de requisitos; até cinco candidatos com identificação repo@commit. writeDecisionMarkdown aceita bundles opcionais para compatibilidade dos chamadores anteriores. Replay renderiza as projeções preservadas. Sem modificar decisão canônica.

0047: se candidato for promovido, implementation_search0047 no checkpoint local; manter replay congelado e sem promoção de capacidades.

0048: checkpoint local implementation_search0048; replay conserva interpretação e restrições congeladas.

0048: decision.md também renderiza QueryInterpretation de cada bundle antes de Requirements, para expor modo e símbolos como analyze. Usar WriteQueryInterpretation do renderizador L2, mantendo limites de candidatos/linhas existentes.

0050: checkpoint local implementation_search0050 e versão de evidência atualizada distinguem coleta com suporte remoto. Replay conserva artefatos da origem.

0050: versionar implementation_search tanto em local como live para impedir retomada de checkpoint remoto anterior sem busca de implementação. Replay continua usando bundles congelados.

0051: implementation_search0051 em local/live pela mudança de recuperação documental. Replay preserva bundles.

0052: checkpoint implementation_search0052 em local/live para novo catálogo lexical; replay conserva bundles congelados.

## Passo 0055 — Retomada
Aplicar 00_nucleo/passos/0055-adendo-checkpoint.md: discovery_query 0055-v1 somente modo live; comparar chave na elegibilidade de checkpoint para impedir reuso de descoberta obsoleta.

Revisão sintaxe 0055: a chave discovery_query live passa para 0055-v2; invalida retomada da primeira serialização.

## Passo 0056
Aplicar 00_nucleo/passos/0056-explicar-cobertura-e-perdas.md: report exibe cobertura desconhecida, motivos e limites sem inferir ausência de capacidade; checkpoint live discovery_query passa0056-v1.

Correção de alinhamento0056: identificador exato de discovery_query no modo live é `live0056-v1`, como no contrato/teste congelado; a primeira implementação usava0056-v1 e falhou no gate. Não alterar testes ou resultados de busca.


## Passo 0057
Checkpoint implementation_search passa a 0057 em local/live. discovery_query permanece live0056-v1.


## Adendo 0057 — Origem documental
Checkpoint implementation_search passa a 0057-r1 em local/live; descoberta permanece intacta.

## Passo 0058
Checkpoint implementation_search passa a 0058 em local/live. Preservar
discovery_query e a descoberta; replay conserva os bundles congelados.

## Passo 0059
Checkpoint implementation_search passa a 0059 em local/live. discovery_query permanece
inalterado e replay conserva bundles congelados.

## Passo 0060
Checkpoint implementation_search passa a 0060 em local/live para a tradução
lexical de caminho na consulta. discovery_query permanece inalterado e replay
conserva bundles congelados.


## Passo 0061
Checkpoint implementation_search = 0061 local/live; descoberta inalterada. Apresentar WriteActionAssessments ao lado de requisitos no Markdown.

## Modelo de peneira opt-in

Adicionar `--screening-model <arquivo>` a `idea recommend`. Sem a flag, preservar
exatamente o fluxo, versões, `InputHashes`, pré-imagem do `RunID` e elegibilidade
de retomada existentes. Com a flag, antes de criar output ou chamar dependências:

- ler bytes uma vez e fazer strict-decode de um único `domain.ScreeningModel`;
- rejeitar trailing data/campos desconhecidos, números não finitos, schema ou
  thresholds inválidos, status diferente de `uncalibrated` e calibrador presente;
- calcular SHA-256 dos bytes exatos e gravá-lo em
  `InputHashes["screening_model"]`;
- usar `RuleVersions["screening"]="screening-opt-in-v1"`;
- inserir `screening_model_hash` depois de `target_hash` e antes de
  `rule_versions` na pré-imagem opt-in do `RunID`;
- exigir presença/hash/versão idênticos para `--resume`;
- injetar o modelo somente por `IdeaRunDependencies.ScreeningModel`.

No `live` sem target/checkpoint, a baseline usa `RecommendationService`, que não
expõe injeção desse modelo. Somente no caminho opt-in, não selecione esse serviço:
use o ramo direto já existente de `RunIdeaPipeline` com o modelo em `deps` e depois
o ranking de alternativas já existente. Não altere L1 nem adicione setter. Sem a
flag, preserve exatamente a condição e o serviço atuais.

Mesmo conteúdo em outro caminho mantém identidade; qualquer byte diferente muda
identidade. Presença versus ausência nunca é compatível para retomada. Erro não
deixa diretório/manifesto parcial. Falha de leitura retorna somente `não foi
possível ler o modelo de screening`; conteúdo inválido retorna somente `modelo de
screening inválido`; nunca interpolar path, conteúdo, campo, valor ou erro JSON.
Não ler modelo de env, default oculto ou `lab/`; não alterar `DefaultScreeningModel`,
`ScreenDiscovery`, schemas públicos nem o pipeline L1.

## Segurança de entrada

Usar readIdeaInput nonblocking/regular/16MiB para hashes de replay. JSON de ideia, policy, rules, checkpoints e modelo passa pela guarda compartilhada antes da desserialização; duplicatas incluindo capitalização/escapes são recusadas. Preservar mensagens de modelo inválido, hashes de arquivos e certificados históricos. Contrato em security_json_input.md.

0096: hash local usa security_tree_hash.md, formato tree-sha256-v2 e orçamento agregado compartilhado entre raízes. local_snapshot participa da identidade e comparação de checkpoint; versão antiga/ausente não é reutilizada. Iniciar contexto único com timeout positivo antes da admissão; só criar dossiê após todos os hashes válidos. Não manter a pré-imagem ambígua v1.

Metadata Git para remapeamento usa readLocalGitConfig: os.Root único, caminhos/pointer gitdir/commondir somente internos, recusa de symlink/hardlink/FIFO e validação pós-leitura. Unsafe/indisponível=>nil identities; linked worktrees com administração externa deixam identidade remota desconhecida. Metadata local é declaração, não autenticação da origem. Contrato em security_git_metadata.md.

0097: modo local captura na mesma leitura do hash e usa SnapshotProvider sem reler origem. Contrato local_snapshot_capture.md; local_collection=frozen-local-v1 participa dos checkpoints. local_metadata participa de InputHashes e RunID por IDs consultivos congelados quando há target. Stages locais registram checkpoint_reused0/1. Hash v2, pesquisa simples e live/replay preservados; limpar capturas em qualquer saída.

## Passo 0098 — Trace opt-in

--trace gera trace-graph.json consultivo e estágio trace com hash dos bytes.
Preservar bytes causais, RunID/RuleVersions e critérios de checkpoint. Validar
saída contra entradas e aliases antes do dossiê; não criar grafo sem flag. Falha
do trace marca execução FAILED. Retomada sem flag preserva grafo antigo no disco,
não o consome/anuncia no manifest e avisa no stdout. Manifest/hashes, não simples
existência do arquivo, identificam pertencimento à execução. Hash não autentica.
## Passo 0099 — Índice legível opt-in

A mesma flag --trace agora gera também trace-index.json via ExplainDecisionTrace,
usando exatamente o build input do grafo. Sem flag não criar índice. Não adicionar
flags, inputs causais, RuleVersions ou gate de checkpoint. Preflight ambos filenames
contra idea/policy/rules/screening, replay discovery/bundles quando consumidos,
checkpoints discovery/bundles em resume, e mutuamente entre os sidecars; recusar
aliases/hardlinks/symlinks antes de prepareDossier/run.json. Revalidar na escrita.
Computar ambos antes de escrever; falha de explicação/escrita/hash deixa run/trace
FAILED com outputs vazios. Stage complete contract decision-trace-projection-v2,
outputs grafo primeiro e índice segundo com hashes de arquivo e counts entries/
locations. Escrita individual atômica não equivale a transação de dois arquivos.
Resume sem flag preserva ambos sidecars antigos no disco, exclui stage trace e
avisa que a projeção não pertence à execução atual. Decide sidecar continua grafo
somente. Declarações legíveis não são validação de verdade/completude/autoria.
