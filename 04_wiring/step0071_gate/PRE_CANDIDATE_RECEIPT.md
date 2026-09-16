# Recibo pré-candidato — passo 0071

- Regime: protocolo completo, fase adversarial pré-candidato.
- Papel: adversário; escrita limitada a `04_wiring/step0071_gate/`.
- Implementação candidata: não lida e ainda ausente.
- Comando: `04_wiring/step0071_gate/run.sh`.
- Instante causal: após congelamento dos contratos 0071 e antes da implementação.
- Resultado esperado observado: **FAIL pré-candidato** (`exit 1`).
- Inputs congelados: todos os 21 registros de `frozen-inputs.sha256` retornaram `OK`.
- Causa dominante: símbolos públicos 0071 ainda não materializados.
- Primeiro diagnóstico: `undefined: domain.EvidenceRuleBasis`.
- Diagnósticos seguintes: campo `CapabilityRule.Basis`, tipos de cobertura e constantes
  `CatalogCoverage*` ausentes.
- Mutation score: não executável antes do candidato; selo não emitido.
- Política de Unknown: somente casos `O*` aceitam Unknown; nenhum default é sucesso.
- Determinismo exigido: vinte permutações, igualdade byte a byte e imutabilidade dos inputs.

Hashes dos artefatos executados estão em `artifacts.sha256`. Este recibo registra falha
pré-candidato, não conformidade funcional, selo ou certificado final.

