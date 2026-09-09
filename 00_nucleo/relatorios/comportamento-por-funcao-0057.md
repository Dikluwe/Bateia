# Execução do passo 0057

## Entrega
Implementadas observações AST por função Go efetivamente recuperada. Cada registro depende de bytes extraídos com SHA256 correspondente, arquivo único, parse válido e declaração com nome/intervalo exatos. Referências ausentes, divergentes, linguagens não suportadas ou arquivos truncados com tamanho declarado maior ficam not_evaluated. Métodos usam a convenção existente de nome simples; arquivo/hash/intervalo distinguem declarações homônimas.

O parser registra atribuições, incrementos/decrementos e retornos, com expressão-alvo, trecho exato e linhas físicas. Não percorre corpos de closures aninhadas. A integração recompõe Behavior após validar os hits, descartando projeções antigas. JSON preserva o conjunto; Markdown mostra até 10 funções e 8 fatos por função, informa omissões e apresenta limitações junto aos fatos.

Analisador evidence-v19-function-behavior-origin e checkpoint implementation_search0057-r1 para local/live invalidam projeções antigas. Consulta, ranking, scores e admissão pública não foram alterados. Nenhuma decisão automática passa a considerar esses fatos como satisfação da intenção.

## Resultado local
O corpus 0053 congelado produziu 121 registros: 120 observados e um não avaliado por origem documental. Os observados cobrem 92 intervalos únicos e contêm 639 fatos sintáticos. Há repetição de fontes entre perguntas; essas contagens não são capacidades comprovadas nem funções pertinentes à pergunta. As três negativas não tiveram observações.

A recuperação permaneceu em 1/6 funções-alvo e 5/6 arquivos-alvo extraídos. Há mais informação sobre funções recuperadas, mas não ganho demonstrado na recuperação de alvos. Descoberta pública não foi repetida e está não avaliada nesta rodada. Nenhum código dos projetos foi executado.

Exemplo pertinente: LinearRgbToXyz, colors.go:446–451, SHA256 a75c3edb5dd7e957c49694a6ac3c8901ec7aba2fb5dededca552510eb533e37c. Foram observadas três atribuições a x/y/z nas linhas 447–449 e return na 450. O recibo contém as expressões completas. Isso descreve a sintaxe; não verifica correção matemática ou comportamento em execução. [Saída da análise](../benchmark/0057/collected-r1/case-05/cli.json).

[Métricas](../benchmark/0057/comparison.json), [recibos e hashes dos inputs](../benchmark/0057/collected-r1/receipts.json).

## Limites da conclusão
Fato observado: instrução identificada pelo parser nos bytes e intervalos citados. Não é prova de que o caminho executa, de mutação observável pelo chamador ou de equivalência funcional. Aliases, chamadas, efeitos indiretos e fluxos de controle permanecem não verificados. Uma função sem escrita direta não ganha classificação de imutável/pura; retorno de uma variável não prova que ela é nova ou independente.

As observações também existem para funções não pertinentes recuperadas pelo ranking. Não é correto apresentar639 fatos como melhora de qualidade da resposta à intenção. Este passo acrescenta material auditável para uma verificação posterior; a relação com os requisitos originais continua explícita como não verificada.

## Verificação e Tekt
A/B sem atestação física. Contrato e prompts antecederam código; candidato, apresentação e testes foram congelados antes dos gates. O autor A teve exposição incidental a trechos de testes preexistentes por erro em um glob de busca; não leu testes 0057/corpus/relatórios. Essa ocorrência impede alegar ausência total de leitura de testes. B não leu produção/patch/corpus. Não há selo de refinamento ou prova formal de independência.

Os controles independentes cobrem fatos e snippets exatos, identidade/hash/range, métodos homônimos, callbacks não percorridos, comentários enganosos, Go inválido, C++ não avaliado, arquivo ausente, recomputação sem fatos antigos, não mutação da entrada e limites Markdown/JSON. A revisão de B que permite limitações em texto combinado, sem exigir quatro elementos de slice, já está no hash efetivamente congelado e testado.

Corpus e todos os artefatos 0053 permanecem byte a byte intactos. Reavaliação offline reproduziu as métricas. Corpus conhecido, não holdout; tempos são telemetria do Bateia, não benchmarks dos projetos.

[Verificação final](../benchmark/0057/verification.json), [reprodução](../benchmark/0057/README.md). O próximo ganho de relevância ainda exige recuperar a implementação pertinente; a verificação de efeitos precisará de contrato específico de alias/fluxo/tipos ou execução controlada, sem inferir propriedades pela ausência de achados.

## Revisão por falha de regressão
A primeira suíte completa detectou que uma função encontrada apenas por documentação passou a mostrar fatos do corpo, contrariando o canal documental0051. O teste legado permaneceu intacto. O adendo0057-origem-documental passou a impedir esse caso: status not_evaluated, zero fatos. Um controle independente verifica que outro hit lexical válido da mesma função continua independente. Coleta e log r0 foram preservados; a rodada final local é collected-r1. Não houve rede. O contrato final inclui esta correção posterior à primeira execução, sem alegar que toda a revisão estava congelada antes das observações.

Verificação final: suíte Go completa, compilação, linteiro sem violações e nove controles do avaliador passaram. Hashes finais e integridade do corpus conferidos.
