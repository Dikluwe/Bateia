# Prompt R2 — oráculos ampliados do FileStore

## Causa

Refinamento pré-candidato contra os sobreviventes do recibo discriminatório R1. O layout privado não é congelado: inspeções descobrem genericamente arquivos/componentes sob uma raiz temporária.

## Testemunhos adicionados

- `CreateOrResolve` e CAS reais competem atrás de barreiras;
- limites e segredo são validados em todas as fronteiras;
- matriz terminal completa permanece byte-estável após reopen;
- payload/MAC truncado, alterado, forjado ou com padding é recusado; cursor válido continua após reopen e troca de `limit`;
- nomes externos não viram componentes de caminho e chave clara não persiste;
- symlink, hardlink, FIFO e sentinela externa são exercitados black-box;
- arquivos superiores não confirmados nunca viram artefato publicado; eles não são apresentados como fault injection de crash;
- publicações concorrentes respeitam quota e preservam runs anteriores.

Garantias de power-loss/rename do dispositivo continuam exclusivamente `Unknown` nos casos O01/O02.
