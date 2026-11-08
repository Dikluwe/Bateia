# Entradas de arquivo limitadas

Ler arquivos selecionados explicitamente pelo usuário, recusando dispositivos,
diretórios, FIFOs e mais de 16 MiB. Conferir descritor e limitar leitura
durante crescimento. Para leitores de decisões, bundles e intake, limitar
profundidade a 64 e elementos a 100000 antes de desserializar. Caminhos
explicitamente selecionados não são o confinamento de repositório remoto.

O descritor é aberto O_NONBLOCK antes de conferir arquivo regular: FIFOs não bloqueiam a admissão. readIdeaDescriptor permite reutilizar a mesma identidade já aberta em radar. A guarda JSON delega ao parser.ValidateBoundedJSON compartilhado, incluindo UTF8, exatamente um valor, 16 MiB, profundidade64, tokens100000 e duplicatas/aliases casefold. Não normalizar bytes usados por hashes históricos.

O teto é por arquivo; não representa orçamento agregado para árvores de múltiplos arquivos. Entrada explicitamente selecionada não é confinamento de repositório.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
