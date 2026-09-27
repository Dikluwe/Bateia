# Prompt: corpus adversarial do provedor GitHub

Camada L3. Gerar `03_infra/github/security_test.go`.

Testes offline com transport sintético devem cobrir positivos fixados, raw sem token e no mesmo commit, symlink/submodule/modo ausente recusados antes de Contents, caminhos inválidos sem rede, base64 inválido, caminho divergente, envelope/arquivo/total acima dos tetos e redirect nunca seguido.

Criar fixture Git local real com LICENSE symlink apontando para dado sintético externo e clone local. Herdar PATH shim e config Git hostis com tokens fictícios durante chamada ao provedor. Em falha API/raw, garantir erro sem extração e sem execução Git, marcador ou leitura externa. A fixture é criada pelo teste; nenhuma rede/segredo real é usado.

Todos os subprocessos para preparar a fixture têm ambiente explícito mínimo; são infraestrutura de teste, não caminho de produção. Não alegar sandbox do Git: teste demonstra ausência de subprocessos do provedor.
