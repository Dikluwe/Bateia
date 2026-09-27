# Integração opt-in de rastreabilidade — 0098

Replay sem rede com e sem trace deve preservar bytes causais e RunID/inputHashes/
RuleVersions. Verificar artefato/estágio somente opt-in e hash de bytes; retomada
sem flag não consome/anuncia grafo antigo. Impacto EXPIRED solicita revisão sem
mutar decisão histórica/grafo. Decide sidecar mantém bytes do report e rejeita
destinos que identificam entradas por path/hardlink ou saída principal, antes
de escrever qualquer stdout. Permitir saída principal ainda inexistente.
0099: sem flag não produzir índice; com trace registrar índice ligado ao graph_hash,
mesmo número de entries/nodes e count locations, segundo output com hash dos bytes.
Resume sem flag preserva bytes dos dois sidecars e aviso aponta ambos. Estágio usa
projection-v2, não altera RunID/RuleVersions/checkpoint ou os bytes causais.
