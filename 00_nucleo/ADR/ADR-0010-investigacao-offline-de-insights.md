# ADR-0010 Investigação offline de insights

**Status**: ACEITO PARA O RECORTE 0105  
**Data**: 2026-10-06  
**Relacionado**: ADR-0008, passo 0103 e [passo 0105](../passos/0105-investigar-insights-sem-decisao-de-reuso.md)

## Contexto

A análise pura de síntese foi separada de seu runtime, mas sua jornada pública ainda exigia artefatos e IDs de runs de recomendação. Isso restringia a investigação de mecanismos a uma saída de reúso. O extrator atual admite apenas quatro classes e usa avaliações fornecidas, cuja autenticidade e semântica não são automaticamente refeitas.

## Decisão

1. Acrescentar um pedido independente com pergunta, contexto do alvo, propriedade, efeito desejado e classes de mecanismo. Não exigir IdeaSpec, decisão/política de reúso ou run_id.
2. No primeiro recorte, receber fontes congeladas com bundles e avaliações declaradas. Verificar integridade dos bytes, identidade imutável e bindings pelo extrator ativo; não afirmar releitura remota ou avaliação independente.
3. Reutilizar o extrator e os tipos de hipóteses/planos existentes, com um caso de uso próprio de investigação. A síntese durável anterior mantém seus formatos e algoritmo.
4. Preservar alvo UNKNOWN, licença de adoção UNRECOGNIZED, transferência REFERENCE_ONLY e NOT_AN_ADOPTION_DECISION. Mecanismo observado segundo o catálogo não comprova efeito no alvo; avaliações limitadas ou fonte incompleta mantêm hipótese INCONCLUSIVE.
5. Formular plano de baseline, aplicabilidade e experimento, sem executar nenhuma dessas ações. A pergunta e o efeito desejado são intenção, não evidência.
6. L1 possui tipos e análise pura. L2 expõe `idea insights analyze`; L4 despacha antes de ambiente/cache/provedores. Este recorte não precisa de adaptador L3 ou runtime novo.

## Consequências e limites

Uma fonte que não passa por política de adoção pode orientar uma hipótese investigativa, sem receber autorização de uso de código. A proveniência declara avaliações não refeitas; hash é vínculo aos bytes e não autenticação da origem externa. O relatório pode ser útil como plano de investigação mesmo inconclusivo, mas a utilidade real exige avaliação própria.

Nenhuma busca automática, enriquecimento por IA, validação de alvo ou execução de candidato é incluída. A classificação em arquiteturas A1 continua em experimento; acrescentar essa capacidade não cria por si uma arquitetura independente.

## Critérios

Integridade/identidade, limites, não mutação e determinismo; fonte documental/incompleta/unbound e licença declarada; CLI estrita com proteção de entradas; binário em ambiente sem home/cache/runtime. Preservar contratos e capturas anteriores e executar a pilha Tekt completa.
