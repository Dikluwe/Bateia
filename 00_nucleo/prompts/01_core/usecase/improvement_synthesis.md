# Improvement synthesis use case

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/improvement_synthesis.go`

Implemente validação, normalização, resolução pré-fila pelas portas duráveis já
certificadas, hash causal, idempotência, execução e projeção do contrato 0028.
Reutilize o mesmo store, fila, run ID, cursor, quota, CAS, cancelamento, recovery e
publicação do Passo 0026. Não crie runtime paralelo.

Hipóteses por regra usam somente taxonomia congelada; propostas externas são dados
não confiáveis e nunca criam evidência. Aplique tetos por mecanismo, alvo,
licença, contraevidência e compatibilidade; lifecycle `SUCCEEDED` nunca implica
`SUPPORTED`. Preserve buscas limitadas, riscos, premissas, unknowns e planos
falsificáveis. Canonicalize e publique os cinco artefatos na ordem contratada.
Zero I/O concreto, MCP, filesystem, rede, env, relógio, goroutine ou mutex.

## Separação de responsabilidades no passo 0103

Este arquivo passa a possuir o protocolo de execução de síntese: admissão, resolução de fontes pelas portas, transições, cancelamento e publicação. Delegar a derivação e o estado de negócio ao analisador puro de `improvement_analysis.go`, transferindo apenas os mecanismos validados e sua degradação. Preservar ordem das transições, cinco artefatos, envelopes, hashes e defaults públicos. O runtime não promove hipótese ao registrar sucesso operacional.
