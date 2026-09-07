# Diagnóstico 0042 — Promessas, engenharia e limites do Bateia

Data: 2026-09-07. Estado: diagnóstico concluído; defeitos de produção abaixo ainda não corrigidos. [Plano](../passos/0042-auditar-promessas-e-limites-do-projeto.md), [reproduções e inventário](../benchmark/0042/README.md).

A evidência não sustenta dizer que todo o projeto é ilusório ou que seu autor não conseguiu fazer engenharia. Há implementações reais, infraestrutura de execução e salvaguardas úteis. O problema transversal observado é a falta de um contrato uniforme de comprovação entre módulos: uma etapa conserva a incerteza, outra aceita um rótulo ou rejeita uma evidência produzida pela etapa anterior. Documentação abrangente demais esconde essas diferenças.

Este é um diagnóstico dirigido do mesmo agente que participou de alterações anteriores. Não é certificação independente, auditoria de cada linha ou revalidação dos oráculos. Estados como “implementado” nos passos históricos não significam “validado externamente”. Os arquivos inventariados foram hashados; somente os intervalos do `review-ledger.json` representam a amostra de leitura registrada. O snapshot foi capturado durante o diagnóstico, não atomicamente no início.

## Defeitos reproduzidos e prioridades

### P1 — Negação e autoridade da evidência não são uniformes

**Fato observado:** `AssessEvidence` atribuiu `PROVEN` à operação local com a única frase `This tool cannot run locally.` no README. A frase `This tool never uploads your code to a cloud service.` virou `CONFLICT`. O controle positivo `This tool can run locally.` também virou `PROVEN`, embora seja somente declaração documental. Como controle separado, a triagem retornou conflito zero para `This tool does not upload code.` com não objetivo `upload code`; esse caso não reproduziu o defeito de operação.

