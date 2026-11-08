# 0159 — Fechar cancelamento, escopo e recibos

## Intenção anterior à materialização

Corrigir os três contraexemplos reproduzidos após0158: releitura técnica continuava após cancelamento e emitia VERIFIED; recuperação solicitava arquivos além dos alvos congelados; expansão admitia recibos COLLECTED sem bundle. Responsáveis em L1: technical_recheck, evidence_recovery/coletor e insight_neighborhood_sources. Preservar formatos e políticas de evidência, os registros/hashes já emitidos e os corpos dos testes anteriores.

## Comportamento esperado e limites

Releitura técnica: rejeitar contexto nil/cancelado; observar contexto antes/depois de cada porta e antes/depois da verificação final. Cancelamento/deadline devolve erro de contexto e relatório vazio, sem VERIFIED nem nova porta após a interrupção. Candidatos sem achado ativo não dispensam a verificação do contexto.

Recuperação: congelar os caminhos explícitos ou os que foram derivados pelos padrões dos targets na etapa preparatória; coleta posterior não amplia essa seleção pela árvore, recuperação de implementação ou arquivos extras do provedor. Reutilizar o coletor por modo interno exclusivo de targets, sem novo campo público ou alteração do modo de busca geral. Arquivo elegível conhecido fora do recorte mantém incompletude explícita; erros de identidade e caminhos extras em resposta são rejeitados sem conteúdo útil. Metadados e enumeração para conferir escopo não autorizam leituras adicionais. Cancelamento não vira falha operacional recuperável.

Expansão: correspondência um-para-um entre recibos primários COLLECTED e bundles GitHub/contextos por sítio, identificada pelo par connector_id/candidate_id. Rejeitar órfão, recibo nil, status interno/externo incoerente e duplicidade antes de qualquer leitor. Recibos de followup continuam no próprio followup, com as regras anteriores. Status não coletado sem bundle continua válido. Atualizar SRC-01 sem declarar cobertura completa das demais combinações.

## Linter e preservação

Verificar a mudança enabled por regra disponível no checkout limpo do tekt-linter e atualizar o pin ativo do Bateia para a revisão 20eeb9b, que inclui 6c4bf23. Preservar source.tar.gz e copiar seu manifest anterior sem reescrever bytes; novo arquivo de fonte tem nome/revisão próprios. Manter o patch do parser Go, severidades e controles anteriores. Configurar V0–V27 explicitamente ativas no Bateia e emitir alcance efetivo; a seleção deve ser comprovada por controles em que regra ligada acusa mutação e desligada não acusa. Não usar flags de reparação de hashes/snapshots. A fonte nova, hashes, runner e política pertencem ao Lab; nada disso integra o produto.

## Aceitação

Regressões com cancelamento/deadline disparados pelas portas; alvos explícitos e derivados com arquivo fora do plano e provedor que devolve extras; recibos primários órfãos/nulos/duplicados/incoerentes em entradas resseladas; controles positivos anteriores preservados. Executar casos antes/depois, controles on/off reais do linter, make verify, links/linhagem e comparação SHA256 do acervo histórico. Logs e binários ficam em diretório temporário autorizado, nunca em L0. Gates selados permanecem no commit e módulo próprios.

## Resultado e verificação

Os três contraexemplos falharam nas regressões antes da correção e passaram depois. Recheck interrompe após árvore/extração cancelada ou com deadline, sem relatório utilizável; recuperação solicita somente targets explícitos/derivados e rejeita extras; a guarda de recibos rejeita órfãos, nil, duplicatas e status incoerentes antes de leitores. Controles positivos conservam o modo geral de coleta e históricos sem coleta atual.

A auditoria0145 inicialmente recusou a mudança do corpo congelado de insightNeighborhoodOrigins. Esse corpo foi restaurado integralmente e a nova guarda de correspondência foi colocada no mesmo responsável de fontes, sendo delegada pela admissão antes do verificador histórico. O gate, seus hashes e consumidores não foram alterados. A execução posterior aprovou a auditoria de preservação e os demais checks arquiteturais.

O pin anterior e53de6f não continha o switch enabled por regra. O gate atual usa 20eeb9b2f54b1c55833c0edbae1a3f9d8813a8bf, com fonte própria e o mesmo patch Go. source.tar.gz permaneceu intocado e manifest-e53de6f.json preserva byte a byte o manifest anterior. A política do Bateia declara V0–V27 enabled=true; o alcance efetivo confirmou 28 regras selecionadas, disabledRules vazio e analysisComplete=true, sem erros de leitura/sintaxe. Durante a preparação, o linter detectou o prompt ainda sem consumidor; após a materialização, terminou sem violações. Os 11 controles anteriores e os dois controles reais V3 on/off passaram. Também passaram os 9 testes upstream analysis_policy_cli da fonte fixada, incluindo sobrescrita de seleção, regras fatais, configuração inválida e análise incompleta. Seleção não implica aplicabilidade de toda regra a toda linguagem nem prova semântica do produto.

A execução final de make verify passou em lint, arquitetura, go test ./..., go vet ./..., build e diff-check, com loopback autorizado para httptest. Logs finais e binário estão em /tmp/bateia-0159-verification-final; controles anteriores/focados em /tmp/bateia-0159-before.log e /tmp/bateia-0159-focused.log; testes do switch em /tmp/bateia-0159-rule-policy-tests.log. São diagnósticos temporários locais, não novos recibos selados.

A comparação SHA256 confirmou 4.680 arquivos históricos de passos, laudos e Lab intactos, incluindo a fonte anterior e o patch do linter. O índice de passos e os arquivos ativos de configuração/runner/documentação do gate foram excluídos dessa contagem por serem parte autorizada da mudança. O checkout upstream do linter permaneceu limpo. Nenhum hash de resultado antigo foi recalculado; gates selados continuam vinculados a seus commits e módulos próprios, sem reprodução automática pelo gate raiz.
