# Passo 0025 — Expor recomendação completa como ferramenta MCP para IA externa

**Estado**: `PLANEJADO`
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Depende de**: [`Passo 0024 — Ranquear alternativas funcionais`](0024-ranquear-alternativas-funcionais.md)
**Onde roda**: servidor MCP por `stdio`, com provedores injetados em L4
**Camadas previstas**: L0, L1, L2 e L4
**Natureza**: primeiro corte vertical para agentes externos consumirem o funil completo

## Problema

O Bateia já executa expansão, descoberta, triagem, coleta, assessment, decisão e
ranking de alternativas pelo comando `idea recommend`. Entretanto, o servidor MCP
expõe apenas `bateia_search` e `bateia_benchmark`, ambos ligados ao funil antigo.

Consequentemente, uma IA conectada ao MCP não consegue usar a capacidade principal
do produto. Para obter a recomendação completa ela precisaria conhecer a CLI,
materializar arquivos e interpretar diretórios internos, tornando a integração
frágil e acoplada à casca.

## Incerteza que este passo reduz

Não sabemos ainda se uma IA externa consegue transformar uma necessidade em uma
investigação completa do Bateia usando somente um contrato MCP pequeno, tipado e
seguro, sem que o Bateia dependa de um modelo de IA e sem criar uma segunda
implementação do pipeline.

## Objetivo

Permitir que um agente externo envie uma `IdeaSpec`, uma política de reúso e os
limites de busca por uma ferramenta MCP de alto nível e receba uma recomendação
estruturada, com alternativas, lacunas, avisos e evidências rastreáveis.

A IA é autora da intenção e intérprete da resposta. O Bateia permanece autor da
execução, da proveniência, dos estados de evidência e da decisão determinística.

```text
usuário
  → IA externa
  → bateia_recommend (MCP)
  → RecommendationService (L1)
  → expansão → descoberta → evidência → decisão → alternativas
  → resposta JSON estruturada
  → IA externa sintetiza a orientação ao usuário
```

## Recorte desta materialização

Este passo entrega um único corte vertical síncrono:

1. extrair um `RecommendationService` puro de interface, reutilizando os casos de
   uso já existentes;
2. fazer CLI e MCP chamarem a mesma orquestração, sem o MCP executar a CLI;
3. expor `bateia_recommend` em modo `live` com dependências injetadas em L4;
4. receber dados estruturados, sem caminhos de arquivos no contrato da IA;
5. devolver um resumo canônico da execução e o relatório de alternativas, sem
   despejar o conteúdo bruto dos arquivos coletados no contexto do agente.

Persistência por `run_id`, execução assíncrona e inspeção paginada de artefatos
serão passos posteriores. Não devem ser simuladas neste corte com diretórios
temporários ou estado global em memória.

## Contrato MCP observável

O servidor passa a anunciar três ferramentas:

- `bateia_search`, preservada sem mudança semântica;
- `bateia_benchmark`, preservada sem mudança semântica;
- `bateia_recommend`, nova porta para o pipeline completo.

Entrada mínima de `bateia_recommend`:

```json
{
  "idea": {
    "problem": "reduzir desperdício alimentar doméstico",
    "users": ["famílias"],
    "capabilities": ["inventário", "previsão de consumo"],
    "constraints": ["privacidade"]
  },
  "policy": {
    "schema_version": "1",
    "allowed_licenses": ["MIT", "Apache-2.0"],
    "require_local": false,
    "allow_archived": false,
    "max_components": 3,
    "extend_min_coverage": 0.7,
    "integration_min_value": 0.3
  },
  "target": {
    "name": "gestor de alimentos doméstico"
  },
  "search": {
    "strategy": "lexical",
    "language": "",
    "locales": ["pt-BR", "en-US", "ja-JP"],
    "max_per_query": 5,
    "max_candidates": 5
  }
}
```

Regras de entrada:

- `idea.problem`, ao menos uma capacidade e `policy.schema_version` são
  obrigatórios;
