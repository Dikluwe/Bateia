# Market discovery wiring

**Camada:** L4  
**Arquivo proprietário:** `04_wiring/market_discovery.go`

Componha os dois conectores repo-backed basais, o único conector Apple, replay,
registry, serviço L1, dispatcher e tools L2. Injete EXATAMENTE as mesmas instâncias
de store, fila, gerador de run ID, raiz e cursor secret de 0026; preserve o mesmo
`RecommendationService`. Apenas composição/configuração: nenhuma validação,
budget, hash, merge, transição, retry, parsing ou regra de negócio em L4.


## Admissão por registro no passo 0104

Delegar a admissão não vazia e a validação de capacidades/IDs ao construtor L1 NewMarketDiscoveryRegistry. Não impor cardinalidade fixa em L4. Aceitar conjuntos válidos unitários ou ampliados, mantendo a mesma configuração e lifecycle; recusar configuração inválida antes de criar recursos.

## Gitee e AtomGit — 0126

Incluir ambos nos conjuntos padrão e federado, preservando os três conectores basais no início. Ler BATEIA_GITEE_TOKEN (opcional) e BATEIA_ATOMGIT_TOKEN (exigido pelo provedor para busca), passando-os aos adaptadores L3. Registrar AtomGit mesmo sem token para tornar a capacidade visível; a falta de credencial é reportada pelo adaptador ao consultar. Não alterar o teto L1 de dez fontes selecionadas por pedido.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
