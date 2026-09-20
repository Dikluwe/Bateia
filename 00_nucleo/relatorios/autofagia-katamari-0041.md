# Autofagia e Katamari — engenharia do passo 0041

Data: 2026-09-07. Implementação concluída sobre o commit-base registrado nos artefatos; sem incorporar dependências ou código externo.

## O que mudou

Autofagia continua respondendo o que pode ser aproveitado internamente. Agora explicita o escopo das revisões examinadas e oferece uma opção por capacidade: manter, revisar sobreposição, investigar possibilidade de completar internamente ou comparar construção com reúso. Uma avaliação faltante ou conflitante não é encoberta por outra PROVEN. Repetir a mesma avaliação não produz uma nova implementação. Múltiplas avaliações comprovadas geram REVIEW_OVERLAP, nunca CONSOLIDATE automático.

Katamari continua procurando composição entre capacidades externas. A entrada passou a conferir restrições, falhas, licença nas duas políticas, manutenção, operação exigida e fontes declaradas. Uma fonte precisa corresponder à revisão, conter caminho, regra, localização e SHA256; licença precisa de documento com autoridade. Rastreabilidade estrutural não reexecuta nem certifica o comportamento alegado, e isso aparece nas advertências.

Evidência parcial e condições desconhecidas geram opções de investigação com fontes e próxima verificação. Essas opções não entram na cobertura. Capacidades existentes podem originar comparação de substituição; não são absorvidas automaticamente. Gaps internos também geram comparação com implementação interna, sem presumir que adicionar uma dependência é melhor.

A seleção remove duplicatas, rejeita avaliações divergentes da mesma revisão e impede compor duas revisões do mesmo projeto. A enumeração mantém limite de 20 candidatos únicos elegíveis, maximiza cobertura declarada e desempata por número de projetos, redundância e identidade. Não mede esforço: o custo 100 + 10 por capacidade foi removido. Cada componente proposto recebe NOT_ESTIMATED e verificações de fronteira de integração, compatibilidade, termos, adaptação/operação e aceitação.

## Estados para consumidores

Inventário e proposta retornam `schema_version: 2`. O envelope competitivo e a entrada de políticas continuam na versão 1. Inventários legados sem findings e escopo continuam legíveis, mas listas isoladas não justificam cobertura ou adoção; geram investigação.

| Estado da proposta | Significado |
| --- | --- |
| COVERED_AWAITING_VALIDATION | Candidatos elegíveis fornecidos cobrem os gaps declarados; integração pendente |
| PARTIAL_COVERAGE | Há proposta de cobertura, com gaps remanescentes |
| INVESTIGATE | Há ações de verificação ou comparação; não há composição qualificada |
| NO_CHANGE_REQUIRED | O inventário conhecido atende às capacidades consideradas; substituição ainda pode ser investigada |
| INCONCLUSIVE | Não há base suficiente para uma conclusão ou opção delimitada |

`coverage_complete` refere-se apenas aos gaps considerados. `minimal` só é verdadeiro quando uma composição foi selecionada; `minimality_scope` informa o critério e universo fornecido. Isso não prova que nenhum projeto melhor exista fora da amostra.

O campo legado `feasible` permanece falso: este módulo não executa nem recebe um protocolo de integração validada. `components[].decision` permanece INCONCLUSIVE e `integration_status` NOT_VERIFIED. Consumidores devem usar `status`, cobertura, opções e verificações pendentes, e não promover componentes diretamente a adoção. AllowUnknown e RequireImmutableRef não podem afrouxar os requisitos de cobertura: desconhecimento permite investigar, não qualificar.

A unidade selecionada ainda corresponde a uma revisão de repositório com fontes de implementação. A fronteira concreta do componente é uma verificação pendente. Logo MaxComponents limita essas unidades propostas; não conta automaticamente módulos extraídos ou dependências transitivas.

## Validação e resultados

A suíte `go test ./...` passou, o linteiro retornou “No violations found” e o pacote wiring compilou. Os testes cobrem cenários positivos do seletor com assessments fornecidos e regressões de autoridade/bloqueios, determinação, duplicação e escopo.

Na CLI, dois inventários foram gerados e seis propostas executadas. Dois candidatos complementares produziram COVERED_AWAITING_VALIDATION; acrescentar um candidato que cobre ambos reduziu a seleção para um. Um candidato cobrindo só parte gerou PARTIAL_COVERAGE. Conflito impediu seleção; evidência parcial preservou a oportunidade investigativa; zero candidatos não virou sucesso. Todos mantiveram `feasible=false`.

As fontes PROVEN dos casos positivos são fixtures declaradas para testar a seleção. Não representam comprovação comportamental executada. Os critérios de comprovação lexical do passo 0040 permaneceram conservadores, e os artefatos anteriores não foram reescritos.

## Limites restantes

O motor não demonstra equivalência entre implementações, não estima custo real de integração e não prova compatibilidade entre componentes. Ele entrega o plano e as verificações necessárias para isso. O próximo avanço depende de protocolos delimitados por capacidade e interface, com resultado de aceitação rastreável. Não haverá um estado “integração validada” obtido apenas pela soma de pontuações ou pela declaração PROVEN de um candidato.

[Artefatos e comandos reproduzíveis](../../lab/bancada/history/benchmarks/0041/README.md).
