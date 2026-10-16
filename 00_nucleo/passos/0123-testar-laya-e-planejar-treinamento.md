# 0123 — Testar Laya e planejar treinamento

Autorização: executar diagnósticos do modelo e pensar treinamento. Sob Tekt A0,
instalação isolada em /tmp, código Laya no commit da pesquisa 0122, revisão dos
pesos resolvida e fixada antes de inferência. Registrar versões, hardware/dtype,
hashes, orçamento de tokens, respostas brutas, tempo e desvios dos quatro
rótulos diagnósticos autorais. Não instalar em ambiente global nem promover
modelo ao produto. Não treinar nesta fatia antes de referência adequada.

Executar checkpoint inglês em CPU sobre os quatro pedidos conhecidos; registrar
repetição e permutações das opções sem ajustar a pergunta após resultados.
Conferir entrada não truncada. Casos não são holdout nem medida de precisão
semântica independente; não foram escolhidos para representar toda a tarefa.

Plano de treinamento: obrigação estreita de relação alegação/declaração,
categorias exclusivas, fonte/versão/condições e referência independente.
Dividir por projeto/família antes de gerar paráfrases; reservar teste final e
calibração, medir contradições/negações/insuficiência e falso suporte. Treinar
somente após conjunto suficiente, mantendo abstenção/cobertura e identificação
do modelo separadas da autoridade factual. Resultados no Lab, síntese em L0.

Verificar hashes e limites, lint, testes raiz, vet e build. Produto/observadores
não alterados; gates selados não reproduzidos.

Executado em CPU com pesos reais fixados: três concordâncias em quatro casos
autorais; ordens/repetição estáveis em 28 execuções correlacionadas. Não houve
truncamento. Protocolo de treino preparado, sem treino/calibração ou conjunto
novo rotulado. [Laudo](../laudos/0123-teste-real-do-laya-e-plano-de-treino.md).
