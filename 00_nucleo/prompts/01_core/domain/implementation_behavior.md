# Passo 0057 — Observações sintáticas

ImplementationBehavior contém Source ImplementationReference, EndLine int, Status, Reason string, Facts []BehaviorFact e Limitations []string. BehaviorFact contém Kind, Target, Snippet string e Line, EndLine int. Todos campos JSON snake_case presentes, incluindo listas vazias e Target vazio. Nenhuma afirmação de execução ou satisfação.
