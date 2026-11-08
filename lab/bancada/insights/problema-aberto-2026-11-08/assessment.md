# Parecer independente: pacotes X e Y

Há exclusivas qualificadas observadas em ambos os pacotes: X-A1 e X-A3 em X, Y-A2 em Y. A contribuição incremental da pesquisa é defensável, com confiança moderada, para as revisões X-A3 e Y-A2; não está demonstrada para X-A1. X-A2/Y-A1 formam um núcleo equivalente forte. Isso não estabelece superioridade global, nova família descoberta nem utilidade de campo comprovada.

A avaliação usa somente o problema, a rubrica e os dois pacotes fornecidos. Não houve busca externa, execução, inspeção do repositório ou inferência de identidade. Notas ordinais não foram somadas. Evidência 3 significa mecanismo e limite conferíveis em fonte estática, sem prova de eficácia.

## Notas por alternativa

| Alternativa | Relevância / 4 | Evidência / 3 | Especificidade / 2 | Qualificada | Baseline |
|---|---:|---:|---:|---|---|
| X-A1 | 4 | 3 | 2 | Sim | mesma família |
| X-A2 | 4 | 3 | 2 | Sim | mesma família |
| X-A3 | 3 | 3 | 2 | Sim | revisão material |
| Y-A1 | 4 | 3 | 2 | Sim | mesma família |
| Y-A2 | 3 | 3 | 2 | Sim | revisão material |

### X-A1

**Relevância:** Ficha local, reconciliação humana, arquivo de observações separado e regras comuns de custódia respondem aos problemas centrais; teste pequeno tem critérios de perda e tempo de resolução.

**Evidência:** Documentação específica comprova inspeção explícita de conflitos e limite de retenção de revisões, além de armazenamento, retry e anexos. Não comprova durabilidade em Android nem adequação operacional.

**Especificidade:** O estado corrente é editado como ficha; resolver conflitos de ficha é trabalho explícito, sem derivar todo estado de eventos.

**Trechos usados:** P-conflict: opção {conflicts: true}; P-compact: non-leaf revisions will no longer be available.

**Referências no pacote:** X:A1:P-local; X:A1:P-conflict; X:A1:P-compact; X:A1:P-sync; X:A1:P-photo; X:A1:P-fetch; X:A1:P-quota.

**Limites pendentes:** Durabilidade local por 72h, espaço e fotos; Escala de 45.000 peças e custos; Tempo e disponibilidade de conciliador; Exportador e recuperação próprios.

**Baseline (X-B1, X-B4):** B1 já continha observações append-only e reconciliação explícita. A opção por ficha corrente é específica, mas fila e preservação separada não demonstram alteração de uma garantia anterior: o baseline não propunha usar revisões como arquivo. A alegação de revisão material do pacote é mais forte que a comparação disponível. Confiança moderada.

### X-A2

**Relevância:** Histórico imutável, projeções reconstruíveis, transferência em duas etapas e confirmação pela custódia tratam perda de observações e promessas falsas; teste inclui reenvio, etapa ausente e reconstrução.

**Evidência:** Estratégia append-only e atomicidade por documento sustentam mecanismo e limite decisório. A regra de custódia é desenho operacional, não garantia da biblioteca.

**Especificidade:** Correção compensa registro anterior; autoridade confirma disponibilidade e trânsito bloqueia promessa.

**Trechos usados:** P-delta: never update or remove existing documents; P-atomic: atomic unit is the document.

**Referências no pacote:** X:A2:P-delta; X:A2:P-compact; X:A2:P-sync; X:A2:P-photo; X:A2:P-atomic; X:A2:P-fetch.

**Limites pendentes:** Contradições físicas e adesão à custódia; Projeções, migrações e exportador próprios; Durabilidade, escala e custo não medidos.

**Baseline (X-B2, X-B3):** Log, projeções, fotos separadas e custódia já eram explícitos. Unidade por comando e tratamento de recebimento ausente são refinamentos importantes da implementação, sem nova família ou mudança demonstrada da regra operacional central. Confiança alta.

### X-A3

**Relevância:** Alternativa condicional para concorrência descritiva frequente, com decisões físicas separadas e teste de descarte. App, persistência, transporte e saída ainda exigem muito trabalho de dois generalistas; benefício para esta cooperativa não foi demonstrado.

