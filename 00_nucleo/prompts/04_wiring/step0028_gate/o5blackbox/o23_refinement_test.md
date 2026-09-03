# O23 — opacidade estrutural binding-free do protocolo O22

Materialize a tabela pública O22 somente para `ARTIFACT_READ` de
`hypothesis_gaps`. Exija documento canônico de uma entrada, ordem fixa
`hypothesis_id,code,source_ref`, identidade sintaticamente canônica, par exato
de código controlado e relação `source_ref`, além das relações de integridade e
publicação. O valor da identidade e os valores de hash nunca controlam a
classificação; somente seus formatos e relações contratuais são validados.

Os self-tests leem o protocolo público O22, exercitam os sete pares e alteram
um requisito por vez. Toda violação continua precedendo `Unknown`; controles
vizinhos permanecem `Preserved`. Não leia material adversarial executável,
records, subjects, configs, candidato, produção, receipt, seal ou veredito.
