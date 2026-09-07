# Passo 0041 — Tornar autofagia e Katamari planos de evolução verificáveis

Estado: IMPLEMENTADO — regressões e CLI verificadas.

## Intenção e contrato

Autofagia inventaria aproveitamento interno por capacidade, com escopo explícito e ação seguinte. Katamari compara oportunidades externas e propõe experimentos delimitados. Não instala, copia, remove ou executa código de candidatos.

- Agregação interna conserva conflitos e avaliações faltantes; resultados repetidos não criam implementações novas. Múltiplas implementações sustentam revisão de sobreposição, não consolidação automática. Evidência parcial origina investigação/completar internamente como hipótese, não lacuna comprovada.
- Status PROVEN de assessments externos exige fontes com identidade, caminho, regra, linha/trecho e hash; isso valida rastreabilidade declarada, não reexecuta a prova. Fonte de licença exige autoridade; todos os bloqueios de restrição, falha, manutenção e operação devem ser respeitados. Unknown nunca entra como cobertura.
- Propostas elegíveis ainda exigem validação de integração. Separar cobertura completa/parcial, investigação, nenhuma mudança e inconclusão; Feasible não será atalho para integração validada. Informar objetivo e universo da minimalidade. Retirar custo numérico fictício.
- Evidências parciais geram oportunidades com próxima verificação; comparar completar internamente, integrar e investigar substituição sem afirmar ganho não medido.
- Deduplicação, revisão única por projeto, limite enumerável e hashes determinísticos. Não confundir candidatos ausentes com nenhuma melhoria necessária.
- Saída versão 2; entrada legada permanece legível, sem confiança automática em listas de gaps nem em bools antigos. Sem verificador de integração, nenhum resultado será marcado como integração validada.

## Validação

Regressões públicas: conflito com boa cobertura, licença sem fonte, fonte de outro repo/hash inválido, manutenção desconhecida, resultado parcial não factível, no-change justificado, ausência opaca, fontes duplicadas, confronto entre revisões, enumeração limitada, permutações e ponte de evidência parcial para investigação. Preservar controles positivos com assessments fornecidos e rastreáveis. Exercitar CLI e pipeline sem rede; go test ./..., crystalline-lint e build do pacote wiring.

## Limites

Planos e verificações propostas não são estimativas de esforço nem integração executada. Prova comportamental geral e equivalência entre implementações permanecem trabalhos específicos. Não há alegação de auditoria independente.

## Resultado

Autofagia produz escopo e opções internas, conserva conflito/ausência de avaliação e substitui consolidação automática por revisão de sobreposição. Katamari aplica bloqueios e rastreabilidade, oferece oportunidades investigativas e distingue cobertura declarada de validação de integração. Saídas internas de inventário/proposta usam schema 2; envelope e entradas existentes permanecem compatíveis.

Testes completos, linteiro e compilação aprovados. Duas execuções CLI de autofagia e seis de Katamari reproduziram os controles. Veja [relatório e migração](../relatorios/autofagia-katamari-0041.md) e [artefatos](../benchmark/0041/README.md). Prova comportamental geral, cálculo real de esforço e execução de integração não foram implementados nem simulados como sucesso.
