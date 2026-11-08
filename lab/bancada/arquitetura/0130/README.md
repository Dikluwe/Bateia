# Regras de alternativas e retomada — 0130

Contrato anterior à materialização: `00_nucleo/passos/0130-extrair-regras-de-alternativas-e-retomada.md` e ADR-0013. Recorte iniciado sobre a worktree existente, preservando as alterações anteriores 0126–0129 e a consolidação das instruções.

## Propriedade das decisões

- `01_core/usecase/alternative_target.go`: regras de alvo local, classificação Stale quando hints faltam e remapeamento consultivo do candidato/assessment/fontes. Leitor restrito conserva o lookup congelado, a ordem e a quantidade das consultas. Não autentica origem.
- `01_core/usecase/screening_model_policy.go`: admissão do perfil opt-in. Sem alteração de ScreenDiscovery, defaults, calibradores ou perfis calibrados. JSON estrito, leitura limitada, hash dos bytes e mensagens sanitizadas continuam na CLI.
- `01_core/usecase/evidence_checkpoint.go`: causalidade e elegibilidade de checkpoint; controla os dois caminhos lógicos permitidos, ordem, duplicatas e short-circuit de leitura. O leitor retorna somente SHA-256; o adaptador CLI resolve a raiz física e utiliza o host existente.

Não se alteram JSON público, RunID, regras de ranking, versões ou critérios de retomada. Mudança somente de política/alvo pode reaproveitar evidência compatível e recalcular a decisão. O novo leitor de checkpoint ausente recusa a consulta; os consumidores antigos sempre fornecem o adaptador.

## Evidência

`baseline.json` congela HEAD e hashes anteriores. `cli-before.log` e `cli-after.log` exercitam os mesmos contratos ativos de recomendação/retomada/alvo/modelo/captura/versão. O primeiro ensaio herdou assinatura SSH de commits de fixture e falhou antes do cenário de produto; `cli-before-environment.log` conserva essa falha. A execução ajustada desativa `commit.gpgsign` somente no processo, acrescentando uma entrada a GIT_CONFIG_COUNT. Nenhuma configuração global foi alterada.

Os seis fontes listados em `preserved-sources.json` mantêm os mesmos hashes anteriores: os três consumidores CLI e as implementações de recommendation, ranking e screening. Os testes e recibos selados não foram reescritos. `lineage.json` registra prompts e SHA-256 das sete materializações do recorte.

`core-policy-tests.log` cobre: hints/múltiplas worktrees e preservação das entradas/fontes; perfil opt-in e números não finitos; estados, causalidade, oito versões, presença do modelo, mudanças de policy/target, artefatos obrigatórios, duplicatas, leitura inesperada/erro/hash alterado, ordem e short-circuit. São testes de regra com leitores em memória, sem relógio ou filesystem.

`checks.json` e os logs correspondentes registram lint, auditor AST, suíte raiz, vet, build e diff-check. Lint usa o parser Go validado e reproduzível de 0129 no PATH do processo; a instalação global permanece intacta. A auditoria produtiva agora inclui os três novos arquivos L1. Build em `/tmp/bateia-0130` substitui `/dev/shm` neste ambiente. Os gates selados possuem outro módulo e não são promovidos pela suíte raiz.

## Reprodução

Na raiz, conforme o guia de engenharia:

```bash
bash lab/bancada/arquitetura/0129/reproduce-lint.sh
go test ./01_core/usecase -run 'Test(AlternativeTarget|OptInScreeningProfile|EvidenceCheckpoint)'
go test ./02_shell/cli -run 'Test(IdeaRecommend|Resume|LocalCaptureFrozenGitHints|TreeHashCheckpointVersionRequired)'
go test -v ./lab/bancada/arquitetura/0129/boundaries
go test ./...
go vet ./...
go build -o /tmp/bateia-0130 ./04_wiring
```

Fixtures com Git podem exigir desativar assinatura apenas no processo; testes com httptest exigem loopback. Pode-se usar GOCACHE temporário permitido. A reprodução do linter depende da revisão/patch e das ferramentas declaradas na bancada 0129.

## Limites e continuidade

A extração não encerra a jornada do dossiê nem move todas as regras da CLI. Construção de identidade/manifesto e composição condicional da recomendação ainda estão em L2. Veredicto experimental web e investigação/runtime de mercado também permanecem fora deste recorte. Aprovação dos checks não comprova a utilidade de hipóteses nem conformidade semântica universal com Tekt.
