# Adendo0057 — Não promover origem documental

O gate legado0051 detectou que fatos AST do corpo passaram a aparecer para hit com Retrieval=documentation. Mesmo com arquivo extraído, a seleção documental não é seleção de implementação. Preservar a separação anterior.

InspectImplementationBehavior deve retornar not_evaluated, zero Facts e motivo explícito para hit Retrieval=documentation. Hits de outro canal com identidade válida continuam independentes e podem ser observados. Não modificar o teste legado. Esta é uma correção após primeira execução, que ficará preservada. Atualizar EvidenceVersion evidence-v19-function-behavior-origin e checkpoint implementation_search0057-r1 local/live para invalidar projeções iniciais. Repetir somente validação local, sem rede.
