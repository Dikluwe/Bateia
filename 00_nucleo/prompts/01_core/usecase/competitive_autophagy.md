# Prompt — motor determinístico de autofagia competitiva e Katamari

Implemente um caso de uso puro e offline que avalie os bundles internos pelas
regras existentes, classifique cada capacidade como `KEEP`, `CONSOLIDATE`, `GAP`
ou `UNKNOWN`, e enumere a menor composição de concorrentes capaz de cobrir apenas
gaps comprovados. Licença, operação, arquivo, falhas e commit imutável são gates,
não pesos compensáveis. Ordene por cobertura, quantidade de projetos, custo,
redundância e identidade estável. Preserve lacunas e rejeições explicitamente.

## Revisão 0041

Aplicar `00_nucleo/passos/0041-reestruturar-autofagia-e-katamari.md`. Saída versão 2 distingue inventário, oportunidades e cobertura de integração validada. CONSOLIDATE não pode ser inferido por contagem. Não estimar custo sem dados. Preservar contratos antigos como campos compatíveis, com semântica conservadora explícita.
