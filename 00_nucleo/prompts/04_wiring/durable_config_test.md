# Verificação da configuração durável local

Garanta que uma chave de cursor hardlinkada seja recusada antes de qualquer
mudança de permissão no inode externo. Verifique também tamanho e modo da chave
persistida criada localmente. Uma raiz ou componente symlinkado deve ser recusado
antes de criar arquivos ou alterar metadados no destino externo.
