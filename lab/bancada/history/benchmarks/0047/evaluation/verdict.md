# Avaliação B0047

Candidata satisfaz o critério de promoção congelado. Nenhuma revisão solicitada e nenhum esperado alterado após resultados. Autor B não leu código de produção candidata nem diffs; usou contrato, API pública, testes0045 preexistentes e fontes externas do cache para rotular holdout.

| Corpus / medida | Baseline0046 | Candidata0047 |
|---|---:|---:|
| Real0046 recall rotulado@3 | 0/6 | 2/6 |
| Real0046 recall rotulado@12 | 1/6 | 4/6 |
| Real0046 negativas com indícios | 1/3 | 1/3 |
| Holdout recall rotulado@3 | 2/2 | 2/2 |
| Holdout recall rotulado@12 | 2/2 | 2/2 |
| Holdout negativas com indícios | 1/1 | 1/1 |
| Controle0045 positivos recuperados | 4/4 | 4/4 |
| Controle0045 negativas com indícios | 0/2 | 0/2 |

Quatro testes públicos novos passam: discriminação de termo raro no mesmo campo, normalização de corpo longo, contribuições finitas e determinismo sob duplicação/reordenação, identificador exato e exclusão de comentários/literais. Seis testes0045 permanecem passando, incluindo chamadas diretas e resolução insegura.

A amostra real0046 preserva manifestos e perguntas byte a byte. O runner original foi copiado com apenas ajuste da raiz pelo novo diretório e GOCACHE solicitado. O holdout usa github.com/google/uuid v1.6.0, commit 0f11ee6918f41a04c201eceeadf612a377bc7fbc, com todas as fontes Go não teste do ZIP; origem, hashes e justificativas manuais estão congelados em holdout/. Nenhum código externo foi executado. Os adaptadores executam somente Bateia.

As medidas @3/@12 são caminhos distintos entre até24 hits, não extração executada. Precisão é contra arquivos esperados rotulados, não julgamento exaustivo de irrelevância. O holdout tem somente duas perguntas positivas, ambas já recuperadas no baseline, e a negativa produz indício em ambos os braços: isto preserva não regressão, mas não comprova ausência de falsos positivos. A amostra de conveniência e os ganhos no corpus conhecido não demonstram generalização. Regime A/B sem atestação física.

Reprodução: `python3 evaluation/real0046/run.py final`, `python3 evaluation/holdout/run.py final`, `GOCACHE=/tmp/bateia-0047-cache go test ./01_core/parser -run 'TestCorpus0047|TestFunctional0045' -v`. Os runners recusam sobrescrever observações existentes e verificam hashes de entradas congeladas. O comparador exige também captura de controle0045 com o mesmo rótulo. Suite global, CLI e revalidação são verificados pelo coordenador/outros escopos.
