# Resolução das fontes congeladas da síntese

**Camada:** L1
**Proprietário:** `01_core/usecase/improvement_sources.go`
**Contrato:** passo 0127, especialização do protocolo 0028 e corte 0103.

Possuir resolvedSynthesisSources e as funções de resolução, leitura paginada, consulta do índice publicado e validação de referências de propostas. Receber uma porta interna synthesisArtifactReader com somente Get e ReadArtifact. Proibido depender de RecommendationRunService, fila, IDs, transições ou publicação; efeitos concretos continuam em L3.

Extrair o comportamento vigente de improvement_synthesis.go sem redesenhar validações: conferir identidade, tipo de run, terminalidade e hash do índice; ler páginas de 50 até 1000, conferir RunID/Artifact/ArtifactHash em todas as páginas; validar JSON, copiar bytes e recusar cursor vazio/repetido. Preservar erros das portas. Indexar mecanismos sem duplicidade, gaps e observações; marcar como degradada fonte terminal sem SUCCEEDED. Validar referências de propostas contra o estado resolvido.

Não normalizar pedidos ou calcular hash causal, não derivar hipóteses e não gerenciar lifecycle. Os dois consumidores no protocolo durável usam a mesma função de resolução. Preservar enums e estruturas públicas; a porta e o estado deste componente são internos ao pacote.
