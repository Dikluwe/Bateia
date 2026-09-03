# Prompt — O20 public end-to-end replay self-tests

Carregue integralmente o corpus público end-to-end O19, derive em memória os
campos de integridade como o replay público e entregue os bytes completos a
`Evaluate`. Exija totais sete `Unknown`, zero `Preserved`, zero `Violated`.
Construa vizinhos conformes adicionando validation plan não vazio e exija
`Preserved`. Verifique precedência de violação, JSON estrito e manifesto, sem
usar semantic_state, identidade de caso, hashes externos, marker, expected,
texto humano ou conteúdo reservado.

Após uma revisão sucessora, trate `manifest-o20.json` apenas como âncora
histórica imutável pelo seu próprio digest; não compare seus hashes antigos com
arquivos de uma suite append-only posterior.
