# Prompt — controles documentais independentes O8

Derive controles somente do contrato R4, catálogo v3 e protocolo O8 públicos.
Uma regra `DOCUMENTARY` deve ligar literalmente o catálogo/hash, exigir source
kind `documentation`, produzir witness kind `DOCUMENTATION` e materializar no
máximo `PARTIAL`. Ela jamais conta como implementação ou testemunha independente
de `PROVEN`, mesmo sob status pai `PROVEN`. Não use ID, record, marker, expected,
payload adversarial ou implementação candidata.

No refinamento O14, o self-test documental calcula `mechanism_id` com o ID do
registro `MechanismRule` do catálogo. O RuleID documental continua em
`extraction_rule_ids` e no witness, mas não ocupa o campo `rule_id` do preimage.
