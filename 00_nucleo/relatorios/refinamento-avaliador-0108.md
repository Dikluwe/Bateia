# Pesquisa 0108 — Refinamento antes do avaliador estrutural

**Data**: 2026-10-06  
**Estado**: PESQUISA CONCLUÍDA; PILOTO AINDA NÃO IMPLEMENTADO  
**Base**: [pesquisa 0107](pesquisa-classificador-0107.md)  
**Evidência**: [bancada 0108](../../lab/bancada/pesquisa/0108/README.md)

## O que mudou na proposta

O piloto deve começar por observações estruturais delimitadas e separar
explicitamente essas observações das conclusões sobre comportamento. Não basta
encontrar make, select, go ou ctx.Done para afirmar que uma investigação tem
capacidade limitada, cancelamento completo ou ausência de vazamento.

A unidade recomendada é propriedade + operação/recurso + função/caminho +
configuração + evidência. O repositório inteiro não deve receber uma classificação
global por causa de um trecho. A categoria do mecanismo continua opcional quando
há uma observação útil fora do catálogo atual. Esta é uma proposta de contrato,
sem alteração de tipos ou regras produtivas.

## Semântica que o piloto precisa respeitar

A capacidade de um canal limita seu buffer. Um canal sem buffer depende da
disponibilidade de emissor e receptor; um canal nil não está pronto para comunicar.
Da capacidade do buffer não se deduz o número de produtores bloqueados, workers
ou tarefas pendentes fora dele. Portanto, registrar capacidade de canal e limite
de trabalho admitido como propriedades distintas.
[Especificação Go](https://go.dev/ref/spec#Channel_types).

CancelFunc sinaliza abandono, mas não espera o trabalho terminar. A propriedade
de cancelamento precisa indicar seu alcance: chamada de admissão, worker, operação
de I/O e encerramento de recursos não são equivalentes. Um ramo ctx.Done pode
permitir que o pai retorne enquanto o filho permanece bloqueado enviando uma
resposta. A primeira afirmação é contrato de biblioteca; a segunda é um padrão
documentado no estudo GCatch, não uma acusação sobre o Bateia.
[Context](https://pkg.go.dev/context#CancelFunc),
[GCatch](https://songlh.github.io/paper/gcatch.pdf).

Um canal bloqueante pode expressar backpressure intencional. Não deve ser
classificado como defeito apenas por bloquear. Também não se pode inferir que
uma cláusula default rejeita trabalho: é necessário verificar o que esse ramo
faz, qual recurso está envolvido e quais efeitos aconteceram antes dele.
Essas são obrigações propostas para nossa avaliação, derivadas da distinção
entre forma sintática e efeito observado.

## Até onde análise estática ajuda

GCatch modela operações de canais, sincronização e caminhos para encontrar
possíveis bloqueios. O artigo também registra falsos positivos ligados a alias
e grafo de chamadas. Goat usa interpretação abstrata e fragmentos localizados
para conter o espaço de estados. A lição para nós é declarar modelo e alcance
de cada análise, em vez de atribuir autoridade universal ao fato de usar AST.
[GCatch, limites](https://songlh.github.io/paper/gcatch.pdf),
[Goat, autores](https://cs.au.dk/~amoeller/papers/goat/).

AST permite observar construções. Tipos ajudam a distinguir objetos homônimos;
SSA permite representar fluxo e identidades tipadas. A documentação de SSA
descreve instruções de canais/select e um grafo de controle com dominância.
Esses recursos apoiam uma verificação delimitada, mas não fornecem sozinhos
uma prova de interação entre todas as goroutines.
[SSA oficial](https://pkg.go.dev/golang.org/x/tools/go/ssa).

Recomendação incremental: começar com AST e observações de forma explícitas;
usar resolução de tipos quando a identidade da operação for necessária; depois
avaliar SSA para propriedades de ordem e fluxo. Construções não suportadas,
aliases não resolvidos, wrappers desconhecidos e dependências ausentes devem
gerar lacuna específica. Não preencher essas lacunas por nomes de símbolos.
Não foi escolhida ou instalada uma versão de x/tools nesta pesquisa.

## Configuração também faz parte da evidência

Arquivos Go disponíveis podem diferir dos arquivos compilados por tags, plataforma,
cgo e variantes de teste. O loader go/packages usa a ferramenta de build e
distingue GoFiles de CompiledGoFiles. Um snapshot suficiente para procurar texto
pode ser insuficiente para analisar tipos e comportamento.
[Loader oficial](https://pkg.go.dev/golang.org/x/tools/go/packages).

O recibo futuro deve conservar fonte/commit/hash, arquivos selecionados, versão
do Go e do analisador, GOOS/GOARCH, tags relevantes, modo de análise, dependências
e falhas de parsing/tipos. A ausência desses dados restringe o alcance da conclusão.
Para preservar Tekt, parsing e regras sobre dados em memória podem ser L1;
loader com filesystem/ferramentas externas e execução concreta pertencem a L3.
L4 continua composição. Não colocar go/packages.Load no núcleo puro.

## Bancada externa sem referência circular

O GoBench correto é timmyyuan/gobench. Seu acervo descreve bugs reais e kernels
extraídos, incluindo sobreposição de origem entre conjuntos. Isso oferece
contraexemplos externos, mas não um catálogo completo de filas e cancelamento.
Não tratar kernels, seus projetos de origem e variantes do mesmo bug como
amostras independentes. Separar por origem/família antes de dividir desenvolvimento
e avaliação, e adicionar exemplos corretos e casos fora do alcance.
[Repositório primário](https://github.com/timmyyuan/gobench),
[artigo](https://lujie.ac.cn/files/papers/GoBench.pdf).

Os rótulos de bug não respondem automaticamente às nossas propriedades. Para
cada caso, precisamos registrar a pergunta, condições e testemunho. Um bug
conhecido pode coexistir com capacidade finita; um exemplo corrigido pode ainda
ter um comportamento bloqueante intencional. Revisão dos rótulos deve ocorrer
separadamente do produtor das avaliações. Nenhum caso foi executado nesta rodada.

O pacote go/analysis oferece passes, fatos entre pacotes e o harness analysistest
para conferir diagnósticos esperados. Esse harness pode verificar o contrato
do futuro analisador; a qualidade semântica das expectativas ainda exige revisão.
[API oficial](https://pkg.go.dev/golang.org/x/tools/go/analysis).

## Ensaios de concorrência sem janelas frágeis

O detector de races observa execuções: não cobre caminhos não executados e
não constitui prova de ausência de deadlock. Goleak verifica goroutines restantes,
mas sua documentação registra problemas de atribuição em testes paralelos.
São observáveis diferentes e devem ter recibos e escopos separados.
[Race detector](https://go.dev/doc/articles/race_detector),
[goleak](https://github.com/uber-go/goleak).

testing/synctest é uma opção para futuros kernels isolados de concorrência:
tempo virtual e espera por bloqueio durável evitam inferir sucesso pelo tempo
de parede. Bloqueio de mutex, I/O e eventos externos não recebe a mesma garantia;
essa técnica não substitui o teste de persistência do runtime. As APIs Test e
Wait foram confirmadas no Go1.26.5 local. A documentação atual também mostra
APIs posteriores, que não devem ser exigidas pelo piloto.
[Synctest oficial](https://pkg.go.dev/testing/synctest).

## Resultado do uso do Bateia

idea plan produziu dez expansões. Quatro consultas manuais ao catálogo tiveram
respostas SUCCEEDED: uma vazia e três com candidatos. Nenhum candidato teve
commit resolvido. A implementação atual ignora o erro de metadata ao compor
o resultado de catálogo: sucesso da busca não equivale a resolução imutável.
Não houve coleta de código nesta rodada.

O catálogo trouxe system-pclub/GCatch, cmpxchg16/gobench e mosajjal/sniproxy.
O primeiro é referência pertinente pelo artigo; o segundo é um load tester
homônimo, diferente do benchmark de bugs; o terceiro é uma aplicação que
correspondeu aos termos, não metodologia de avaliação. A triagem marcou os
três REJECT. Sua pontuação não substituiu a inspeção manual de relevância.

A identificação do GoBench correto, Goat e as premissas de Go foram conferidas
em fontes primárias externas pelo agente. Não atribuir esse aprofundamento
ao classificador do Bateia. Os dois estágios offline de expansão e triagem
foram reproduzidos byte a byte; referências e seleção estão registradas no Lab.

## Proposta atual para a primeira implementação

Construir um observador exploratório no Lab para formas locais de canal e
select, com resultados por propriedade e escopo. Em casos selecionados,
distinguir capacidade do canal, limite de workers, comportamento da admissão,
configuração obrigatória/opcional e alcance do cancelamento. Quando a análise
não consegue decidir uma dessas propriedades, conservar as observações que
conseguiu localizar e registrar a lacuna, sem promover a categoria inteira.

As primeiras métricas devem separar correção da localização, correção da
propriedade sob condições, conclusões excessivas, cobertura e autoridade
atribuída. Usar pares equivalentes, mudanças que alteram a propriedade,
exemplos corretos, contraexemplos e casos não suportados. A família de cada
caso conta na avaliação; aumentar variantes não cria independência estatística.

Essa rodada ajustou o escopo e forneceu referências para implementação. Não
comprovou qualidade do classificador, não executou ferramentas de candidatos
nem implementou o observador. O código Go produtivo permanece inalterado.
