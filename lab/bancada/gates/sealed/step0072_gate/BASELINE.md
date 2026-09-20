# Identidade congelada do baseline

Commit autoritativo: `03fd5f1`.

| Caminho no commit | SHA-256 do conteúdo |
|---|---|
| `01_core/domain/evidence.go` | `a3387e44d328ef0a1e1239226fa16c3ecac2e136a293bab6af4ffe47dd230e8b` |
| `01_core/domain/interfaces.go` | `8292e9a0b1425aabdad38fbf6a6c4b6964d146116e4b8ba5c1f3b982aee50ad4` |
| `01_core/usecase/evidence.go` | `2659cf5db5fba30e5003064aa3371aca11ee3673d7665db4b035f6d1d9c56050` |
| `03_infra/local/client.go` | `29b8176ab30bc07ec9ea818e7a2a6608fa0409b8d0255fa5b12b6268581bb012` |

O runner verifica esses blobs no commit, não os caminhos do working tree que pertencem à
implementação futura. Assim, uma mudança candidata legítima não invalida a identidade do
baseline nem impede a execução pós-candidato.
