# Herança de G7 — reprodução e rastreamento do pipeline

**Resultado:** a associação da nota anterior a um pin novo foi reproduzida no catálogo gerado pela release oficial rc.101 e no site gerado pelo código original. **Não demonstramos esse problema no catálogo público.** Há um bloqueio independente de formato no caminho de publicação examinado.

## Reprodução controlada

Usamos um corpus sintético `mini-g7`, com pins A (`aaaa…`) e B (`bbbb…`), pacotes assinados com uma chave descartável e uma avaliação G7 criada somente antes da troca. Os snapshots são opacos, suficientes para o leitor de metadados do catálogo; não são índices restaurados nem implementações cuja qualidade foi realmente medida. A nota 5/5 é um dado de fixture, não um julgamento de relevância feito nesta sessão.

O harness executa os programas originais: reescrita do manifesto de re-pin, validação de registro, builder de catálogo do binário oficial, verificação de assinatura, gate pré-publicação e gerador completo de site. Nenhuma etapa de deploy foi executada. O fallback offline do gerador está explícito (`XERJ_HUB_OFFLINE=1`); não se presume que a fixture possui um release público.

| Estado | Pin no catálogo | Nota no catálogo rc.101 | Página gerada | Avaliação |
| --- | --- | --- | --- | --- |
| Antes | A | G7 5.0/5 | 5.0/5 | Criada para o cenário A |
| Após re-pin | B | G7 5.0/5 | Pin B e 5.0/5 | Mesmos bytes, sem nova avaliação |
| Sem o arquivo G7 | B | Ausente | Sem a nota 5.0/5 | Removida como controle |

A validação do registro passou nos três estados. O catálogo e suas assinaturas foram construídos e verificados nos três estados. O hash da avaliação é idêntico antes/depois. Um pin propositalmente divergente foi recusado pelo gate e um pacote adulterado foi recusado pelo builder por assinatura inválida. Esses controles estão em [results/result.json](results/result.json) e nos respectivos logs.

Isso confirma **herança sem vínculo causal exigido no leitor**, e não demonstra que o número herdado seria necessariamente incorreto para uma implementação real. Uma alteração equivalente pode conservar a mesma qualidade. O risco é apresentar o número como avaliação daquela versão sem um recibo que sustente essa associação.

## O pipeline completo tem uma distinção adicional

| Etapa examinada | Comportamento relevante |
| --- | --- |
| `hub_packs.py rewrite_manifest` | Atualiza SHA e pin_date; o checklist de PR conserva revisão de licença antiga e pede nova revisão. Não remove nem vincula avaliações G7. |
| Planejamento e criação do pack | Seleção por pin/release; assinatura e metadados do pack não são uma avaliação de recuperação. O ensaio fornece packs sintéticos, sem executar indexação/restauração. |
| Builder na branch corpus-hub | O fonte nesse commit deixa `quality_grade: None`; não é o mesmo comportamento de main/release. Essa distinção foi inspecionada estaticamente. |
| Builder da release oficial rc.101 | Executado: emite formato 2 e associa a nota por slug, preservando a nota após troca de pin. |
| `check-catalogue` das duas branches | Executados: exigem formato 1 e recusam a saída original da release por `format_version is 2, expected 1`. Esse bloqueio foi mantido, sem corrigir o programa original. |
| Controle isolado de formato 1 | Uma cópia sintética do envelope, explicitamente convertida para formato 1, passa no gate com a nota herdada. É um controle do gate, não uma saída original do builder nem prova de deploy. |
| `gen.py` e `validate_hub.py` de corpus-hub | Executados sem alteração: geram página com pin B e nota antiga e não recusam o registro por essa associação. |
| Catálogo público assinado | Consultado em 2026-11-08: formato 1, 514 entradas e nenhuma nota G7. A assinatura foi validada com a chave pública do repositório. |

O [snapshot público](results/live-catalogue.json), sua [assinatura](results/live-catalogue.json.sig) e a [síntese da consulta](live-check.json) delimitam o que está publicado. A leitura pública não prova que houve ou deixou de haver avaliações internas. Tampouco permite associar o incidente sintético a algum corpus real do serviço.

## Versões e limites de execução

- main: `1557de2f5be981d9db441969201c92a8af6d2de3`.
- corpus-hub: `708e974516a4ba1ba9bfa312aff164041997cfd6`.
- Binário: release oficial `v1.0.0-rc.101`, checksum do arquivo conferido contra o publicado; identidade em [release.json](release.json).
- O teste Rust offline do builder não iniciou: o cache não tinha `axum-server`. A reprodução dinâmica usou o binário oficial, não uma reimplementação Python do leitor Rust.
- O Python do harness somente constrói a fixture, assina com chave descartável e coordena os programas originais. Chaves privadas não são preservadas.
- Não executamos Actions, deploy, criação de releases ou restauração de índices. Portanto, o achado não é uma afirmação de publicação incorreta em produção.

## Repetir

Com os checkouts nos commits acima e o binário oficial extraído em `/tmp`:

```sh
python3 lab/bancada/insights/xerj-2026-11-08/g7-reproduction/reproduce.py \
  --hub /tmp/bateia-xerj-hub-20261108 \
  --main /tmp/bateia-xerj-source-20261108 \
  --binary /tmp/xerj-official-rc101/unpacked/xerj-1.0.0-rc.101-x86_64-unknown-linux-gnu/xerj \
  --output /tmp/xerj-g7-reproduction-results
```

Requer Python, Git, OpenSSL, xxd, zstd e Bash. O script recusa os pins e o binário divergentes. A chave descartável muda em cada execução, portanto hashes das assinaturas/pacotes podem mudar; o esperado é conservar as transições e os resultados dos controles. Código executado e seus hashes estão no resultado.

## Como apresentar ao Ivan

“Reproduzimos localmente que uma avaliação G7 anterior continua associada ao corpus após troca de pin, no gerador do site e no builder rc.101. Não identificamos uma ocorrência no catálogo público. Existe alguma etapa operacional que exige reavaliar ou declarar a validade da nota para o novo pin? Poderíamos vincular o recibo ao pin, pack e protocolo de consulta?”

Essa formulação comunica o observável e permite que Ivan explique uma garantia operacional fora do repositório. O assunto da incompatibilidade de formato 2/1 deve ser tratado separadamente.
