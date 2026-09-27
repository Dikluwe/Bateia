# Simulação do intake avançado com 10 agentes Luna — passo 0083

Data: 2026-09-24.

## Escopo e método

Dez agentes `gpt-6-luna`, cada um com uma persona e um tema distintos, simularam pessoas
usando o executável real produzido pelo worktree. Cada ensaio recebeu uma entrada JSON,
executou `idea intake prepare`, inspecionou as necessidades e, quando aplicável, respondeu
com `USER` confirmado por `idea intake apply`. Um caso também exercitou a ferramenta MCP
avançada. Todos os artefatos de execução ficaram em `/tmp/bateia-human-01` a
`/tmp/bateia-human-10`; nenhum ensaio iniciou recomendação, rede ou execução de candidato.

Isto é um ensaio sintético de contrato e usabilidade, não pesquisa com seres humanos. Os
agentes conhecem software melhor que parte das personas representadas, de modo que a
facilidade observada provavelmente superestima a facilidade para usuários reais.

## Resultado agregado

- 6 entradas incompletas produziram 7 necessidades humanas; todas chegaram a `READY`
  após respostas tipadas válidas.
- 2 entradas completas chegaram diretamente a `READY` e permaneceram iguais após
  `apply` com lista vazia.
- 1 entrada com valor explícito inválido chegou a `INVALID`, sem pergunta ou caminho de
  correção no protocolo.
- 1 entrada com campo desconhecido foi rejeitada antes do intake; a pessoa só se
  recuperou removendo manualmente o campo e repetindo o comando.
- A tentativa sem `confirmed: true` foi rejeitada corretamente.
- Valor vazio explícito foi preservado como `EXPLICIT/USER`, não confundido com default.
- Respostas em ordem diferente da ordem das perguntas foram aplicadas corretamente por
  `need_id`.

| Caso | Persona e tema | Estado inicial | Perguntas | Resultado |
| --- | --- | --- | ---: | --- |
| 01 | Coordenadora de clínica; filas e faltas | `NEEDS_INPUT` | 1 | `READY` após capacidades |
| 02 | Gerente de restaurante; desperdício | `NEEDS_INPUT` | 1 | `READY` após capacidades |
| 03 | Engenheira de campo; inspeção offline | `NEEDS_INPUT` | 1 | `READY` após problema; resposta não confirmada rejeitada |
| 04 | Diretora de ONG; CRM de doadores | `READY` | 0 | Preservado; `apply []` foi no-op |
| 05 | Produtora de cursos; QA de legendas | `NEEDS_INPUT` | 1 | `READY` após problema |
| 06 | Pesquisadora; reprodutibilidade | `NEEDS_INPUT` | 1 | `READY`; vazio explícito preservado |
| 07 | Família; finanças privadas | `NEEDS_INPUT` | 2 | `READY`; respostas fora de ordem aceitas |
| 08 | Educadora comunitária; idioma offline | `READY` | 0 | Sem alerta para contradição operacional |
| 09 | Engenheiro industrial; manutenção | `INVALID` | 0 | Sem correção por resposta tipada |
| 10 | Advogada; triagem sob LGPD | erro de schema, depois `READY` | 0 | Recuperação exigiu edição manual |

## O que funcionou

O protocolo se absteve de completar problema ou capacidades por conta própria. As
perguntas carregaram caminho, tipo, regra de validação, autoridade e efeito; respostas
confirmadas alteraram apenas o campo alvo e geraram novo hash. A proveniência distinguiu
valores explícitos, defaults versionados e alvo derivado. Nenhuma preferência humana foi
promovida a evidência e nenhum efeito foi iniciado pelo intake.

O comportamento foi estável em casos omitidos, vazios explícitos e completos. A exigência
de `USER` confirmado foi efetiva, e a identificação por `need_id` evitou dependência da
ordem conversacional.

## Problemas encontrados

