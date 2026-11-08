# 0127 — Separar resolução de fontes da síntese

**Data:** 2026-10-28
**Capacidade:** resolução íntegra de fontes congeladas para síntese, A1-INV sob Tekt.
**Responsáveis:** regras e orquestração por portas em L1; leitura concreta pelo store L3; lifecycle continua no runtime existente.

## Problema e observável

A análise de hipóteses já foi separada no passo 0103, mas o protocolo durável ainda implementa leitura paginada, verificação de identidade/hash e interpretação dos artefatos de origem como métodos de RecommendationRunService. Essa responsabilidade precisa somente de Get e ReadArtifact; receber o serviço completo acopla investigação a recursos de mutação e execução desnecessários.

Extrair essa resolução para improvement_sources.go, em L1, com porta interna de leitura que exponha somente Get e ReadArtifact. O protocolo chama o resolvedor na admissão e execução. A análise recebe o mesmo estado resolvido. Não criar serviço, pacote ou runtime paralelo.

## Preservação

Manter validação de tipo de run, terminalidade, índice publicado, fallback entre Hash e ArtifactHash, hash causal, paginação de 50 entradas e teto de 1000 páginas. Preservar falhas públicas, cópia dos bytes, detecção de cursor inválido, mecanismos únicos, lacunas e referências observadas. Fontes terminais sem sucesso continuam degradadas; essa marca não promove hipótese. Preservar normalização, idempotência, fila, cancelamento, transições, formatos e ordem dos cinco artefatos.

## Verificação

Escrever e executar caracterização antes da extração: projeção resolvida e saída da análise congeladas por digest; falhas de fonte e paginação com portas em memória. Reexecutar os mesmos testes após extração e manter contratos de integração do módulo raiz. Executar lint, testes, vet e build; registrar bloqueios preexistentes sem reescrever seus resultados históricos.

## Limites

A separação é de dependência e responsabilidade dentro do pacote usecase. Não remove relógio de todo L1, não extrai persistência da CLI, não generaliza o dispatcher e não conclui a avaliação da hierarquia A0–A2. Mudanças locais 0126 e artefatos de outras investigações permanecem preservados. Resultados brutos em lab/bancada/arquitetura/0127.
