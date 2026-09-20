# Passo 0042 — Analisar e sanear promessas, evidência e limites do Bateia

Estado: CONCLUÍDO — diagnóstico e saneamento documental; correções de produção pendentes. Natureza: diagnóstico transversal de engenharia, com saneamento documental de afirmações contraditas por evidência. Não é certificação independente.

## Pergunta

Onde a intenção do produto foi implementada e validada, onde há somente aproximação ou contrato, e onde a saída induz confiança maior que a evidência? A hipótese de fragilidade deve ser examinada por subsistema; não será aplicada indiscriminadamente ao projeto nem atribuída à capacidade de seu autor.

## Escopo pré-declarado

Inventariar produção e documentação pública. Examinar amostras dirigidas dos fluxos: promessa/uso, descoberta e expansão, pontuação, coleta e suporte linguístico, comprovação, decisão/alternativas, autofagia/Katamari, mecanismos/hipóteses, validações/benchmarks, rastreabilidade/replay/persistência e arquitetura/lint. Registrar arquivos e intervalos efetivamente lidos; inventário não significa revisão integral de cada arquivo.

## Classes de conclusão

- Fato observado: código, trecho ou execução identificados.
- Inferência: consequência fundamentada, ainda não medida diretamente.
- Hipótese: suspeita que exige experimento adicional.
- Não verificado: escopo não exercitado, sem nota de aprovação.

Estados de produto separados: especificado, implementado, conectado à jornada, exercitado em fixture, exercitado em caso local real, validado externamente. Aprovação de testes, hashes e nomes de gates não promove automaticamente o último estado.

## Execução

1. Congelar identidade do checkout e inventário de caminhos; preservar alterações do 0041.
2. Registrar matriz promessa → mecanismo → evidência → limite por área.
3. Reproduzir contraexemplos pequenos e controles, sem rede ou execução de candidatos; não acessar oráculos protegidos ou reescrever resultados históricos.
4. Sanear documentação atual quando a contradição for direta; preservar artefatos históricos com qualificação explícita.
5. Entregar achados priorizados e sequência de correções, separando limitação honesta de defeito confirmado. Não reengenheirar silenciosamente todos os subsistemas neste diagnóstico.

## Critério de conclusão

Cada área tem achado sustentado ou marcação de não verificado. Toda conclusão forte possui referência; comandos e limites estão registrados. Documentação corrigida não pode alegar correção de comportamento inexistente. Produto, segurança operacional, descoberta pública, múltiplas linguagens e independência de benchmarks não serão certificados por uma amostra local.

## Resultado

[Relatório](../relatorios/sanitizacao-promessas-0042.md) e [evidências reproduzíveis](../../lab/bancada/history/benchmarks/0042/README.md). Contraexemplos confirmados em negação, autoridade de decisão, integração de mecanismos, contenção local e métricas heurísticas. README e índice de relatórios saneados; hashes de produção permaneceram iguais durante este diagnóstico. Probe offline, lint estrutural e verificação de whitespace concluídos. Não houve nova certificação do produto nem execução de oráculos.
