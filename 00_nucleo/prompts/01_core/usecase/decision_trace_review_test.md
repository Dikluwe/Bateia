# 0101 — testes do registro consultivo

Testar estados/default, contagens, sem tarefa diferente de revisão completa,
COMPLETED sem aprovação, vinculação hashes/nó, mudança de motivos/new_fingerprint,
canonicidade e não mutação. Recusar notas/evidência insuficientes, hashes ou estados
inválidos, registros repetidos/desconhecidos, hash grafo adulterado, esquemas,
UTF8 e quotas agregadas, 1000/1001 registros e 20/21 referências. Guardar saída
expandida e retornar zero em erro. Testes puros L1, sem I/O/rede/relógio.
