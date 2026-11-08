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

## Admissão por registro no passo 0104

Delegar a admissão não vazia e a validação de capacidades/IDs ao construtor L1 NewMarketDiscoveryRegistry. Não impor cardinalidade fixa em L4. Aceitar conjuntos válidos unitários ou ampliados, mantendo a mesma configuração e lifecycle; recusar configuração inválida antes de criar recursos.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

0131: admitir catálogo com NewConfiguredMechanismExtractor; não executar ExtractMechanisms com entrada fictícia. Compor executor de mercado explícito e preservar a ordem de lifecycle.
