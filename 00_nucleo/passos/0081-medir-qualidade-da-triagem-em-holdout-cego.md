# Passo 0081 — Medir a qualidade da triagem

## Objetivo

Avaliar a triagem em uma amostra separada do desenvolvimento, com rótulos independentes,
sem usar os rótulos durante a previsão e sem alterar o produto para favorecer o corpus.

## Medidas

- precisão e cobertura de `ADVANCE`;
- perdas causadas por `DROP`;
- taxa de auditoria;
- qualidade do ranking;
- calibração quando houver probabilidades válidas.

## Resultado

Na amostra OpenCascade de 24 casos, o modelo preservou todos os relevantes porque enviou
tudo para auditoria: cobertura de automação `0`, taxa de auditoria `1` e nenhum `DROP`.
O ranking apresentou algum sinal, mas abaixo do critério definido; a calibração foi pior
que a prevalência. O resultado não sustenta promoção do modelo nem alegação de
generalização. Os números completos estão no laudo 0039.
