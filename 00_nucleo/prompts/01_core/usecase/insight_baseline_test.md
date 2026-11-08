# Testes de integridade e limites da baseline

**Camada:** L1. **Contrato:** passo 0135.

Consumir CompareInsightBaseline com a fixture pura de comparação. Verificar rejeição de schema/contexto, identidade de repositório, manifesto/hash/conteúdo/path/UTF-8/range e duplicatas inválidos, sem relatório parcial. Conferir fronteiras exatas de 64 arquivos, 256 KiB por arquivo e 1 MiB total e seus overflows. Sem I/O. Este arquivo acompanha o normalizador insight_baseline.go; os testes relacionais/determinismo ficam no consumidor insight_comparison_test.go.
