# Passo 0038 — Protótipo Clang AST + CFG + pré-processamento

## Objetivo e contrato
Implementar adaptador opcional LibTooling, integrado ao provedor Clang, que exporta JSON versionado: funções e localizações AST, blocos CFG, sucessores, terminadores, chamadas diretas/não resolvidas e intervalos ignorados pelo pré-processador. Escopo intraprocedural no arquivo principal de cada unidade. Não interpretar arestas como prova de viabilidade, dominância ou execução. Opções de construção do CFG explícitas.

## Integração
Binário auxiliar explicitamente configurado por `BATEIA_CLANG_FLOW_TOOL`; sem descoberta em repositórios analisados. Invocação sem shell com flags já filtradas, limites do provedor, hashes das entradas conferidos após AST e CFG. Registrar hash do auxiliar e argumentos. Ausência/falha do auxiliar: fluxo não avaliado, sem destruir evidência AST válida. Núcleo Go puro; C++ restrito à infraestrutura.

## Evidência e apresentação
Intervalos ignorados não sustentam seleção lexical positiva na unidade cujo fluxo foi validado. Apresentar blocos/arestas/terminadores ligados às definições selecionadas, com contexto e limites; manter JSON integral. Não atribuir polaridade por posição do sucessor sem contrato próprio.

## Aceitação
Clang real: retorno antecipado, if/else, for/while, break/continue, curto-circuito, chamada indireta, #if 0 dentro de função ativa. Confirmar saída estruturada e localização; fluxo ausente é desconhecido. Validar Coin e isolamento da configuração. Executar lint, testes Go e build, além da compilação do auxiliar.

## Limites
Não há análise interprocedural, resolução de ponteiros, prova de independência de ordem, interpretação geral de algoritmos ou execução gráfica. Protótipo vinculado às bibliotecas LLVM 18 disponíveis. Sem certificação independente. Localizações registram expansão e origem (spelling) inicial; histórico completo de expansões aninhadas permanece fora deste incremento.

## Execução concluída
Implementado, integrado e executado no Coin. [Resultados e limites](../benchmark/0038/README.md): CFG de 71 funções no renderizador; suíte completa com auxiliar, teste real, lint e builds aprovados. Auxiliar opcional e explicitamente configurado; nenhuma prova de execução foi atribuída ao grafo.
