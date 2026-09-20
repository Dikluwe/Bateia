# Matriz discriminatória — passo 0072

| Classe | Casos | Resultado exigido |
|---|---|---|
| Positivo | `span`, `bits/basic_string.h`, duas quotas justas, arquivo compartilhado, Go pequeno, árvore >20.000 | caminhos dirigidos examinados; deduplicação; legado preservado; enumeração completa |
| Negativo | `spanner`, `span.txt`, `basic_string.md`, `NOTICE`, conteúdo sem símbolo | não candidato/não coletado; nunca `PROVEN`, `PARTIAL` ou componente por nome |
| Opaco | bytes esgotados, falha de árvore, regra ausente/hipótese | `COLLECTION_INCOMPLETE` ou limitação de catálogo; nunca ausência/sucesso implícito |
| Determinismo | 20 permutações de árvore, consultas, termos, catálogo e requisitos | JSON canônico idêntico e entradas imutáveis |

## Mutantes

| ID | Predicado que o elimina |
|---|---|
| M01 | sucesso local inclui entrada posterior à antiga fronteira 20.000 |
| M02 | primeiros caminhos solicitados pertencem às reservas dirigidas |
| M03 | `NOTICE` extensionless não é candidato nem código |
| M04 | `span` extensionless é candidato e é examinado |
| M05 | `spanner` não aparece em candidatos de `span` |
| M06 | ambos os requisitos obtêm uma leitura sob quota unitária |
| M07 | conteúdo irrelevante em `span` não promove status/componente |
| M08 | bytes insuficientes resultam em `COLLECTION_INCOMPLETE` |
| M09 | `span` omitido aparece em `skipped_paths` |
| M10 | coleta incompleta não aparece em `CandidateAssessment.Gaps` |
| M11 | vinte permutações produzem observação byte a byte idêntica |
| M12 | alvo compartilhado ocorre uma vez na requisição, itens e cobrança |
| M13 | lacuna e hipótese não geram consulta declarada |
| M14 | `graph.go` continua coletado sem consultas dirigidas |

Mutation score requerido: `14/14 = 1.0`. O resultado pré-candidato não calcula score:
ele apenas confirma que o gate discrimina o baseline ainda não materializado.
