# Laudo 0006 — Endurecimento e primeira calibração real

**Data:** 2026-08-24
**Estado:** implementação e gate concluídos pelo replay do Passo 0007

## Correções executadas

- identidade `owner/repo` separada da revisão SHA;
- preservação da decisão da peneira depois da resolução do commit;
- estados terminais por alvo de coleta;
- completude derivada de árvore, falhas, orçamento e alvos relevantes;
- árvore GitHub recursiva no commit, sem fallback fictício;
- planejamento orientado por `IdeaSpec` e regras;
- símbolos conservadores com nome e linha;
- testes remotos descritos como implementação encontrada, não executada;
- evidência auxiliar limitada a força auxiliar;
- componentes emitidos somente a partir de símbolo de implementação;
- licença reconhecida por marcadores exatos;
- operação local positiva e upload externo por frase contextual;
- metadados de manutenção, requisitos operacionais e non-goals separados;
- coleta de um candidato não aborta os demais;
- validação de hashes e referências antes do assessment;
- cálculo offline das seis métricas do passo.
- inclusão explícita de controles de calibração por `--include owner/repo`, sem
  alterar score ou decisão da peneira;
- normalização Unicode determinística para termos acentuados;
- evidência negativa limitada à mesma linha do sinal relacionado;
- sinais fracos nunca convertidos em `ABSENT`;
- conflito de manutenção declarado prevalecendo sobre atividade recente.

O pipeline permanece sem LLM, embeddings, reranker neural ou serviço de IA.

## Execução real

O fluxo `plan → collect → assess → evaluate` foi executado sobre quatro
ideias do benchmark, com a versão de regras `2026-08-24.4`. Os controles explícitos
permitiram avaliar projetos relevantes e irrelevantes que a peneira original não
colocava juntos, sem lhes conceder artificialmente uma decisão melhor.

```text
ideias avaliadas:              4
candidatos avaliados:         22
capacidades PROVEN emitidas:  4
precisão dos PROVEN:          100% nas ideias em que foi mensurável
taxa de falsa ausência:       0% nas quatro ideias
rastreabilidade das fontes:   100%
conclusão dentro do orçamento: 100%
amostra mínima de PROVEN:     não atingida
```

A execução revelou um falso conflito operacional: palavras como “send”, “API” e
“cloud” em regiões diferentes do mesmo arquivo eram combinadas. A regra foi
restringida a transferência explícita de repositório/código para destino externo
na mesma linha. Também revelou que comentários em arquivos de teste podiam receber
força executável; somente declarações de teste relacionadas mantêm força forte e
matches auxiliares passaram a ser limitados a `0,40`.

A segunda rodada encontrou outros quatro modos de falha e os transformou em
regressões: acentos impediam a seleção da regra; uma negação distante contaminava
o arquivo inteiro; metadados recentes escondiam a declaração "no longer
maintained"; e sinais fracos podiam terminar como ausência. Também removeu-se
`lint` como padrão isolado, pois ele aceitava ferramentas que apenas formatavam a
saída de outro analisador.

## Resultado por ideia

| Ideia | `evidence_precision` | `false_absence_rate` | `inconclusive_rate` | `PROVEN` | Amostra suficiente |
|---|---:|---:|---:|---:|---|
| impacto estrutural | 100% | 0% | 0% | 2 | não |
| análise estática | 100% | 0% | 25,0% | 2 | não |
| eventos assíncronos | não mensurável | 0% | 28,6% | 0 | não |
| organização de receitas | não mensurável | 0% | 50,0% | 0 | não |

Os dois últimos resultados são deliberadamente conservadores. `stackus/edat`
ficou `INCONCLUSIVE`: o título declara um toolkit de arquitetura orientada a
eventos, mas o recorte coletado não trouxe implementação forte. `reaper47/recipya`
ficou `PARTIAL`: README, metadados e um módulo de integração convergem, mas não
atingem o limiar estrutural de `PROVEN`. Esses casos medem falta de cobertura sem
transformá-la em falsa certeza.

Os julgamentos manuais das quatro ideias estão versionados em
`02_shell/cli/testdata/evidence_calibration/`. Os bundles completos permanecem
artefatos de replay da execução, pois contêm snapshots extensos de código de
terceiros e ainda não devem ser incorporados ao repositório.

## Gate

Ainda não há dez capacidades `PROVEN` revisadas. Há quatro emissões `PROVEN`
observadas; portanto, a precisão de 100% é apenas descritiva e não libera o gate.
`component_precision` continua não mensurável porque nenhum julgamento de
componente reutilizável foi concluído nesta amostra.

Portanto:

- a implementação do endurecimento está pronta;
- o corpus e a boca de métricas estão prontos;
- o Passo 0006 permanece `EM CALIBRAÇÃO`;
- o Passo 0007 pode ser especificado, mas não implementado antes do gate.

## Encerramento posterior pelo Passo 0007

O diagnóstico de cobertura de 2026-08-26 substituiu a conclusão estatística
provisória acima. Depois de deduplicar aliases por par repositório–regra, corrigir
o planejamento do problema, priorizar paths relacionados e ampliar o corpus com
controles positivos e negativos, o replay atingiu onze unidades independentes
`PROVEN` revisadas.

As métricas finais observadas foram: precisão 100%, recall 100%, falsa ausência
0%, rastreabilidade 100%, conclusão de orçamento 100% e precisão de componentes
100% na primeira amostra emitida. O manifesto e os hashes reproduzíveis estão em
`02_shell/cli/testdata/evidence_calibration/calibration-run-2026-08-26.manifest.json`.

Com isso, o Passo 0006 passa a `IMPLEMENTADO` e libera o Passo 0008.

## Verificação automática

Os testes cobrem ausência insegura, árvore falsa, símbolo distante, componente,
licença por substring, operação, identidade, hashes e reprodutibilidade das
métricas. A verificação final da suíte está registrada na entrega desta execução.
