# Transporte de snapshots íntegros0168

L2. [Passo0168](../../../passos/0168-preservar-integridade-e-determinismo-do-oraculo.md), ASS-01 do [contrato](../../../contratos/integridade-de-entradas-e-replay-v1.md#ass-01--snapshot-íntegro-antes-da-união-cega). CLI prepare com snapshot adulterado/sem selo deve falhar sem criar saída. Preparar controle selado explicitamente; helper de escrita não repara negativos. Efeitos de fixtures pelo host configurado, sem mudar assertivas anteriores ou registros históricos.
