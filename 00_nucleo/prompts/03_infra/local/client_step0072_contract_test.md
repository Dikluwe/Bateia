# Contrato independente — enumeração local e header C++, passo 0072

Derivar testes somente do passo 0072 e do baseline público congelado. Não ler código
candidato.

Em diretório temporário confinado, verificar que `FetchTree` não corta silenciosamente a
superfície no antigo limite de 20.000 entradas e que inclui um arquivo posterior ao
limite. Verificar também que a coleta pública aceita `span` extensionless quando dirigido
por consulta, mas não classifica arquivo extensionless irrelevante como código.

Não seguir aliases externos, não relaxar limites de leitura e não depender da ordem do
filesystem.
