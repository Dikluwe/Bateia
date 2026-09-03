# Prompt proprietário — `coverage_test.go`

## Papel causal

Materializar a auditoria estrutural pública da matriz de cobertura O2 do Passo
0028, derivada exclusivamente do contrato R3 e do feedback sanitizado O1. O
teste pertence ao autor segregado de oráculos e não consulta implementação
candidata, ataques, rótulos privados, receipts, seals, certificado ou veredito.

## Obrigações

- decodificar estritamente o feedback público e a matriz pública;
- expandir ranges inclusivos e demonstrar cobertura exata, sem duplicata nem
  lacuna, de `M001` a `M160`;
- demonstrar cobertura exata de `P01` a `P12`;
- exigir capability não vazia, operação do harness público e pelo menos um case
  ID executável para cada grupo;
- exigir que cada positivo esteja ligado a seu próprio case ID declarado;
- vincular a matriz aos hashes públicos autorizados do feedback O1;
- manter os casos `O01`–`O10` fora deste teste e fora do mutation score.

## Restrições

O teste não abre a fixture protegida de casos de conformidade nem qualquer
conteúdo adversarial. Ele prova completude estrutural da cobertura declarada,
não aprovação do gate, mutation score, correção do candidato ou emissão de
veredito.
