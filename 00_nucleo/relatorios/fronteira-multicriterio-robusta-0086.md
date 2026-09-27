# Fronteira multicritério robusta — passo 0086

## Resultado

O Bateia agora separa três decisões: admissibilidade por gates, não dominância por
Pareto e preferência sob cenários explicitamente declarados. Nenhum candidato vence
por ordem alfabética na superfície avançada e nenhum peso global é inferido.

## Método

- `frontier`: candidatos sem dominador certo;
- `robust_frontier`: candidatos que nenhum rival pode dominar dentro dos intervalos;
- `INCOMPARABLE`: cada candidato é melhor em critérios diferentes;
- `EVIDENCE_UNCERTAIN`: os intervalos impedem concluir a ordem;
- vencedor necessário: o menor score possível não perde para o maior score rival;
- vencedor possível: o maior score possível não perde para o menor score rival.

Scores são normalizados pelas escalas declaradas e existem somente dentro de cada
cenário nomeado. As contribuições por critério permanecem no resultado.

## Superfície

```text
bateia idea decide frontier \
  --input multicriteria-request.json \
  --output multicriteria-report.json
```

Medições `PROVEN` são pontuais e exigem fonte. Medições `PARTIAL` usam intervalo e
exigem fonte. Estado ausente, conflitante ou inconclusivo ocupa toda a escala do
critério, evitando transformar desconhecimento em desempenho ruim.

## Limites

- esta etapa avalia cenários finitos declarados; ainda não resolve regiões contínuas
  de pesos por programação linear;
- ELECTRE/PROMETHEE não foram introduzidos: os gates existentes cumprem o papel de
  veto e a relação de Pareto preserva incomparabilidade;
- critérios de custo, manutenção e adaptação precisam ser medidos ou estimados com
  proveniência pelo chamador.

## Verificação

- crystalline-lint --fail-on warning .: sem violações;
- go test ./...: aprovado;
- go vet ./...: aprovado;
- go build -o /dev/shm/bateia ./04_wiring: aprovado;
- git diff --check: aprovado.

A bancada cobre Pareto, incomparabilidade, gates, incerteza, cenários e ordem.
