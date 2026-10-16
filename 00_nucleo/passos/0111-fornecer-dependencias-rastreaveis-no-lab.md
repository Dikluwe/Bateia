# 0111 — Fornecer dependências rastreáveis no Lab

Sob Tekt A0, capacidade investigativa no Lab. Preparação com IO separada da
transformação pura de observação; nenhum import produtivo ou nova dependência.
Prompts de ambiente, resolução, observador, CLI e testes precederam código.

Nove fontes congeladas de conc e 64 exports da biblioteca padrão local são
fornecidos por dados em memória. Importer sem fallback, validação de SHA256,
limites, binding, ciclos e falhas. environment_sha256 identifica o ambiente;
fontes conservam bytes e commit do coletor 0106. Contexto de build permanece
limitado, sem alegação de build completo do repositório.

Com dependências fornecidas, checagem de tipos sem diagnósticos e limiter
reconhecido em make(limiter,n). Remover dependências suprime crédito e adulterar
hash impede publicação. Qualidade semântica independente e fluxo global ainda
pendentes. [Bancada e uso](../../lab/bancada/concorrencia/0111/README.md).
