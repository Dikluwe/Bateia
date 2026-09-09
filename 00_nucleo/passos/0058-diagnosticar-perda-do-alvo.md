# Passo 0058 — Diagnosticar a perda de uma implementação conhecida

## Experimento anterior à correção
Usar a pergunta congelada bitset-p1 e snapshot0053. Comparar recuperação sem corte de24: pergunta original, termos da hipótese produzida pelo planejador existente e símbolo conhecido. O símbolo é controle dirigido, não descoberta nem ganho de recall. Registrar ranking, termos, campos, hashes e limite. Não alterar corpus, alvos, catálogo ou métricas. Escrever contrato da correção somente após localizar a perda.

## Regime
Tekt A/B para a correção: implementação e testes independentes a partir do contrato congelado, sem atestação física. Diagnóstico é do integrador e inclui alvo conhecido; não holdout. Nenhuma busca pública necessária se o defeito estiver no fluxo local. Não adicionar novas funcionalidades apenas para descrever o problema.

## Diagnóstico observado
No snapshot bitset, pergunta original produz55hits sem Intersection (limite10000); termos do plano existente colocam Intersection em7º e símbolo explícito em1º. Os bytes/AST permitem localizá-la. O ponto de perda é a representação da consulta local, que não recebe a hipótese de ação/objeto existente. Não é um corte no top24. diagnostic-before.json preserva comparações; consulta pelo nome é controle dirigido.

## Contrato da correção, congelado antes do candidato
- Em InterpretImplementationQuery, depois de determinar Mode, para cada linha Mode=natural, consultar PlanDiscoveryQuery da mesma linha original. Somente status hypothesis autoriza adicionar seus ActionTerms/ObjectTerms a Terms, com deduplicação/ordem determinística. Não alterar Original, Symbols, Constraints, nem a consulta enviada ao provedor. Mode mixed/symbol não recebe essa nova expansão.
- Registrar duas ImplementationQueryExpansion (ação e objeto), Source=linha original completa, Kind=translation, Status=hypothesis, RuleID=`0058-v1:` + plan.Version + `:` + respectivo RuleID do plano (action.* e object.*). Terms contém exatamente os termos daquele grupo do plano. Não inventar sinônimos, nomes de funções, repositórios ou novos gatilhos. Não expandir recursivamente. Planos literal_only (incluindo negação/ambiguidade/sintaxe) não recebem expansão0058. Expansões0052 existentes continuam intactas.
- Nenhuma mudança de pesos/score/admissão/limites/AST/extração ou seleção documental. Novos termos são hipóteses lexicais, nunca satisfação do requisito. Expansões já são exibidas no relatório pelo contrato existente.
- EvidenceVersion evidence-v20-local-query-plan e checkpoint implementation_search0058 local/live; manter descoberta/checkpointdiscovery_query.
- A implementa somente query_interpretation.go, evidence.go, checkpointCLI e respectivos prompts. B escreve controles independentes de integração do plano (planos já existentes são dependência estável), positivos fora da pergunta de diagnóstico, limites negados/ambíguos/símbolos, provenance/determinismo e não mutação. Gate de ranking com fixtures sem nomes dos projetos do corpus verifica que a operação nomeada passa a ser recuperada a partir da ação+objeto. Não alterar expectativas legadas para fabricar ganho.
- Repetir diagnóstico e corpus0053 local sem trocar perguntas/alvos. Mostrar resultado por caso, não só média. Esses casos são calibração, não holdout. Não executar rede nem código dos projetos. Se não houver ganho no alvo diagnosticado, não acrescentar outro mecanismo: revisar o diagnóstico.
