# Laudo 0001 — Materialização de `IdeaSpec → ExpansionPlan`

**Data**: 2026-08-24  
**Passo verificado**: [0001 — Materializar `IdeaSpec → ExpansionPlan`](../passos/0001-materializar-ideaspec-e-plano-de-expansao.md)  
**Resultado**: `APROVADO`

## O que foi implementado

- contrato de domínio L1 para `IdeaSpec`, `SearchExpansion` e `ExpansionPlan`;
- expansor puro e determinístico que normaliza dados, deduplica queries e mantém
  cobertura rastreável por campo da ideia;
- preservação explícita de problema, capacidades, restrições, não objetivos e
  stack preferida, sem transformar `NonGoals` ou stack em identidade positiva;
- comando `bateia idea plan --input ... [--output ...]` em L2;
- despacho do comando antes da leitura de `.env` e da construção de cache e
  clientes remotos em L4;
- prompts L0 e linhagem cristalina para os componentes novos;
- fixture e golden file que tornam o formato público deliberado.

## Evidências

### Comportamento do núcleo

Os testes cobrem:

- estabilidade byte a byte para a mesma entrada;
- ausência de mutação da entrada;
- normalização e deduplicação por caixa e espaços;
- rejeição de ideia vazia, composta somente por stack ou somente por não objetivo;
- não promoção de `NonGoals`;
- preservação conjunta de `Covers` quando problema e capacidade geram a mesma
  query;
- separação entre a ideia-alvo e uma query externa de descoberta.

### Contrato da CLI

O golden test compara integralmente a saída JSON. Testes adicionais verificam
gravação por `--output`, `stdout` vazio nesse modo, JSON inválido, ideia vazia e
ausência de `--input`.

### Execução offline

O binário foi construído e executado em `/tmp`, sem variáveis de ambiente:

```bash
env -i /tmp/bateia-step-0001 idea plan \
  --input /home/dikluwe/Documentos/Antigravity/Bateia/02_shell/cli/testdata/idea_spec.json
```

Resultado: código zero e somente o `ExpansionPlan` em `stdout`. Não houve leitura
visível de `.env`, inicialização de Bleve/cache, solicitação de token ou acesso à
rede.

### Verificações automatizadas

```text
go test ./...          PASS
go test -race ./...    PASS
go vet ./...           PASS
crystalline-lint .     PASS — No violations found
git diff --check       PASS
```

A primeira tentativa de `go test ./...` dentro do sandbox falhou apenas porque
ele proibia portas de loopback usadas por `httptest` nos adaptadores GitHub e
Sourcegraph. A mesma suíte passou integralmente com essa permissão habilitada.

## Medida observável

Para a fixture de referência, uma ideia produz seis expansões estáveis: uma de
problema, três de capacidades e duas de restrições. Cada expansão possui ID,
faceta, justificativa e origem em `Covers`. A stack e o não objetivo permanecem
na ideia, mas não geram buscas positivas.

## O que permanece aberto

- interpretar ideia em texto livre e preencher `IdeaSpec`;
- gerar sinônimos ou expansões semânticas;
- executar o plano em provedores de descoberta;
- agregar resultados e medir diversidade/precisão;
- comparar candidatos com a ideia-alvo e recomendar reuso ou construção.

O próximo experimento pode agora avaliar estratégias de expansão semântica sem
voltar a misturar a definição da ideia com as queries externas.
