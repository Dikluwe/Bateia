# screening_model_policy

**Camada:** L1. **Contrato:** passo 0130 e ADR-0013.

Expor IsValidOptInScreeningModel(domain.ScreeningModel) bool com o contrato exato do perfil opt-in da CLI: schema 1, status uncalibrated, calibrador nil, Weights não nil, thresholds/taxa válidos e todos os valores numéricos finitos. Não desserializar JSON nem alterar ScreenDiscovery, defaults ou perfis calibrados. Não impor restrições novas a nomes de sinais, priors, seed ou pesos negativos.
