# 0110 — Canais nomeados sob checagem restrita

A implementação ativa evoluiu para [0111](../0111/README.md). Resultados e
recibos deste diretório conservam o ensaio histórico v2, sem ambiente externo.

Evolução ativa do piloto 0109, mantida no mesmo pacote. Versão
`go-concurrency-syntax-0110-v2`. Código e efeitos continuam no Lab; produto
não importa o piloto. Nenhuma dependência externa foi adicionada.

`go/types` recebe somente ASTs fornecidas, agrupadas por diretório/pacote,
com tamanhos gc/amd64 fixos e sem importer. Não carrega dependências, escolhe
build ou executa fontes. Qualquer erro no grupo, incluindo parsing parcial,
suprime atribuições por tipos; os fatos sintáticos continuam disponíveis.

Um make builtin com expressão de tipo canal e resultado canal em grupo sem
erros emite `channel_type_identity=RESOLVED_CHANNEL_TYPE`, autoridade
`GO_TYPES_SUPPLIED_PACKAGE`. Sua capacidade continua uma observação
`SYNTAX_ONLY`; identidade de tipo não prova identidade de instância, limite de
workers/trabalho ou cancelamento. Canais nomeados, aliases e declarações locais
ou em outro arquivo são reconhecidos. Slice/map checados não recebem crédito.

Lacunas novas: `SUPPLIED_PACKAGE_TYPECHECK_FAILED` e
`NAMED_CHANNEL_TYPE_UNRESOLVED`. Na v2, a lacuna histórica
`TYPECHECK_NOT_PERFORMED` foi substituída por
`SYNTAX_PROPERTIES_NOT_TYPECHECKED`, porque a trilha sintática continua sem
exigir sucesso de tipos. A versão no relatório diferencia os contratos.

## Reprodução

Na raiz:

```sh
go build -o /tmp/bateia-observer-0110 ./lab/bancada/concorrencia/0109/cmd/observer
/tmp/bateia-observer-0110 evaluate --corpus lab/bancada/concorrencia/0110/testdata/baseline-v2.json
/tmp/bateia-observer-0110 evaluate --corpus lab/bancada/concorrencia/0110/testdata/types-corpus.json
/tmp/bateia-observer-0110 observe --input lab/bancada/concorrencia/0109/testdata/conc-input.json
```

O corpus original e results do 0109 são evidência histórica v1. A cópia
baseline-v2 conserva seus 15 perfis e atualiza apenas a lacuna exigida.
O corpus novo contém 16 casos autorais em 15 famílias. Concordância desses
perfis mede conformidade exploratória; não precisão semântica independente.
Testes adicionais cobrem grupos, parsing parcial, reordenação e hashes.

Nos dois arquivos congelados de sourcegraph/conc, as sete observações
sintáticas permanecem. A resolução de limiter é suprimida por erros do pacote
incompleto/importações ausentes. Isso é uma limitação registrada, não evidência
de que limiter não é canal. p.handle.Go continua fora do recorte. A próxima
fatia precisa fornecer dependências em memória e registrar seleção de build,
antes de analisar admissão/fluxo.

Recibos e logs atuais pertencem a results; nenhuma fonte de candidato é
executada. A CLI mantém publicação exclusiva de arquivo novo e JSON estrito.
