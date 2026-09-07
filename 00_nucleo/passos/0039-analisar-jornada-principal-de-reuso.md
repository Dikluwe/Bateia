# Passo 0039 — Analisar a jornada principal de reúso

Estado: CONCLUÍDO — achados de qualidade reproduzidos. Natureza: diagnóstico de produto, sem alterar lógica de produção e sem alegar auditoria independente.

## Incerteza
O fluxo IdeaSpec → descoberta → evidência → decisão → dossiê preserva a intenção e ajuda a decidir o que reutilizar em casos diferentes? Quais limitações se originam no catálogo de regras, na coleta, na avaliação ou na apresentação?

## Escopo e casos declarados antes das execuções
1. Ideia pública do próprio Bateia (0012), executada no checkout atual pelo CLI local.
2. Ideia de exportação de notas Markdown, com implementação Python pequena e inspecionável: observar cobertura do catálogo fora do domínio predominante.
3. Processamento assíncrono de eventos, com candidato Go lexicalmente atraente que apenas devolve a entrada: verificar se nomes e testes triviais viram capacidade comprovada.
4. Variante com documentação que nega um não objetivo: verificar tratamento de negação e contraevidência.
5. Descoberta vazia e falha de provedor em RecommendationService e MCP, por adaptador controlado: não produzir BUILD ou sucesso indevido.
6. Replay de um dossiê local: verificar preservação de decisão e artefatos causais, sem exigir igualdade de telemetria/tempos.

Caso sintético não mede recall real nem qualidade geral; caso local não executa descoberta pública. Usar somente material público/documentado e novos exemplos abertos; não ler conjuntos TEST ocultos, oráculos protegidos ou alterar resultados históricos.

## Observáveis
Input original, expansões, candidatos, planos/arquivos coletados, estados e fontes por capacidade, restrições/não objetivos, decisão, alternativas e ação seguinte. Registrar comandos, retornos, hashes e limites. Diferenciar testes de interface, análise de código e execuções completas.

## Entrega
Relatório com achados priorizados, evidência reproduzível e recomendação do próximo incremento. Preservar código de produção durante a análise. Não atribuir ganho de produto a quantidade de símbolos ou ao caso Coin. Coin não integra esta amostra.

## Resultado da execução

Cinco recomendações locais, um replay e dois cenários controlados de descoberta pelo serviço e MCP executados. Testes existentes de usecase, CLI e MCP aprovados. O diagnóstico reproduziu comprovação indevida por nomes/testes, fonte de licença sem autoridade, metadados locais derivados da consulta, cobertura insuficiente e falhas de negação/projeção. As proteções de descoberta inconclusiva e replay foram preservadas nos casos examinados.

Veja o [relatório completo](../relatorios/analise-jornada-principal-0039.md) e os [artefatos reproduzíveis](../benchmark/0039/README.md). Nenhuma correção de produção foi aplicada neste passo; a próxima prioridade recomendada é o contrato de comprovação e autoridade da evidência na jornada principal.
