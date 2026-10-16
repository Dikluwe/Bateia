# 0125 — Contexto de compilação e investigação Mesa

Data: 2026-10-09. Capacidade: leitura semântica C/C++ no adaptador Clang (L3),
contrato proprietário em prompts/03_infra/clang/client.md. Consumidores: funil
local e relatórios CLI, sem mudança de schema ou critérios de evidência em L1.

Esperado: usar o banco do build irmão Mesa, preservar seu hash e fechar fontes
e headers antes de analisar AST. Suportar a representação command do Meson
como argumentos estáticos filtrados, sem executar shell ou compilador candidato.
Banco acima de 16 MiB e configurações inválidas permanecem não avaliados.

Refazer separadamente clipping/contorno e GS→sem GS. Ampliar a pesquisa para
projetos externos com commits e bytes rastreáveis. Resultado: arquivos, hipóteses,
controles discriminantes e limites, distinguindo coleta automática de interpretação
manual. Sem novos draws, instalação de driver ou alteração do Mesa privado.

Verificação: controles reais Clang, rejeições de contexto executável, lint,
testes raiz, vet e build. Resultados brutos pertencem a lab/bancada/pesquisa/0125.
