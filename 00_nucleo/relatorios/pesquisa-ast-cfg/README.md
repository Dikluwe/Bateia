# Pesquisa: AST, CFG e pré-processamento para o Bateia

Pesquisa de 2026-09-07 em fontes primárias LLVM/Clang. Código produtivo do Bateia não alterado. Documentação online consultada identifica-se como 24.0.0git; demonstração local executada com Clang 18.1.3. Exemplos C++ de LibTooling abaixo são referências de estudo, não componentes compilados nesta pesquisa.

## Exemplos de implementação consultados

| Referência oficial | O que estudar | Aplicação proposta |
| --- | --- | --- |
| [Tutorial LibTooling/ASTMatchers](https://clang.llvm.org/docs/LibASTMatchersTutorial.html) | Ferramenta loop-convert, callbacks e seleção de ForStmt | Capturar inicialização, condição e incremento do laço |
| [RecursiveASTVisitor](https://clang.llvm.org/docs/RAVFrontendAction.html) | ASTFrontendAction, ASTConsumer e visitante | Percorrer estrutura e conservar pai, função e localização |
| [AST Matchers](https://clang.llvm.org/docs/LibASTMatchers.html) | Predicados, binding e acesso aos nós C++ | Relacionar condições e chamadas por estrutura |
| [CFG: implementação](https://clang.llvm.org/doxygen/CFG_8cpp_source.html) | CFG::buildCFG e construção dos blocos | Caminhos, retornos antecipados e laços |
| [AnalysisDeclContext](https://clang.llvm.org/doxygen/AnalysisDeclContext_8cpp_source.html) | Construção e cache de CFG; opções de poda | Registrar opções da análise e evitar confundir arestas podadas com ausência de código |
| [pp-trace: implementação](https://clang.llvm.org/extra/doxygen/PPCallbacksTracker_8cpp_source.html) | Registro dos callbacks do pré-processador | Separar intervalos ignorados e expansões de macros |
| [PPCallbacks](https://clang.llvm.org/doxygen/classclang_1_1PPCallbacks.html) | SourceRangeSkipped | Registrar intervalos excluídos por diretivas condicionais |
| [Introdução a dataflow](https://clang.llvm.org/docs/DataFlowAnalysisIntro.html) | Propagação de fatos e condições de fluxo | Etapa posterior: fatos sobre valores e condições implicadas |

## Demonstração local

Exemplo próprio, criado para esta pesquisa:

```cpp
void draw_layer(int);
void compose();
void fallback();
void render(bool supported, int passes) {
  if (!supported) { fallback(); return; }
  for (int i = 0; i < passes; ++i) { draw_layer(i); }
  compose();
#if 0
  disabled_render();
#endif
}

```

Comandos efetivamente executados:

```sh
clang -std=c++17 -fsyntax-only -Xclang -ast-dump=json /tmp/bateia-ast-research.cpp
clang -cc1 -std=c++17 -analyze -analyzer-checker=debug.DumpCFG /tmp/bateia-ast-research.cpp
```

A segunda invocação produz um dump de depuração, não uma interface JSON de produção. Saída preservada em `example-cfg.txt`; contagens AST em `example-validation.json`.

Resultados observados: AST com IfStmt, ForStmt, ReturnStmt e três CallExpr; chamada desativada ausente. CFG: B7 bifurca para B6 (fallback e return) ou B5 (inicialização do laço); B4 testa o laço; B3 chama draw_layer; B2 incrementa e retorna a B4; B1 chama compose.

Interpretação deste exemplo: compose não está sintaticamente dentro do if, mas só pode ser alcançado pela alternativa que não retorna antecipadamente. Procurar um ancestral IfStmt da chamada seria insuficiente. Número real de passagens e conclusão das chamadas não foram verificados por execução.

## Recomendação para o Bateia

1. Estender o contrato existente com relações estruturais: condição, ramo, corpo de laço, chamada e intervalos. Diferenciar relação AST de dependência de controle; não tratar ancestralidade como prova de caminho.
2. Prototipar um adaptador C++ isolado de LibTooling que exporte AST + CFG + eventos de pré-processamento em JSON próprio e versionado. O núcleo Go continuaria puro e consumiria uma porta. Preservar argumentos efetivos, hashes, limites e contexto único da versão atual.
3. Usar SourceRangeSkipped para excluir trechos inativos da evidência positiva, inclusive dentro de funções ativas. Não inferir inatividade apenas da ausência de um nó AST.
4. Usar CFG intraprocedural para explicar retornos, ramos e laços; controle e alcançabilidade não comprovam comportamento funcional. Adiar análise interprocedural, ponteiros e despacho virtual até haver contrato e testes próprios.
5. Acrescentar dataflow apenas quando houver perguntas concretas sobre valores/condições. Não é necessário implementar um analisador simbólico completo para melhorar a próxima saída.

O binário auxiliar é uma proposta, não foi implementado ou instalado. A compatibilidade das APIs deve ser validada contra a versão LLVM escolhida. Não recomendo consumir o texto de debug.DumpCFG como contrato permanente.

## Critérios para o próximo passo

Retorno antecipado; if/else; for/while com break/continue; curto-circuito &&/||; #if 0 dentro de função ativa; macros com origem e expansão; múltiplas configurações; chamadas indiretas não resolvidas explicitamente desconhecidas. No Coin, provar relações estruturais entre modo, suporte, loop e composição; não prometer correção visual ou independência de ordem apenas pelo nome do mecanismo.
