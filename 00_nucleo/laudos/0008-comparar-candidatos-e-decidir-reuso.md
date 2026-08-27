# Laudo 0008 — Comparar candidatos e decidir reuso

**Data**: 2026-08-26  
**Veredito**: `IMPLEMENTADO`  
**Regime**: A/B limitado, executado sem atestação de isolamento

## Resultado

O Bateia agora transforma `IdeaSpec`, `AssessmentReport` e política versionada
em uma das decisões `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE`,
`BUILD` ou `INCONCLUSIVE`, sem LLM, embeddings ou acesso à rede. O relatório
expõe matriz por capacidade/candidato, política e hash, seleção, componentes,
lacunas, bloqueios, razões com IDs de regra, alternativas e esforço observável.

`BUILD` somente é emitido depois de ausência conclusiva. Falha de execução,
capacidade obrigatória incompleta ou gate ambíguo produz `INCONCLUSIVE`.
Licença, manutenção, operação e conflitos ficam fora da cobertura funcional.

## Materialização

- domínio: `01_core/domain/reuse_decision.go`;
- motor: `01_core/usecase/reuse_decision.go`;
- CLI: `bateia idea decide --idea --assessment --policy [--output]`;
- testes adversariais: `01_core/usecase/reuse_decision_test.go`;
- replay de produto: `01_core/usecase/reuse_decision_product_test.go`;
- teste offline da boca de inspeção: `02_shell/cli/idea_decide_test.go`.

Cada novo consumidor Go possui prompt L0 próprio, inclusive os arquivos de
teste, conforme o contrato V15 do Tekt.

## Ataques e gates

Os testes cobrem conflito contra `FORK`, falha/inconclusão contra `BUILD`, sete
saídas, ordem invariável, irrelevância do score de screening, rastreabilidade,
mudança isolada de política, imutabilidade da entrada e CLI sem rede.

Replay congelado:

| Ideia | Julgamento | Resultado |
|---|---:|---:|
| busca local | `FORK` | `FORK` |
| analisador estrutural | `EXTEND` | `EXTEND` |
| pipeline de eventos | `COMPOSE` | `COMPOSE` |
| editor de receitas | `INTEGRATE` | `INTEGRATE` |
| simulador inédito | `BUILD` | `BUILD` |

Utilidade registrada: `5/5 = 100%` (gate: ≥ 80%).

Comandos reproduzíveis:

```bash
GOCACHE=/tmp/bateia-go-cache go test ./...
/repos/Antigravity/tekt-linter/target/debug/crystalline-lint .
git diff --check
```

## Papéis e limitação de independência

- autoria da obrigação: Passo 0008 previamente congelado;
- autoria do contrato e implementação: sessão atual;
- autoria dos ataques: sessão atual, escrita contra o contrato público;
- executor e veredito: sessão atual.

Houve separação por artefatos, ordem contrato → ataques → implementação → gates,
mas não houve isolamento entre autores, executor e árbitro. Portanto este laudo
não alega verificação independente nem materialização segregada atestada.
