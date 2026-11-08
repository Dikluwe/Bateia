# 0153 — Ligar coleta à revisão de alternativas

## Contrato ex-ante

Corrigir as lacunas observadas no ensaio 0152 sem inserir produtos, algoritmos ou fontes esperadas. Capacidade: investigação; ownership L1 para coleta coordenada e contexto, L2 para flags/publicação. Novo modo opt-in `evidence collect --question` devolve bundles e exploração atualizada numa única saída. Sem question, preservar a lista histórica. Não chamar modelo nem executar buscas adicionais implicitamente.

O modo de revisão permite até três passagens não sobrepostas por arquivo, sob o mesmo teto global; primeiro conserva a representação dos arquivos escolhidos pelo seletor vigente. Tarefas pedem alternativas dentro da mesma fonte e uma revisão final de cobertura: opções sustentadas, rejeitadas e lacunas ainda pesquisáveis. Não afirmar cobertura completa ou novidade. `explore --review-alternatives` oferece o mesmo comportamento sobre bundles existentes; padrão preservado.

Validar entrada antes de chamar portas; recusar bundles adulterados; preservar cancelamento, limites e erros de resolução como falhas de bundles. Falha ao preparar contexto não publica resultado parcial disfarçado. Prompts e consumidores são atualizados antes do código; testes devem mostrar que a segunda passagem deixa de ser descartada sem perder a primeira, e que ausência de evidência permanece explícita.

## Novo ensaio

Agentes novos, mesma pergunta do 0152, sem histórico nem referências. Baseline anterior às ferramentas, aquisições e origem auditadas. Mesmos limites para ambos; revisão de cobertura exigida igualmente na entrega, sem indicar mecanismos. Bateia recebe somente documentação genérica da interface nova. Comparar decisões, delta do baseline, descoberta versus memória, uso real da correção e limites. Não reescrever o ensaio anterior ou supor vantagem.

Artefatos de execução, parecer e evidências ficam no Lab. Gate canônico: make verify.

## Resultado e evidência

Materialização e gate canônico concluídos. O replay manteve saída padrão idêntica e ampliou passagens no modo opt-in. Dois agentes novos executaram o reteste; um terceiro avaliou pacotes anônimos. Os dois braços cobriram as mesmas políticas operacionais e consideraram réplicas colaborativas. Não houve família exclusiva nova atribuível à pesquisa; descoberta nativa permaneceu limitada. Registro, parecer, auditorias e limitações na [bancada do reteste](../../lab/bancada/insights/problema-aberto-reteste-2026-11-08/README.md).