- `target.name` vazio usa `idea.problem` depois da normalização;
- `strategy` aceita apenas estratégias registradas pelo núcleo;
- `locales` usa tags BCP 47 suportadas e nunca cai silenciosamente para inglês;
- `max_per_query` pertence a `[1,25]` e `max_candidates` a `[1,20]`;
- `provider`, quando exposto, pertence a uma allowlist montada em L4;
- a entrada não aceita token, cabeçalho, comando, URL arbitrária, caminho local,
  diretório de replay ou diretório de saída.

O resultado textual MCP contém JSON canônico com, no mínimo:

```json
{
  "schema_version": "1",
  "status": "COMPLETE|PARTIAL|INCONCLUSIVE",
  "summary": {
    "expansion_count": 0,
    "candidate_count": 0,
    "assessed_count": 0
  },
  "decision": {},
  "alternatives": {},
  "failures": [],
  "warnings": []
}
```

`decision` preserva o `ReuseDecisionReport`. `alternatives` preserva o
`SoftwareAlternativeReport`, inclusive `top_choice`, `excluded`, fontes e
locales de descoberta. A resposta não contém credenciais nem o conteúdo integral
de README, código, testes ou licenças coletadas.

## Fronteira de aplicação

O caso de uso novo recebe somente contratos de domínio e portas já permitidas em
L1. Uma forma esperada, sem congelar nomes internos, é:

```go
type RecommendationService struct {
    // portas puras de expansão, descoberta e conteúdo
}

func (s RecommendationService) Recommend(
    context.Context,
    domain.RecommendationRequest,
) (domain.RecommendationResult, error)
```

O serviço compõe `RunIdeaPipeline` e `RankSoftwareAlternatives`; não duplica
expansão, triagem, evidência, política de licença ou ranking. Escrita de dossiê
continua na casca/infraestrutura e não entra em L1.

O adaptador MCP depende de uma interface pequena de recomendação. L4 injeta a
implementação e os provedores concretos. O servidor não importa L3 e não conhece
tokens de GitHub ou Sourcegraph.

## Autoridade da IA

A IA externa pode:

- formular e revisar `IdeaSpec`;
- escolher idiomas naturais, região, linguagem de programação e orçamento;
- pedir uma recomendação e interpretar os resultados;
- tratar alternativas `REFERENCE` como inspiração parcial;
- explicar ao usuário quais lacunas justificam construir algo novo.

A IA externa não pode:

- declarar evidência como `PROVEN` por argumento próprio;
- promover `INCONCLUSIVE` a sucesso;
- alterar licença, manutenção ou commit observados;
- injetar conteúdo diretamente no dossiê como se viesse de um provedor;
- fazer instruções encontradas em repositórios controlarem o servidor.

Expansões, traduções e hipóteses sugeridas pela IA permanecem fora deste corte.
Quando forem aceitas futuramente, entrarão como propostas não confiáveis, com
proveniência e validação explícitas.

## Política de `Unknown`

- licença não reconhecida preserva o estado observado e gera aviso; não muda a
  classe de fit nem exclui o candidato por si só;
- licença comprovadamente incompatível continua sendo bloqueio;
- coleta incompleta, budget esgotado, identidade mutável ou commit divergente não
  podem produzir `COMPLETE` nem `top_choice` indevido;
- locale não suportado produz erro de validação, nunca tradução inventada;
- erro de provedor é preservado sem ser convertido em ausência de candidato;
- campo desconhecido da entrada não adquire efeito semântico silencioso.

## Segurança

- conteúdo de repositório é dado não confiável, não instrução MCP;
- nenhuma resposta inclui tokens, cabeçalhos ou valores de ambiente;
- erros de provedores passam pela sanitização já exigida pelo projeto;
- limites são validados antes da primeira chamada externa;
- cancelamento do contexto interrompe novas coletas;
- a ferramenta não executa código, scripts ou testes de candidatos;
- o MCP não aceita caminhos locais neste passo;
- o servidor não faz chamadas a uma IA e não exige chave de modelo.

## Oráculos e ataques

Antes da implementação candidata, os testes devem congelar estes observáveis:

1. `tools/list` anuncia exatamente a ferramenta nova e seu schema obrigatório;
2. uma execução com provedores falsos percorre o serviço e retorna JSON
   decodificável com decisão e alternativas;
