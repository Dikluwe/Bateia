# Laudo 0020 — Busca autenticada com orçamento observável

**Data**: 2026-08-31
**Estado do passo**: `IMPLEMENTADO`
**Veredito externo**: `INCONCLUSIVE_RATE_LIMIT`
**Regime Tekt**: `EXECUTADO_SEM_ATESTACAO_DE_ISOLAMENTO`

## Materializado

- operação separada `idea web authenticated-experiment`;
- autorização redigida e token recebido somente pela composição L4;
- orçamento congelado de 36 queries, 10 resultados, timeout por tentativa e total;
- estados fechados para sucesso, vazio, cota, timeout, resposta opaca, privado e budget;
- visibilidade pública comprovada antes de persistir hits;
- contadores de rate limit reduzidos a inteiros permitidos;
- request, response e snapshot selados para replay offline;
- varredura de credencial antes da persistência;
- receipt com contagens e causa de parada.

## Execução externa

O token estava presente, mas nunca foi impresso, serializado ou incluído em hash.
A primeira query retornou resposta válida vazia e cota restante dentro da margem
segura declarada. O orquestrador encerrou novas chamadas e materializou todas as
posições restantes como `SKIPPED_BUDGET`.

| Observável | Resultado |
|---|---:|
| queries planejadas | 36 |
| requests iniciadas | 1 |
| `EMPTY` | 1 |
| `SKIPPED_BUDGET` | 35 |
| retries | 0 |
| snapshots | 3 |
| credenciais encontradas | 0 |

Distribuição por braço:

| Braço | `EMPTY` | `SKIPPED_BUDGET` |
|---|---:|---:|
| N | 1 | 11 |
| F1 | 0 | 12 |
| F2 | 0 | 12 |

## Reprodutibilidade

Cada snapshot contém doze posições. Dois replays offline por braço foram
idênticos byte a byte e marcaram `replayable=true`. Hashes:

```text
N  2357943aa87f9360ed37396d33bf2d668b474472fd2a843517abd9b87bc21ac7
F1 380fd848c1fbcb11830fc9a354a33ed05e8df7910126051b9d007b640ab3d48e
F2 7bae0700bedd0711d4c39f7c2b1b424902091163bca057d76bd57eff021c596c
```

## Limite Tekt

Contrato e ataques foram escritos por autoridades lógicas separadas. O gate
pré-implementação permaneceu bloqueado em `4/5`: uma obrigação já satisfeita pelo
baseline foi contada como mutação sobrevivente. A implementação não altera esse
registro nem reivindica selo independente. O resultado é, portanto, execução com
testes independentes, mas sem atestação completa de isolamento.

## Veredito

Segurança, orçamento e replay funcionaram. O gate de exposição não passou: cada
braço precisava concluir ao menos 80%, e somente uma de 36 queries concluiu. Não
é permitido calcular precisão, recall ou escolher gramática com esse run.

O próximo run válido deve ocorrer em nova janela de cota, mantendo queries,
ordem, orçamento e margem. Este resultado não deve ser apagado nem tratado como
ausência de candidatos.

## Retomada posterior — run-5

Foi adicionada retomada segura por snapshots: respostas `SUCCESS`/`EMPTY` são
preservadas e somente casos opacos ou pulados são repetidos. Também foi corrigido
um falso positivo do filtro de segredo que rejeitava descrições públicas contendo
as palavras comuns `authorization` ou `cookie`.

A cadeia `run-2 → run-3 → run-4 → run-5` concluiu 36/36 queries, 12/12 em cada
braço, sem `UNKNOWN` ou `SKIPPED_BUDGET`. O token exato não aparece nos artefatos
e o replay offline foi idêntico por braço. O veredito atualizado da descoberta é
`READY_FOR_BLINDED_JUDGMENT`; os runs anteriores permanecem preservados.
