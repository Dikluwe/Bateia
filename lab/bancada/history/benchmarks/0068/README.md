# Benchmark 0068

Benchmark offline para o compilador finito de uma capacidade geométrica. Separa dois
corpos de evidência:

- `real_replay`: hashes e observações dos dossiês produzidos antes do candidato; serve
  para reproduzir a causa, não para alegar ganho novo;
- `counterfactual_cases`: respostas controladas para discriminar tradução, queries
  técnicas, admissão same-file, união, custo e contaminação sem rede.

O harness futuro deve materializar fixtures a partir de `baseline.json`, executar literal
e candidata com os mesmos budgets e calcular recall@5 e contaminação somente sobre
julgamentos explícitos. Ausência de denominador é `Unknown`. O campo
`evidence_addressable_capabilities` significa apenas que a busca retornou trecho apto a
ser coletado/avaliado; não é comprovação da capacidade.

Não atualizar hashes, alvos ou julgamentos após ler resultados candidatos. Rede pública,
tempo de resposta e generalização multilíngue estão fora do gate.

Estado: implementação verificada em modo compatível. A frase original permanece no
plano, três expansões técnicas inglesas são emitidas e a admissão exige os dois
discriminadores no mesmo arquivo. No contrafactual congelado, recall@5 passou de 0 para
1 e contaminação@5 de 1 para 0,5; 16/16 mutações foram eliminadas. `OBSERVED` vale apenas
para esse replay offline, não para o índice Sourcegraph atual nem como prova da capacidade.
