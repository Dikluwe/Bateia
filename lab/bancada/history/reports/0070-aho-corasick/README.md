# História do Aho-Corasick no Bateia

Este dossiê é uma execução real do modo `algorithm history` sobre o Aho-Corasick. Ele separa relações apoiadas pelas fontes de relações que continuam historicamente indeterminadas.

## Execução

```bash
GOCACHE=/tmp/bateia-aho-history-go-cache \
  go run ./04_wiring algorithm history analyze \
  --input relatorios/algorithm-history/aho-corasick/input.json \
  --output relatorios/algorithm-history/aho-corasick/report.json
```

Resultado: `COMPLETE`, com 4 algoritmos/artefatos, 7 proveniências, 5 alegações, 8 evidências e 1 inferência. O relatório classificou 1 relação como `CORROBORATED`, 2 como `SUPPORTED` e 2 como `UNKNOWN`.

## Leitura histórica

1. `SUPPORTED` — O autômato Aho-Corasick de 1975 melhora a varredura independente de cada palavra-chave. O artigo original descreve uma máquina finita, processamento do texto em uma passagem e ganho empírico sobre a abordagem direta.
2. `SUPPORTED` — O prefilter inspecionado do Gitleaks deriva de Aho-Corasick: o código importa uma biblioteca com esse algoritmo, constrói sua trie e usa `Match`.
3. `CORROBORATED` — O componente atual do Bateia se parece mecanicamente com a abordagem direta anterior ao Aho-Corasick. Seu `MatchAll` percorre os padrões e chama `strings.Contains` para cada um; a biblioteca padrão confirma que `Contains` delega a uma busca individual por `Index`.
4. `UNKNOWN` — O prompt do Bateia pede explicitamente Aho-Corasick, mas o artefato entregue não contém trie, estados, transições de falha ou função de saída. O nome e a intenção não bastam para comprovar descendência mecânica.
5. `UNKNOWN` — O prompt declara que a ideia foi encontrada numa varredura que incluiu Gitleaks, mas as fontes preservadas ainda não demonstram a cadeia exata de transferência de Gitleaks para o código do Bateia.

## Hipótese de desenvolvimento

O próximo experimento útil é comparar a implementação atual com um autômato Aho-Corasick verdadeiro usando os mesmos conjuntos de padrões e textos. As medidas mínimas são: tempo de construção, tempo de busca, alocações, memória e ponto de equilíbrio conforme cresce o número de padrões. Até esse benchmark existir, trocar a implementação seria prematuro: a versão simples pode ser melhor para conjuntos pequenos, enquanto o autômato tende a ganhar valor quando muitos padrões compartilham prefixos ou são consultados repetidamente.

## Artefatos

- `input.json`: afirmações, cronologia e proveniência auditáveis.
- `report.json`: vereditos determinísticos emitidos pelo Bateia.

SHA-256 da entrada: `bf3431c219a1510768c47860edead806fb319ed403191f272bee2735f30028cc`.

SHA-256 do relatório: `fab58d1d96d77617e15979d21cb7ba6a9771846b175967960fff5e951151544d`.
