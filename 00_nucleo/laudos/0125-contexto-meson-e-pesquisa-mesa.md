# 0125 — Contexto Meson e pesquisa Mesa

Data: 2026-10-09. Implementação em L3 Clang, sem alteração de schema público,
ownership L1, critérios do catálogo ou fontes Mesa. O adaptador procura o banco
na raiz, build local e build irmão, nessa ordem; o primeiro existente governa.
Limite do banco e da releitura: 16 MiB. Directory relativo é resolvido pela
localização do banco. Command é tokenizado como argv literal e filtrado, sem
shell ou execução do compilador indicado. GCC flags concretas suportadas e duas
flags exclusivamente de geração de código descartadas são declaradas no prompt.

Testes reais verificam banco irmão de mais de 4 MiB, header gerado, macro com
espaços, caminho/hash, ausência de objeto e de execução do compilador candidato,
precedência e diagnóstico de configuração inválida/excessiva, rejeição de shell,
substituição, plugins e response files. Testes raiz passaram com loopback e
assinatura Git desativada só no processo. Lint, vet e build passaram. Cache e
binário usam /tmp por restrição de escrita. Gates selados não executados.

As duas consultas originais Mesa agora têm configuração rastreável: 12 unidades
selecionadas em cada, de 1.328 elegíveis, mas nenhuma AST concluída sob o limite
real de 32 MiB. Configuração e fechamento de inputs foram capturados. O ensaio
encontrou bypass de limite por bytes.Buffer.ReadFrom, corrigido com teste próprio.
Resultados anteriores completos foram preservados como superados. Complemento
com oito arquivos, via porta existente, selecionou cinco unidades clipping e onze
GS, também limitadas. Os achados são manuais, não resolução AST concluída. Pesquisa externa
fixou fontes SwiftShader e ANGLE: primeiro sem contexto Clang, segundo descartado
na triagem GN. Seus achados são manuais e possuem hashes/URLs próprios.

[Pesquisa e testes discriminantes](../relatorios/pesquisa-mesa-0125.md).
[Bancada, comandos e dados](../../lab/bancada/pesquisa/0125/README.md).
Não há novo draw, captura GPU, execução de fixture dos candidatos ou patch
Mesa validado. A utilidade demonstrada é recuperação de contexto e localização
estática; causas e intervenções no driver permanecem hipóteses.
