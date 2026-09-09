# 0046 — Revalidar antes de projetar requisitos

Origem: passo0046 congelado. Antes de implementação.
ValidateImplementationRequirements(query,report,files) chama ValidateImplementationHits e substitui Requirements pela projeção parser.DescribeImplementationRequirements. Nunca manter Requirements antigas após perder fontes na releitura. Core puro via funções existentes, sem IO. Analyze e CollectEvidence usam este helper com a query original ou requisitos agregados, sem alterar avaliação de capacidades.

0048: depois de ValidateImplementationHits, substituir QueryInterpretation e Requirements usando query original. Permite relatório antigo/provedor não suportado com interpretação explícita sem fabricar fontes. Preservar status parcial/não avaliado e restrições not_verified.


## Passo 0057
Após validar hits e bytes, substituir Behavior via parser.InspectImplementationBehavior(out, files), descartando projeções antigas; preservar demais regras.


## Passo 0061
Depois da revalidação, recalcular ActionAssessments via parser.AssessImplementationActions(query,out,files); descartar projeções anteriores.
