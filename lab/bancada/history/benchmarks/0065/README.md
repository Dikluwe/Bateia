# Experimento 0065 — Coalizões em modo sombra

Contrato autoritativo em
`../../passos/0065-analisar-coalizoes-de-evidencia-em-modo-sombra.md`; recorte canônico
em `contract.md`. Esta pasta nasce antes de qualquer implementação Go, teste executado,
oráculo ou mutante.

O experimento mede classificação conservadora de coalizões explicitamente anotadas,
preservação de conflito/Unknown, determinismo, limites e ausência de mudança em
vereditos baseline. Não mede independência causal real nem autoriza integração do
resultado na decisão do Bateia.

`baseline.json` identifica a implementação existente relevante; `manifest.json`
registra hashes, observáveis e capacidades; `budget.json` fixa condições de parada. O
autor de oráculos deve criar corpus e mutantes depois do congelamento e antes do
implementador. O verificador produzirá selo/recibos sem editar entradas protegidas.

Estado atual: contrato candidato, ainda não selado nem executado. Regime completo como
ensaio, filesystem compartilhado e sem atestação física de isolamento.

