# Improvement synthesis wiring

**Camada:** L4  
**Arquivo proprietário:** `04_wiring/improvement_synthesis.go`

Componha extrator, catálogo fechado, serviço L1, dispatcher e tool L2 do contrato
0028. Injete EXATAMENTE as mesmas instâncias certificadas de store, fila, gerador
de run ID, raiz, cursor secret, quota, workers, recovery e shutdown dos Passos
0026–0027. O dispatcher único encaminha `IMPROVEMENT_SYNTHESIS` sem criar worker,
store, índice de idempotência ou configuração paralelos.

Carregue o catálogo empacotado ou injetado como bytes imutáveis cujo SHA-256
confere com o contrato; falha impede inicialização. Apenas composição: nenhuma
regra, validação, hash, resolução de source, ranking, status, parsing de conteúdo
ou política em L4. O pacote multifile deve compilar com
`go build -o /tmp/bateia-step0028 ./04_wiring`.
