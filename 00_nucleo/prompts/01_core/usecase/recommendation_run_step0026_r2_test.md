# Prompt R2 — oráculos ampliados do coordenador de runs

## Causa

Refinamento pré-candidato autorizado pelo recibo discriminatório R1 e pelos casos M01–M60. A implementação candidata permanece proibida.

## Testemunhos adicionados

- os sete campos-conjunto são permutados/deduplicados isolada e conjuntamente, e o par base/reordenado larga concorrentemente;
- cada campo semântico participa do hash, enquanto a chave operacional não participa;
- retry retorna lifecycle/version corrente sem reservar, criar ou executar novamente;
- matriz exaustiva de chave e pedido inválidos não toca porta alguma;
- todas as operações sobre fixtures terminais com a matriz íntegra exigida são no-op byte-estável e views filtram `PENDING`;
- progresso percorre os oito marcos e rejeita duplicação, regressão, limites e texto livre;
- dezesseis cancelamentos concorrentes produzem uma transição e um sinal;
- a ordem pública de artefatos é a allowlist normativa, não ordenação lexical acidental.

Fakes continuam puramente em memória e sincronizados; nenhum I/O concreto entra em L1.
