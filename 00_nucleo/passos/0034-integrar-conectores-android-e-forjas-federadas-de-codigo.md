# Passo 0034 — Integrar Conectores Android (Google Play & F-Droid) e Forjas Federadas de Código Aberto (Gitee, Codeberg, GitLab)

**Estado**: `IMPLEMENTADO`  
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)  
**Depende de**: [`Passo 0027 — Generalizar descoberta multilíngue e multifuente`](0027-generalizar-descoberta-multilingue-e-multifuente.md), [`Passo 0028 — Extrair mecanismos e formular hipóteses de melhoria`](0028-extrair-mecanismos-e-formular-hipoteses-de-melhoria.md) e [`Passo 0033 — Expandir descoberta global federada de mercados e locales`](0033-expandir-descoberta-global-federada-de-mercados-e-locales.md)  
**Baselines obrigatórios**: certificados `GLOBAL_MARKET_DISCOVERY_VERIFIED`, `MARKET_DISCOVERY_DURABLE_READY` e `BATEIA_E2E_INTEGRATION_VERIFIED`  
**Onde roda**: L1 (domínio e interfaces puras), L2 (esquemas e ferramentas MCP), L3 (adaptadores de rede e parsers de catálogo/forjas) e L4 (composição e injeção)  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: expansão multissistema do Bateia para cobrir o ecossistema móvel Android (comercial e FOSS) e forjas soberanas/regionais de código aberto, preservando determinismo, rastreabilidade causal e isolamento de I/O.

---

## 1. Problema

Historicamente, os subsistemas de descoberta do Bateia possuíam dois pontos cegos estruturais:

1. **Monocultura Móvel Apple (iOS/macOS):**
   No Passo 0027 e Passo 0033, a descoberta de mercado foi consolidada através do conector `apple-itunes-search-software`. Embora o iTunes ofereça cobertura em 170+ territórios, ele reflete apenas o ecossistema iOS. A plataforma **Android detém mais de 70% da participação de mercado global** (e mais de 85% em nações da América Latina, Ásia e África). Ignorar o ecossistema Android distorce a análise de mercado de software, omitindo soluções dominantes em países emergentes.

2. **Desconexão entre Aplicativo Móvel e Código Aberto Upstream:**
   Usuários e agentes que buscam software muitas vezes encontram uma solução móvel, mas não conseguem comprovar se ela é de código aberto ou reutilizável. O ecossistema Android possui no **F-Droid** um repositório formal onde 100% dos aplicativos são FOSS e contêm metadados explícitos vinculando o pacote compilado (`applicationId`) ao repositório de código upstream no GitHub, GitLab ou Codeberg. Essa ponte causal não existia no Bateia.

3. **Centralização Monolítica no GitHub:**
   A descoberta de código e mecanismos de arquitetura no Bateia (`03_infra/github/`) assume o GitHub como fonte quase exclusiva. No entanto, por razões geopolíticas, de soberania digital ou de preferência filosófica por software livre estrito:
   - Na **China e Leste Asiático**, grande volume de código aberto e projetos comunitários reside no **Gitee** (`gitee.com`);
   - Na **Europa e na comunidade FOSS estrita**, o **Codeberg** (`codeberg.org`, baseado em Forgejo) é amplamente adotado como alternativa sem fins lucrativos e sem telemetria corporativa;
   - Em **infraestruturas científicas e corporativas**, instâncias federadas de **GitLab** abrigam projetos estruturantes de domínio não indexados pelo GitHub.

---

## 2. Incerteza que este passo reduz

1. **Viabilidade e Resiliência da Google Play Store:**
   Não sabemos como consultar a Google Play Store de forma sustentável, idempotente e sem custos proibitivos de API privada, garantindo a extração de pacotes Android por geolocalização (`gl=ISO-3166`) e idioma (`hl=BCP-47`) sem ser bloqueado por sistemas antibot da Google.

2. **Consistência do Elo Causal F-Droid → Código Upstream:**
   Não sabemos se o catálogo estruturado do F-Droid (`index-v2.json` / pacotes) pode ser ingerido pelo Bateia sob controle estrito de memória e cota de I/O em $L_3$, gerando automaticamente candidatos para o Funil de Redução de 3 Estágios sem intervenção manual.

3. **Generalização da Interface de Forjas de Código:**
   Não sabemos se a interface abstrata de busca de código (`SearchProvider` / `ProjectCatalog`) em $L_1$ acomoda as variações de payload e contratos de paginação da Gitee API v5, Codeberg API v1 e GitLab API v4 sem vazamento de tipos concretos de infraestrutura para o núcleo puro.

---

## 3. Objetivo Arquitetural

Implementar a arquitetura federada de plataformas móveis e forjas no reticulado de 5 camadas:

```text
L4 (Wiring)  ──> Injeta conectores Google Play, F-Droid, Gitee, Codeberg e GitLab
L2 (Shell)   ──> Expõe novas ferramentas e opções no MCP (bateia_start_market_discovery)
L3 (Infra)   ──> Implementa clientes HTTP seguros com rate-limiting, cache e replay
L1 (Core)    ──> Interfaces puras: MarketCatalogProvider, FOSSLinkage, ForgeSearchProvider
L0 (Nucleus) ──> Prompts e especificações espelhadas com linhagem causal estrita
```

