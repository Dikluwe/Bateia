# Pesquisa — Como verificar efeitos no Bateia

Pesquisa em fontes primárias, consultadas em 2026-09-08. Proposta de engenharia; nenhum verificador externo foi instalado ou executado nesta pesquisa.

## Diagnóstico local

O parser implementation_behavior.go observa atribuições, incrementos e retornos, mas não resolve aliases, fluxo de controle ou efeitos indiretos. action_assessment.go compara nomes de declarações/chamadas com famílias de ação. Esses módulos não demonstram efeitos. Há infraestrutura experimental de controles/mutantes nos gates, mas isso não constitui um verificador geral de projetos encontrados.

## Alternativas e força da evidência

1. **Tipos, SSA e fluxo de dados.** A representação SSA de Go contém operações de memória e identidade de chamadas, em vez de depender apenas de nomes. Permite investigar se uma escrita alcança parâmetro, receiver ou estado global e seguir chamadas estáticas. SSA não prova automaticamente alcance em execução ou ausência de efeitos; chamadas dinâmicas e código opaco exigem análise/modelos adicionais. Fonte: [Go SSA](https://pkg.go.dev/golang.org/x/tools/go/ssa).

2. **Análise interprocedural.** CodeQL modela fluxo local/global e distingue propagação de valores de taint. É referência para implementar rastreamento parâmetro→operação, com atenção a aliases e destinos de chamadas. Um caminho possível não demonstra que o efeito acontece em toda execução. Fonte: [CodeQL — data flow](https://codeql.github.com/docs/writing-codeql-queries/about-data-flow-analysis/).

3. **Contratos executáveis e fuzzing.** Formular pré-condições, resultado e estado antes/depois permite testar propriedades com entradas variadas e guardar contraexemplos mínimos. O fuzzing Go usa cobertura para explorar entradas. Passar sustenta somente o domínio e as execuções exercitados, não uma prova universal. Fonte: [Go Fuzzing](https://go.dev/doc/security/fuzz/).

4. **Mutation testing.** Alterar deliberadamente a operação e exigir que o teste rejeite a variante verifica se o oráculo distingue o efeito. Exemplos propostos: trocar AND por OR, eliminar uma escrita, escrever na entrada, ignorar cancelamento. Mutantes sobreviventes exigem investigar oráculo insuficiente, domínio restrito ou equivalência; não são todos defeitos automaticamente. Fonte: [Stryker — mutation testing](https://stryker-mutator.io/docs/). Stryker é referência metodológica, não recomendação de backend Go.

5. **Verificação formal para Go.** Gobra recebe programas anotados com contratos e invariantes e usa Viper/SMT. Por padrão prova correção parcial; terminação precisa de obrigação adicional. Pode fundamentar propriedades de heap/permissões e pós-condições no fragmento suportado, sob especificações e modelos fornecidos. Não descobre sozinho o contrato correto em linguagem natural. Fonte: [tutorial Gobra](https://github.com/viperproject/gobra/blob/master/docs/tutorial.md).

6. **C/C++.** CBMC verifica assertions e segurança por model checking com desenrolamento de loops; limites e sua suficiência precisam constar no recibo. Frama-C/WP verifica contratos ACSL para C. São adaptadores distintos por linguagem, não regras de Go aplicadas globalmente. Fontes: [CBMC](https://www.cprover.org/cbmc/), [Frama-C/WP](https://www.frama-c.com/fc-plugins/wp.html).

## Particularidade de efeitos assíncronos

Ausência de evento durante uma espera não prova ausência futura. Go synctest permite controlar tempo e esperar quiescência em operações suportadas dentro da sua bolha; I/O real e syscalls não se tornam controláveis por isso. O race detector só encontra corridas exercitadas durante execução. Fontes: [Testing Time](https://go.dev/blog/testing-time), [synctest](https://pkg.go.dev/testing/synctest), [race detector](https://go.dev/doc/articles/race_detector).

Para cancelamento, propor contratos separados: registro removido do estado interno; estímulos novos não produzem notificações após uma barreira definida; eventos já enfileirados têm semântica explicitada. Um backend simulado valida o modelo, não o comportamento do sistema operacional. Teste de integração real precisa declarar plataforma, backend, sincronização e janela observada.

## Caminho recomendado para o projeto

Minha recomendação é um piloto de **contrato de efeito + execução reproduzível + mutantes**, começando por transformações determinísticas de dados, com pouca dependência externa. Isso cria uma medida independente de qualidade para avaliar futuramente SSA ou verificadores formais.

Exemplo de obrigação: uma transformação retorna o resultado especificado e preserva os valores de entrada após retornar. Isso é diferente de provar que nunca houve uma escrita transitória. Propriedades de não compartilhamento de memória também precisam ser separadas. Oráculo simples e independente, snapshots profundos e entradas com aliasing quando permitido pelo contrato evitam comparar somente ponteiros ou repetir a implementação.

Antes do candidato: congelar contrato, domínio de entradas, política de erros, observáveis e mutantes. Executar em ambiente isolado com orçamento, versão/commit/hash e sem credenciais; observar efeitos reais em vez de aceitar somente exit code zero. Esses são requisitos propostos para o futuro executor, não uma ação realizada nesta pesquisa.

Em paralelo ou depois do piloto, adicionar resumos estáticos por identidade de objeto: escrita possível em parâmetro/receiver/global, alocação, retorno e chamadas não resolvidas. Ausência de efeitos não é inferida de análise incompleta. Carregamento de pacotes/compilação fica na infraestrutura; regras e comparação de observações no núcleo puro, respeitando as camadas Tekt.

## Contrato de saída proposto

Não comprimir evidências numa nota única. Registrar separadamente:
- propriedade e versão do contrato;
- mecanismo: sintaxe, análise estática, execução ou prova formal;
- resultado: sustentado no escopo, contraexemplo ou desconhecido;
- escopo: entradas, caminhos/modelos, plataforma, dependências, limites e hipóteses;
- evidências: estado antes/depois, rastros, contraexemplo, ferramenta/versão, comandos e hashes.

Prova formal pode ser geral dentro do modelo; teste bem-sucedido permanece evidência das execuções. Timeout, linguagem sem suporte e chamada opaca não viram aprovação. Cobertura de linhas e percentual de mutantes detectados não viram probabilidade de correção.

O critério de ganho do piloto deve ser: rejeitar candidatos que violam o efeito, inclusive candidatos com nomes convincentes, e preservar Unknown quando o mecanismo não consegue decidir. Novas perguntas devem ser congeladas antes da calibração; recuperar 6/6 alvos conhecidos continua insuficiente para medir essa capacidade.
