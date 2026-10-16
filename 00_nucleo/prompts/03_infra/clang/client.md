# 03_infra/clang/client

## Passo 0035 — Análise semântica C/C++

Cumprir `00_nucleo/passos/0035-corrigir-lint-e-analisar-c-cpp-com-clang.md`. Porta L1 SemanticProvider, execução Clang limitada em L3, injeção L4. Contexto explícito em compile_commands.arguments, sem executar command, plugins ou response files; capturar hashes de configuração, fonte e headers, argumentos, versão, símbolos e chamadas com intervalos. Falta de evidência nunca vira complete. Exposição L2 separa resolução estática e adequação à consulta.

## Passo 0038
Aplicar contrato do protótipo AST+CFG: adaptador opcional explicitamente configurado, hashes e limites, fluxo intraprocedural e estados não avaliados.

## Passo 0125 — Build irmão e contexto Meson

Buscar, nesta ordem, compile_commands.json na raiz, em build/ e em ../build/.
Somente ausência permite fallback; configuração presente ilegível, excessiva ou inválida
deve produzir diagnóstico, sem trocar silenciosamente de contexto. Limite do banco:
16 MiB, aplicado também à releitura de integridade. Resolver directory relativo ao
diretório do banco, file relativo a directory e selecionar somente fontes da raiz.
Preservar hash e caminho do banco nos inputs das unidades, incluindo falhas.

Aceitar command somente como texto de argumentos quando arguments estiver ausente.
Tokenizar espaços, aspas simples/duplas e escapes, sem shell, expansão, pipelines,
redirecionamento ou execução do compilador registrado. Recusar sintaxe de shell,
response files, plugins e opções fora da lista explícita. arguments tem precedência.
Aceitar flags concretas de contexto Meson/GCC compatíveis com Clang; descartar somente
-mtls-dialect=gnu2 e -flifetime-dse=1 (geração de código GCC, sem efeito no AST).
Não promover falha de Clang, orçamento excedido ou AST truncada a complete.
Garantir que a via io.Copy/ReadFrom usada por os/exec também respeite o limite
de saída, sem herdar bytes.Buffer.ReadFrom que contorne Write.
Quando nenhuma unidade conclui, o motivo deve declarar ausência de resolução,
separando captura de configuração/inputs de análise AST completa.
