# Pesquisa Mesa: clipping e transição GS

Data: 2026-10-09. Duas perguntas independentes, sem novos draws na GPU.
Resultado: fila de investigação; nenhum patch Mesa foi validado.

O Bateia agora lê o banco Meson do build irmão. As consultas originais analisaram
12 unidades cada de 1.328 elegíveis. O banco tem 6.069.324 bytes e SHA256
`5ed538cf999e09364a4aa0221ec6a9d9063b75b6a373b8cd13cded59e0f14725`.
Clang, argumentos, caminhos e hashes estão nos relatórios. A análise AST final
é NOT_EVALUATED: todas as unidades selecionadas excederam o limite real de
32 MiB de saída. O motivo já não é ausência de configuração: o banco e o
fechamento de fontes/headers foram efetivamente encontrados e capturados. A seleção
lexical também trouxe unidades pouco pertinentes; não prova adequação. Rodada
complementar usa Paths explícitos na porta SemanticProvider para oito arquivos
centrais: cinco unidades na consulta clipping e onze na GS (nove comandos para
si_state_draw.cpp, mais dois arquivos). Também excederam o limite de saída.
Paths preserva variantes de compilação; não significa um comando por arquivo.

A conferência final verificou 762 inputs únicos, incluindo banco, fontes e
headers: todos os hashes capturados coincidem com os bytes após as consultas.
Os três arquivos de fontes congeladas também conferem com seus manifests.
Isso demonstra integridade neste recorte, não resolução AST ou efeito GPU.

[Bancada e reprodução](../../lab/bancada/pesquisa/0125/README.md): resultados
completos comprimidos, projeções legíveis e fontes congeladas com manifests.
As hipóteses abaixo são interpretação manual dos bytes congelados; o catálogo de
insights não ganhou classes de clipping ou estado GPU.

## Clipping: arquivos, hipóteses e controles

| ID / Mesa 25.2.8 | Evidência estática | Hipótese | Teste discriminante sem GPU, proposto |
| --- | --- | --- | --- |
| C1: `draw_pipe_clip.c:404,480`; `draw_pipe_unfilled.c:127,157`; `draw_pipe_validate.c:217,261` | O clipper conserva flags por primitiva e vértice, marca verdadeira a aresta nova do plano do usuário e precede unfilled. LINE/POINT exigem ambas as flags. | A trilha acelerada perde criação ou transporte das flags do contorno. | Fixture CPU do pipeline draw: quad atravessando plano do usuário, modos LINE/POINT, flags externas verdadeiras. Capturar vértices/flags entre clip e unfilled. Controles: pré-recortado e flags falsas. Extremos/aresta nova devem aparecer quando elegíveis. Distingue geração de geometria de seleção de arestas. |
| C2: `draw_pipe_clip.c:488,493` | O fonte distingue plano do usuário de frustum: para frustum a nova edgeflag é falsa, com comentário explícito. | Generalizar o oracle do plano do usuário aos seis planos do frustum confundiria escopos deste caminho. | Trocar somente a classe do plano na fixture C1 e confrontar as flags com o contrato OpenGL aplicável. O contraste restringe o escopo; não identifica o bloco de hardware defeituoso. |
| C3: `st_atom_rasterizer.c:138`; `radeonsi/si_state.c:1083` | State tracker traduz FrontMode/BackMode; driver codifica POLY_MODE/PTYPE. Configuração não demonstra criação de interseções. | O modo pode estar correto enquanto a geometria/flags após recorte são insuficientes. | Fixture CPU das traduções e builder de rasterizer: mesmos modos, face e flags para original/pré-recortado, comparando bits com a saída de C1. Igualdade exclui tradução simples nesse recorte, sem comprovar efeito do hardware. |

C2 é uma restrição nova da referência de software. O oracle desta campanha
cobre plano fixo do usuário e gl_ClipDistance. Correção geral exige controles
próprios por classe de plano, atributos e identidade das arestas externas.

## GS → sem GS: arquivos, hipóteses e controles

| ID / Mesa 25.2.8 | Evidência estática e referência anterior | Hipótese | Teste discriminante sem GPU, proposto |
| --- | --- | --- | --- |
| G1: `si_state_shaders.cpp:3800,3824`; `si_descriptors.c:2059` | O fonte contém a chamada de bind_gs a notify. Notify muda bases de userdata e as_ls/as_es/as_ngg. Logs anteriores confirmam NULL→GS→NULL. | Estado derivado pode divergir mesmo com vínculo restaurado. | Fixture CPU: comparar contexto sem GS desde início com GS→NULL; observar chave VS, bases userdata, dirty bits e variante corrente. Divergência localiza restauração; igualdade desloca foco à emissão. Verificar só GS NULL não distingue a hipótese. |
| G2: `si_state_shaders.cpp:1689,1768`; `si_state_draw.cpp:43` | VGT_GS_MODE é emitido pelo estado VS. O comentário exige troca de VS/copy shader entre pipelines, enquanto GS pode não ser reenviado. | Variante VS ou cache PM4 incorreto conserva configuração do GS no draw seguinte. | Command sink CPU no emissor ou dumps existentes: comparar programa VS, VGT_GS_MODE, VGT_SHADER_STAGES_EN e userdata. Aceite: estado coerente com sem GS, observando exceção PrimID/GS_SCENARIO_A. Bytes/variante corretos enfraquecem essa hipótese. Não submeter o draw à GPU. |
| G3: `si_state_draw.cpp:59`; `docs/relnotes/25.2.6.rst:164` | A propagação VS→GS combinado em GFX9+ já está em 25.2.8 e foi incluída em 25.2.6. | Não há suporte para tratá-la como correção posterior ausente do baseline. | Fixture CPU da seleção: VS dirty com HAS_GS ligado/desligado; conferir previous_stage combinado e seleção normal ao sair. Título de release não distingue essas situações. |

