# Resultado0058 — Perda da implementação na consulta local

O passo0058 corrigiu a ligação entre o planejador existente e a interpretação da consulta local. Não adicionou vocabulário ou alterou pesos. Expansões de ação/objeto são registradas como hipóteses, com origem e versão; não demonstram satisfação do requisito. Consultas negadas, ambíguas, com sintaxe explícita ou modos mixed/symbol continuam fora dessa expansão.

## Evidência observada
A pergunta bitset-p1 não recuperava Intersection mesmo com limite10000:55hits. O plano existente produzia termos que a colocavam em7º; consulta pelo símbolo a colocava em1º. Isso localizou a perda na representação da pergunta antes do ranking, sem indicar falha do parser ou corte no top24. Após a correção, a pergunta original a recupera em7º entre87hits. Fonte: benchmark/0058/diagnostic-before.json e diagnostic-after.json; bitset.go:950–959, SHA2563db62f1020df04368bdcaa16be403aa1410900bfd688a602aa9c63fb851d7686.

## Comparação no mesmo corpus

| Caso | Função esperada | Antes | Depois |
| --- | --- | --- | --- |
| fsnotify-p1 | Add | não recuperada | não recuperada |
| fsnotify-p2 | Remove | não recuperada | não recuperada |
| colorful-p1 | BlendLab | não recuperada | não recuperada |
| colorful-p2 | LinearRgbToXyz | recuperada | recuperada |
| bitset-p1 | Intersection | não recuperada | recuperada |
| bitset-p2 | Union | não recuperada | recuperada |

Recuperação, confirmação de identidade e corpo disponível:1/6→3/6. Arquivos encontrados:6/6; extraídos:5/6, sem mudança. Nos três negativos:0 indicações e0 conclusões indevidas, antes/depois. Fonte: benchmark/0058/comparison.json e collected/summary.json. Os três casos positivos restantes são lacunas observadas; a causa de cada perda requer diagnóstico próprio.

## Validação e limites
Build, testes independentes TestBridge0058, suíte Go completa, nove controles do runner, crystalline-lint e diff-check passaram. Replay das observações reproduziu o resumo; hashes do corpus e candidato permaneceram iguais aos congelados. Recibos e logs estão em benchmark/0058.

A/B sem atestação física: A e B separados, testes congelados antes da execução. A teve leitura incidental do parágrafo diagnóstico no contrato, registrada em verification.json; não acessou corpus/testes. B não leu produção candidata. Integrador conhece os alvos. Portanto, ganho de calibração, não validação holdout.

Não verificado: comportamento funcional, preservação das entradas e descoberta pública. Encontrar Intersection/Union não prova que resolvem integralmente a pergunta. Tempos registrados são telemetria da ferramenta. A próxima investigação deve rastrear um dos três alvos ainda perdidos, antes de propor novos mecanismos.
