# 0111 — Ambiente de tipos fornecido e rastreável

O código ativo evoluiu para [0112](../0112/README.md). Os recibos deste
diretório conservam o ensaio v3; os dados congelados continuam consumíveis.

Piloto v3 em lab/bancada/concorrencia/0109. A CLI aceita --environment; o
observador importa apenas fontes e export data presentes na memória. Lookup
de go/importer é explícito, sem fallback para disco. Nenhum pacote candidato
é compilado ou executado. Produto não importa Lab.

## Fronteiras

prepare.py realiza efeitos de bancada: reutiliza o bundle congelado 0106,
confere os hashes e o commit, seleciona nove fontes não-test de pool, panics
e raiz, e obtém dados de tipos da biblioteca padrão com go list -deps -export.
Apenas pacotes Standard são aceitos. Conserva uma seção __.PKGDEF por archive,
sem código objeto; o ambiente contém 64 exports, 1.539.504 bytes de dados.
A preparação recusa substituir congelamentos existentes. Reproduzir em um
novo diretório/cópia de bancada; a seleção não certifica um build do repositório.

O ambiente schema 1 declara toolchain/goos/goarch, packages
(import_path/directory) e exports (import_path/origin/data base64/sha256).
Aceita apenas linux/amd64 nesta fatia, até 32 bindings, 128 exports, 32MiB de
dados e 64MiB JSON. Integridade, caminhos, duplicatas e colisões são validados.
Declarações e hashes não autenticam a origem externa; o recibo local vincula os
bytes ao bundle anterior e ao runtime usado. Contexto de build continua limitado.

ObserveWithEnvironment é uma transformação pura. O importer de fontes resolve
pacotes recursivamente; ciclos, importação ausente, parsing e erros impedem
atribuição por tipos aos dependentes. A checagem também usa dados gc fornecidos
em memória. Não busca arquivos/pacotes no ambiente instalado. A CLI faz somente
leitura/publicação local; a preparação é uma operação separada.

O relatório liga environment_sha256 canônico aos bytes fornecidos. A nova
autoridade é GO_TYPES_SUPPLIED_ENVIRONMENT, restrita à identidade de tipo.
Observações de capacidade permanecem SYNTAX_ONLY. type_diagnostics localiza
falhas da entrada por posição física/hash; elas não são veredito sobre defeitos
do projeto. BUILD_SELECTION_NOT_VERIFIED e limites de workers/admissão/
cancelamento permanecem lacunas.

## Execução offline

Na raiz:

```sh
go build -o /tmp/bateia-observer-0111 ./lab/bancada/concorrencia/0109/cmd/observer
/tmp/bateia-observer-0111 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json --environment lab/bancada/concorrencia/0111/testdata/environment.json
/tmp/bateia-observer-0111 observe --input lab/bancada/concorrencia/0111/testdata/conc-input.json
/tmp/bateia-observer-0111 evaluate --corpus lab/bancada/concorrencia/0110/testdata/types-corpus.json
```

A preparação usa toolchain/cache local e somente biblioteca padrão; a análise
usa os JSON congelados. Não requer rede para replay. --output continua exclusivo
para arquivo novo. evaluate não aceita ambiente nesta fatia.

## Resultado

Sem ambiente: oito observações sintáticas, nove lacunas de checagem de pacote
e zero atribuições de tipos. Com ambiente: dez observações, zero diagnósticos
de tipos, uma atribuição de tipo — limiter em pool/pool.go:96. Sua capacidade
é CAPACITY_EXPRESSION, expressão n. Remover bindings fonte suprime a atribuição;
export com hash adulterado falha antes da publicação. Replays byte idênticos.

Os 15 casos da baseline e os 16 contrastes de tipos passaram. São referências
autorais, sem precisão semântica independente. Testes adicionais cobrem ciclos,
dependências inválidas/ausentes, exports corrompidos, ordem, integridade e CLI.
Resultados/recibos estão em results e a preparação em testdata.

Nenhum limite de workers foi demonstrado. Ainda falta relacionar a reserva no
select, a chamada tipada p.handle.Go, o lançamento go interno e a liberação
adiada do token. Esse será um recorte de fluxo com condições explícitas,
antes de atribuir uma propriedade global de bounded concurrency.
