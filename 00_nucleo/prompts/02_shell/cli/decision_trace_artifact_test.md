# Integração opt-in de rastreabilidade — 0098

Replay sem rede com e sem trace deve preservar bytes causais e RunID/inputHashes/
RuleVersions. Verificar artefato/estágio somente opt-in e hash de bytes; retomada
sem flag não consome/anuncia grafo antigo. Impacto EXPIRED solicita revisão sem
mutar decisão histórica/grafo. Decide sidecar mantém bytes do report e rejeita
destinos que identificam entradas por path/hardlink ou saída principal, antes
de escrever qualquer stdout. Permitir saída principal ainda inexistente.
