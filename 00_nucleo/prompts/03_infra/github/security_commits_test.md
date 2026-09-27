# Prompt: revisão imutável resolvida pelo GitHub

Camada L3. Gerar `03_infra/github/security_commits_test.go`.

Manter HEAD mutável no legado. A API de commits só pode resolver HEAD para SHA hexadecimal completo de 40 caracteres; branch, abreviação, path traversal e query injetada devem retornar erro sem promover CommitID. Validar branch padrão antes de interpolar em requisição. Testes offline com transport sintético, sem rede/credenciais reais.
