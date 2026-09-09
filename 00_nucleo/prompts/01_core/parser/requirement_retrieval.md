# 0046 — Projeção lexical por requisito

Origem: passo0046 congelado. Antes de implementação.
DescribeImplementationRequirements é puro: linhastrim/ordenadas/deduplicadas da query; tokens do vocabulário lexical existente. Para cada requisito, somente hits Retrieval lexical com caminho relativo seguro, hash SHA256 hexadecimal64 e Line positiva fornecem fonte. Interseção direta dos termos do requisito com MatchedTerms lower-case do hit; não extrair correspondência de Reason, documento ou RelatedTo. direct_call sozinho não corresponde. MatchedTerms/MissingTerms ordenados; Sources deduplicadas/ordenadas path,line,symbol,hash. Sem fontes status not_verified; com fontes evidence_located explicitamente inferência lexical, não satisfação. Ausência/budget nunca prova inexistência. Sem termos úteis preserva not_verified. Input não mutado. Projeção deve receber somente hits revalidados ao ser usada em produção.

0048: para cada linha usar InterpretImplementationQuery, termos positivos interpretados e Constraints originais not_verified. Marcador natural não aparece como correspondência ou termo faltante: é restrição não verificada separada. Sempre dizer que evidência lexical não verifica cumprimento nem escopo da restrição, mesmo sem MissingTerms. Não modificar status para PROVEN/ABSENT. Fontes seguem validação anterior.

0051: projetar hits já revalidados lexical e documentation separadamente em Sources e DocumentationSources. Documentação só para interpretação natural/mixed e com intervalo Doc positivo anterior à declaração e EndLine válido. MatchedTerms/MissingTerms unem os canais. evidence_located prevalece se há código, documentation_located quando só há Doc, not_verified quando nenhum. Ordenar/deduplicar fontes por canal; Reason distingue indício documental de código e preserva comportamento não verificado.

0052 r2: antes de associar cada hit a uma linha, reaplicar admissão por ação do
adendo aos MatchedTerms normalizados desse hit. Linha inadmissível não recebe
fonte nem termos por interseções secundárias de outra linha. Preservar consulta
símbolo, canais isolados e estados anteriores; Reason explicita filtro quando
aplicável. Não extrair evidência de Reason ou de outros hits para completar ação.
