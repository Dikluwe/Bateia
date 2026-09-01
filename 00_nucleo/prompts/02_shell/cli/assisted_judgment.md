# Assisted judgment — contrato da CLI

Exponha na CLI a materialização e a inspeção offline do oráculo assistido do Passo 0021.
A CLI recebe caminhos explícitos para snapshot selado, julgamentos da IA e revisões
humanas, além de destino de saída. Ela não incorpora nem invoca LLM, não usa rede e não
aceita texto de candidato como comando ou configuração.

## Comportamento observável

- validar integralmente entradas antes de publicar saída;
- produzir pacote cego sem braço, query, rank/posição original, métricas parciais,
  registry protegido ou expectativa de promoção;
- listar para revisão os motivos determinísticos de escalonamento, incluindo confiança
  `< 0.80`, `UNCERTAIN`/`CONFLICT`, fonte inválida, ambiguidade e amostra de 10%;
- mostrar ao revisor apenas dados necessários ao caso, sem métricas agregadas antes da
  decisão;
- aceitar `CONFIRMED`, `OVERRIDDEN` e `RESOLVED`; exigir motivo em substituição e validar
  o veredito humano;
- preservar no arquivo final julgamento original, revisão, efetivo, versão da rubrica e
  histórico; nunca atualizar destrutivamente uma decisão anterior;
- usar escrita atômica: erro de validação não deixa oráculo parcial nem corrompe artefato
  existente;
- oferecer saída legível por máquina, determinística e sem segredos. Erros vão para
  stderr e retorno não zero; dados estruturados vão para stdout/arquivo sem mistura.

Casos escalonados sem revisão são emitidos como `UNKNOWN` e resultam em
`INCONCLUSIVE_HUMAN_REVIEW`. Violações de cegamento resultam em `FAIL_BLINDING`; fontes
inválidas em `FAIL_EVIDENCE_GROUNDING`; baixa concordância ou override excessivo em
`FAIL_AI_RELIABILITY`.

Para corpus vazio, a CLI deve criar um oráculo vazio válido, sem solicitar revisão, sem
divisão por zero e com métricas indisponíveis representadas explicitamente, nunca como
zero enganoso, `NaN` ou infinito. Repetir o comando com as mesmas entradas produz bytes
canônicos equivalentes, salvo metadados explicitamente excluídos do hash.