Fontes: [operação](../../01_core/usecase/evidence.go#L1147), [triagem](../../01_core/usecase/probabilistic_screening.go#L163); resultados `negative_offline`, `negative_upload`, `positive_offline`, `negative_non_goal_screening` do probe. A regra de descarte usa `NonGoalConflict > 0` na linha 281. **Inferência:** a avaliação de operação pode tanto aceitar quanto bloquear candidatos incorretamente. O caso isolado de triagem não demonstra cobertura geral de negações.

### P1 — A decisão principal aceita status sem comprovação

**Fato observado:** `DecideReuse` recebeu uma avaliação sintética com status `PROVEN`, sem fontes nem partes reutilizáveis, e produziu `FORK`, cobertura de 100% e esforço `LOW`. Fonte: [seleção e gates](../../01_core/usecase/reuse_decision.go#L32), `decision_from_sourceless_assessment`.

O experimento exercita a entrada pública do seletor com avaliações fornecidas; não demonstra que a coleta lexical atual gera essa entrada normalmente. **Inferência:** consumidores, replays ou produtores de avaliações podem contornar a prudência do coletor. Os gates de evidência precisam existir na fronteira que autoriza a decisão, compartilhados com autofagia/Katamari.

### P1 — Avaliação e extração de mecanismos têm contratos incompatíveis

**Fato observado:** um bundle sintético passou pelo `AssessEvidence` real e foi entregue ao extrator real com o catálogo distribuído. O extrator retornou dois `EVIDENCE_CONTENT_HASH_MISMATCH`, sem mecanismos. O produtor usa hexadecimal de 64 caracteres; o consumidor exige `sha256:`. Fontes: [hash do produtor](../../01_core/usecase/evidence.go#L546), [vínculo do consumidor](../../01_core/usecase/mechanism_extraction.go#L277), [ligação do serviço](../../01_core/usecase/recommendation_run.go#L455).

A comparação dos catálogos públicos registrados também encontrou interseção vazia entre os 14 IDs de avaliação e os 32 IDs de evidência aceitos pelo catálogo de mecanismos. Ver `catalog-compatibility.json`. **Inferência:** corrigir apenas o prefixo não conecta o catálogo padrão à jornada. Isso não prova que entradas customizadas sejam incapazes de produzir mecanismos. Tampouco autoriza traduzir rótulos lexicais em provas de mecanismo: é necessário implementar os produtores das evidências exigidas.

### P1 — A raiz local não contém todas as leituras

**Fato observado:** em diretórios temporários sintéticos, `FetchManifest` leu um `go.mod` simbólico externo à raiz. `ExtractRepoFiles` bloqueou o link final direto, mas leu `linked/marker.go` por um diretório intermediário simbólico externo. Fonte: [adaptador local](../../03_infra/local/client.go#L49), `local_symlink_boundary`.

O resultado limita-se à leitura do manifesto e de caminhos explicitamente solicitados. Não demonstra que a caminhada padrão siga todos os links. Nenhum arquivo privado foi usado. Corrigir resolução/contenção em todos os pontos de leitura, considerando também mudança do caminho entre verificação e abertura; a resistência a essa corrida ainda não foi avaliada.

### P2 — Métricas e pureza alegadas excedem o que é calculado

- **Fato observado:** complexidade passa de 1 para 6 ao adicionar somente o comentário `// if for while switch case`. [Calculador](../../01_core/parser/cyclomatic.go#L21). É contagem lexical, não cálculo do grafo de controle; deve ser nomeada como aproximação ou substituída onde houver parser adequado.
- **Fato observado:** esforço `LOW` resulta de contagens de componentes/lacunas, sem ensaio de integração. [Estimador](../../01_core/usecase/reuse_decision.go#L301). Não representa horas, custo ou viabilidade comprovada.
- **Fato observado:** duas chamadas idênticas de `ParseCsProj` retornam instantes distintos em `ConsultadoEm`. [Parser](../../01_core/parser/csharp.go#L29). A afirmação de pureza de todo L1 é excessiva. Isso não prova falha em replay de artefatos já congelados.

## Matriz transversal: capacidade existente e limite

| Área | Evidência e estado | Limite da conclusão |
| --- | --- | --- |
| Propósito e uso | README e D-001/D-002 definem decisão de reúso; saneamento documental neste passo | Visão desejada não equivale a entrega validada |
| Descoberta e fallback | Sourcegraph, linhas 172–220, registra consulta e relaxamentos | Não houve pesquisa remota nova; recall público não avaliado |
| Expansão semântica | `semantic_expander.go`, linhas 70–140, valida derivações; wiring usa expansão lexical | Nome “semântico” não demonstra compreensão geral da intenção |
| Pontuação | `probabilistic_screening.go`, linhas 106–113 e 250–284: pesos e estado `uncalibrated` explícitos | Pesos não são probabilidades validadas; controle isolado de negação não certifica o classificador |
| Calibração | `screening_calibration.go`, linhas 211–271: ajuste e comparação de Brier por grupos | Critério local não demonstra representatividade externa |
| Extração e linguagens | Dispatcher de manifestos, AST Go e adaptador Clang implementados | Ler manifesto não é analisar implementação; Clang depende de configuração e limites, sem cobertura equivalente em todas as linguagens |
| Comprovação e resposta | `query_evidence.go`, trecho 220–260, explicita classes de evidência; histórico 0040 limita prova lexical | Avaliação de operação e seletor principal ainda divergem desse contrato |
| Decisão e alternativas | Seletor implementado; sem avaliações retorna inconclusivo | `FORK` sem fontes reproduzido; composição por cobertura não verifica compatibilidade de interfaces |
| Autofagia/Katamari | [0041](autofagia-katamari-0041.md) registra inventário, gates e propostas com integração não verificada | Evidência anterior, não nova validação independente; não demonstra integração funcionando |
| Mecanismos e melhorias | Extração valida testemunhos; síntese declara hipóteses e `NOT_AN_ADOPTION_DECISION` | Incompatibilidade reproduzida; propósito/estratégia vêm do catálogo, não de reconstrução completa da arquitetura |
| Arquitetura e segurança | Linter retorna não avaliado sem regras; gerador qualifica penalidades; scanner tem base local pequena e datada | Não são certificado de saúde nem scanner completo de vulnerabilidades |
| Validação de produto | `product_validation.go`, linhas 16–69, deixa downstream não avaliado sem métricas; benchmark executa pipeline | Utilidade no benchmark vem de julgamento fornecido (`end_to_end_benchmark.go`, linhas 153–204), não de benefício medido automaticamente |
| Rastreabilidade e execução | Serviço possui persistência/idempotência; 0040 registra replay delimitado | Não se reexecutaram testes de falhas, concorrência, durabilidade ou toda a superfície MCP; limite local por symlink confirmado |

As referências de código dessa matriz, com hashes e intervalos, estão no [registro de leitura](../benchmark/0042/review-ledger.json). Resultados de testes 0041 são evidência histórica de regressão daquele checkout, não prova de precisão do produto. Tempos de execução do Bateia são telemetria da ferramenta; não medem desempenho dos candidatos.

## Saneamento executado

O README agora separa leitura de manifestos, AST e comportamento; troca “necessidades ocultas” por observações/hipóteses, qualifica pureza e aprovação de testes, registra limites conhecidos e compila o pacote inteiro de wiring. O índice de relatórios qualifica registros históricos. Não foram reescritos resultados antigos nem alterado código de produção neste passo. As mudanças pendentes do 0041 foram preservadas.

## Sequência recomendada de correções

1. Conter leituras locais e uniformizar evidência na decisão: bloquear classificações sem testemunhos válidos; distinguir afirmação, negação e declaração documental em todas as etapas. Regressões: frases acima, fonte ausente, fonte incompatível e caminhos simbólicos.
2. Conectar avaliação e mecanismos por contrato versionado de identidade, tipo, hash, regra e autoridade. Validar produtor real → consumidor real, incluindo ausência de prova; não fabricar regras ou provas para preencher a saída.
3. Qualificar esforço/complexidade e retirar relógio dos cálculos que prometem determinismo. Separar metadados de coleta dos valores analisados.
4. Só depois medir a utilidade das decisões em ideias e ecossistemas variados, com julgamentos separados da implementação. Medir falsa certeza, abstinência correta e qualidade da decisão, além de recall.

**Não verificado:** qualidade global, independência de todos os benchmarks históricos, precisão em todos os ecossistemas, provedores ao vivo, durabilidade sob falhas e integração real das propostas. **Hipótese para investigação posterior:** hashes capturados em momentos diferentes podem não representar um único estado atômico da worktree. Nenhuma nota global foi atribuída.
