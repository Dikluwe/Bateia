# 0164 — Preservar UTF-8 no truncamento e replay

## Intenção anterior à materialização

Fechar o único achado confirmado na revisão de d826623: cortar um arquivo textual UTF-8 no meio de um caractere gera bytes inválidos, substituídos na serialização JSON sem alterar o hash anterior. Responsável L1: evidence, compartilhado pelos fluxos geral, coordenado e de recuperação. Preservar schemas, hash dos controles ASCII/completos e todos os registros históricos. Não reescrever bundles antigos nem atualizar hashes para mascarar o defeito.

## Comportamento esperado

O conteúdo textual bruto admitido precisa ser UTF-8 válido. Recusar arquivo recebido já inválido com FAILED/falha de codificação, sem item, em vez de fabricar bytes com substituição. Conservar as verificações de tamanho bruto anteriores ao corte.

Quando o saldo de bytes limita uma fonte válida, conservar o maior prefixo que termina numa fronteira de caractere UTF-8 dentro desse saldo. Não preencher com U+FFFD, normalizar texto ou ultrapassar o teto. O hash e os bytes contabilizados referem-se exatamente ao prefixo preservado. Se o primeiro caractere não couber, registrar SKIPPED_BUDGET/Bytes0 sem criar item vazio; o saldo não usado pode servir a outro arquivo. Fragmentos continuam incompletos, sem fonte dirigida examinada ou promoção de força. ASCII e conteúdo completo conservam bytes/hashes.

Serializar/reabrir bundle emitido não altera o conteúdo nem seus hashes: reavaliar a cópia deve continuar aceitando a integridade, inclusive quando o conteúdo é apenas fragmento. Preservar SKIPPED_BUDGET, motivo, histórico de retomada e métricas0163. Não prometer completar grafemas ou funções: a fronteira é de codificação UTF-8, não de semântica do trecho.

## Aceitação e preservação

Regressões públicas em todas as posições de corte de texto com ASCII e caracteres de2/3/4 bytes, incluindo marca combinante, controles antes/depois do caractere e conteúdo completo. Esperados independentes por prefixos de runes; conferir JSON em múltiplas reaberturas, igualdade de conteúdo/hash, contagem de bytes e reavaliação. Primeiro caractere sem espaço não cria item; saldo restante alimenta arquivo ASCII posterior. Rejeitar fonte originalmente UTF-8 inválida sem sanitização silenciosa. Conferir jornadas coordenada e de recuperação e preservação dos inputs. Testes anteriores intactos.

Executar regressões antes/depois e make verify com V0–V27 selecionadas, disabledRules vazio e análise completa. Comparar SHA256 do acervo, conferir links/linhagem e diff. Logs/binários no /tmp; gates selados continuam vinculados aos próprios commits/módulos, sem reprodução ou reparação de resultados históricos.

## Resultado e verificação

O contraexemplo falhou nas regressões antes do patch e passou após a correção. O coletor valida a codificação bruta, recua cortes em fronteiras de caractere antes do hash e não cria item quando nenhum caractere cabe. Só os bytes efetivamente preservados são cobrados; o saldo pode atender outro arquivo. Fragmentos continuam SKIPPED_BUDGET/incompletos e retomáveis, sem ganhar autoridade.

As regressões percorrem todas as posições de corte da fixture com ASCII, caracteres de2/3/4 bytes e marca combinante; cada bundle é serializado/reaberto e reavaliado três vezes, com prefixos, hashes e bytes estáveis. Controles completos e ASCII permanecem, fonte originalmente inválida é recusada sem substituição, e as jornadas coordenada/recuperação preservam publicação, reavaliação e histórico. Os inputs e todos os testes anteriores permaneceram intactos.

A reaplicação do probe original confirmou os controles22/24/40 e corrigiu o limite23: antes ele dividia `é`, gerava UTF-8 inválido e falhava após JSON; agora preserva o prefixo de22 bytes e mantém bytes/hash aceitos antes e depois da reabertura. Logs temporários em /tmp/bateia-0164-before.log, /tmp/bateia-0164-focused.log, /tmp/bateia-0164-core.log e /tmp/bateia-0164-after-probe.log. São diagnósticos locais, não recibos selados.

make verify passou em lint, arquitetura, go test ./..., go vet ./..., build e diff-check, com loopback autorizado para fixtures HTTP e assinatura Git desativada somente no processo das fixtures pelo runner canônico. Logs finais/binário em /tmp/bateia-0164-verification. O alcance efetivo confirmou V0–V27 selecionadas, disabledRules vazio, analysisComplete=true e zero violações/erros de leitura/sintaxe. Os 11 controles do parser e dois controles V3 ligado/desligado passaram; exclusões históricas explícitas permaneceram. Passagem estrutural não autentica origem nem comprova a interpretação dos trechos.

A comparação SHA256 confirmou 4.691 arquivos históricos de Lab, passos e laudos intactos; o índice ativo de passos foi excluído por integrar esta atualização. Fontes/patches do linter e acervo de gates mantiveram hashes. Nenhum teste anterior foi editado; links, linhagem e diff foram conferidos. Não foram reparados bundles antigos, recalculados hashes/snapshots históricos ou reproduzidos módulos selados como resultados novos.
