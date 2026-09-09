# Avaliação B — 0050

Resultado: **PASS no recorte testado**, seis testes principais e 21 subcasos; primeira execução, sem revisão de oráculos após observar a candidata. Os hashes de `tests-freeze.json` foram congelados antes da execução e permaneceram iguais depois. Log integral em `independent-tests.log`; metadados e SHA-256 do log em `independent-evaluation.json`.

Comando: `GOCACHE=/tmp/bateia-0048-cache go test ./03_infra/github -run TestStep0050 -count=1 -v`, saída 0. Execução com autorização de sandbox para httptest em loopback.

O positivo compara hits completos e fontes examinadas de dois arquivos inéditos entre clientes local/remoto, incluindo pesos, linhas e SHA-256; também exige hit para consulta funcional cujo caminho não repete a pergunta. Router encaminha local/remoto e mantém provedor legado not_evaluated. Seleção prioriza arquivo recuperado; releitura íntegra preserva hits e releitura divergente/ausente os rejeita.

Os adversos cobrem identidade antes de I/O, SHA/tamanho/encoding/base64/conteúdo inválidos, bytes decodificados descartados, árvores truncadas/opacas/duplicadas, caminho inseguro, symlink, limites menores, teto de tentativas, 403/429, cancelamento prévio e envelopes excessivos.

## Limitações

Regime Tekt A/B **executado sem atestação de isolamento**: filesystem/permissões compartilhados. B não leu implementação candidata GitHub, corpos do Router, ranking/parser ou fixtures antigas. Busca de assinaturas de API mostrou nomes de testes antigos, sem seus corpos. Tipos públicos, construtores e API de seleção/revalidação do usecase foram consultados. Paridade compartilha ranking existente e não demonstra pertinência geral.

Este harness não cobre cancelamento durante streaming, todas as variantes de caminhos inseguros, limites exatos dos envelopes, saturação real de 4 MiB ou 24 hits distintos. Tetos de bytes/resultados são verificados no budget efetivo declarado. O integrador realiza suíte completa, lint, build e ensaio público real separadamente. Não houve mutações calibradas; o resultado não constitui conformidade com protocolo completo ou capacidade geral comprovada.
