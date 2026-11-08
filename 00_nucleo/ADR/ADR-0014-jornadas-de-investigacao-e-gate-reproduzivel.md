# ADR-0014 — Jornadas de investigação e verificação reproduzível

**Status:** ACEITO PARA O RECORTE 0131. **Data:** 2026-10-28.

L1 possui a jornada do dossiê, identidade causal, estágios, manifesto e decisões de execução, usando portas de leitura/publicação/captura. L2 trata flags, strict JSON, mensagens e renderização; seus adaptadores não escolhem política de negócio. L4 injeta a fachada e as dependências, inclusive quando regras de evidência são entradas dinâmicas. O caminho compatível da API CLI delega à mesma jornada L1.

A investigação web possui admissão, intercalação, orçamento, retomada e veredicto em L1. O resultado do experimento não comprova benefício ou adoção. A investigação de mercado implementa a porta ObservableMarketDiscoveryExecutor em L1; o runtime durável conserva admissão operacional, CAS, lifecycle e publicação. Não há runtime paralelo nem efeitos concretos em L1.

Testes unitários seguem seus proprietários; verificadores gerais e evidência bruta ficam no Lab. A fonte e o patch fixados do linter são dependência de desenvolvimento, isolados do produto e reproduzíveis pelo gate do projeto. Mudanças no verificador preservam regras/severidades e controles negativos.

Os formatos, hashes, versões e políticas anteriores continuam vigentes. Esta decisão substitui ownership histórico de orquestração nas cascas apenas para as capacidades contratadas no passo 0131.
