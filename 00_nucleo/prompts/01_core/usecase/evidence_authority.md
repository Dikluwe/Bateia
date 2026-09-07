# Autoridade comum de evidência — passo 0043

## Contrato
Funções puras verificam ponteiros no escopo repo@revisão, caminho relativo limpo, SHA-256 (hex ou sha256:hex), tipo conhecido e regra identificada. Hash identifica conteúdo; não valida a afirmação sem leitura. Fontes de implementação exigem caminho de código e tipo de implementação. Documentação isolada sustenta somente declaração parcial, jamais operação comprovada.

SanitizeAssessmentEvidence copia avaliações e reduz PROVEN incompatível para PARTIAL quando há declaração válida ou INCONCLUSIVE quando não há testemunho admissível. Gates licença e manutenção exigem seus tipos próprios. Partes só conservam capacidades comprovadas com fonte de implementação correspondente à avaliação; nenhuma fonte é fabricada. Estados de conflito são conservados quando rastreáveis. Consumidores devem chamar a mesma fronteira.

Negação de operação local mantém conflito delimitado; negação de envio externo não comprova envio. Classificação lexical permanece limitada ao trecho, sem inferir comportamento geral.

## Vínculo limitado de pertinência — revisão de integração 0043
Para IDs pertencentes ao DefaultEvidenceRules, um PROVEN fornecido só conserva autoridade se a capacidade seleciona a mesma regra pelo procedimento findRule já utilizado pelo avaliador. Texto vazio não seleciona regra. IDs customizados continuam avaliações fornecidas com identidade declarada: a admissão não valida conteúdo nem inventa catálogo externo. Incompatibilidade conhecida reduz para INCONCLUSIVE com motivo explícito.

## Ausência fornecida — revisão final 0043
Status ABSENT fornecido ao consumidor só conserva ausência delimitada quando suas fontes são admissíveis. Sem fonte ou com fonte incompatível, reduzir para INCONCLUSIVE e registrar lacuna uma vez; chamadas repetidas da sanitização devem ser idempotentes. Esta admissão não converte silêncio documental em ausência geral.
