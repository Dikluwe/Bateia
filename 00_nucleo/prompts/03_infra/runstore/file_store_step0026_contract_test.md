# Prompt espelhado — oráculo do FileStore 0026

## Autoridade

Teste independente do adaptador L3, derivado do Passo 0026 e contrato R1. Não pressupõe layout privado: usa apenas a porta pública e inspeções adversariais genéricas da raiz temporária.

## Obrigações

- `CreateOrResolve` atômico, CAS/versionamento, idempotência e reopen durável.
- Permissões privadas, rejeição de symlink e hardlink inesperado.
- Publicação confirmada com hash do documento completo, tamanho, contagem e imutabilidade.
- Paginação repetível; cursor HMAC preso a run/nome/versão/hash e independente de `limit`.
- `PENDING`/ausente não vira lista vazia; quota é verificada sem overshoot residual.
- Corrupção é `INTEGRITY_ERROR`; recovery interrompe órfãos e não inventa resultado.

## Ambiente

Todos os casos usam `t.TempDir`, segredo fixo de teste e zero rede. Atomicidade além das garantias do filesystem exercitado permanece caso opaco `Unknown`.
