# Auxiliar opcional AST + CFG

Requer headers/bibliotecas Clang e LLVM 18. Exemplo para o ambiente validado:

```sh
clang++ -std=c++17 -I/usr/lib/llvm-18/include 03_infra/clang/flowtool/main.cpp /usr/lib/llvm-18/lib/libclang-cpp.so.18.1 -L/usr/lib/llvm-18/lib -lLLVM-18 -o /tmp/bateia-clang-flow
BATEIA_CLANG_FLOW_TOOL=/tmp/bateia-clang-flow ./bateia analyze -path /caminho/projeto -query 'consulta' -lang cpp
BATEIA_TEST_FLOW_TOOL=/tmp/bateia-clang-flow go test ./03_infra/clang -run 0038 -v
```

O auxiliar recebe arquivo principal seguido de flags já filtradas pelo adaptador Go. Não executa comandos de compile_commands.json. A saída `bateia-flow-v1` contém CFG intraprocedural, localizações AST e intervalos ignorados do arquivo principal. Opções e limitações são explícitas. O dump de depuração textual do Clang não é utilizado como contrato.
