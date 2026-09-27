## Passo 0095 — coerência de bundles

Aplicar orçamento de API pura antes do hash; rejeitar bundles repetidos da mesma
identidade/revisão. Schema desconhecido não pode gerar fonte. Unicidade de licença
conta paths raiz independentemente de Kind. Normalizar conflito LICENSE_ID com a
mesma equivalência do verificador. Chaves estruturadas sem colisões NUL. SBOM passa
por verificador limitado 1.6 e guarda JSON embutida, sem promoção além de declaração.

## Passo 0094 — garantias distintas

Separar autenticidade da origem, integridade dos bytes, base da interpretação e
comportamento. Offline origem não estabelecida; hash válido não prova semântica.
Promover licença somente com integridade e FULL_TEXT_TEMPLATE. Não transformar
achados importados ou declarações SBOM em prova de execução segura.

# Verificação determinística de achados técnicos

Ligar achado a bundle completo na revisão fixada, recomputar digest dos bytes,
rejeitar ambiguidade, vencimento e contradição. Reusar o avaliador de licença
existente para `LICENSE_ID`; conferir declaração CycloneDX para
`SBOM_COMPONENT`. Somente licença verificada pode preencher gate antes
inconclusivo; não substituir evidência já comprovada, inferir compatibilidade
jurídica ou presença no artefato executado.
