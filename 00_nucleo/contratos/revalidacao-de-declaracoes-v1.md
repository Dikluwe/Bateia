# Revalidação de declarações v1

Contrato ativo do [passo0161](../passos/0161-corrigir-transicoes-historico-e-revalidacao.md), especializando a releitura0044–0051 sem reescrever seus registros. L1 implementation_retrieval valida os bytes e a identidade AST de hits antes da projeção de requisitos; efeitos continuam nas portas L3.

Depois de hash e limites numéricos, conferir símbolo e início/fim exatos da declaração Go. O conjunto coincide com os nós do índice atual: FuncDecl, TypeSpec e nomes de ValueSpec globais. Métodos conservam nome simples e faixa própria. Campos Kind/Retrieval ausentes legados não dispensam a identidade da declaração. Arquivo não Go, gerado ou com parse inválido não confirma declaração. Declarações são indexadas uma vez por arquivo na releitura. O canal documental continua exigindo seu grupo Doc, e referências de chamada conservam as validações anteriores.

Declaração falsa é removida com status partial e diagnóstico; Requirements não a reapresenta como fonte localizada. Preservar Examined/procedência do índice e entrada do chamador. Não recalcular pontuação lexical nem promover localização a prova de requisito, comportamento ou efeito. Não modificar resultados históricos.

Critérios DECL-01: símbolo ou faixa adulterados com bytes iguais são recusados; DECL-02: controles reais de funções/métodos/tipos/valores/legado permanecem; DECL-03: hash, linguagem, parse e arquivo gerado são recusados. Consumidor: [regressões de identidade](../../01_core/usecase/implementation_identity_regression_test.go).
