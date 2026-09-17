# Teste proprietário — avaliação local da cobertura dirigida 0072

Verificar que uma amostra geral incompleta não apaga indícios positivos encontrados
quando a cobertura do requisito está `TARGETED_EXAMINED`. A mesma incompletude deve
continuar impedindo `ABSENT`, e bundles legados sem `RequirementCoverage` ou com
cobertura dirigida incompleta devem permanecer `INCONCLUSIVE/COLLECTION_INCOMPLETE`.

O teste usa somente bundles e regras em memória, sem rede, filesystem ou relógio.
