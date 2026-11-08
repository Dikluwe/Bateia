# Prompt Espelhado — 04_wiring/step0033_global_market_test.go
## Passo 0033: Expandir Descoberta Global Federada de Mercados e Locales

### Propósito
Oráculo de teste integrado (L4) que valida as mudanças introduzidas pelo Passo 0033:
substituição dos enums fechados de 10 países e 10 locales por validadores normativos
ISO 3166-1 alpha-2 e BCP-47 em L1, L2 e L3.

### Contrato de Teste
1. `TestStep0033ISO3166AndBCP47ValidatorsAcceptGlobalMarkets`: confirma que os validadores
   de L1 aceitam qualquer código ISO 3166-1 alpha-2 válido e qualquer tag BCP-47 no
   subconjunto "xx" ou "xx-YY".
2. `TestStep0033InvalidInputsAreRejected`: confirma que entradas malformadas continuam
   sendo rejeitadas (lowercase, underscore, 3+ chars, numérico, etc.).
3. `TestStep0033MCPSchemaUsesPatternNotEnum`: confirma que o registry L1 aceita conectores
   com os novos campos globais.
4. `TestStep0033MultiBatchWithEmptyStorefrontDoesNotFail`: valida resiliência de lote —
   um mercado sem catálogo (AttemptEmptyObserved) não impede os demais.
5. `TestStep0033ItunesConnectorAcceptsGlobalMarkets`: confirma que o conector iTunes L3
   aceita mercados globais via validador itunesValidISO3166Alpha2.

### Restrições
- Zero I/O real — todos os testes usam stubs e mocks.
- Sem dependências externas além de L1 (usecase) e L3 (itunes).
- Conforme Lei da Gravidade: L4 importa L1 e L3, nunca L2.

## Passo 0129
O teste iTunes e seu transporte pertencem agora ao teste do adaptador L3. Aqui permanecem contratos de composição de domínio/MCP.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.

## Consolidação 0131
Os dois testes puros de validadores foram transferidos, sem mudar esperados, para 01_core/usecase/global_market_validation_test.go. Aqui permanecem schema MCP e integração de composição.
