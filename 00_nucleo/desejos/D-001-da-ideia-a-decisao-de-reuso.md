# D-001 — Da ideia à decisão de reuso

**Estado**: `DESEJO`  
**Registrado em**: 2026-08-23  
**Origem**: intenção inicial do Bateia, reafirmada após a autoanálise e a busca
poliglota por projetos semelhantes  
**Relacionado**: `caso-de-uso-motor-de-busca-v2.md`

## O desejo

Queremos descrever uma ideia de software ao Bateia e receber evidência suficiente
para **não começar do zero quando o problema, uma parte dele ou uma fundação
reutilizável já existir**.

O resultado principal não deveria ser apenas uma lista de repositórios que contêm
as mesmas palavras. Deveria ajudar a decidir:

- reutilizar um projeto quase inteiro;
- estender um núcleo existente;
- integrar uma biblioteca ou componente;
- compor capacidades encontradas em projetos diferentes;
- usar um projeto apenas como referência;
- construir do zero porque nenhuma base adequada foi encontrada;
- abandonar ou reformular a ideia porque ela já está suficientemente resolvida.

## O problema preservado

Antes de iniciar um projeto, descobrir o que já existe exige transformar uma ideia
vaga em muitas perguntas, buscar com vocabulários diferentes, eliminar resultados
que apenas repetem palavras, ler documentação e código, comparar capacidades e
avaliar licença, maturidade e esforço de adaptação.

O Bateia nasceu para comprimir esse trabalho numa decisão fundamentada. A
implementação atual já descobre e extrai repositórios, mas ainda aproxima
"mesma ideia" por correspondência lexical. A execução de 2026-08-23 mostrou o
desvio: 26 candidatos, maior similaridade de 18%, resultados irrelevantes com
relevância alta e nenhuma recomendação automática útil.

## A expansão desejada

Uma descrição não deve virar uma única query longa. O Bateia deveria primeiro
decompô-la numa especificação de intenção e depois gerar **várias buscas pequenas,
independentes e rastreáveis**.

Exemplo de ideia:

> Ferramenta local que mostra o impacto de alterar um símbolo e entrega contexto
> estrutural verificável para agentes de programação.

Expansões possíveis:

```text
symbol dependency graph
code impact analysis
blast radius source code
repository memory MCP
architecture dependency matrix
incremental AST index
diff aware code intelligence
callers and transitive dependencies
```

Cada expansão deve declarar qual parte da ideia tenta encontrar. Uma expansão que
não produz sinal pode ser alargada; uma que produz ruído pode ser descartada sem
contaminar as demais.

## Forma desejada da ideia

A forma exata ainda é decisão futura, mas queremos preservar pelo menos estas
dimensões:

```go
type IdeaSpec struct {
    Problem        string
    Users          []string
    Capabilities   []string
    Constraints    []string
    NonGoals       []string
    PreferredStack []string
}
```

- **Problema**: qual dor ou decisão a ideia resolve;
- **Usuários**: para quem a solução existe;
- **Capacidades**: o que precisa ser possível fazer;
- **Restrições**: local, offline, licença, linguagem, privacidade, desempenho;
- **Não objetivos**: semelhanças aparentes que não devem atrair candidatos;
- **Stack preferida**: preferência de integração, não definição da identidade da
  ideia.

## Funil desejado

```text
descrição da ideia
    ↓
IdeaSpec: problema, usuários, capacidades, restrições e não objetivos
    ↓
expansões de descoberta por problema, capacidade e vocabulário alternativo
    ↓
candidatos reunidos e deduplicados
    ↓
extração de evidências: README, docs, manifesto, estrutura, símbolos e código
    ↓
comparação entre capacidades desejadas e capacidades comprovadas
    ↓
decisão de reuso: FORK · EXTEND · COMPOSE · INTEGRATE · REFERENCE · BUILD
```

## Evidência, não adivinhação

Uma capacidade só deveria contar no score quando o relatório puder apontar a
evidência que a sustenta, por exemplo:

