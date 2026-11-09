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

## Passo 0068 — transporte e união

Aplicar `00_nucleo/passos/0068-compilar-capacidade-em-consultas-tecnicas.md`. Copiar
profundamente `SearchExpansion.Derivation` para `DiscoveryRequest.Derivation`, sem mutar
entrada. Executar todas as expansões técnicas independentes sob os budgets existentes;
resultado não vazio de uma não cancela as demais. Deduplicar candidatos como hoje e
conservar todas as expansion IDs, queries e origens técnicas canônicas. Não interpretar
discriminadores nem avaliar capacidade em L1.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.


## Integridade 0169

Antes de unir candidatos, copiar e ordenar respostas com ordem total: comparador de metadados existente e desempate pelo JSON completo. Escolher o menor representante por essa ordem, sem combinar campos. Preservar todos os matches como origens distintas, ordenadas/desduplicadas; match único e fallback sem match mantêm a projeção anterior. Não modificar slices da porta. Evidências de cada match têm ordem total por JSON. Não reescrever snapshots históricos.

Copiar também Attempts antes de preencher procedência regional; enriquecimento não modifica a resposta retida pela porta.
