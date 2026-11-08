# Bateia investigando melhorias para o XERJ — 2026-11-08

**Alvo:** XERJ. **Resultado:** três propostas fundamentadas, sem alteração no XERJ e sem efeito validado. A pergunta do usuário é o que o Bateia pode trazer de ideia para o XERJ. O Bateia permanece no workspace local.

Leia primeiro as [ideias e seus limites](ideias-para-xerj.md). O [registro estruturado](proposals.json) associa cada proposta a cartões, trechos, commits, hashes, condições e experimento de rejeição.

**Continuação — passo0149:** a [reprodução da herança G7](g7-reproduction/README.md) e o [teste de um agente usando somente XERJ](xerj-agent/README.md) acrescentam evidências às hipóteses iniciais, sem reescrever seus artefatos. A reprodução local confirmou herança na geração, com bloqueio independente de formato na publicação; não demonstrou ocorrência no catálogo público. O agente chegou a propostas comparáveis em duas direções, sem receber as anteriores.

**Jornada completa — passo0150:** a [comparação com dois agentes sem referências prontas](../jornada-completa-2026-11-08/README.md) inclui descoberta nativa e aquisição das fontes. Os dois percursos produziram propostas fundamentadas; a [ampliação em sete linguagens](../jornada-completa-2026-11-08/multilinguagem/README.md) é uma fase posterior, separada dessa comparação.

## O que o Bateia executou

1. `idea insights discover`: cinco consultas ao catálogo GitHub; todas responderam, com dez candidatos distintos. Quatro referências foram selecionadas; dois READMEs foram coletados e dois excederam o limite de conteúdo. O resultado é `PARTIAL` e preserva uma perspectiva sem referência selecionada. Não foram fornecidas credenciais nos comandos nem houve publicação no Hub.
2. `idea evidence collect`: leitura local de 17 arquivos selecionados do XERJ, em duas branches. O harness comparou os bytes a objetos Git nos commits fixados antes da coleta. As cópias temporárias contêm somente a seleção; completude dos bundles não significa cobertura integral do repositório.
3. `idea insights explore`: quatro bundles de entrada, 16 cartões, cinco tarefas de exploração, três identidades de repositório selecionadas. Quatro arquivos sofreram limite de prefixo de 1024 linhas. O resultado é `CONTEXT_PREPARED`, novidade `NOT_ASSESSED` e efeito `UNVERIFIED`.
4. A IA desta sessão interpretou o contexto, consultou trechos adicionais dos mesmos arquivos congelados e escreveu as propostas. O executável não produziu sozinho as conclusões. Não foram forçadas classes do catálogo fechado de analyze/compare.

Os dois bundles da descoberta permanecem incompletos porque cobrem somente README. A descoberta foi lexical e trouxe homônimos e candidatos pouco pertinentes; as duas fontes aproveitadas oferecem pistas documentais, não mecanismos executados nesta sessão.

## Fontes do alvo

| Branch | Commit | Recorte |
| --- | --- | --- |
| main | `1557de2f5be981d9db441969201c92a8af6d2de3` | Tipos do Hub, builder, assinatura/pacotes, lockfile, feedback, routing, arquitetura e Query-1 |
| corpus-hub | `708e974516a4ba1ba9bfa312aff164041997cfd6` | Contratos de curadoria, template, exemplo de avaliação, candidato e preservação de pedidos de autores |

As origens, hashes e links imutáveis estão em [sources.json](sources.json). C-ares e Feynman foram encontrados pelo Bateia; commits e bytes estão em [discovery.json](discovery.json). A origem consultiva do provedor local não é apresentada como autenticação Git: a conferência adicional pertence ao harness do Lab.

## Reprodução

Na raiz do Bateia:

```sh
GOCACHE=/tmp/bateia-go-cache go build -o /tmp/bateia-xerj ./04_wiring
```

O replay offline não precisa dos clones nem da rede:

```sh
python3 lab/bancada/insights/xerj-2026-11-08/verify.py --binary /tmp/bateia-xerj
```

Para repetir a coleta, obtenha checkouts públicos nos dois commits acima e passe seus diretórios ao harness. O script recusa HEAD divergente ou arquivo alterado. Use outro diretório para conservar o ensaio original:

```sh
python3 lab/bancada/insights/xerj-2026-11-08/investigate.py \
  --binary /tmp/bateia-xerj \
  --main /tmp/bateia-xerj-source-20261108 \
  --hub /tmp/bateia-xerj-hub-20261108 \
  --discovery lab/bancada/insights/xerj-2026-11-08/discovery.json \
  --output /tmp/bateia-xerj-repeticao
```

Para uma nova descoberta, execute a partir de um diretório sem `.env` do executável e grave outro resultado:

```sh
/tmp/bateia-xerj idea insights discover \
  --question 'provenance evidence verification reproducible software corpus' \
  --familiar 'BM25' --familiar 'vector search' --familiar 'commit pinning' \
  --max-per-query 2 --max-references 4 --timeout 90s \
  --output /tmp/bateia-xerj-nova-descoberta.json
```

A rede muda: a nova descoberta não promete os mesmos candidatos. O replay usa os bytes preservados.

## Artefatos e interpretação

| Arquivo | Papel |
| --- | --- |
| `discovery.json` | Consultas, candidatos, tentativas, READMEs e contexto inicial reais |
| `collection-plan.json`, `target-bundles.json` | Plano e coleta limitada do alvo; caminhos temporários do plano são recriados pelo harness |
| `sources.json`, `case.json` | Identidade das fontes, hashes e pergunta |
| `exploration-bundles.json`, `exploration.json` | Entrada e contexto produzidos pelo executável |
| `proposals.json`, `ideias-para-xerj.md` | Interpretação da IA; não resultado automático de avaliação semântica |
| `replay-check.json` | Integridade dos vínculos e igualdade do replay |
| `verification/` | Logs e códigos de saída da verificação do Bateia |

O gate canônico final `make verify` passou em lint, auditorias, suíte raiz, vet, build e diff-check. A primeira tentativa falhou por armazenamento temporário e sockets bloqueados no sandbox; os registros originais permanecem em `verification/`. A execução final usou `GOFLAGS=-p=2`, cache temporário e loopback permitido. Veja [os resultados finais](verification/verification.json) e [o contexto do ambiente](verification/environment.json). Gates selados com módulo próprio não foram executados.

No ensaio inicial0148, não executamos o XERJ nem testes dos projetos encontrados. A continuação0149 executou a release oficial e os scripts de pipeline, conforme o recorte próprio acima. Não enviamos feedback, email, issue ou PR. Os experimentos de utilidade das propostas continuam pendentes; não foi medido ganho de criatividade, precisão, tempo ou tokens.
