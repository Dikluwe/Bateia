# Prompt espelhado — gate discriminatório pré-candidato do Passo 0065

## Linhagem e autoridade

Espelha `04_wiring/step0065_gate/discriminatory_test.go`. Este artefato pertence ao
papel ADVERSÁRIO e nasce antes do candidato. Suas únicas entradas normativas são o
Passo 0065, `benchmark/0065/{contract.md,baseline.json,budget.json,manifest.json}` e a
skill `tekt-materializacao-segregada`. É proibido ler ou importar a futura
implementação `evidence_coalition.go`, editar contrato/oráculos ou corrigir a solução.

## Obrigação executável

Ler `00_nucleo/benchmark/0065/adversarial-cases.json` e validar estritamente:

- corpus contratual positivo, negativo e opaco, incluindo os cinco valores baseline;
- M01–M12 uma única vez, cada qual com referência pública preservada, mutante violado e
  testemunha ligada por SHA-256 ao ID, regra, estímulo e diferenças observáveis;
- diferenças reais exatamente nos caminhos declarados, sem aceitar texto narrativo como
  prova de rejeição;
- `mutation_score = rejected_valid_mutations / valid_mutations = 12/12 = 1.0`;
- opacos somente como `Unknown`, fora do denominador de mutações;
- vinte hashes canônicos idênticos para repetição/reordenação e igualdade `nil`/vazio;
- operações iguais a `normalized_witnesses + required_role_pairs`, e M12 rejeitada por
  violar a contagem ou declarar enumeração de subconjuntos;
- hashes SHA-256 congelados das entradas L0 antes de aceitar a suíte.

O gate é um teste do poder discriminatório do contrato/modelo público. Ele não chama,
aprova ou mede a implementação candidata e não emite selo nem certificado final.

