# Passo 0033 — Expandir descoberta global federada de mercados e locales

**Estado**: `PLANEJADO`  
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)  
**Depende de**: [`Passo 0027 — Generalizar descoberta multilíngue e multifuente`](0027-generalizar-descoberta-multilingue-e-multifuente.md), [`Passo 0028 — Extrair mecanismos e formular hipóteses de melhoria`](0028-extrair-mecanismos-e-formular-hipoteses-de-melhoria.md) e [`Passo 0032 — Executar teste end-to-end de validação do Bateia`](0032-executar-teste-end-to-end-de-validacao-do-bateia.md)  
**Baselines obrigatórios**: certificados `BATEIA_E2E_INTEGRATION_VERIFIED` e `MARKET_DISCOVERY_DURABLE_READY`  
**Onde roda**: validação e resolução conservadora em L1, adaptadores de catálogo em L3, esquemas MCP em L2 e injeção em L4  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: expansão da capacidade de observação geográfica e linguística do motor de mercado, preservando pureza, isolamento e controle estrito de orçamento

---

## 1. Problema

No Passo 0027, o Bateia introduziu o subsistema durável de **Market Discovery** (`bateia_start_market_discovery`) e o conector público de catálogo de aplicações (`apple-itunes-search-software`). Para garantir um baseline inicial e seguro, a validação de entrada em $L_1$ (`usecase/market_discovery.go`), os esquemas de validação do MCP em $L_2$ (`mcp/market_discovery.go`) e a implementação em $L_3$ (`itunes/market_discovery.go`) foram delimitados a um subconjunto rígido de **10 mercados piloto** (`BR, DE, ES, FR, IN, JP, KR, RU, SA, US`) e **10 locales canônicos** (`pt-BR, de-DE, es-ES, fr-FR, hi-IN, ja-JP, ko-KR, ru-RU, ar-SA, en-US`).

Essa restrição artificial cria uma barreira geográfica e de produto severa para qualquer investigação global:
1. **Invisibilidade de Mercados Globais:** Mais de 160 territórios da norma ISO 3166-1 alpha-2 com storefronts ativos (ex: Itália `IT`, Grécia `GR`, México `MX`, Argentina `AR`, Mongólia `MN`, Cazaquistão `KZ`, Nigéria `NG`, África do Sul `ZA`, Egito `EG`, Polônia `PL`, Suécia `SE`, Turquia `TR`, etc.) são sumariamente rejeitados com `ErrInvalidRunInput`.
2. **Invisibilidade Linguística e Cultural:** Idiomas locais e variantes regionais BCP-47 com centenas de milhões de falantes (ex: italiano `it-IT`, grego `el-GR`, polonês `pl-PL`, turco `tr-TR`, vietnamita `vi-VN`, tailandês `th-TH`, sueco `sv-SE`, árabe egípcio `ar-EG`, mongol `mn-MN`, cazaque `kk-KZ`) são bloqueados no enum fechado do esquema JSON Schema do MCP e no validador puro de $L_1$.
3. **Fragilidade em Lotes Multi-Território:** Ao investigar uma hipótese de software em múltiplos mercados simultâneos, a ausência de um catálogo normativo e de isolamento de falha por mercado faz com que um território indisponível, com catálogo reduzido ou vazio contamine o lote ou interrompa a auditoria de proveniência de toda a execução durável.

---

## 2. Incerteza que este passo reduz

Não sabemos se abrir a validação para o catálogo completo de territórios ISO 3166-1 e tags de idioma BCP-47 preserva a integridade, o determinismo e as cotas do Bateia sem:
- Degradar o tempo de execução do worker ou esgotar limites de memória (`maximumBodyBytes`);
- Sofrer injeção de caracteres de controle ou URLs canônicas maliciosas de storefronts internacionais não testados;
- Romper os esquemas fechados do MCP certificados no Passo 0028 (`o5blackbox`);
- Gerar corridas duráveis com falhas catastróficas quando um país específico retornar 404, rate limit ou payload vazio;
- Violar a Lei da Gravidade e a pureza estrita de $L_1$ (zero I/O).

---

## 3. Objetivo Arquitetural

Generalizar a validação, roteamento e extração de catálogos de mercado no Bateia para suporte global normativo:

1. **Normatização de Mercados em $L_1$:**
   - Adotar um validador canônico em $L_1$ baseado no padrão normativo **ISO 3166-1 alpha-2** (conjunto oficial de códigos de 2 letras maiúsculas), substituindo o mapa estático de 10 países.
   - Adotar um validador de locales conforme a especificação **BCP-47** (combinações de 2 ou 3 letras de idioma e opcionalmente 2 letras de território, ex: `^[a-z]{2,3}(-[A-Z]{2})?$`), garantindo que consultas em idiomas e grafias nativas não sejam rejeitadas arbitrariamente.

