# Um agente usando somente XERJ consegue chegar às ideias?

**Sim, neste ensaio chegou a duas direções semelhantes e trouxe uma terceira complementar.** O XERJ fez recuperação de trechos; um agente independente interpretou os resultados. A nota G7 herdada após re-pin não foi identificada espontaneamente por ele.

## Protocolo

O agente começou sem histórico do chat e sem acesso às propostas anteriores. Recebeu a mesma pergunta orientada da exploração do Bateia, as mesmas quatro abordagens familiares e acesso ao mesmo corpus congelado de 19 arquivos públicos. Só podia consultar o wrapper `search.py`, que chama `xerj search` no nó local. Não podia ler diretamente os arquivos, buscar na web ou consultar relatórios do Bateia.

Pergunta fornecida:

> How can XERJ corpus evidence provenance verification distinguish retrieval relevance from tested behavior, expose applicability conditions and support reproducible counterexamples without confusing signatures or feedback with correctness?

Essa pergunta já sugere os temas de aplicabilidade, verificação e contraexemplos. O ensaio mede se o agente consegue desenvolver propostas fundamentadas a partir desses temas, não descoberta espontânea de temas nem vantagem causal de uma ferramenta. O avaliador conhece as propostas anteriores; não houve avaliação cega de qualidade.

Usamos a release oficial rc.101, nó apenas em loopback, diretório de dados temporário e feedback desativado. O indexador concluiu em 2,2 segundos, registrando 18 arquivos e 415 registros pesquisáveis a partir da pasta com 19 arquivos. O inventário de entrada está em [corpus.json](corpus.json), e a contagem efetiva do indexador em [index.log](index.log); não presumimos igualdade entre inventário e registros indexados.

O agente realizou sete buscas lexicais BM25. Consultas, comandos e respostas integrais ficam em [queries.jsonl](queries.jsonl). Todas terminaram com código 0. O relatório independente foi preservado em [independent-result.md](independent-result.md), sem reescrita pelo avaliador. O servidor foi encerrado ao concluir.

## Comparação com as propostas anteriores

| Proposta do ensaio Bateia + IA | Resultado do agente + XERJ |
| --- | --- |
| Vincular a nota G7 ao pin/protocolo | Não encontrou o problema específico de herança. Falou de vínculos por pack/revisão para replay, uma direção relacionada, mas insuficiente para contar como o mesmo achado. |
| Recibo do que foi efetivamente verificado | Propôs separar integridade, autenticidade, relevância, execução da alegação e feedback em recibo por alegação. Coincidência substancial. |
| Condições de aplicação e contraexemplos junto à referência | Propôs condições verificáveis de entidade/operação/versão/ambiente e uma exportação separada para replay de contraexemplo mínimo. Coincidência de direção, com desenvolvimento diferente. |

O agente também encontrou na documentação Query-1 um contraexemplo de contagem de pessoas/dispositivos e usou limites de auditabilidade das contagens para desenvolver sua proposta de replay. O próprio documento já reconhece esse problema e anuncia trabalho para rc.102; o agente não o apresentou como descoberta nova.

O resultado reduz a justificativa para apresentar aquelas duas direções como uma exclusividade do Bateia. A recuperação do XERJ foi suficiente para que um agente chegasse a propostas comparáveis sobre o mesmo material. O primeiro achado exigiu a investigação adicional e a reprodução descritas na [bancada G7](../g7-reproduction/README.md).

## O que este teste não estabelece

É uma execução de um agente e uma tarefa, sem repetição estatística ou medida comparável de tempo/tokens de raciocínio. O corpus foi escolhido antes; o Bateia selecionou 16 cartões e limitou quatro arquivos ao prefixo, enquanto o XERJ indexou o conteúdo fornecido com seu próprio recorte. Não há uma conclusão geral sobre superioridade, novidade ou produtividade de qualquer ferramenta.

O XERJ não gerou as propostas por si só neste teste. Seu Query-1 responde consultas estruturadas sobre fatos; a interpretação aberta de engenharia foi feita pelo agente. As propostas não comprovam que os recursos já existam no produto ou que tragam ganho ao implementá-los.

Para reproduzir a recuperação, reconstrua a pasta do corpus a partir de `../exploration-bundles.json`, preservando os caminhos registrados em `corpus.json`. Inicie um nó local isolado com feedback desativado, indexe com `--prefix agenttest --workers 2` e repita as consultas do transcript com o mesmo binário. Os arquivos arquivados em `search.py` registram os caminhos temporários usados no ensaio. Repetir as buscas não garante que outro agente produza o mesmo texto.
