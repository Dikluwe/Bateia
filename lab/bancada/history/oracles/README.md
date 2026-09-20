# Oráculos históricos supersedidos

Este diretório preserva byte a byte oráculos incompatíveis com decisões
arquiteturais posteriores. Eles continuam vinculados aos selos dos passos de
origem, mas não são compilados contra a topologia atual.

O passo 0079 arquiva os oráculos L4 dos passos 0027 e 0028 porque ambos exigiam
por contagem textual que store e fila fossem construídos em L4. A obrigação de
unicidade permanece ativa no gate AST 0079, agora com ownership em L3; a cobertura
funcional Apple/replay foi transportada para um teste sucessor independente.
