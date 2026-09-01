# Laudo 0022 — Calibração de regras de evidência

**Data**: 2026-08-31
**Estado**: `IMPLEMENTADO`
**Veredito**: `PROMOTE_EVIDENCE_CATALOG`
**Regime Tekt**: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

## Resultado

O catálogo avançou para o schema 2 e passou a declarar família, proveniência,
alvos justificáveis e necessidade de fontes independentes. Regras por identidade,
termos genéricos, duplicações entre famílias e hash divergente são rejeitados.
Versões anteriores continuam legíveis para replay.

O corpus congelado contém 4 casos `CONTRACT`, 3 `VALIDATION` e 3 `TEST`. Precisão
e recall de `PROVEN` foram 100% em todos os splits; não houve promoção opaca,
fonte inválida ou perda sem diagnóstico. Fontes `PROVEN` são verificadas contra
conteúdo presente no bundle do mesmo `repo@commit` e caminho — não por um booleano
autodeclarado.

As 12 mutações registradas foram rejeitadas, produzindo mutation score 1,0. A
suíte completa, `go vet`, `staticcheck 2026.2.1`, Crystalline Lint e diff check
passaram.

## Limite

Os papéis foram separados logicamente e os artefatos são reproduzíveis, mas o
ambiente compartilhado não comprova isolamento forte do split `TEST`. O laudo
atesta somente o fragmento observável deste corpus offline e não generalização
para outros domínios.
