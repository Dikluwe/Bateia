# Contrato de evidência

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.

## Passo 0037
Aplicar o contrato de resposta por mecanismo em `00_nucleo/passos/0037-resposta-fundamentada-por-mecanismo.md`. Contexto único de compilação, trechos exatos, limites explícitos e estados desconhecidos preservados.

## Passo 0038
Aplicar contrato do protótipo AST+CFG: adaptador opcional explicitamente configurado, hashes e limites, fluxo intraprocedural e estados não avaliados.