### 1. `READY` estrutural pode parecer autorização ou segurança operacional

O caso 08 ficou `READY` com `execution.mode=local`, `provider=sourcegraph`,
`local_paths=[]` e `policy.require_local=false`. Não houve issue nem pergunta, apesar da
restrição explícita de soberania comunitária dos dados. Nos casos de saúde, finanças e
LGPD, `live`, Sourcegraph e ausência de exigência local também foram materializados como
defaults sem consentimento específico.

As constraints em texto são preservadas, mas não causam política nem validação. Assim,
`READY` significa apenas “campos mínimos válidos”, significado fácil de confundir com
“seguro para executar” ou “adequado à LGPD”. Esta foi a falha mais importante do ensaio.

### 2. Entrada inválida não é corrigível pelo protocolo de perguntas

No caso 09, `search.max_candidates=0` foi corretamente preservado e diagnosticado com a
faixa `[1,20]`, mas o intake retornou `questions=[]`. `apply []` saiu com código zero e
gravou outro intake `INVALID` idêntico. A pessoa precisa editar o rascunho fora do
protocolo, mesmo sendo uma correção tipável e local.

### 3. A projeção textual ainda não é uma pergunta humana

Os assuntos emitidos foram “capacidades obrigatórias da solução” e “problema que a
recomendação deve resolver”. A issue plural usa “não informado”. Pessoas menos técnicas
consideraram “capacidades obrigatórias” abstrato. Quando ambos faltam, capacidades vêm
antes do problema por desempate lexical de ID, uma ordem pouco natural para elicitação.

### 4. Onboarding exige conhecer o envelope JSON

Texto livre, campos `problem/capabilities` no topo e a chave intuitiva `tema` foram
rejeitados. A validação estrita é desejável, mas a mensagem não mostra o formato aceito
nem como gerar um template. Um agente chegou a acionar a busca simples antes de encontrar
o caminho avançado correto.

### 5. Defaults são auditáveis, porém pouco explicados

Proveniência e `source_ref` tornam os defaults rastreáveis, mas o JSON longo não explica
suas consequências. `live`, `sourcegraph`, idiomas `pt-BR/en-US`, política permissiva de
licenças e limiares aparecem sem resumo humano. `apply` com respostas vazias conclui em
silêncio, inclusive quando o intake permanece inválido.

## Próximas correções recomendadas

1. Adicionar validações cruzadas determinísticas: modo `local` requer caminhos locais e
   não pode herdar provedor remoto silenciosamente.
2. Estruturar autorização de rede, egressão de dados e requisito de processamento local;
   se ausentes antes de um modo `live`, emitir `REQUEST_AUTHORIZATION`, não inferir a
   intenção a partir de texto livre.
3. Distinguir externamente `READY_FOR_RECOMMENDATION` de autorização para executar e
   emitir um resumo de defaults e riscos antes de qualquer efeito.
4. Converter erros editáveis em `InformationNeed` tipada, incluindo resposta inteira e
   faixa permitida para limites numéricos.
5. Projetar perguntas por templates humanos: primeiro o problema, depois “quais funções
   são indispensáveis?”, mantendo os metadados como contrato interno.
6. Fazer `apply` reportar no-op e retornar erro quando respostas vazias não resolvem um
   estado `INVALID` ou `NEEDS_INPUT`.
7. Oferecer `idea intake template` ou aceitar um `IdeaSpec` simples como atalho, sem
   remover a leitura estrita do envelope avançado.

## Veredito

A primeira fatia prova o mecanismo central: a IA pode receber necessidades tipadas e
perguntar somente o que exige autoridade humana, enquanto a aplicação das respostas
permanece determinística. Ela ainda não está pronta para mediar execução em contextos de
privacidade, offline ou soberania de dados, porque completude estrutural e segurança
operacional continuam misturadas na projeção `READY`. O próximo incremento deve atacar
essa fronteira antes de conectar perguntas a runs duráveis.
