# Fronteiras arquiteturais — 0129

Contrato: `00_nucleo/passos/0129-concluir-fronteiras-arquiteturais.md` e ADR-0012. Evidência desta mudança local, sobre a worktree que já continha 0126–0128. Não é promoção dos gates selados nem prova da utilidade de hipóteses.

## Responsabilidades verificadas

- L1: portas temporais/arquivos, regras e estado; containers de capacidades separados do lifecycle; parser do alvo e hints Git puros.
- L2: entrada, JSON, projeção, filtros de segredo e comunicação por portas; stdio na fachada do app. Sem driver físico de snapshots ou chamadas concretas a relógio/disco.
- L3: relógio/deadline/sinais/ambiente, armazenamento, exports, trace, hash, metadata e captura local. Retém descritores, limites, modos, identidade, proteção de aliases e cleanup.
- L4: injeção de host/runtime/provider, composição e lifecycle. Interpretação do limite MCP delegada a L2. Testes de protocolo dos conectores estão em L3; composição permanece em L4.
- Lab: auditor AST e controles adversariais; nenhuma dependência produtiva do Lab. L0 sem material executável no escopo inspecionado.

`baseline.json` registra o HEAD e os hashes anteriores ao corte. Os logs por capacidade registram a execução durante a migração; logs finais e `verification.json` determinam o resultado final. `go-test-initial.log` preserva a divergência de ordem física de `Examined` encontrada pela suíte. A correção ordena a projeção sem trocar a varredura limitada em lotes; o oráculo de paridade local/remoto permanece intacto. O hash certificado de `recommendation.go` foi preservado mediante extensão aditiva em `recommendation_runtime.go`.

## Lint reproduzível

```bash
bash lab/bancada/arquitetura/0129/reproduce-lint.sh
```

O script extrai a revisão `e53de6fd2bff9f3a87a3e8a076f5e470ce4bb7a0` do tekt-linter para `/tmp`, aplica `crystalline-go-parser.patch`, executa testes da biblioteca da ferramenta e 11 controles no binário real e invoca `crystalline-lint --fail-on warning .` com o binário validado no PATH do processo. Não altera checkout externo, regras, severidades ou instalação global. Necessita Rust/Cargo e dependências em cache; `BATEIA_LINTER_SOURCE` pode indicar outro clone com a revisão disponível.

O patch corrige fatos sintáticos: funções Go livres não são impls; métodos produtivos (inclusive privados) continuam fiscalizados; interfaces e aliases são classificados pelos tipos reais; `_test.go` recebe origem de teste, como cfg(test) em Rust. `check_test_imports=true` continua detectando inversões em testes. Os controles exigem rejeição de métodos, interfaces, aliases e imports produtivos ilegais. A instalação global anterior conserva os defeitos registrados em `lint-initial.log`; a aprovação final é da ferramenta reproduzida, não uma alegação de que o binário global foi atualizado.

## Auditor de fronteiras

```bash
go test -v ./lab/bancada/arquitetura/0129/boundaries
```

346 arquivos Go produtivos inspecionados por AST, com controles para aliases de imports, referência indireta a relógio, deadlines, disco, gravidade, isolamento do Lab e tipos de contrato em L4. A auditoria também recusa fontes, executáveis e fontes mascaradas em L0. Testes Go não integram o grafo produtivo; `cli_test` compõe explicitamente o mesmo host L3 para os testes da CLI. Um auditor sintático não demonstra toda a semântica de negócio, pureza referencial, autenticação de origens nem segurança universal contra filesystem hostil.

## Compatibilidade e limites

O executável instala o host antes de despachar comandos, inclusive offline. Bibliotecas CLI precisam de `ConfigureCommandHost` no startup. `WithExecutionRuntime` configura fachadas L1; runtime ausente não inventa data/duração e não satisfaz deadline solicitado. O seam global da CLI não deve ser mutado durante comandos concorrentes.

A suíte raiz usa loopback; assinatura Git é desativada apenas no ambiente do processo de teste, acrescentando uma entrada à configuração existente. Build em `/tmp/bateia-0129` substitui `/dev/shm` neste ambiente. Gates selados possuem módulo próprio e não foram promovidos por esta execução.

A revisão automática rejeitou uma reescrita ampla por regex durante a preparação. A migração foi feita em extrações por capacidade, preservando corpos de segurança e executando seus testes. Também rejeitou substituir a entrada protegida `.agents/AGENTS.md` por uma referência canônica; esse arquivo permanece intacto. Não se concede nova política ou exceção estrutural com este registro.