### 3.1. Conector Google Play Store (`google-play-apps`)
- **Localização:** `03_infra/playstore/`
- **Capacidades:**
  - Suporte aos parâmetros canônicos globais `gl` (ISO 3166-1 alpha-2) e `hl` (BCP-47).
  - Extração de: `packageId` (ex: `com.example.app`), título, desenvolvedor, score/avaliações, descrição resumida e categoria.
  - Sanitização de texto (`safeCatalogText`), conformidade com orçamento observável de tempo/bytes e fallback gracioso em caso de país sem catálogo ou restrição de distribuição.

### 3.2. Conector F-Droid FOSS Android (`fdroid-foss-apps`)
- **Localização:** `03_infra/fdroid/`
- **Capacidades:**
  - Consulta ao índice estruturado do F-Droid com cache local por digest SHA-256.
  - Rastreamento causal: mapeia `packageId` diretamente para `SourceCodeURL`, `License` e `IssueTracker`.
  - Ponte de autofagia: todo resultado do F-Droid alimenta diretamente a entrada do estágio de Triagem do Bateia para verificação de código em $L_1$.

### 3.3. Forjas Federadas de Código Aberto
- **Gitee (`03_infra/gitee/`):** Conector para a API pública v5 do Gitee, permitindo buscas em língua chinesa (`zh-CN`) e repositórios hospedados na China.
- **Codeberg (`03_infra/codeberg/`):** Conector para a API Forgejo/Gitea do Codeberg, permitindo buscas em repositórios da comunidade de software livre europeia.
- **GitLab (`03_infra/gitlab/`):** Conector para busca de repositórios públicos via GitLab API v4.

---

## 4. Matriz de Componentes e Camadas

| Camada | Caminho do Componente | Responsabilidade e Restrições |
| :--- | :--- | :--- |
| **$L_0$** | `00_nucleo/prompts/03_infra/playstore/market_discovery.md` | Especificação formal do conector Google Play Store. |
| **$L_0$** | `00_nucleo/prompts/03_infra/fdroid/market_discovery.md` | Especificação formal do conector F-Droid e elo causal FOSS. |
| **$L_0$** | `00_nucleo/prompts/03_infra/forge/federated_search.md` | Especificação dos adaptadores de forjas (Gitee, Codeberg, GitLab). |
| **$L_1$** | `01_core/domain/market_types.go` | Novos tipos canônicos: `PlatformType` (`iOS`, `Android`, `Web`), `FOSSLinkage` (URL upstream, licença). |
| **$L_1$** | `01_core/domain/interfaces.go` | Interfaces `MarketCatalogProvider` e `ForgeProvider` (ZERO I/O). |
| **$L_1$** | `01_core/usecase/market_discovery.go` | Validação de provedores múltiplos e agregação federada de resultados. |
| **$L_2$** | `02_shell/mcp/market_discovery.go` | Atualização do JSON Schema das ferramentas MCP com catálogo de provedores expandido. |
| **$L_3$** | `03_infra/playstore/market_discovery.go` | Implementação do conector Google Play com parsing seguro e cache. |
| **$L_3$** | `03_infra/fdroid/market_discovery.go` | Implementação do conector F-Droid com resolução de código-fonte. |
| **$L_3$** | `03_infra/forge/gitee_client.go` | Cliente de busca Gitee v5. |
| **$L_3$** | `03_infra/forge/codeberg_client.go` | Cliente de busca Codeberg (Forgejo). |
| **$L_3$** | `03_infra/forge/gitlab_client.go` | Cliente de busca GitLab v4. |
| **$L_4$** | `04_wiring/main.go` | Registro e injeção de todos os novos provedores no composition root. |
| **$L_4$** | `04_wiring/step0034_android_and_forges_test.go` | Oráculo de integração e conformidade de todos os conectores. |

---

## 5. Critérios de Aceite e Verificação

1. **Pureza Absoluta de $L_1$:**
   - `01_core/` não importa net/http, disco ou bibliotecas de terceiros.
   - Testes unitários em $L_1$ executam 100% offline em memória com mocks puros.

2. **Resiliência e Isolamento em $L_3$:**
   - Nenhuma falha em uma forja ou storefront (ex: 404, rate limit, timeout) pode derrubar o lote ou corromper a execução durável.
   - Todas as respostas de rede passam por validação de tamanho máximo de payload (`maximumBodyBytes`) e digest SHA-256.

3. **Oráculo de Teste End-to-End ($L_4$):**
   - O teste `TestStep0034AndroidAndForgesIntegration` deve validar com sucesso:
     a) Descoberta na Google Play Store com parâmetros `gl=BR` e `hl=pt-BR`;
     b) Descoberta no F-Droid extraindo o link de repositório de código-fonte aberto;
     c) Descoberta federada de código no Codeberg e Gitee;
     d) Agregação multi-provedor no MCP sem quebra de esquema JSON.

4. **Verificação de Linteiro:**
   - Passagem limpa no `crystalline-lint .` com cabeçalhos obrigatórios de linhagem causal em todos os arquivos Go.
   - `go test ./...` passando com 100% de sucesso.
   - Compilação do binário Bateia sem avisos ou erros.

---

## 6. Certificado Conquistado ao Final

A implementação deste passo concede ao repositório o certificado:
```text
CERTIFICATE: BATEIA_FEDERATED_MOBILE_AND_FORGES_READY
Hash: SHA256(Passo-0034-Lattice-Artifacts)
```
Declarando o Bateia formalmente capaz de auditar mercados móveis (iOS e Android) e ecossistemas de código aberto em qualquer território ou forja do planeta.
