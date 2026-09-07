# Passo 0043 — Delimitar promessas e planejar buscas verificáveis

Estado: PLANO DOCUMENTADO — execução de correções de produção pendente.
Data: 2026-09-07.
Origem: esclarecimento do usuário após o 0042: “ilusão” significa uma promessa que excede o que o procedimento e a evidência conseguem sustentar. Não é julgamento sobre o projeto inteiro nem sobre seu autor.

## Objetivo de produto

Preservar D-001 e D-002: ajudar a decidir o que reutilizar antes de construir. Converter cada promessa em uma pergunta delimitada, um procedimento, evidência necessária e uma saída permitida. As buscas devem encontrar candidatos e reduzir dúvidas específicas para produzir um plano de investigação/adaptação. Encontrar algo parecido não demonstra substituição ou integração viável.

O [0042](../relatorios/sanitizacao-promessas-0042.md) fornece contraexemplos de implementação. Este passo acrescenta o contrato de alcance e a sequência de trabalho; não declara corrigidos aqueles defeitos.

## Classificação das promessas

Aplicar a cada função, sem rotular um subsistema inteiro como “não computável”:

- **Cálculo delimitado:** resultado verificável sobre entradas finitas e declaradas, como hashes, caminhos, conjunto de candidatos e cobertura registrada.
- **Estimativa:** aproximação sujeita a erro, como relevância ou provável esforço. Exige método, hipóteses e limites; não autoriza prova.
- **Verificação condicionada:** propriedade que depende de configuração, contratos, ambiente ou experimento. Sem essas condições, manter desconhecido.
- **Promessa geral não assumida:** completude de busca em todo software existente, equivalência irrestrita ou sucesso de qualquer integração. Reformular como pergunta de escopo finito; não afirmar impossibilidade matemática sem demonstração própria.

Limite de orçamento, falta de evidência, defeito de código e impossibilidade geral são causas diferentes. A saída deve dizer qual delas impede uma conclusão.

## Contratos por função

| Função / promessa a corrigir | Pergunta executável | Evidência exigida e conclusão permitida | Limite / próxima ação |
| --- | --- | --- | --- |
| Descobrir “a solução existente” | Quais candidatos foram encontrados nos provedores, consultas e janela registrados? | Origem, consulta efetiva, filtros, orçamento e falhas; candidato encontrado | Zero resultados não prova inexistência; ampliar estratégia ou encerrar com escopo explícito |
| Buscar pela intenção | Que parte da necessidade motivou cada consulta? | ID de requisito, termos originais/expandidos, motivo da expansão; pista de busca | Expansão lexical ou sugerida por IA não vira capacidade do candidato |
| Relevância / “solução exata” | Quais requisitos possuem evidência pertinente neste candidato? | Campos e termos de busca separados de testemunhos de implementação | Score ordena investigação; não autoriza sozinho recomendação |
| Análise semântica | Que relações o analisador resolve nestes arquivos e configuração? | Linguagem, backend, arquivos/intervalos, símbolos resolvidos e limites | AST não comprova comportamento geral; relações não resolvidas permanecem desconhecidas |
| Capacidade / operação | A fonte declara, implementa ou demonstra o requisito delimitado? | Tipo de testemunho, polaridade, identidade e vínculo com o requisito | README prova declaração; um teste só demonstra a condição exercitada |
| Alternativas / reúso | Há base suficiente para sugerir esta ação sob restrições conhecidas? | Evidência admissível por requisito e gates explícitos | Sugestão condicionada; bloqueio desconhecido não vira aprovação |
| Katamari | Quais conjuntos de candidatos cobrem os requisitos registrados sob limites de busca? | Universo finito, cobertura, orçamento, interfaces conhecidas e desconhecidas | Cobertura não significa composição executável; produzir plano de integração e ensaios |
| Autofagia | Que responsabilidades se sobrepõem e podem ser investigadas para consolidação? | Inventário, chamadas/contratos, invariantes, consumidores e contraexemplos | Sem validar preservação dos comportamentos definidos, produzir proposta; não remover automaticamente |
| Mecanismos / melhorias | Há testemunhos dos mecanismos definidos no catálogo? | Produtor de evidência compatível, regra, fonte vinculada e limites do catálogo | Hipótese de transferência requer ensaio no alvo; nome semelhante não demonstra mecanismo |
| Esforço / saúde | O que foi medido, quais regras se aplicaram e o que falta avaliar? | Método, cobertura e fatores observados | Não converter ausência de achados em saúde perfeita, nem contagem em custo comprovado |

## Plano das buscas

### 1. Transformar a ideia em obrigações verificáveis

Antes das consultas, representar requisitos obrigatórios, preferências e não objetivos separadamente. Cada requisito recebe ID estável, pergunta observável e condição de aceitação. Ambiguidade relevante fica registrada como hipótese ou pergunta ao usuário; o motor não inventa um requisito obrigatório.

Exemplo: “renderização OIT em OpenGL” deve procurar uma implementação relacionada a transparência independente de ordem e identificar técnica, API e limitações. Encontrar `USalign` ou `BinDiff` por termos amplos não satisfaz esse requisito. O caso é regressão; não define o domínio do produto.

### 2. Separar descoberta de investigação de evidência

Consultas de descoberta procuram produtos, capacidades e vocabulários alternativos. Consultas de investigação procuram implementação, API, testes, limitações e contraindicações dentro de candidatos. Preservar para ambas: requisito de origem, consulta original, expansão, consulta enviada, filtros, fallback, provedor, resultados, erros e limites.

