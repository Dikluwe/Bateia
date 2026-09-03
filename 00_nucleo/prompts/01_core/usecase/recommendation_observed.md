# Progresso observável da recomendação

Estenda `RecommendationService` com a porta observável usada pelo worker
durável. Preserve a projeção do resultado síncrono e emita somente estágios e
contadores canônicos controlados pelo núcleo, nunca mensagens de providers.
