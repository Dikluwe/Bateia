# insight_reference_resolver

**Contrato:** passo0142.

ResolveRepository para trilhas: validar identidade canônica/pública, metadata limitada e commit da branch. Nenhum raw fallback/redirect/host arbitrário; erros só códigos.

O FollowReader GitHub é um adaptador dedicado sobre Client: árvore raiz não recursiva limitada e leitura direta Contents sem o fallback raw do extrator legado. Até4 GETs (metadata/commit/árvore/arquivo), corpo256KiB/64KiB, arquivo32KiB com tamanho/base64/blob verificados. O leitor primário legado permanece intacto.
