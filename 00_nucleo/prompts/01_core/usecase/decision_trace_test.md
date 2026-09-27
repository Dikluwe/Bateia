# Contrato independente do núcleo de rastreabilidade — 0098

Testar canonicalização determinística, não mutação de entradas, congelamento idempotente, rejeição de adulteração/hash/ciclos/arestas inválidas/duplicadas/IDs/kinds/scope/decisão inválidos e limites antes de marshal integral. Fontes conflitantes por conteúdo/snippet/força devem ser rejeitadas, enquanto content_hash ausente permanece uma declaração não autenticada.

Construção registra dependências sem promover documentação/PATH_ONLY, sem recalcular/alterar primary/status/decisão do relatório. Compor fixture com fontes em dois ramos independentes: alterar fonte A revisa somente seu ramo e a decisão registrada, não fonte/claim B. REFUTES também propaga impacto. Duas alterações convergentes produzem tarefas com ambos ancestrais ordenados e resultados estáveis sob permutação.

Verificar noop por CONTENT_CHANGED com fingerprint idêntico, EXPIRED/REVOKED mesmo fingerprint exigindo revisão, schema/unknown/duplicatechange/reason/fingerprint inválidos, alterações de claims/decision, mudança policy/input conservadora. Nenhuma execução de candidatos/rede ou mutação histórica é permitida.

Fixture de amplificação: grafo válido de 100 fontes convergindo a um nó intermediário, 700 alegações descendentes e uma decisão, dentro de 2 MiB de entrada. Alterar as 100 fontes expande ancestrais e deve retornar erro controlado da quota de saída, artefato vazio, antes de marshal integral do plano. Alterar somente uma fonte do mesmo grafo permanece controle positivo válido.
