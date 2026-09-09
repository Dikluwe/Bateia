# Recuperação de implementações GitHub

Implementar `ImplementationSearcher` em L3 conforme contrato congelado do passo 0050, usando exclusivamente árvore recursiva no commit SHA completo e blobs identificados por SHA. Validar owner/repo antes de HTTP, rejeitar caminhos inseguros, modos não regulares e identidades ambíguas; usar o filtro Go existente. Deduplicar e ordenar caminhos antes de coleta.

Aplicar defaults e tetos de 64 tentativas de blob, 4 MiB totais, 512 KiB por arquivo e 24 resultados; pedidos positivos menores prevalecem. Limitar árvore a 8 MiB e JSON de blob ao tamanho base64 permitido mais 64 KiB de envelope. Não seguir redirects, clonar ou relaxar commit. Interromper em cancelamento e HTTP 403/429. Respostas incompletas ou inválidas, árvore truncada e arquivos opacos produzem relatório partial; limites que excluem arquivos marcam Limited.

Decodificar base64 contabilizando inclusive bytes descartados; somente admitir arquivo completo com tamanho da árvore, tamanho do blob, SHA anunciado e SHA-1 do objeto Git consistentes. Encaminhar arquivos íntegros a RankGoImplementations sem alterar ranking. Registrar somente fontes admitidas em Examined com SHA-256, preservando BytesRead de toda decodificação. Testes independentes pertencem ao autor B; autor A não consulta testes ou oráculos.
