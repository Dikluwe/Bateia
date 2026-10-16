# 0117 — Cenários determinísticos do modelo

Usar testing/synctest e canais/rendezvous, sem sleeps como oráculo. Testar:
configure antes de initialize é aceita; configure depois de initialize é
rejeitada e conserva valor; pausa depois do check, initialize, retoma write
aceita escrita após initialize no modelo; controle oposto initialize antes da
checagem rejeita e não chama hook. Checar estado apenas após sincronização.
Cada resultado emitir WITNESS JSON com cenário/aceitação/estado, scope MODEL_ONLY.
Nomear interleaving artificial e origem autoral. Não declarar teste independente,
execução de candidato ou garantia de proteção concorrente. Detector race deve
passar no modelo e repetição deve produzir mesmos witnesses, sem avaliar tempo.
