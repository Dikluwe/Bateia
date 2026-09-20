# Protocolo de execução — Passo 0011

**Regime**: protocolo Tekt completo executado sem atestação de isolamento  
**Baseline declarado**: `df14ec138261094a9fd33c2e887ef83db228b4a4`  
**Data**: 2026-08-26

## Entradas congeladas nesta execução

| Artefato | SHA-256 |
|---|---|
| `00_nucleo/passos/0011-congelar-benchmark-ponta-a-ponta-e-auditar-utilidade.md` | `7746d03fa20b38067fa43ce5125c9484e356dfc2e21683e1eba37098d3e0985b` |
| `01_core/domain/end_to_end_benchmark.go` | `06f84e7f59278522d2b757900ce777c070e538cf041f07cf4e4c8b198aa0d6b8` |
| `01_core/usecase/end_to_end_benchmark.go` | `cc17e3b419955f1f236e4abdfc04a3ecce26eb39cdeef9f7cb385c0a0ada6495` |
| `02_shell/cli/benchmark.go` | `48e207fc8eb9a70f59b46ce78d634f777bd1a050aa6c46c173d97daf5d223162` |

O baseline é uma referência imutável do repositório, mas a worktree já continha
mudanças dos passos anteriores. Estes hashes registram a execução observada; não
alegam que contrato, oráculos, implementação e veredito foram produzidos por
autoridades tecnicamente isoladas.

## Política operacional de `UNKNOWN`

Um case é `UNKNOWN` quando a identidade causal não pode ser provada, inclusive
commit `HEAD`, regra sem versão/hash, proveniência incompleta ou bundle opaco.
Ele não entra como acerto em precisão e precisa de recaptura. Divergência
comprovada entre ideia, candidato, commit, conteúdo ou manifesto é `VIOLATED`.

## Recibos de papéis

| Papel | Entrada | Saída | Limitação |
|---|---|---|---|
| obrigação | desejo, passos 0009/0010 | passo 0011 | mesma sessão |
| contrato | passo 0011 | tipos e correspondência L1 | pôde observar a worktree |
| adversário | contrato | ataques em testes | mesma autoridade do implementador |
| implementador | contrato e ataques | runner e CLI | sem isolamento técnico |
| verificador | worktree materializada | testes, linter e laudo | não é verificador independente |

Não há recibo de autor de oráculos humanos: nenhum corpus `TEST` foi selado. Essa
ausência impede promover o resultado a certificação de utilidade.

## Gate de parada aplicado

Não foi fabricado um corpus de quinze casos nem um split `TEST` a partir de
expectativas escolhidas pela própria implementação. O mecanismo está pronto para
receber oráculos correspondentes; a avaliação de produto permanece aberta até
que esses oráculos sejam materializados e revisados independentemente.
