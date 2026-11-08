# Oráculo de Integração do Passo 0034

**Camada**: L4
**Responsabilidade**:
- Validar as capacidades declaradas e execução determinística de todos os conectores do Passo 0034:
  - Google Play Store (`google-play-store-apps`)
  - F-Droid (`fdroid-foss-apps`) com extração de elo FOSS
  - Gitee (`gitee-repository-catalog`)
  - Codeberg (`codeberg-repository-catalog`)
  - GitLab (`gitlab-repository-catalog`)
- Validar a aceitação de todos os novos conectores no esquema MCP `bateia_start_market_discovery`.

## Evolução 0126

A composição federada inclui AtomGit; validar presença e unicidade das fontes pelo registro, sem congelar cardinalidade histórica. Testes de transporte específicos pertencem a L3.

## Passo 0129
Ownership aqui restrito à composição federada e Bitbucket. Asserções dos adaptadores foram transferidas para testes L3 de playstore, fdroid e forge, sem mudança de oráculo.
