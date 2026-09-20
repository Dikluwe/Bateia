# CLI idea recommend adversarial tests

Ataque replay sem rede, determinismo, cadeia de hashes, criação do dossiê,
validação antes de provedores e serialização de segredos. Os testes usam somente
arquivos temporários e contratos públicos; não executam código de candidatos.

## Contrato 0040

Correspondências lexicais, declarações e testes vazios são evidência parcial, não comprovação comportamental. As expectativas anteriores de PROVEN e componentes derivados apenas desses sinais foram substituídas. Preservar verificações de fontes, localização, determinismo, licença e rejeição; identidade de execução muda com a versão de regras.

## Migração da identidade 0043 por testador separado
Zero locales preserva evidence hash e bytes de bundles, mas a identidade da execução incorpora versões decision=0043, alternatives=0043 e analyzer=0043 para impedir reúso de decisão antiga. Expansion/screening seguem 1 e evidence segue DefaultEvidenceRules().Version. Antes da execução, congelar expectativa por fórmula pública: JSON compacto Go struct com campos na ordem evidence_hash, policy_hash, target_hash, rule_versions (map keys lexicais); SHA256 hexadecimal, primeiros 16 caracteres. As entradas são hashes evidence/policy/target do manifesto; comparar RuleVersions com mapa contratual explícito. Não copiar o runID observado. Preservar golden do evidenceHash, hash dos bytes de bundles e ausência de search_locales; repetir replay e exigir a mesma identidade. Os casos de mudança de locale permanecem intactos.

## Modelo de screening explícito

Verificar que `--screening-model` aceita o modelo declarativo em
`00_nucleo/modelos/screening-operacional-v1.json`, registra a identidade de seus bytes,
usa a versão `screening-opt-in-v1` e não concede `DROP` ou probabilidade a modelo não
calibrado. Campos desconhecidos, estado calibrado e conteúdo JSON adicional devem ser
rejeitados antes da criação do diretório de saída. Alterar os bytes do modelo deve
alterar a identidade da execução.
