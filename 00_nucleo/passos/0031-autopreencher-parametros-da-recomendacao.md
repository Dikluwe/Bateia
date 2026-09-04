# Passo 0031 — Autopreencher parâmetros da recomendação MCP

**Estado**: `EM MATERIALIZAÇÃO`
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Corrige**: atrito de chamada observado na execução do Passo 0030
**Baseline obrigatório**: implementação e testes existentes antes deste passo
**Camadas previstas**: L0, L1 e L2; L3/L4 apenas para regressão
**Natureza**: refinamento semântico retrocompatível para entradas completas

## Problema

`bateia_recommend` e `bateia_start_recommendation` exigem hoje que a IA monte
integralmente `policy`, `target` e `search`. Isso obriga o chamador a decidir
parâmetros operacionais que pertencem ao produto e aumenta a chance de uma
chamada inválida antes de o Bateia conseguir auxiliar a busca.

O Bateia deve aceitar a intenção mínima de negócio e materializar defaults
seguros e observáveis. Autopreencher a política de licenças significa escolher
uma **allowlist de licenças aceitas**, nunca atribuir ou presumir uma licença a
um candidato cuja licença é desconhecida.

## Entrada mínima

Permanece obrigatório informar:

```json
{
  "idea": {
    "problem": "problema a resolver",
    "capabilities": ["capacidade necessária"]
  }
}
```

O Bateia não inventa o problema nem as capacidades desejadas. Os objetos
`policy`, `target` e `search`, e cada campo interno deles, tornam-se opcionais.

## Defaults versionados

Na ausência do campo correspondente, o perfil `BATEIA_RECOMMEND_DEFAULTS_V1`
materializa:

```json
{
  "policy": {
    "schema_version": "1",
    "allowed_licenses": [
      "0BSD",
      "Apache-2.0",
      "BSD-2-Clause",
      "BSD-3-Clause",
      "BSL-1.0",
      "ISC",
      "MIT",
      "Unlicense",
      "Zlib"
    ],
    "require_local": false,
    "allow_archived": false,
    "max_components": 5,
    "extend_min_coverage": 0.5,
    "integration_min_value": 0.5
  },
  "target": {
    "name": "<idea.problem>"
  },
  "search": {
    "strategy": "lexical",
    "language": "",
    "locales": ["pt-BR", "en-US"],
    "max_per_query": 5,
    "max_candidates": 8
  }
}
```

A allowlist é ampla entre licenças SPDX permissivas usuais, mas não reivindica
ser juridicamente exaustiva. Licenças copyleft e licenças de fonte disponível
com restrições adicionais ficam fora do perfil V1. Alterá-la exige nova versão
do perfil e novo contrato.

## Regras de precedência

1. Campo omitido recebe o default V1 correspondente.
2. Campo informado pelo chamador, inclusive `false`, `0`, string vazia ou lista
   vazia, não é substituído silenciosamente.
3. Valor explícito inválido continua produzindo erro de validação.
4. `allowed_licenses` omitido recebe a allowlist V1; `allowed_licenses: []`
   permanece lista vazia explícita e não recebe default.
5. Campos desconhecidos continuam rejeitados.
6. Uma requisição completa já válida preserva o mesmo significado.
7. Os caminhos síncrono e assíncrono materializam bytes semanticamente
   equivalentes antes de chamar o caso de uso.

## Observabilidade

O schema MCP anuncia os defaults como anotações, mas o servidor — não o cliente
nem o schema — é a autoridade que os aplica. O resultado continua expondo a
política efetivamente usada, permitindo distinguir a decisão materializada.

## Refinamento do contrato anterior

Este passo substitui somente a obrigatoriedade exata congelada no Passo 0025:

- no schema de `bateia_recommend`, apenas `idea` permanece obrigatória na raiz;
- no schema de `bateia_start_recommendation`, apenas `idempotency_key` e `idea`
  permanecem obrigatórios na raiz;
- a forma e a semântica de entradas completas válidas são preservadas;
- o teste histórico de schema exato deve ser refinado de modo explícito, sem
  apagar a evidência de que o contrato anterior era mais estrito.

## Ataques obrigatórios

O oráculo independente deve discriminar pelo menos:

- requisição mínima versus requisição completa equivalente;
- omissão de cada objeto e de cada campo separadamente;
- lista de licenças omitida versus lista explicitamente vazia;
- booleano omitido versus `false` explícito;
- zero explícito inválido versus inteiro omitido;
- campo desconhecido;
- paridade síncrona/assíncrona;
- mutações da allowlist, dos limites e dos thresholds;
- falsa atribuição de licença a candidato com licença desconhecida.

## Fora de escopo

- inferir `idea.problem` ou `idea.capabilities` de campos arbitrários na raiz;
- aceitar o payload plano usado incorretamente no Passo 0030;
- escolher ou modificar a licença de um repositório;
- aconselhamento jurídico sobre compatibilidade de licenças;
- afirmar que uma licença permitida existe sem evidência coletada;
- mudar descoberta, scoring ou critérios de evidência.

## Critério de conclusão

O passo conclui quando contrato e oráculos segregados estiverem selados, a
implementação aplicar os defaults nos dois caminhos MCP, ataques relevantes
forem mortos, testes focais e regressões passarem e um verificador independente
emitir veredito reproduzível.
