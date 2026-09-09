# Passo 0050 — Recuperar implementações de candidatos remotos

2026-09-08. Obrigação derivada do diagnóstico 0049: a busca pública precisa alcançar definições Go no commit do candidato, e não parar na seleção por nomes de arquivos.

## Contrato observável congelado

O adaptador GitHub implementará a porta existente ImplementationSearcher. O Router delegará ao provedor selecionado, local ou remoto, quando este implementar a porta. Provedor sem suporte permanece not_evaluated.

A recuperação GitHub exige owner/repo válidos e commit SHA hexadecimal completo de 40 caracteres; identidade inválida produz erro/partial antes de I/O. Ler árvore recursiva nesse commit e blobs por SHA da árvore. Validar objetos Git blob (SHA-1 de `blob <tamanho>\0<bytes>`), tamanho declarado e conteúdo completo antes de admitir arquivo. Árvores truncadas, respostas inválidas, arquivos opacos e falhas permanecem partial, nunca ausência comprovada. Não fazer clone, executar código ou relaxar referência para HEAD.

Usar caminhos Go aceitos pelo filtro existente, apenas blobs regulares (100644/100755), ordenados e sem duplicatas ambíguas. Rejeitar caminhos inseguros. Aplicar defaults/tetos remotos de 64 tentativas de arquivo, 4 MiB de fonte total, 512 KiB por arquivo e 24 hits; orçamento solicitado menor prevalece. Uma requisição de árvore e no máximo MaxFiles requisições de blob. HTTP 403/429 interrompe coleta; contexto cancelado interrompe. Limitar corpo de árvore a 8 MiB e corpos JSON de blobs de acordo com orçamento de conteúdo base64 mais envelope limitado. Contabilizar conteúdo decodificado recebido, inclusive descartado, e não indexar bytes parciais ou identidade divergente. Marcar Limited quando limite exclui candidatos.

Arquivos válidos alimentam RankGoImplementations sem mudança de ranking, glossário ou interpretação. Examined registra somente arquivos íntegros admitidos, hash SHA-256 e bytes; BytesRead registra bytes de fonte recebidos. Releitura posterior e validação de hits existentes permanecem obrigatórias.

## Aceitação

Teste independente via servidor GitHub simulado: mesmo conjunto Go local/remoto produz mesmos hits, pesos, linhas e hashes quando todos os arquivos cabem. Consulta funcional recupera implementação remota cujo caminho não repete a pergunta. Verificar encaminhamento pelo Router e seleção/revalidação na jornada usecase. Casos de commit inválido, blob divergente, árvore truncada, symlink, limites, 403/429 e cancelamento não viram complete com falsa ausência. Preservar controles locais 0048 e testes existentes.

Build, lint e suíte completa. Ensaio público real se rede disponível, reportando falhas como limitações. Ganho exigido neste passo: passar de nenhuma implementação remota examinada a fontes remotas verificáveis no mesmo snapshot. Não alegar ganho geral de pertinência nem capacidade comprovada.

## Regime e escopo

Tekt A/B sem atestação de isolamento físico. Contrato do integrador, implementação A e testes B separados; filesystem compartilhado. Congelar contrato e hashes de baseline antes da implementação. Índice documental, interpretação de ação/objeto/condições e cobertura do catálogo ficam para obrigações posteriores: não são resolvidos por habilitar a coleta remota.
