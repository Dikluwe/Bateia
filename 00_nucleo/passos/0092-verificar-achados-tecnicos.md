# Passo 0092 — Verificar achados técnicos antes de usá-los

## Escopo

Achados livres continuam consultivos. Dois predicados estruturados agora têm
verificadores determinísticos sobre `EvidenceBundle` coletado na revisão
fixada do candidato:

- `LICENSE_ID`, com técnica `LICENSE_COMPATIBILITY`: confere SHA-256 dos bytes,
  revisão, caminho raiz `LICENSE`/`COPYING` e identificador detectado pelo
  avaliador de licença existente. **Identifica a licença; não determina
  compatibilidade jurídica de uma composição.**
- `SBOM_COMPONENT`, com técnica `SBOM_SCA`: confere os mesmos vínculos e lê um
  manifesto JSON CycloneDX para confirmar `name@version`. **Confirma apenas a
  declaração no SBOM; não a presença no binário nem explorabilidade de falha.**

Cada achado estruturado usa `predicate`, `value`, `source.ref` (caminho),
`source.revision` (commit fixado) e `source.digest` (`sha256:<hex>`). O bundle
deve estar completo, sem falhas, e declarar a mesma revisão no plano, no
relatório e no item. Achados vencidos, contraditos, sem hash válido, com
conteúdo adulterado ou com técnica incompatível não são promovidos. A
verificação estabelece **consistência interna dos bytes coletados**; um JSON
forjado não autentica por si a origem remota. A coleta deve vir de um provedor
confiável do funil.

## Comandos

`bateia idea evidence verify --input achados.json --bundles bundles.json --as-of 2026-09-24T00:00:00Z --output revisao.json`

`bateia idea decide --idea ideia.json --assessment assessment.json --policy reuse-policy.json --technical-evidence achados.json --technical-bundles bundles.json --evidence-as-of 2026-09-24T00:00:00Z`

`--technical-bundles` é opcional, mas exige `--technical-evidence`. Sem ele,
o comportamento do passo 0091 permanece consultivo. Com ele, apenas um
`LICENSE_ID` verificado pode preencher um gate de licença antes
`INCONCLUSIVE`; licença já comprovada nunca é substituída. O gate ainda
aplica `allowed_licenses` da política. SBOM e demais técnicas não mudam a
decisão nesta etapa.

## Validação

Testes cobrem identificação positiva, digest alterado, revisão divergente,
bundle incompleto, licença ambígua, SBOM declarativo, contradição por
predicado, vencimento e mudança de `INCONCLUSIVE` para `FORK` somente com
licença verificada. A busca simples permanece inalterada.
