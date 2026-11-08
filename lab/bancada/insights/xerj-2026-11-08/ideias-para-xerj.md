# Ideias para o XERJ a partir da investigação com o Bateia

O resultado mais concreto é **vincular cada avaliação do Corpus Hub à versão exata avaliada**. As outras duas propostas são oferecer um recibo do que foi verificado e entregar condições de aplicação junto à referência. São hipóteses desta sessão, com fontes conferíveis; nenhuma foi implementada ou validada no XERJ.

**Evidência posterior:** a [continuação0149](g7-reproduction/README.md) reproduziu a associaçãoG7 após re-pin em fixture local e identificou limites do caminho de publicação. O [agente usando XERJ](xerj-agent/README.md) chegou independentemente a direções semelhantes de recibo e aplicabilidade, com a mesma pergunta dirigida. As seções abaixo conservam a interpretação inicial; os novos relatos qualificam suas conclusões.

## 1. Uma nota de busca deve dizer a que versão se refere

O XERJ já verifica pacotes assinados, registra o commit e oferece `xerj.lock`. Não falta fixação de versão do código.

O ponto diferente está na avaliação G7. No recorte examinado, `read_graded_suites` lê as avaliações e devolve notas indexadas pelo nome do corpus. Arquivos são processados em ordem lexical e uma ocorrência posterior substitui a anterior. O builder associa essa nota ao corpus pelo nome. Esse leitor não exige no recibo o commit avaliado, o hash do pacote ou o protocolo da consulta.

Isso permite formular uma pergunta concreta: **quando o conteúdo muda, como garantir que a nota continua descrevendo aquele conteúdo?** A inspeção estática não prova que o Hub publicou uma associação incorreta; outra etapa do processo pode impor essa relação.

**Proposta:** um recibo de G7 com identidade do corpus, commit/hash do pacote, versão do indexador, hash das consultas, data e responsável. A nota atual só usa um recibo compatível; avaliações anteriores continuam disponíveis como históricas. É uma aplicação do vínculo entre evidência e versão usado no Bateia.

**Experimento:** construir dois pacotes com o mesmo nome e commits diferentes, mantendo apenas a avaliação do primeiro. Conferir se o catálogo do segundo herda a nota. Depois fornecer avaliação compatível. Rejeitar ou reduzir a proposta se o pipeline completo já impedir a associação por outro contrato comprovado.

Fontes: [leitura das avaliações](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-hub/src/build_catalogue.rs#L418-L494), [tipos do catálogo](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-common/src/hub.rs#L141-L171). A evidência congelada e a associação no builder estão em `XERJ-01` de [proposals.json](proposals.json). O cartão automático do arquivo apontou um teste; a análise do leitor foi uma leitura adicional da IA sobre o arquivo já coletado.

## 2. Mostrar o que foi verificado em cada solução

G7 avalia relevância dos trechos recuperados. Assinatura verifica a origem/integridade do pacote. Feedback registra experiência de uso. O código examinado já distingue alguns desses escopos: por exemplo, um download bem-sucedido não recebe automaticamente uma nota de utilidade.

Para um agente recebendo uma solução, ainda seria útil distinguir **“o teste existe”, “o upstream declarou que passou” e “alguém executou esse teste neste ambiente”**. Os tipos `Answer` e `Hit` examinados não carregam esse recibo funcional. Isso é uma observação do formato, não prova de que ninguém executa testes na curadoria.

**Proposta:** um anexo opcional por contribuição com propriedade avaliada, origem da avaliação, versão, comando, ambiente e resultado localizado. Ausência fica como não executado/desconhecido; não vira falha nem validação. A extensão requer versionamento explícito, pois os tipos atuais recusam campos desconhecidos. Não é necessário executar código recebido durante a indexação.

O Bateia encontrou o Feynman, cujo README separa auditoria, plano de replicação e proveniência. Esse é um precedente documental para a separação; não o executamos para confirmar o comportamento descrito. O c-ares também documenta um recibo de atestação de release, útil como analogia de vinculação da evidência ao artefato, sem equivaler a prova funcional.

**Experimento:** dez referências com escopos de evidência diferentes, comparando cegamente agentes com e sem recibo. Medir afirmações indevidas de validação, precisão da escolha e custo de contexto. Rejeitar a proposta se não reduzir a confusão ou se o custo de manutenção superar o benefício.

Fontes: [gates do Hub](https://github.com/xerj-org/xerj/blob/708e974516a4ba1ba9bfa312aff164041997cfd6/tools/xerj-code/hub/PROGRAM-100.md#L26-L42), [formato da resposta](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-common/src/hub.rs#L280-L315), cartões de Feynman e c-ares em [exploration.json](exploration.json). Registro `XERJ-02` em [proposals.json](proposals.json).

## 3. Entregar também quando a solução não serve

O XERJ já pode recuperar testes e documentação. A proposta é relacioná-los explicitamente à solução oferecida: propriedade, precondições, teste e contraexemplo, com localizadores independentes.

Por exemplo, duas filas podem ter código parecido e contratos diferentes: uma bloqueia quando cheia; outra rejeita a submissão. Um trecho pertinente à palavra “fila” não basta para escolher entre elas. Esse é um exemplo hipotético de avaliação, não um comportamento atribuído ao XERJ ou a um projeto coletado.

**Proposta:** um modo opcional de referência contextual que agrupe implementação, teste especificador e condições de aplicação. Relações sugeridas por IA são identificadas como declarações até revisão. Quando não há contraexemplo conhecido, registrar a lacuna e uma pergunta de investigação. Encontrar o teste não significa que ele passou.

É a contribuição mais direta do fluxo de exploração do Bateia: as tarefas pedem inversão de premissas, recombinação e contraexemplos, em vez de tratar toda referência encontrada como candidata pronta para adoção.

**Experimento:** cinco problemas com implementações parecidas e contratos diferentes. Comparar recuperação atual com o pacote contextual, medindo escolhas e adaptações incorretas, explicação das restrições e tokens adicionais. Rejeitar a proposta se os agentes já resolvem essas distinções com as passagens atuais ou se o contexto extra não ajuda.

Fontes: tarefas e instruções de [exploration.json](exploration.json), formato de `Hit` e gates relacionados em `XERJ-03` de [proposals.json](proposals.json). O esquema de aplicabilidade é proposto; sua ausência em todo o ecossistema XERJ não foi demonstrada.

## Limites da investigação

A seleção automática produziu passagens úteis e outras genéricas, inclusive cabeçalhos e código de teste do próprio Hub. Foi necessário ler mais do arquivo congelado para chegar ao primeiro achado. Isso também evidencia trabalho restante no Bateia: preparar contexto não elimina interpretação especializada.

Não propomos ao XERJ “adicionar busca híbrida”, “fixar commits”, “assinar pacotes”, “ter memória” ou “coletar feedback”: essas capacidades já aparecem no projeto. Query-1 continua uma superfície distinta de consulta sobre fatos; não presumimos que resolva aplicabilidade de código ou execução de testes de uma contribuição.

O próximo passo de maior valor é reproduzir `XERJ-01` no pipeline completo. Só depois caberia apresentar um defeito confirmado ou materializar uma mudança. As outras duas propostas precisam de um pequeno ensaio de produto antes de ampliar esquemas e curadoria.
