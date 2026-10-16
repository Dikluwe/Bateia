# 0105 Investigar insights sem decisão de reúso

**Estado**: IMPLEMENTADO E VERIFICADO NO RECORTE OFFLINE  
**Data**: 2026-10-06  
**Relacionado**: ADR-0008, passos 0028/0103/0104

## Resultado e incerteza

Permitir perguntar como um mecanismo observado poderia melhorar um alvo, produzir uma hipótese e um plano de validação sem exigir IdeaSpec, política, decisão, run_id ou artefatos de recomendação. O recorte inicial é offline e recebe bundles e avaliações congeladas selecionadas pelo usuário. Não descobre fontes nem reavalia semântica das avaliações fornecidas.

## Contrato

Entrada schema_version 1: question, target (name, context, property, desired_effect, constraints), mechanism_classes, sources (bundle, assessments), max_hypotheses opcional. Campos textuais obrigatórios de até 1000 bytes, nomes/propriedades até 200; constraints até 20 itens de 200 bytes. Classes devem pertencer ao catálogo ativo de quatro mecanismos, sem repetição; fontes 1..8; itens e testemunhas até 512 agregados; avaliações até 128; pedido serializado até 2 MiB; max_hypotheses 1..12, default 5. Conteúdo de código pode ter quebras de linha; textos de controle não são admitidos nos campos da pergunta/alvo.

Commit de cada bundle deve ser SHA hexadecimal imutável; todos os itens precisam pertencer à identidade do bundle e conferir com seus hashes, ter paths válidos/sem ambiguidades. Não confiar em mecanismos prontos nem conceder status com base no texto da pergunta. Executar o extrator ativo, que vincula testemunhas aos bytes/localizadores e ao catálogo fechado. As avaliações continuam declaradas e são explicitamente identificadas como não refeitas.

O relatório insight-investigation-v1 conserva pergunta/alvo, input hash, hash do catálogo e proveniência das fontes. Seleciona as classes solicitadas e emite hipóteses de referência com TargetState UNKNOWN, estado no máximo PROPOSED (INCONCLUSIVE se fonte ou mecanismo limitado), licença UNRECOGNIZED e NOT_AN_ADOPTION_DECISION. Plano: medir baseline do alvo, conferir premissas/restrições e desenhar experimento limitado que possa rejeitar o efeito desejado. Não encontrar mecanismo compatível resulta em INCONCLUSIVE, não em ausência demonstrada.

## Camadas e superfície

L1 declara tipos, validação, extração e análise. L2 lê/projeta JSON limitado por `idea insights analyze --input arquivo [--output arquivo]`. L4 despacha antes de .env, vocabulário, cache e clientes. Reutilizar proteções de entrada/saída existentes. Nenhum adaptador novo de L3 é necessário para este recorte offline.

## Verificação e limites

Testar integridade, identidade, limites, classes, hipóteses conservadoras, determinismo, não mutação, licença sem gate de adoção, saída protegida e startup offline do binário. Compilar, vet, lint e suíte completa. Preservar toda a síntese durável/recomendação existente. Exemplos e medições ficam no Lab.

O marco não comprova benefício real para o alvo, autenticidade externa de fontes, cobertura de todo o ecossistema, busca automática ou ganho da IA. Avaliar uma fonte real de forma independente é o próximo marco, após verificar este contrato.


## Execução e preservação

A [síntese de execução](../laudos/0105-investigacao-independente-de-insights.md) registra os cenários e limites. A bancada ilustrativa fica em `lab/bancada/insights/0105`. Lint, vet, build, ensaio do binário e suíte completa com `go test -p 1 ./...` passaram. A execução concorrente teve duas falhas temporais em testes existentes, descritas no laudo; nenhum gate ou critério foi afrouxado. Utilidade para um alvo real e avaliações externas independentes permanecem pendentes.
