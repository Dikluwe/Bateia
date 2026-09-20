# Avaliação independente B — passo 0045

Regime A/B executado sem atestação de isolamento físico. Autor B não leu a implementação candidata, diffs ou corpos dos testes existentes de recuperação. Uma busca inicial de assinaturas com glob excessivo retornou nomes e algumas linhas de testes antigos não relacionados; ocorrência comunicada ao coordenador. Não houve uso deliberado desse material para os casos.

Corpus quantitativo e esperados congelados antes da candidata, com hashes em manifest.json. Baseline executado contra fonte de worktree 0044 antes da liberação do parser. Casos adicionais de chamadas e referências foram congelados após liberação da implementação, ainda sem acesso a ela. O congelamento pré-candidata se aplica ao corpus quantitativo, não à totalidade dos testes. Não houve alteração de fixtures ou expectativas após execução candidata.

Escopo: três projetos sintéticos independentes, seis consultas (três funcionais PT, um símbolo exato, ausência e sintaxe inválida). Por consulta, orçamento de nove hits e primeiros três caminhos distintos. Denominadores agregados por microcontagem; negativos não entram no denominador de recall. Precisão e recall sem denominador são null. Os resultados não constituem benchmark global nem evidência sobre projetos reais.

| Execução | Esperados recuperados | Retornados | Recall | Precisão |
| --- | --- | --- | --- | --- |
| Baseline0044 | 3/4 | 10 | 75% | 30% |
| Candidata r1 | 4/4 | 4 | 100% | 100% |

Ganho observado de 25 pontos percentuais de recall e 70 de precisão, sem perda de recall. Controles ausência/parser inválido e prioridade de símbolo exato preservados. A comparação baseline contra si mesma foi rejeitada pelo runner, como esperado por ausência de ganho. Logs e JSON: baseline.log/json, candidate-r1-tests.log, candidate.json, comparison.json; comandos de captura/comparação em compare.py.

Os seis testes parser independentes passam: métricas/negativos/exato, invariância a ordenação/duplicação, prioridade operacional, helper sem termos e profundidade/limites, resolução insegura e diversificação por linha. A primeira execução de validação de referências falhou em um de sete casos: hash e arquivo corretos com referência na linha1 (`package q`), fora da função chamadora, preservaram hit exclusivamente direct_call. Falha pública mínima transmitida ao coordenador; nenhuma expectativa alterada. Aceitação de r1: ganho quantitativo confirmado, contrato de referência ainda não aceito.

Colisão de artefato: B escreveu baseline.json métrico sem verificar existência e sobrescreveu registro de identidade do coordenador. Coordenador informado; identidade de produção será recuperada em contract-baseline.json a partir do registro0044 anterior, sem alegar novo congelamento prévio. Essa ocorrência é limitação de custódia, não correção de métrica.

## Reexecução r2 e aceitação limitada

Após correção de produção pelo coordenador, os mesmos testes passam integralmente: seis funções de teste parser e sete subcasos de referências. O contraexemplo r1 agora é rejeitado; caso válido e sobrevivência lexical continuam passando. Métricas r2 permanecem 4/4 recall e 4/4 precisão (100%/100%). Hashes congelados de corpus, runner e ambos arquivos de testes foram revalidados sem mudança.

Aceitação B: **aprovado no recorte sintético congelado e nos controles independentes executados**. A melhoria agregada atende ao critério sem perda de recall e os negativos/exato foram preservados. Não há atestação física, equivalência semântica nem afirmação de ganho em repositórios reais. Suíte completa, lint, build e consultas reais são responsabilidade do coordenador e não foram executados por B.

Comando r2: `GOCACHE=/tmp/bateia-0045-cache go test ./01_core/parser ./01_core/usecase -run '^TestFunctional0045' -count=1 -v` (saída candidate-tests.log, código0). Uma rodada de correção de produção, ganho discriminatório: um caso rejeitado em r1 passa em r2; nenhuma revisão do contrato/corpus/testes.
