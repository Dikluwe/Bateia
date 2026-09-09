# Bateia em relação ao objetivo do produto — após0064

Avaliação dirigida em 2026-09-08: leitura de propósito, relatórios medidos, ligações de produção e execução local atual. Não é auditoria independente de todas as linhas, nova avaliação pública ou certificação semântica. Não houve alteração de produção, testes ou parâmetros de ranking.

## Veredito

A engenharia de evidência está mais desenvolvida que a demonstração de utilidade da decisão final. Há infraestrutura real de busca, recuperação, avaliação conservadora, planos e rastreabilidade. O objetivo de receber uma ideia e decidir com evidência o que reutilizar continua parcialmente implementado e insuficientemente validado como produto.

A intenção de referência é D-001/D-002 e a abertura atual do README: evitar começar do zero quando existe base reutilizável, comparar alternativas e decidir entre reúso, extensão, integração, composição, referência ou inconclusão. Contar arquivos, aprovar testes e produzir hashes são meios. O antigo caso-de-uso-v2 conserva um escopo histórico diferente (pesquisa arquitetural Go-first e ausência de análise semântica); deve ser lido junto às delimitações e desejos posteriores, não como retrato único do produto atual.

## O que melhorou e o que permanece limitado

| Dimensão | Evidência de avanço | Limite atual |
| --- | --- | --- |
| Autoridade e linguagem | 0040/0043 retiram comprovação lexical automática, admitem fontes na fronteira da decisão e qualificam ausência/saúde/esforço | Admissão de ponteiro não demonstra comportamento; detectores e consumidores precisam validar o alcance |
| Consulta e descoberta | 0054 preserva intenção/filtros; 0055/0056 expõem planejamento e perdas | Última comparação pública examinada, 0056: 0/3 repositórios-alvo conhecidos; índice/candidatos alternativos não avaliados exaustivamente |
| Recuperação local | Corpus0053: 1/6 funções-alvo inicialmente; 0060/0061: 6/6 alvos encontrados, quatro por código e dois por documentação | Amostra passou a ser calibração. Não há nova validação independente que demonstre generalização desse ganho |
| Entendimento da ação | 0061 distingue pistas alinhadas/opostas de início/interrupção | Oito das nove perguntas ficam fora da ação suportada; nomes/chamadas não provam efeitos |
| Efeitos executados | 0062 verifica ordenação/preservação em funções Go autônomas, detectando cinco defeitos controlados | Um único contrato/assinatura. Não resolve linguagens, projetos arbitrários ou capacidades de produto |
| Rastreabilidade | 0063 vincula arquivo/declaração; 0064 identifica produtor e parâmetros | Avanço do piloto effects; não é identidade de todo worktree, assinatura autêntica ou rastreio uniforme de toda jornada |
| Autofagia/Katamari | Inventário e seleção com gates, limites e verificações pendentes | Propostas; integração real não verificada. Testes positivos de seleção usam avaliações fornecidas |
| Recomendação principal | Dossiê, alternativas, planos e INCONCLUSIVE implementados | Falta demonstrar decisões úteis em ideias novas, com capacidade específica sustentada e custo de adaptação observado |

Fontes: relatórios0040,0041,0043,0053,0054,0056,0059,0060,0061 e0062–0064. Resultados antigos não foram reinterpretados como novas execuções. Os defeitos diagnosticados no0042 e corrigidos no0043 não são apresentados como pendências atuais sem nova reprodução.

## Ensaio atual da jornada principal

Compilado o pacote ./04_wiring para /tmp/bateia-auditoria-atual. Executado o exemplo local documentado no README com IdeaSpec/política do benchmark0012, alvo Bateia e worktree atual. Saída: dossiê978dac098c99a801, decisão INCONCLUSIVE.

- Dez capacidades: dez INCONCLUSIVE; nove indicam ausência de regra versionada, uma indica coleta incompleta.
- Plano: 24 requisitos e 24 ações de investigação.
- Bundle: 44 itens, 404726 bytes lidos, complete=false; target_results registra 43 COLLECTED, 1002 SKIPPED_BUDGET e quatro NOT_FOUND. Contagens referem-se a observáveis diferentes, não a 44 implementações comprovadas.
- Alternativas: zero elegíveis, uma excluída. Como é autoanálise local, não mede descoberta de alternativas nem autoriza dizer que elas não existem.

Artefatos completos: ../benchmark/auditoria-produto-pos0064/dossie/. Comando, binário e hashes em receipt.json. A saída foi coletada em /tmp e copiada após a execução, para não inserir o próprio dossiê na coleta.

**Inferência:** o gargalo observado aqui é ligar requisitos a detectores/evidências pertinentes, além da cobertura de coleta. O algoritmo prudente sabe suspender a conclusão, mas muitas ações devolvem ao usuário a tarefa genérica de definir a verificação. Isso não prova ausência das capacidades no projeto nem demonstra que elevar limites resolva o problema.

## Desvio recente de foco

O wiring despacha effects separadamente em 04_wiring/main.go:66–71. O serviço principal é construído separadamente com expander lexical, descoberta e conteúdo. Não há consumo de EffectReport/EffectVerifier no fluxo de recommendation examinado. O avaliador lexical em evidence.go também preserva comportamento não verificado.

Portanto, 0062–0064 melhoraram um instrumento auxiliar real, mas não demonstraram melhora da decisão de reúso. Continuar ampliando esse piloto sem escolher uma capacidade da jornada pode aumentar a infraestrutura e os artefatos sem aproximar proporcionalmente o produto do objetivo.

A ponte necessária não é converter supported_in_scope em PROVEN genérico. Precisa preservar requisito, fonte exata, dependências, condições e limites; execução de um exemplo não comprova uma capacidade ampla ou viabilidade de integração.

## Prioridades sugeridas pela evidência

1. Voltar à jornada inteira: congelar um pequeno conjunto novo de ideias reais e respostas de engenharia esperadas antes de ajustar regras. Incluir alternativas úteis, semelhanças enganosas, restrições e casos inconclusivos. O corpus conhecido deve permanecer regressão.
2. Medir separadamente descoberta, recuperação da implementação, suporte a cada requisito e utilidade da decisão. Reportar falsos positivos, abstenções úteis e excessivas, alternativas relevantes e trabalho que ainda ficou para o usuário. Sem denominador misto de arquivo/documentação/comportamento.
3. Escolher uma capacidade estreita desses casos e completar uma cadeia útil: necessidade → candidato → trecho/implementação → evidência adequada → comparação → decisão limitada. Documentação pode sustentar declaração; execução só quando apropriada e dentro do recorte autorizado.
4. Fazer o plano de investigação apontar uma próxima ação concreta baseada na lacuna encontrada. Não criar detectores ad hoc para reconhecer o Bateia ou trocar a falta de prova por mais palavras no catálogo.
5. Só então expandir contratos de effects e integração com autofagia/Katamari quando a necessidade medida justificar. Preservar as salvaguardas e rastreabilidade já entregues.

Critério de êxito: o usuário consegue tomar uma decisão melhor fundamentada e identificar menos trabalho manual remanescente, em casos que não guiaram o ajuste do motor. Sem esse ensaio, uma nota global seria falsa precisão.

## Verificação desta análise

Build atual e execução local concluídos com código0. A suíte completa/lint aprovados no0064 são evidência da revisão anterior, não uma nova execução nesta análise. Não foram executados projetos externos, testes de integração de propostas ou descoberta pública. As conclusões arquiteturais são de leitura dirigida; não há atestação de independência. Tempos são telemetria da ferramenta. Sem commit.
