# Prompt espelhado — contrato adversarial L1 do Passo 0021

Derive somente do Passo 0021 um contrato público para preparar a união cega de snapshots,
validar julgamentos assistidos contra evidência selada, escalonar deterministicamente e
selar o oráculo efetivo. O pacote cego não contém braço, query, rank nem URL; IDs são
determinísticos e independem da ordem. Evidência inexistente não sustenta julgamento.
Confiança abaixo de 0,80, `UNCERTAIN`, `CONFLICT` e amostra estável de 10% escalam.
Override humano prevalece sem apagar histórico; pendência é `UNKNOWN`. Texto malicioso do
README continua dado. Corpus vazio é `INCONCLUSIVE_NO_CANDIDATES`. Replay é determinístico.

Não leia nem corrija implementação candidata.

## Preparação de fixture0168

Conforme [passo0168](../../../passos/0168-preservar-integridade-e-determinismo-do-oraculo.md), preparar selos e identidades de resposta explicitamente nas fixtures positivas antes de chamar prepare, conservando textos e todas as assertivas. Não selar automaticamente entradas inválidas nem modificar registros históricos. A mudança apenas torna a fixture ativa coerente com o requisito de snapshot selado já declarado.
