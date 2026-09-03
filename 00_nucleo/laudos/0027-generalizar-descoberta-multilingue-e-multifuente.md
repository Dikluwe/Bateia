# Laudo do Passo 0027 — descoberta multilíngue e multifuente

**Veredito:** `MULTISOURCE_MARKET_DISCOVERY_READY`  
**Atestação:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`  
**Contrato ativo:** R4, SHA-256 `44980c809a12d74eafe919b88220f7a47f0987d648e6f4edfa0ea77310619725`  
**Selo ativo:** `seal-r4.json`, SHA-256 `a8dd85c3aa567f31249e1906195f14b47dec4bbd7bdebbd8e678a97492b22bc4`  
**Certificado:** SHA-256 `acf05df728b05ccdc5bc15f78a18e49f19e743b2016ab1141e553ca2417c7e43`

## Conclusão

O candidato do Passo 0027 satisfaz o fragmento observável congelado para descoberta de produtos por consultas multilíngues e mercados independentes, com um catálogo público Apple não baseado em repositório, preservação das superfícies MCP e duráveis dos Passos 0025/0026, replay offline determinístico e projeções públicas conservadoras.

O resultado se limita aos hashes, ao conector `apple-itunes-search-software` versão `itunes-search-doc-1.0.3-response-v2`, à matriz de dez mercados, ao corpus, às fixtures e ao runtime local single-writer registrados. Não estabelece equivalência funcional geral nem disponibilidade futura do provedor.

## Cadeia causal e refinamento R4

O R4 substituiu exclusivamente o alvo inviável de build de arquivo isolado `./04_wiring/main.go` pelo build do pacote Go `./04_wiring`. A reconstrução em memória da linha anterior recuperou exatamente o contrato R3 SHA-256 `ac69fcaf3668097432d36137477f791ec0e95c193689fbbe32b1442960e6c037`. Os demais 22 inputs protegidos auditados não mudaram.

- invalidação de binding R4: `8268a9a8d6f893f6c62de5b873101652fa848b9c29ea7b4301aeed277ff5cac7`;
- receipt discriminatório R4: `954d4d1a29cc17e9c0745988eb513a0c02e11a744d22bad1b175e8d52641dee8`;
- selo R4: `a8dd85c3aa567f31249e1906195f14b47dec4bbd7bdebbd8e678a97492b22bc4`;
- receipt do candidato: `115750f24d17e67e9970ac3945677d402a9011eb03c38e8bb3703228cd917d0a`;
- receipt de verificação: `cbd88c68cdabb5b1683dc73e7d546c561939b0f3557c525d4ca72697c89930c8`.

Os selos R1, R2 e R3 permanecem preservados como história causal invalidada; nenhum foi sobrescrito.

## Gates executados

| Gate | Resultado |
|---|---|
| JSON estrito e hashes protegidos | PASS |
| `go test -count=1 ./...`, execução 1 pós-R4 | PASS |
| `go test -count=1 ./...`, execução 2 pós-R4 | PASS |
| race de L1, MCP, L3 e L4 | PASS |
| `go vet ./...` | PASS |
| build literal corrigido `./04_wiring` | PASS; binário `62d3c995a280ee7b2d3a68d789a13913c1d28a0580ec57ff55870e88cc967c29` |
| `crystalline-lint .` | PASS com avisos não bloqueantes registrados |
| `git diff --check` | PASS |
| três suítes de oráculo protegidas | PASS |
| MCP compilado real via stdio | PASS; nove ferramentas e três conectores canônicos |
| replay offline repetido | PASS; zero rede e bytes idênticos |
| execução pública Apple limitada | PASS; uma chamada, um mercado, um resultado, sem retry |
| refinamento do Passo 0025 | PRESERVED |
| refinamento do Passo 0026 | PRESERVED |

O replay produziu SHA-256 `fc37bcd625e493b312a2b5e2820308f03233562ebf32d8c91984053f371edc4f` nas duas serializações. A execução pública Apple retornou HTTP 200 e um resultado strict JSON; o body não retido teve SHA-256 `5ce95b9d989be1dec5415c4153ca9a9c89022d3f73513cda0e51e8be2c2bb5ed`. Nenhum segredo ou conteúdo bruto proibido foi registrado.

## Discriminação e métricas

- mutações negativas válidas: 90;
- classificadas `Violated`: 90;
- sobreviventes: 0;
- mutation score: `1.0`;
- controles positivos `Preserved`: 5;
- casos opacos `Unknown`: 5, excluídos do denominador.

A exposição observada foi 100%. Precisão top 10 por família, recall de `VALIDATION+TEST`, precisão de merge, preservação de produto sem repositório e identidade byte a byte do replay foram 100%. O ganho marginal Apple sobre o baseline repo-backed foi 100 pontos percentuais no registro congelado, com candidato incremental relevante em 6/6 ideias. Inferência dimensional indevida e influência de popularidade, ratings ou downloads foram zero.

## Arquitetura e segurança

L1 permaneceu sem rede, filesystem, ambiente, goroutine ou mutex; L2 não importou implementações concretas L3/L4; L3 não criou dependências ascendentes; L4 reutilizou a única instância de store e queue duráveis. O adaptador Apple restringiu método, esquema, host, porta, path, parâmetros, tamanho e redirecionamento; falhas externas permaneceram `UNKNOWN_PROVIDER`, e cancelamento impediu novas chamadas.

As dimensões `query_locale`, `listing_locale`, `observed_market` e `publisher_country` permaneceram independentes. O conector Apple não inventou listing locale ou país do publisher, não converteu descrição em capacidade provada e preservou produtos sem repository link.

## Limitações e invalidação

O trabalho ocorreu em workspace compartilhado, e R3/R4 são posteriores à existência do candidato. Portanto, a atestação não reivindica isolamento temporal forte. A consulta pública comprova apenas viabilidade no instante registrado; não comprova disponibilidade futura, qualidade semântica de traduções, idioma de listing, origem do publisher ou instalabilidade.

Qualquer alteração em hash causal, input protegido, candidato, predecessor, conector, política pública, matriz, corpus, fixture, testemunha, política de `Unknown` ou resultado de gate invalida este laudo e o certificado.