- trecho de documentação que declara a capacidade;
- dependência ou manifesto compatível;
- símbolo, pacote ou módulo que a implementa;
- teste ou exemplo que demonstra o comportamento;
- caminho exato e referência `repo@commit`;
- ausência explícita quando a capacidade não foi encontrada.

"Parece semelhante" pode ajudar na descoberta, mas não deveria sustentar sozinho
uma recomendação de reuso.

## Comparação desejada

O score deveria priorizar a proximidade do problema e a cobertura de capacidades,
não a repetição de termos ou a coincidência de stack. Uma hipótese inicial para
experimento, não uma decisão:

```text
35% problema resolvido
30% capacidades comprovadas
15% arquitetura e modo de operação compatíveis
10% maturidade, licença e manutenção
10% stack e custo de integração
```

O relatório deve manter os vetores separados. Um total alto não pode esconder uma
restrição eliminatória ou ausência da capacidade central.

## Saída desejada

Para cada candidato relevante, queremos responder:

| Pergunta | Resposta esperada |
| :------- | :---------------- |
| O que ele já resolve? | Capacidades comprovadas e respectivas evidências |
| O que podemos reutilizar? | Projeto, biblioteca, módulo, algoritmo ou padrão |
| O que falta? | Lacunas em relação à `IdeaSpec` |
| Qual o custo provável? | Compatibilidade, dependências e esforço de adaptação |
| Qual a decisão sugerida? | `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou `BUILD` |
| Por quê? | Cadeia de evidências e incertezas declaradas |

O relatório final deveria ser uma **decisão de partida**, não uma página de
resultados de busca.

## Exemplo do valor pretendido

Para a ideia de contexto estrutural para agentes, um resultado útil poderia dizer:

```text
Não comece do zero.

- Use o grafo resolvido e o raio de impacto do tekt-cargo-dsm como núcleo.
- Adote o padrão MCP e a memória de repositório observados no Sverklo.
- Considere a indexação incremental e as consultas de símbolos do CodeGraph.
- Construa apenas a integração entre essas capacidades e a camada poliglota.
```

Cada frase precisaria apontar para evidência extraída, além de declarar limites de
licença, linguagem e maturidade.

## Sinais de que o desejo foi atendido

Ainda não são critérios de aceite de uma versão. São sinais de valor:

1. Uma ideia pode gerar várias estratégias de descoberta compreensíveis e
   auditáveis.
2. O alvo e os candidatos são fluxos separados: o alvo define a necessidade; a
   descoberta encontra alternativas.
3. Projetos lexicalmente parecidos, mas funcionalmente divergentes, recebem score
   baixo com motivo explícito.
4. Cada capacidade atribuída a um candidato possui evidência rastreável.
5. O relatório mostra componentes reutilizáveis e lacunas, não apenas dependências.
6. O usuário consegue tomar uma decisão `reusar versus construir` sem ler todos os
   repositórios candidatos integralmente.
7. Quando não há evidência suficiente, o Bateia diz "inconclusivo" em vez de
   fabricar precisão.

## Não decisões

Este desejo não decide:

- usar ou não embeddings ou modelos de linguagem;
- qual provedor fará expansão semântica;
- pesos definitivos do score;
- formato definitivo da `IdeaSpec`;
- número de queries ou candidatos;
- se o motor será local, remoto ou híbrido;
- em qual versão isso entrará.

Essas escolhas devem nascer de experimentos comparáveis. Tree-sitter, embeddings,
grafos e reranking são meios possíveis; nenhum deles substitui a separação
conceitual entre **ideia, descoberta, evidência e decisão de reuso**.

## Delimitação operacional — 2026-09-07

O [plano 0043](../passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md) delimita as promessas deste desejo. Documentação sustenta declarações; manifestos indicam dependências; código e testes sustentam conclusões somente no escopo examinado. Capacidade não encontrada permanece não verificada, salvo evidência específica de ausência. Nenhum resultado de busca prova inexistência global de alternativas. Cobertura funcional não garante substituição ou integração, e esforço permanece estimativa até haver avaliação correspondente. O plano preserva a intenção acima; a [execução 0043](../relatorios/execucao-promessas-buscas-0043.md) registra as correções, os ensaios e a validação externa ainda pendente.
