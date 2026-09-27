# Prompt: limites e fontes regulares do provedor GitHub

Camada L3. Gerar `03_infra/github/security.go`.

Impor limites antes de JSON/base64: conteúdo remoto 2 MiB por arquivo, envelope Contents 4 MiB, árvores/descobertas 8 MiB, metadados 1 MiB. Verificar Content-Length e ler no máximo teto+1 mesmo em chunked/compressão. JSON inválido ou limite excedido é terminal: não degradar para outro caminho de leitura mais permissivo.

Validar owner/repo, revisões (HEAD permitido no legado, nunca tratado como imutável), caminhos relativos sem travessia/absolutos/controles. Fazer escaping explícito. Bloquear redirecionamentos; nunca enviar token para raw. A base customizada existe apenas como configuração explícita do adaptador/testes, não como URL fornecida por candidatos.

Contents pode desreferenciar symlinks internos. Antes de extração exigir árvore Git completa com modos explícitos 100644/100755, caminhos únicos e regulares; recusar symlink, submodule, diretório, árvore truncada, arquivo ausente e tamanho acima do teto. Exigir caminho/tipo/encoding seguros no envelope Contents e medir bytes realmente decodificados.

No máximo 100 arquivos e 8 MiB admitidos por extração; pré-recusar soma declarada da árvore acima do orçamento. Tamanho falso não permite ultrapassar o teto final. Todos os erros deixam extração sem evidência admissível, não retorno parcial disfarçado de sucesso.

Não clonar, executar checkout, hooks, filtros, nem herdar ambiente/configuração Git. Não há sandbox de subprocessos Git prometido: eles foram removidos deste provedor. Quando API da árvore/Contents/raw falharem, retornar erro para coleta inconclusiva. Esta escolha troca disponibilidade por limite previsível de recursos; não permite extração offline por clone e adiciona releitura da árvore antes da extração.
