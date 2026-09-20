# Benchmark 0069 — radar exportável para Archify

Este diretório congela o contrato, os ataques e um exemplar reproduzível da exportação
do radar de ideias para Archify Dataflow v1. O Bateia continua responsável pelo
significado; o Archify recebe somente a projeção visual.

## Uso

```bash
go run ./04_wiring idea radar export \
  --input 00_nucleo/benchmark/0069/example-input.json \
  --output radar.archify.json
```

Sem `--output`, o JSON é escrito em stdout. A saída pode ser entregue ao Archify:

```bash
archify deliver dataflow radar.archify.json radar.archify.html
```

## Artefatos mensuráveis

- `example-input.json`: entrada canônica do radar;
- `example.archify.json`: saída determinística do Bateia;
- `example.archify.html`: render produzido pelo Archify no commit congelado;
- `example.archify.visual-check.json`: receipt do navegador em quatro viewports e dois
  temas;
- `adversarial-fixtures*.json` e `04_wiring/step0069_gate`: 27 mutações contratuais;
- `verification.json` e `seal.json`: veredito e hashes finais.

O exemplar passou os nove checks do Archify em perfil `standard`, com zero erros e cinco
avisos: quatro de distância entre label e rota e um corredor compartilhado. O check de
navegador passou em 1440x900, 1600x1000, 1920x1080 e 2048x1320, sem overflow. A inspeção
humana confirmou legibilidade funcional, mas o diagrama mantém espaço vazio à direita e
rotas convergentes densas; por isso ele é evidência de integração, não alegação de
qualidade `showcase`.

A versão externa está fixada em Archify
`6db72a9aea3d0f67a6a034e41f8a5491476a11c1`. O produto não executa Node, Archify, rede
ou filesystem no apresentador puro. A segregação Tekt foi procedural em filesystem
compartilhado, sem atestação física de isolamento.
