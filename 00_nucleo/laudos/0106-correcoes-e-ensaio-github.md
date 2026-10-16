# 0106 Correções e ensaio com GitHub

**Data**: 2026-10-06  
**Estado**: VERIFICADO NO RECORTE  
**Contrato**: [passo 0106](../passos/0106-corrigir-esperas-e-investigar-github.md)

## Resultado

Foram corrigidas as duas esperas frágeis registradas no laudo 0105. O teste
do container continua exigindo SUCCEEDED persistido e reabertura bem-sucedida,
mas trata o prazo como watchdog de até 30 segundos, respeita o deadline do
teste, reduz polling para 10 ms e falha em outro estado terminal com diagnóstico.
Isso não estabelece SLA de latência nem altera o runtime produtivo.

O gate ativo R5 observa M28 por rendezvous no caminho que espera recovery
ou conclusão real de Start com a barreira fechada. A ausência de conclusão
em 100 ms deixou de decidir o veredito. O watchdog agora só falha o teste.
Referência e mutante continuam produzindo observáveis comparados; 60/60
mutações permanecem Violated nas três ordens e três repetições. Receitas,
goldens e recibos históricos não foram reescritos.

## Fonte real e investigação

O adaptador existente GitHub foi exercitado pelo binário usando coleta pública
sem token. O repositório sourcegraph/conc foi fixado no commit
5f936abd7ae87036af1f75c95fb9d0daaf00116b. A coleta registrou 29 itens/95.340
bytes, complete=true e nenhuma falha no escopo do coletor. Os bytes dos dois
arquivos testemunha foram relidos online diretamente nesse SHA e coincidiram.

A [bancada reproduzível](../../lab/bancada/insights/0106/README.md) conserva
planos, bundle, avaliações, observações, resultado, controles e logs. A
implementação configura reserva finita antes de criar workers; seu teste
observa concorrência máxima. Porém o padrão permite concorrência ilimitada,
Go bloqueia submissão e não recebe contexto de cancelamento. O mecanismo é
referência parcial para sobrecarga; não equivale à fila durável do Bateia.

As avaliações foram declaradas após inspeção manual e continuam
DECLARED_NOT_REASSESSED no relatório. O código dos testes da fonte foi lido,
sem executá-los. O comando produziu mecanismo PARTIAL, hipótese INCONCLUSIVE,
alvo UNKNOWN e plano de baseline/aplicabilidade/experimento. Nenhum benefício
para o alvo, aprovação de adoção ou execução de GitHub Actions foi demonstrado.
As narrativas de estratégia são projeções genéricas do catálogo: os detalhes
específicos observados permanecem na revisão da fonte, não são inferidos por
um avaliador automático.

## Verificação

Suíte completa com paralelismo padrão e count=1 passou; pacotes corrigidos
passaram com count=3 e com detector de corridas. Lint cristalino, vet, build
e diff check passaram. Replay offline foi idêntico byte a byte; adulteração
de bytes foi recusada preservando saída; controle sem avaliações produziu
INCONCLUSIVE sem mecanismo. Ambiente de fixtures usou loopback e desativou
assinatura Git somente no processo. Gates selados com módulo separado não
foram reproduzidos. A aprovação vale para esta árvore e estas execuções.

## Próximo refinamento evidenciado

A infraestrutura online funciona, mas a avaliação semântica e o teste de
aplicabilidade ao alvo ainda são externos à jornada independente. O passo
seguinte deve preservar observações condicionais da fonte no relatório e
medir o runtime atual antes de propor mudança: o Bateia já possui fila
limitada, reserva e erros de sobrecarga. Não há justificativa neste ensaio
para substituir sua política por submissão bloqueante.
