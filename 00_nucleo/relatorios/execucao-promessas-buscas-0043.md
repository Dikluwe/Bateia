# Execução 0043 — Promessas, buscas e planos verificáveis

Data: 2026-09-07. Origem: [plano 0043 congelado](../passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md). Evidências: [registro da execução](../benchmark/0043/README.md).

## Resultado e alcance

As etapas A–F receberam implementação ou integração com os controles existentes. A etapa G foi exercitada em um ensaio sintético público de três domínios e controles positivos separados; a utilidade para usuários reais e a qualidade da descoberta pública continuam não verificadas. Não se atribuiu nota global ao Bateia.

A mudança central é um contrato de alcance: busca descreve onde um candidato foi encontrado; fonte documental sustenta declaração; evidência admissível sustenta uma avaliação delimitada; integração e substituição continuam exigindo ensaio no alvo. O plano agora acompanha a decisão e expõe pergunta, lacuna, próxima ação e critério de aceitação.

## Engenharia entregue

| Etapa | Implementação / ligação | Limite preservado |
| --- | --- | --- |
| A — raiz local | `03_infra/local/client.go`: `os.Root` na abertura, leitura e caminhada; limite efetivo de bytes; hash dos bytes retornados e snapshot | Git pode consultar metadados de gitdir externo legítimo. Não há snapshot atômico contra escrita concorrente; mounts/hardlinks não são uma fronteira isolada |
| B — autoridade | `evidence_authority.go`: admissão comum em decisão, alternativas e gates de planejamento; PROVEN/ABSENT sem fonte admissível são inconclusivos; tipos e referência verificados | Admissão de ponteiro não é execução do candidato. Regras customizadas fornecidas não são certificadas automaticamente; consumidores de bundles também precisam vincular conteúdo |
| B — polaridade | Negação de operação local não vira operação comprovada; negação de upload não afirma transferência; declaração local positiva fica PARTIAL | Regras de linguagem natural são delimitadas; não há compreensão universal de negações |
| C — busca por requisito | IDs estáveis por tipo/texto entram na requisição ao provedor e ficam no snapshot; plano preserva consulta original/efetiva e mudança de query | Alteração de query não permite inventar quais filtros foram relaxados; informações ausentes permanecem explícitas. A triagem e avaliação usam a intenção original, não o nome do candidato |
| D — investigação e mecanismos | Plano diferencia falta de coleta e ausência de detector; extrator aceita digest hexadecimal do produtor, verifica tipo/localização/snippet e mantém gaps para regras fora do catálogo | Não se criaram provas nem IDs artificiais para preencher o catálogo. Correspondência textual de símbolo não é resolução semântica universal |
| E — plano condicionado | Decisão recebe `investigation_plan`; Katamari/autofagia reutilizam admissão e mantêm propostas e verificações do 0041 | Cobertura de capacidades não demonstra integração viável nem preservação geral do comportamento |
| F — apresentação e identidade | JSON CLI/MCP transporta o plano; Markdown mostra perguntas, lacunas e fontes; esforço identifica método heurístico; versão analyzer0043 invalida retomada anterior | O gateway certificado foi preservado; texto público é projetado ao construir o plano. Tempos são telemetria da ferramenta |
| F — complexidade | Estimativa lexical Go usa scanner e ignora comentários/literais; API antiga preservada como wrapper | Não é cálculo de McCabe por CFG, nem suporte equivalente a todas as linguagens; a função não possuía consumidor de produção |
| G — ensaio | 9 casos públicos sintéticos: renderização, receitas e eventos; documentação, irrelevância e busca vazia | Não mede recall ao vivo, sensibilidade a prova comportamental positiva ou utilidade humana |

Os controles já existentes para saúde não avaliada, linguagem dos relatórios, limites de extração e propostas de Katamari continuam cobertos pela suíte do projeto. Não foram reimplementados apenas para produzir mais alterações.

## Tekt e autoria

Regime A/B, executado sem atestação de isolamento. O plano e hashes do baseline foram congelados antes da implementação candidata. Prompts L0 precederam as alterações; produção mantém linhagem e camadas. Implementadores e testador receberam contextos distintos; o filesystem e suas permissões são compartilhados, portanto isso não demonstra isolamento físico.

O testador não leu a implementação candidata. Fixtures públicas antigas foram migradas por ele: testemunhos válidos substituíram status sem fontes nos controles positivos, sem retirar as verificações de decisão e matriz. Expectativas que contradiziam diretamente o novo contrato foram alteradas explicitamente: README deixa de provar execução; decisão de entrada opaca deixa de ser FORK; identidade de execução muda com analyzer0043. Duas fixtures independentes inicialmente inválidas foram corrigidas e arquivadas antes de nova execução, sem mudar suas asserções.

Testes de complexidade têm autoria compartilhada com essa implementação, indicada separadamente. Os gates históricos não tiveram certificados ou hashes reescritos. Uma alteração inicial no arquivo certificado de recomendação foi retirada; a projeção nova ficou na construção do plano, preservando o gateway anterior. Isso reduz a superfície alterada, sem alegar que o certificado histórico cobre a nova funcionalidade.

## Evidência e limitações

O ensaio G registrou 0/9 falsas comprovações, 9/9 abstenções corretas, 9/9 planos com ação/aceitação segundo critério mecânico e 9/9 resultados determinísticos. Esses denominadores pertencem somente aos casos negativos/opacos definidos previamente. Controles positivos separados verificam arquivos internos legíveis, ponteiros admissíveis e decisões de reúso com fixtures tipadas; eles também não demonstram utilidade humana.

A execução local real do Bateia produz um plano com lacunas explícitas quando o catálogo não sabe avaliar uma capacidade. Essa abstenção não prova que o projeto não implemente a capacidade. A investigação agora orienta completar a coleta ou definir um detector conforme o motivo observado, em vez de prometer que qualquer busca adicional produzirá comprovação.

A suíte completa (`go test ./...`), o lint e a compilação passaram. Na execução local houve 24 ações de investigação e decisão INCONCLUSIVE; assessment, decisão e alternativas foram idênticos byte a byte no replay. Comandos finais, versões, hashes, logs e resultados estão registrados em `validation.json`. Aprovação da suíte e do lint é verificação desta revisão, não certificado de computabilidade, precisão ou ausência de todos os defeitos.

## Trabalho que exige avaliação posterior

- Medir utilidade com decisões de usuários e casos externos, mantendo julgamentos separados dos ajustes do motor.
- Implementar detectores novos somente a partir de requisitos e testemunhos definidos; não preencher lacunas por similaridade nominal.
- Verificar integração real das propostas escolhidas em ambiente e contratos específicos.
- Ampliar testes de concorrência da worktree e de construções de linguagem natural fora das regras atuais.

O passo congelado continua como especificação de origem; este relatório e o índice registram o estado de execução. Não foram feitos commits nesta execução.
