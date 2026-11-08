# alternative_target

**Camada:** L1. **Contrato:** passo 0130 e ADR-0013.

Possuir remapeamento de alvo local antes do ranking, sem I/O concreto. Expor AlternativeInputsForTarget e a porta AlternativeIdentityReader.RepositoryIdentities(string). Reproduzir modos/guardas e ordem das consultas atuais, copiar candidatos/assessments/fontes antes de modificar; hints vazios marcam Stale, identidade não correspondente não remapeia, somente o alvo correspondente recebe owner/repo canônicos em todas as referências já tratadas. Metadata é consultiva. Implementação não autentica origem nem altera RankSoftwareAlternatives.
