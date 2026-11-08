# 0160 — Preservar revisões e orçamentos de coleta

## Intenção anterior à materialização

Corrigir os quatro contraexemplos da revisão de2669ead. A capacidade é coleta/recuperação da investigação, nos responsáveis L1 evidence e evidence_recovery. Conservar schemas, política probatória, registros e hashes históricos. Atualizar somente contratos ativos e acrescentar regressões próprias, sem reescrever testes anteriores ou gates selados.

## Comportamento esperado

Diagnóstico, agrupamento de planos e união de bundles usam repo@commit completo. Duas revisões do mesmo repositório são coletadas e avaliadas separadamente, inclusive quando compartilham path ou hash. A árvore preparatória usa o pin do target, não o commit antigo do candidato. A união conserva os itens antigos no próprio commit e não altera slices do chamador. HEAD pode resolver, mas não autoriza fundir conteúdos de revisões diferentes; o registro original continua explícito.

Na coleta geral, conferir contexto imediatamente após FetchTree e antes/depois da busca opcional de implementação. Cancelamento ou deadline interrompem sem iniciar a próxima porta. Manter as regras de integridade e seleção anteriores.

AdditionalBytes limita a soma de novos BytesRead da rodada, incluindo metadados. Distribuir o saldo em ordem canônica de repo@commit. O saldo restante limita a próxima coleta; saldo zero produz bundle incompleto com targets SKIPPED_BUDGET, sem portas de coleta para esse grupo. Targets omitidos por falta de saldo não entram em TargetsExecuted. Não cobrar os bytes originais nem interpretar omissão como ausência. As quotas de tipos de arquivos conservam o comportamento anterior; este recorte corrige o teto agregado de bytes.

TotalBytes inclui metadados. Admitir seu JSON integral apenas quando couber; caso contrário, omitir o item, registrar SKIPPED_BUDGET e Complete=false, sem truncar JSON ou produzir hash de objeto incompleto. Os demais conteúdos continuam usando o saldo disponível. BytesRead nunca excede o teto admitido. Não reserializar dossiês antigos.

## Aceitação e limites

Regressões públicas com plano produzido pelo planejador para duas revisões, mesmo hash/path entre revisões e original preservado; diagnóstico com falha de uma revisão isolada; cancelamento/deadline após árvore e busca; orçamento agregado em dois/três candidatos, incluindo saldo zero sem portas; metadados abaixo, iguais e acima do limite, verificando JSON/hash e contabilização. Demonstrar falhas antes do patch e passagem depois, sem modificar consumidores históricos.

Executar make verify com V0–V27 selecionadas, disabledRules vazio e análise completa; conservar exclusões históricas explícitas. Comparar SHA256 do acervo, conferir links/linhagem e git diff. Logs temporários no /tmp. Gates selados não são reproduzidos pelo módulo raiz nem promovidos a resultados novos. Passagem estrutural não prova semântica dos provedores, efeito ou utilidade dos insights.

## Resultado e verificação

Os quatro contraexemplos falharam antes da correção. As regressões próprias agora passam, incluindo separação de revisões com/sem originais e hashes iguais, isolamento das causas, pin na preparação, cancelamento/deadline em árvore/busca, distribuição canônica do saldo, candidato derivado omitido sem enumeração, bytes originais não cobrados e fronteiras de tamanho do JSON de metadados. O cenário de memória sem árvore foi ajustado para não solicitar um README vazio incidental; os testes consumidores anteriores não foram alterados.

O diagnóstico, a preparação e a união usam chaves repo@commit. A preparação de padrões agora ocorre por grupo somente quando há saldo, antes da coleta restrita desse grupo. O teto adicional passa ao coletor como saldo; grupos omitidos conservam os targets e motivo de orçamento. A métrica TargetsExecuted exclui esses grupos. O coletor omite JSON de metadados excessivo, mantém incompletude e observa contexto entre portas. Não foram relaxadas regras de evidência ou reescritos dossiês antigos.

A execução de make verify em /tmp/bateia-0160-verification passou em lint, arquitetura, go test ./..., go vet ./..., build e diff-check, com loopback autorizado para fixtures HTTP e assinatura Git desativada somente no processo de fixtures pelo runner canônico. Alcance efetivo: V0–V27 selecionadas, disabledRules vazio, analysisComplete=true, sem violações, erros de leitura ou sintaxe. Os 11 controles de parser e os dois controles V3 ligado/desligado também passaram. As exclusões históricas explícitas não mudaram; seleção de regras não implica prova semântica.

Logs temporários: /tmp/bateia-0160-before.log e /tmp/bateia-0160-focused.log registram as regressões; /tmp/bateia-0160-after-probes.log reaplica os contraexemplos da revisão. Foram lidas ambas as revisões, a busca após árvore cancelada teve zero chamadas, o bundle com limite1 ficou dentro de1 byte e a rodada de400 bytes admitiu400, em vez de800. Esses logs são diagnósticos locais, não recibos selados.

A comparação SHA256 confirmou 4.687 arquivos históricos de Lab, passos e laudos intactos. Somente o índice ativo de passos foi excluído dessa comparação; os hashes das fontes/patches do linter e gates antigos permaneceram iguais. O diff conserva os testes anteriores. Links locais dos documentos alterados foram conferidos. Não houve atualização de hashes/snapshots por ferramentas de reparação nem reprodução dos módulos históricos selados como resultados novos.
