# 0132 — Diagnóstico de retomada

Contrato ex-ante: [passo 0132](../../../../00_nucleo/passos/0132-explicar-retomada-de-evidencias.md). Prioridade escolhida pelo usuário: confiabilidade e clareza dos resultados.

## Reprodução

Na raiz do projeto:

```sh
GOCACHE=/tmp/bateia-review-go-cache go test ./01_core/usecase -run 'Test(EvidenceCheckpoint|CheckpointDiagnostic|RecommendationResume|RecommendationDossier)' -count=1
GOCACHE=/tmp/bateia-review-go-cache go test ./02_shell/cli -run 'Test(ResumeDiagnostic|IdeaRecommendResumeDiagnostic)' -count=1
BATEIA_VERIFY_OUTPUT="$PWD/lab/bancada/confiabilidade/0132/final" GOCACHE=/tmp/bateia-review-go-cache make verify
```

Os testes de descoberta usam provedores em memória; replay usa evidência autoral, sem APIs remotas. O gate mantém a configuração de assinatura Git apenas no processo dos testes e registra comandos/saídas. Gates selados conservam sua reprodução independente.

## Evidência

- `baseline.json`: hashes anteriores à alteração, incluindo oráculos existentes.
- `lineage.json`: hashes das sete materializações e de seus prompts; os oráculos anteriores permanecem intactos.
- `core-tests.log`: política anterior e diagnóstico/jornada em memória.
- `final/verification.json` e logs correspondentes: pilha completa.
- `sandbox-attempt/`: primeira pilha, com duas fixtures HTTP bloqueadas por `socket: operation not permitted`; lint, auditorias, vet, build e diff passaram.
- `loopback-tests.log`: os mesmos dois pacotes passaram com loopback autorizado, sem mudança de teste ou produto para contornar o ambiente.
- `verification.json`: resultado consolidado, preservação e limites.

O teste CLI confirma diagnóstico compatível com trabalho efetivo: mudança só de política não repete descoberta/coleta; hash adulterado impede reutilização e causa nova descoberta/coleta. A matriz replay distingue entrada, hash, artefato ilegível, versão e manifesto ausente/malformado. A mensagem não afirma coleta nova em replay. O teste puro mantém ordem, primeira duplicata não vazia, short-circuit e primeiro impedimento, sem expor erro/valores. O retorno de uma jornada que falha depois da avaliação conserva o diagnóstico sem anunciar conclusão.

## Observação do ensaio e limites

Uma primeira comparação de bytes entre execução inicial e retomada encontrou `requirement_coverage: null` versus `[]` no bundle da fixture sem cobertura. O comportamento decorre do `MarshalJSON` já existente em `01_core/domain/evidence.go`, que este passo não altera. A saída dessa tentativa está em `artifact-comparison-attempt.log`. A comparação de bytes final usa uma cobertura registrada não vazia; os testes históricos com cobertura vazia continuam intactos. Este passo não certifica idempotência universal da serialização nem corrige essa normalização legada.

O diagnóstico é transitório: não modifica schema dos artefatos, identidade ou versões. É apresentado pela CLI após publicação bem-sucedida e antes do cleanup; falhas não recebem anúncio de conclusão. MANIFEST_UNAVAILABLE cobre leitura/interpretação, sem presumir arquivo ausente. Somente o primeiro impedimento é relatado, sem leituras adicionais; não é inventário de todos os defeitos, autenticação de origem, garantia contra mudanças concorrentes ou promoção de evidência/hipótese. Não houve avaliação com usuários ou provedores públicos.
