# Prompt proprietário — `03_infra/local/client_test.go`

Origem L0 exclusiva de `03_infra/local/client_test.go`. Este contrato foi separado de `00_nucleo/prompts/03_infra/cache/cache_store.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `writeFixture`, `TestClientAnalyzesLocalWorktree`, `TestClientRejectsPathTraversal`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/03_infra/cache/cache_store.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.
