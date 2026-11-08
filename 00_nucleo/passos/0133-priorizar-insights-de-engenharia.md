# 0133 — Priorizar insights de engenharia

**Data:** 2026-10-28. **Estado:** FOCO DE PRODUTO REGISTRADO; evolução funcional ainda depende de contrato e materialização próprios.

## Intenção do usuário

Priorizar insights de engenharia derivados da investigação de software. A pergunta orientadora é: o que podemos aprender com outras implementações para conceber ou melhorar uma solução? O usuário reafirmou esse propósito após a comparação com outros projetos. Esta decisão não transforma o Bateia em auditor de cópias ou de infrações de propriedade intelectual.

Reúso continua como capacidade existente e possível consequência da investigação. A prioridade não revoga suas políticas, exclui comandos ou muda a força da evidência. Uma fonte pode ensinar uma estratégia sem ser elegível para adoção de código.

## Entrega de produto desejada

Partir de pergunta/problema e, quando disponível, contexto de um alvo. Localizar abordagens pertinentes, examinar fontes e produzir aprendizados com:

- estratégia concreta e observações localizadas em fonte/versionamento;
- condições necessárias, limites e comportamento em situações adversas;
- trade-offs e diferenças entre abordagens quando houver evidência suficiente;
- relação com a pergunta e com as restrições do alvo, distinguindo inferência de observação;
- hipótese de melhoria e próximo experimento que possa rejeitá-la.

Uma conclusão de não transferibilidade também pode ser um insight útil. Não inventar comparações, efeitos ou novidade quando a evidência não os sustenta. A falta de decisão de adoção ou de benefício medido não elimina o valor de uma observação delimitada.

IA pode propor consultas, interpretações e hipóteses; cada observação precisa de fonte conferível e escopo. Interpretação gerada não autentica a origem nem comprova efeito. Evitar confundir descrição documental, observação estática e execução. O critério é o aprendizado útil com suas condições, não o número de fontes, símbolos ou relatórios.

## Estado existente e lacuna prioritária

O ADR-0010 e o passo 0105 já fornecem `idea insights analyze`, com entrada independente, fontes congeladas, catálogo fechado e avaliações fornecidas. A implementação valida vínculos/bytes; não refaz a avaliação semântica, busca automaticamente o ecossistema ou valida o alvo.

O ensaio 0106 coletou e revisou manualmente uma fonte real. Observações como limite de workers, submissão bloqueante e ausência de contexto na assinatura ficaram em `source-review.json`. A estratégia projetada no relatório permanece genérica do catálogo. Essa diferença está explicitamente documentada na bancada, sem execução dos testes do projeto externo ou benefício medido no alvo.

O próximo corte funcional prioritário deve carregar observações específicas e condicionais da fonte até o relatório, com proveniência e limites, em vez de tratar o nome de uma classe como descrição suficiente da implementação. Esse corte exige prompt proprietário e testes antes de código. Não é materializado por este passo documental. Depois, avaliar comparação entre fontes, pertinência à pergunta e assistência de IA sob contratos próprios.

## Ownership e preservação

L1 possui critérios de observação, hipóteses e análise em `insight_investigation.go`/`improvement_analysis.go`; L2 possui entrada e apresentação; L3 possui busca/coleta; L4 compõe. Tekt e os ADRs existentes permanecem governantes. Não criar novo runtime ou ampliar o catálogo por mera mudança de foco.

Neste passo atualizar apenas a orientação pública, a visão de produto e a navegação. Preservar todos os arquivos Go, schemas, hashes históricos, avaliações, testes e políticas. Registrar baseline e verificação em `lab/bancada/insights/0133`, sem promover o ensaio histórico ou atribuir utilidade já demonstrada ao novo foco.

## Critério para a primeira investigação

Escolher uma pergunta concreta e um alvo ou contexto; registrar fontes e expectativas antes da avaliação. Julgar se os aprendizados são específicos, verificáveis e úteis para orientar uma escolha ou experimento. Medir tempo e trabalho restante; comparar com investigação assistida por ferramentas existentes antes de afirmar ganho. A escolha da primeira pergunta não está fixada por este documento.

## Primeira investigação escolhida pelo usuário

O usuário escolheu melhorar o próprio Bateia usando outros projetos como referência. Antes do novo ensaio, fixar a pergunta: o relatório de insights atual preserva as condições e diferenças relevantes de uma implementação real para orientar uma melhoria do Bateia?

Usar o corpus já conhecido 0106 (`sourcegraph/conc` no commit registrado), sem alegar holdout. Conferir hashes e localizadores dos bytes congelados; reproduzir o comando atual e comparar as descrições produzidas com a revisão manual. Examinar a fila atual do Bateia e seus testes de capacidade/cancelamento para evitar recomendar algo já existente. Distinguir nome de classe, estratégia específica, pré-condições operacionais e estado observado do alvo.

Critérios antes do ensaio: preservar replay/integridade; não executar o código da fonte externa; não transformar avaliações manuais em avaliação independente; registrar quais condicionantes estão apenas na revisão manual e quais chegam ao relatório. Formular aprendizados e um próximo contrato funcional com limites, sem mudar código, catálogo, estados ou recibos anteriores. Este ensaio é diagnóstico de um caso conhecido, não medida geral de precisão nem demonstração de benefício no alvo.

## Resultado registrado

Consultar o [laudo 0133](../laudos/0133-foco-em-insights-e-primeira-investigacao.md) e a [bancada](../../lab/bancada/insights/0133/README.md). O foco está registrado e a primeira investigação do próprio Bateia foi realizada. As melhorias funcionais derivadas permanecem propostas, sem alteração do código Go ou promoção de hipóteses.
