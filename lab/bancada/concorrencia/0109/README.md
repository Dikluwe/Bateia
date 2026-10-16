# 0109 — Observador exploratório de concorrência Go

A implementação ativa evoluiu para v8 no [passo 0116](../0116/README.md).
Corpus e results deste diretório registram a v1; os limites e resultados abaixo
se referem ao ensaio histórico. Para executar a avaliação atual, usar os
corpora v2 do 0110.

Piloto ativo no Lab, com transformação pura em observer.go, avaliação em
evaluation.go e efeitos locais na CLI cmd/observer. Nenhuma camada produtiva
importa este pacote. Não há promoção de mecanismo ou alteração do catálogo.

## Uso

Na raiz do Bateia:

```sh
go build -o /tmp/bateia-observer ./lab/bancada/concorrencia/0109/cmd/observer
/tmp/bateia-observer observe --input lab/bancada/concorrencia/0109/testdata/buffered-input.json
/tmp/bateia-observer evaluate --corpus lab/bancada/concorrencia/0110/testdata/baseline-v2.json
/tmp/bateia-observer observe --input lab/bancada/concorrencia/0109/testdata/conc-input.json
```

Omitir --output envia JSON ao stdout. --output aceita somente arquivo novo e
publica por link exclusivo de um temporário completo, com permissão 0600; não
sobrescreve arquivos ou aliases existentes. Evaluation com expectativa divergente
publica o relatório e retorna exit code 1. Requer filesystem com hard links
(no ensaio, Linux). Fontes não são compiladas ou executadas.

## Entrada e observações

Request schema_version 1 declara origin, build (go_version, goos, goarch, tags)
e files (path, content, sha256 hexadecimal sem prefixo). Origem/build não são
autenticados. Até 32 arquivos .go relativos únicos, 256KiB por arquivo e 2MiB
para fontes agregadas e pedido serializado; CLI limita entrada a 2MiB. Conteúdo
deve ser UTF-8 e conferir com SHA256. JSON duplicado/desconhecido/múltiplo é erro.

Observações SYNTAX_ONLY têm property/value, resource textual, função com
posição física/receptor, condições sintáticas e localização física com hash
do arquivo. Diretivas //line não mudam a proveniência. Closures têm escopo
próprio, sem atribuir condições de execução do pai à sua análise interna.

As propriedades localizadas são channel_buffer_syntax, channel_send,
channel_receive, select_shape, goroutine_launch e cancellation_candidate.
A sintaxe de make com chan T identifica forma sem buffer, capacidade literal
ou expressão. Não valida tipos ou validade de capacidade para GOARCH;
make sombreado é lacuna. .Done() recebido é candidato textual estruturado,
sem identidade context.Context verificada.

Lacunas incluem build/tipos não verificados, parsing falho, limite de workers,
trabalho admitido e conclusão de cancelamento não resolvidos. Parsing falho
em um arquivo conserva observações dos outros com status PARTIAL_PARSE.
Essas lacunas não significam ausência do mecanismo ou defeito demonstrado.
Condições if/else/loop/case são localizadores sintáticos, não condições de
caminhos cuja viabilidade foi provada. O corpus compara perfis completos de
property/value e lacunas exigidas; não avalia a semântica das condições.

## Evidência e limites

O corpus autoral tem 15 casos em 12 famílias: controles mínimos, variantes
e parsing malformado. Perfil e lacunas esperados coincidiram nos 15 casos;
expectativa errada e observação extra produziram divergência. A métrica
syntax_profile_agreement=1 pertence a esse corpus exploratório, NÃO é
precisão semântica independente. independent_semantic_accuracy_measured=false.
Expectativas estão congeladas em testdata/corpus.json e vinculadas no recibo.

conc-input.json conserva dois arquivos do SHA
5f936abd7ae87036af1f75c95fb9d0daaf00116b de sourcegraph/conc, coletados no 0106.
O piloto produziu sete observações sintáticas; nenhuma propriedade global
de limite/cancelamento foi comprovada. make(limiter,n) não é resolvido porque
limiter é um tipo nomeado; p.handle.Go não é tratado como go statement.
O default observado não demonstra rejeição: nessa fonte, ele participa do
caminho que pode iniciar outro worker. É necessária análise de tipos/fluxo
para estudar essas relações, não inferência por nomes.

Limites adicionais: nenhuma seleção de build, análise de alcance, aliases,
chamadas interprocedurais, identidade de canal, compilabilidade, repositório
completo ou efeito executado. Funções sem construções reconhecidas não
comprovam ausência de concorrência. Identidades/posições mudam com edição de
bytes; perfis de propriedade podem conservar equivalência sob renomeação.

results contém exemplos, replay/controles e logs da suíte/race. go test ./...,
go vet ./..., build do produto/observador, lint e diff check passaram.
A suíte usou assinatura Git desativada somente no processo e loopback para
fixtures. Gates selados de módulo separado não foram reproduzidos.

Próxima avaliação: referência revisada independentemente e um recorte de
resolução de tipos/fluxo com condições, antes de integração no produto.
