# D-002 — Alternativas de software com evidência

**Estado**: `DESEJO`
**Registrado em**: 2026-08-31
**Origem**: clarificação do produto — “algo parecido com AlternativeTo”
**Relacionado**: `D-001 — Da ideia à decisão de reúso`

## O desejo

Queremos informar um software conhecido, ou uma necessidade quando ainda não há
alvo, e receber uma lista comparável de alternativas. O gesto deve ser tão simples
quanto “alternativas para X”, mas o ranking do Bateia deve ser sustentado por
compatibilidade funcional comprovada, não apenas por votos, estrelas ou palavras.

## Diferença de produto

Serviços comunitários organizam alternativas por aplicação, plataforma, licença,
popularidade e opinião. O Bateia usa os mesmos eixos de decisão, mas acrescenta:

- capacidades comprovadas em documentação, código, testes e manifestos;
- identidade imutável `owner/repo@commit`;
- lacunas funcionais e bloqueios explícitos;
- operação local/self-hosted e licença como filtros verificáveis;
- `INCONCLUSIVE` quando a coleta não sustenta comparação.

Popularidade pode desempatar candidatos funcionalmente equivalentes, mas nunca
promover um candidato com menor cobertura ou evidência opaca.

## Jornada principal

```text
software alvo ou necessidade
  → perfil funcional (IdeaSpec)
  → descoberta de candidatos
  → comprovação de capacidades
  → filtros de licença, manutenção e operação
  → alternativas ordenadas com cobertura, lacunas e evidências
  → recomendação de substituição, extensão, integração ou referência
```

O próprio alvo nunca aparece como sua alternativa. Uma alternativa só é destacada
como primeira escolha quando passa pelos filtros e possui evidência suficiente.

## Saída mínima

Para cada candidato:

- repositório, commit e URL;
- porcentagem de capacidades comprovadas;
- capacidades atendidas e ausentes;
- licença, estado de manutenção e operação local;
- classificação `RECOMMENDED`, `PARTIAL`, `REFERENCE`, `BLOCKED`,
  `INCONCLUSIVE` ou `NOT_ALTERNATIVE`;
- razões, bloqueios e evidências rastreáveis.

O relatório separa alternativas elegíveis de candidatos excluídos e mantém ordem
determinística. Estrelas são apenas desempate final.

`alternatives` contém somente `RECOMMENDED`, `PARTIAL` e `REFERENCE`.
`excluded` preserva `BLOCKED`, `INCONCLUSIVE`, `NOT_ALTERNATIVE` e o próprio
alvo. `top_choice` nunca aponta para `REFERENCE` ou para `excluded`.

## Sinais de atendimento

1. “Alternativas para owner/projeto” exclui `owner/projeto` do resultado.
2. Um candidato com cobertura comprovada vence outro apenas popular.
3. Filtros de licença, arquivamento e operação local são explicáveis.
4. Evidência ou identidade opaca produz `INCONCLUSIVE`, não recomendação.
5. Reordenar os mesmos inputs não muda a saída.
6. O dossiê de `idea recommend` inclui `alternatives.json` e
   `alternatives.md` sem criar um segundo motor de descoberta.
