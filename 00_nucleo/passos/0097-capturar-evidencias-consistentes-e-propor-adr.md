# Passo 0097 — Evidências consistentes e ADR sem aprovação automática

Implementação das duas primeiras prioridades da rodada de melhorias. O modo
simples de pesquisa continua disponível; live e replay mantêm seu fluxo.
Não foram criadas assinaturas, execução de candidatos, permissões novas ou um
piloto humano. Trabalharam agentes em frentes separadas de ADR, provedor congelado
e revisão read-only, com captura/integração/testes na tarefa principal.

## ADR é proposta, não aprovação

O gerador antes escrevia Aceito automaticamente e prometia redução de dívida
técnica com padrões testados em produção. Agora escreve PROPOSTO, descreve o
primeiro resultado como candidato à avaliação e apresenta benefício como
hipótese. Score não demonstra comportamento, qualidade ou aprovação humana.
Dependências são observadas, não aprovadas. Referências vazias/HEAD são marcadas
como não fixadas. Resultado vazio não prova inexistência de alternativas.
Consulta, contagens, score, justificativa, URL e revisão válida são preservados.

## Mesmos bytes para identidade e análise local

`idea recommend --mode local` captura os arquivos em diretórios privados antes
de criar o dossiê. A mesma leitura escreve conteúdo na cópia e no hash v2;
nunca se calcula hash da origem para depois reler/copiar outra versão.

O novo SnapshotProvider recebe bindings exatos de owner/repo/URL/CommitID para
as cópias e atende manifesto, árvore, extração, busca de implementação e metadata
somente ali. Local desconhecido ou revisão diferente dá erro: não há fallback
para origem ou fakeprovider. Identidades públicas mantêm URL original, não a
temporária. Erros/diagnósticos não expõem o caminho privado.

`tree-sha256-v2` e seus hashes permanecem iguais para o mesmo inventário admitido.
Adiciona-se `local_collection=frozen-local-v1` ao manifesto e à comparação de
checkpoint. Checkpoints anteriores sem essa versão são recalculados. Os stages
locais discover/collect registram `checkpoint_reused`0/1 para auditoria.
Certificados históricos não foram reescritos.

## Recursos e ciclo de vida

Continuam os limites0096:20 raízes,20mil entradas,64MiB agregados,16MiB/arquivo,
profundidade64 e paths4096. Hash/captura compartilham o orçamento entre raízes.
Coleta posterior preserva seus próprios limites mais restritivos. A captura
pode conter arquivos que a coleta não chega a examinar; isso não demonstra
cobertura de análise total.

Container MkdirTemp e diretórios0700; arquivos0400 sem preservar permissões de
execução. Symlinks/hardlinks/FIFO não são capturados como arquivos admitidos.
Diretórios vazios são retidos, mas continuam fora do digest. Ignorados e output
descendente seguem a política0096. Raiz que inclui o armazenamento temporário
(como `/tmp`, se esse for o temp do host) é recusada para evitar autocaptura.
Também se recusam identidades owner/repo duplicadas: o sufixo curto derivado do
path não pode criar vínculo ambíguo entre raízes.

Capturas são limpas em sucesso, retomada e qualquer erro, inclusive recusa do
output. A limpeza valida a identidade do container criado antes de removê-lo;
alvo substituído não é apagado. Falha de limpeza é reportada. Nenhum artefato
preexistente é removido pela admissão. Não se executa código capturado.

## Metadados consultivos congelados

`.git` diretório não integra hash/cópia; BaseCommit não é inventado. O provedor
declara metadata Git indisponível e bytes sem autenticação por Git.

Se há target, os remotes consultivos são observados uma vez na admissão, usando
o leitor confinado com config até64KiB por raiz (mais ponteiros de até4KiB).
Indisponível/inseguro/excessivo resulta em identidade desconhecida. O mapa de IDs
normalizados é congelado e usado após pipeline sem reler a origem. Seu fingerprint
`local_metadata` participa de InputHashes/RunID, mas não é prova de origem.

Uma mudança de hints com arquivos idênticos pode reutilizar evidências de arquivos
e recalcular ranking com os hints novos; a identidade da execução muda. Mudanças
irrelevantes de config que não alteram os IDs normalizados não mudam esse fingerprint.

## Limitações

Consistência dos bytes do workflow não é snapshot de ponto único: arquivos podem
ter sido observados em instantes diferentes durante a captura. Também não é
autenticação, equivalência com commit Git ou imutabilidade contra outro processo
com mesmo UID/administrador. Diretórios vazios/metadata completa não integram a
identidade dos arquivos. Cancelamento permanece cooperativo, sem deadline de SO
para syscall já bloqueada. Políticas/raízes são escolhas explícitas do usuário;
os temporários não substituem uma sandbox de execução.

## Verificação

Corpus parser ADR; SnapshotProvider (incluindo limites de árvore e fallback);
CLI local_snapshot_capture/cli e regressões de hash/retomada. Inclui controles
positivos, mutação/remoção da origem, hash estável, privacidade, quotas, limpeza,
Git consultivo e versões de checkpoint.

Validação final em 2026-09-26: go test ./..., go vet ./... e
git diff --check passaram. go test -race passou no corpus ADR, captura/CLI,
SnapshotProvider e regressões locais de hash, retomada e identidade.
