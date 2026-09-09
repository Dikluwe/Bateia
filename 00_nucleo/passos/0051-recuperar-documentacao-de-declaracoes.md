# Passo 0051 — Recuperar documentação vinculada a declarações

2026-09-08. O diagnóstico 0049 mostrou palavras funcionais na documentação de funções que o índice ignora. Este passo acrescenta um canal documental rastreável, preservando prioridade e pontuação da recuperação de código. Não amplia glossário para memorizar os exemplos.

## Contrato congelado

RankGoImplementations deve recuperar termos de comentários Doc vinculados pelo AST a funções/métodos Go. Comentários de pacote, soltos, internos ao corpo, literais, tipos e arquivos gerados não originam este canal. Usar implementação lexical existente para tokenizar documentação e InterpretImplementationQuery para termos; consultas somente de símbolo não consultam documentação. Consultas naturais/mistas podem localizar documentação, sem alegar execução ou cumprimento, inclusive diante de uma frase documental negativa.

Preservar integralmente a ordenação/pesos dos hits de código existentes e a expansão de chamadas. Acrescentar candidatos documentais somente nos slots restantes após hits de código/expansão, sem duplicar a mesma declaração. Score documental = número de termos distintos correspondentes, limitado a 99; desempate por caminho/linha/símbolo. Somente comentários da declaração lida integralmente são elegíveis. Não computar novos campos no IDF do código. Hits documentais usam Retrieval=documentation, Kind=definition, Line/EndLine da declaração, SHA256 do arquivo, MatchedTerms e Reason explicitando indício documental. Acrescentar DocumentationLine e DocumentationEndLine (JSON documentation_line/documentation_end_line, omitempty), referentes ao grupo Doc anexado. Nenhum TermWeights de código inventado. Limite que omite candidatos documentais marca Limited/partial.

Revalidar hits documentais depois da extração: hash, declaração com nome/intervalo correspondente e grupo Doc AST exatamente no intervalo informado. Intervalos positivos arbitrários não bastam. Fonte documental não se transforma em fonte lexical de implementação. RequirementRetrieval ganha DocumentationSources []ImplementationReference (documentation_sources,omitempty). Sources continua código lexical. Status documentation_located quando só há fontes documentais; evidence_located quando há código, mesmo se também houver documentação; not_verified quando nenhum canal tem fonte. MatchedTerms/MissingTerms consideram ambos canais, com Reason explicitando a distinção e o comportamento não verificado. Não alterar avaliação de capacidades nem promovê-las.

Markdown mostra indicação documental e trecho de comentário (no máximo 20 linhas) para hits documentais; não exibir corpo como se fosse evidência da frase. Projeção por requisito separa fontes documentais e de código. JSON preserva todos intervalos/hashes. Atualizar versão de evidência e checkpoint local/live; canal funciona local/remoto via parser compartilhado.

## Aceitação

Autor B congela novos exemplos independentes antes de executar. Positivo: pergunta funcional natural localiza função de nome/corpo neutros por Doc; negativo: mesmo texto solto, em literal, corpo ou arquivo gerado não localiza; documentação negativa não promove capacidade. Símbolo exato preservado, código anterior não deslocado, limite e duplicatas explícitos. Releitura alterada e DocLine forjada removem fonte. Verificar campos na projeção e Markdown. Ensaios do integrador em fonte real antes/depois demonstram localização documental antes ausente, sem inferir ganho PT/EN universal. Suíte completa, lint e build.

## Regime e limites

Tekt A/B sem atestação física. Contrato/integrador, implementação A e testes B separados, filesystem compartilhado. Sem selo de refinamento. Documentação é declaração do autor, não prova; orçamento cheio pode impedir este canal. Tradução, análise de ação/objeto/condições, seleção conjunta e cobertura do catálogo continuam pendentes. Critério de ganho é recuperar fontes documentais vinculadas antes invisíveis, não simplesmente aumentar número de hits.