**Evidência:** Código/documentação de edição local, conflitos de campos, histórico e exigência de transporte sustentam mecanismo e limite; não são resultados executados.

**Especificidade:** Mesclar mapas/listas descritivos reduz a unidade da conciliação em relação à ficha inteira, mas preserva revisão de valores concorrentes e custódia para estoque.

**Trechos usados:** A-conflict: arbitrarily (but deterministically; A-history: getAllChanges(doc); A-transport: reliable in-order stream.

**Referências no pacote:** X:A3:A-local; X:A3:A-conflict; X:A3:A-history; X:A3:A-transport; X:A3:A-save.

**Limites pendentes:** Frequência real de colaboração descritiva; Custo de transporte/persistência/app e operação; Memória e CPU com 45.000 fichas; Exportação legível e mídia; Teste proposto de 500 fichas não verifica sozinho escala final.

**Baseline (X-B1, X-B3):** A mudança identificável é editar e mesclar estruturas descritivas com conciliação por campo, em vez de apenas registrar observações e reconciliar edições/eventos. A-local e A-conflict sustentam essa mudança e a necessidade de expor alternativas. Não é a mera menção de CRDT. O baseline genérico impede atribuição causal forte. Confiança moderada.

### Y-A1

**Relevância:** Diário imutável, custódia, reserva física, trânsito e exportação reconstruível respondem ao problema; piloto futuro cobre 72h, 45.000 metadados, fotos e falhas com rejeição concreta.

**Evidência:** Estratégia de documentos imutáveis, vencedor determinístico, retry, anexos e compactação sustentam mecanismo e limites. Não há prova de campo.

**Especificidade:** Promessas dependem de separação e confirmação pela oficina custodiante; conflitos semânticos viram exceção humana.

**Trechos usados:** conflicts.md: never update or remove existing documents; compact-and-destroy.md: non-leaf revisions will no longer be available.

**Referências no pacote:** Y:A1:docs/guides/conflicts.md:140-154; Y:A1:docs/guides/conflicts.md:74-90; Y:A1:docs/guides/replication.md:109-116; Y:A1:docs/guides/attachments.md:8-16,73-83; Y:A1:docs/guides/compact-and-destroy.md:26-46.

**Limites pendentes:** Durabilidade e espaço dos Android; Custódia aceita pelos voluntários; Custo de manutenção e infraestrutura; Implementação de projeções/exportador.

**Baseline (Y-B1, Y-B2):** Append-only, resolução explícita e handoffs já constavam no baseline. As fontes tornam o desenho conferível e delimitam revisões, sem acrescentar família ou alterar a autoridade operacional central. Confiança alta.

### Y-A2

**Relevância:** Formulários Android offline, estado por peça, custódia e revisão podem atender com menos captura/exportação próprios. Consulta de fotos históricas offline, volume e telas de frescor permanecem condições decisivas; piloto testa essas condições.

**Evidência:** Documentação específica sustenta representação local de entidades, ordem de aplicação, histórico, exportação com mídia e independência entre editar submissão e estado. Alguns trechos são curtos demais para comprovar alegações adicionais.

**Especificidade:** Configurar formulários e revisar estado recebido com política last-write-wins; correções exigem manter submissão histórica e estado corrente coerentes, em vez de reconstruir toda disponibilidade de log próprio.

**Trechos usados:** central-entities.rst: offline representation of Entity Lists; entities-intro.rst: not automatically applied to the related Entities; central-submissions.rst: multimedia submitted to the form.

**Referências no pacote:** Y:A2:docs/collect-intro.rst:13-17; Y:A2:docs/central-entities.rst:648-654,637-643,520-527,418-441; Y:A2:docs/entities-intro.rst:197-203,104-106,217-220; Y:A2:docs/central-submissions.rst:134-136; Y:A2:docs/central-install.rst:45-48; Y:A2:docs/form-audit-log.rst:4-16.

**Limites pendentes:** Foto/história acessível offline por peça; Escala, filtros e aparente divergência documental; Versões mínimas e retenção de cinco dias não demonstradas pelos excertos exibidos; Sem contexto, Not yet não comprova sozinho ausência de mídia em Entity; Custo, telas adicionais e esforço de formulários não medidos.

**Baseline (Y-B3, Y-B2):** Coleta offline já era família prevista. Representação local de entidades para seguimento e a independência entre editar submissão e atualizar estado impõem um fluxo concreto de revisão/correção; isso altera responsabilidade e garantia, além de nomear uma ferramenta. Ressalva: parte dos limites alegados não é conferível nos trechos curtos. Confiança moderada.

## Equivalência por decisões e garantias

- **G1 — X-A1:** Ficha corrente mutável, conciliação humana de conflitos documentais e arquivo histórico separado. Difere de G2 na origem do estado corrente e no trabalho de conciliação; compartilhar transporte não torna os desenhos equivalentes.
- **G2 — X-A2, Y-A1:** Operações imutáveis, projeção reconstruível, custódia e handoff explícitos com promessas provisórias. Mesmas garantias e responsabilidades centrais; atomicidade mais explícita em X e piloto mais amplo em Y não criam famílias exclusivas.
- **G3 — X-A3:** Colaboração descritiva com mescla de estruturas e exposição de conflitos de campo; autoridade física separada. A unidade de colaboração/conciliação muda; compartilhar custódia com G2 não apaga essa diferença.
- **G4 — Y-A2:** Captura por formulários e estado por entidade aplicado na recepção, com revisão e correção explícita de estado e submissões. Responsabilidade deslocada para configuração/administração e revisão de estado do sistema existente; não é apenas outra marca para o log reconstruível de G2.

## Exclusividade observada e contribuição incremental

A exclusividade descreve apenas o conteúdo desta execução. Não demonstra capacidade exclusiva de um braço, nem novidade no mundo.

- **X-A1:** contribuição da pesquisa: **não demonstrada**. Exclusiva no resultado, mas reconciliação e observações imutáveis já estavam em X-B1; a pesquisa documenta riscos sem evidenciar mudança material do baseline.
- **X-A3:** contribuição da pesquisa: **revisão material sustentada com confiança moderada**. Mescla descritiva e conciliação por campo acrescentam decisão identificável; benefício e causalidade da descoberta permanecem não comprovados.
- **Y-A2:** contribuição da pesquisa: **revisão material sustentada com confiança moderada**. A família de coleta offline já existia, mas seguimento por entidades e correção independente do estado alteram o fluxo operacional.

## Exclusivas apenas aparentes

- **X-A2 ou Y-A1 seriam exclusivos:** São equivalentes em decisões e garantias.
- **O nome de uma biblioteca/plataforma prova família nova:** Nenhuma família nova foi demonstrada só por nome. X-A3 e Y-A2 são revisões materiais condicionais, não novidade no mundo.
- **Fila humana e arquivo de observações de X-A1 foram acrescentados pela pesquisa:** O baseline já previa reconciliação explícita e observações append-only; atribuição incremental insuficiente.

## Incertezas e avaliação antes de decidir

- Cinco candidatas qualificadas não equivalem a cinco soluções prontas ou eficazes em campo.
- Notas são condicionais ao desenho e aos trechos; não houve execução dos projetos nem piloto da cooperativa.
- Locators imprecisos em X não anulam literais cuja checagem no pin foi informada; não foi feita nova conferência externa.
- Baselines breves limitam a separação entre elaboração e mudança material; não há sequência causal detalhada de decisões.
- G1/G3/G4 partilham regras de custódia e preservação de observações; exclusividade é da decisão adicional especificada, não de todo o desenho.
- Ausência de equivalente em um pacote não demonstra incapacidade do outro braço.
- Sem custos medidos ou validação de escala, durabilidade, usabilidade e adesão, não há vencedor de adequação prática.
- Nenhuma identidade dos produtos foi inferida.

O piloto mais informativo usa o cenário comum descrito em Y-A1/Y-A2: dois Android baratos, duas oficinas, 200 peças reais, 45.000 metadados e 72h sem rede. Comparar perda de observações/fotos após reinício e exportação, promessas incorretas, reenvios, correções e transferências incompletas, tempo dos voluntários e esforço recorrente dos mantenedores. Para G1 medir conciliação de ficha; para G3 medir ganho real da edição descritiva e extrapolar além das 500 fichas propostas; para G4 verificar seguimento local, correção do estado e consulta de fotos offline. Esses são experimentos futuros, não resultados.

Não há base para declarar vencedor global. Há um núcleo equivalente bem sustentado e revisões condicionais distintas em ambos os lados.
