# Especificação de Forjas Federadas de Código Aberto

**Camada**: L3
**Identificadores**:
- `gitee-repository-catalog` (v5 API - China)
- `codeberg-repository-catalog` (Forgejo/Gitea API - FOSS UE)
- `gitlab-repository-catalog` (GitLab API v4)
**Responsabilidade**:
- Buscar repositórios públicos nas forjas federadas regionais e soberanas.
- Mapear para tipos canônicos de `ProductCandidate` com nível de identidade `owner_repo` e `commit`.
