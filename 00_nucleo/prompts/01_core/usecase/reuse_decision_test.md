# Reuse decision adversarial tests

Ataque o contrato do decisor sem compartilhar sua implementação: cobertura total
com conflito não pode virar FORK; falha ou inconclusão não pode virar BUILD;
exercite cada decisão; permute candidatos e scores; altere apenas a política;
compare a entrada antes/depois; exija regra e evidência nas decisões positivas.

## Migração de fixtures 0043 por testador separado
Preservar decisões positivas/negativas e assertions existentes. Os positivos agora carregam fontes com identidade, hash SHA256 dos snippets literais, linha, caminho, regra e kind compatível. Symbol demonstra implementação estática delimitada; teste independente auxilia capacidade/operação. Gates licença e manutenção recebem license e repository_metadata. Não basta Status=PROVEN sem fonte. Esta migração altera validade das entradas, não a política esperada dos cenários históricos.

Os componentes dos controles positivos vinculam as mesmas fontes das capacidades que alegam cobrir. Uma fonte genérica "component" sem correspondente avaliação não demonstra parte reutilizável; corrigir essa ligação preservando EXTEND/COMPOSE/INTEGRATE esperados.
