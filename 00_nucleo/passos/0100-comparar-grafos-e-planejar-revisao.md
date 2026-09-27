# 0100 — Comparar grafos e planejar revisão incremental

## Resultado

Comparar dois snapshots congelados, listar diferenças registradas e produzir
um plano de revisão ancorado no grafo anterior. Não substituir a decisão
histórica, aprovar a nova declaração ou consultar/executar software candidato.
A pesquisa simples e os artefatos padrão de recommend permanecem inalterados.

```sh
./bateia idea trace compare --input comparacao.json --output revisao.json
```

Entrada JSON: schema_version "1", before e after com grafos completos
decision-trace-v1 e graph_hash não vazio. Sem --output, emitir JSON em stdout.
O diretório de saída deve existir. Não é uma comparação de arquivos de fontes;
os hashes comprovam coerência das representações tipadas, não autenticidade.

## Deltas e impacto

A saída decision-trace-comparison-v1 contém before_graph_hash, after_graph_hash,
scope RECORDED_DEPENDENCIES_ONLY, historical_decision e candidate_decision.
candidate_decision é somente o veredito declarado no segundo snapshot;
decision_impact permanece NONE. impact.decision_before e decision_after continuam
iguais à decisão do primeiro grafo. Mudanças exigem revisão, não promoção automática.

node_changes registra ADDED, REMOVED ou DECLARATION_CHANGED, com os nós completos
antes/depois disponíveis. Comparar também tipo, rótulo e fonte: fingerprints
iguais não tornam metadados iguais. added_edges e removed_edges incluem todas as
relações, inclusive REFUTES. Ordem de nós, arestas e limitações não cria diferenças.
metadata_changed identifica mudanças nas limitações declaradas.

Nós alterados/removidos iniciam revisão no grafo anterior. Diferenças de arestas
iniciam revisão pelo consumidor anterior, ou pela decisão histórica quando o
consumidor é novo. Nós adicionados e mudanças de metadados iniciam revisão
conservadora da decisão histórica. Novos nós aparecem nos deltas, mas não são
falsamente apresentados como ancestrais existentes no grafo antigo.

Uma fonte ausente no segundo snapshot não é declarada REVOKED nem EXPIRED:
ausência pode significar mudança na modelagem ou seleção. Os eventos reais
continuam sendo entradas explícitas de idea trace impact. Também não se infere
completude das dependências. O construtor mantém a dependência global do relatório,
portanto diferenças podem exigir uma revisão ampla; não a ocultar para aparentar
precisão incremental.

## Segurança e limites

Validar hashes, DAG e esquema dos dois grafos, sem alterar as entradas. Reutilizar
JSON estrito, proteção de entradas/aliases/hardlinks e escrita atômica individual
da CLI trace. Não carregar .env ou provedores no despacho offline.

Limites agregados de entrada/saída: 2 MiB de texto/JSON, 100000 valores,
profundidade 64; cada grafo até 1000 nós e 4000 arestas. No máximo 100 raízes
distintas para o plano de impacto. Excesso falha explicitamente, sem truncar
deltas ou anunciar plano parcial. Guardar a saída antes da serialização integral
e não sobrescrever uma saída existente quando a comparação é inválida.

## Verificação

Testar igualdade com reordenação, rótulo sem mudança de fingerprint, remoção e
adição de fonte, arestas REFUTES, alteração declarada de decisão, hashes inválidos,
limites agregados/expansão, não mutação e CLI offline com entradas protegidas.
O relatório é uma lista operacional de revisão, não uma nova evidência factual.
