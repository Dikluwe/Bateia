# Passo 0094 — Endurecer segurança e garantias

## Correções

O corpus adversarial da revisão mostrou seis falhas. Este passo implementa
controles determinísticos e regressões permanentes, preservando busca simples
e mantendo o modo avançado opcional.

1. **GitHub sem clone automático.** Contents API e raw continuam disponíveis,
   mas não há subprocesso Git, checkout, hooks, filtros ou leitura local de
   arquivos do candidato. A extração exige arquivos Git regulares e rejeita
   symlinks/submodules. Indisponibilidade dos dois canais produz falha
   explícita, nunca tentativa de acesso local ou certeza fabricada.
2. **Respostas limitadas antes da decodificação.** O provedor limita arquivo a
   2 MiB, envelope de conteúdo a 4 MiB, árvore a 8 MiB e metadados a 1 MiB.
   Extração aceita até cem arquivos e 8 MiB de conteúdo total. Redirecionamentos
   não são seguidos. Falhas de segurança não provocam fallback mais permissivo.
3. **Relatórios HTML offline.** Conteúdo Mermaid é textual, sem CDN/JavaScript.
   Há escapes separados para labels, HTML e atributos, política de URL e CSP
   defensiva. Links perigosos viram texto. Markdown permanece dependente de
   renderer que desabilite HTML bruto e use Mermaid em modo estrito.
4. **Um despacho MCP.** Extensões registram handlers; nenhuma consome o
   transporte antes do guarda de orçamento/contexto. Todas as chamadas usam
   o mesmo contador. Framing limita mensagens a 64 KiB antes de bufferização;
   limites estruturais são 32 níveis e 8192 tokens. Mensagem oversized encerra
   a sessão sem drenar fluxo ilimitado. Reader genérico bloqueado não oferece
   cancelamento de leitura portável; cancelamento inicial e operações
   cooperativas continuam sendo respeitados.
5. **Entradas CLI limitadas.** Os leitores de entradas de ideia protegidos
   recusam arquivos não regulares e mais de 16 MiB. Decisão, bundles e intake
   também têm orçamento estrutural de JSON. Isso não transforma qualquer
   caminho escolhido pelo usuário numa sandbox de arquivos.
6. **Licença conservadora.** Apenas texto MIT completo reconhecido pelo
   template normalizado pode preencher o gate automaticamente nesta etapa.
   Marcadores, textos modificados, citações e cláusulas extras são
   inconclusivos. Cabeçalho copyright tem gramática fechada: ano/range e um
   identificador ASCII de titular, sem cauda livre. Titulares multiword ou
   não ASCII ficam inconclusivos, um custo deliberado contra ressalvas ocultas.
   Apache, MPL e GPL permanecem sinais consultivos até existirem
   verificadores de texto completo. Coleta incompleta e múltiplas licenças
   raiz também impedem promoção. Não há análise jurídica de composição.

## Melhorias importadas

- Cada verificação separa `source_authentication`, `integrity_verified`,
  `interpretation_basis` e `behavior_verified`. Offline, a origem é
  `NOT_ESTABLISHED`; releitura pode registrar `PROVIDER_OBSERVED`. Esses campos
  são observações, não assinatura de exportações nem prova de titularidade.
- Relatórios técnicos declaram `UNTRUSTED_DATA_NO_AUTHORITY`. Achados e
  documentação não mudam permissões. O indicador é descritivo; controles
  determinísticos de gate e executor continuam sendo a barreira real.
- Investigação registra `security_risk` e justificativa. Entre opções
  pivotais, leitura limitada precede execução de terceiros; ganho de
  informação não autoriza execução. Executor de terceiros continua indisponível.
- O scanner heurístico expõe `NOT_SCANNED`, `SCANNED_NO_MATCH` ou
  `VULNERABILITY_FOUND`, com escopo e data da base. Sem correspondência não
  significa seguro. A base local é parcial e não resolve transitivos.

## Scanner opcional de binários Go

`bateia security scan --binary /caminho/absoluto/binario-go --output scan.json`
consulta `govulncheck` previamente instalado. Nunca instala ferramentas,
recompila o projeto, executa o candidato ou aceita modo source. O despacho
antecede `.env`, caches e clientes de busca. O scanner instalado é uma
dependência confiável do host e pode consultar `https://vuln.go.dev`.

O artefato regular de até 256 MiB é copiado para snapshot privado; SHA-256
identifica os bytes analisados. Leitura verifica identidade e rejeita arquivos
não regulares, inclusive troca para FIFO. Scanner não pode ser o próprio
artefato. Ambiente mínimo não herda tokens, proxies ou configuração Go.
Prazo padrão é um minuto, máximo cinco; stdout limitado a 8 MiB e stderr a
64 KiB. Saída/JSON parcial ou falha resultam em `NOT_SCANNED`.

JSON é consultivo (`ADVISORY_ONLY`), com scanner, base, data quando disponível,
versão de protocolo e escopo. `SCANNED_NO_MATCH` significa apenas nenhuma
correspondência na análise concluída, não ausência de vulnerabilidades.
Achados binários não estabelecem alcançabilidade/explorabilidade em runtime.
Arquivo de saída não pode ser alias do binário; publicação é atômica. Testes
usam scanners sintéticos offline, incluindo excesso de saída, timeout,
protocolo inválido, ambiente, symlinks/FIFO e proteção contra sobrescrita.

## Corpus de regressão

Testes de GitHub usam transporte offline e fixtures Git locais com symlinks,
PATH/configuração hostis, respostas grandes, caminhos adulterados e redirects.
Testes de relatório exercitam fechamento de tags, atributos, URLs perigosas
e labels Mermaid. MCP cobre combinações de extensões, budget compartilhado,
mensagens com/sem LF, profundidade, tokens e contexto. Core cobre licença
negada, reconhecimento incompleto, garantias separadas e risco de investigação.

As fixtures positivas antigas que continham apenas duas frases de licença
foram substituídas por texto MIT completo: o objetivo positivo permanece, mas
os testes deixam de exigir a promoção insegura que a auditoria demonstrou.
Oráculos de hashes do MCP foram supersedidos explicitamente nos prompts;
schemas/goldens e certificados históricos não foram reescritos.

## Limites

Este passo não certifica ausência de vulnerabilidades. O executor `effects`
continua sem sandbox rígida e não deve receber código hostil. Não adicionamos
execução automática de candidatos, instalação automática de scanners ou
autoridade derivada de texto de repositórios. Exportações JSON permanecem
adulteráveis sem armazenamento/autenticação externos confiáveis.

## Validação em 2026-09-26

- `go test ./...` e `go vet ./...` passaram.
- `go test -race` passou para usecase, CLI, MCP, report, GitHub e vulnerability;
  corpus de licença/técnica foi repetido após fechar a gramática de copyright.
- `git diff --check` passou.
- Contraprovas originais agora respeitam budget MCP, rejeitam frame de 128 KiB,
  escapam script em HTML e recusam promoção de licença negada.
- Despacho CLI real, sem govulncheck instalado, retornou NOT_SCANNED e não
  instalou ferramentas. Scans positivos/falhos foram testados com executáveis
  sintéticos offline; não houve varredura real na base remota nesta validação.
- Hash de código congelado no gate R4 também foi supersedido explicitamente;
  certificados históricos, corpus de ataques e controles causais preservados.
