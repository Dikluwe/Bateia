# 0105 Investigação independente de insights

**Data**: 2026-10-06  
**Contrato**: [passo 0105](../passos/0105-investigar-insights-sem-decisao-de-reuso.md)  
**Decisão**: [ADR-0010](../ADR/ADR-0010-investigacao-offline-de-insights.md)  
**Estado**: IMPLEMENTADO E VERIFICADO NO RECORTE OFFLINE

## Resultado

Foi acrescentada a jornada `idea insights analyze`, independente de IdeaSpec, política de adoção, decisão de reúso, run_id e artefatos de recomendação. O pedido declara pergunta, contexto do alvo, propriedade, efeito desejado, restrições, classes de mecanismo e fontes congeladas.

L1 valida e normaliza dados limitados, confere commits/bytes/hashes e usa o extrator determinístico existente. O relatório conserva proveniência, identidade causal, mecanismos selecionados, lacunas e hipóteses. O alvo permanece UNKNOWN e nenhuma hipótese se torna SUPPORTED. O estado de negócio permanece PARTIAL ou INCONCLUSIVE, pois a aplicabilidade ao alvo ainda não foi avaliada.

L2 lê JSON estrito com limites estruturais e preserva entradas/aliases e saídas em erros; usa escrita atômica e bloqueio de segredos existentes. L4 despacha antes de .env, vocabulário, clientes, cache e runtime. A síntese durável e a recomendação anteriores não foram alteradas.

## Experimento ilustrativo

A [bancada](../../lab/bancada/insights/0105/README.md) registra pergunta sobre sobrecarga do Bateia usando uma fixture documental de fila limitada. Owner/repo/SHA são identidades ilustrativas, não uma origem externa autenticada. A licença e a avaliação PARTIAL são declaradas.

O binário compilado retornou mecanismo PARTIAL, hipótese INCONCLUSIVE, alvo UNKNOWN e três etapas propostas: baseline, aplicabilidade e experimento. Nenhuma adoção foi autorizada e nenhum efeito executado. O comando passou com HOME ausente, raiz de runtime inválida e .env no diretório temporário, sem criar cache ou runtime. A saída completa e as condições do ensaio ficam no Lab.

## Cenários verificados

Testes cobrem fonte com implementação/testemunha independente, documentação, coleção incompleta, classe não observada, avaliações ausentes ou sem binding, hashes adulterados, identidade cruzada, duplicatas, classes/limites inválidos, cancelamento e orçamento de bytes. Pergunta/contexto alteram identidade; pedido original permanece inalterado e repetições são determinísticas.

Uma licença restritiva declarada não impede gerar hipótese de referência e não concede direito de transferir código. A força do mecanismo e a força da hipótese/estado do alvo permanecem distintas. CLI cobre JSON desconhecido, múltiplo ou duplicado, ausência de publicação em falha, input/alias protegidos, arquivo/stdout equivalentes e short writes.

## Verificação técnica

Lint, `go vet ./...` e build de `./04_wiring` passaram. Todos os testes focados de domínio/caso de uso/CLI e o ensaio do binário offline passaram. A suíte concorrente apresentou duas falhas em testes existentes com janelas temporais: lifecycle do container ficou RUNNING no prazo e o gate R5 não distinguiu M28 na janela de 100 ms.

O teste de lifecycle passou três vezes isoladamente. Na repetição da suíte com `go test -p 1 ./...`, todos os pacotes e gates ativos passaram, incluindo container e R5, sem mudança de código ou critério. Essa execução aprovada não apaga as falhas temporais observadas na execução concorrente; a causa de scheduling/carga continua uma hipótese, não uma prova. Loopback foi permitido e assinatura Git desativada somente no processo de fixtures. Gates selados do módulo separado não foram reproduzidos.

## Limites

A jornada recebe avaliações fornecidas; verifica seu vínculo com os bytes, não refaz automaticamente suas conclusões semânticas ou autenticidade externa. O catálogo permanece fechado às quatro classes existentes. Sem fontes selecionadas e congeladas não há descoberta automática nesta interface.

O efeito desejado e as restrições são intenção, não observações. O plano ainda precisa de métricas/limiares específicos e validação de relevância/aplicabilidade para um caso real. O experimento ilustrativo demonstra funcionamento e limites da jornada, não utilidade comprovada para o alvo ou ganho da IA.