2. **Extensão do Conector em $L_3$ (`itunes` e futuros catálogos):**
   - Atualizar a verificação `validRequest` do conector de catálogo para aceitar qualquer código ISO 3166-1 suportado pelo endpoint público internacional.
   - Manter a sanitização rigorosa de texto (`safeCatalogText`), verificação de digest SHA-256 e parsing determinístico com `json.Number`.
   - Implementar resiliência por storefront: se um país específico não tiver catálogo de software ou retornar vazio (`ResultCount == 0`), a resposta deve emitir status `AttemptEmptyObserved` com evidência estruturada sem derrubar o restante do lote de mercados solicitados.

3. **Esquema Flexível e Seguro no Servidor MCP ($L_2$):**
   - Atualizar `marketStartSchema` em `02_shell/mcp/market_discovery.go`:
     - O campo `markets` passa a validar itens com padrão `^[A-Z]{2}$` (ISO 3166-1 alpha-2) com `maxItems: 50` por requisição (ou valor calibrado por orçamento).
     - O campo `locale` em `query_proposals` passa a validar tags BCP-47 válidas.
   - Garantir compatibilidade retroativa total com os oráculos do Passo 0028 (`o5blackbox` e testes de conformidade MCP).

4. **Isolamento e Pureza ($L_1$ & $L_4$):**
   - $L_1$ continua sem realizar qualquer chamada HTTP ou I/O.
   - A composição em $L_4$ injeta o registro completo de conectores e suas novas capacidades declaradas.

---

## 4. Matriz de Componentes e Alterações

| Camada | Arquivo / Módulo | Natureza da Modificação |
| :--- | :--- | :--- |
| **$L_0$ (Núcleo)** | `00_nucleo/prompts/01_core/usecase/market_discovery.md`<br>`00_nucleo/prompts/02_shell/mcp/market_discovery.md`<br>`00_nucleo/prompts/03_infra/itunes/market_discovery.md` | Atualização dos prompts espelhados contendo os contratos de validação ISO 3166-1 e BCP-47. |
| **$L_1$ (Core)** | `01_core/usecase/market_discovery.go`<br>`01_core/usecase/market_discovery_test.go` | Generalização dos validadores `marketLocales` e `marketCountries` para conformidade com padrões abertos; testes unitários cobrindo territórios de todos os continentes. |
| **$L_2$ (Shell)** | `02_shell/mcp/market_discovery.go`<br>`02_shell/mcp/market_discovery_test.go` | Atualização dos validadores JSON Schema de `bateia_start_market_discovery` mantendo os gates de segurança contra injeção e tipagem estrita. |
| **$L_3$ (Infra)** | `03_infra/itunes/market_discovery.go`<br>`03_infra/itunes/market_discovery_test.go` | Remoção do hardcode de 10 países; sanitização de URL para todos os domínios nacionais suportados (`apps.apple.com/<country>/app/...`); suíte de testes com mocks de respostas de múltiplos países. |
| **$L_4$ (Wiring)** | `04_wiring/market_discovery.go`<br>`04_wiring/step0033_global_market_test.go` | Oráculo integrado de validação multi-território determinístico sem I/O real. |

---

## 5. Critérios de Aceitação e Verificação

1. **Validação ISO 3166-1 e BCP-47:** O Bateia aceita requisições duráveis de prospecção para qualquer mercado ISO válido (ex: `IT`, `GR`, `KZ`, `MN`, `NG`, `ZA`, `PL`, `JP`, `BR`) e qualquer locale BCP-47 válido (ex: `it-IT`, `el-GR`, `kk-KZ`, `mn-MN`, `pt-BR`, `ja-JP`), sem lançar `ErrInvalidRunInput`.
2. **Rejeição Estrita de Entradas Inválidas:** Códigos de país inexistentes (ex: `XYZ`, `12`, `br_SP`) ou locales malformados continuam sendo sumariamente rejeitados na borda com erro estruturado.
3. **Resiliência de Lote:** Uma requisição com `markets: ["BR", "IT", "KZ", "GR"]` executa todas as sondagens de forma independente; falha ou ausência de dados em um mercado não corrompe os demais nem impede a geração do artefato final `market_snapshot`.
4. **Preservação de Gates Anteriores:** Os testes dos Passos 0027, 0028 e 0032 devem continuar passando 100% verdes (`go test ./...`).
5. **Linteiro Cristalino:** Execução de `crystalline-lint .` com zero violações.
6. **Emissão de Laudo:** Registro formal em `00_nucleo/laudos/0033-expandir-descoberta-global-federada-de-mercados-e-locales.md`.