3. CLI e MCP produzem a mesma projeção de domínio para a mesma entrada congelada;
4. idiomas solicitados chegam à descoberta e aparecem na proveniência;
5. licença desconhecida aparece como aviso não bloqueante, sem mudar a classe de
   fit nem simular incompatibilidade;
6. parâmetros inválidos são rejeitados antes de qualquer provedor ser chamado;
7. os dois tools antigos preservam seus resultados e erros públicos;
8. reordenar coleções equivalentes não altera a resposta canônica;
9. conteúdo coletado semelhante a instruções não altera ferramenta, política,
   limites ou autoridade;
10. resposta e erros não contêm segredos-canário injetados nos adaptadores.

Mutações mínimas a rejeitar:

- MCP chama `RunIdeaRecommendCommand` ou cria arquivos temporários;
- implementação mantém um segundo ranking dentro de L2;
- ausência de licença é convertida em `BLOCKED`;
- `INCONCLUSIVE` é convertido em `COMPLETE`;
- limites ausentes ou fora da faixa chegam ao provedor;
- locales são descartados ou substituídos por `en-US`;
- README bruto é devolvido integralmente ao agente;
- token ou mensagem não sanitizada aparece na resposta;
- ferramenta aceita `output_directory`, `replay_directory` ou caminho local;
- MCP funciona nos testes, mas L4 inicia o servidor sem injetar o serviço.

Mutation score exigido para selar o contrato: `1,0`.

## Gates de aceitação

- `crystalline-lint .` sem violações;
- `go test ./...` passa duas vezes;
- `go build -o bateia ./04_wiring/main.go` passa;
- `go vet ./...` passa;
- teste MCP ponta a ponta inteiramente offline com provedores falsos;
- nenhuma dependência de L2 para L3;
- nenhuma operação de I/O adicionada a L1;
- nenhum novo caminho de credencial no schema MCP;
- respostas determinísticas para entradas e fixtures equivalentes;
- ferramentas MCP anteriores permanecem retrocompatíveis;
- `bateia mcp` montado por L4 consegue anunciar e chamar
  `bateia_recommend`.

## Regime Tekt previsto

Aplicar o protocolo completo. Autor da intenção, autor do contrato, autor dos
oráculos, implementador, adversário e verificador devem registrar entradas,
capacidades e artefatos. O contrato e os testes A/B são congelados antes do patch
candidato. Casos opacos esperam `INCONCLUSIVE`; somente mutações explicitamente
opacas podem receber `Unknown` no gate discriminatório.

Se não houver isolamento material entre autores, a execução pode produzir um
ensaio reproduzível, mas o laudo deve dizer
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO` e não alegar selo independente.

## Fora de escopo

- incorporar um modelo de IA ao processo do Bateia;
- permitir que IA determine estados de evidência;
- execução assíncrona, cancelamento persistente e fila de jobs;
- `run_id`, `bateia_get_run`, `bateia_get_artifact` e paginação;
- modos MCP `local`, `replay` e `resume`;
- conectores para app stores, SaaS, registries ou datasets regionais;
- tradução automática dentro do Bateia;
- download, instalação, fork ou modificação de candidatos;
- autenticação multiusuário e transporte MCP por rede.

## Como pode falhar

- o contrato MCP vira apenas um wrapper de arquivos da CLI;
- a resposta é grande demais para uso prático por agentes;
- validação diverge entre CLI e MCP;
- a ferramenta omite proveniência para economizar contexto;
- a composição do serviço duplica regras existentes;
- a fiação real não injeta os provedores usados nos testes;
- conteúdo não confiável atravessa a fronteira como instrução;
- falhas parciais desaparecem no resumo.

## Critério de conclusão

O passo termina quando um cliente MCP, usando somente o schema anunciado e
provedores falsos offline, consegue submeter uma ideia completa e receber a mesma
decisão e o mesmo ranking produzidos pelos casos de uso de L1; e quando a fiação
real expõe essa ferramenta sem criar dependência de IA, I/O em L1 ou autoridade
nova sobre evidências.

O veredito positivo é `MCP_RECOMMENDATION_GATEWAY_READY`. Ele atesta somente o
modo síncrono `live` e não implica persistência, cobertura geral de mercados ou
qualidade de tradução multilíngue.
