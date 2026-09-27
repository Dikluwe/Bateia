# Passo 0086 — Fronteira multicritério robusta

## Estado e propósito

- Estado: implementado e verificado em 2026-09-24.
- Depende dos passos 0010, 0083, 0084 e 0085.
- Objetivo: comparar candidatos elegíveis sem esconder trade-offs em um score único.

## Entregas

1. reutilizar os gates eliminatórios do decisor de reúso;
2. aceitar critérios direcionais com escalas explícitas;
3. representar medições pontuais, intervalares e ausentes;
4. calcular dominância certa, possível e fronteiras de Pareto;
5. explicar incomparabilidade e incerteza par a par;
6. avaliar apenas cenários de preferência declarados;
7. mostrar vencedor necessário, possível e sensibilidade entre cenários.

## Invariantes

- licença, operação, arquivamento e restrições não são compensados por pesos;
- medição ausente ocupa toda a escala declarada, nunca recebe zero;
- um score só existe dentro de um cenário de preferência nomeado;
- a fronteira de Pareto não depende de score de screening, estrelas ou ordem de entrada;
- a saída pode abster e formular perguntas de preferência ou de evidência;
- `idea decide` simples mantém seu contrato anterior.

## Superfície

`idea decide frontier --input request.json --output report.json` recebe JSON estrito e
executa offline.

## Gate

Testar Pareto exato, trade-off, gate não compensável, intervalo ausente, cenários
antagônicos, ordem invariável, CLI estrita, lint, testes, vet e build.