As [notas do Mesa 25.2.6](https://docs.mesa3d.org/relnotes/25.2.6.html)
confirmam o backport. “call si_shader_change_notify when vs bind” aparece nas
[notas do 26.0](https://docs.mesa3d.org/relnotes/26.0.0.html) no conjunto de
mudanças de task/mesh. Seu diff completo não foi auditado aqui; permanece pista
secundária, sem vínculo comprovado com Renoir ou indicação de cherry-pick.

Oracle e logs anteriores foram congelados em
[oracle-manifest.json](../../lab/bancada/pesquisa/0125/oracle-manifest.json):
falha do contorno no radeonsi, passagem no llvmpipe, draw suplementar isolado
funcionando, reset na sequência e VGT_FLUSH insuficiente. Foram copiados sem
reexecução. Não há dump PM4 completo nos arquivos selecionados; G2 permanece
proposta para fixture CPU ou análise de captura já disponível.

## Outros projetos

| Projeto / commit / arquivos | Evidência | Hipótese transferível e teste discriminante proposto |
| --- | --- | --- |
| SwiftShader `c21aa26e0256a3efed81740d6215acb393b61d9e`; `src/Device/Renderer.cpp:917,1129`, `Clipper.cpp` | [Simplify clipping of complex polygons](https://swiftshader.googlesource.com/SwiftShader/+/c21aa26e0256a3efed81740d6215acb393b61d9e) troca flags intermediárias por clipping completo. O fonte expande linhas/pontos em quads e recorta antes de setupRoutine. | Referência de organização para conservar geometria após recorte. Controle CPU: primitiva expandida junto ao frustum; comparar polígono e interseções antes/depois. Escopo diferente: Vulkan e primitivas expandidas; não implementa GL_POLYGON/LINE/POINT nem prova correção Mesa. |
| ANGLE `38cac5971fef0a7b9b450a920ed538d9efb6fb88`; `StateManager11.cpp:2117,2139,2670`, `Blit11.cpp:1127` | [Fonte fixado](https://chromium.googlesource.com/angle/angle/+/38cac5971fef0a7b9b450a920ed538d9efb6fb88/src/libANGLE/renderer/d3d/d3d11/StateManager11.cpp): GS NULL tem serial zero; setGeometryShader envia GSSetShader e invalida shaders. syncDrawShaders restaura conjunto antes de limpar dirty. Blits também usam o state manager. | Referência para distinguir estado aplicado/desejado após operação interna. Mock CPU de device context: GS→NULL→mesmo GS, esperando três binds e invalidação nas transições, inclusive após blit. No Mesa incluir VS/copy shader e PM4. D3D11 não prova comportamento VGT/radeonsi. |

Commits, URLs e hashes: [manifest externo](../../lab/bancada/pesquisa/0125/external/manifest.json).
São recortes de fontes, não checkouts completos. SwiftShader passou coleta/triagem
na consulta Bateia, mas não teve contexto Clang. ANGLE foi descartado pela triagem:
GN não é reconhecido nesta jornada. Descarte não demonstra irrelevância; o achado
ANGLE vem da inspeção manual dos bytes. Não se criaram manifestos falsos nem bancos
de compilação sintéticos. Piglit também foi procurado, mas os arquivos específicos
não foram recuperados; nenhum achado lhe é atribuído.

## Verificação e limites

Testes reais Clang cobrem banco irmão, macro com espaços, header gerado, banco
maior que 4 MiB, precedência, JSON inválido, limite de 16 MiB e rejeição de
shell/plugins/response files. Não geram objeto nem executam compilador candidato.
Lint, testes raiz, vet e build passaram. O teste inicial sob sandbox falhou nas
fixtures HTTP por loopback bloqueado e passou ao habilitá-lo. Assinatura Git
desativada só no processo; build/cache em /tmp. Gates selados não executados.
O ensaio também encontrou bypass de orçamento: bytes.Buffer.ReadFrom herdado
permitia que io.Copy de os/exec contornasse bounded.Write. O adaptador passou
a aplicar o mesmo limite por ReadFrom, com teste de regressão. Resultados
anteriores que alegavam unidades completas estão em before-output-bound como
história superada, não como comprovação de cobertura sob o contrato. Uma rodada
complementar GS anterior foi interrompida após minutos decodificando AST enorme;
a pilha capturada registra a limitação. A rodada final com limite real terminou.
As tabelas propõem testes: esta entrega não executou fixtures dos projetos
investigados, não produziu dumps novos e não validou correção de driver.
