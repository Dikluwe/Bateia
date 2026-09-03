# Prompt R2 — composição, arquitetura e gate 3×3

## Composição observável

A fábrica exportada recebe um único executor que implementa `Recommend` e `RecommendObserved`. O teste invoca as vias síncrona e assíncrona e prova que o mesmo objeto recebeu ambas. Dezesseis starts verificam formato, unicidade e ausência de correlação com cada campo textual controlado pelo usuário; isso não promove a estimativa de entropia opaca O04.

## Auditoria estrutural

AST/imports percorrem todos os arquivos Go de produção das camadas, rejeitam I/O e dependências ascendentes em L1, infraestrutura em L2 e política de transição/hash em L4. A fábrica deve executar recovery antes de iniciar workers.

## Gate offline

O antigo gate R2 baseado em presença de bindings foi removido integralmente. O gate executável independente R3 vive em `04_wiring/r3_gate`, executa fixtures/mutantes/faults concretos e conserva o protocolo 3×3.
