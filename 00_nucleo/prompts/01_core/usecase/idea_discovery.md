# Prompt Spec L0: Descoberta Rastreável de Ideias

## Objetivo

Executar cada expansão por uma porta observável, preservar tentativas e falhas e
deduplicar candidatos sem perder suas origens.

## Restrições

- L1 puro: sem rede, disco, relógio ou provedor concreto.
- Concorrência limitada e cancelamento respeitado.
- Ordem canônica segue o plano, não a conclusão das goroutines.
- Falhas parciais preservam resultados; falha total retorna erro.
- Identidade de repositório é `owner/repo`, sem diferença de caixa.

## Verificação

- Uma chamada por expansão.
- Origens mescladas na deduplicação.
- JSON idêntico para diferentes ordens de conclusão.
## Extensão operacional do Passo 0013

Execute expansões sob uma política limitada, propague attempt timeout ao provedor
e devolva snapshot parcial junto com erro de deadline/cancelamento. Preserve
candidatos e attempts concluídos, marque expansões não iniciadas, ordene a saída
pela ordem do plano e não deixe workers ativos após o retorno.

## Contrato 0043 — consultas ligadas à intenção

Antes da chamada ao provedor, cada expansão recebe RequirementIDs derivados dos Covers da IdeaSpec original, com IDs estáveis por conteúdo e tipo. Enviar os IDs na DiscoveryRequest e preservá-los na expansão do snapshot. Sem vínculo válido não fabricar associação: lista vazia exige investigação. Estes IDs descrevem a pergunta, não evidência do candidato.
