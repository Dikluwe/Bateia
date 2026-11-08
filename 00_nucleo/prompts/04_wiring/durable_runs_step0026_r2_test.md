# Prompt R2 — composição, arquitetura e gate 3×3

## Composição observável

A fábrica exportada recebe um único executor que implementa `Recommend` e `RecommendObserved`. O teste invoca as vias síncrona e assíncrona e prova que o mesmo objeto recebeu ambas. Dezesseis starts verificam formato, unicidade e ausência de correlação com cada campo textual controlado pelo usuário; isso não promove a estimativa de entropia opaca O04.

## Auditoria estrutural

AST/imports percorrem todos os arquivos Go de produção das camadas, rejeitam I/O e dependências ascendentes em L1, infraestrutura em L2 e política de transição/hash em L4. A fábrica deve executar recovery antes de iniciar workers.

## Gate offline

O antigo gate R2 baseado em presença de bindings foi removido integralmente. O gate executável independente R3 vive em `04_wiring/r3_gate`, executa fixtures/mutantes/faults concretos e conserva o protocolo 3×3.

## Portas declarativas — 0129

Distinguir import de io/io/fs usado exclusivamente para declarar interfaces/tipos/constantes de operação concreta. Permitir esses dois imports somente em arquivos sem funções, variáveis inicializadas ou chamadas AST. Manter recusa de os, net/http, dependências ascendentes e qualquer arquivo executável que use io. Adicionar controles para interface declarativa, função e inicializador com efeito; não liberar import por nome de arquivo.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.

## Consolidação 0131
O auditor AST geral e seu controle de portas declarativas foram transferidos para lab/bancada/arquitetura/0131/contracts/legacy_architecture_test.go, preservando asserções e raízes. Este arquivo possui apenas testes de composição, publicação, privacidade e lifecycle.
