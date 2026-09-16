# Contrato independente — catálogo fornecido no fluxo integrado, passo 0071

Testar `idea recommend` sem rede, a partir de replay congelado. `--rules` deve ser lido e
validado antes de provedores, injetado no pipeline/serviço e refletido em versão, hash,
identidade causal, assessment e `evidence/catalog-coverage.json`. O catálogo default não
pode substituir silenciosamente o fornecido.

Catálogo inválido deve falhar antes de qualquer chamada externa. Repetição e reordenação
do catálogo produzem artifacts semanticamente e byte a byte estáveis. Não ler
implementação candidata ou ajustar fixtures ao patch.
