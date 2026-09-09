# Testes independentes — avaliador sort-preserve-v1

Raiz causal: `00_nucleo/passos/0062-verificar-efeito-executado.md`, posterior ao gate discriminatório selado. Autor B sem leitura da implementação A. Materializar testes externos `usecase_test` em `01_core/usecase/effect_verification_test.go`.

Transcrever os 20 controles congelados do gate: 6 positivos, 10 violações semânticas e 4 opacos. Verificar status e não mutação de todas as observações. Preservação de valores antes/depois, ordenação não decrescente, multiplicidades exatas, extremos int64, nil/vazio equivalentes e prioridade de Completed=false sobre Failure. Razão deve explicar todo veredito, sem fixar vocabulário da implementação.

Acrescentar controle de alias: slices de entrada, estado posterior e retorno podem compartilhar armazenamento e avaliação não deve alterá-los. Exercitar exaustivamente saídas de tamanho 0..4 sobre {-1,0,1}, confrontando preservação do multiset com contagem independente e pares adjacentes; não ordenar resultados esperados nem usar algoritmo de produção. Não executar testes antes da integração pelo verificador.
