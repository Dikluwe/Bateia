# Integridade do produtor web, retomada e união cega

Materializar testes puros para `ExecuteWebExperiment`, `ReplayWebSnapshot` e
`PrepareBlindCandidateUnion`, por suas APIs públicas. Usar runtime e provedor em
memória, sem rede ou arquivos. As expansões possuem 12 queries por braço N/F1/F2,
fornecidas em ordem inversa, autorização pública explícita e margem segura.

A saída real autenticada deve conter snapshots schema2 replayáveis e seus 36 hits
públicos completos devem entrar no pacote cego schema1. O legado schema1 íntegro
continua admitido. Schema ausente/desconhecido, adulteração, item privado ou campo
material ausente produz ErrIntegrity e pacote zero; não alterar os inputs.

Na retomada, validar todos os snapshots antes de buscar qualquer query faltante.
Exercitar um braço inicial parcial e um snapshot final inválido; variar RunID,
braço, provedor, grammar/query-set hashes, autorização, orçamento e schema.
Recalcular os selos dos adversários sem corrigir seu conteúdo: query, IDs,
MaxResults, locale, provedor/braço/posição/tentativa da resposta, query duplicada e
campo material. Também adulterar bytes sem recalcular selos. Exigir erro de
integridade/reprodutibilidade, saída zero e nenhuma nova busca.

Retomada completa e parcial deve aceitar schema1 e schema2 íntegros; schema1 pode
omitir autorização/orçamento, mas metadado presente divergente é rejeitado.
Reutilizar SUCCESS/EMPTY; consultar novamente UNKNOWN e queries ausentes,
preservando a sequência natural intercalada e os IDs 0020 por posição. Conservar
snapshots/selos no controle autenticado completo. Recusar nomes de braço
repetidos/desconhecidos e IDs repetidos/vazios ou texto vazio antes de buscar.
Ordenar cópias: nenhum pedido válido ou inválido pode modificar as queries do
chamador. Os testes não mudam rubrica, política, formatos ou hashes históricos.
