# 0158 — Validar origem da coleta e estabilizar bundles

## Intenção e observáveis antes da implementação

Corrigir três divergências observadas após0157: respostas de extração com outra identidade/revisão eram atribuídas ao plano; planos de schema desconhecido podiam executar portas; coverage vazio alternava entre null e array, mudando IDs após serialização. Domínio L1 possui formato; casos de uso L1 possuem admissão e vínculo de origem; efeitos concretos continuam em L3. Não alterar política de reúso, classificação semântica ou força de evidência.

## Admissão e origem

CollectInsightEvidence deve admitir todos os planos antes do primeiro efeito, incluindo schema exatamente1. CollectEvidence compartilha a admissão, preservando schema ausente somente na API legada, mas rejeitando versões desconhecidas. Exigir identidade owner/repo válida e igual à do Candidate, quotas não negativas e TotalBytes positivo. Revisões opacas abreviadas da API genérica legada permanecem aceitas; a jornada de insights continua exigindo pin imutável para produzir contexto. HEAD/vazio autorizam somente a resolução explícita existente. Não introduzir novos tetos ou renomear campos.

O plano fixado governa o pin de leitura: metadados atuais não podem trocar esse commit. Metadados resolvidos ou recebidos com outro owner/repo são ErrIntegrity antes de árvore/extração. Resultado nil ou Ref de extração diferente do pin efetivo é ErrIntegrity antes de criar itens de conteúdo. Propagar ErrIntegrity pela coleta coordenada como erro fatal sem relatório utilizável; não convertê-lo em falha de resolução nem fonte não explorada. Falhas operacionais de resolução mantêm o tratamento histórico. Contexto cancelado conserva precedência após portas. Erros não expõem resposta bruta.

## Serialização e compatibilidade

EvidenceBundle.MarshalJSON deve emitir requirement_coverage como array, tanto para nil como para slice vazia não nula, sem mutar o chamador. Após repetidos roundtrips, bytes e hashes de novos bundles devem estabilizar. Consumidores de formulação/exploração devem conservar IDs para a mesma entrada salva/reaberta, inclusive quando recebam null legado.

Nenhum registro, recibo ou hash já emitido será reescrito. Representações antes instáveis convergem para array nas novas serializações; não há como conservar simultaneamente os dois IDs alternantes como um único ID canônico. Bundles com cobertura populada, que já eram estáveis, devem conservar seus bytes e IDs; congelar controles antes da implementação. Replay de gates históricos continua usando seu commit original, sem regenerar evidência para a árvore atual. Registrar a distinção no ADR0022.

## Aceitação e verificação

Controles negativos para owner/repo/commit de resposta, nil, metadado divergente, schema vazio/desconhecido em insights e schema desconhecido na API genérica; nenhum efeito quando um plano inválido estiver no fim do lote. Controles positivos para pin mantido apesar de metadado HEAD mais recente, resolução HEAD válida, ausência de mutação e falha operacional ainda explícita. Testar nil/vazio/populado com múltiplos roundtrips e entrada null legada, identidade de formulação e controles congelados populados.

Preservar os testes anteriores e ajustar somente doubles que não cumpram a identidade da porta, sem mudar seus oráculos funcionais. Acrescentar regressões com prompts próprios antes do código. Executar testes focados, make verify, conferência dos hashes históricos e diff-check. Logs e executáveis ficam fora de L0.

## Resultado e limites da verificação

As regressões de resposta divergente/nula, pin, schemas e roundtrip falharam antes da correção e passaram depois. A admissão coordenada agora recusa o lote inteiro antes de efeitos; a API genérica conserva schema ausente e rejeita desconhecido. Falhas de integridade deixam de produzir conteúdo atribuível ao plano. Metadados de branch não trocam o pin de leitura.

Coverage nil/vazio/null legado estabiliza como array nas novas serializações. Os controles populados foram congelados antes da mudança: digest JSON 6fc91c3dfcf827ad4bcbf0820cc6fbe71fd82d2d58bac8f4c9919d5e4e81598f e ID algorithms_16a7f584dcdd7468cb65d8502de60f2fad2b35e0a8b62d9ab286fa7aa7506e77 permaneceram iguais. A correção das formas vazias não é uma reemissão dos registros históricos alternantes; aplica-se a novas projeções conforme ADR0022.

A primeira execução completa revelou doubles sem Ref em evidence_test, idea_run_test e recommendation_step0025_contract_test. Seus retornos foram ajustados à identidade do candidato solicitado; conteúdos, assertivas e corpos de Test* foram preservados. Também houve bloqueio de loopback no teste httptest do GitHub. A execução final de make verify com loopback autorizado passou: lint sem violações e 11 controles do parser; arquitetura; go test ./...; go vet ./...; build de ./04_wiring em destino temporário; diff-check. Nenhum teste de integração foi desativado para contornar o ambiente.

Logs anteriores/finais e executáveis estão em /tmp/bateia-0158-before.log, /tmp/bateia-0158-focused-final.log e /tmp/bateia-0158-verification-final. São diagnósticos temporários locais, não evidência selada. A comparação SHA256 confirmou 4.683 arquivos históricos de passos, laudos e Lab intocados, excluído somente o índice de passos atualizado para navegação. Gates selados com módulo próprio não foram reproduzidos; permanecem vinculados a seus commits originais.
