# Laudo 0005 — Triagem determinística por evidências

**Data:** 2026-08-24  
**Resultado:** implementação concluída; calibração manual com candidatos remotos permanece como atividade observacional.

## Entregue

- contratos de plano, bundle congelado, fontes, estados, avaliações e componentes;
- catálogo determinístico versionado com nove capacidades iniciais;
- seleção reproduzível de `ADVANCE`, `AUDIT` e amostra de `DROP`;
- orçamento de documentos, manifestos, código, testes, bytes e profundidade;
- coleta deduplicada por hash e congelada em commit imutável;
- resolução de commit remoto pelo GitHub e local por `git rev-parse HEAD`;
- avaliação offline com explicações por templates;
- vetores separados para capacidade, restrição, licença, manutenção e operação;
- comandos `idea evidence plan`, `collect` e `assess`;
- composição que só inicializa o provedor para `collect`.

O pipeline não contém LLM, embeddings, prompts de classificação, reranker neural
ou chamada a serviço de IA.

## Controles automáticos

- documentação isolada não produz `PROVEN`;
- símbolo implementado e teste relacionado produzem prova forte;
- falha de coleta produz `INCONCLUSIVE`, nunca `ABSENT`;
- ruído em `vendor` não produz `PROVEN`;
- ordem das entradas não altera a avaliação;
- conteúdo duplicado é deduplicado por hash;
- plano aceita os snapshots existentes com `HEAD` pendente, mas a coleta só começa
  depois de resolver um commit imutável;
- CLI executa planejamento e avaliação inteiramente a partir de JSON local.

## Verificação

```text
go test -race ./...  PASS
go vet ./...         PASS
crystalline-lint .   ✓ No violations found
go build              PASS
git diff --check      PASS
```

O planejamento foi executado contra
`structural-impact.screening.json` e produziu candidatos limitados com orçamento
explícito, incluindo Enola. Como a revisão humana dos repositórios remotos não faz
parte dos testes determinísticos, `evidence_precision`, `false_absence_rate` e
`source_trace_rate` ainda não devem ser publicados como métricas observadas. Essa
amostra deve ser coletada e julgada sem alterar as regras durante o replay.

## Limites conhecidos

- o catálogo inicial reconhece apenas as nove famílias previstas no passo;
- metadados temporais de manutenção ainda não fazem parte do bundle;
- reconhecimento de símbolos usa declarações lexicais conservadoras, não AST para
  todas as linguagens;
- um repositório local precisa estar dentro de uma worktree Git para fornecer uma
  referência imutável real.