Relaxar filtros pode aumentar o conjunto de candidatos, mas não altera silenciosamente o requisito original. Todo candidato novo é reavaliado contra esse requisito. Dependência, popularidade e coincidência de nome permanecem sinais de triagem.

### 3. Buscar para resolver a próxima dúvida

Priorizar arquivos e consultas que respondam a uma lacuna concreta: “há implementação?”, “depende de serviço externo?”, “qual interface permite integrar?”. Registrar a pergunta que cada coleta tenta resolver. Procurar também restrições, negações e incompatibilidades; evidência contraditória não pode ser descartada para aumentar cobertura.

Manter separados: arquivo localizado, conteúdo lido, relação estática examinada e comportamento exercitado. Registrar limites de extração e impedir que truncamento de apresentação seja confundido com cobertura da análise.

### 4. Encerrar com motivo e plano

Encerrar por orçamento, falha, saturação das estratégias configuradas ou evidência suficiente para a decisão delimitada. Saturação local não significa busca global completa. Para cada lacuna, fornecer ação seguinte, evidência esperada e condição de parada. Tempos de busca são telemetria do Bateia, nunca benchmark do candidato.

## Contrato mínimo da saída planejada

Cada conclusão precisa carregar, no formato de domínio a definir na implementação:

- requisito e pergunta respondida;
- classe: fato observado, inferência, hipótese ou não verificado;
- fonte admissível: arquivo/intervalo/hash ou execução identificada, com proveniência;
- escopo da análise e método empregado;
- contrapontos, condições e lacunas;
- ação sugerida, verificação necessária e critério de aceitação.

Uma inferência aponta para fatos que a sustentam. Uma hipótese aponta para sua origem e para o experimento que a pode refutar. “Não verificado” aponta para a lacuna ou falha de coleta; não requer fabricar uma fonte positiva. A existência de campos preenchidos não valida o conteúdo: o consumidor que autoriza a decisão deve verificar seus vínculos.

## Sequência de implementação

As etapas abaixo são backlog deste plano, não funcionalidades concluídas. Cada etapa exige especificação espelhada quando alterar Go, implementação e regressões pertinentes antes de mudar a promessa pública.

| Ordem | Entrega | Dependência | Aceitação objetiva |
| --- | --- | --- | --- |
| A | Conter leituras locais e registrar identidade efetivamente lida | Independente das demais; prioridade imediata pelo 0042 | Links finais/intermediários externos rejeitados; arquivos internos válidos continuam legíveis; mudança de conteúdo altera identidade |
| B | Contrato comum de autoridade de evidência na decisão e nos planejadores | Base para D e E | Status sem fonte, fonte incompatível e documentação isolada não comprovam execução; frases negativas do 0042 não viram provas positivas |
| C | Plano de consultas por requisito, fallback explícito e reavaliação | IDs de requisitos; conectar a B | Consulta efetiva e relaxamento auditáveis; candidato irrelevante de fallback não ganha autoridade; zero resultados mantém limite explícito |
| D | Coleta orientada à pergunta e ligação real avaliação → mecanismos | B; usar perguntas de C | Teste com produtor e consumidor reais valida hash/tipo/regra; catálogo sem evidência retorna desconhecido; implementações relevantes ou cobertura insuficiente |
| E | Plano de reúso, Katamari e autofagia condicionado a contratos | B e D | Cobertura não vira integração garantida; proposta lista invariantes, interfaces, ensaios e lacunas; ausência de ensaio mantém viabilidade não verificada |
| F | Relatórios e métricas de alcance coerentes entre CLI/MCP | B–E | Mesma autoridade nos formatos; complexidade lexical e esforço heurístico identificados; sem regras aplicáveis, não avaliado |
| G | Avaliação transversal da utilidade | Contratos anteriores estabilizados | Casos e julgamentos definidos antes da execução, separados de ajuste; resultados publicados por requisito/domínio, incluindo erros e abstenções |

A aprovação em G não deve depender só de uma média: publicar falsa comprovação, rejeição indevida, cobertura de evidência, abstenção correta e utilidade do plano. Definir tamanho da amostra, denominadores e limiares antes de executar, sem ajustá-los aos resultados. Nenhuma medição nova foi feita neste passo.

## Regressões obrigatórias do plano

- OIT retorna USalign/BinDiff: manter candidato sem evidência pertinente, sem “solução exata”.
- Renderização em C/C++: extrair implementação pertinente ou declarar cobertura insuficiente; sem convenções de testes Go.
- Nenhuma regra aplicável / zero candidatos: não avaliado / nenhuma observação do ecossistema.
- README nega operação local: não promover a capacidade; README nega envio externo: não afirmar envio só por palavras.
- Avaliação sem fontes: seletor não autoriza FORK com capacidade comprovada.
- Evidência real → mecanismos: contrato compatível, sem fabricar IDs de regras para forçar resultado.
- Katamari cobre capacidades, mas interfaces são desconhecidas: plano com lacunas, sem integração confirmada.
- Autofagia encontra sobreposição nominal: nenhuma garantia de substituição sem invariantes e validação delimitada.
- Worktree muda / links escapam da raiz: identidade muda / leitura fora de escopo rejeitada.

## Entrega deste passo

Plano registrado e vinculado aos desejos e README, com alcance das promessas esclarecido. Não implementa novas buscas nem altera decisões em produção. Preserva mudanças anteriores e relatórios históricos. A primeira execução de engenharia deve começar por A e B, seguida de C; acrescentar mais provedores ou maior volume de busca antes disso pode apenas aumentar a quantidade de evidência mal interpretada.
