# Ambiente de tipos em memória — 0111

Regime Lab. Environment schema 1 contém bindings import_path/directory para
pacotes de fontes do Request e export data Go gc base64 ([]byte), cada um com
SHA256, import_path e origem declarada; toolchain/goos/goarch declarados.
Até 32 bindings, 128 exports, 32MiB agregados de export data, 64MiB JSON;
recusar duplicatas, caminhos inseguros e colisão fonte/export. Somente linux
amd64 nesta fatia. Hash canônico ordenado do ambiente no relatório, sem mutação.
Resolver recursivamente apenas fontes/exports fornecidos; lookup gc em bytes
por io.NopCloser, nunca fallback para disco/rede/go tool. Importações ausentes,
ciclos, parsing ou erros de pacote suprimem observações por tipos. Não autorizar
compilabilidade/build real, efeitos ou limites globais. Sem ambiente, contrato
0110 continua. Com ambiente, autoridade GO_TYPES_SUPPLIED_ENVIRONMENT.
