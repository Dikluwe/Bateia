# Laudo 0126 — Gitee e AtomGit

**Data:** 2026-10-18
**Contrato:** [passo 0126](../passos/0126-integrar-gitee-e-atomgit.md)

## Resultado

Gitee e AtomGit entram na composição padrão usada pelo executável MCP e na composição federada. Os schemas e a listagem derivam do registro; a seleção continua limitada a dez fontes. Gitee recebe token opcional por BATEIA_GITEE_TOKEN; AtomGit recebe BATEIA_ATOMGIT_TOKEN e reporta AUTH_REQUIRED sem rede se ausente.

O adaptador compartilhado publica somente identidade de catálogo, com commit vazio e vínculo não verificado, hash SHA-256 prefixado e instante UTC. Preserva os formatos públicos e distingue falhas de busca vazia. Flags públicas/privadas booleanas e inteiras, namespaces aninhados e URLs legadas GitCode são tratados explicitamente. Nenhum token foi lido de .env, registrado ou versionado nesta verificação.

A execução durável agora consulta fontes globais com orçamento. O plano e as lacunas preservam MARKET_NOT_FILTERED; a projeção não inventa país observado nos candidatos nem em sua proveniência. Isso também libera outros catálogos globais já registrados. O teste da porta em memória comprova despacho, orçamento e artefatos; testes de composição exercitam mercado e síntese offline.

## Verificação

- `go test ./...`: aprovado, com assinatura Git desativada somente no processo das fixtures e acesso ao cache/loopback. Nova rodada dos três pacotes alterados aprovada após ajuste final do prefixo do hash.
- `go vet ./...`: aprovado na versão final.
- `go build -o /tmp/bateia-0126 ./04_wiring`: aprovado. Destino temporário autorizado substitui /dev/shm.
- `git diff --check`: aprovado.
- `crystalline-lint --fail-on warning .`: bloqueado por incompatibilidade preexistente de `languages.go = true` com o binário instalado, que exige tabela. Não houve alteração da configuração do projeto.
- Diagnóstico adicional com cópia temporária da configuração, usando `go = { grammar = "tree-sitter-go", enabled = true }`: zero erros, 46 avisos V12 preexistentes e nenhum novo texto de diagnóstico frente ao HEAD original. Isso não equivale à aprovação do gate exigido. O lint instalado considera várias funções de composição/teste como lógica em L4; remediar essa compatibilidade é outro recorte.

Resultados brutos, configuração temporária e commit de comparação: [bancada 0126](../../lab/bancada/forjas/0126/verification.json).

## Limites

A consulta real sem credenciais ao Gitee retornou HTTP 200 e array vazio; demonstra resposta do endpoint, não qualidade ou cobertura. AtomGit retornou HTTP 401 sem credencial. Sucesso autenticado no AtomGit não foi validado remotamente; protocolo, parsing e falhas foram verificados por testes determinísticos. Nenhuma fonte foi promovida a evidência de código, hipótese sustentada ou autorização de reúso. Gates selados e hashes históricos não foram alterados.
